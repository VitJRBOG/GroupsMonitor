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
		fmt.Printf("[%s] Addition a new access token: An access token with this name already exists...\n",
			tools.GetCurrentDateAndTime())
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

	fmt.Printf("[%s] Addition a new access token: New access token added successfully...\n",
		tools.GetCurrentDateAndTime())
}

type updAccessToken db.AccessToken

func updExistsAccessToken(accessTokenName string) {
	var accessToken db.AccessToken
	accessToken.SelectByName(accessTokenName)
	if accessToken.ID != 0 {
		var a updAccessToken
		a.ID = accessToken.ID
		a.setName()
		nameIsUnique := a.uniquenessCheck(accessTokenName)
		if nameIsUnique {
			a.setValue()
			a.saveToDB()
		} else {
			fmt.Printf("[%s] Access token update: An access token with this name already exists...\n",
				tools.GetCurrentDateAndTime())
			updExistsAccessToken(accessTokenName)
		}
	} else {
		fmt.Printf("[%s] Access token update: Access token with this name was not found...\n",
			tools.GetCurrentDateAndTime())
	}
}

func (a *updAccessToken) setName() {
	fmt.Print("--- Enter a new name for the access token and press «Enter»... ---\n> ")
	var userInput string
	_, err := fmt.Scan(&userInput)
	if err != nil {
		tools.WriteToLog(err, debug.Stack())
		panic(err.Error())
	}
	a.Name = userInput
}

func (a *updAccessToken) uniquenessCheck(origName string) bool {
	accessTokens := db.SelectAccessTokens()
	nameIsUnique := true
	for _, accessToken := range accessTokens {
		if a.Name != origName {
			if accessToken.Name == a.Name {
				nameIsUnique = false
				break
			}
		}
	}
	return nameIsUnique
}

func (a *updAccessToken) setValue() {
	fmt.Print("--- Enter a new value for the access token and press «Enter»... ---\n> ")
	var userInput string
	_, err := fmt.Scan(&userInput)
	if err != nil {
		tools.WriteToLog(err, debug.Stack())
		panic(err.Error())
	}
	a.Value = userInput
}

func (a *updAccessToken) saveToDB() {
	var accessToken db.AccessToken
	accessToken.ID = a.ID
	accessToken.Name = a.Name
	accessToken.Value = a.Value

	accessToken.UpdateInDB()

	fmt.Printf("[%s] Access token update: Access token updated successfully...\n",
		tools.GetCurrentDateAndTime())
}
