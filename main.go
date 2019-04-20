package main

import (
	"log"
)

func main() {
	threads := MakeThreads()
	if err := ListenUserCommands(threads); err != nil {
		log.Fatal(err)
	}
}
