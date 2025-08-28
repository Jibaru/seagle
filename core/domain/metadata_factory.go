package domain

import (
	"fmt"
)

type MetadataFactory struct {
	serviceFactory *ServiceFactory
}

func NewMetadataFactory(serviceFactory *ServiceFactory) *MetadataFactory {
	return &MetadataFactory{serviceFactory: serviceFactory}
}

func (s *MetadataFactory) NewConnectionMetadata(conn *Connection) (*ConnectionMetadata, error) {
	dbService, err := s.serviceFactory.NewDatabaseService(conn)
	if err != nil {
		return nil, fmt.Errorf("failed to create database service: %w", err)
	}

	if err := dbService.Connect(conn); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	defer dbService.Disconnect(conn)

	metadata := NewConnectionMetadata(conn.ID(), conn.Host(), conn.Port())

	// Get all databases
	databases, err := dbService.GetDatabaseNames(conn)
	if err != nil {
		return nil, fmt.Errorf("failed to get database list: %w", err)
	}

	// Analyze each database
	for _, dbName := range databases {
		dbMetadata, err := s.analyzeDatabaseMetadata(conn, dbName)
		if err != nil {
			// Log error but continue with other databases
			continue
		}
		metadata.AddDatabase(dbMetadata)
	}

	return metadata, nil
}

// analyzeDatabaseMetadata analyzes a specific database and returns its metadata
func (s *MetadataFactory) analyzeDatabaseMetadata(conn *Connection, databaseName string) (*DatabaseMetadata, error) {
	// Create a connection copy for the specific database
	cpy := CopyConnection(conn, databaseName)
	dbService, err := s.serviceFactory.NewDatabaseService(cpy)
	if err != nil {
		return nil, fmt.Errorf("failed to create database service: %w", err)
	}

	if err := dbService.Connect(cpy); err != nil {
		return nil, fmt.Errorf("failed to connect to database %s: %w", databaseName, err)
	}
	defer dbService.Disconnect(cpy)

	metadata := NewDatabaseMetadata(databaseName)

	// Get all tables in the database
	tables, err := s.getTableList(cpy, dbService)
	if err != nil {
		return nil, fmt.Errorf("failed to get table list for database %s: %w", databaseName, err)
	}

	// Analyze each table
	for _, tableInfo := range tables {
		tableMetadata, err := dbService.GetTableMetadata(cpy, tableInfo.Name, tableInfo.Schema)
		if err != nil {
			// Log error but continue with other tables
			continue
		}
		metadata.AddTable(tableMetadata)
	}

	return metadata, nil
}

// TableInfo holds basic table information
type TableInfo struct {
	Name   string
	Schema string
}

// getTableList retrieves all tables in a database
func (s *MetadataFactory) getTableList(conn *Connection, dbService DatabaseService) ([]TableInfo, error) {
	var query string
	switch conn.Vendor() {
	case "postgresql":
		query = `
			SELECT table_name, table_schema
			FROM information_schema.tables 
			WHERE table_schema = 'public' 
			AND table_type = 'BASE TABLE'
			ORDER BY table_name
		`
	case "mysql":
		query = `
			SELECT table_name, table_schema
			FROM information_schema.tables 
			WHERE table_schema = DATABASE()
			AND table_type = 'BASE TABLE'
			ORDER BY table_name
		`
	default:
		return nil, fmt.Errorf("unsupported vendor: %s", conn.Vendor())
	}

	data, err := dbService.ExecQuery(conn, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query tables: %w", err)
	}

	var tables []TableInfo
	for _, row := range data.Rows {
		tableName, ok := row[0].(string)
		if !ok {
			return nil, fmt.Errorf("failed to convert table name to string")
		}

		schemaName, ok := row[1].(string)
		if !ok {
			return nil, fmt.Errorf("failed to convert schema name to string")
		}

		tables = append(tables, TableInfo{Name: tableName, Schema: schemaName})
	}

	return tables, nil
}
