package eb

import (
	"encoding/base64"
	"fmt"
	"github.com/go-redis/redis/v8"
	"io/ioutil"
	"log"
	"runtime/debug"
	"smcp/disk"
	"smcp/gdrive"
	"smcp/tb"
	"smcp/utils"
	"strings"

	"gopkg.in/tucnak/telebot.v2"
)

func handlePanic() {
	if r := recover(); r != nil {
		fmt.Println("RECOVER", r)
		debug.PrintStack()
	}
}

type DetectedImage struct {
	FileName    string  `json:"file_name"`
	Base64Image *string `json:"base64_image"`
}

// EventHandler needs desperately generics
type EventHandler interface {
	Handle(event *redis.Message) (interface{}, error)
}

type DiskEventHandler struct {
	FolderManager *disk.FolderManager
}

func (d *DiskEventHandler) Handle(event *redis.Message) (interface{}, error) {
	defer handlePanic()

	var de = DetectedImage{}
	utils.DeserializeJson(event.Payload, &de)

	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(*de.Base64Image))
	defer ioutil.NopCloser(reader)
	fileBytes, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Println("DiskEventHandler: Reading base64 message error: " + err.Error())
		return nil, err
	}

	fileFullPath, err := d.FolderManager.SaveFile(de.FileName, fileBytes)
	if err != nil {
		log.Println("DiskEventHandler: Saving base64 file error: " + err.Error())
		return nil, err
	}

	log.Println("DiskEventHandler: image saved successfully as " + de.FileName)

	return fileFullPath, nil
}

type TelegramEventHandler struct {
	*tb.TelegramBotClient
}

func (t *TelegramEventHandler) Handle(event *redis.Message) (interface{}, error) {
	defer handlePanic()

	var de = DetectedImage{}
	utils.DeserializeJson(event.Payload, &de)

	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(*de.Base64Image))
	defer ioutil.NopCloser(reader)

	tbFile := telebot.FromReader(reader)
	tbFile.UniqueID = de.FileName
	tbPhoto := &telebot.Photo{File: tbFile, Caption: de.FileName}

	users := t.Repository.GetAllUsers()
	for _, user := range users {
		msg, sendErr := t.Bot.Send(user, tbPhoto)
		if sendErr != nil {
			log.Println("TelegramEventHandler: Send error for " + msg.Caption + ". The error is " + sendErr.Error())
			return nil, sendErr
		}
	}

	log.Println("TelegramEventHandler: image send successfully as " + de.FileName + " the message is " + de.FileName)

	return nil, nil
}

type GdriveEventHandler struct {
	*gdrive.FolderManager
}

func (g *GdriveEventHandler) Handle(event *redis.Message) (interface{}, error) {
	defer handlePanic()

	var de = DetectedImage{}
	utils.DeserializeJson(event.Payload, &de)

	file, err := g.UploadImage(de.FileName, de.Base64Image)
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
