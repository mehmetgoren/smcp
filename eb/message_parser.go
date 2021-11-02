package eb

import (
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"log"
)

type ObjectDetectionParser struct {
}

func (m ObjectDetectionParser) Parse(redisMsg *redis.Message) *Message {
	var msg Message
	err := json.Unmarshal([]byte(redisMsg.Payload), &msg)
	if err != nil {
		log.Println("json conversation has been failed due to " + err.Error())
		return nil
	}
	return &msg
}
