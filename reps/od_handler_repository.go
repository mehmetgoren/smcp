package reps

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"path"
	"smcp/models"
	"smcp/utils"
	"strconv"
	"strings"
	"time"
)

type OdHandlerRepository struct {
	Config *models.Config
}

func (o *OdHandlerRepository) Save(de *models.ObjectDetectionModel) error {
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(de.Base64Image))
	defer ioutil.NopCloser(reader)
	imageBytes, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Println("DiskEventHandler: Reading base64 message error: " + err.Error())
		return err
	}

	createdAt := utils.StringToTime(de.CreatedAt)
	timeIndex := TimeIndex{}
	timeIndex.SetValuesFrom(&createdAt)

	//save the image first
	// creates an indexed data directory
	rootPathImage := utils.GetOdImagesPathBySourceId(o.Config, de.SourceId)
	fullPathImage := timeIndex.GetIndexedPath(rootPathImage)
	err = utils.CreateDirectoryIfNotExists(fullPathImage)
	if err != nil {
		log.Println("an error occurred during the creating indexed image directory, ", err)
	}
	//
	// write a file as jpeg
	fullFileNameImage := path.Join(fullPathImage, de.CreateFileName()+".jpg")
	err = ioutil.WriteFile(fullFileNameImage, imageBytes, 0777)
	if err != nil {
		log.Println("an error occurred during the writing image file, ", err)
	}
	log.Println("DiskEventHandler: image saved successfully as " + fullFileNameImage)
	//

	// and then json file
	// creates an indexed data directory
	rootPathData := utils.GetOdDataPathBySourceId(o.Config, de.SourceId)
	fullPathData := timeIndex.GetIndexedPath(rootPathData)
	err = utils.CreateDirectoryIfNotExists(fullPathData)
	if err != nil {
		log.Println("an error occurred during the creating indexed data directory, ", err)
	}
	//write a file as json
	baseObj := models.ObjectDetectionJsonBaseObject{Id: de.Id, SourceId: de.SourceId, CreatedAt: de.CreatedAt,
		DetectedObjects: de.DetectedObjects, AiClipEnabled: de.AiClipEnabled}
	baseObj.ImageFileName = strings.Replace(fullFileNameImage, o.Config.General.RootFolderPath+"/", "", -1)
	fullFileNameData := path.Join(fullPathData, de.CreateFileName()+".json")
	baseObj.DataFileName = strings.Replace(fullFileNameData, o.Config.General.RootFolderPath+"/", "", -1)
	jsonObj := models.ObjectDetectionJsonObject{ObjectDetection: &baseObj, Video: &models.VideoClipJsonObject{}}
	bytes, _ := json.Marshal(jsonObj)

	err = ioutil.WriteFile(fullFileNameData, bytes, 0777)
	if err != nil {
		log.Println("an error occurred during the writing json data file, ", err)
	}
	log.Println("DiskEventHandler: json data saved successfully as " + fullFileNameData)
	//

	return nil
}

type TimeIndex struct {
	Year  string
	Month string
	Day   string
	Hour  string
}

func (i *TimeIndex) SetValuesFrom(t *time.Time) *TimeIndex {
	i.Year = strconv.Itoa(t.Year())
	i.Month = strconv.Itoa(int(t.Month()))
	i.Day = strconv.Itoa(t.Day())
	i.Hour = strconv.Itoa(t.Hour())
	return i
}

func (i *TimeIndex) GetIndexedPath(rootPath string) string {
	return path.Join(rootPath, i.Year, i.Month, i.Day, i.Hour)
}
