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

func (a *AiClipQueueRepository) Add(aiType int, payloadJson *string) error {
	queueModel := &models.AiClipQueueModel{AiType: aiType}
	var err error
	switch aiType {
	case models.Od:
		od := &models.ObjectDetectionModel{}
		err = json.Unmarshal([]byte(*payloadJson), od)
		if err != nil {
			log.Println("an error occurred while deserializing an od model")
			return err
		}
		queueModel.Od = od
		break
	case models.Fr:
		fr := &models.FaceRecognitionModel{}
		err = json.Unmarshal([]byte(*payloadJson), fr)
		if err != nil {
			log.Println("an error occurred while deserializing an fr model")
			return err
		}
		queueModel.Fr = fr
		break
	case models.Alpr:
		alpr := &models.AlprResponse{}
		err = json.Unmarshal([]byte(*payloadJson), alpr)
		if err != nil {
			log.Println("an error occurred while deserializing an alpr model")
			return err
		}
		queueModel.Alpr = alpr
		break
	}
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
