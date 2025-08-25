package domain

import (
	"database/sql"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

type Connection struct {
	id        string
	host      string
	port      int
	database  string
	username  string
	password  string
	arguments map[string]string

	conn *sql.DB
}

func NewConnection(id, host string, port int, database, username, password string, arguments map[string]string) *Connection {
	return &Connection{
		id:        id,
		host:      host,
		port:      port,
		database:  database,
		username:  username,
		password:  password,
		arguments: arguments,
	}
}

func CopyConnection(conn *Connection, database string) *Connection {
	return &Connection{
		id:        conn.id,
		host:      conn.host,
		port:      conn.port,
		database:  database,
		username:  conn.username,
		password:  conn.password,
		arguments: conn.arguments,
	}
}

func NewConnectionFromString(id, connStr string) (*Connection, error) {
	arguments := make(map[string]string)

	// Parse PostgreSQL URI format: postgresql://username:password@host:port/database?param=value
	parsedURL, err := url.Parse(connStr)
	if err != nil {
		return nil, fmt.Errorf("invalid connection string format: %v", err)
	}

	if parsedURL.Scheme != "postgresql" && parsedURL.Scheme != "postgres" {
		return nil, fmt.Errorf("unsupported scheme: %s", parsedURL.Scheme)
	}

	// Extract host and port
	host := parsedURL.Hostname()
	if host == "" {
		return nil, fmt.Errorf("host is required in connection string")
	}

	port := 5432 // default PostgreSQL port
	if parsedURL.Port() != "" {
		if p, err := strconv.Atoi(parsedURL.Port()); err == nil {
			port = p
		}
	}

	// Extract username and password
	var username, password string
	if parsedURL.User != nil {
		username = parsedURL.User.Username()
		if pass, ok := parsedURL.User.Password(); ok {
			password = pass
		}
	}

	if username == "" {
		return nil, fmt.Errorf("username is required in connection string")
	}

	// Extract database name
	database := strings.TrimPrefix(parsedURL.Path, "/")
	if database == "" {
		return nil, fmt.Errorf("database is required in connection string")
	}

	// Extract query parameters as arguments
	for key, values := range parsedURL.Query() {
		if len(values) > 0 {
			arguments[key] = values[0]
		}
	}

	return NewConnection(id, host, port, database, username, password, arguments), nil
}

func NewConnectionFromMap(data map[string]interface{}) *Connection {
	arguments := make(map[string]string)
	if args, ok := data["arguments"].(map[string]interface{}); ok {
		for k, v := range args {
			if strVal, ok := v.(string); ok {
				arguments[k] = strVal
			}
		}
	}

	return &Connection{
		id:        data["id"].(string),
		host:      data["host"].(string),
		port:      int(data["port"].(float64)),
		database:  data["database"].(string),
		username:  data["username"].(string),
		password:  data["password"].(string),
		arguments: arguments,
	}
}

func (c *Connection) Map() map[string]interface{} {
	return map[string]interface{}{
		"id":        c.id,
		"host":      c.host,
		"port":      c.port,
		"database":  c.database,
		"username":  c.username,
		"password":  c.password,
		"arguments": c.arguments,
	}
}

func (c *Connection) connectionString() string {
	dbURL := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s", url.QueryEscape(c.username), url.QueryEscape(c.password), c.host, c.port, c.database)

	if len(c.arguments) > 0 {
		args := make([]string, 0, len(c.arguments))
		for k, v := range c.arguments {
			args = append(args, fmt.Sprintf("%s=%s", url.QueryEscape(k), url.QueryEscape(v)))
		}
		dbURL += "?" + strings.Join(args, "&")
	}

	return dbURL
}

func (c *Connection) Connect() error {
	if c.conn == nil {
		db, err := sql.Open("postgres", c.connectionString())
		if err != nil {
			return fmt.Errorf("failed to open database connection: %w", err)
		}

		c.conn = db
	}

	if err := c.conn.Ping(); err != nil {
		c.conn.Close()
		return fmt.Errorf("failed to ping database: %w", err)
	}

	return nil
}

func (c *Connection) Disconnect() error {
	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			return fmt.Errorf("failed to close database connection: %w", err)
		}
		c.conn = nil
	}
	return nil
}

func (c *Connection) IsConnected() bool {
	return c.conn != nil
}

func (c *Connection) DB() *sql.DB {
	return c.conn
}
