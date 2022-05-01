package main

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"os"
	"smcp/eb"
	"smcp/models"
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

func createHeartbeatRepository(client *redis.Client, serviceName string, config *models.Config) *reps.HeartbeatRepository {
	var heartbeatRepository = reps.HeartbeatRepository{Client: client, TimeSecond: int64(config.General.HeartbeatInterval), ServiceName: serviceName}

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
	utils.SetPool(redisClient)

	var configRep = reps.ConfigRepository{Connection: redisClient}
	config, _ := configRep.GetConfig()

	serviceName := "cloud_integration_service"
	heartbeat := createHeartbeatRepository(redisClient, serviceName, config)
	go heartbeat.Start()

	serviceRepository := createServiceRepository(redisClient)
	go func() {
		_, err := serviceRepository.Add(serviceName)
		if err != nil {
			log.Println("An error occurred while registering process id, error is:" + err.Error())
		}
	}()

	go listenOdEventHandlers(redisClient, config, host, port)
	go listenAlprEventHandler(config, host, port)
	listenFrEventHandler(config, host, port)
}

func listenOdEventHandlers(mainConn *redis.Client, config *models.Config, host string, port int) {
	handlerList := make([]eb.EventHandler, 0)
	ohr := &reps.OdHandlerRepository{Config: config}
	var diskHandler = eb.OdDiskEventHandler{Ohr: ohr}
	handlerList = append(handlerList, &diskHandler)

	//detection series handler
	var vch = eb.OdVideoClipsEventHandler{Connection: mainConn}
	handlerList = append(handlerList, &vch)

	//telegramBotClient, botErr := tb.CreateTelegramBot(&rep)
	//if botErr != nil {
	//	log.Println("telegram bot connection couldn't be created, the operation is now exiting")
	//	return
	//}
	//var tbHandler eb.EventHandler = &eb.OdTelegramEventHandler{TelegramBotClient: &telegramBotClient}
	//handlerList = append(handlerList, tbHandler)

	//var fm = &gdrive.FolderManager{}
	//fm.Redis = redisClient
	//fm.Gdrive = &gdrive.GdriveClient{}
	//fm.Gdrive.Repository = &rep
	//var gHandler = &eb.OdGdriveEventHandler{FolderManager: fm}
	//handlerList = append(handlerList, gHandler)

	// starts video clips processor
	odqRep := reps.OdQueueRepository{Connection: mainConn}
	streamRep := reps.StreamRepository{Connection: mainConn}
	vcp := vc.VideoClipProcessor{Config: config, OdqRep: &odqRep, StreamRep: &streamRep}
	go vcp.Start()
	// ends video clips processor

	var handler = &eb.ComboEventHandler{
		EventHandlers: handlerList,
	}

	var e = eb.EventBus{PubSubConnection: createRedisClient(host, port, EVENTBUS), Channel: "od_service"}
	e.Subscribe(handler)
}

func listenFrEventHandler(config *models.Config, host string, port int) {
	handlerList := make([]eb.EventHandler, 0)
	fhr := &reps.FrHandlerRepository{Config: config}
	var diskHandler = eb.FrDiskEventHandler{Fhr: fhr}
	handlerList = append(handlerList, &diskHandler)

	var handler = &eb.ComboEventHandler{
		EventHandlers: handlerList,
	}

	var e = eb.EventBus{PubSubConnection: createRedisClient(host, port, EVENTBUS), Channel: "fr_service"}
	e.Subscribe(handler)
}

func listenAlprEventHandler(config *models.Config, host string, port int) {
	handlerList := make([]eb.EventHandler, 0)
	ahr := &reps.AlprHandlerRepository{Config: config}
	var diskHandler = eb.AlprDiskEventHandler{Ahr: ahr}
	handlerList = append(handlerList, &diskHandler)

	var handler = &eb.ComboEventHandler{
		EventHandlers: handlerList,
	}

	var e = eb.EventBus{PubSubConnection: createRedisClient(host, port, EVENTBUS), Channel: "alpr_service"}
	e.Subscribe(handler)
}
