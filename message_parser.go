package main

import (
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"log"
)

// todo:need generic, changes implementation after 1.18
type MessageParser interface {
	Parse(redisMsg *redis.Message) *Message // todo: Message needs to be generic after 1.18
}

type ObjectDetectionParser struct {
	FileName string `json:"file_name"`
	Image []byte `json:"image"`
}

func (m ObjectDetectionParser) Parse(redisMsg *redis.Message) *Message  {
	var msg Message
	err := json.Unmarshal([]byte(redisMsg.Payload), &msg)
	if err != nil {
		log.Println("json conversation has been failed due to " + err.Error())
		return nil
	}
	return &msg
}
