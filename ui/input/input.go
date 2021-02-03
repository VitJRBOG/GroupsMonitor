package input

import (
	"bufio"
	"github.com/VitJRBOG/GroupsMonitor/tools"
	"os"
	"runtime/debug"
)

func GetDataFromUser() string {
	var userInput string
	in := bufio.NewReader(os.Stdin)
	userInput, err := in.ReadString('\n')
	if err != nil {
		tools.WriteToLog(err, debug.Stack())
		panic(err.Error())
	}
	return userInput[:len(userInput)-1]
}
