package utils

import (
	"github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	redsyncredis "github.com/go-redsync/redsync/v4/redis"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
	"log"
	"os"
	"path"
	"smcp/abstract"
	"smcp/models"
)

func GetRecordPath(dirPath string) string {
	return path.Join(dirPath, "record")
}

func getAiPath(dirPath string) string {
	return path.Join(dirPath, "ai")
}

func GetAiImagesPathBySourceId(config *models.Config, sourceId string) string {
	sourceDirPath := getSourceDirPath(config, sourceId)
	return path.Join(getAiPath(sourceDirPath), sourceId, "images")
}

func GetAiClipPathBySourceId(config *models.Config, sourceId string) string {
	sourceDirPath := getSourceDirPath(config, sourceId)
	return path.Join(GetRecordPath(sourceDirPath), sourceId, "ai")
}

var pool redsyncredis.Pool = nil

func SetPool(conn *redis.Client) {
	pool = goredis.NewPool(conn)
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

var sourceDirMaps *TTLMap[*models.StreamModel] = nil
var streamRepository abstract.Repository[*models.StreamModel] = nil

func SetDirParameters(sdm *TTLMap[*models.StreamModel], sr abstract.Repository[*models.StreamModel]) {
	sourceDirMaps = sdm
	streamRepository = sr
}

// getDefaultDirPath The first item of the dirPaths is default and the SenseAI backup file also use it
func getDefaultDirPath(config *models.Config) string {
	dirPaths := config.General.DirPaths
	if dirPaths == nil || len(dirPaths) == 0 || dirPaths[0] == "" {
		log.Fatal("Config.General.DirPaths is empty, the program will be terminated")
	}
	rootPath := dirPaths[0]
	return rootPath
}

func getSourceDirPath(config *models.Config, sourceId string) string {
	stream := sourceDirMaps.Get(sourceId)
	var err error
	if stream == nil {
		stream, err = streamRepository.Get(sourceId)
		if err != nil || stream == nil {
			log.Fatal("An error occurred on getting source dir path, the process is now being ended, err: " + err.Error())
			return ""
		}
		sourceDirMaps.Put(sourceId, stream)
	}
	sourceDirPath := ""
	if len(stream.RootDirPath) > 0 {
		sourceDirPath = stream.RootDirPath
	} else {
		sourceDirPath = getDefaultDirPath(config)
	}

	return sourceDirPath
}
