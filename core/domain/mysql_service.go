package domain

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type MySQLService struct {
	pool map[string]*sql.DB
}

func NewMySQLService() *MySQLService {
	return &MySQLService{
		pool: make(map[string]*sql.DB),
	}
}

func (s *MySQLService) pooledDBConn(c *Connection) *sql.DB {
	if db, exists := s.pool[c.ID()]; exists {
		return db
	}
	return nil
}

func (s *MySQLService) Connect(c *Connection) error {
	dbConn := s.pooledDBConn(c)
	if dbConn == nil {
		db, err := sql.Open("mysql", s.buildConnectionString(c))
		if err != nil {
			return fmt.Errorf("failed to open database connection: %w", err)
		}
		dbConn = db
	}

	if err := dbConn.Ping(); err != nil {
		dbConn.Close()
		return fmt.Errorf("failed to ping database: %w", err)
	}

	s.pool[c.ID()] = dbConn
	return nil
}

func (s *MySQLService) Disconnect(c *Connection) error {
	dbConn := s.pooledDBConn(c)

	if dbConn != nil {
		if err := dbConn.Close(); err != nil {
			return fmt.Errorf("failed to close database connection: %w", err)
		}
		delete(s.pool, c.ID())
	}
	return nil
}

func (s *MySQLService) buildConnectionString(c *Connection) string {
	// MySQL connection string format: user:password@tcp(host:port)/database?params
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", c.username, c.password, c.Host(), c.Port(), c.database)
	
	if len(c.arguments) > 0 {
		connStr += "?"
		first := true
		for k, v := range c.arguments {
			if !first {
				connStr += "&"
			}
			connStr += fmt.Sprintf("%s=%s", k, v)
			first = false
		}
	}
	
	return connStr
}

func (s *MySQLService) GetDatabaseNames(c *Connection) ([]string, error) {
	query := `SHOW DATABASES`

	dbConn := s.pooledDBConn(c)
	if dbConn == nil {
		return nil, fmt.Errorf("no active connection found")
	}

	rows, err := dbConn.Query(query)
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
		// Filter out MySQL system databases
		if dbName != "information_schema" && dbName != "mysql" && dbName != "performance_schema" && dbName != "sys" {
			databases = append(databases, dbName)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating database results: %w", err)
	}

	return databases, nil
}

func (s *MySQLService) GetTableNames(c *Connection, databaseName string) ([]string, error) {
	query := `
		SELECT table_name 
		FROM information_schema.tables 
		WHERE table_schema = ? 
		AND table_type = 'BASE TABLE'
		ORDER BY table_name
	`

	dbConn := s.pooledDBConn(c)
	if dbConn == nil {
		return nil, fmt.Errorf("no active connection found")
	}

	rows, err := dbConn.Query(query, databaseName)
	if err != nil {
		return nil, fmt.Errorf("failed to query tables for database %s: %w", databaseName, err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return nil, fmt.Errorf("failed to scan table name: %w", err)
		}
		tables = append(tables, tableName)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating table results: %w", err)
	}

	return tables, nil
}

func (s *MySQLService) GetTableColumns(c *Connection, databaseName, tableName string) ([]ColumnMetadata, error) {
	query := `
		SELECT column_name, data_type, is_nullable, column_default, ordinal_position
		FROM information_schema.columns
		WHERE table_schema = ? AND table_name = ?
		ORDER BY ordinal_position		
	`
	dbConn := s.pooledDBConn(c)
	if dbConn == nil {
		return nil, fmt.Errorf("no active connection found")
	}

	rows, err := dbConn.Query(query, databaseName, tableName)
	if err != nil {
		return nil, fmt.Errorf("failed to query columns for table %s: %w", tableName, err)
	}
	defer rows.Close()

	type Col struct {
		Name         string
		DataType     string
		IsNullable   string
		DefaultValue sql.NullString
		Position     int
	}

	var columns []ColumnMetadata
	for rows.Next() {
		var col Col
		if err := rows.Scan(&col.Name, &col.DataType, &col.IsNullable, &col.DefaultValue, &col.Position); err != nil {
			return nil, fmt.Errorf("failed to scan column metadata: %w", err)
		}

		column := ColumnMetadata{
			name:         col.Name,
			dataType:     col.DataType,
			isNullable:   col.IsNullable == "YES",
			defaultValue: col.DefaultValue.String,
			position:     col.Position,
		}
		columns = append(columns, column)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating column results: %w", err)
	}

	return columns, nil
}

func (s *MySQLService) ExecQuery(c *Connection, query string) (*QueryResult, error) {
	dbConn := s.pooledDBConn(c)
	if dbConn == nil {
		return nil, fmt.Errorf("no active connection found")
	}

	start := time.Now()

	// Check if it's a SELECT query or DML/DDL
	rows, err := dbConn.Query(query)
	if err != nil {
		// If Query fails, try Exec for DML/DDL statements
		result, execErr := dbConn.Exec(query)
		if execErr != nil {
			return nil, fmt.Errorf("failed to execute query: %w", execErr)
		}

		rowsAffected, _ := result.RowsAffected()
		duration := time.Since(start).Milliseconds()

		return &QueryResult{
			Columns:      []string{},
			Rows:         [][]interface{}{},
			RowsAffected: rowsAffected,
			Duration:     duration,
		}, nil
	}
	defer rows.Close()

	// Get column information
	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to get columns: %w", err)
	}

	// Prepare result structure
	var resultRows [][]interface{}

	for rows.Next() {
		// Create a slice of interface{} to hold the values
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))

		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		// Convert []byte to string for display
		row := make([]interface{}, len(values))
		for i, val := range values {
			if b, ok := val.([]byte); ok {
				row[i] = string(b)
			} else {
				row[i] = val
			}
		}

		resultRows = append(resultRows, row)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	duration := time.Since(start).Milliseconds()

	return &QueryResult{
		Columns:      columns,
		Rows:         resultRows,
		RowsAffected: int64(len(resultRows)),
		Duration:     duration,
	}, nil
}

func (s *MySQLService) GetTableMetadata(c *Connection, tableName, schemaName string) (*TableMetadata, error) {
	db := s.pooledDBConn(c)
	if db == nil {
		return nil, fmt.Errorf("no active connection found")
	}

	metadata := NewTableMetadata(tableName, schemaName)

	query := `
		SELECT 
			column_name,
			data_type,
			is_nullable = 'YES' as is_nullable,
			COALESCE(column_default, '') as column_default,
			ordinal_position
		FROM information_schema.columns 
		WHERE table_schema = ? 
		AND table_name = ?
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