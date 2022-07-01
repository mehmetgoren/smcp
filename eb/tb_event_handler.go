package eb

import (
	"encoding/base64"
	"errors"
	"github.com/go-redis/redis/v8"
	"gopkg.in/tucnak/telebot.v2"
	"io/ioutil"
	"log"
	"smcp/models"
	"smcp/tb"
	"smcp/utils"
	"strconv"
	"strings"
)

type ImageInfo struct {
	FileName    string
	Base64Image *string
}

func CreateImageInfo(event *redis.Message, aiType int) (*ImageInfo, error) {
	ret := &ImageInfo{}
	switch aiType {
	case ObjectDetection:
		var de = models.ObjectDetectionModel{}
		err := utils.DeserializeJson(event.Payload, &de)
		if err != nil {
			return nil, err
		}
		ret.FileName = de.CreateFileName()
		if de.DetectedObjects != nil && len(de.DetectedObjects) > 0 {
			arr := make([]string, 0)
			for _, do := range de.DetectedObjects {
				arr = append(arr, do.PredClsName)
			}
			ret.FileName = strings.Join(arr, ", ")
		}
		ret.Base64Image = &de.Base64Image
		break
	case FaceRecognition:
		var fr = models.FaceRecognitionModel{}
		err := utils.DeserializeJson(event.Payload, &fr)
		if err != nil {
			return nil, err
		}
		ret.FileName = fr.CreateFileName()
		if fr.DetectedFaces != nil && len(fr.DetectedFaces) > 0 {
			arr := make([]string, 0)
			for _, df := range fr.DetectedFaces {
				arr = append(arr, df.PredClsName)
			}
			ret.FileName = strings.Join(arr, ", ")
		}
		ret.Base64Image = &fr.Base64Image
		break
	case PlateRecognition:
		alpr := models.AlprResponse{}
		err := utils.DeserializeJson(event.Payload, &alpr)
		if err != nil {
			return nil, err
		}
		ret.FileName = alpr.CreateFileName()
		if alpr.Results != nil && len(alpr.Results) > 0 {
			arr := make([]string, 0)
			for _, r := range alpr.Results {
				arr = append(arr, r.Plate)
			}
			ret.FileName = strings.Join(arr[:], ", ")
		}
		ret.Base64Image = &alpr.Base64Image
	default:
		return nil, errors.New("Not Supported " + strconv.Itoa(aiType))
	}
	return ret, nil
}

type OdTelegramEventHandler struct {
	*tb.TelegramBotClient
	AiType int
}

func (t *OdTelegramEventHandler) Handle(event *redis.Message) (interface{}, error) {
	defer utils.HandlePanic()

	ii, err := CreateImageInfo(event, t.AiType)
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
			log.Println("TelegramEventHandler: Send error for " + msg.Caption + ". The error is " + sendErr.Error())
			return nil, sendErr
		}
	}

	log.Println("TelegramEventHandler: image send successfully as " + ii.FileName + " the message is " + ii.FileName)

	return nil, nil
}
