package vc

import (
	. "github.com/ahmetb/go-linq/v3"
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
	"strconv"
	"time"
)

type AiClipProcessor struct {
	Config    *models.Config
	AiQuRep   *reps.AiClipQueueRepository
	StreamRep *reps.StreamRepository

	Factory *cmn.Factory
}

func (v *AiClipProcessor) getAiRecordPath(sourceId string) string {
	return utils.GetAiClipPathBySourceId(v.Config, sourceId)
}

func (v *AiClipProcessor) getIndexedSourceVideosPath(clip *AiClipObject) string {
	rootPath := v.getAiRecordPath(clip.SourceId)
	ti := utils.TimeIndex{}
	ti.SetValuesFrom(&clip.CreatedAtTime)
	return ti.GetIndexedPath(rootPath)
}

var multiplier = 5
var minFileCount = 2

var emptyFileInfos = make([]fs.FileInfo, 0)

func (v *AiClipProcessor) getAiVideoFolders(sourceId string) []fs.FileInfo {
	aiRootPath := v.getAiRecordPath(sourceId)
	temp, _ := ioutil.ReadDir(aiRootPath)
	videoFiles := make([]fs.FileInfo, 0)
	for _, f := range temp {
		if f.IsDir() {
			continue
		}
		videoFiles = append(videoFiles, f)
	}
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
	allDetectedObjects, _ := v.AiQuRep.PopAll()
	if len(allDetectedObjects) == 0 {
		return hasDetectionVideoClips, nil
	}

	temp, err := v.StreamRep.GetAll()
	if err != nil {
		log.Println("an error has been occurred while getting a stream object, err: " + err.Error())
		return hasDetectionVideoClips, err
	}
	streamDic := make(map[string][]fs.FileInfo)
	for _, stream := range temp {
		if !stream.AiClipEnabled {
			continue
		}
		aiVideoFiles := v.getAiVideoFolders(stream.Id)
		if len(aiVideoFiles) == 0 {
			continue
		}
		streamDic[stream.Id] = aiVideoFiles
	}

	for sourceId, aiClipFiles := range streamDic {
		allSourceDetectedObjects := make([]*models.AiClipQueueModel, 0)
		From(allDetectedObjects).Where(func(de interface{}) bool {
			return predicateSourceId(de.(*models.AiClipQueueModel), sourceId)
		}).ToSlice(&allSourceDetectedObjects)
		if len(allSourceDetectedObjects) == 0 {
			continue
		}
		for _, aiVideoFi := range aiClipFiles {
			aiFileName := aiVideoFi.Name()
			vci := AiClipObject{}
			vci.SourceId = sourceId
			vci.AiQueueModels = make([]*models.AiClipQueueModel, 0)
			vci.FileName = aiFileName
			vci.Duration = duration
			vci.SetupDateTimes()
			for _, detectedObject := range allSourceDetectedObjects {
				createdAtTime := utils.StringToTime(getCreatedAt(detectedObject))
				if vci.IsInTimeSpan(createdAtTime) {
					vci.AiQueueModels = append(vci.AiQueueModels, detectedObject)
				}
			}
			if len(vci.AiQueueModels) > 0 {
				hasDetectionVideoClips = append(hasDetectionVideoClips, &vci)
			} else {
				//delete the non-object detection containing video files
				aiRootPath := v.getAiRecordPath(sourceId)
				fullAiFileName := path.Join(aiRootPath, aiFileName)
				os.Remove(fullAiFileName)
				log.Println("an ai video file deleted: " + fullAiFileName)
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

		if clip.AiQueueModels == nil || len(clip.AiQueueModels) == 0 {
			continue
		}
		for _, queueModel := range clip.AiQueueModels {
			aiClipModel := data.AiClip{}
			aiClipModel.Setup(newLocation, clip.CreatedAt, clip.LastModified, clip.Duration)
			rep.SetAiClipFields(getId(queueModel), &aiClipModel)
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

func predicateSourceId(queueModel *models.AiClipQueueModel, sourceId string) bool {
	switch queueModel.AiType {
	case models.Od:
		return queueModel.Od.SourceId == sourceId
	case models.Fr:
		return queueModel.Fr.SourceId == sourceId
	case models.Alpr:
		return queueModel.Alpr.SourceId == sourceId
	}
	log.Println("unsupported type has been detected, value is " + strconv.Itoa(queueModel.AiType))
	return false
}

func getCreatedAt(queueModel *models.AiClipQueueModel) string {
	switch queueModel.AiType {
	case models.Od:
		return queueModel.Od.CreatedAt
	case models.Fr:
		return queueModel.Fr.CreatedAt
	case models.Alpr:
		return queueModel.Alpr.CreatedAt
	}
	log.Println("unsupported type has been detected, value is " + strconv.Itoa(queueModel.AiType))
	return ""
}

func getId(queueModel *models.AiClipQueueModel) string {
	switch queueModel.AiType {
	case models.Od:
		return queueModel.Od.Id
	case models.Fr:
		return queueModel.Fr.Id
	case models.Alpr:
		return queueModel.Alpr.Id
	}
	log.Println("unsupported type has been detected, value is " + strconv.Itoa(queueModel.AiType))
	return ""
}
