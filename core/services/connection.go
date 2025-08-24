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
