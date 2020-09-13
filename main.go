package main

import (
	"fmt"
	"log"
	"time"
)

func main() {
	err := CheckFiles()
	if err != nil {
		date := UnixTimeStampToDate(int(time.Now().Unix()))
		log.Fatal(fmt.Errorf("> [%v] WARNING! Error: %v", date, err))
	}

	RunGui()
}
