package gdrive

import (
	"github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	redsyncredis "github.com/go-redsync/redsync/v4/redis"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
	"google.golang.org/api/drive/v3"
	"log"
	"strconv"
	"time"
)

var rootDirectoryName string = "Smart Machines"

type GoogleDriveFolders struct {
	SmartMachines *drive.File
	Year          *drive.File
	Months        map[string]*drive.File
	Days          map[string]map[string]*drive.File
}

type FolderManager struct {
	Gdrive  *GdriveClient
	Redis   *redis.Client
	Pool    redsyncredis.Pool
	Folders GoogleDriveFolders
}

func parseMonth(month time.Month) (int, string) {
	monthInt := int(month)
	monthStr := strconv.Itoa(monthInt)
	return monthInt, monthStr
}

var mutexName = "global-mutex-today"

func (d *FolderManager) getTodayFolder() (*drive.File, error) {
	if d.Pool == nil {
		d.Pool = goredis.NewPool(d.Redis) // or, pool := redigo.NewPool(...)
	}

	rs := redsync.New(d.Pool)
	mutex := rs.NewMutex(mutexName)

	var err error
	if err = mutex.Lock(); err != nil {
		log.Println("An error occurred on mutex lock: " + err.Error())
		return nil, err
	}
	defer func(mutex *redsync.Mutex) {
		if ok, err := mutex.Unlock(); !ok || err != nil {
			log.Println("An error occurred on mutex unlock: " + err.Error())
		}
	}(mutex)

	var smartMachineFolder *drive.File
	smartMachineFolder = d.Folders.SmartMachines
	if smartMachineFolder == nil {
		smartMachineFolder, err = d.Gdrive.FindFolderByName(rootDirectoryName) // look first if ti was created before.
		if err != nil {
			log.Fatalf("Unable to parse client secret file to config: %v", err)
			return nil, err
		}
		if smartMachineFolder == nil {
			smartMachineFolder, err = d.Gdrive.CreateFolder(rootDirectoryName)
			if err != nil {
				log.Fatalf("An error occurred on create Smart Machine folder: " + err.Error())
				return nil, err
			}
			log.Println("Smart Machine folder has been created first time")
		} else {
			log.Println("Smart Machine folder has been created before, no need to create again")
		}
		d.Folders.SmartMachines = smartMachineFolder
	} else {
		log.Println("Smart Machine folder has been already read and is being used now")
	}

	now := time.Now()
	var getOrCreateChildFolder = func(parentFolder *drive.File, childFolder *drive.File, childFolderName string) (*drive.File, error) {
		if childFolder == nil { //if it was not created before.
			var err error
			childFolder, err = d.Gdrive.FindChildFolder(parentFolder, childFolderName) // look if it was created by other process before.
			if err != nil {
				log.Fatalf("An error occurred on find folder: " + childFolderName + ", err: " + err.Error())
				return nil, err
			}

			if childFolder == nil { // if it wasn't created, lets' dot it for first time
				log.Println(childFolderName + " Child folder is now creating first time: ")
				childFolder, err = d.Gdrive.CreateChildFolder(parentFolder, childFolderName)
				if err != nil {
					log.Fatalf("An error occurred on create folder: " + childFolderName + ", err: " + err.Error())
					return nil, err
				}
			} else {
				log.Println(childFolderName + " Child folder has been created before,no need to create again: ")
			}
		} else {
			log.Println(childFolderName + " Child folder has been already read and is being used now: ")
		}

		return childFolder, nil
	}

	var todayFolder *drive.File = nil
	yearStr := strconv.Itoa(now.Year())
	var yearFolder *drive.File
	yearFolder, err = getOrCreateChildFolder(smartMachineFolder, d.Folders.Year, yearStr) // d.Folders.Year
	if err == nil {
		d.Folders.Year = yearFolder
		if d.Folders.Months == nil {
			d.Folders.Months = make(map[string]*drive.File)
		}
		_, monthStr := parseMonth(now.Month())
		var monthFolder *drive.File
		monthFolder, err = getOrCreateChildFolder(yearFolder, d.Folders.Months[monthStr], monthStr)
		if err == nil {
			d.Folders.Months[monthStr] = monthFolder

			if d.Folders.Days == nil {
				d.Folders.Days = make(map[string]map[string]*drive.File)
			}
			if d.Folders.Days[monthStr] == nil {
				d.Folders.Days[monthStr] = make(map[string]*drive.File)
			}
			dayStr := strconv.Itoa(now.Day())
			_, exists := d.Folders.Days[monthStr][dayStr]
			if !exists {
				d.Folders.Days[monthStr][dayStr] = nil
			}

			todayFolder, err = getOrCreateChildFolder(monthFolder, d.Folders.Days[monthStr][dayStr], dayStr)
			if err == nil {
				d.Folders.Days[monthStr][dayStr] = todayFolder
			}
		}
	}

	return todayFolder, nil
}

func (d *FolderManager) UploadImage(fileName string, base64Image *string) (*drive.File, error) {
	todayFolder, err := d.getTodayFolder()
	if err != nil {
		return nil, err
	}
	file, err := d.Gdrive.CreateImageFile(todayFolder.Id, fileName, base64Image)
	if err != nil {
		log.Fatalf("Unable to upload image on drive service: %v", err)
		return nil, err
	}

	return file, nil
}
