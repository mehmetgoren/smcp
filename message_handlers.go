package main

import (
	"encoding/base64"
	tb "gopkg.in/tucnak/telebot.v2"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type Message struct{
	FileName string
	Base64Image *string
}

type MessageHandler interface{
	Handle(message *Message) (interface{}, error)
}

type DiskInfoHandler struct {
	RootFolder string // "/home/gokalp/Documents/shared_codes/resources/delete_later/"
}

func (d DiskInfoHandler) Handle(message *Message)  (interface{}, error) {
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(*message.Base64Image))
	defer ioutil.NopCloser(reader)
	fileBytes, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Println("reading base64 message error: " + err.Error())
		return nil, err
	}

	fileFullPath := d.RootFolder + message.FileName
	file, err := os.Create(d.RootFolder + message.FileName)
	if err != nil {
		log.Println("creating the image file error: " + err.Error())
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Println("closing the image file error: " + err.Error())
		}
	}(file)

	if _, err := file.Write(fileBytes); err != nil {
		log.Println("writing the image file error: " + err.Error())
		return nil, err
	}
	if err := file.Sync(); err != nil {
		log.Println("syncing the image file error: " + err.Error())
		return nil, err
	}

	log.Println("image saved successfully as " + message.FileName)

	return fileFullPath, nil
}

///////////////////////////////////////////////////////telegram \\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\
type TelegramAndDiskHandler struct {
	*DiskInfoHandler
	*TelegramBotClient
}

func (td TelegramAndDiskHandler) Handle(message *Message) (interface{}, error)  {
	fileFullPath, err := td.DiskInfoHandler.Handle(message)
	if err != nil {
		log.Println("writing image file error: " + err.Error())
		return nil, err
	}

	tbFile := tb.FromDisk(fileFullPath.(string))
	tbFile.UniqueID = message.FileName
	tbPhoto:=&tb.Photo{File: tbFile, Caption: message.FileName}

	msg, sendErr := td.Bot.Send(td.Sender, tbPhoto)
	if sendErr != nil{
		log.Println("send error is " + sendErr.Error())
		return nil, sendErr
	}

	log.Println("image send successfully as " + message.FileName + " the message is " + msg.Text)

	return msg, nil
}

///////////////////////////////////////////////////////telegram \\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\
type TelegramHandler struct {
	*TelegramBotClient
}

func (t TelegramHandler) Handle(message *Message) (interface{}, error)   {
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(*message.Base64Image))
	defer ioutil.NopCloser(reader)

	tbFile := tb.FromReader(reader)
	//tbFile.FileID = "aaa.jpg"
	//tbFile.FileLocal = "aaa.jpg"
	//tbFile.FilePath = "aaa.jpg"
	tbFile.UniqueID = message.FileName
	tbPhoto:=&tb.Photo{File: tbFile, Caption: message.FileName}

	msg, sendErr := t.Bot.Send(t.Sender, tbPhoto)
	if sendErr != nil{
		log.Println("send error is " + sendErr.Error())
		return nil, sendErr
	}

	log.Println("image send successfully as " + message.FileName + " the message is " + msg.Text)

	return msg, nil
}

