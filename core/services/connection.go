package services

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"seagle/core/services/types"
)

// ConnectionService manages database connections
type ConnectionService struct {
	connection *types.DatabaseConnection
}

// NewConnectionService creates a new ConnectionService instance
func NewConnectionService() *ConnectionService {
	return &ConnectionService{}
}

// Connect establishes a connection to PostgreSQL database
func (cs *ConnectionService) Connect(config types.DatabaseConfig) (*types.DatabaseConnection, error) {
	connStr, err := cs.buildConnectionString(config)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	connection := &types.DatabaseConnection{
		Config:      config,
		IsConnected: true,
		DB:          db,
	}

	cs.connection = connection
	return connection, nil
}

// TestConnection tests the database connection with given parameters
func (cs *ConnectionService) TestConnection(config types.DatabaseConfig) error {
	connStr, err := cs.buildConnectionString(config)
	if err != nil {
		return err
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to open database connection: %w", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	return nil
}

// Disconnect closes the current database connection
func (cs *ConnectionService) Disconnect() error {
	if cs.connection != nil && cs.connection.DB != nil {
		if err := cs.connection.DB.Close(); err != nil {
			return fmt.Errorf("failed to close database connection: %w", err)
		}
		cs.connection.IsConnected = false
		cs.connection = nil
	}
	return nil
}

// buildConnectionString constructs the PostgreSQL connection string
func (cs *ConnectionService) buildConnectionString(config types.DatabaseConfig) (string, error) {
	if config.UseConnectionString {
		if config.ConnectionString == "" {
			return "", fmt.Errorf("connection string cannot be empty")
		}
		return config.ConnectionString, nil
	}

	// Validate required fields for form-based connection
	if config.Host == "" || config.Database == "" || config.Username == "" {
		return "", fmt.Errorf("host, database, and username are required")
	}

	sslMode := config.SSLMode
	if sslMode == "" {
		sslMode = "require"
	}

	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.Username, config.Password, config.Database, sslMode), nil
}

// GetDatabases returns a list of databases available in the connected PostgreSQL instance
func (cs *ConnectionService) GetDatabases() ([]string, error) {
	if cs.connection == nil || cs.connection.DB == nil {
		return nil, fmt.Errorf("no active database connection")
	}

	query := `
		SELECT datname 
		FROM pg_database 
		WHERE datistemplate = false
		ORDER BY datname
	`

	rows, err := cs.connection.DB.Query(query)
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

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating database results: %w", err)
	}

	return databases, nil
}

// GetTables returns a list of tables for a specific database
func (cs *ConnectionService) GetTables(databaseName string) ([]string, error) {
	if cs.connection == nil || cs.connection.DB == nil {
		return nil, fmt.Errorf("no active database connection")
	}

	// First, connect to the specific database
	config := cs.connection.Config
	config.Database = databaseName
	
	connStr, err := cs.buildConnectionString(config)
	if err != nil {
		return nil, fmt.Errorf("failed to build connection string: %w", err)
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database %s: %w", databaseName, err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database %s: %w", databaseName, err)
	}

	query := `
		SELECT table_name 
		FROM information_schema.tables 
		WHERE table_schema = 'public' 
		AND table_type = 'BASE TABLE'
		ORDER BY table_name
	`

	rows, err := db.Query(query)
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

// GetTableColumns returns the columns for a specific table in a database
func (cs *ConnectionService) GetTableColumns(databaseName, tableName string) ([]types.TableColumn, error) {
	if cs.connection == nil || cs.connection.DB == nil {
		return nil, fmt.Errorf("no active database connection")
	}

	// Connect to the specific database
	config := cs.connection.Config
	config.Database = databaseName
	
	connStr, err := cs.buildConnectionString(config)
	if err != nil {
		return nil, fmt.Errorf("failed to build connection string: %w", err)
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database %s: %w", databaseName, err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database %s: %w", databaseName, err)
	}

	query := `
		SELECT 
			column_name,
			data_type,
			is_nullable = 'YES' as is_nullable,
			COALESCE(column_default, '') as column_default
		FROM information_schema.columns 
		WHERE table_schema = 'public' 
		AND table_name = $1
		ORDER BY ordinal_position
	`

	rows, err := db.Query(query, tableName)
	if err != nil {
		return nil, fmt.Errorf("failed to query columns for table %s.%s: %w", databaseName, tableName, err)
	}
	defer rows.Close()

	var columns []types.TableColumn
	for rows.Next() {
		var col types.TableColumn
		if err := rows.Scan(&col.Name, &col.DataType, &col.IsNullable, &col.DefaultValue); err != nil {
			return nil, fmt.Errorf("failed to scan column: %w", err)
		}
		columns = append(columns, col)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating column results: %w", err)
	}

	return columns, nil
}
