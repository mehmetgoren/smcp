package rd

import (
	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
	"log"
)

type RedisSubPubOptions struct {
	*RedisOptions
	Channel string // "obj_detection"
}


type RedisListener interface{
	Listen(func(message *redis.Message))
}

func (r RedisSubPubOptions) Listen(onMessageReceived func(message *redis.Message)){
	pong, err := r.Client.Ping(context.Background()).Result()
	if err != nil {
		log.Println("ping has been failed, exiting now...")
		return
	}

	log.Println("ping: " + pong)
	log.Println("redis pubsub is listening")

	channel := r.Channel
	subscribe := r.Client.Subscribe(context.Background(), channel)
	subscriptions := subscribe.ChannelWithSubscriptions(context.Background(), 1)
	for {
		select {
		case sub := <-subscriptions:
			var message, isRedisMessage = sub.(*redis.Message)
			if !isRedisMessage {
				continue
			}
			go onMessageReceived(message)
		}
	}
}
