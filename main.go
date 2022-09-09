package main

import (
	"github.com/go-redis/redis/v8"
	"log"
	"smcp/data/cmn"
	"smcp/eb"
	"smcp/gdrive"
	"smcp/models"
	"smcp/reps"
	"smcp/tb"
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

	factory := &cmn.Factory{Config: config}
	factory.Init()
	defer factory.Close()

	notifier := &eb.NotifierPublisher{PubSubConnection: pubSubConn}
	cloudRep := &reps.CloudRepository{Connection: mainConn}
	var tbc *tb.TelegramBotClient = nil
	if cloudRep.IsTelegramIntegrationEnabled() {
		tbc, _ = tb.CreateTelegramBot(cloudRep)
	}

	pars := &ListenParams{Config: config, Tcb: tbc, CloudRep: cloudRep, MainConn: mainConn, PubSubConn: pubSubConn, Factory: factory, Notifier: notifier}
	startVideoClipProcessor(pars)
	go listenOdEventHandlers(pars)
	go listenFrEventHandler(pars)
	go listenAlprEventHandler(pars)
	go listenVideoFilesEventHandlers(pars)
	listenNotifyFailedEventHandler(pars)
}

func startVideoClipProcessor(pars *ListenParams) {
	odqRep := reps.AiClipQueueRepository{Connection: pars.MainConn}
	streamRep := reps.StreamRepository{Connection: pars.MainConn}
	vcp := vc.AiClipProcessor{Config: pars.Config, AiQuRep: &odqRep, StreamRep: &streamRep, Factory: pars.Factory}
	go vcp.Start()
}

func createCloudEventHandlers(pars *ListenParams, aiType int) ([]eb.EventHandler, error) {
	handlerList := make([]eb.EventHandler, 0)
	if pars.CloudRep.IsTelegramIntegrationEnabled() {
		var tbHandler eb.EventHandler = &eb.OdTelegramEventHandler{TelegramBotClient: pars.Tcb, AiType: aiType}
		handlerList = append(handlerList, tbHandler)
	}

	if pars.CloudRep.IsGdriveIntegrationEnabled() {
		var fm = &gdrive.FolderManager{Redis: pars.MainConn, Client: &gdrive.Client{}}
		fm.Client.Repository = pars.CloudRep
		var gHandler = &eb.GdriveEventHandler{FolderManager: fm, AiType: aiType}
		handlerList = append(handlerList, gHandler)
	}

	return handlerList, nil
}

func listenOdEventHandlers(pars *ListenParams) {
	handlerList := make([]eb.EventHandler, 0)
	var diskHandler = eb.OdEventHandler{Factory: pars.Factory, Notifier: pars.Notifier}
	handlerList = append(handlerList, &diskHandler)

	//detection series handler
	var ace = eb.AiClipEventHandler{Connection: pars.MainConn, AiType: models.Od}
	handlerList = append(handlerList, &ace)

	cloudHandlers, err := createCloudEventHandlers(pars, eb.ObjectDetection)
	if err == nil && cloudHandlers != nil && len(cloudHandlers) > 0 {
		for _, ch := range cloudHandlers {
			handlerList = append(handlerList, ch)
		}
	} else {
		log.Println("No Cloud Provider has been register for Object Detection")
	}

	var comboHandler = &eb.ComboEventHandler{
		EventHandlers: handlerList,
	}

	var e = eb.EventBus{PubSubConnection: pars.PubSubConn, Channel: "od_service"}
	e.Subscribe(comboHandler)
}

func listenFrEventHandler(pars *ListenParams) {
	handlerList := make([]eb.EventHandler, 0)
	var diskHandler = eb.FrEventHandler{Factory: pars.Factory, Notifier: pars.Notifier}
	handlerList = append(handlerList, &diskHandler)

	//detection series handler
	var ace = eb.AiClipEventHandler{Connection: pars.MainConn, AiType: models.Fr}
	handlerList = append(handlerList, &ace)

	cloudHandlers, err := createCloudEventHandlers(pars, eb.FaceRecognition)
	if err == nil && cloudHandlers != nil && len(cloudHandlers) > 0 {
		for _, ch := range cloudHandlers {
			handlerList = append(handlerList, ch)
		}
	} else {
		log.Println("No Cloud Provider has been register for Face Recognition")
	}

	var handler = &eb.ComboEventHandler{
		EventHandlers: handlerList,
	}

	var e = eb.EventBus{PubSubConnection: pars.PubSubConn, Channel: "fr_service"}
	e.Subscribe(handler)
}

func listenAlprEventHandler(pars *ListenParams) {
	handlerList := make([]eb.EventHandler, 0)
	var diskHandler = eb.AlprEventHandler{Factory: pars.Factory, Notifier: pars.Notifier}
	handlerList = append(handlerList, &diskHandler)

	//detection series handler
	var ace = eb.AiClipEventHandler{Connection: pars.MainConn, AiType: models.Alpr}
	handlerList = append(handlerList, &ace)

	cloudHandlers, err := createCloudEventHandlers(pars, eb.PlateRecognition)
	if err == nil && cloudHandlers != nil && len(cloudHandlers) > 0 {
		for _, ch := range cloudHandlers {
			handlerList = append(handlerList, ch)
		}
	} else {
		log.Println("No Cloud Provider has been register for License Plate Recognition")
	}

	var handler = &eb.ComboEventHandler{
		EventHandlers: handlerList,
	}

	var e = eb.EventBus{PubSubConnection: pars.PubSubConn, Channel: "alpr_service"}
	e.Subscribe(handler)
}

func listenVideoFilesEventHandlers(pars *ListenParams) {
	var vfiHandler = &eb.VfiResponseEventHandler{Factory: pars.Factory}
	var e = eb.EventBus{PubSubConnection: pars.PubSubConn, Channel: "vfi_response"}
	go e.Subscribe(vfiHandler)

	var vfmHandler = &eb.VfmResponseEventHandler{Factory: pars.Factory}
	var e2 = eb.EventBus{PubSubConnection: pars.PubSubConn, Channel: "vfm_response"}
	e2.Subscribe(vfmHandler)
}

func listenNotifyFailedEventHandler(pars *ListenParams) {
	var nfh = &eb.NotifyFailedHandler{Notifier: pars.Notifier}
	var e = eb.EventBus{PubSubConnection: pars.PubSubConn, Channel: "notify_failed"}
	e.Subscribe(nfh)
}

type ListenParams struct {
	MainConn   *redis.Client
	PubSubConn *redis.Client

	Factory *cmn.Factory

	Notifier *eb.NotifierPublisher
	Config   *models.Config
	CloudRep *reps.CloudRepository
	Tcb      *tb.TelegramBotClient
}
