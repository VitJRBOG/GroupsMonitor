package main

import (
	"fmt"
)

// OutputMessage выводит сообщение в консоль
func OutputMessage(sender string, message string) {
	fmt.Println("COMPUTER [" + sender + "]: " + message)
}
