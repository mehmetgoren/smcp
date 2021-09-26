package main

import (
	"context"
	"encoding/base64"
	"github.com/go-redis/redis/v8"
	tb "gopkg.in/tucnak/telebot.v2"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var ctx = context.Background()
var rootFolder = "/home/gokalp/Documents/shared_codes/resources/delete_later/"

func startListen() {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	pong, err := client.Ping(ctx).Result()
	if err != nil {
		log.Println("ping has been failed, exiting now...")
		return
	}
	log.Println("ping: " + pong)
	log.Println("redis pubsub is listening")

	bi := TelegramBotOptions{Token: "1944447440:AAF8C0vJ2rjd__9CWT7PVcg9cON8QixdAMs"}
	si := SenderOptions{Id: 1608330992, Username:"ionian_gokalp", LanguageCode: "en"}
	c, botErr := CreateTelegramBot(&bi, &si)
	if botErr != nil{
		log.Println("telegram bot connection couldn't be created")
		return
	}

	channel := "obj_detection"
	subscribe := client.Subscribe(ctx, channel)
	subscriptions := subscribe.ChannelWithSubscriptions(ctx, 1)
	for {
		select {
		case sub := <-subscriptions:
			var message, isRedisMessage = sub.(*redis.Message)
			if !isRedisMessage {
				continue
			}
			arr := strings.Split(message.Payload, "■")
			if len(arr) == 2 {
				fileName := arr[0]
				base64Image := arr[1]
				saveMessage(c.Bot, c.Sender, fileName, &base64Image)
			}
			log.Println(arr[0])
		}
	}
}

func saveMessage(bot *tb.Bot, sender *tb.User,  fileName string, base64Image *string) {
	var reader io.Reader = nil
	var file *os.File = nil
	defer func() {
		recoverMe := func(){
			if r := recover(); r!= nil {
				log.Println("recovered from ", r)
			}else{
				log.Println("no panic so far")
			}
		}
		if reader != nil{
			ioutil.NopCloser(reader)
			log.Println("ioutil.NopCloser was closed")
		}
		if file != nil{
			err := file.Close()
			if err == nil{
				log.Println("file was closed")
			}else{
				log.Println("an error occured dur,ng the closing a file " + err.Error())
			}
		}
		defer func() {
			recoverMe()
		}()
		recoverMe()
	}()

	reader = base64.NewDecoder(base64.StdEncoding, strings.NewReader(*base64Image))
	fileBytes, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Println("reading base64 message error is " + err.Error())
		return
	}

	file, err = os.Create(rootFolder + fileName)
	if err != nil {
		log.Println("creating the image file error is " + err.Error())
		return
	}

	if _, err := file.Write(fileBytes); err != nil {
		log.Println("writing the image file error is " + err.Error())
		return
	}
	if err := file.Sync(); err != nil {
		log.Println("syncing the image file error is " + err.Error())
		return
	}

	tbFile := tb.FromReader(reader)//tb.FromReader(reader)//tb.FromDisk(file.Name())
	tbFile.FileLocal = file.Name()
	//tbFile.FileID = fileName
	tbFile.UniqueID = fileName
	tbPhoto:=&tb.Photo{File: tbFile, Caption: fileName}
	//bx := tb.Bot{}
	//bx.Raw()

	msg, sendErr := bot.Send(sender, tbPhoto)
	if sendErr != nil{
		log.Println("send error is " + sendErr.Error())
		return
	}

	log.Println("image saved as " + fileName + " the message is " + msg.Text)
}
