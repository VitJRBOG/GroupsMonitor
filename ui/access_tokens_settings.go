package ui

import (
	"fmt"
	"github.com/VitJRBOG/GroupsMonitor/db"
	"github.com/VitJRBOG/GroupsMonitor/tools"
	"runtime/debug"
)

func addNewAccessToken() {
	var a newAccessToken
	a.setName()
	nameIsUnique := a.uniquenessCheck()
	if nameIsUnique {
		a.setValue()
		a.saveToDB()
	} else {
		output := fmt.Sprintf("[%s] Addition a new access token: An access token with this name already exists...",
			tools.GetCurrentDateAndTime())
		fmt.Println(output)
		addNewAccessToken()
	}
}

type newAccessToken db.AccessToken

func (a *newAccessToken) setName() {
	fmt.Print("--- Enter a name for the new access token and press «Enter»... ---\n> ")
	var userInput string
	_, err := fmt.Scan(&userInput)
	if err != nil {
		tools.WriteToLog(err, debug.Stack())
		panic(err.Error())
	}
	a.Name = userInput
}

func (a *newAccessToken) uniquenessCheck() bool {
	accessTokens := db.SelectAccessTokens()
	nameIsUnique := true
	for _, accessToken := range accessTokens {
		if accessToken.Name == a.Name {
			nameIsUnique = false
			break
		}
	}
	return nameIsUnique
}

func (a *newAccessToken) setValue() {
	fmt.Print("--- Enter a value for the new access token and press «Enter»... ---\n> ")
	var userInput string
	_, err := fmt.Scan(&userInput)
	if err != nil {
		tools.WriteToLog(err, debug.Stack())
		panic(err.Error())
	}
	a.Value = userInput
}

func (a *newAccessToken) saveToDB() {
	var accessToken db.AccessToken
	accessToken.ID = a.ID
	accessToken.Name = a.Name
	accessToken.Value = a.Value

	accessToken.InsertIntoDB()

	output := fmt.Sprintf("[%s] Addition a new access token: New access token added successfully...",
		tools.GetCurrentDateAndTime())
	fmt.Println(output)
}
