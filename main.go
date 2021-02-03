package main

import (
	"github.com/VitJRBOG/GroupsObserver/db"
	"github.com/VitJRBOG/GroupsObserver/tools"
	"github.com/VitJRBOG/GroupsObserver/ui"
)

func main() {
	tools.LogFileInitialization()
	dbHasBeenInitialized := db.Initialization()
	ui.ShowUI(dbHasBeenInitialized)
}
