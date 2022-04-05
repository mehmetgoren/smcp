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
)

type FolderManager struct {
	RootFolderPath string
	Redis          *redis.Client
	pool           redsyncredis.Pool
}

func (fm *FolderManager) CreateFolderIfNotExists(provider IndexedFolderInfoProvider) (string, error) {
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

	folderPath := provider.Create().GetIndexedPath(fm.RootFolderPath)
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		err := os.MkdirAll(folderPath, 0755)
		if err != nil {
			log.Println("an error occurred during the creating hour folder: " + err.Error())
			return "", err
		}
	}
	return folderPath, nil
}

func (fm *FolderManager) SaveFile(provider IndexedFolderInfoProvider, fileName string, file []byte) (string, error) {
	folderPath, _ := fm.CreateFolderIfNotExists(provider)
	filePath := path.Join(folderPath, fileName)
	err := ioutil.WriteFile(filePath, file, 0644)
	if err != nil {
		log.Println("An error occurred on FolderManager ioutil.WriteFile: " + err.Error())
		return "", err
	}

	return filePath, nil
}
