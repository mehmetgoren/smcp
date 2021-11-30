package main

import (
	"fmt"
	"log"
	"os"
	"smcp/disk"
	"smcp/eb"
	"smcp/gdrive"
	"smcp/rd"
	"smcp/tb"
	"strconv"

	"github.com/go-redis/redis/v8"
)

func createRedisClient(host string, port int) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     host + ":" + strconv.Itoa(port),
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
	host := os.Getenv("REDIS_HOST")
	fmt.Println("Redis host: ", host)
	if len(host) == 0 {
		host = "127.0.0.1"
	}
	portStr := os.Getenv("REDIS_PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Println("An error occurred while converting Redis port value:" + err.Error())
		port = 6379
	}
	redisClient := createRedisClient(host, port)
	redisOptions := rd.RedisOptions{Client: redisClient}
	var rep = createRedisRepository(&redisOptions)

	heartbeat := createHeartbeat(&rep)
	go heartbeat.Start()

	handlerList := make([]eb.EventHandler, 0)
	var parser eb.ObjectDetectionParser
	var diskHandler = eb.DiskEventHandler{}
	diskHandler.FolderManager = &disk.FolderManager{SmartMachineFolderPath: "/go/src/smcp/images/"}
	diskHandler.FolderManager.Redis = redisClient
	handlerList = append(handlerList, &diskHandler)

	telegramBotClient, botErr := tb.CreateTelegramBot(&rep)
	if botErr != nil {
		log.Println("telegram bot connection couldn't be created, the operation is now exiting")
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
