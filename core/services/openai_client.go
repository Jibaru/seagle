package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
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
func (c *OpenAIClient) GenerateQuery(userPrompt string, metadata *domain.ConnectionMetadata) (string, error) {
	if c.apiKey == "" {
		return "", fmt.Errorf("OpenAI API key not configured")
	}

	// Build the system prompt with database metadata
	systemPrompt := c.buildSystemPrompt(metadata)

	// Create the request
	request := OpenAIRequest{
		Model: "gpt-4o",
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
		MaxTokens:   1000,
		Temperature: 0.1,
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

	return response.Choices[0].Message.Content, nil
}

// buildSystemPrompt creates a system prompt with database schema information
func (c *OpenAIClient) buildSystemPrompt(metadata *domain.ConnectionMetadata) string {
	prompt := `You are a PostgreSQL expert. Generate SQL queries based on the user's natural language request.

Database Schema Information:
`

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

	prompt += `
Rules:
1. Generate only PostgreSQL-compatible SQL queries
2. Use proper table and column names from the schema above
3. Include appropriate WHERE clauses, JOINs, and other SQL constructs as needed
4. Return only the SQL query without explanation or markdown formatting
5. Ensure the query is syntactically correct and executable
6. Use meaningful aliases when joining tables
7. Consider performance best practices

Respond with only the SQL query.`

	return prompt
}

// makeRequest makes an HTTP request to OpenAI API
func (c *OpenAIClient) makeRequest(request OpenAIRequest) (*OpenAIResponse, error) {
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

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

	var response OpenAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}
