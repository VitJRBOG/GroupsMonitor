package handler

import (
	"net/http"
	"strconv"

	"github.com/VitJRBOG/GroupsObserver/data_manager"
	"github.com/gorilla/mux"
)

func homePageHandler(w http.ResponseWriter, _ *http.Request) {

	wards := data_manager.SelectWards()

	t := getHtmlTemplates()
	err := t.ExecuteTemplate(w, "index", wards)
	if err != nil {
		panic(err.Error())
	}
}

func wardObservationTogglePageHandler(w http.ResponseWriter, r *http.Request) {
	wardID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
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

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
