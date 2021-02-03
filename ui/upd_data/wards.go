package upd_data

import (
	"fmt"
	"github.com/VitJRBOG/GroupsObserver/data_manager"
	"github.com/VitJRBOG/GroupsObserver/tools"
	"github.com/VitJRBOG/GroupsObserver/ui/input"
	"runtime/debug"
	"strings"
)

func AddNewWard() {
	activity := "Addition a new ward"

	var ward data_manager.Ward
	setWardName(&ward, activity)
	setWardVkID(&ward, activity)
	setWardGetAccessToken(&ward, activity)
	ward.SaveToDB()

	var w data_manager.Ward
	err := w.SelectFromDB(ward.Name)
	if err != nil {
		tools.WriteToLog(err, debug.Stack())
		panic(err.Error())
	}
	addNewObserver(&w)

	fmt.Printf("[%s] %s: New ward added successfully...\n",
		tools.GetCurrentDateAndTime(), activity)
}

func UpdExistsWard(wardName string) {
	activity := "Ward update"

	ward := selectDataOfExistsWard(wardName, activity)
	setWardName(ward, activity)
	setWardVkID(ward, activity)
	setWardGetAccessToken(ward, activity)

	ward.UpdateInDB()

	fmt.Printf("[%s] %s: Ward updated successfully...\n",
		tools.GetCurrentDateAndTime(), activity)
}

func selectDataOfExistsWard(wardName, activity string) *data_manager.Ward {
	var ward data_manager.Ward
	err := ward.SelectFromDB(wardName)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "no such ward found") {
			fmt.Printf("[%s] %s: a ward with this name does not exist...\n",
				tools.GetCurrentDateAndTime(), activity)
			return nil
		} else {
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
	}
	return &ward
}

func setWardName(ward *data_manager.Ward, activity string) {
	fmt.Print("--- Enter a name for the new ward and press «Enter»... ---\n> ")
	name := input.GetDataFromUser()
	err := ward.SetName(name)
	if err != nil {
		switch true {
		case strings.Contains(strings.ToLower(err.Error()), "string length is zero"):
			fmt.Printf("[%s] %s: You must enter a name...\n",
				tools.GetCurrentDateAndTime(), activity)
			setWardName(ward, activity)
			return
		case strings.Contains(strings.ToLower(err.Error()), "ward with this name already exists"):
			fmt.Printf("[%s] %s: A ward with this name already exists...\n",
				tools.GetCurrentDateAndTime(), activity)
			setWardName(ward, activity)
			return
		default:
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
	}
}

func setWardVkID(ward *data_manager.Ward, activity string) {
	fmt.Print("--- Enter the VK ID for the new ward and press «Enter»... ---\n> ")
	strVkID := input.GetDataFromUser()
	err := ward.SetVkID(strVkID)
	if err != nil {
		switch true {
		case strings.Contains(strings.ToLower(err.Error()), "string length is zero"):
			fmt.Printf("[%s] %s: You must enter a VK ID...\n",
				tools.GetCurrentDateAndTime(), activity)
			setWardVkID(ward, activity)
			return
		case strings.Contains(strings.ToLower(err.Error()), "vk id starts with zero"):
			fmt.Printf("[%s] %s: VK ID should not start with zero...\n",
				tools.GetCurrentDateAndTime(), activity)
			setWardVkID(ward, activity)
			return
		case strings.Contains(strings.ToLower(err.Error()), "invalid syntax"):
			fmt.Printf("[%s] %s: VK ID must be integer...\n",
				tools.GetCurrentDateAndTime(), activity)
			setWardVkID(ward, activity)
			return
		case strings.Contains(strings.ToLower(err.Error()), "vk group id positive"):
			fmt.Printf("[%s] %s: VK group ID must be negative...\n",
				tools.GetCurrentDateAndTime(), activity)
			setWardVkID(ward, activity)
			return
		default:
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
	}
}

func setWardGetAccessToken(ward *data_manager.Ward, activity string) {
	fmt.Print("--- Enter the name of the access key to new ward and press «Enter»... ---\n> ")
	accessTokenName := input.GetDataFromUser()
	err := ward.SetAccessToken(accessTokenName)
	if err != nil {
		switch true {
		case strings.Contains(strings.ToLower(err.Error()), "string length is zero"):
			fmt.Printf("[%s] %s: You must enter a name of the access token...\n",
				tools.GetCurrentDateAndTime(), activity)
			setWardGetAccessToken(ward, activity)
			return
		case strings.Contains(strings.ToLower(err.Error()), "no such access token found"):
			fmt.Printf("[%s] %s: an access token with this name does not exist...\n",
				tools.GetCurrentDateAndTime(), activity)
			setWardGetAccessToken(ward, activity)
			return
		default:
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
	}
}
