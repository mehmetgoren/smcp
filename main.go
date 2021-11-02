package main

import (
	"github.com/go-redis/redis/v8"
	"log"
	"smcp/eb"
	"smcp/gdrive"
	"smcp/rd"
	"smcp/tb"
)

func createRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

func createRedisRepository(opts *rd.RedisOptions) rd.RedisRepository {
	return rd.RedisRepository{RedisOptions: opts, Key: "users"}
}

func createHeartbeat(rep *rd.RedisRepository) *rd.HeartbeatClient {
	var heartbeat rd.HeartbeatClient = rd.HeartbeatClient{Repository: rep, TimeSecond: 5}

	return &heartbeat
}

func main() {
	redisClient := createRedisClient()
	redisOptions := rd.RedisOptions{Client: redisClient}
	var rep = createRedisRepository(&redisOptions)

	heartbeat := createHeartbeat(&rep)
	go heartbeat.Start()

	var rc rd.RedisListener = rd.RedisSubPubOptions{RedisOptions: &redisOptions, Channel: "obj_detection"}
	token := "****"
	telegramBotClient, botErr := tb.CreateTelegramBot(token, &rep)
	if botErr != nil {
		log.Println("telegram bot connection couldn't be created")
		return
	}

	handlerList := make([]eb.EventHandler, 0)
	var parser eb.ObjectDetectionParser
	var diskHandler = eb.DiskEventHandler{RootFolder: "/home/gokalp/Documents/shared_codes/object_detector/resources/delete_later/"}
	handlerList = append(handlerList, &diskHandler)
	var tbHandler eb.EventHandler = &eb.TelegramEventHandler{TelegramBotClient: &telegramBotClient} //TelegramAndDiskHandler{&diskInfo, &telegramBotClient}
	handlerList = append(handlerList, tbHandler)

	var fm = &gdrive.FolderManager{}
	fm.Redis = redisClient
	fm.Gdrive = &gdrive.GdriveClient{}
	var gHandler = &eb.GdriveEventHandler{FolderManager: fm}
	handlerList = append(handlerList, gHandler)

	var handler = eb.ComboEventHandler{
		EventHandlers: handlerList,
	}

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
