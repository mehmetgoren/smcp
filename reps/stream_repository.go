package reps

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
	"smcp/models"
)

type StreamRepository struct {
	Connection *redis.Client
}

var redisKeyStream = "streams:"

func (r *StreamRepository) GetAll() ([]*models.StreamModel, error) {
	conn := r.Connection
	keys, err := conn.Keys(context.Background(), redisKeyStream+"*").Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			emptyList := make([]*models.StreamModel, 0)
			conn.Set(context.Background(), redisKeyStream, emptyList, 0)
			return emptyList, nil
		} else {
			log.Println("Error getting all stream from redis: ", err)
			return nil, err
		}
	}
	list := make([]*models.StreamModel, 0, 5)
	for _, key := range keys {
		var p models.StreamModel
		err := conn.HGetAll(context.Background(), key).Scan(&p)
		if err != nil {
			log.Println("Error getting stream from redis: ", err)
			return nil, err
		}
		list = append(list, &p)
	}
	return list, nil
}
