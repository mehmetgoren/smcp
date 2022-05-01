package disk

import (
	"path"
	"strconv"
	"time"
)

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
