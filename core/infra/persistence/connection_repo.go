package persistence

import (
	"encoding/json"
	"os"

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
	data, err := loadDataFromFile(r.filename)
	if err != nil {
		data = make(map[string]interface{})
	}

	var connections []map[string]interface{}
	if existingConnections, ok := data["connections"].([]interface{}); ok {
		for _, conn := range existingConnections {
			if connMap, ok := conn.(map[string]interface{}); ok {
				connections = append(connections, connMap)
			}
		}
	}

	connections = append(connections, connection.Map())
	data["connections"] = connections

	return saveDataToFile(r.filename, data)
}

func saveDataToFile(filename string, data map[string]interface{}) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, jsonData, 0644)
}

func loadDataFromFile(filename string) (map[string]interface{}, error) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return map[string]interface{}{"connections": []interface{}{}}, nil
	}

	fileData, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	if len(fileData) == 0 {
		return map[string]interface{}{"connections": []interface{}{}}, nil
	}

	var data map[string]interface{}
	err = json.Unmarshal(fileData, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
