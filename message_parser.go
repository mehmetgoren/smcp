package main

import (
	"github.com/go-redis/redis/v8"
	"log"
	"strings"
)

// todo:need generic, changes implementation after 1.18
type MessageParser interface {
	Parse(redisMsg *redis.Message) *Message // todo: Message needs to be generic after 1.18
}

type ObjectDetectionMessage struct {
	Separator string
}

func (m ObjectDetectionMessage) Parse(redisMsg *redis.Message) *Message  {
	if len(m.Separator) == 0{
		m.Separator = "■"
	}
	arr := strings.Split(redisMsg.Payload, m.Separator)
	if len(arr) == 2 {
		fileName := arr[0]
		base64Image := arr[1]
		log.Println(arr[0])

		return &Message{FileName: fileName, Base64Image: &base64Image}
	}

	return nil
}
