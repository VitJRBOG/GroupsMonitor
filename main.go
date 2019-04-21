package main

import (
	"log"
)

func main() {
	threads, err := MakeThreads()
	if err != nil {
		ErrorHandler(err)
	}
	if err := ListenUserCommands(threads); err != nil {
		ErrorHandler(err)
	}
}

func ErrorHandler(err error) {
	log.Fatal(err)
}
