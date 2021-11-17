package main

import (
	"log"
	"smcp/disk"
	"smcp/eb"
	"smcp/gdrive"
	"smcp/rd"
	"smcp/tb"

	"github.com/go-redis/redis/v8"
)

func createRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

func createRedisRepository(opts *rd.RedisOptions) rd.RedisRepository {
	return rd.RedisRepository{RedisOptions: opts}
}

func createHeartbeat(rep *rd.RedisRepository) *rd.HeartbeatClient {
	var heartbeat = rd.HeartbeatClient{Repository: rep, TimeSecond: 5}

	return &heartbeat
}

func main() {
	redisClient := createRedisClient()
	redisOptions := rd.RedisOptions{Client: redisClient}
	var rep = createRedisRepository(&redisOptions)

	heartbeat := createHeartbeat(&rep)
	go heartbeat.Start()

	handlerList := make([]eb.EventHandler, 0)
	var parser eb.ObjectDetectionParser
	var diskHandler = eb.DiskEventHandler{}
	diskHandler.FolderManager = &disk.FolderManager{SmartMachineFolderPath: "/home/gokalp/Pictures/detected/"}
	diskHandler.FolderManager.Redis = redisClient
	handlerList = append(handlerList, &diskHandler)

	token := "1944447440:AAF8C0vJ2rjd__9CWT7PVcg9cON8QixdAMs"
	telegramBotClient, botErr := tb.CreateTelegramBot(token, &rep)
	if botErr != nil {
		log.Println("telegram bot connection couldn't be created")
		return
	}
	var tbHandler eb.EventHandler = &eb.TelegramEventHandler{TelegramBotClient: &telegramBotClient}
	handlerList = append(handlerList, tbHandler)

	var fm = &gdrive.FolderManager{}
	fm.Redis = redisClient
	fm.Gdrive = &gdrive.GdriveClient{}
	fm.Gdrive.Repository = &rep
	var gHandler = &eb.GdriveEventHandler{FolderManager: fm}
	handlerList = append(handlerList, gHandler)

	var handler = eb.ComboEventHandler{
		EventHandlers: handlerList,
	}

	var rc rd.RedisListener = rd.RedisSubPubOptions{RedisOptions: &redisOptions, Channel: "detect"}
	rc.Listen(func(message *redis.Message) {
		msg := parser.Parse(message)
		if msg == nil {
			log.Println("Message parsing returned nil")
			return
		}
		_, err := handler.Handle(msg)
		if err != nil {
			log.Println("An error occurred on handle: " + err.Error())
		}
	})
}
