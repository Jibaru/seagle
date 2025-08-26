package services

import (
	"fmt"
	"time"

	_ "github.com/lib/pq"

	"seagle/core/domain"
	"seagle/core/services/types"
)

// ConnectionService manages database connections
type ConnectionService struct {
	currentConnection *domain.Connection
	repo              domain.ConnectionRepo
	metadataRepo      domain.MetadataRepo
}

// NewConnectionService creates a new ConnectionService instance
func NewConnectionService(repo domain.ConnectionRepo, metadataRepo domain.MetadataRepo) *ConnectionService {
	return &ConnectionService{
		repo:         repo,
		metadataRepo: metadataRepo,
	}
}

// Connect establishes a connection using DatabaseConfig
func (cs *ConnectionService) Connect(config types.DatabaseConfig) (*types.DatabaseConnection, error) {
	domainConn, err := cs.configToDomainConnection(cs.repo.NextID(), config)
	if err != nil {
		return nil, err
	}

	if err := cs.repo.Save(domainConn); err != nil {
		return nil, fmt.Errorf("failed to save connection: %w", err)
	}

	if err := domainConn.Connect(); err != nil {
		return nil, err
	}

	cs.currentConnection = domainConn

	return &types.DatabaseConnection{
		Config:      config,
		IsConnected: true,
	}, nil
}

func (cs *ConnectionService) ConnectByID(id string) (*types.DatabaseConnection, error) {
	conn, err := cs.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to find connection by ID: %w", err)
	}
	if conn == nil {
		return nil, fmt.Errorf("connection with ID %s not found", id)
	}

	if err := conn.Connect(); err != nil {
		return nil, err
	}

	cs.currentConnection = conn

	return &types.DatabaseConnection{
		IsConnected: true,
	}, nil
}

// TestConnection tests the database connection with given parameters
func (cs *ConnectionService) TestConnection(config types.DatabaseConfig) error {
	domainConn, err := cs.configToDomainConnection(cs.repo.NextID(), config)
	if err != nil {
		return err
	}

	if err := domainConn.Connect(); err != nil {
		return err
	}

	return nil
}

// Disconnect closes the current database connection
func (cs *ConnectionService) Disconnect() error {
	return cs.currentConnection.Disconnect()
}

// HasActiveConnection checks if there is an active database connection
func (cs *ConnectionService) HasActiveConnection() bool {
	return cs.currentConnection != nil && cs.currentConnection.IsConnected()
}

// GetDatabases returns a list of databases available in the connected PostgreSQL instance
func (cs *ConnectionService) GetDatabases() ([]string, error) {
	if !cs.HasActiveConnection() {
		return nil, fmt.Errorf("no active database connection")
	}

	query := `
		SELECT datname 
		FROM pg_database 
		WHERE datistemplate = false
		ORDER BY datname
	`

	rows, err := cs.currentConnection.DB().Query(query)
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
	if !cs.HasActiveConnection() {
		return nil, fmt.Errorf("no active database connection")
	}

	cpy := domain.CopyConnection(cs.currentConnection, databaseName)
	if err := cpy.Connect(); err != nil {
		return nil, fmt.Errorf("failed to connect to database %s: %w", databaseName, err)
	}
	defer cpy.Disconnect()

	query := `
		SELECT table_name 
		FROM information_schema.tables 
		WHERE table_schema = 'public' 
		AND table_type = 'BASE TABLE'
		ORDER BY table_name
	`

	rows, err := cpy.DB().Query(query)
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
	if !cs.HasActiveConnection() {
		return nil, fmt.Errorf("no active database connection")
	}

	cpy := domain.CopyConnection(cs.currentConnection, databaseName)
	if err := cpy.Connect(); err != nil {
		return nil, fmt.Errorf("failed to connect to database %s: %w", databaseName, err)
	}
	defer cpy.Disconnect()

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

	rows, err := cpy.DB().Query(query, tableName)
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

// ExecuteQuery executes a SQL query against a specific database and returns the results
func (cs *ConnectionService) ExecuteQuery(databaseName, query string) (*types.QueryResult, error) {
	if !cs.HasActiveConnection() {
		return nil, fmt.Errorf("no active database connection")
	}

	cpy := domain.CopyConnection(cs.currentConnection, databaseName)
	if err := cpy.Connect(); err != nil {
		return nil, fmt.Errorf("failed to connect to database %s: %w", databaseName, err)
	}
	defer cpy.Disconnect()

	start := time.Now()

	// Check if it's a SELECT query or DML/DDL
	rows, err := cpy.DB().Query(query)
	if err != nil {
		// If Query fails, try Exec for DML/DDL statements
		result, execErr := cpy.DB().Exec(query)
		if execErr != nil {
			return nil, fmt.Errorf("failed to execute query: %w", execErr)
		}

		rowsAffected, _ := result.RowsAffected()
		duration := time.Since(start).Milliseconds()

		return &types.QueryResult{
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

	return &types.QueryResult{
		Columns:      columns,
		Rows:         resultRows,
		RowsAffected: int64(len(resultRows)),
		Duration:     duration,
	}, nil
}

// ListConnections returns a list of saved connections
func (cs *ConnectionService) ListConnections() ([]types.ConnectionSummary, error) {
	connections, err := cs.repo.List()
	if err != nil {
		return nil, fmt.Errorf("failed to list connections: %w", err)
	}

	summaries := make([]types.ConnectionSummary, len(connections))
	for i, conn := range connections {
		summaries[i] = types.ConnectionSummary{
			ID:   conn.ID(),
			Host: conn.Host(),
			Port: conn.Port(),
		}
	}

	return summaries, nil
}

// AnalyzeConnectionMetadata analyzes the current connection and persists the metadata
func (cs *ConnectionService) AnalyzeConnectionMetadata() error {
	if !cs.HasActiveConnection() {
		return fmt.Errorf("no active database connection")
	}

	// Use the domain method to analyze metadata
	domainMetadata, err := domain.AnalyzeConnectionMetadata(cs.currentConnection)
	if err != nil {
		return fmt.Errorf("failed to analyze connection metadata: %w", err)
	}

	// Persist the metadata using the repository
	if err := cs.metadataRepo.Save(domainMetadata); err != nil {
		return fmt.Errorf("failed to persist connection metadata: %w", err)
	}

	return nil
}

// Helper function to convert types.DatabaseConfig to domain.Connection
func (cs *ConnectionService) configToDomainConnection(id string, config types.DatabaseConfig) (*domain.Connection, error) {
	if config.UseConnectionString && config.ConnectionString != "" {
		return domain.NewConnectionFromString(id, config.ConnectionString)
	}

	arguments := make(map[string]string)
	if config.SSLMode != "" {
		arguments["sslmode"] = config.SSLMode
	}

	return domain.NewConnection(id, config.Host, config.Port, config.Database, config.Username, config.Password, arguments), nil
}
