package upd_data

import (
	"fmt"
	"github.com/VitJRBOG/GroupsMonitor/data_manager"
	"github.com/VitJRBOG/GroupsMonitor/tools"
	"github.com/VitJRBOG/GroupsMonitor/ui/input"
	"runtime/debug"
	"strings"
)

func AddNewAccessToken() {
	operation := "Addition a new access token"

	var accessToken data_manager.AccessToken
	setAccessTokenName(&accessToken, operation)
	setAccessTokenValue(&accessToken, operation)
	accessToken.SaveToDB()

	fmt.Printf("[%s] %s: New access token added successfully...\n",
		tools.GetCurrentDateAndTime(), operation)
}

func UpdExistsAccessToken(accessTokenName string) {
	operation := "Access token update"

	accessToken := selectDataOfExistsAccessToken(accessTokenName, operation)
	if accessToken == nil {
		return
	}
	setAccessTokenName(accessToken, operation)
	setAccessTokenValue(accessToken, operation)
	accessToken.UpdateInDB()

	fmt.Printf("[%s] %s: Access token updated successfully...\n",
		tools.GetCurrentDateAndTime(), operation)
}

func selectDataOfExistsAccessToken(accessTokenName, operation string) *data_manager.AccessToken {
	var accessToken data_manager.AccessToken
	err := accessToken.SelectFromDB(accessTokenName)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "no such access token found") {
			fmt.Printf("[%s] %s: an access token with this name does not exist...\n",
				tools.GetCurrentDateAndTime(), operation)
			return nil
		} else {
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
	}
	return &accessToken
}

func setAccessTokenName(accessToken *data_manager.AccessToken, operation string) {
	fmt.Print("--- Enter a name for the new access token and press «Enter»... ---\n> ")
	name := input.GetDataFromUser()
	err := accessToken.SetName(name)
	if err != nil {
		switch true {
		case strings.Contains(strings.ToLower(err.Error()), "string length is zero"):
			fmt.Printf("[%s] %s: You must enter a name...\n",
				tools.GetCurrentDateAndTime(), operation)
			setAccessTokenName(accessToken, operation)
			return
		case strings.Contains(strings.ToLower(err.Error()), "access token with this name already exists"):
			fmt.Printf("[%s] %s: An access token with this name already exists...\n",
				tools.GetCurrentDateAndTime(), operation)
			setAccessTokenName(accessToken, operation)
			return
		default:
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
	}
}

func setAccessTokenValue(accessToken *data_manager.AccessToken, operation string) {
	fmt.Print("--- Enter a value for the new access token and press «Enter»... ---\n> ")
	value := input.GetDataFromUser()
	err := accessToken.SetValue(value)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "string length is zero") {
			fmt.Printf("[%s] %s: You must enter a value...\n",
				tools.GetCurrentDateAndTime(), operation)
			setAccessTokenValue(accessToken, operation)
			return
		} else {
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
	}
}
