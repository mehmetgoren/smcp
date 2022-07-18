package vc

import (
	"github.com/go-co-op/gocron"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path"
	"smcp/data"
	"smcp/data/cmn"
	"smcp/models"
	"smcp/reps"
	"smcp/utils"
	"time"
)

type AiClipProcessor struct {
	Config    *models.Config
	OdqRep    *reps.OdQueueRepository
	StreamRep *reps.StreamRepository

	Factory *cmn.Factory
}

func (v *AiClipProcessor) getAiRecordPath(sourceId string) string {
	return path.Join(utils.GetRecordPath(v.Config), sourceId, "ai")
}

func (v *AiClipProcessor) getIndexedSourceVideosPath(clip *AiClipObject) string {
	rootPath := utils.GetOdVideosPathBySourceId(v.Config, clip.SourceId)
	ti := utils.TimeIndex{}
	ti.SetValuesFrom(&clip.CreatedAtTime)
	return ti.GetIndexedPath(rootPath)
}

var multiplier = 5
var minFileCount = 2

var emptyFileInfos = make([]fs.FileInfo, 0)

func (v *AiClipProcessor) getAiVideoFolders(sourceId string) []fs.FileInfo {
	aiRootPath := v.getAiRecordPath(sourceId)
	videoFiles, _ := ioutil.ReadDir(aiRootPath)
	if len(videoFiles) < 2 {
		log.Println("no video clips found on the ai folder")
		return emptyFileInfos
	}
	//remove the last one
	clipsCount := len(videoFiles)
	videoFiles = videoFiles[0 : clipsCount-minFileCount]
	return videoFiles
}

func (v *AiClipProcessor) createVideoClipInfos() ([]*AiClipObject, error) {
	hasDetectionVideoClips := make([]*AiClipObject, 0)

	duration := v.Config.Ai.VideoClipDuration
	allDetectedObjects, _ := v.OdqRep.PopAll()
	streams, _ := v.StreamRep.GetAll()
	for _, stream := range streams {
		if !stream.AiClipEnabled {
			continue
		}
		sourceId := stream.Id
		aiVideoFiles := v.getAiVideoFolders(sourceId)
		for _, aiVideoFi := range aiVideoFiles {
			aiFileName := aiVideoFi.Name()
			vci := AiClipObject{}
			vci.SourceId = sourceId
			vci.ObjectDetectionModels = make([]*models.ObjectDetectionModel, 0)
			vci.FileName = aiFileName
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
				aiRootPath := v.getAiRecordPath(sourceId)
				os.Remove(path.Join(aiRootPath, aiFileName))
				log.Println("an ai video file deleted: " + aiFileName)
			}
		}
	}

	return hasDetectionVideoClips, nil
}

func (v *AiClipProcessor) move(clips []*AiClipObject) error {
	defer utils.HandlePanic()

	rep := v.Factory.CreateRepository()
	for _, clip := range clips {
		//move video clips' to persistent folder
		oldLocation := path.Join(v.getAiRecordPath(clip.SourceId), clip.FileName)

		indexedSourceVideosPath := v.getIndexedSourceVideosPath(clip)
		err := utils.CreateDirectoryIfNotExists(indexedSourceVideosPath)
		if err != nil {
			log.Println("an error occurred during the creating indexed data directory, ", err)
			continue
		}
		newLocation := path.Join(indexedSourceVideosPath, clip.FileName)
		os.Rename(oldLocation, newLocation) //moves the short video clip file

		if clip.ObjectDetectionModels == nil || len(clip.ObjectDetectionModels) == 0 {
			continue
		}
		for _, od := range clip.ObjectDetectionModels {
			aiClipModel := data.AiClip{}
			aiClipModel.Setup(newLocation, clip.CreatedAt, clip.LastModified, clip.Duration)
			rep.SetOdVideoClipFields(od.Id, &aiClipModel)
		}
	}
	return nil
}

func (v *AiClipProcessor) check() {
	defer utils.HandlePanic()
	log.Println("Video Clip Processor checking has been started at " + utils.TimeToString(time.Now(), true))
	clips, _ := v.createVideoClipInfos()
	v.move(clips)
}

func (v *AiClipProcessor) Start() {
	defer utils.HandlePanic()

	s := gocron.NewScheduler(time.UTC)

	s.Every(v.Config.Ai.VideoClipDuration * multiplier).Seconds().Do(v.check)

	s.StartAsync()
}
