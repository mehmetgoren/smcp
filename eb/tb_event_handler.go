package eb

import (
	"encoding/base64"
	"github.com/go-redis/redis/v8"
	"gopkg.in/tucnak/telebot.v2"
	"io/ioutil"
	"log"
	"smcp/models"
	"smcp/tb"
	"smcp/utils"
	"strings"
)

type ImageInfo struct {
	FileName    string
	Base64Image *string
}

func CreateImageInfo(event *redis.Message) (*ImageInfo, error) {
	ret := &ImageInfo{}
	var adm = models.AiDetectionModel{}
	err := utils.DeserializeJson(event.Payload, &adm)
	if err != nil {
		return nil, err
	}
	ret.FileName = adm.CreateFileName()
	if adm.Detections != nil && len(adm.Detections) > 0 {
		arr := make([]string, 0)
		for _, do := range adm.Detections {
			arr = append(arr, do.Label)
		}
		ret.FileName = strings.Join(arr, ", ")
	}
	ret.Base64Image = &adm.Base64Image

	return ret, nil
}

type TelegramEventHandler struct {
	*tb.TelegramBotClient
}

func (t *TelegramEventHandler) Handle(event *redis.Message) (interface{}, error) {
	defer utils.HandlePanic()

	ii, err := CreateImageInfo(event)
	if err != nil {
		return nil, err
	}

	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(*ii.Base64Image))
	defer ioutil.NopCloser(reader)
	tbFile := telebot.FromReader(reader)
	tbFile.UniqueID = ii.FileName
	tbPhoto := &telebot.Photo{File: tbFile, Caption: ii.FileName}

	users, err := t.CloudRepository.GetTelegramUsers()
	if err != nil {
		return nil, err
	}
	for _, user := range users {
		msg, sendErr := t.Bot.Send(user.MapTo(), tbPhoto)
		if sendErr != nil {
			caption := ""
			if msg != nil {
				caption = msg.Caption
			}
			log.Println("TelegramEventHandler: Send error for " + caption + ". The error is " + sendErr.Error())
			return nil, sendErr
		}
	}

	log.Println("TelegramEventHandler: image send successfully as " + ii.FileName + " the message is " + ii.FileName)

	return nil, nil
}
