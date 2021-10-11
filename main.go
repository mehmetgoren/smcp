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

func createRedisRepository(opts *RedisOptions) RedisRepository {
	return RedisRepository{opts, "users"}
}

func createHeartbeat(rep *RedisRepository) *HeartbeatClient{
	var heartbeat HeartbeatClient = HeartbeatClient{rep, 5}

	return &heartbeat
}

func main() {
	redisClient := createRedisClient()
	redisOptions := RedisOptions{BackgroundContext, redisClient}
	var rep = createRedisRepository(&redisOptions)

	heartbeat := createHeartbeat(&rep)
	go heartbeat.Start()

	var rc RedisListener = RedisSubPubOptions{&redisOptions,  "obj_detection"}
	token := "1944447440:AAF8C0vJ2rjd__9CWT7PVcg9cON8QixdAMs"
	telegramBotClient, botErr := CreateTelegramBot(token, &rep)
	if botErr != nil{
		log.Println("telegram bot connection couldn't be created")
		return
	}

	var parser ObjectDetectionParser
	//var handler = DiskInfoHandler{RootFolder: "/home/gokalp/Documents/shared_codes/object_detector/resources/delete_later/"}
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