package utils

import (
	"path"
	"smcp/models"
)

func GetRecordFolderPath(config *models.Config) string {
	return path.Join(config.General.RootFolderPath, "record")
}

func GetOdFolderPath(config *models.Config) string {
	dir := path.Join(config.General.RootFolderPath, "od")
	return dir
}
