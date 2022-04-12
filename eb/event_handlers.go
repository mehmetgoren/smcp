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

// EventHandler desperately needs generics
type EventHandler interface {
	Handle(event *redis.Message) (interface{}, error)
}

type DiskEventHandler struct {
	Ohr *reps.OdHandlerRepository
}

func (d *DiskEventHandler) Handle(event *redis.Message) (interface{}, error) {
	defer utils.HandlePanic()

	var de = models.ObjectDetectionModel{}
	utils.DeserializeJson(event.Payload, &de)
	d.Ohr.Save(&de)

	return nil, nil
}

type VideoClipsEventHandler struct {
	Connection *redis.Client
}

func (v *VideoClipsEventHandler) Handle(event *redis.Message) (interface{}, error) {
	defer utils.HandlePanic()
	rep := reps.ObjectDetectionQueueRepository{Connection: v.Connection}
	rep.Add(&event.Payload)

	return true, nil
}

type TelegramEventHandler struct {
	*tb.TelegramBotClient
}

func (t *TelegramEventHandler) Handle(event *redis.Message) (interface{}, error) {
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

type GdriveEventHandler struct {
	*gdrive.FolderManager
}

func (g *GdriveEventHandler) Handle(event *redis.Message) (interface{}, error) {
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
