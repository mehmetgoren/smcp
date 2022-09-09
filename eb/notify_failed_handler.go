package eb

import "github.com/go-redis/redis/v8"

type NotifyFailedHandler struct {
	Notifier *NotifierPublisher
}

func (n *NotifyFailedHandler) Handle(event *redis.Message) (interface{}, error) {
	err := n.Notifier.Publish(&event.Payload, NotifyFailed)
	return nil, err
}
