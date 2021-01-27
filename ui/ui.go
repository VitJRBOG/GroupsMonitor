package ui

import (
	"bufio"
	"fmt"
	"github.com/VitJRBOG/GroupsMonitor/observer"
	"github.com/VitJRBOG/GroupsMonitor/tools"
	"os"
	"runtime/debug"
	"strings"
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
		fmt.Printf("[%s] %s is %s: «%s».\n", tools.GetCurrentDateAndTime(),
			params.Name, params.Status, msg)
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
			fmt.Printf("[%s]: All observers is stopped. Exit from program...\n",
				tools.GetCurrentDateAndTime())
			return
		}
	}
}

func listenUserCommands(params []*observer.ModuleParams) {
	for {
		var userInput string
		in := bufio.NewReader(os.Stdin)
		userInput, err := in.ReadString('\n')
		if err != nil {
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
		needExit := consoleCommandHandler(userInput[:len(userInput)-1], params)
		if needExit {
			return
		}
	}
}

func consoleCommandHandler(userInput string, params []*observer.ModuleParams) bool {
	command := strings.Split(userInput, " ")
	if len(command) > 0 {
		switch command[0] {
		case "add":
			if len(command) > 1 {
				switch command[1] {
				case "":
					fmt.Printf("[%s]: Argument is empty. Available arguments: access_token...\n",
						tools.GetCurrentDateAndTime())
				case "access_token":
					addNewAccessToken()
				default:
					fmt.Printf("[%s]: Unknown argument. Available arguments: access_token...\n",
						tools.GetCurrentDateAndTime())
				}
			} else {
				fmt.Printf("[%s]: Argument is empty. Available arguments: access_token...\n",
					tools.GetCurrentDateAndTime())
			}
		case "upd":
			if len(command) > 1 {
				switch command[1] {
				case "":
					fmt.Printf("[%s]: Argument is empty. Available arguments: access_token...\n",
						tools.GetCurrentDateAndTime())
				case "access_token":
					if len(command) > 2 {
						if len(command[2]) > 0 {
							updExistsAccessToken(command[2])
						} else {
							fmt.Printf("[%s]: Access token name is empty.\n",
								tools.GetCurrentDateAndTime())
						}
					} else {
						fmt.Printf("[%s]: Access token name is empty.\n",
							tools.GetCurrentDateAndTime())
					}
				default:
					fmt.Printf("[%s]: Argument is empty. Available arguments: access_token...\n",
						tools.GetCurrentDateAndTime())
				}
			} else {
				fmt.Printf("[%s]: Argument is empty. Available arguments: access_token...\n",
					tools.GetCurrentDateAndTime())
			}
		case "exit":
			for _, p := range params {
				if p != nil {
					p.BrakeFlag = true
				}
			}
			return true
		default:
			fmt.Printf("[%s]: Unknown command...\n",
				tools.GetCurrentDateAndTime())
		}
	}
	return false
}
