package main

import (
	"fmt"
	"time"
)

// OutputMessage выводит сообщение в консоль
func OutputMessage(sender string, message string) {
	fmt.Printf("> [%v] [%v]: %v\n", UnixTimeStampToDate(int(time.Now().Unix())), sender, message)
}
