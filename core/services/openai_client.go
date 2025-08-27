package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"seagle/core/domain"
)

// OpenAIClient handles communication with OpenAI API
type OpenAIClient struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

// NewOpenAIClient creates a new OpenAI client
func NewOpenAIClient(apiKey string) *OpenAIClient {
	return &OpenAIClient{
		apiKey:  apiKey,
		baseURL: "https://api.openai.com/v1",
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// OpenAIRequest represents the request structure for OpenAI API
type OpenAIRequest struct {
	Model       string          `json:"model"`
	Messages    []OpenAIMessage `json:"messages"`
	MaxTokens   int             `json:"max_tokens"`
	Temperature float64         `json:"temperature"`
}

// OpenAIMessage represents a message in the conversation
type OpenAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// OpenAIResponse represents the response from OpenAI API
type OpenAIResponse struct {
	ID      string         `json:"id"`
	Object  string         `json:"object"`
	Created int64          `json:"created"`
	Model   string         `json:"model"`
	Choices []OpenAIChoice `json:"choices"`
	Usage   OpenAIUsage    `json:"usage"`
	Error   *OpenAIError   `json:"error,omitempty"`
}

// OpenAIChoice represents a choice in the response
type OpenAIChoice struct {
	Index        int           `json:"index"`
	Message      OpenAIMessage `json:"message"`
	FinishReason string        `json:"finish_reason"`
}

// OpenAIUsage represents token usage information
type OpenAIUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// OpenAIError represents an error from OpenAI API
type OpenAIError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Code    string `json:"code"`
}

// GenerateQuery generates a SQL query based on natural language input and database metadata
func (c *OpenAIClient) GenerateQuery(userPrompt string, metadata *domain.ConnectionMetadata, connection *domain.Connection) (string, error) {
	if c.apiKey == "" {
		return "", fmt.Errorf("OpenAI API key not configured")
	}

	// Build the system prompt with database metadata
	systemPrompt := c.buildSystemPrompt(metadata, connection)

	// Create the request
	request := OpenAIRequest{
		Model: "gpt-4o", // Using GPT-4o for better SQL generation
		Messages: []OpenAIMessage{
			{
				Role:    "system",
				Content: systemPrompt,
			},
			{
				Role:    "user",
				Content: userPrompt,
			},
		},
		MaxTokens:   1500, // Increased for complex queries
		Temperature: 0.5,  // Low temperature for consistent SQL generation
	}

	// Make the API call
	response, err := c.makeRequest(request)
	if err != nil {
		return "", fmt.Errorf("failed to call OpenAI API: %w", err)
	}

	if response.Error != nil {
		return "", fmt.Errorf("OpenAI API error: %s", response.Error.Message)
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no response from OpenAI API")
	}

	// Clean the response to remove code block markers
	cleanedQuery := c.cleanSQLResponse(response.Choices[0].Message.Content)
	return cleanedQuery, nil
}

// cleanSQLResponse removes code block markers and extra formatting from AI response
func (c *OpenAIClient) cleanSQLResponse(response string) string {
	// Remove leading/trailing whitespace
	cleaned := strings.TrimSpace(response)

	// Remove SQL code block markers (```sql and ```)
	// Pattern matches: ```sql\nCODE\n``` or ```\nCODE\n```
	sqlBlockPattern := `(?s)^\s*` + "```" + `(?:sql)?\s*\n?(.*?)\n?\s*` + "```" + `\s*$`
	sqlBlockRegex := regexp.MustCompile(sqlBlockPattern)
	if match := sqlBlockRegex.FindStringSubmatch(cleaned); len(match) > 1 {
		cleaned = match[1]
	}

	// Remove any remaining backticks at start/end
	cleaned = strings.Trim(cleaned, "`")

	// Remove extra whitespace and newlines
	cleaned = strings.TrimSpace(cleaned)

	// Ensure the query ends with semicolon if it doesn't already
	if !strings.HasSuffix(cleaned, ";") {
		cleaned += ";"
	}

	return cleaned
}

// buildSystemPrompt creates a system prompt with database schema information
func (c *OpenAIClient) buildSystemPrompt(metadata *domain.ConnectionMetadata, connection *domain.Connection) string {
	databaseType := connection.Vendor()
	var expertType string

	switch databaseType {
	case "postgresql":
		expertType = "PostgreSQL"
	case "mysql":
		expertType = "MySQL"
	default:
		expertType = "SQL"
	}

	prompt := fmt.Sprintf(`You are a %s expert. Generate SQL queries based on the user's natural language request.

Database Type: %s
Database Schema Information:
`, expertType, expertType)

	for _, db := range metadata.Databases() {
		prompt += fmt.Sprintf("\nDatabase: %s\n", db.Name())

		for _, table := range db.Tables() {
			prompt += fmt.Sprintf("  Table: %s.%s\n", table.Schema(), table.Name())

			for _, col := range table.Columns() {
				nullable := "NOT NULL"
				if col.IsNullable() {
					nullable = "NULL"
				}

				defaultValue := ""
				if col.DefaultValue() != "" {
					defaultValue = fmt.Sprintf(" DEFAULT %s", col.DefaultValue())
				}

				prompt += fmt.Sprintf("    - %s: %s %s%s\n",
					col.Name(), col.DataType(), nullable, defaultValue)
			}
			prompt += "\n"
		}
	}

	// Add database-specific rules
	var rules string
	switch databaseType {
	case "postgresql":
		rules = `
Rules:
1. Generate only PostgreSQL-compatible SQL queries
2. Use proper table and column names from the schema above
3. Include appropriate WHERE clauses, JOINs, and other SQL constructs as needed
4. Return ONLY the raw SQL query without any markdown formatting, code blocks, or explanations
5. DO NOT use code block markers or any markdown formatting
6. Ensure the query is syntactically correct and executable
7. Use meaningful aliases when joining tables
8. Consider PostgreSQL-specific features like SERIAL, LIMIT, and case-sensitive identifiers
9. Use double quotes for identifiers if needed
10. Consider performance best practices

Respond with only the raw SQL query, no formatting.`
	case "mysql":
		rules = `
Rules:
1. Generate only MySQL-compatible SQL queries
2. Use proper table and column names from the schema above
3. Include appropriate WHERE clauses, JOINs, and other SQL constructs as needed
4. Return ONLY the raw SQL query without any markdown formatting, code blocks, or explanations
5. DO NOT use code block markers or any markdown formatting
6. Ensure the query is syntactically correct and executable
7. Use meaningful aliases when joining tables
8. Consider MySQL-specific features like AUTO_INCREMENT, LIMIT, and backticks for identifiers
9. Use backticks for identifiers if needed
10. Consider performance best practices

Respond with only the raw SQL query, no formatting.`
	default:
		rules = `
Rules:
1. Generate standard SQL queries compatible with the database type
2. Use proper table and column names from the schema above
3. Include appropriate WHERE clauses, JOINs, and other SQL constructs as needed
4. Return ONLY the raw SQL query without any markdown formatting, code blocks, or explanations
5. DO NOT use code block markers or any markdown formatting
6. Ensure the query is syntactically correct and executable
7. Use meaningful aliases when joining tables
8. Consider performance best practices

Respond with only the raw SQL query, no formatting.`
	}

	prompt += rules

	return prompt
}

// makeRequest makes an HTTP request to OpenAI API
func (c *OpenAIClient) makeRequest(request OpenAIRequest) (*OpenAIResponse, error) {
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	log.Printf("Making OpenAI request with model: %s, max_tokens: %d", request.Model, request.MaxTokens)

	req, err := http.NewRequest("POST", c.baseURL+"/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make HTTP request: %w", err)
	}
	defer resp.Body.Close()

	// Check HTTP status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OpenAI API returned status code: %d", resp.StatusCode)
	}

	var response OpenAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Log token usage for monitoring
	if response.Usage.TotalTokens > 0 {
		log.Printf("OpenAI API usage - Prompt tokens: %d, Completion tokens: %d, Total tokens: %d",
			response.Usage.PromptTokens, response.Usage.CompletionTokens, response.Usage.TotalTokens)
	}

	return &response, nil
}
