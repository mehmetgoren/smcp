package utils

import (
	"fmt"
	"path/filepath"
	"runtime/debug"
)

func HandlePanic() {
	if r := recover(); r != nil {
		fmt.Println("RECOVER", r)
		debug.PrintStack()
	}
}

func GetFileNameWithoutExtension(fileName string) string {
	fileName = filepath.Base(fileName)
	extension := filepath.Ext(fileName)
	fileName = fileName[0 : len(fileName)-len(extension)]
	return fileName
}
