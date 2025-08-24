package types

import "database/sql"

// DatabaseConfig holds database connection parameters
type DatabaseConfig struct {
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
	DB          *sql.DB        `json:"-"`
}

// TableColumn represents a column in a database table
type TableColumn struct {
	Name         string `json:"name"`
	DataType     string `json:"dataType"`
	IsNullable   bool   `json:"isNullable"`
	DefaultValue string `json:"defaultValue,omitempty"`
}