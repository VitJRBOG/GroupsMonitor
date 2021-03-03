package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/VitJRBOG/GroupsObserver/data_manager"
	"github.com/gorilla/mux"
)

func observersPageHandler(w http.ResponseWriter, r *http.Request) {

	wards := data_manager.SelectWards()

	http.Redirect(w, r, fmt.Sprintf("/observers/%d", wards[0].ID), http.StatusSeeOther)
}

type observersData struct {
	WardID           int
	Wards            []data_manager.Ward
	Observers        []data_manager.Observer
	ObserversTypesRu []string
}

func observerControlPageHandler(w http.ResponseWriter, r *http.Request) {
	wardID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		panic(err.Error())
	}

	var obs observersData

	obs.WardID = wardID

	obs.Wards = data_manager.SelectWards()

	obs.ObserversTypesRu = []string{"Посты на стене", "Комментарии на стене", "Фото в альбомах",
		"Комментарии под фото", "Видео в альбомах", "Комментарии под видео", "Обсуждения"}

	observersTypes := []string{
		"wall_post", "wall_reply", "photo", "photo_comment", "video", "video_comment", "board_post",
	}

	for _, item := range observersTypes {
		var o data_manager.Observer

		o.SelectFromDB(item, wardID)

		obs.Observers = append(obs.Observers, o)
	}

	t := getHtmlTemplates()
	err = t.ExecuteTemplate(w, "observer_control", obs)
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

	var o data_manager.Observer
	o.SelectFromDB(observerName, wardID)

	if o.UnderObservation == 1 {
		o.SetObservationFlag(false)
	} else {
		o.SetObservationFlag(true)
	}

	o.UpdateInDB()

	http.Redirect(w, r, fmt.Sprintf("/observers/%d", wardID), http.StatusSeeOther)
}
