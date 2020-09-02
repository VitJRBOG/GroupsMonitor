package main

import (
	"fmt"
	"log"
	"time"
)

// OutputMessage выводит сообщение в консоль
func OutputMessage(sender string, message string) {
	fmt.Printf("> [%v] [%v]: %v\n", UnixTimeStampToDate(int(time.Now().Unix())), sender, message)
	OutputToTextFile(sender, message)
}

// OutputToTextFile сохраняет сообщение в текстовый файл
func OutputToTextFile(sender string, message string) {
	textToOutput := fmt.Sprintf("> [%v] [%v]: %v\n", UnixTimeStampToDate(int(time.Now().Unix())), sender, message)

	logText, err := ReadTextFile("log.txt")
	if err != nil {
		date := UnixTimeStampToDate(int(time.Now().Unix()))
		log.Fatal(fmt.Errorf("> [%v] WARNING! Error: %v", date, err))
	}

	textToWrite := logText + textToOutput

	err = WriteTextFile("log.txt", textToWrite)
	if err != nil {
		date := UnixTimeStampToDate(int(time.Now().Unix()))
		log.Fatal(fmt.Errorf("> [%v] WARNING! Error: %v", date, err))
	}
}
