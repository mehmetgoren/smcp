package utils

import (
	"github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	redsyncredis "github.com/go-redsync/redsync/v4/redis"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
	"log"
	"os"
	"path"
	"smcp/models"
)

func GetRecordPath(config *models.Config) string {
	return path.Join(config.General.RootFolderPath, "record")
}

func getOdPath(config *models.Config) string {
	return path.Join(config.General.RootFolderPath, "od")
}

func GetOdImagesPathBySourceId(config *models.Config, sourceId string) string {
	return path.Join(getOdPath(config), sourceId, "images")
}

func GetAiClipPathBySourceId(config *models.Config, sourceId string) string {
	return path.Join(GetRecordPath(config), sourceId, "ai")
}

var pool redsyncredis.Pool = nil

func SetPool(conn *redis.Client) {
	pool = goredis.NewPool(conn) // or, pool := redigo.NewPool(...)
}

func CreateDirectoryIfNotExists(directoryPath string) error {
	rs := redsync.New(pool)
	mutex := rs.NewMutex("mutex-disk-manager")

	var err error
	if err = mutex.Lock(); err != nil {
		log.Println("An error occurred on FolderManager mutex lock: " + err.Error())
		return err
	}
	defer func(mutex *redsync.Mutex) {
		if ok, err := mutex.Unlock(); !ok || err != nil {
			log.Println("An error occurred on FolderManager mutex unlock: " + err.Error())
		}
	}(mutex)

	if _, err := os.Stat(directoryPath); os.IsNotExist(err) {
		err := os.MkdirAll(directoryPath, 0755)
		if err != nil {
			log.Println("an error occurred during the creating hour folder: " + err.Error())
			return err
		}
	}
	return nil
}

func getFrPath(config *models.Config) string {
	return path.Join(config.General.RootFolderPath, "fr")
}

func GetFrImagesPathBySourceId(config *models.Config, sourceId string) string {
	return path.Join(getFrPath(config), sourceId, "images")
}

// alpr starts
func getAlprPath(config *models.Config) string {
	return path.Join(config.General.RootFolderPath, "alpr")
}

func GetAlprImagesPathBySourceId(config *models.Config, sourceId string) string {
	return path.Join(getAlprPath(config), sourceId, "images")
}
