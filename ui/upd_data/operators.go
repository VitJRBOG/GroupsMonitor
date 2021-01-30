package upd_data

import (
	"fmt"
	"github.com/VitJRBOG/GroupsMonitor/data_manager"
	"github.com/VitJRBOG/GroupsMonitor/tools"
	"runtime/debug"
	"strings"
)

func AddNewOperator() {
	var o data_manager.Operator

	fmt.Print("--- Enter a name for the new operator and press «Enter»... ---\n> ")
	name := getDataFromUser()
	o.SetName(name)

	fmt.Print("--- Enter the VK ID for the new operator and press «Enter»... ---\n> ")
	strVkID := getDataFromUser()
	err := o.SetVkID(strVkID)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "invalid syntax") {
			fmt.Printf("[%s] Addition a new operator: VK ID must be integer...\n",
				tools.GetCurrentDateAndTime())
			AddNewOperator()
			return
		} else {
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
	}

	o.SaveToDB()

	fmt.Printf("[%s] Addition a new operator: New operator added successfully...\n",
		tools.GetCurrentDateAndTime())
}
