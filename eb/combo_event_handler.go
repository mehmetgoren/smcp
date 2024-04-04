package eb

import (
	"github.com/go-redis/redis/v8"
	"log"
)

type ComboEventHandler struct {
	EventHandlers []EventHandler
}

func (c *ComboEventHandler) Handle(event *redis.Message) (interface{}, error) {
	for _, ev := range c.EventHandlers {
		go func(eventHandler EventHandler, evt *redis.Message) {
			_, err := eventHandler.Handle(evt)
			if err != nil {
				log.Println("ComboEventHandler: An error occurred during the combo event handling")
				log.Println(err)
			}
		}(ev, event)
	}

	return nil, nil
}
