package eb

import (
	"github.com/go-redis/redis/v8"
	"log"
	"smcp/data/cmn"
	"smcp/models"
	"smcp/utils"
)

type AiEventHandler struct {
	Factory  *cmn.Factory
	Notifier *NotifierPublisher
}

func (d *AiEventHandler) Handle(event *redis.Message) (interface{}, error) {
	defer utils.HandlePanic()

	var adm = models.AiDetectionModel{}
	err := utils.DeserializeJson(event.Payload, &adm)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	err = d.Factory.CreateRepository().AiSave(&adm)

	if err == nil {
		go func() {
			err := d.Notifier.Publish(&event.Payload, AiDetection)
			if err != nil {
				log.Println(err.Error())
			}
		}()
	} else {
		log.Println(err.Error())
	}

	return nil, err
}
