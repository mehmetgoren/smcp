package main

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
)

var BackgroundContext = context.Background()

func createRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:	  "localhost:6379",
		Password: "", // no password set
		DB:		  0,  // use default DB
	})
}

func main() {
	redisClient := createRedisClient()
	redisOptions := RedisOptions{BackgroundContext, redisClient}
	var rc RedisListener = RedisSubPubOptions{&redisOptions,  "obj_detection"}

	bi := TelegramBotOptions{Token: "1944447440:AAF8C0vJ2rjd__9CWT7PVcg9cON8QixdAMs"}
	si := SenderOptions{Id: 1608330992, Username:"ionian_gokalp", LanguageCode: "en"}
	telegramBotClient, botErr := CreateTelegramBot(&bi, &si)
	if botErr != nil{
		log.Println("telegram bot connection couldn't be created")
		return
	}

	//diskInfo := DiskInfoHandler{RootFolder: "/home/gokalp/Documents/shared_codes/resources/delete_later/"}
	var parser MessageParser = ObjectDetectionMessage {}
	var handler MessageHandler = TelegramHandler{&telegramBotClient} //TelegramAndDiskHandler{&diskInfo, &telegramBotClient}
	rc.Listen(func(message *redis.Message) {
		msg := parser.Parse(message)
		if msg == nil{
			log.Println("Message parsing returned nil")
			return
		}
		_, err := handler.Handle(msg)
		if err != nil {
			log.Println("An error occurred on handle: " + err.Error())
		}
	})
}

//defer func() {
//	fmt.Println("not today asshole")
//	if r := recover(); r!= nil {
//		log.Println("recovered from ", r)
//	}
//}()
//
//y := 01
//x := 12 / y
//
//panic("yo mf")
//
//print(x)

//defer func() {
//	fmt.Println(" not today asshole")
//	if r := recover(); r!= nil {
//		log.Println("recovered from ", r)
//	}
//}()

//startListen()