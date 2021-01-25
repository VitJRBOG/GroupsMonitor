package main

import (
	"github.com/VitJRBOG/GroupsMonitor/db"
	"github.com/VitJRBOG/GroupsMonitor/tools"
	"github.com/VitJRBOG/GroupsMonitor/ui"
)

func main() {
	tools.LogFileInitialization()
	dbHasBeenInitialized := db.Initialization()
	ui.ShowUI(dbHasBeenInitialized)
}
