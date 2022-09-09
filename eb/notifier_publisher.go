package eb

import (
	"encoding/base64"
	"encoding/json"
	"github.com/go-redis/redis/v8"
)

const (
	ObjectDetection  = 0
	FaceRecognition  = 1
	PlateRecognition = 2
	NotifyFailed     = 3
)

type NotifierInfo struct {
	Base64Object string `json:"base_64_object"`
	Type         int    `json:"type"`
}

func (n NotifierInfo) MarshalBinary() ([]byte, error) {
	return json.Marshal(n)
}

type NotifierPublisher struct {
	PubSubConnection *redis.Client
	eb               *EventBus
}

func (n *NotifierPublisher) Publish(payload *string, eventType int) error {
	if n.eb == nil {
		n.eb = &EventBus{PubSubConnection: n.PubSubConnection, Channel: "notifier"}
	}
	b64 := base64.StdEncoding.EncodeToString([]byte(*payload))
	err := n.eb.Publish(NotifierInfo{Base64Object: b64, Type: eventType})
	return err
}
