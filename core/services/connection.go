package services

import (
	"fmt"

	"seagle/core/domain"
	"seagle/core/services/types"
)

// ConnectionService manages database connections
type ConnectionService struct {
	currentConnection *domain.Connection
	repo              domain.ConnectionRepo
	metadataRepo      domain.MetadataRepo
	connectionService *domain.ConnectionService
	metadataFactory   *domain.MetadataFactory
	openaiClient      *OpenAIClient
}

// NewConnectionService creates a new ConnectionService instance
func NewConnectionService(
	repo domain.ConnectionRepo,
	metadataRepo domain.MetadataRepo,
	connectionService *domain.ConnectionService,
	metadataFactory *domain.MetadataFactory,
	openaiClient *OpenAIClient,
) *ConnectionService {
	return &ConnectionService{
		repo:              repo,
		metadataRepo:      metadataRepo,
		connectionService: connectionService,
		metadataFactory:   metadataFactory,
		openaiClient:      openaiClient,
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

	if err := cs.connectionService.Connect(domainConn); err != nil {
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

	if err := cs.connectionService.Connect(conn); err != nil {
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

	if err := cs.connectionService.Connect(domainConn); err != nil {
		return err
	}

	return nil
}

// Disconnect closes the current database connection
func (cs *ConnectionService) Disconnect() error {
	return cs.connectionService.Disconnect(cs.currentConnection)
}

// HasActiveConnection checks if there is an active connection
func (cs *ConnectionService) HasActiveConnection() bool {
	return cs.currentConnection != nil
}

// GetDatabases returns a list of databases available in the connected PostgreSQL instance
func (cs *ConnectionService) GetDatabases() ([]string, error) {
	if !cs.HasActiveConnection() {
		return nil, fmt.Errorf("no active database connection")
	}

	return cs.connectionService.GetDatabaseNames(cs.currentConnection)
}

// GetTables returns a list of tables for a specific database
func (cs *ConnectionService) GetTables(databaseName string) ([]string, error) {
	if !cs.HasActiveConnection() {
		return nil, fmt.Errorf("no active database connection")
	}

	cpy := domain.CopyConnection(cs.currentConnection, databaseName)
	if err := cs.connectionService.Connect(cpy); err != nil {
		return nil, fmt.Errorf("failed to connect to database %s: %w", databaseName, err)
	}
	defer cs.connectionService.Disconnect(cpy)

	return cs.connectionService.GetTableNames(cpy, databaseName)
}

// GetTableColumns returns the columns for a specific table in a database
func (cs *ConnectionService) GetTableColumns(databaseName, tableName string) ([]types.TableColumn, error) {
	if !cs.HasActiveConnection() {
		return nil, fmt.Errorf("no active database connection")
	}

	cpy := domain.CopyConnection(cs.currentConnection, databaseName)
	if err := cs.connectionService.Connect(cpy); err != nil {
		return nil, fmt.Errorf("failed to connect to database %s: %w", databaseName, err)
	}
	defer cs.connectionService.Disconnect(cpy)

	columns, err := cs.connectionService.GetTableColumns(cpy, databaseName, tableName)
	if err != nil {
		return nil, fmt.Errorf("failed to get columns for table %s: %w", tableName, err)
	}

	result := make([]types.TableColumn, len(columns))
	for i, col := range columns {
		result[i] = types.TableColumn{
			Name:         col.Name(),
			DataType:     col.DataType(),
			IsNullable:   col.IsNullable(),
			DefaultValue: col.DefaultValue(),
		}
	}

	return result, nil
}

// ExecuteQuery executes a SQL query against a specific database and returns the results
func (cs *ConnectionService) ExecuteQuery(databaseName, query string) (*types.QueryResult, error) {
	if !cs.HasActiveConnection() {
		return nil, fmt.Errorf("no active database connection")
	}

	cpy := domain.CopyConnection(cs.currentConnection, databaseName)
	if err := cs.connectionService.Connect(cpy); err != nil {
		return nil, fmt.Errorf("failed to connect to database %s: %w", databaseName, err)
	}
	defer cs.connectionService.Disconnect(cpy)

	res, err := cs.connectionService.ExecQuery(cpy, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return &types.QueryResult{
		Columns:      res.Columns,
		Rows:         res.Rows,
		RowsAffected: res.RowsAffected,
		Duration:     res.Duration,
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
	domainMetadata, err := cs.metadataFactory.NewConnectionMetadata(cs.currentConnection)
	if err != nil {
		return fmt.Errorf("failed to analyze connection metadata: %w", err)
	}

	// Persist the metadata using the repository
	if err := cs.metadataRepo.Save(domainMetadata); err != nil {
		return fmt.Errorf("failed to persist connection metadata: %w", err)
	}

	return nil
}

// GenerateQuery generates a SQL query using AI based on natural language input
func (cs *ConnectionService) GenerateQuery(request types.GenerateQueryRequest) (*types.GenerateQueryResult, error) {
	if !cs.HasActiveConnection() {
		return nil, fmt.Errorf("no active database connection")
	}

	// Check if metadata exists for the current connection
	connectionID := cs.currentConnection.ID()
	metadata, err := cs.metadataRepo.FindByConnectionID(connectionID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing metadata: %w", err)
	}

	// If metadata doesn't exist, generate it first
	if metadata == nil {
		if err := cs.AnalyzeConnectionMetadata(); err != nil {
			return nil, fmt.Errorf("failed to generate metadata: %w", err)
		}

		// Retrieve the newly generated metadata
		metadata, err = cs.metadataRepo.FindByConnectionID(connectionID)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve generated metadata: %w", err)
		}

		if metadata == nil {
			return nil, fmt.Errorf("metadata generation failed")
		}
	}

	// Use OpenAI to generate the query
	generatedQuery, err := cs.openaiClient.GenerateQuery(request.Prompt, metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to generate query with AI: %w", err)
	}

	return &types.GenerateQueryResult{
		GeneratedQuery: generatedQuery,
		OriginalPrompt: request.Prompt,
	}, nil
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

	return domain.NewConnection(id, config.Vendor, config.Host, config.Port, config.Database, config.Username, config.Password, arguments)
}
