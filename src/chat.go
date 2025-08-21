package src

import (
	"encoding/json"
	"log"
)

func Chat(data map[string]interface{}) (interface{}, error) {
    _, err := json.MarshalIndent(data, "", "  ")
    if err != nil {
        log.Printf("Error marshaling data: %v", err)
		return nil,nil
    }
	return nil,nil
}