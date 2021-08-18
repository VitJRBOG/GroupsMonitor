package upd_data

import (
	"fmt"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/VitJRBOG/Watcher/data_manager"
	"github.com/VitJRBOG/Watcher/tools"
	"github.com/VitJRBOG/Watcher/ui/cli/input"
)

func addNewObserver(ward *data_manager.Ward) {
	observersTypes := []string{
		"wall_post", "wall_reply", "photo", "photo_comment", "video", "video_comment", "board_post",
	}
	for _, observerType := range observersTypes {
		activity := fmt.Sprintf("Addition a new observer of the %s",
			strings.ReplaceAll(observerType, "_", " "))
		fmt.Printf("[%s] %s:...\n", tools.GetCurrentDateAndTime(), activity)

		var observer data_manager.Observer
		observer.SetName(observerType)
		observer.SetWardID(ward.ID)
		setObserverOperator(&observer, activity)
		setObserverSendAccessToken(&observer, activity)
		if observerType == "wall_post" {
			setObserverWallPostAdditionalParams(&observer, activity)
		}
		observer.SaveToDB()
	}
}

func UpdExistsObserver(wardName string) {
	activity := "Observer update"

	ward := selectDataOfExistsWard(wardName, activity)
	observer := selectDataOfExistsObserver(ward, activity)
	setObserverOperator(observer, activity)
	setObserverSendAccessToken(observer, activity)
	if observer.Name == "wall_post" {
		setObserverWallPostAdditionalParams(observer, activity)
	}
	observer.UpdateInDB()

	fmt.Printf("[%s] %s: Observer updated successfully...\n",
		tools.GetCurrentDateAndTime(), activity)
}

func selectDataOfExistsObserver(ward *data_manager.Ward, activity string) *data_manager.Observer {
	var observer data_manager.Observer
	fmt.Print("Observers of this ward:\n")
	observersTypes := []string{
		"wall_post", "wall_reply", "photo", "photo_comment", "video", "video_comment", "board_post",
	}
	for i, item := range observersTypes {
		fmt.Printf("%d - %s\n", i+1, item)
	}
	fmt.Print("--- Enter the number of name the observer of this ward and press «Enter»... ---\n> ")
	strObserverTypeNumber := input.GetDataFromUser()
	observerTypeNumber, err := strconv.Atoi(strObserverTypeNumber)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "invalid syntax") {
			fmt.Printf("[%s] %s: Number of name the observer must be integer...\n",
				tools.GetCurrentDateAndTime(), activity)
			return selectDataOfExistsObserver(ward, activity)
		} else {
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
	}
	if observerTypeNumber > 0 && observerTypeNumber <= len(observersTypes) {
		err := observer.SelectFromDB(observersTypes[observerTypeNumber-1], ward.ID)
		if err != nil {
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
	} else {
		fmt.Printf("[%s] %s: Number of name the observer must be in the range from 1 to %d...\n",
			tools.GetCurrentDateAndTime(), activity, len(observersTypes))
		return selectDataOfExistsObserver(ward, activity)
	}
	return &observer
}

func setObserverOperator(observer *data_manager.Observer, activity string) {
	fmt.Print("--- Enter the name of the operator for new observer and press «Enter»... ---\n> ")
	operatorName := input.GetDataFromUser()
	err := observer.SetOperator(operatorName)
	if err != nil {
		switch true {
		case strings.Contains(strings.ToLower(err.Error()), "string length is zero"):
			fmt.Printf("[%s] %s: You must enter a name...\n",
				tools.GetCurrentDateAndTime(), activity)
			setObserverOperator(observer, activity)
			return
		case strings.Contains(strings.ToLower(err.Error()), "no such operator found"):
			fmt.Printf("[%s] %s: an operator with this name does not exist...\n",
				tools.GetCurrentDateAndTime(), activity)
			setObserverOperator(observer, activity)
			return
		default:
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
	}
}

func setObserverSendAccessToken(observer *data_manager.Observer, activity string) {
	fmt.Print("--- Enter the name of the access token for sending messages to this operator and press «Enter»... ---\n> ")
	accessTokenName := input.GetDataFromUser()
	err := observer.SetAccessToken(accessTokenName)
	if err != nil {
		switch true {
		case strings.Contains(strings.ToLower(err.Error()), "string length is zero"):
			fmt.Printf("[%s] %s: You must enter a name...\n",
				tools.GetCurrentDateAndTime(), activity)
			setObserverSendAccessToken(observer, activity)
			return
		case strings.Contains(strings.ToLower(err.Error()), "no such access token found"):
			fmt.Printf("[%s] %s: an access token with this name does not exist...\n",
				tools.GetCurrentDateAndTime(), activity)
			setObserverSendAccessToken(observer, activity)
			return
		default:
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
	}
}

func setObserverWallPostAdditionalParams(observer *data_manager.Observer, activity string) {
	fmt.Print("Filters of the wall posts VK:\n")
	filters := []string{
		"post", "suggest", "postponed",
	}
	for i, item := range filters {
		fmt.Printf("%d - %s\n", i+1, item)
	}
	fmt.Print("--- Enter the number of name the filter of wall posts for this observer and press «Enter»... ---\n> ")
	strFilterNumber := input.GetDataFromUser()
	filterNumber, err := strconv.Atoi(strFilterNumber)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "invalid syntax") {
			fmt.Printf("[%s] %s: Number of name the filter must be integer...\n",
				tools.GetCurrentDateAndTime(), activity)
			setObserverWallPostAdditionalParams(observer, activity)
			return
		} else {
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
	}
	if filterNumber > 0 && filterNumber <= len(filters) {
		observer.SetAdditionalParams(filters[filterNumber-1])
	} else {
		fmt.Printf("[%s] %s: Number of name the filter must be in the range from 1 to %d...\n",
			tools.GetCurrentDateAndTime(), activity, len(filters))
		setObserverWallPostAdditionalParams(observer, activity)
		return
	}
}
