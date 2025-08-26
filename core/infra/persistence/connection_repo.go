package persistence

import (
	"seagle/core/domain"

	"github.com/google/uuid"
)

type ConnectionRepo struct {
	filename string
}

func NewConnection(filename string) *ConnectionRepo {
	return &ConnectionRepo{
		filename: filename,
	}
}

func (r *ConnectionRepo) NextID() string {
	return uuid.NewString()
}

func (r *ConnectionRepo) Save(connection *domain.Connection) error {
	connections, err := r.load()
	if err != nil {
		return err
	}

	exists := false
	for i, conn := range connections {
		if conn.ID() == connection.ID() {
			connections[i] = connection
			exists = true
			break
		}
	}

	if !exists {
		connections = append(connections, connection)
	}

	return saveDataToFile(r.filename, r.toDataMap(connections))
}

func (r *ConnectionRepo) List() ([]*domain.Connection, error) {
	return r.load()
}

func (r *ConnectionRepo) load() ([]*domain.Connection, error) {
	data, err := loadDataFromFile(r.filename)
	if err != nil {
		return nil, err
	}

	connections := []*domain.Connection{}
	if existingConnections, ok := data["connections"].([]interface{}); ok {
		for _, conn := range existingConnections {
			if connMap, ok := conn.(map[string]interface{}); ok {
				connections = append(connections, domain.NewConnectionFromMap(connMap))
			}
		}
	}

	return connections, nil
}

func (r *ConnectionRepo) toDataMap(connections []*domain.Connection) map[string]interface{} {
	connMaps := make([]map[string]interface{}, len(connections))
	for i, conn := range connections {
		connMaps[i] = conn.Map()
	}
	return map[string]interface{}{
		"connections": connMaps,
	}
}
