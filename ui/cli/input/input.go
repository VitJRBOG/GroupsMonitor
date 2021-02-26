package input

import (
	"bufio"
	"os"
	"runtime/debug"
	"strings"

	"github.com/VitJRBOG/GroupsObserver/tools"
)

func GetDataFromUser() string {
	var userInput string
	in := bufio.NewReader(os.Stdin)
	userInput, err := in.ReadString('\n')
	if err != nil {
		tools.WriteToLog(err, debug.Stack())
		panic(err.Error())
	}

	if len(userInput) > 0 {
		u := strings.Split(userInput, "")
		u = u[:len(u)-1]
		userInput = strings.Join(u, "")
	}

	return userInput
}
