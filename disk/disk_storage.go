package disk

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

type Storage[T any] struct {
	Provider FeatureProvider[T]
}

func (s *Storage[T]) Save(config *models.Config, model interface{}) error {
	s.Provider.SetConfigAndModel(config, model)

	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(*s.Provider.GetBase64Image()))
	defer ioutil.NopCloser(reader)
	imageBytes, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Println("AlprHandlerRepository: Reading base64 message error: " + err.Error())
		return err
	}

	createdAt := utils.StringToTime(s.Provider.GetCreatedAt())
	timeIndex := utils.TimeIndex{}
	timeIndex.SetValuesFrom(&createdAt)

	//save the image first
	// creates an indexed data directory
	rootPathImage := s.Provider.GetImagesPathBySourceId()
	fullPathImage := timeIndex.GetIndexedPath(rootPathImage)
	err = utils.CreateDirectoryIfNotExists(fullPathImage)
	if err != nil {
		log.Println("an error occurred during the creating indexed image directory, ", err)
	}
	//
	// write a file as jpeg
	fullFileNameImage := path.Join(fullPathImage, s.Provider.GetFileName()+".jpg")
	err = ioutil.WriteFile(fullFileNameImage, imageBytes, 0777)
	if err != nil {
		log.Println("an error occurred during the writing image file, ", err)
	}
	log.Println("DiskEventHandler: image saved successfully as " + fullFileNameImage)
	//

	// and then json file
	// creates an indexed data directory
	rootPathData := s.Provider.GetDataPathBySourceId()
	fullPathData := timeIndex.GetIndexedPath(rootPathData)
	err = utils.CreateDirectoryIfNotExists(fullPathData)
	if err != nil {
		log.Println("an error occurred during the creating indexed data directory, ", err)
	}
	//write a file as json
	jsonObj := s.Provider.CreateJsonObject()
	s.Provider.SetImageFileName(jsonObj, strings.Replace(fullFileNameImage, config.General.RootFolderPath+"/", "", -1))
	fullFileNameData := path.Join(fullPathData, s.Provider.GetFileName()+".json")
	s.Provider.SetDataFileName(jsonObj, strings.Replace(fullFileNameData, config.General.RootFolderPath+"/", "", -1))
	bytes, _ := json.Marshal(jsonObj)
	err = ioutil.WriteFile(fullFileNameData, bytes, 0777)
	if err != nil {
		log.Println("an error occurred during the writing json data file, ", err)
	}
	log.Println("DiskEventHandler: json data saved successfully as " + fullFileNameData)
	//

	return err
}
