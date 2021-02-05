package ui

import (
	"github.com/VitJRBOG/GroupsObserver/observer"
	"github.com/VitJRBOG/GroupsObserver/ui/cli"
)

func ShowUI(dbHasBeenInitialized bool) {
	cli.ShowDBStatus(dbHasBeenInitialized)
	params := observer.MakeObservers()
	cli.StartObservers(params)
	go cli.ListenUserCommands(params)
	cli.CheckObservers(params)
}
