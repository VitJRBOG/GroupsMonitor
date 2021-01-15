package main

import (
	"github.com/VitJRBOG/GroupsMonitor_new/db"
	"github.com/VitJRBOG/GroupsMonitor_new/tools"
	"github.com/VitJRBOG/GroupsMonitor_new/ui"
)

func main() {
	tools.LogFileInitialization()
	dbHasBeenInitialized := db.Initialization()
	ui.ShowUI(dbHasBeenInitialized)
}
