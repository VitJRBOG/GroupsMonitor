package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/VitJRBOG/GroupsObserver/data_manager"
	"github.com/gorilla/mux"
)

func observersPageHandler(w http.ResponseWriter, _ *http.Request) {

	wards := data_manager.SelectWards()

	t := getHtmlTemplates()
	err := t.ExecuteTemplate(w, "observers", wards)
	if err != nil {
		panic(err.Error())
	}
}

func observerControlPageHandler(w http.ResponseWriter, r *http.Request) {
	wardID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		panic(err.Error())
	}

	var ward data_manager.Ward
	ward.SelectFromDBByID(wardID)

	t := getHtmlTemplates()
	err = t.ExecuteTemplate(w, "observer_control", ward)
	if err != nil {
		panic(err.Error())
	}
}

func observerTurnOnPageHandler(w http.ResponseWriter, r *http.Request) {
	wardID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		panic(err.Error())
	}

	var ward data_manager.Ward
	ward.SelectFromDBByID(wardID)
	ward.SetObservationFlag(true)

	ward.UpdateInDB()

	http.Redirect(w, r, fmt.Sprintf("/observers/%d", ward.ID), http.StatusSeeOther)
}

func observerTurnOffPageHandler(w http.ResponseWriter, r *http.Request) {
	wardID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		panic(err.Error())
	}

	var ward data_manager.Ward
	ward.SelectFromDBByID(wardID)
	ward.SetObservationFlag(false)

	ward.UpdateInDB()

	http.Redirect(w, r, fmt.Sprintf("/observers/%d", ward.ID), http.StatusSeeOther)
}
