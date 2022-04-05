package utils

import (
	"fmt"
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

func StringToTime(dateString string, includeNanoSec bool) time.Time {
	splits := strings.Split(dateString, "_")

	year, _ := strconv.Atoi(splits[0])
	month, _ := strconv.Atoi(splits[1])
	day, _ := strconv.Atoi(splits[2])
	hour, _ := strconv.Atoi(splits[3])
	minute, _ := strconv.Atoi(splits[4])
	second, _ := strconv.Atoi(splits[5])
	nanoSec := 0
	if includeNanoSec {
		nanoSec, _ = strconv.Atoi(splits[6])
	}

	return time.Date(year, time.Month(month), day, hour, minute, second, nanoSec, time.UTC)
}

func TimeToString(time time.Time, includeNanoSec bool) string {
	arr := make([]string, 0)
	arr = append(arr, strconv.Itoa(time.Year()))
	arr = append(arr, strconv.Itoa(int(time.Month())))
	arr = append(arr, strconv.Itoa(time.Day()))
	arr = append(arr, strconv.Itoa(time.Hour()))
	arr = append(arr, strconv.Itoa(time.Minute()))
	arr = append(arr, strconv.Itoa(time.Second()))
	if includeNanoSec {
		arr = append(arr, strconv.Itoa(time.Nanosecond()))
	}

	return strings.Join(arr, "_")
}

func GetFileNameWithoutExtension(fileName string) string {
	fileName = filepath.Base(fileName)
	extension := filepath.Ext(fileName)
	fileName = fileName[0 : len(fileName)-len(extension)]
	return fileName
}
