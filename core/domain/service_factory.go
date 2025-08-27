package domain

import "fmt"

// ServiceFactory creates appropriate database services based on vendor
type ServiceFactory struct{}

// NewServiceFactory creates a new service factory
func NewServiceFactory() *ServiceFactory {
	return &ServiceFactory{}
}

// NewDatabaseService creates the appropriate database service for the given connection
func (f *ServiceFactory) NewDatabaseService(c *Connection) (DatabaseService, error) {
	switch c.Vendor() {
	case "postgresql":
		return NewPostgreSQLService(), nil
	case "mysql":
		return NewMySQLService(), nil
	default:
		return nil, fmt.Errorf("unsupported database vendor: %s", c.Vendor())
	}
}
