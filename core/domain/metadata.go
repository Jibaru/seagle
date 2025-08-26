package domain

import (
	"database/sql"
	"fmt"
)

// ColumnMetadata represents metadata for a single column
type ColumnMetadata struct {
	name         string
	dataType     string
	isNullable   bool
	defaultValue string
	position     int
}

// NewColumnMetadata creates a new ColumnMetadata instance
func NewColumnMetadata(name, dataType string, isNullable bool, defaultValue string, position int) *ColumnMetadata {
	return &ColumnMetadata{
		name:         name,
		dataType:     dataType,
		isNullable:   isNullable,
		defaultValue: defaultValue,
		position:     position,
	}
}

// Name returns the column name
func (c *ColumnMetadata) Name() string {
	return c.name
}

// DataType returns the column data type
func (c *ColumnMetadata) DataType() string {
	return c.dataType
}

// IsNullable returns whether the column is nullable
func (c *ColumnMetadata) IsNullable() bool {
	return c.isNullable
}

// DefaultValue returns the column default value
func (c *ColumnMetadata) DefaultValue() string {
	return c.defaultValue
}

// Position returns the column position
func (c *ColumnMetadata) Position() int {
	return c.position
}

// TableMetadata represents metadata for a single table
type TableMetadata struct {
	name    string
	schema  string
	columns []*ColumnMetadata
}

// NewTableMetadata creates a new TableMetadata instance
func NewTableMetadata(name, schema string) *TableMetadata {
	return &TableMetadata{
		name:    name,
		schema:  schema,
		columns: make([]*ColumnMetadata, 0),
	}
}

// Name returns the table name
func (t *TableMetadata) Name() string {
	return t.name
}

// Schema returns the table schema
func (t *TableMetadata) Schema() string {
	return t.schema
}

// Columns returns the table columns
func (t *TableMetadata) Columns() []*ColumnMetadata {
	return t.columns
}

// AddColumn adds a column to the table metadata
func (t *TableMetadata) AddColumn(column *ColumnMetadata) {
	t.columns = append(t.columns, column)
}

// DatabaseMetadata represents the complete metadata structure for a database
type DatabaseMetadata struct {
	name   string
	tables []*TableMetadata
}

// NewDatabaseMetadata creates a new DatabaseMetadata instance
func NewDatabaseMetadata(name string) *DatabaseMetadata {
	return &DatabaseMetadata{
		name:   name,
		tables: make([]*TableMetadata, 0),
	}
}

// Name returns the database name
func (d *DatabaseMetadata) Name() string {
	return d.name
}

// Tables returns the database tables
func (d *DatabaseMetadata) Tables() []*TableMetadata {
	return d.tables
}

// AddTable adds a table to the database metadata
func (d *DatabaseMetadata) AddTable(table *TableMetadata) {
	d.tables = append(d.tables, table)
}

// ConnectionMetadata represents the complete metadata for a connection
type ConnectionMetadata struct {
	connectionID string
	host         string
	port         int
	databases    []*DatabaseMetadata
}

// NewConnectionMetadata creates a new ConnectionMetadata instance
func NewConnectionMetadata(connectionID, host string, port int) *ConnectionMetadata {
	return &ConnectionMetadata{
		connectionID: connectionID,
		host:         host,
		port:         port,
		databases:    make([]*DatabaseMetadata, 0),
	}
}

// ConnectionID returns the connection ID
func (c *ConnectionMetadata) ConnectionID() string {
	return c.connectionID
}

// Host returns the connection host
func (c *ConnectionMetadata) Host() string {
	return c.host
}

// Port returns the connection port
func (c *ConnectionMetadata) Port() int {
	return c.port
}

// Databases returns the connection databases
func (c *ConnectionMetadata) Databases() []*DatabaseMetadata {
	return c.databases
}

// AddDatabase adds a database to the connection metadata
func (c *ConnectionMetadata) AddDatabase(database *DatabaseMetadata) {
	c.databases = append(c.databases, database)
}

