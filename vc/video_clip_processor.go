package vc

import (
	"encoding/json"
	"github.com/go-co-op/gocron"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path"
	"smcp/models"
	"smcp/reps"
	"smcp/utils"
	"strings"
	"time"
)

type VideoClipProcessor struct {
	Config    *models.Config
	OdqRep    *reps.ObjectDetectionQueueRepository
	StreamRep *reps.StreamRepository
}

func (v *VideoClipProcessor) getTempRecordPath(sourceId string) string {
	return path.Join(utils.GetRecordPath(v.Config), sourceId, "temp")
}

func (v *VideoClipProcessor) getIndexedSourceVideosPath(clip *VideoClipObject) string {
	rootPath := utils.GetOdVideosPathBySourceId(v.Config, clip.SourceId)
	ti := reps.TimeIndex{}
	ti.SetValuesFrom(&clip.CreatedAtTime)
	return ti.GetIndexedPath(rootPath)
}

func (v *VideoClipProcessor) getIndexedSourceDataPath(clip *VideoClipObject) string {
	rootPath := path.Join(utils.GetOdDataPathBySourceId(v.Config, clip.SourceId))
	ti := reps.TimeIndex{}
	ti.SetValuesFrom(&clip.CreatedAtTime)
	return ti.GetIndexedPath(rootPath)
}

var multiplier = 5
var minFileCount = 3

var emptyFileInfos = make([]fs.FileInfo, 0)

func (v *VideoClipProcessor) getTempVideoFolders(sourceId string) []fs.FileInfo {
	tempRootPath := v.getTempRecordPath(sourceId)
	videoFiles, _ := ioutil.ReadDir(tempRootPath)
	if len(videoFiles) < 2 {
		log.Println("no video clips found on the temp folder")
		return emptyFileInfos
	}
	//remove the last one
	clipsCount := len(videoFiles)
	videoFiles = videoFiles[0 : clipsCount-minFileCount]
	return videoFiles
}

func (v *VideoClipProcessor) createVideoClipInfos() ([]*VideoClipObject, error) {
	hasDetectionVideoClips := make([]*VideoClipObject, 0)

	duration := v.Config.Ai.VideoClipDuration
	allDetectedObjects, _ := v.OdqRep.PopAll()
	streams, _ := v.StreamRep.GetAll()
	for _, stream := range streams {
		if !stream.VideoClipEnabled {
			continue
		}
		sourceId := stream.Id
		tempVideoFiles := v.getTempVideoFolders(sourceId)
		for _, tempVideoFi := range tempVideoFiles {
			tempFileName := tempVideoFi.Name()
			vci := VideoClipObject{}
			vci.SourceId = sourceId
			vci.ObjectDetectionModels = make([]*models.ObjectDetectionModel, 0)
			vci.FileName = tempFileName
			vci.Duration = duration
			vci.SetupDateTimes()
			for _, detectedObject := range allDetectedObjects {
				createdAtTime := utils.StringToTime(detectedObject.CreatedAt)
				if vci.IsInTimeSpan(createdAtTime) {
					vci.ObjectDetectionModels = append(vci.ObjectDetectionModels, detectedObject)
				}
			}

			if len(vci.ObjectDetectionModels) > 0 {
				hasDetectionVideoClips = append(hasDetectionVideoClips, &vci)
			} else {
				//delete the non-object detection containing video files
				tempRootPath := v.getTempRecordPath(sourceId)
				os.Remove(path.Join(tempRootPath, tempFileName))
				log.Println("a temp video file deleted: " + tempFileName)
			}
		}
	}

	return hasDetectionVideoClips, nil
}

func (v *VideoClipProcessor) move(clips []*VideoClipObject) error {
	defer utils.HandlePanic()

	for _, clip := range clips {
		//move video clips' to persistent folder
		oldLocation := path.Join(v.getTempRecordPath(clip.SourceId), clip.FileName)

		indexedSourceVideosPath := v.getIndexedSourceVideosPath(clip)
		err := utils.CreateDirectoryIfNotExists(indexedSourceVideosPath)
		if err != nil {
			log.Println("an error occurred during the creating indexed data directory, ", err)
			continue
		}
		newLocation := path.Join(indexedSourceVideosPath, clip.FileName)
		os.Rename(oldLocation, newLocation) //moves the short video clip file

		//and also create a json file next to the video clip file for metadata
		indexedSourceDataPath := v.getIndexedSourceDataPath(clip)
		fileInfos, _ := ioutil.ReadDir(indexedSourceDataPath)
		for _, fileInfo := range fileInfos {
			splits := strings.Split(fileInfo.Name(), "_")
			id := strings.Split(splits[len(splits)-1], ".")[0]
			odModel := findOdModel(clip.ObjectDetectionModels, id)
			if odModel == nil {
				continue // if it doesn't match, do not mutate the file.
			}

			jsonDataFileName := path.Join(indexedSourceDataPath, fileInfo.Name())

			//read json file
			fileBytes, _ := ioutil.ReadFile(jsonDataFileName)
			jo := &models.ObjectDetectionJsonObject{}
			json.Unmarshal(fileBytes, jo)

			//change the data
			jo.Video.FileName = strings.Replace(newLocation, v.Config.General.RootFolderPath+"/", "", -1)
			jo.Video.CreatedAt = clip.CreatedAt
			jo.Video.LastModifiedAt = clip.LastModified
			jo.Video.Duration = clip.Duration

			//write json file
			objectBytes, _ := json.Marshal(jo)
			ioutil.WriteFile(jsonDataFileName, objectBytes, 0777)

		}
	}
	return nil
}

func findOdModel(list []*models.ObjectDetectionModel, id string) *models.ObjectDetectionModel {
	for _, item := range list {
		if item.Id == id {
			return item
		}
	}
	return nil
}

func (v *VideoClipProcessor) check() {
	defer utils.HandlePanic()
	log.Println("Video Clip Processor checking has been started at " + utils.TimeToString(time.Now(), true))
	clips, _ := v.createVideoClipInfos()
	v.move(clips)
}

func (v *VideoClipProcessor) Start() {
	defer utils.HandlePanic()

	s := gocron.NewScheduler(time.UTC)

	s.Every(v.Config.Ai.VideoClipDuration * multiplier).Seconds().Do(v.check)

	s.StartAsync()
}
