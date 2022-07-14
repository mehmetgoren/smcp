package eb

import (
	"github.com/go-redis/redis/v8"
	"log"
	"smcp/data/cmn"
	"smcp/models"
	"smcp/utils"
)

type FrEventHandler struct {
	Factory  *cmn.Factory
	Notifier *NotifierPublisher
}

func (d *FrEventHandler) Handle(event *redis.Message) (interface{}, error) {
	defer utils.HandlePanic()

	var fr = &models.FaceRecognitionModel{}
	err := utils.DeserializeJson(event.Payload, fr)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	err = d.Factory.CreateRepository().FrSave(fr)
	if err == nil {
		go func() {
			err := d.Notifier.Publish(&event.Payload, FaceRecognition)
			if err != nil {
				log.Println(err.Error())
			}
		}()
	} else {
		log.Println(err.Error())
	}

	return nil, err
}
