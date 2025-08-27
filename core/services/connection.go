package services

import (
	"fmt"

	"seagle/core/domain"
	"seagle/core/services/types"
)

// ConnectionService manages database connections
type ConnectionService struct {
	repo            domain.ConnectionRepo
	metadataRepo    domain.MetadataRepo
	serviceFactory  *domain.ServiceFactory
	metadataFactory *domain.MetadataFactory
	openaiClient    *OpenAIClient
}

// NewConnectionService creates a new ConnectionService instance
func NewConnectionService(
	repo domain.ConnectionRepo,
	metadataRepo domain.MetadataRepo,
	serviceFactory *domain.ServiceFactory,
	metadataFactory *domain.MetadataFactory,
	openaiClient *OpenAIClient,
) *ConnectionService {
	return &ConnectionService{
		repo:            repo,
		metadataRepo:    metadataRepo,
		serviceFactory:  serviceFactory,
		metadataFactory: metadataFactory,
		openaiClient:    openaiClient,
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

	dbService, err := cs.serviceFactory.NewDatabaseService(domainConn)
	if err != nil {
		return nil, fmt.Errorf("failed to create database service: %w", err)
	}

	if err := dbService.Connect(domainConn); err != nil {
		return nil, err
	}

	return &types.DatabaseConnection{
		ID:          domainConn.ID(),
		Config:      config,
		IsConnected: true,
	}, nil
}

func (cs *ConnectionService) ConnectByID(id string) (*types.DatabaseConnection, error) {
	conn, dbService, err := cs.lookup(id)
	if err != nil {
		return nil, err
	}

	if err := dbService.Connect(conn); err != nil {
		return nil, err
	}

	return &types.DatabaseConnection{
		IsConnected: true,
	}, nil
}

func (cs *ConnectionService) lookup(id string) (*domain.Connection, domain.DatabaseService, error) {
	conn, err := cs.repo.FindByID(id)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to find connection by ID: %w", err)
	}
	if conn == nil {
		return nil, nil, fmt.Errorf("connection with ID %s not found", id)
	}

	dbService, err := cs.serviceFactory.NewDatabaseService(conn)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create database service: %w", err)
	}

	return conn, dbService, nil
}

// TestConnection tests the database connection with given parameters
func (cs *ConnectionService) TestConnection(config types.DatabaseConfig) error {
	domainConn, err := cs.configToDomainConnection(cs.repo.NextID(), config)
	if err != nil {
		return err
	}

	dbService, err := cs.serviceFactory.NewDatabaseService(domainConn)
	if err != nil {
		return fmt.Errorf("failed to create database service: %w", err)
	}

	if err := dbService.Connect(domainConn); err != nil {
		return err
	}

	return dbService.Disconnect(domainConn)
}

// Disconnect closes the current database connection
func (cs *ConnectionService) Disconnect(id string) error {
	conn, err := cs.repo.FindByID(id)
	if err != nil {
		return fmt.Errorf("failed to find connection by ID: %w", err)
	}
	if conn == nil {
		return fmt.Errorf("connection with ID %s not found", id)
	}

	dbService, err := cs.serviceFactory.NewDatabaseService(conn)
	if err != nil {
		return fmt.Errorf("failed to create database service: %w", err)
	}

	if err := dbService.Disconnect(conn); err != nil {
		return err
	}

	return nil
}

// GetDatabases returns a list of databases available in the connected database instance
func (cs *ConnectionService) GetDatabases(id string) ([]string, error) {
	conn, dbService, err := cs.lookup(id)
	if err != nil {
		return nil, err
	}

	if err := dbService.Connect(conn); err != nil {
		return nil, err
	}

	return dbService.GetDatabaseNames(conn)
}

// GetTables returns a list of tables for a specific database
func (cs *ConnectionService) GetTables(originalID string, databaseName string) ([]string, error) {
	conn, dbService, err := cs.lookup(originalID)
	if err != nil {
		return nil, err
	}

	cpy := domain.CopyConnection(conn, databaseName)

	if err := dbService.Connect(cpy); err != nil {
		return nil, fmt.Errorf("failed to connect to database %s: %w", databaseName, err)
	}
	defer dbService.Disconnect(cpy)

	return dbService.GetTableNames(cpy, databaseName)
}

// GetTableColumns returns the columns for a specific table in a database
func (cs *ConnectionService) GetTableColumns(originalID, databaseName, tableName string) ([]types.TableColumn, error) {
	conn, dbService, err := cs.lookup(originalID)
	if err != nil {
		return nil, err
	}

	cpy := domain.CopyConnection(conn, databaseName)

	if err := dbService.Connect(cpy); err != nil {
		return nil, fmt.Errorf("failed to connect to database %s: %w", databaseName, err)
	}
	defer dbService.Disconnect(cpy)

	columns, err := dbService.GetTableColumns(cpy, databaseName, tableName)
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
func (cs *ConnectionService) ExecuteQuery(originalID, databaseName, query string) (*types.QueryResult, error) {
	conn, dbService, err := cs.lookup(originalID)
	if err != nil {
		return nil, err
	}

	cpy := domain.CopyConnection(conn, databaseName)

	if err := dbService.Connect(cpy); err != nil {
		return nil, fmt.Errorf("failed to connect to database %s: %w", databaseName, err)
	}
	defer dbService.Disconnect(cpy)

	res, err := dbService.ExecQuery(cpy, query)
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
func (cs *ConnectionService) AnalyzeConnectionMetadata(id string) error {
	conn, err := cs.repo.FindByID(id)
	if err != nil {
		return fmt.Errorf("failed to find connection by ID: %w", err)
	}

	// Use the domain method to analyze metadata
	domainMetadata, err := cs.metadataFactory.NewConnectionMetadata(conn)
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
func (cs *ConnectionService) GenerateQuery(id string, request types.GenerateQueryRequest) (*types.GenerateQueryResult, error) {
	conn, err := cs.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to find connection by ID: %w", err)
	}

	metadata, err := cs.metadataRepo.FindByConnectionID(conn.ID())
	if err != nil {
		return nil, fmt.Errorf("failed to check existing metadata: %w", err)
	}

	// If metadata doesn't exist, generate it first
	if metadata == nil {
		if err := cs.AnalyzeConnectionMetadata(conn.ID()); err != nil {
			return nil, fmt.Errorf("failed to generate metadata: %w", err)
		}

		// Retrieve the newly generated metadata
		metadata, err = cs.metadataRepo.FindByConnectionID(conn.ID())
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
