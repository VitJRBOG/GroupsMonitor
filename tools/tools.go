package tools

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func GetPath(localPath string) string {
	pathToRootOfProgram, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err.Error())
	}
	path := filepath.FromSlash(fmt.Sprintf("%s/%s", pathToRootOfProgram, localPath))
	return path
}

func GetCurrentDateAndTime() string {
	return ConvertUnixTimeToDate(int(time.Now().Unix()))
}

func ConvertUnixTimeToDate(ut int) string {
	t := time.Unix(int64(ut), 0)
	dateFormat := "02.01.2006 15:04:05"
	date := t.Format(dateFormat)
	return date
}
