package domain

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
