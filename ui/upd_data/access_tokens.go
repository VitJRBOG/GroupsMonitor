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
	var a data_manager.AccessToken

	fmt.Print("--- Enter a name for the new access token and press «Enter»... ---\n> ")
	name := input.GetDataFromUser()
	err := a.SetName(name)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "access token with this name already exists") {
			fmt.Printf("[%s] Addition a new access token: An access token with this name already exists...\n",
				tools.GetCurrentDateAndTime())
			AddNewAccessToken()
			return
		} else {
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
	}

	fmt.Print("--- Enter a value for the new access token and press «Enter»... ---\n> ")
	value := input.GetDataFromUser()
	a.SetValue(value)

	a.SaveToDB()

	fmt.Printf("[%s] Addition a new access token: New access token added successfully...\n",
		tools.GetCurrentDateAndTime())

}

func UpdExistsAccessToken(accessTokenName string) {
	var a data_manager.AccessToken

	a.SelectFromDB(accessTokenName)

	fmt.Print("--- Enter a new name for the access token and press «Enter»... ---\n> ")
	name := input.GetDataFromUser()
	err := a.SetName(name)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "access token with this name already exists") {
			fmt.Printf("[%s] Access token update: An access token with this name already exists...\n",
				tools.GetCurrentDateAndTime())
			UpdExistsAccessToken(accessTokenName)
			return
		} else {
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
	}

	fmt.Print("--- Enter a new value for the access token and press «Enter»... ---\n> ")
	value := input.GetDataFromUser()
	a.SetValue(value)

	a.UpdateInDB()

	fmt.Printf("[%s] Access token update: Access token updated successfully...\n",
		tools.GetCurrentDateAndTime())
}
