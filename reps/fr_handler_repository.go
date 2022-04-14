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

type FrHandlerRepository struct {
	Config *models.Config
}

func (o *FrHandlerRepository) Save(fr *models.FaceRecognitionModel) error {
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(fr.Base64Image))
	defer ioutil.NopCloser(reader)
	imageBytes, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Println("DiskEventHandler: Reading base64 message error: " + err.Error())
		return err
	}
	//
	createdAt := utils.StringToTime(fr.CreatedAt)
	timeIndex := TimeIndex{}
	timeIndex.SetValuesFrom(&createdAt)

	//save the image first
	// creates an indexed data directory
	rootPathImage := utils.GetFrImagesPathBySourceId(o.Config, fr.SourceId)
	fullPathImage := timeIndex.GetIndexedPath(rootPathImage)
	err = utils.CreateDirectoryIfNotExists(fullPathImage)
	if err != nil {
		log.Println("an error occurred during the creating indexed image directory, ", err)
	}

	// write a file as jpeg
	fullFileNameImage := path.Join(fullPathImage, fr.CreateFileName()+".jpeg")
	err = ioutil.WriteFile(fullFileNameImage, imageBytes, 0777)
	if err != nil {
		log.Println("an error occurred during the writing image file, ", err)
	}
	log.Println("DiskEventHandler: image saved successfully as " + fullFileNameImage)
	//

	// and then json file
	// creates an indexed data directory
	rootPathData := utils.GetFrDataPathBySourceId(o.Config, fr.SourceId)
	fullPathData := timeIndex.GetIndexedPath(rootPathData)
	err = utils.CreateDirectoryIfNotExists(fullPathData)
	if err != nil {
		log.Println("an error occurred during the creating indexed data directory, ", err)
	}
	//write a file as json
	baseObj := models.FaceRecognitionJsonBaseObject{Id: fr.Id, SourceId: fr.SourceId, CreatedAt: fr.CreatedAt,
		DetectedFaces: fr.DetectedFaces, VideoClipEnabled: fr.VideoClipEnabled}
	baseObj.ImageFileName = strings.Replace(fullFileNameImage, o.Config.General.RootFolderPath+"/", "", -1)
	fullFileNameData := path.Join(fullPathData, fr.CreateFileName()+".json")
	baseObj.DataFileName = strings.Replace(fullFileNameData, o.Config.General.RootFolderPath+"/", "", -1)
	jsonObj := models.FaceRecognitionJsonObject{FaceRecognition: &baseObj, Video: &models.VideoClipJsonObject{}}
	bytes, _ := json.Marshal(jsonObj)

	err = ioutil.WriteFile(fullFileNameData, bytes, 0777)
	if err != nil {
		log.Println("an error occurred during the writing json data file, ", err)
	}
	log.Println("DiskEventHandler: json data saved successfully as " + fullFileNameData)
	//

	return nil
}
