package main

import (
	"fmt"
	"log"
	"os"
	"smcp/disk"
	"smcp/eb"
	"smcp/rd"
	"strconv"

	"github.com/go-redis/redis/v8"
)

const (
	MAIN     = 0
	SERVICE  = 1
	SOURCES  = 2
	EVENTBUS = 3
)

func createRedisClient(host string, port int, db int) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     host + ":" + strconv.Itoa(port),
		Password: "", // no password set
		DB:       db, // use default DB (0)
	})
}

//func createRedisRepository(opts *rd.RedisOptions) rd.RedisRepository {
//	return rd.RedisRepository{RedisOptions: opts}
//}

func createHeartbeatRepository(opts *rd.RedisOptions) *rd.HeartbeatRepository {
	var heartbeatRepository = rd.HeartbeatRepository{RedisOptions: opts, TimeSecond: 10}

	return &heartbeatRepository
}

func createPidRepository(host string, port int) *rd.PidRepository {
	client := createRedisClient(host, port, SERVICE) // 1 is SERVICE db
	opts := rd.RedisOptions{Client: client}
	var pidRepository = rd.PidRepository{RedisOptions: &opts}

	return &pidRepository
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
	redisClientMain := createRedisClient(host, port, MAIN)
	redisOptions := rd.RedisOptions{Client: redisClientMain}
	//var rep = createRedisRepository(&redisOptions)

	heartbeat := createHeartbeatRepository(&redisOptions)
	go heartbeat.Start()

	pid := createPidRepository(host, port)
	go func() {
		_, err := pid.Add()
		if err != nil {
			log.Println("An error occurred while registering process id, error is:" + err.Error())
		}
	}()

	handlerList := make([]eb.EventHandler, 0)
	var parser eb.ObjectDetectionParser
	var diskHandler = eb.DiskEventHandler{}
	diskHandler.FolderManager = &disk.FolderManager{SmartMachineFolderPath: "/home/gokalp/Pictures/detected/"}
	diskHandler.FolderManager.Redis = redisClientMain
	handlerList = append(handlerList, &diskHandler)

	//telegramBotClient, botErr := tb.CreateTelegramBot(&rep)
	//if botErr != nil {
	//	log.Println("telegram bot connection couldn't be created, the operation is now exiting")
	//	return
	//}
	//var tbHandler eb.EventHandler = &eb.TelegramEventHandler{TelegramBotClient: &telegramBotClient}
	//handlerList = append(handlerList, tbHandler)

	//var fm = &gdrive.FolderManager{}
	//fm.Redis = redisClient
	//fm.Gdrive = &gdrive.GdriveClient{}
	//fm.Gdrive.Repository = &rep
	//var gHandler = &eb.GdriveEventHandler{FolderManager: fm}
	//handlerList = append(handlerList, gHandler)

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
