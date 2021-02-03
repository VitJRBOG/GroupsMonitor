package upd_data

import (
	"fmt"
	"github.com/VitJRBOG/GroupsObserver/data_manager"
	"github.com/VitJRBOG/GroupsObserver/tools"
	"github.com/VitJRBOG/GroupsObserver/ui/input"
	"runtime/debug"
	"strings"
)

func AddNewAccessToken() {
	activity := "Addition a new access token"

	var accessToken data_manager.AccessToken
	setAccessTokenName(&accessToken, activity)
	setAccessTokenValue(&accessToken, activity)
	accessToken.SaveToDB()

	fmt.Printf("[%s] %s: New access token added successfully...\n",
		tools.GetCurrentDateAndTime(), activity)
}

func UpdExistsAccessToken(accessTokenName string) {
	activity := "Access token update"

	accessToken := selectDataOfExistsAccessToken(accessTokenName, activity)
	if accessToken == nil {
		return
	}
	setAccessTokenName(accessToken, activity)
	setAccessTokenValue(accessToken, activity)
	accessToken.UpdateInDB()

	fmt.Printf("[%s] %s: Access token updated successfully...\n",
		tools.GetCurrentDateAndTime(), activity)
}

func selectDataOfExistsAccessToken(accessTokenName, activity string) *data_manager.AccessToken {
	var accessToken data_manager.AccessToken
	err := accessToken.SelectFromDB(accessTokenName)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "no such access token found") {
			fmt.Printf("[%s] %s: an access token with this name does not exist...\n",
				tools.GetCurrentDateAndTime(), activity)
			return nil
		} else {
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
	}
	return &accessToken
}

func setAccessTokenName(accessToken *data_manager.AccessToken, activity string) {
	fmt.Print("--- Enter a name for the new access token and press «Enter»... ---\n> ")
	name := input.GetDataFromUser()
	err := accessToken.SetName(name)
	if err != nil {
		switch true {
		case strings.Contains(strings.ToLower(err.Error()), "string length is zero"):
			fmt.Printf("[%s] %s: You must enter a name...\n",
				tools.GetCurrentDateAndTime(), activity)
			setAccessTokenName(accessToken, activity)
			return
		case strings.Contains(strings.ToLower(err.Error()), "access token with this name already exists"):
			fmt.Printf("[%s] %s: An access token with this name already exists...\n",
				tools.GetCurrentDateAndTime(), activity)
			setAccessTokenName(accessToken, activity)
			return
		default:
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
	}
}

func setAccessTokenValue(accessToken *data_manager.AccessToken, activity string) {
	fmt.Print("--- Enter a value for the new access token and press «Enter»... ---\n> ")
	value := input.GetDataFromUser()
	err := accessToken.SetValue(value)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "string length is zero") {
			fmt.Printf("[%s] %s: You must enter a value...\n",
				tools.GetCurrentDateAndTime(), activity)
			setAccessTokenValue(accessToken, activity)
			return
		} else {
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
	}
}
