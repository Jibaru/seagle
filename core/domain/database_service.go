package domain

// DatabaseService defines the interface for database operations
type DatabaseService interface {
	// Connection management
	Connect(c *Connection) error
	Disconnect(c *Connection) error

	// Database metadata operations
	GetDatabaseNames(c *Connection) ([]string, error)
	GetTableNames(c *Connection, databaseName string) ([]string, error)
	GetTableColumns(c *Connection, databaseName, tableName string) ([]ColumnMetadata, error)
	GetTableMetadata(c *Connection, tableName, schemaName string) (*TableMetadata, error)

	// Query execution
	ExecQuery(c *Connection, query string) (*QueryResult, error)
}

// QueryResult represents the result of a query execution
type QueryResult struct {
	Columns      []string
	Rows         [][]interface{}
	RowsAffected int64
	Duration     int64
}