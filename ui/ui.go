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
	listControllerParams := observer.RunControllerObservers()
	time.Sleep(1 * time.Second) // чтобы программа успела в фоне создать структуры с параметрами
	go listenUserCommands(listControllerParams)
	checkObserversStatus(listControllerParams)
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

func checkObserversStatus(l []*observer.ControllerParams) {
	firstOutputObserversStatus(l)
	for true {
		allIsNil := true
		for i := 0; i < len(l); i++ {
			if l[i] != nil {
				allIsNil = false

				if len(l[i].Message) > 0 {
					outputMessage := fmt.Sprintf("[%s] %s: «%s». %s is %s.",
						tools.GetCurrentDateAndTime(), l[i].Name, l[i].Message, l[i].Name, l[i].Status)
					fmt.Println(outputMessage)
					l[i].Message = ""
				}

				if l[i].Status == "stopped" {
					l[i] = nil
				}
			}
		}
		if allIsNil {
			fmt.Println("All observers was stopped. Exit...")
			return
		}
		time.Sleep(5 * time.Second)
	}
}

func firstOutputObserversStatus(listControllerParams []*observer.ControllerParams) {
	for _, item := range listControllerParams {
		outputMessage := fmt.Sprintf("%s is %s.", item.Name, item.Status)
		fmt.Println(outputMessage)
	}
}

func listenUserCommands(listControllerParams []*observer.ControllerParams) {
	for true {
		fmt.Print("> ")
		var userInput string
		_, err := fmt.Scan(&userInput)
		if err != nil {
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
		success := consoleCommandHandler(userInput, listControllerParams)
		if success {
			return
		}
	}
}

func consoleCommandHandler(userInput string, listControllerParams []*observer.ControllerParams) bool {
	switch userInput {
	case "":
		return false
	case "exit":
		for _, item := range listControllerParams {
			if item != nil {
				item.BrakeFlag = true
			}
		}
		return true
	default:
		fmt.Println("Unknown command...")
		return false
	}
}
