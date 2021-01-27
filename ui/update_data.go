package ui

import (
	"fmt"
	"github.com/VitJRBOG/GroupsMonitor/data_manager"
	"github.com/VitJRBOG/GroupsMonitor/tools"
	"runtime/debug"
	"strings"
)

func addNewAccessToken() {
	var a data_manager.AccessToken

	fmt.Print("--- Enter a name for the new access token and press «Enter»... ---\n> ")
	name := getDataFromUser()
	err := a.SetName(name)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "access token with this name already exists") {
			fmt.Printf("[%s] Addition a new access token: An access token with this name already exists...\n",
				tools.GetCurrentDateAndTime())
			addNewAccessToken()
			return
		} else {
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
	}

	fmt.Print("--- Enter a value for the new access token and press «Enter»... ---\n> ")
	value := getDataFromUser()
	a.SetValue(value)

	a.SaveToDB()

	fmt.Printf("[%s] Access token update: New access token added successfully...\n",
		tools.GetCurrentDateAndTime())

}

func updExistsAccessToken(accessTokenName string) {
	var a data_manager.AccessToken

	a.SelectFromDB(accessTokenName)

	fmt.Print("--- Enter a new name for the access token and press «Enter»... ---\n> ")
	name := getDataFromUser()
	err := a.SetName(name)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "access token with this name already exists") {
			fmt.Printf("[%s] Addition a new access token: An access token with this name already exists...\n",
				tools.GetCurrentDateAndTime())
			updExistsAccessToken(accessTokenName)
			return
		} else {
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
	}

	fmt.Print("--- Enter a new value for the access token and press «Enter»... ---\n> ")
	value := getDataFromUser()
	a.SetValue(value)

	a.UpdateInDB()

	fmt.Printf("[%s] Access token update: Access token updated successfully...\n",
		tools.GetCurrentDateAndTime())
}

func getDataFromUser() string {
	var userInput string
	_, err := fmt.Scan(&userInput)
	if err != nil {
		tools.WriteToLog(err, debug.Stack())
		panic(err.Error())
	}
	return userInput
}
