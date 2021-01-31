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
	var o data_manager.Operator

	fmt.Print("--- Enter a name for the new operator and press «Enter»... ---\n> ")
	name := input.GetDataFromUser()
	err := o.SetName(name)
	if err != nil {
		switch true {
		case strings.Contains(strings.ToLower(err.Error()), "string length is zero"):
			fmt.Printf("[%s] Addition a new operator: You must enter a name...\n",
				tools.GetCurrentDateAndTime())
			AddNewOperator()
			return
		case strings.Contains(strings.ToLower(err.Error()), "operator with this name already exists"):
			fmt.Printf("[%s] Addition a new operator: An operator with this name already exists...\n",
				tools.GetCurrentDateAndTime())
			AddNewOperator()
			return
		default:
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
	}

	fmt.Print("--- Enter the VK ID for the new operator and press «Enter»... ---\n> ")
	strVkID := input.GetDataFromUser()
	err = o.SetVkID(strVkID)
	if err != nil {
		switch true {
		case strings.Contains(strings.ToLower(err.Error()), "string length is zero"):
			fmt.Printf("[%s] Addition a new operator: You must enter a VK ID...\n",
				tools.GetCurrentDateAndTime())
			AddNewOperator()
			return
		case strings.Contains(strings.ToLower(err.Error()), "vk id starts with zero"):
			fmt.Printf("[%s] Addition a new operator: VK ID should not start with zero...\n",
				tools.GetCurrentDateAndTime())
			AddNewOperator()
			return
		case strings.Contains(strings.ToLower(err.Error()), "invalid syntax"):
			fmt.Printf("[%s] Addition a new operator: VK ID must be integer...\n",
				tools.GetCurrentDateAndTime())
			AddNewOperator()
			return
		default:
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
	}

	o.SaveToDB()

	fmt.Printf("[%s] Addition a new operator: New operator added successfully...\n",
		tools.GetCurrentDateAndTime())
}

func UpdExistsOperator(operatorName string) {
	var o data_manager.Operator

	o.SelectFromDB(operatorName)

	fmt.Print("--- Enter a new name for the operator and press «Enter»... ---\n> ")
	name := input.GetDataFromUser()
	err := o.SetName(name)
	if err != nil {
		switch true {
		case strings.Contains(strings.ToLower(err.Error()), "string length is zero"):
			fmt.Printf("[%s] Operator update: You must enter a name...\n",
				tools.GetCurrentDateAndTime())
			UpdExistsOperator(operatorName)
			return
		case strings.Contains(strings.ToLower(err.Error()), "operator with this name already exists"):
			fmt.Printf("[%s] Operator update: An operator with this name already exists...\n",
				tools.GetCurrentDateAndTime())
			UpdExistsOperator(operatorName)
			return
		default:
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
	}

	fmt.Print("--- Enter a new VK ID for the operator and press «Enter»... ---\n> ")
	strVkID := input.GetDataFromUser()
	err = o.SetVkID(strVkID)
	if err != nil {
		switch true {
		case strings.Contains(strings.ToLower(err.Error()), "string length is zero"):
			fmt.Printf("[%s] Operator update: You must enter a VK ID...\n",
				tools.GetCurrentDateAndTime())
			UpdExistsOperator(operatorName)
			return
		case strings.Contains(strings.ToLower(err.Error()), "vk id starts with zero"):
			fmt.Printf("[%s] Operator update: VK ID should not start with zero...\n",
				tools.GetCurrentDateAndTime())
			UpdExistsOperator(operatorName)
			return
		case strings.Contains(strings.ToLower(err.Error()), "invalid syntax"):
			fmt.Printf("[%s] Operator update: VK ID must be integer...\n",
				tools.GetCurrentDateAndTime())
			UpdExistsOperator(operatorName)
			return
		default:
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
	}

	o.UpdateIdDB()

	fmt.Printf("[%s] Operator update: Operator updated successfully...\n",
		tools.GetCurrentDateAndTime())
}
