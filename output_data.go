package main

import (
	"fmt"
	"time"
)

// OutputMessage выводит сообщение в консоль
func OutputMessage(sender string, message string) {
	fmt.Printf("> [%v] [%v]: %v\n", UnixTimeStampToDate(int(time.Now().Unix())), sender, message)
}

// ToLogFile сохраняет сообщение в текстовый файл
func ToLogFile(errorMessage, trace string) {
	textToOutput := fmt.Sprintf("[%v]: %v\n%v\n", UnixTimeStampToDate(int(time.Now().Unix())), errorMessage, trace)

	logText, err := ReadTextFile("log.txt")
	if err != nil {
		panic(err.Error())
	}

	textToWrite := logText + textToOutput

	err = WriteTextFile("log.txt", textToWrite)
	if err != nil {
		panic(err.Error())
	}
}
