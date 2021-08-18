package upd_data

import (
	"fmt"
	"runtime/debug"
	"strings"

	"github.com/VitJRBOG/Watcher/data_manager"
	"github.com/VitJRBOG/Watcher/tools"
	"github.com/VitJRBOG/Watcher/ui/cli/input"
)

func AddNewOperator() {
	activity := "Addition a new operator"

	var operator data_manager.Operator
	setOperatorName(&operator, activity)
	setOperatorVkID(&operator, activity)
	operator.SaveToDB()

	fmt.Printf("[%s] %s: New operator added successfully...\n",
		tools.GetCurrentDateAndTime(), activity)
}

func UpdExistsOperator(operatorName string) {
	activity := "Operator update"

	operator := selectDataOfExistsOperator(operatorName, activity)
	setOperatorName(operator, activity)
	setOperatorVkID(operator, activity)

	operator.UpdateIdDB()

	fmt.Printf("[%s] %s: Operator updated successfully...\n",
		tools.GetCurrentDateAndTime(), activity)
}

func selectDataOfExistsOperator(operatorName, activity string) *data_manager.Operator {
	var operator data_manager.Operator
	err := operator.SelectFromDB(operatorName)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "no such operator found") {
			fmt.Printf("[%s] %s: an operator with this name does not exist...\n",
				tools.GetCurrentDateAndTime(), activity)
			return nil
		} else {
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
	}
	return &operator
}

func setOperatorName(operator *data_manager.Operator, activity string) {
	fmt.Print("--- Enter a name for the new operator and press «Enter»... ---\n> ")
	name := input.GetDataFromUser()
	err := operator.SetName(name)
	if err != nil {
		switch true {
		case strings.Contains(strings.ToLower(err.Error()), "string length is zero"):
			fmt.Printf("[%s] %s: You must enter a name...\n",
				tools.GetCurrentDateAndTime(), activity)
			setOperatorName(operator, activity)
			return
		case strings.Contains(strings.ToLower(err.Error()), "operator with this name already exists"):
			fmt.Printf("[%s] %s: An operator with this name already exists...\n",
				tools.GetCurrentDateAndTime(), activity)
			setOperatorName(operator, activity)
			return
		default:
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
	}
}

func setOperatorVkID(operator *data_manager.Operator, activity string) {
	fmt.Print("--- Enter the VK ID for the new operator and press «Enter»... ---\n> ")
	strVkID := input.GetDataFromUser()
	err := operator.SetVkID(strVkID)
	if err != nil {
		switch true {
		case strings.Contains(strings.ToLower(err.Error()), "string length is zero"):
			fmt.Printf("[%s] %s: You must enter a VK ID...\n",
				tools.GetCurrentDateAndTime(), activity)
			setOperatorVkID(operator, activity)
			return
		case strings.Contains(strings.ToLower(err.Error()), "vk id starts with zero"):
			fmt.Printf("[%s] %s: VK ID should not start with zero...\n",
				tools.GetCurrentDateAndTime(), activity)
			setOperatorVkID(operator, activity)
			return
		case strings.Contains(strings.ToLower(err.Error()), "invalid syntax"):
			fmt.Printf("[%s] %s: VK ID must be integer...\n",
				tools.GetCurrentDateAndTime(), activity)
			setOperatorVkID(operator, activity)
			return
		default:
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
	}
}
