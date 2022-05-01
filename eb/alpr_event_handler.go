package eb

import (
	"github.com/go-redis/redis/v8"
	"smcp/models"
	"smcp/reps"
	"smcp/utils"
)

type AlprEventHandler struct {
	Ahr *reps.AlprHandlerRepository
}

func (a *AlprEventHandler) Handle(event *redis.Message) (interface{}, error) {
	defer utils.HandlePanic()

	var ar = models.AlprResponse{}
	utils.DeserializeJson(event.Payload, &ar)
	a.Ahr.Save(&ar)

	return nil, nil
}
