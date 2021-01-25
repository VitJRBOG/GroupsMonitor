package ui

import (
	"fmt"
	"github.com/VitJRBOG/GroupsMonitor_new/observer"
	"github.com/VitJRBOG/GroupsMonitor_new/tools"
	"runtime/debug"
	"time"
)

func ShowUI(dbHasBeenInitialized bool) {
	showDBStatus(dbHasBeenInitialized)
	params := observer.MakeObservers()
	startObservers(params)
	go listenUserCommands(params)
	checkObservers(params)
}

func showDBStatus(dbHasBeenInitialized bool) {
	if dbHasBeenInitialized {
		fmt.Printf("DB is empty.\n--- Press «Enter for exit... ---»")
		_, err := fmt.Scan()
		if err != nil {
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
	}
}

func startObservers(params []*observer.ModuleParams) {
	for _, p := range params {
		go observer.StartObserver(p)
		go receivingMessagesFromObserver(p)
	}
}

func receivingMessagesFromObserver(params *observer.ModuleParams) {
	for {
		msg := <-params.Message
		output := fmt.Sprintf("[%s] %s: «%s». %s is %s.", tools.GetCurrentDateAndTime(),
			params.Name, msg, params.Name, params.Status)
		fmt.Println(output)
		if params.Status == "stopped" {
			return
		}
	}
}

func checkObservers(params []*observer.ModuleParams) {
	for {
		time.Sleep(3 * time.Second)
		allObserversIsStopped := true
		for _, p := range params {
			if p.Status != "stopped" {
				allObserversIsStopped = false
				break
			}
		}
		if allObserversIsStopped {
			output := fmt.Sprintf("[%s]: All observers is stopped. Exit from program...",
				tools.GetCurrentDateAndTime())
			fmt.Println(output)
			return
		}
	}
}

func listenUserCommands(params []*observer.ModuleParams) {
	for {
		var userInput string
		_, err := fmt.Scan(&userInput)
		if err != nil {
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
		success := consoleCommandHandler(userInput, params)
		if success {
			return
		}
	}
}

func consoleCommandHandler(userInput string, params []*observer.ModuleParams) bool {
	switch userInput {
	case "":
		return false
	case "exit":
		for _, p := range params {
			if p != nil {
				p.BrakeFlag = true
			}
		}
		return true
	default:
		fmt.Println("Unknown command...")
		return false
	}
}
