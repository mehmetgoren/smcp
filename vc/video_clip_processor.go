package vc

import (
	"encoding/json"
	"github.com/go-co-op/gocron"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path"
	"smcp/disk"
	"smcp/models"
	"smcp/reps"
	"smcp/utils"
	"strings"
	"time"
)

type VideoClipProcessor struct {
	Config    *models.Config
	DoRep     *DetectedObjectQueueRepository
	StreamRep *reps.StreamRepository
}

var pathPrefix = "vcs"

func (v *VideoClipProcessor) getRootPath(sourceId string) string {
	return path.Join(utils.GetRecordFolderPath(v.Config), sourceId, pathPrefix)
}

func (v *VideoClipProcessor) getTempRootPath(sourceId string) string {
	return path.Join(utils.GetRecordFolderPath(v.Config), sourceId, pathPrefix, "temp")
}

var emptyFileInfos = make([]fs.FileInfo, 0)

func (v *VideoClipProcessor) getTempVideoFolders(sourceId string) []fs.FileInfo {
	tempRootPath := v.getTempRootPath(sourceId)
	videoFiles, _ := ioutil.ReadDir(tempRootPath)
	if len(videoFiles) < 2 {
		log.Println("no video clips found on the temp folder")
		return emptyFileInfos
	}
	//remove the last one
	clipsCount := len(videoFiles)
	videoFiles = videoFiles[0 : clipsCount-1]
	return videoFiles
}

func (v *VideoClipProcessor) createVideoClipInfos() ([]*VideoClipJsonObject, error) {
	hasDetectionVideoClips := make([]*VideoClipJsonObject, 0)

	duration := v.Config.Ai.VideoClipDuration
	allDetectedObjects, _ := v.DoRep.PopAll()
	streams, _ := v.StreamRep.GetAll()
	for _, stream := range streams {
		if !stream.VideoClipEnabled {
			continue
		}
		sourceId := stream.Id
		tempVideoFiles := v.getTempVideoFolders(sourceId)
		for _, tempVideoFi := range tempVideoFiles {
			tempFileName := tempVideoFi.Name()
			vci := VideoClipJsonObject{}
			vci.FileName = tempFileName
			vci.Duration = duration
			vci.SetupDateTimes()
			deleteVideoFile := true
			for _, detectedObject := range allDetectedObjects {
				createdAtTime := utils.StringToTime(detectedObject.CreatedAt, false)
				if vci.IsInTimeSpan(createdAtTime) {
					vci.DetectedImage = detectedObject
					hasDetectionVideoClips = append(hasDetectionVideoClips, &vci)
					deleteVideoFile = false
					break
				}
			}

			if deleteVideoFile {
				//delete the non-object detection containing video files
				tempRootPath := v.getTempRootPath(sourceId)
				os.Remove(path.Join(tempRootPath, tempFileName))
				log.Println("a temp video file deleted: " + tempFileName)
			}
		}
	}

	return hasDetectionVideoClips, nil
}

func (v *VideoClipProcessor) move(clips []*VideoClipJsonObject) error {
	defer utils.HandlePanic()

	for _, clip := range clips {
		fm := &disk.FolderManager{Redis: v.DoRep.Connection, RootFolderPath: v.getRootPath(clip.DetectedImage.SourceId)}
		//move video clips' to persistent folder
		oldLocation := path.Join(v.getTempRootPath(clip.DetectedImage.SourceId), clip.FileName)
		provider := disk.FileNameIndexedFolderInfoProvider{FileName: clip.FileName}
		folderPath, _ := fm.CreateFolderIfNotExists(provider) //creates year/month/day/hour folder
		newLocation := path.Join(folderPath, clip.FileName)
		os.Rename(oldLocation, newLocation) //moves the short video clip file

		//and also create a json file next to the video clip file for metadata
		clip.FileName = path.Join(clip.DetectedImage.SourceId, pathPrefix, strings.Replace(newLocation, fm.RootFolderPath, "", -1))
		videoFileExt := path.Ext(newLocation)
		jsonFullFileName := strings.Replace(newLocation, videoFileExt, ".json", -1)
		jsonListBytes, _ := json.Marshal(clip)
		ioutil.WriteFile(jsonFullFileName, jsonListBytes, 0777)
		log.Println("a object detected video file has been moved to a indexed folder")
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

	s.Every(v.Config.Ai.VideoClipDuration * 2).Seconds().Do(v.check)

	s.StartAsync()
}
