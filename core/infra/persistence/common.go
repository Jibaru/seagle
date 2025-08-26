package persistence

import (
	"encoding/json"
	"os"
)

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
