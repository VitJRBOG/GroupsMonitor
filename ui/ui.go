package ui

import (
	"github.com/VitJRBOG/Watcher/observer"
	"github.com/VitJRBOG/Watcher/ui/cli"
	"github.com/VitJRBOG/Watcher/ui/webview"
)

func ShowUI(dbHasBeenInitialized bool) {
	cli.ShowDBStatus(dbHasBeenInitialized)
	params := observer.MakeObservers()
	cli.StartObservers(params)
	go cli.ListenUserCommands(params)
	go webview.InitWebview()
	cli.CheckObservers(params)
}
