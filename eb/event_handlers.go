package eb

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime/debug"
	"smcp/gdrive"
	"smcp/tb"
	"strings"

	"gopkg.in/tucnak/telebot.v2"
)

func handlePanic() {
	if r := recover(); r != nil {
		fmt.Println("RECOVER", r)
		debug.PrintStack()
	}
}

type Message struct {
	FileName    string  `json:"file_name"`
	Base64Image *string `json:"base64_image"`
}

type EventHandler interface {
	Handle(message *Message) (interface{}, error)
}

type DiskEventHandler struct {
	RootFolder string // "/home/gokalp/Documents/shared_codes/resources/delete_later/"
}

func (d *DiskEventHandler) Handle(message *Message) (interface{}, error) {
	defer handlePanic()
	if message == nil {
		log.Println("DiskEventHandler: Message passed as null, disk operation won't be executed")
		return nil, nil
	}
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(*message.Base64Image))
	defer ioutil.NopCloser(reader)
	fileBytes, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Println("DiskEventHandler: Reading base64 message error: " + err.Error())
		return nil, err
	}

	fileFullPath := d.RootFolder + message.FileName
	file, err := os.Create(d.RootFolder + message.FileName)
	if err != nil {
		log.Println("DiskEventHandler: Creating the image file error: " + err.Error())
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Println("DiskEventHandler: Closing the image file error: " + err.Error())
		}
	}(file)

	if _, err := file.Write(fileBytes); err != nil {
		log.Println("DiskEventHandler: Writing the image file error: " + err.Error())
		return nil, err
	}
	if err := file.Sync(); err != nil {
		log.Println("DiskEventHandler: Syncing the image file error: " + err.Error())
		return nil, err
	}

	log.Println("DiskEventHandler: image saved successfully as " + message.FileName)

	return fileFullPath, nil
}

type TelegramEventHandler struct {
	*tb.TelegramBotClient
}

func (t *TelegramEventHandler) Handle(message *Message) (interface{}, error) {
	defer handlePanic()
	if message == nil {
		log.Println("TelegramEventHandler: Message passed as null, telegram operation won't be executed")
		return nil, nil
	}

	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(*message.Base64Image))
	defer ioutil.NopCloser(reader)

	tbFile := telebot.FromReader(reader)
	tbFile.UniqueID = message.FileName
	tbPhoto := &telebot.Photo{File: tbFile, Caption: message.FileName}

	users := t.Repository.GetAllUsers()
	for _, user := range users {
		msg, sendErr := t.Bot.Send(user, tbPhoto)
		if sendErr != nil {
			log.Println("TelegramEventHandler: Send error for " + msg.Caption + ". The error is " + sendErr.Error())
			return nil, sendErr
		}
	}

	log.Println("TelegramEventHandler: image send successfully as " + message.FileName + " the message is " + message.FileName)

	return nil, nil
}

type GdriveEventHandler struct {
	*gdrive.FolderManager
}

func (g *GdriveEventHandler) Handle(message *Message) (interface{}, error) {
	defer handlePanic()
	if message == nil {
		log.Println("GdriveEventHandler: Message passed as null, google drive operation won't be executed")
		return nil, nil
	}
	file, err := g.UploadImage(message.FileName, message.Base64Image)
	if err != nil {
		log.Println("GdriveEventHandler: An error occurred during the handling image uploading to google drive")
		return nil, err
	}

	return file, nil
}

type ComboEventHandler struct {
	EventHandlers []EventHandler
}

func (c *ComboEventHandler) Handle(message *Message) (interface{}, error) {
	for _, ev := range c.EventHandlers {
		go func(eventHandler EventHandler, msg Message) {
			_, err := eventHandler.Handle(&msg)
			if err != nil {
				log.Println("ComboEventHandler: An error occurred during the combo event handling")
			}
		}(ev, *message)
	}

	return nil, nil
}
