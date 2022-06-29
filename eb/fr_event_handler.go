package eb

import (
	"github.com/go-redis/redis/v8"
	"log"
	"smcp/models"
	"smcp/reps"
	"smcp/utils"
)

type FrEventHandler struct {
	Fhr      *reps.FrHandlerRepository
	Notifier *NotifierPublisher
}

func (d *FrEventHandler) Handle(event *redis.Message) (interface{}, error) {
	defer utils.HandlePanic()

	var de = models.FaceRecognitionModel{}
	err := utils.DeserializeJson(event.Payload, &de)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	err = d.Fhr.Save(&de)
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
