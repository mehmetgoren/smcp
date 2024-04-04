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

func createServiceRepository(client *redis.Client) *reps.ServiceRepository {
	var pidRepository = reps.ServiceRepository{Client: client}

	return &pidRepository
}

func main() {
	mainConn := reps.CreateRedisConnection(reps.MAIN)
	utils.SetPool(mainConn)
	utils.SetDirParameters(utils.NewTTLMap[*models.StreamModel](0, 60*15), // 15 minutes
		&reps.StreamRepository{Connection: mainConn})

	var configRep = reps.ConfigRepository{Connection: mainConn}
	config, _ := configRep.GetConfig()

	serviceName := "cloud_integration_service"
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
	go listenAiEventHandlers(pars)
	go listenVideoFilesEventHandlers(pars)
	listenNotifyFailedEventHandler(pars)
}

func startVideoClipProcessor(pars *ListenParams) {
	odqRep := reps.AiClipQueueRepository{Connection: pars.MainConn}
	streamRep := reps.StreamRepository{Connection: pars.MainConn}
	vcp := vc.AiClipProcessor{Config: pars.Config, AiQuRep: &odqRep, StreamRep: &streamRep, Factory: pars.Factory}
	go vcp.Start()
}

func createCloudEventHandlers(pars *ListenParams) ([]eb.EventHandler, error) {
	handlerList := make([]eb.EventHandler, 0)
	if pars.CloudRep.IsTelegramIntegrationEnabled() {
		var tbHandler eb.EventHandler = &eb.TelegramEventHandler{TelegramBotClient: pars.Tcb}
		handlerList = append(handlerList, tbHandler)
	}

	if pars.CloudRep.IsGdriveIntegrationEnabled() {
		var fm = &gdrive.FolderManager{Redis: pars.MainConn, Client: &gdrive.Client{}}
		fm.Client.Repository = pars.CloudRep
		var gHandler = &eb.GdriveEventHandler{FolderManager: fm}
		handlerList = append(handlerList, gHandler)
	}

	return handlerList, nil
}

func listenAiEventHandlers(pars *ListenParams) {
	handlerList := make([]eb.EventHandler, 0)
	var diskHandler = eb.AiEventHandler{Factory: pars.Factory, Notifier: pars.Notifier}
	handlerList = append(handlerList, &diskHandler)

	//detection series handler
	var ace = eb.AiClipEventHandler{Connection: pars.MainConn}
	handlerList = append(handlerList, &ace)

	cloudHandlers, err := createCloudEventHandlers(pars)
	if err == nil && cloudHandlers != nil && len(cloudHandlers) > 0 {
		for _, ch := range cloudHandlers {
			handlerList = append(handlerList, ch)
		}
	} else {
		log.Println("No Cloud Provider has been register for AI Detection")
	}

	var comboHandler = &eb.ComboEventHandler{
		EventHandlers: handlerList,
	}

	var e = eb.EventBus{PubSubConnection: pars.PubSubConn, Channel: "smcp_in"}
	e.Subscribe(comboHandler)
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
