package chat

import (
	"encoding/json"
	"log"
)

func Chat(data map[string]interface{}) (interface{}, error) {
    jsonData, err := json.MarshalIndent(data, "", "  ")
    if err != nil {
        log.Printf("Error marshaling data: %v", err)
		return nil,nil
    }
    log.Printf("Chat data:\n%s", string(jsonData))
	return nil,nil
}