package cli

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/VitJRBOG/GroupsObserver/observer"
	"github.com/VitJRBOG/GroupsObserver/tools"
	"github.com/VitJRBOG/GroupsObserver/ui/cli/input"
	"github.com/VitJRBOG/GroupsObserver/ui/cli/upd_data"
)

func ShowDBStatus(dbHasBeenInitialized bool) {
	if dbHasBeenInitialized {
		fmt.Printf("DB is empty.\n--- Press «Enter» for exit... ---")
		_ = input.GetDataFromUser()
		os.Exit(0)
	}
}

func StartObservers(params []*observer.ModuleParams) {
	fmt.Printf("[%s]: Observers in the startup process. Please stand by...\n",
		tools.GetCurrentDateAndTime())
	for _, p := range params {
		go observer.StartObserver(p)
		go receivingMessagesFromObserver(p)
	}
}

func receivingMessagesFromObserver(params *observer.ModuleParams) {
	time.Sleep(1 * time.Second)
	for {
		msg := <-params.Message
		fmt.Printf("[%s] %s is %s: «%s».\n", tools.GetCurrentDateAndTime(),
			params.Name, params.Status, msg)
		if params.Status == "stopped" {
			return
		}
	}
}

func CheckObservers(params []*observer.ModuleParams) {
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

func ListenUserCommands(params []*observer.ModuleParams) {
	for {
		userInput := input.GetDataFromUser()
		needExit := consoleCommandHandler(userInput, params)
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
					fmt.Printf("[%s]: Argument is empty. Available arguments: access_token, operator, ward...\n",
						tools.GetCurrentDateAndTime())
				case "access_token":
					upd_data.AddNewAccessToken()
				case "operator":
					upd_data.AddNewOperator()
				case "ward":
					upd_data.AddNewWard()
				default:
					fmt.Printf("[%s]: Unknown argument. Available arguments: access_token, operator, ward...\n",
						tools.GetCurrentDateAndTime())
				}
			} else {
				fmt.Printf("[%s]: Argument is empty. Available arguments: access_token, operator, ward...\n",
					tools.GetCurrentDateAndTime())
			}
		case "upd":
			if len(command) > 1 {
				switch command[1] {
				case "":
					fmt.Printf("[%s]: Argument is empty. Available arguments: "+
						"access_token, operator, ward, observer...\n",
						tools.GetCurrentDateAndTime())
				case "access_token":
					if len(command) > 2 {
						if len(command[2]) > 0 {
							upd_data.UpdExistsAccessToken(command[2])
						} else {
							fmt.Printf("[%s]: Access token name is empty.\n",
								tools.GetCurrentDateAndTime())
						}
					} else {
						fmt.Printf("[%s]: Access token name is empty.\n",
							tools.GetCurrentDateAndTime())
					}
				case "operator":
					if len(command) > 2 {
						if len(command[2]) > 0 {
							upd_data.UpdExistsOperator(command[2])
						} else {
							fmt.Printf("[%s]: Operator name is empty.\n",
								tools.GetCurrentDateAndTime())
						}
					} else {
						fmt.Printf("[%s]: Operator name is empty.\n",
							tools.GetCurrentDateAndTime())
					}
				case "ward":
					if len(command) > 2 {
						if len(command[2]) > 0 {
							upd_data.UpdExistsWard(command[2])
						} else {
							fmt.Printf("[%s]: Ward name is empty.\n",
								tools.GetCurrentDateAndTime())
						}
					} else {
						fmt.Printf("[%s]: Ward name is empty.\n",
							tools.GetCurrentDateAndTime())
					}
				case "observer":
					if len(command) > 2 {
						if len(command[2]) > 0 {
							upd_data.UpdExistsObserver(command[2])
						} else {
							fmt.Printf("[%s]: Ward name is empty.\n",
								tools.GetCurrentDateAndTime())
						}
					} else {
						fmt.Printf("[%s]: Ward name is empty.\n",
							tools.GetCurrentDateAndTime())
					}
				default:
					fmt.Printf("[%s]: Argument is empty. Available arguments: "+
						"access_token, operator, ward, observer...\n",
						tools.GetCurrentDateAndTime())
				}
			} else {
				fmt.Printf("[%s]: Argument is empty. Available arguments: "+
					"access_token, operator, ward, observer...\n",
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
