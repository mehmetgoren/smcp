package main

import (
	"github.com/go-redis/redis/v8"
	"log"
	"smcp/eb"
	"smcp/models"
	"smcp/reps"
	"smcp/utils"
	"smcp/vc"
)

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

	mainConn := reps.CreateRedisConnection(reps.MAIN)
	utils.SetPool(mainConn)

	var configRep = reps.ConfigRepository{Connection: mainConn}
	config, _ := configRep.GetConfig()

	serviceName := "cloud_integration_service"
	heartbeat := createHeartbeatRepository(mainConn, serviceName, config)
	go heartbeat.Start()

	serviceRepository := createServiceRepository(mainConn)
	go func() {
		_, err := serviceRepository.Add(serviceName)
		if err != nil {
			log.Println("An error occurred while registering process id, error is:" + err.Error())
		}
	}()

	pubSubConn := reps.CreateRedisConnection(reps.EVENTBUS)
	go listenOdEventHandlers(mainConn, pubSubConn, config)
	go listenAlprEventHandler(pubSubConn, config)
	listenFrEventHandler(pubSubConn, config)
}

func listenOdEventHandlers(mainConn *redis.Client, pubSubConn *redis.Client, config *models.Config) {
	handlerList := make([]eb.EventHandler, 0)
	ohr := &reps.OdHandlerRepository{Config: config}
	var diskHandler = eb.OdEventHandler{Ohr: ohr}
	handlerList = append(handlerList, &diskHandler)

	//detection series handler
	var vch = eb.OdAiClipEventHandler{Connection: mainConn}
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
	vcp := vc.AiClipProcessor{Config: config, OdqRep: &odqRep, StreamRep: &streamRep}
	go vcp.Start()
	// ends video clips processor

	var comboHandler = &eb.ComboEventHandler{
		EventHandlers: handlerList,
	}

	var e = eb.EventBus{PubSubConnection: pubSubConn, Channel: "od_service"}
	e.Subscribe(comboHandler)
}

func listenFrEventHandler(pubSubConn *redis.Client, config *models.Config) {
	handlerList := make([]eb.EventHandler, 0)
	fhr := &reps.FrHandlerRepository{Config: config}
	var diskHandler = eb.FrEventHandler{Fhr: fhr}
	handlerList = append(handlerList, &diskHandler)

	var handler = &eb.ComboEventHandler{
		EventHandlers: handlerList,
	}

	var e = eb.EventBus{PubSubConnection: pubSubConn, Channel: "fr_service"}
	e.Subscribe(handler)
}

func listenAlprEventHandler(pubSubConn *redis.Client, config *models.Config) {
	handlerList := make([]eb.EventHandler, 0)
	ahr := &reps.AlprHandlerRepository{Config: config}
	var diskHandler = eb.AlprEventHandler{Ahr: ahr}
	handlerList = append(handlerList, &diskHandler)

	var handler = &eb.ComboEventHandler{
		EventHandlers: handlerList,
	}

	var e = eb.EventBus{PubSubConnection: pubSubConn, Channel: "alpr_service"}
	e.Subscribe(handler)
}
