package disk

import (
	"github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	redsyncredis "github.com/go-redsync/redsync/v4/redis"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
	"time"
)

type FolderManager struct {
	SmartMachineFolderPath string
	Redis                  *redis.Client
	pool                   redsyncredis.Pool
}

func (fm *FolderManager) createCurrentHourPath() string {
	now := time.Now()
	year := strconv.Itoa(now.Year())
	month := strconv.Itoa(int(now.Month()))
	day := strconv.Itoa(now.Day())
	hour := strconv.Itoa(now.Hour())
	return path.Join(fm.SmartMachineFolderPath, year, month, day, hour)
}

func (fm *FolderManager) createFolderIfNotExists() (string, error) {
	folderPath := fm.createCurrentHourPath()
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		err := os.MkdirAll(folderPath, 0755)
		if err != nil {
			log.Println("an error occurred during the creating hour folder: " + err.Error())
			return "", err
		}
	}
	return folderPath, nil
}

func (fm *FolderManager) getFilePath(fileName string) (string, error) {
	if fm.pool == nil {
		fm.pool = goredis.NewPool(fm.Redis) // or, pool := redigo.NewPool(...)
	}

	rs := redsync.New(fm.pool)
	mutex := rs.NewMutex("mutex-disk-manager")

	var err error
	if err = mutex.Lock(); err != nil {
		log.Println("An error occurred on FolderManager mutex lock: " + err.Error())
		return "", err
	}
	defer func(mutex *redsync.Mutex) {
		if ok, err := mutex.Unlock(); !ok || err != nil {
			log.Println("An error occurred on FolderManager mutex unlock: " + err.Error())
		}
	}(mutex)

	folderPath, err := fm.createFolderIfNotExists()
	if err != nil {
		log.Println("An error occurred on FolderManager createFolderIfNotExists: " + err.Error())
		return "", err
	}

	return path.Join(folderPath, fileName), nil
}

func (fm *FolderManager) SaveFile(fileName string, file []byte) (string, error) {
	filePath, err := fm.getFilePath(fileName)
	if err != nil {
		log.Println("An error occurred on FolderManager getFilePath: " + err.Error())
		return "", err
	}

	err = ioutil.WriteFile(filePath, file, 0644)
	if err != nil {
		log.Println("An error occurred on FolderManager ioutil.WriteFile: " + err.Error())
		return "", err
	}

	return filePath, nil
}
