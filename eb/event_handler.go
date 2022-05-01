package eb

import "github.com/go-redis/redis/v8"

// EventHandler desperately needs generics
type EventHandler interface {
	Handle(event *redis.Message) (interface{}, error)
}
