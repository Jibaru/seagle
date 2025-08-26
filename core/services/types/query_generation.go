package types

// GenerateQueryRequest represents a request to generate a SQL query
type GenerateQueryRequest struct {
	Database string `json:"database"`
	Prompt   string `json:"prompt"`
}

// GenerateQueryResult represents the result of query generation
type GenerateQueryResult struct {
	GeneratedQuery string `json:"generatedQuery"`
	OriginalPrompt string `json:"originalPrompt"`
}