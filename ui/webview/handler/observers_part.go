package handler

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"strconv"

	"github.com/VitJRBOG/GroupsObserver/data_manager"
	"github.com/VitJRBOG/GroupsObserver/tools"
	"github.com/VitJRBOG/GroupsObserver/vkapi"
	"github.com/gorilla/mux"
)

func observersPageHandler(w http.ResponseWriter, r *http.Request) {

	wards := data_manager.SelectWards()

	http.Redirect(w, r, fmt.Sprintf("/observers/%d", wards[0].ID), http.StatusSeeOther)
}

func observerControlPageHandler(w http.ResponseWriter, r *http.Request) {
	wardID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		panic(err.Error())
	}

	var ward data_manager.Ward
	err = ward.SelectFromDBByID(wardID)
	if err != nil {
		tools.WriteToLog(err, debug.Stack())
		panic(err.Error())
	}

	var accessToken data_manager.AccessToken
	err = accessToken.SelectFromDBByID(ward.GetAccessTokenID)
	if err != nil {
		tools.WriteToLog(err, debug.Stack())
		panic(err.Error())
	}

	lpApiSettings := vkapi.GetLongPollSettings(accessToken.Value, ward.VkID)

	wards := data_manager.SelectWards()

	type observerControlData struct {
		WardID        int
		Ward          data_manager.Ward
		Wards         []data_manager.Ward
		AccessToken   data_manager.AccessToken
		LpApiSettings vkapi.LongPollApiSettings
	}

	ocd := observerControlData{
		WardID:        wardID,
		Ward:          ward,
		Wards:         wards,
		AccessToken:   accessToken,
		LpApiSettings: lpApiSettings,
	}

	t := getHtmlTemplates()
	err = t.ExecuteTemplate(w, "observer_control", ocd)
	if err != nil {
		panic(err.Error())
	}
}

func observerTogglePageHandler(w http.ResponseWriter, r *http.Request) {
	wardID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		panic(err.Error())
	}

	observerName := mux.Vars(r)["name"]

	observerMode := mux.Vars(r)["mode"]

	var ward data_manager.Ward
	err = ward.SelectFromDBByID(wardID)
	if err != nil {
		tools.WriteToLog(err, debug.Stack())
		panic(err.Error())
	}

	var accessToken data_manager.AccessToken
	err = accessToken.SelectFromDBByID(ward.GetAccessTokenID)
	if err != nil {
		tools.WriteToLog(err, debug.Stack())
		panic(err.Error())
	}

	newParam := map[string]string{
		"event_name": observerName,
		"mode":       observerMode,
	}

	vkapi.SetLongPollSettings(accessToken.Value, ward.VkID, newParam)

	http.Redirect(w, r, fmt.Sprintf("/observers/%d", wardID), http.StatusSeeOther)
}
