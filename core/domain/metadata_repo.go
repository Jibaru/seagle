package domain

// MetadataRepo defines the interface for metadata persistence operations
type MetadataRepo interface {
	// Save persists the connection metadata
	Save(metadata *ConnectionMetadata) error

	// FindByConnectionID retrieves metadata for a specific connection
	FindByConnectionID(connectionID string) (*ConnectionMetadata, error)

	// Exists checks if metadata exists for a connection
	Exists(connectionID string) (bool, error)

	// Delete removes metadata for a connection
	Delete(connectionID string) error

	// List returns all stored connection metadata
	List() ([]*ConnectionMetadata, error)
}
