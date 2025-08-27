package domain

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

var supportedVendors = map[string]bool{
	"postgresql": true,
	"mysql":      true,
}

type Connection struct {
	id        string
	vendor    string
	host      string
	port      int
	database  string
	username  string
	password  string
	arguments map[string]string
}

func NewConnection(id, vendor, host string, port int, database, username, password string, arguments map[string]string) (*Connection, error) {
	if !supportedVendors[vendor] {
		return nil, fmt.Errorf("unsupported vendor: %s", vendor)
	}

	return &Connection{
		id:        id,
		vendor:    vendor,
		host:      host,
		port:      port,
		database:  database,
		username:  username,
		password:  password,
		arguments: arguments,
	}, nil
}

func CopyConnection(conn *Connection, database string) *Connection {
	return &Connection{
		id:        conn.id,
		vendor:    conn.vendor,
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

	// Parse URI format: scheme://username:password@host:port/database?param=value
	parsedURL, err := url.Parse(connStr)
	if err != nil {
		return nil, fmt.Errorf("invalid connection string format: %v", err)
	}

	var vendor string
	var defaultPort int

	switch parsedURL.Scheme {
	case "postgresql", "postgres":
		vendor = "postgresql"
		defaultPort = 5432
	case "mysql":
		vendor = "mysql"
		defaultPort = 3306
	default:
		return nil, fmt.Errorf("unsupported scheme: %s", parsedURL.Scheme)
	}

	// Extract host and port
	host := parsedURL.Hostname()
	if host == "" {
		return nil, fmt.Errorf("host is required in connection string")
	}

	port := defaultPort
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

	return NewConnection(id, vendor, host, port, database, username, password, arguments)
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
		vendor:    data["vendor"].(string),
		host:      data["host"].(string),
		port:      int(data["port"].(float64)),
		database:  data["database"].(string),
		username:  data["username"].(string),
		password:  data["password"].(string),
		arguments: arguments,
	}
}

func (c *Connection) ID() string {
	return c.id
}

func (c *Connection) Vendor() string {
	return c.vendor
}

func (c *Connection) Host() string {
	return c.host
}

func (c *Connection) Port() int {
	return c.port
}

func (c *Connection) Map() map[string]interface{} {
	return map[string]interface{}{
		"id":        c.id,
		"vendor":    c.vendor,
		"host":      c.host,
		"port":      c.port,
		"database":  c.database,
		"username":  c.username,
		"password":  c.password,
		"arguments": c.arguments,
	}
}

func (c *Connection) connectionString() string {
	switch c.vendor {
	case "postgresql":
		dbURL := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s", url.QueryEscape(c.username), url.QueryEscape(c.password), c.host, c.port, c.database)
		if len(c.arguments) > 0 {
			args := make([]string, 0, len(c.arguments))
			for k, v := range c.arguments {
				args = append(args, fmt.Sprintf("%s=%s", url.QueryEscape(k), url.QueryEscape(v)))
			}
			dbURL += "?" + strings.Join(args, "&")
		}
		return dbURL
	case "mysql":
		dbURL := fmt.Sprintf("mysql://%s:%s@%s:%d/%s", url.QueryEscape(c.username), url.QueryEscape(c.password), c.host, c.port, c.database)
		if len(c.arguments) > 0 {
			args := make([]string, 0, len(c.arguments))
			for k, v := range c.arguments {
				args = append(args, fmt.Sprintf("%s=%s", url.QueryEscape(k), url.QueryEscape(v)))
			}
			dbURL += "?" + strings.Join(args, "&")
		}
		return dbURL
	default:
		return ""
	}
}
