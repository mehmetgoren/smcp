package utils

import (
	"fmt"
	"path"
	"path/filepath"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
)

func HandlePanic() {
	if r := recover(); r != nil {
		fmt.Println("RECOVER", r)
		debug.PrintStack()
	}
}

var sep = "_"

func StringToTime(dateString string) time.Time {
	splits := strings.Split(dateString, sep)
	l := len(splits)

	values := make([]int, 0)
	for _, split := range splits {
		value, _ := strconv.Atoi(split)
		values = append(values, value)
	}

	year := values[0]
	month := 1
	if l > 1 {
		month = values[1]
	}
	day := 1
	if l > 2 {
		day = values[2]
	}
	hour := 0
	if l > 3 {
		hour = values[3]
	}
	minute := 0
	if l > 4 {
		minute = values[4]
	}
	second := 0
	if l > 5 {
		second = values[5]
	}
	nanoSec := 0
	if l > 6 {
		nanoSec = values[6]
	}

	return time.Date(year, time.Month(month), day, hour, minute, second, nanoSec, time.UTC)
}

func TimeToString(time time.Time, includeNanoSec bool) string {
	arr := make([]string, 0)
	arr = append(arr, strconv.Itoa(time.Year()))
	arr = append(arr, fixZero(int(time.Month())))
	arr = append(arr, fixZero(time.Day()))
	arr = append(arr, fixZero(time.Hour()))
	arr = append(arr, fixZero(time.Minute()))
	arr = append(arr, fixZero(time.Second()))
	if includeNanoSec {
		arr = append(arr, fixZero(time.Nanosecond()))
	}

	return strings.Join(arr, sep)
}

func GetFileNameWithoutExtension(fileName string) string {
	fileName = filepath.Base(fileName)
	extension := filepath.Ext(fileName)
	fileName = fileName[0 : len(fileName)-len(extension)]
	return fileName
}

func fixZero(val int) string {
	if val < 9 {
		return "0" + strconv.Itoa(val)
	}
	return strconv.Itoa(val)
}

type TimeIndex struct {
	Year  string
	Month string
	Day   string
	Hour  string
}

func (i *TimeIndex) SetValuesFrom(t *time.Time) *TimeIndex {
	i.Year = strconv.Itoa(t.Year())
	i.Month = fixZero(int(t.Month()))
	i.Day = fixZero(t.Day())
	i.Hour = fixZero(t.Hour())
	return i
}

func (i *TimeIndex) GetIndexedPath(rootPath string) string {
	arr := make([]string, 0)
	arr = append(arr, rootPath)
	arr = append(arr, i.Year)
	v, _ := strconv.Atoi(i.Month)
	if v > 0 {
		arr = append(arr, i.Month)
	}
	v, _ = strconv.Atoi(i.Day)
	if v > 0 {
		arr = append(arr, i.Day)
	}
	v, _ = strconv.Atoi(i.Hour)
	arr = append(arr, i.Hour)

	return path.Join(arr...)
}
