package handler

import (
	"net/http"
	"runtime/debug"
	"strconv"

	"github.com/VitJRBOG/Watcher/data_manager"
	"github.com/VitJRBOG/Watcher/tools"
	"github.com/gorilla/mux"
)

func homePageHandler(w http.ResponseWriter, _ *http.Request) {
	t := getHtmlTemplates()
	err := t.ExecuteTemplate(w, "index", nil)
	if err != nil {
		tools.WriteToLog(err, debug.Stack())
		panic(err.Error())
	}
}

func wardObservationModeSwitcher(w http.ResponseWriter, r *http.Request) {
	wardID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		tools.WriteToLog(err, debug.Stack())
		panic(err.Error())
	}

	var ward data_manager.Ward
	ward.SelectFromDBByID(wardID)

	if ward.UnderObservation == 1 {
		ward.SetObservationFlag(false)
	} else {
		ward.SetObservationFlag(true)
	}

	ward.UpdateInDB()
}
