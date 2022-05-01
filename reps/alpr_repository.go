package reps

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"path"
	"smcp/models"
	"smcp/utils"
	"strings"
)

type AlprHandlerRepository struct {
	Config *models.Config
}

func (a *AlprHandlerRepository) Save(ar *models.AlprResponse) error {
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(ar.Base64Image))
	defer ioutil.NopCloser(reader)

	imageBytes, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Println("AlprHandlerRepository: Reading base64 message error: " + err.Error())
		return err
	}

	createdAt := utils.StringToTime(ar.CreatedAt)
	timeIndex := TimeIndex{}
	timeIndex.SetValuesFrom(&createdAt)

	//save the image first
	// creates an indexed data directory
	rootPathImage := utils.GetAlprImagesPathBySourceId(a.Config, ar.SourceId)
	fullPathImage := timeIndex.GetIndexedPath(rootPathImage)
	err = utils.CreateDirectoryIfNotExists(fullPathImage)
	if err != nil {
		log.Println("an error occurred during the creating indexed image directory, ", err)
	}

	// write a file as jpeg
	fullFileNameImage := path.Join(fullPathImage, ar.CreateFileName()+".jpg")
	err = ioutil.WriteFile(fullFileNameImage, imageBytes, 0777)
	if err != nil {
		log.Println("an error occurred during the writing image file, ", err)
	}
	log.Println("DiskEventHandler: image saved successfully as " + fullFileNameImage)
	//

	// and then json file
	// creates an indexed data directory
	rootPathData := utils.GetAlprDataPathBySourceId(a.Config, ar.SourceId)
	fullPathData := timeIndex.GetIndexedPath(rootPathData)
	err = utils.CreateDirectoryIfNotExists(fullPathData)
	if err != nil {
		log.Println("an error occurred during the creating indexed data directory, ", err)
	}

	//write a file as json
	baseObj := models.AlprJsonBaseObject{
		ImgWidth: ar.ImgWidth, ImgHeight: ar.ImgHeight, ProcessingTimeMs: ar.ProcessingTimeMs,
		Results: ar.Results, Id: ar.Id, SourceId: ar.SourceId, CreatedAt: ar.CreatedAt, AiClipEnabled: ar.AiClipEnabled,
	}
	baseObj.ImageFileName = strings.Replace(fullFileNameImage, a.Config.General.RootFolderPath+"/", "", -1)
	fullFileNameData := path.Join(fullPathData, ar.CreateFileName()+".json")
	baseObj.DataFileName = strings.Replace(fullFileNameData, a.Config.General.RootFolderPath+"/", "", -1)
	jsonObj := models.AlprJsonObject{AlprResults: &baseObj, Video: &models.VideoClipJsonObject{}}
	bytes, _ := json.Marshal(jsonObj)

	err = ioutil.WriteFile(fullFileNameData, bytes, 0777)
	if err != nil {
		log.Println("an error occurred during the writing json data file, ", err)
	}
	log.Println("DiskEventHandler: json data saved successfully as " + fullFileNameData)
	//

	return nil
}
