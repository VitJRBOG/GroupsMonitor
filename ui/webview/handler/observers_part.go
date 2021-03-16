package handler

import (
	"encoding/json"
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

type observerControlPageData struct {
	WardID        int                       `json:"ward_id"`
	Ward          data_manager.Ward         `json:"ward"`
	Wards         []data_manager.Ward       `json:"wards"`
	LpApiSettings vkapi.LongPollApiSettings `json:"lp_api_settings"`
}

func observerControlGet(w http.ResponseWriter, r *http.Request) {
	var data observerControlPageData
	var err error

	data.WardID, err = strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		tools.WriteToLog(err, debug.Stack())
		panic(err.Error())
	}

	err = data.Ward.SelectFromDBByID(data.WardID)
	if err != nil {
		tools.WriteToLog(err, debug.Stack())
		panic(err.Error())
	}

	var accessToken data_manager.AccessToken
	err = accessToken.SelectFromDBByID(data.Ward.GetAccessTokenID)
	if err != nil {
		tools.WriteToLog(err, debug.Stack())
		panic(err.Error())
	}

	data.LpApiSettings = vkapi.GetLongPollSettings(accessToken.Value, data.Ward.VkID)

	data.Wards = data_manager.SelectWards()

	d, err := json.Marshal(data)
	if err != nil {
		tools.WriteToLog(err, debug.Stack())
		panic(err.Error())
	}

	_, err = w.Write(d)
	if err != nil {
		tools.WriteToLog(err, debug.Stack())
		panic(err.Error())
	}
}

func observerControlPageHandler(w http.ResponseWriter, r *http.Request) {
	wardID := mux.Vars(r)["id"]

	t := getHtmlTemplates()
	err := t.ExecuteTemplate(w, "observer_control", wardID)
	if err != nil {
		tools.WriteToLog(err, debug.Stack())
		panic(err.Error())
	}
}

func observerTogglePageHandler(w http.ResponseWriter, r *http.Request) {
	wardID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		tools.WriteToLog(err, debug.Stack())
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
}