// AnalyzeConnectionMetadata analyzes a connection and returns complete metadata
func AnalyzeConnectionMetadata(conn *Connection) (*ConnectionMetadata, error) {
	if !conn.IsConnected() {
		return nil, fmt.Errorf("connection is not active")
	}

	metadata := NewConnectionMetadata(conn.ID(), conn.Host(), conn.Port())

	// Get all databases
	databases, err := getDatabaseList(conn.DB())
	if err != nil {
		return nil, fmt.Errorf("failed to get database list: %w", err)
	}

	// Analyze each database
	for _, dbName := range databases {
		dbMetadata, err := analyzeDatabaseMetadata(conn, dbName)
		if err != nil {
			// Log error but continue with other databases
			continue
		}
		metadata.AddDatabase(dbMetadata)
	}

	return metadata, nil
}

// getDatabaseList retrieves all non-template databases
func getDatabaseList(db *sql.DB) ([]string, error) {
	query := `
		SELECT datname 
		FROM pg_database 
		WHERE datistemplate = false
		ORDER BY datname
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query databases: %w", err)
	}
	defer rows.Close()

	var databases []string
	for rows.Next() {
		var dbName string
		if err := rows.Scan(&dbName); err != nil {
			return nil, fmt.Errorf("failed to scan database name: %w", err)
		}
		databases = append(databases, dbName)
	}

	return databases, rows.Err()
}

// analyzeDatabaseMetadata analyzes a specific database and returns its metadata
func analyzeDatabaseMetadata(conn *Connection, databaseName string) (*DatabaseMetadata, error) {
	// Create a connection copy for the specific database
	dbConn := CopyConnection(conn, databaseName)
	if err := dbConn.Connect(); err != nil {
		return nil, fmt.Errorf("failed to connect to database %s: %w", databaseName, err)
	}
	defer dbConn.Disconnect()

	metadata := NewDatabaseMetadata(databaseName)

	// Get all tables in the database
	tables, err := getTableList(dbConn.DB())
	if err != nil {
		return nil, fmt.Errorf("failed to get table list for database %s: %w", databaseName, err)
	}

	// Analyze each table
	for _, tableInfo := range tables {
		tableMetadata, err := analyzeTableMetadata(dbConn.DB(), tableInfo.Name, tableInfo.Schema)
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
func getTableList(db *sql.DB) ([]TableInfo, error) {
	query := `
		SELECT table_name, table_schema
		FROM information_schema.tables 
		WHERE table_schema = 'public' 
		AND table_type = 'BASE TABLE'
		ORDER BY table_name
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query tables: %w", err)
	}
	defer rows.Close()

	var tables []TableInfo
	for rows.Next() {
		var table TableInfo
		if err := rows.Scan(&table.Name, &table.Schema); err != nil {
			return nil, fmt.Errorf("failed to scan table info: %w", err)
		}
		tables = append(tables, table)
	}

	return tables, rows.Err()
}

// analyzeTableMetadata analyzes a specific table and returns its metadata
func analyzeTableMetadata(db *sql.DB, tableName, schemaName string) (*TableMetadata, error) {
	metadata := NewTableMetadata(tableName, schemaName)

	query := `
		SELECT 
			column_name,
			data_type,
			is_nullable = 'YES' as is_nullable,
			COALESCE(column_default, '') as column_default,
			ordinal_position
		FROM information_schema.columns 
		WHERE table_schema = $1 
		AND table_name = $2
		ORDER BY ordinal_position
	`

	rows, err := db.Query(query, schemaName, tableName)
	if err != nil {
		return nil, fmt.Errorf("failed to query columns for table %s.%s: %w", schemaName, tableName, err)
	}
	defer rows.Close()

	for rows.Next() {
		var name, dataType, defaultValue string
		var isNullable bool
		var position int

		if err := rows.Scan(&name, &dataType, &isNullable, &defaultValue, &position); err != nil {
			return nil, fmt.Errorf("failed to scan column metadata: %w", err)
		}

		column := NewColumnMetadata(name, dataType, isNullable, defaultValue, position)
		metadata.AddColumn(column)
	}

	return metadata, rows.Err()
}
