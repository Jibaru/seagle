package persistence

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"seagle/core/domain"
)

// MetadataRepository implements the domain.MetadataRepo interface using JSON file storage
type MetadataRepository struct {
	filePath string
}

// NewMetadataRepository creates a new MetadataRepository instance
func NewMetadataRepository(fileName string) *MetadataRepository {
	return &MetadataRepository{
		filePath: fileName,
	}
}

// metadataFile represents the JSON structure for metadata storage
type metadataFile struct {
	Metadata []metadataRecord `json:"metadata"`
}

// metadataRecord represents a single connection metadata record in JSON
type metadataRecord struct {
	ConnectionID string           `json:"connectionId"`
	Host         string           `json:"host"`
	Port         int              `json:"port"`
	Databases    []databaseRecord `json:"databases"`
}

type databaseRecord struct {
	Name   string        `json:"name"`
	Tables []tableRecord `json:"tables"`
}

type tableRecord struct {
	Name    string         `json:"name"`
	Schema  string         `json:"schema"`
	Columns []columnRecord `json:"columns"`
}

type columnRecord struct {
	Name         string `json:"name"`
	DataType     string `json:"dataType"`
	IsNullable   bool   `json:"isNullable"`
	DefaultValue string `json:"defaultValue,omitempty"`
	Position     int    `json:"position"`
}

// Save persists the connection metadata
func (r *MetadataRepository) Save(metadata *domain.ConnectionMetadata) error {
	file, err := r.loadFile()
	if err != nil {
		return fmt.Errorf("failed to load metadata file: %w", err)
	}

	// Convert domain metadata to persistence record
	record := r.domainToRecord(metadata)

	// Remove existing metadata for this connection if it exists
	for i, existing := range file.Metadata {
		if existing.ConnectionID == metadata.ConnectionID() {
			file.Metadata = append(file.Metadata[:i], file.Metadata[i+1:]...)
			break
		}
	}

	// Add the new metadata
	file.Metadata = append(file.Metadata, record)

	return r.saveFile(file)
}

// FindByConnectionID retrieves metadata for a specific connection
func (r *MetadataRepository) FindByConnectionID(connectionID string) (*domain.ConnectionMetadata, error) {
	file, err := r.loadFile()
	if err != nil {
		return nil, fmt.Errorf("failed to load metadata file: %w", err)
	}

	for _, record := range file.Metadata {
		if record.ConnectionID == connectionID {
			return r.recordToDomain(record), nil
		}
	}

	return nil, nil // Not found
}

// Exists checks if metadata exists for a connection
func (r *MetadataRepository) Exists(connectionID string) (bool, error) {
	metadata, err := r.FindByConnectionID(connectionID)
	if err != nil {
		return false, err
	}
	return metadata != nil, nil
}

// Delete removes metadata for a connection
func (r *MetadataRepository) Delete(connectionID string) error {
	file, err := r.loadFile()
	if err != nil {
		return fmt.Errorf("failed to load metadata file: %w", err)
	}

	for i, record := range file.Metadata {
		if record.ConnectionID == connectionID {
			file.Metadata = append(file.Metadata[:i], file.Metadata[i+1:]...)
			return r.saveFile(file)
		}
	}

	return nil // Not found, but not an error
}

// List returns all stored connection metadata
func (r *MetadataRepository) List() ([]*domain.ConnectionMetadata, error) {
	file, err := r.loadFile()
	if err != nil {
		return nil, fmt.Errorf("failed to load metadata file: %w", err)
	}

	result := make([]*domain.ConnectionMetadata, len(file.Metadata))
	for i, record := range file.Metadata {
		result[i] = r.recordToDomain(record)
	}

	return result, nil
}

// loadFile loads the metadata file or creates an empty one if it doesn't exist
func (r *MetadataRepository) loadFile() (*metadataFile, error) {
	if _, err := os.Stat(r.filePath); os.IsNotExist(err) {
		return &metadataFile{Metadata: []metadataRecord{}}, nil
	}

	data, err := os.ReadFile(r.filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read metadata file: %w", err)
	}

	var file metadataFile
	if err := json.Unmarshal(data, &file); err != nil {
		return nil, fmt.Errorf("failed to unmarshal metadata file: %w", err)
	}

	return &file, nil
}

// saveFile saves the metadata file
func (r *MetadataRepository) saveFile(file *metadataFile) error {
	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(r.filePath), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	data, err := json.MarshalIndent(file, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	if err := os.WriteFile(r.filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write metadata file: %w", err)
	}

	return nil
}

// domainToRecord converts domain metadata to persistence record
func (r *MetadataRepository) domainToRecord(metadata *domain.ConnectionMetadata) metadataRecord {
	record := metadataRecord{
		ConnectionID: metadata.ConnectionID(),
		Host:         metadata.Host(),
		Port:         metadata.Port(),
		Databases:    make([]databaseRecord, len(metadata.Databases())),
	}

	for i, db := range metadata.Databases() {
		record.Databases[i] = databaseRecord{
			Name:   db.Name(),
			Tables: make([]tableRecord, len(db.Tables())),
		}

		for j, table := range db.Tables() {
			record.Databases[i].Tables[j] = tableRecord{
				Name:    table.Name(),
				Schema:  table.Schema(),
				Columns: make([]columnRecord, len(table.Columns())),
			}

			for k, col := range table.Columns() {
				record.Databases[i].Tables[j].Columns[k] = columnRecord{
					Name:         col.Name(),
					DataType:     col.DataType(),
					IsNullable:   col.IsNullable(),
					DefaultValue: col.DefaultValue(),
					Position:     col.Position(),
				}
			}
		}
	}

	return record
}

// recordToDomain converts persistence record to domain metadata
func (r *MetadataRepository) recordToDomain(record metadataRecord) *domain.ConnectionMetadata {
	metadata := domain.NewConnectionMetadata(record.ConnectionID, record.Host, record.Port)

	for _, dbRecord := range record.Databases {
		dbMetadata := domain.NewDatabaseMetadata(dbRecord.Name)

		for _, tableRecord := range dbRecord.Tables {
			tableMetadata := domain.NewTableMetadata(tableRecord.Name, tableRecord.Schema)

			for _, colRecord := range tableRecord.Columns {
				columnMetadata := domain.NewColumnMetadata(
					colRecord.Name,
					colRecord.DataType,
					colRecord.IsNullable,
					colRecord.DefaultValue,
					colRecord.Position,
				)
				tableMetadata.AddColumn(columnMetadata)
			}

			dbMetadata.AddTable(tableMetadata)
		}

		metadata.AddDatabase(dbMetadata)
	}

	return metadata
}
