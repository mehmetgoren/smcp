package utils

import (
	"encoding/json"
	"log"
)

func SerializeJson(data interface{}) (string, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		log.Println("Error serializing json: ", err)
		return "", err
	}

	return string(bytes), nil
}

func DeserializeJson(data string, result interface{}) error {
	err := json.Unmarshal([]byte(data), result)
	if err != nil {
		log.Println("Error deserializing json: ", err)
		return err
	}
	return nil
}
