package data

import (
	"encoding/base64"
	"io/ioutil"
	"log"
	"path"
	"smcp/models"
	"smcp/utils"
	"strings"
)

type SaveImageOptions struct {
	Config *models.Config

	Base64Image *string
	CreatedAt   string
	ImagesPath  string
	FileName    string
}

func (s *SaveImageOptions) getFullPathImagePath() ([]byte, string, error) {
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(*s.Base64Image))
	defer ioutil.NopCloser(reader)
	imageBytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, "", err
	}

	createdAt := utils.StringToTime(s.CreatedAt)
	timeIndex := utils.TimeIndex{}
	timeIndex.SetValuesFrom(&createdAt)

	//save the image first
	// creates an indexed data directory
	rootPathImage := s.ImagesPath
	fullPathImage := timeIndex.GetIndexedPath(rootPathImage)
	err = utils.CreateDirectoryIfNotExists(fullPathImage)
	if err != nil {
		return nil, "", err
	}
	//
	return imageBytes, fullPathImage, nil
}

func (s *SaveImageOptions) SaveImage() (string, error) {
	imageBytes, fullPathImage, err := s.getFullPathImagePath()
	if err != nil {
		return "", err
	}
	// write a file as jpeg
	fullFileNameImage := path.Join(fullPathImage, s.FileName+".jpg")
	err = ioutil.WriteFile(fullFileNameImage, imageBytes, 0777)
	if err != nil {
		return "", err
	}
	log.Println("DiskEventHandler: image saved successfully as " + fullFileNameImage)
	//

	return fullFileNameImage, nil
}

func (s *SaveImageOptions) MapFromAi(source *models.AiDetectionModel) {
	s.Base64Image = &source.Base64Image
	s.CreatedAt = source.CreatedAt
	s.ImagesPath = utils.GetAiImagesPathBySourceId(s.Config, source.SourceId)
	s.FileName = source.CreateFileName()
}
