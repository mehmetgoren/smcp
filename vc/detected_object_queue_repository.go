package vc

import (
	"context"
	"github.com/go-redis/redis/v8"
	"smcp/models"
	"smcp/utils"
)

type DetectedObjectQueueRepository struct {
	Connection *redis.Client
}

var key = "detectionseries"

func (d *DetectedObjectQueueRepository) Add(payloadJson *string) error {
	_, err := d.Connection.RPush(context.Background(), key, *payloadJson).Result()
	return err
}

func (d *DetectedObjectQueueRepository) PopAll() ([]*models.DetectedImage, error) {
	items := make([]*models.DetectedImage, 0)
	c := d.Connection
	for true {
		json, err := c.LPop(context.Background(), key).Result()
		if err != nil {
			if err.Error() == "redis: nil" {
				return items, nil
			}
			return nil, err
		}

		var item = &models.DetectedImage{}
		utils.DeserializeJson(json, item)
		items = append(items, item)
	}
	return items, nil
}
