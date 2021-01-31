package upd_data

import (
	"fmt"
	"github.com/VitJRBOG/GroupsMonitor/data_manager"
	"github.com/VitJRBOG/GroupsMonitor/tools"
	"github.com/VitJRBOG/GroupsMonitor/ui/input"
	"runtime/debug"
	"strings"
)

func AddNewOperator() {
	operation := "Addition a new operator"

	var operator data_manager.Operator
	setOperatorName(&operator, operation)
	setOperatorVkID(&operator, operation)
	operator.SaveToDB()

	fmt.Printf("[%s] %s: New operator added successfully...\n",
		tools.GetCurrentDateAndTime(), operation)
}

func UpdExistsOperator(operatorName string) {
	operation := "Operator update"

	operator := selectDataOfExistsOperator(operatorName, operation)
	setOperatorName(operator, operation)
	setOperatorVkID(operator, operation)

	operator.UpdateIdDB()

	fmt.Printf("[%s] %s: Operator updated successfully...\n",
		tools.GetCurrentDateAndTime(), operation)
}

func selectDataOfExistsOperator(operatorName, operation string) *data_manager.Operator {
	var operator data_manager.Operator
	err := operator.SelectFromDB(operatorName)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "no such operator found") {
			fmt.Printf("[%s] %s: an operator with this name does not exist...\n",
				tools.GetCurrentDateAndTime(), operation)
			return nil
		} else {
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
	}
	return &operator
}

func setOperatorName(operator *data_manager.Operator, operation string) {
	fmt.Print("--- Enter a name for the new operator and press «Enter»... ---\n> ")
	name := input.GetDataFromUser()
	err := operator.SetName(name)
	if err != nil {
		switch true {
		case strings.Contains(strings.ToLower(err.Error()), "string length is zero"):
			fmt.Printf("[%s] %s: You must enter a name...\n",
				tools.GetCurrentDateAndTime(), operation)
			setOperatorName(operator, operation)
			return
		case strings.Contains(strings.ToLower(err.Error()), "operator with this name already exists"):
			fmt.Printf("[%s] %s: An operator with this name already exists...\n",
				tools.GetCurrentDateAndTime(), operation)
			setOperatorName(operator, operation)
			return
		default:
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
	}
}

func setOperatorVkID(operator *data_manager.Operator, operation string) {
	fmt.Print("--- Enter the VK ID for the new operator and press «Enter»... ---\n> ")
	strVkID := input.GetDataFromUser()
	err := operator.SetVkID(strVkID)
	if err != nil {
		switch true {
		case strings.Contains(strings.ToLower(err.Error()), "string length is zero"):
			fmt.Printf("[%s] %s: You must enter a VK ID...\n",
				tools.GetCurrentDateAndTime(), operation)
			setOperatorVkID(operator, operation)
			return
		case strings.Contains(strings.ToLower(err.Error()), "vk id starts with zero"):
			fmt.Printf("[%s] %s: VK ID should not start with zero...\n",
				tools.GetCurrentDateAndTime(), operation)
			setOperatorVkID(operator, operation)
			return
		case strings.Contains(strings.ToLower(err.Error()), "invalid syntax"):
			fmt.Printf("[%s] %s: VK ID must be integer...\n",
				tools.GetCurrentDateAndTime(), operation)
			setOperatorVkID(operator, operation)
			return
		default:
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
	}
}
