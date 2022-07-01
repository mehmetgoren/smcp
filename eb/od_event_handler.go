package eb

import (
	"github.com/go-redis/redis/v8"
	"log"
	"smcp/models"
	"smcp/reps"
	"smcp/utils"
)

type OdEventHandler struct {
	Ohr      *reps.OdHandlerRepository
	Notifier *NotifierPublisher
}

func (d *OdEventHandler) Handle(event *redis.Message) (interface{}, error) {
	defer utils.HandlePanic()

	var de = models.ObjectDetectionModel{}
	err := utils.DeserializeJson(event.Payload, &de)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	err = d.Ohr.Save(&de)

	if err == nil {
		go func() {
			err := d.Notifier.Publish(&event.Payload, ObjectDetection)
			if err != nil {
				log.Println(err.Error())
			}
		}()
	} else {
		log.Println(err.Error())
	}

	return nil, err
}
