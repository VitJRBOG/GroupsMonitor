package main

import (
	"fmt"
	"runtime"
	"time"
)

// Thread - структура для хранения данных о потоке
type Thread struct {
	Name     string
	StopFlag int
	Status   string
}

// MakeThreads создает и запускает потоки
func MakeThreads() []*Thread {
	var thread Thread
	thread.Name = "Testing"
	thread.Status = "alive"
	go testThreadint(&thread)
	var threads []*Thread
	threads = append(threads, &thread)

	//
	// тут нужен поток с функцией проверки жизни остальных потоков
	//

	return threads
}

func testThreadint(threadData *Thread) {
	fmt.Println("Okay, let's do this!")
	for true {
		interval := 10
		for i := 0; i < interval; i++ {
			time.Sleep(1 * time.Second)
			if threadData.StopFlag == 1 {
				threadData.Status = "stopped"
				fmt.Println("I'll be back...")
				runtime.Goexit()
			}
		}
		fmt.Println("Hmm... I'm stil alive...")
	}
}
