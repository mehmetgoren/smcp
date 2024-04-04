package reps

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"log"
	"smcp/models"
	"smcp/utils"
)

type AiClipQueueRepository struct {
	Connection *redis.Client
}

var key = "aiclipseries"

// todo: remove Unmarshal and Marshal operations if it is not necessary
func (a *AiClipQueueRepository) Add(payloadJson *string) error {
	queueModel := &models.AiClipQueueModel{}
	var err error

	ai := &models.AiDetectionModel{}
	err = json.Unmarshal([]byte(*payloadJson), ai)
	if err != nil {
		log.Println("an error occurred while deserializing an ai model")
		return err
	}
	queueModel.Ai = ai
	modelJson, err := json.Marshal(queueModel)
	if err != nil {
		log.Println("an error occurred while ")
		return err
	}
	_, err = a.Connection.RPush(context.Background(), key, string(modelJson)).Result()
	return err
}

func (a *AiClipQueueRepository) PopAll() ([]*models.AiClipQueueModel, error) {
	items := make([]*models.AiClipQueueModel, 0)
	c := a.Connection
	for true {
		json, err := c.LPop(context.Background(), key).Result()
		if err != nil {
			if err.Error() == "redis: nil" {
				return items, nil
			}
			return nil, err
		}

		var queueModel = &models.AiClipQueueModel{}
		utils.DeserializeJson(json, queueModel)
		items = append(items, queueModel)
	}
	return items, nil
}
