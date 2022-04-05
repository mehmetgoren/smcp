package disk

import (
	"path"
	"smcp/utils"
	"strconv"
	"strings"
	"time"
)

type IndexedFolderInfo struct {
	Year  string
	Month string
	Day   string
	Hour  string
}

func (i *IndexedFolderInfo) GetIndexedPath(rootPath string) string {
	return path.Join(rootPath, i.Year, i.Month, i.Day, i.Hour)
}

type IndexedFolderInfoProvider interface {
	Create() *IndexedFolderInfo
}

type CurrentTimeIndexedFolderInfoProvider struct {
}

func timeToIndexedFolderInfo(t *time.Time) *IndexedFolderInfo {
	year := strconv.Itoa(t.Year())
	month := strconv.Itoa(int(t.Month()))
	day := strconv.Itoa(t.Day())
	hour := strconv.Itoa(t.Hour())
	return &IndexedFolderInfo{year, month, day, hour}
}

func (c CurrentTimeIndexedFolderInfoProvider) Create() *IndexedFolderInfo {
	now := time.Now()
	return timeToIndexedFolderInfo(&now)
}

type FileNameIndexedFolderInfoProvider struct {
	FileName string
}

func (f FileNameIndexedFolderInfoProvider) Create() *IndexedFolderInfo {
	t := utils.StringToTime(strings.Split(f.FileName, ".")[0], false)
	return timeToIndexedFolderInfo(&t)
}
