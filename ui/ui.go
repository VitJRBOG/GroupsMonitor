package ui

import (
	"github.com/VitJRBOG/GroupsObserver/observer"
	"github.com/VitJRBOG/GroupsObserver/ui/cli"
	"github.com/VitJRBOG/GroupsObserver/ui/webview"
)

func ShowUI(dbHasBeenInitialized bool) {
	cli.ShowDBStatus(dbHasBeenInitialized)
	params := observer.MakeObservers()
	cli.StartObservers(params)
	go cli.ListenUserCommands(params)
	go webview.InitWebview()
	cli.CheckObservers(params)
}
