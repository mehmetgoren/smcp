package reps

import (
	"context"
	"github.com/go-redis/redis/v8"
	"smcp/models"
	"smcp/utils"
)

type OdQueueRepository struct {
	Connection *redis.Client
}

var key = "odseries"

func (d *OdQueueRepository) Add(payloadJson *string) error {
	_, err := d.Connection.RPush(context.Background(), key, *payloadJson).Result()
	return err
}

func (d *OdQueueRepository) PopAll() ([]*models.ObjectDetectionModel, error) {
	items := make([]*models.ObjectDetectionModel, 0)
	c := d.Connection
	for true {
		json, err := c.LPop(context.Background(), key).Result()
		if err != nil {
			if err.Error() == "redis: nil" {
				return items, nil
			}
			return nil, err
		}

		var item = &models.ObjectDetectionModel{}
		utils.DeserializeJson(json, item)
		items = append(items, item)
	}
	return items, nil
}
