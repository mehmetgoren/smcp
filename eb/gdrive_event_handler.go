package eb

import (
	"github.com/go-redis/redis/v8"
	"log"
	"smcp/gdrive"
	"smcp/utils"
)

type GdriveEventHandler struct {
	*gdrive.FolderManager
}

func (g *GdriveEventHandler) Handle(event *redis.Message) (interface{}, error) {
	defer utils.HandlePanic()

	ii, err := CreateImageInfo(event)

	file, err := g.UploadImage(ii.FileName, ii.Base64Image)
	if err != nil {
		log.Println("GdriveEventHandler: An error occurred during the handling image uploading to google drive")
		return nil, err
	}

	return file, nil
}
