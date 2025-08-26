package types

// DatabaseConfig holds database connection parameters
type DatabaseConfig struct {
	Vendor              string `json:"vendor"`
	Host                string `json:"host"`
	Port                int    `json:"port"`
	Database            string `json:"database"`
	Username            string `json:"username"`
	Password            string `json:"password"`
	SSLMode             string `json:"sslmode"`
	ConnectionString    string `json:"connectionString"`
	UseConnectionString bool   `json:"useConnectionString"`
}

// DatabaseConnection represents a database connection
type DatabaseConnection struct {
	Config      DatabaseConfig `json:"config"`
	IsConnected bool           `json:"isConnected"`
}

// TableColumn represents a column in a database table
type TableColumn struct {
	Name         string `json:"name"`
	DataType     string `json:"dataType"`
	IsNullable   bool   `json:"isNullable"`
	DefaultValue string `json:"defaultValue,omitempty"`
}

// QueryResult represents the result of a SQL query
type QueryResult struct {
	Columns      []string        `json:"columns"`
	Rows         [][]interface{} `json:"rows"`
	RowsAffected int64           `json:"rowsAffected"`
	Duration     int64           `json:"duration"` // in milliseconds
}

type ConnectionSummary struct {
	ID   string `json:"id"`
	Host string `json:"host"`
	Port int    `json:"port"`
}
