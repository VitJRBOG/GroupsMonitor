package main

import (
	"fmt"
	"log"
	"time"
)

func main() {
	err := CheckFiles()
	if err != nil {
		errorHandler(err)
	}

	RunGui()

	// if err = ListenUserCommands(threads); err != nil {
	// 	errorHandler(err)
	// }
}

// errorHandler обработчик ошибок
func errorHandler(err error) {
	date := UnixTimeStampToDate(int(time.Now().Unix()))
	log.Fatal(fmt.Errorf("> [%v] WARNING! Error: %v", date, err))
}
