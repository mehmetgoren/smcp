package eb

import (
	"github.com/go-redis/redis/v8"
	"log"
	"smcp/models"
	"smcp/reps"
	"smcp/utils"
)

type AlprEventHandler struct {
	Ahr      *reps.AlprHandlerRepository
	Notifier *NotifierPublisher
}

func (a *AlprEventHandler) Handle(event *redis.Message) (interface{}, error) {
	defer utils.HandlePanic()

	var ar = models.AlprResponse{}
	err := utils.DeserializeJson(event.Payload, &ar)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	err = a.Ahr.Save(&ar)
	if err == nil {
		go func() {
			err := a.Notifier.Publish(&event.Payload, PlateRecognition)
			if err != nil {
				log.Println(err.Error())
			}
		}()
	} else {
		log.Println(err.Error())
	}

	return nil, nil
}
