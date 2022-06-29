package eb

import (
	"encoding/base64"
	"github.com/go-redis/redis/v8"
	"io/ioutil"
	"log"
	"smcp/gdrive"
	"smcp/models"
	"smcp/reps"
	"smcp/tb"
	"smcp/utils"
	"strings"

	"gopkg.in/tucnak/telebot.v2"
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

type OdAiClipEventHandler struct {
	Connection *redis.Client
}

func (v *OdAiClipEventHandler) Handle(event *redis.Message) (interface{}, error) {
	defer utils.HandlePanic()
	rep := reps.OdQueueRepository{Connection: v.Connection}
	rep.Add(&event.Payload)

	return true, nil
}

type OdTelegramEventHandler struct {
	*tb.TelegramBotClient
}

func (t *OdTelegramEventHandler) Handle(event *redis.Message) (interface{}, error) {
	defer utils.HandlePanic()

	var de = models.ObjectDetectionModel{}
	utils.DeserializeJson(event.Payload, &de)

	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(de.Base64Image))
	defer ioutil.NopCloser(reader)

	fileName := de.CreateFileName()
	tbFile := telebot.FromReader(reader)
	tbFile.UniqueID = fileName
	tbPhoto := &telebot.Photo{File: tbFile, Caption: fileName}

	users := t.Repository.GetAllUsers()
	for _, user := range users {
		msg, sendErr := t.Bot.Send(user, tbPhoto)
		if sendErr != nil {
			log.Println("TelegramEventHandler: Send error for " + msg.Caption + ". The error is " + sendErr.Error())
			return nil, sendErr
		}
	}

	log.Println("TelegramEventHandler: image send successfully as " + fileName + " the message is " + fileName)

	return nil, nil
}

type OdGdriveEventHandler struct {
	*gdrive.FolderManager
}

func (g *OdGdriveEventHandler) Handle(event *redis.Message) (interface{}, error) {
	defer utils.HandlePanic()

	var de = models.ObjectDetectionModel{}
	utils.DeserializeJson(event.Payload, &de)

	file, err := g.UploadImage(de.CreateFileName(), &de.Base64Image)
	if err != nil {
		log.Println("GdriveEventHandler: An error occurred during the handling image uploading to google drive")
		return nil, err
	}

	return file, nil
}

type ComboEventHandler struct {
	EventHandlers []EventHandler
}

func (c *ComboEventHandler) Handle(event *redis.Message) (interface{}, error) {
	for _, ev := range c.EventHandlers {
		go func(eventHandler EventHandler, evt *redis.Message) {
			_, err := eventHandler.Handle(evt)
			if err != nil {
				log.Println("ComboEventHandler: An error occurred during the combo event handling")
			}
		}(ev, event)
	}

	return nil, nil
}
