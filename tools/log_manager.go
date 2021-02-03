package tools

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
)

func LogFileInitialization() {
	path := GetPath("log.txt")
	isExist := checkLogFileExistence(path)
	if isExist {
		return
	}
	createLogFile(path)
}

func checkLogFileExistence(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func createLogFile(path string) {
	writeToTextFile(path, "")
}

func WriteToLog(err error, debugStack []byte) {
	currentTime := GetCurrentDateAndTime()
	text := fmt.Sprintf("[%s]: %s\n%s", currentTime, err.Error(), string(debugStack))
	path := GetPath("log.txt")
	textFromLogFile := readTextFile(path)
	textForWriting := fmt.Sprintf("%s\n%s", text, textFromLogFile)
	writeToTextFile(path, textForWriting)
}

func readTextFile(path string) string {
	file, err := os.Open(path)
	defer func() {
		err := file.Close()
		if err != nil {
			panic(err.Error())
		}
	}()
	if err != nil {
		panic(err.Error())
	}

	scanner := bufio.NewScanner(file)

	var text string
	for scanner.Scan() {
		text += fmt.Sprintf("%v\n", scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		panic(err.Error())
	}

	return text
}

func writeToTextFile(path, text string) {
	data := []byte(text)
	err := ioutil.WriteFile(path, data, 0644)
	if err != nil {
		panic(err.Error())
	}
}
