package main

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"os"
	"smcp/disk"
	"smcp/eb"
	"smcp/reps"
	"smcp/utils"
	"smcp/vc"
	"strconv"
)

const (
	MAIN     = 0
	RQ       = 1
	EVENTBUS = 15
)

func createRedisClient(host string, port int, db int) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     host + ":" + strconv.Itoa(port),
		Password: "", // no password set
		DB:       db, // use default DB (0)
	})
}

func createHeartbeatRepository(client *redis.Client, serviceName string) *reps.HeartbeatRepository {
	var heartbeatRepository = reps.HeartbeatRepository{Client: client, TimeSecond: 10, ServiceName: serviceName}

	return &heartbeatRepository
}

func createServiceRepository(client *redis.Client) *reps.ServiceRepository {
	var pidRepository = reps.ServiceRepository{Client: client}

	return &pidRepository
}

func main() {
	defer utils.HandlePanic()

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
	redisClient := createRedisClient(host, port, MAIN)
	//var rep = createRedisRepository(&redisOptions)

	serviceName := "cloud_integration_service"
	heartbeat := createHeartbeatRepository(redisClient, serviceName)
	go heartbeat.Start()

	serviceRepository := createServiceRepository(redisClient)
	go func() {
		_, err := serviceRepository.Add(serviceName)
		if err != nil {
			log.Println("An error occurred while registering process id, error is:" + err.Error())
		}
	}()

	var configRep = reps.ConfigRepository{Connection: redisClient}
	config, _ := configRep.GetConfig()
	handlerList := make([]eb.EventHandler, 0)
	var diskHandler = eb.DiskEventHandler{}

	diskHandler.FolderManager = &disk.FolderManager{RootFolderPath: config.AiConfig.DetectedFolder}
	diskHandler.FolderManager.Redis = redisClient
	handlerList = append(handlerList, &diskHandler)

	//detection series handler
	var dsHandler = eb.VideoClipsEventHandler{Connection: redisClient}
	handlerList = append(handlerList, &dsHandler)

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

	//starts video clips processor
	dsRep := vc.DetectedObjectQueueRepository{Connection: redisClient}
	streamRep := reps.StreamRepository{Connection: redisClient}
	vcp := vc.VideoClipProcessor{Config: config, DoRep: &dsRep, StreamRep: &streamRep}
	go vcp.Start()
	//ends video clips processor

	var handler = &eb.ComboEventHandler{
		EventHandlers: handlerList,
	}

	var e = eb.EventBus{PubSubConnection: createRedisClient(host, port, EVENTBUS), Channel: "detect_service"}
	e.Subscribe(handler)
}
