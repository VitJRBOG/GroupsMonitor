package main

import (
	"github.com/VitJRBOG/Watcher/db"
	"github.com/VitJRBOG/Watcher/tools"
	"github.com/VitJRBOG/Watcher/ui"
)

func main() {
	tools.LogFileInitialization()
	dbHasBeenInitialized := db.Initialization()
	ui.ShowUI(dbHasBeenInitialized)
}
