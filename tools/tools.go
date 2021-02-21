package tools

import (
	"fmt"
	"io/ioutil"
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

func ReadJSON(path string) ([]byte, error) {
	rawData, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return rawData, nil
}

func WriteJSON(path string, raw []byte) error {
	err := ioutil.WriteFile(path, raw, 0644)
	if err != nil {
		return err
	}

	return nil
}
