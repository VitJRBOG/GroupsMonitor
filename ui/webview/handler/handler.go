package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/VitJRBOG/GroupsObserver/data_manager"
	"github.com/VitJRBOG/GroupsObserver/tools"
	"github.com/gorilla/mux"
)

func InitHandler() {
	rtr := mux.NewRouter()
	rtr.HandleFunc("/", homePageHandler).Methods("GET")
	rtr.HandleFunc("/observers", observersPageHandler).Methods("GET")
	rtr.HandleFunc("/settings", settingsPageHandler).Methods("GET")

	rtr.HandleFunc("/settings/access_tokens", accessTokensPageHandler).Methods("GET", "POST")
	rtr.HandleFunc("/settings/access_tokens/new", accessTokenNewPageHandler).Methods("GET", "POST")
	rtr.HandleFunc("/settings/access_tokens/new/create",
		accessTokenCreateNewPageHandler).Methods("POST")
	rtr.HandleFunc("/settings/access_tokens/{id:[0-9]+}",
		accessTokenSettingsPageHandler).Methods("GET", "POST")
	rtr.HandleFunc("/settings/access_tokens/{id:[0-9]+}/delete",
		accessTokenDeletePageHandler).Methods("POST")
	rtr.HandleFunc("/settings/access_token_settings/{id:[0-9]+}/update",
		accessTokenUpdatePageHandler).Methods("POST")

	rtr.HandleFunc("/settings/operators", operatorsPageHandler).Methods("GET", "POST")
	rtr.HandleFunc("/settings/operators/new", operatorNewPageHandler).Methods("GET", "POST")
	rtr.HandleFunc("/settings/operators/new/create", operatorCreateNewPageHandler).Methods("POST")
	rtr.HandleFunc("/settings/operators/{id:[0-9]+}",
		operatorSettingsPageHandler).Methods("GET", "POST")
	rtr.HandleFunc("/settings/operators/{id:[0-9]+}/delete",
		operatorDeletePageHandler).Methods("POST")
	rtr.HandleFunc("/settings/operator_settings/{id:[0-9]+}/update",
		operatorUpdatePageHandler).Methods("POST")

	rtr.HandleFunc("/settings/wards", wardsPageHandler).Methods("GET", "POST")
	rtr.HandleFunc("/settings/wards/new", wardNewPageHandler).Methods("GET", "POST")
	rtr.HandleFunc("/settings/wards/new/create",
		wardCreateNewPageHandler).Methods("POST")
	rtr.HandleFunc("/settings/wards/{id:[0-9]+}",
		wardSettingsPageHandler).Methods("GET", "POST")
	rtr.HandleFunc("/settings/wards/{id:[0-9]+}/delete",
		wardDeletePageHandler).Methods("POST")
	rtr.HandleFunc("/settings/ward_settings/{id:[0-9]+}/update",
		wardUpdatePageHandler).Methods("POST")

	pathToResourcesWebview := tools.GetPath("ui/webview/")

	http.Handle("/static/", http.StripPrefix("/static/",
		http.FileServer(http.Dir(pathToResourcesWebview+"./static/"))))
	http.Handle("/", rtr)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err.Error())
	}
}

func homePageHandler(w http.ResponseWriter, _ *http.Request) {
	t := getHtmlTemplates()
	err := t.ExecuteTemplate(w, "index", nil)
	if err != nil {
		panic(err.Error())
	}
}

func observersPageHandler(w http.ResponseWriter, _ *http.Request) {
	t := getHtmlTemplates()
	err := t.ExecuteTemplate(w, "observers", nil)
	if err != nil {
		panic(err.Error())
	}
}

type settings struct {
	AccessTokens []data_manager.AccessToken
	Operators    []data_manager.Operator
	Wards        []data_manager.Ward
}

func settingsPageHandler(w http.ResponseWriter, _ *http.Request) {
	t := getHtmlTemplates()
	err := t.ExecuteTemplate(w, "settings", nil)
	if err != nil {
		panic(err.Error())
	}
}

func accessTokensPageHandler(w http.ResponseWriter, _ *http.Request) {
	var s settings
	s.AccessTokens = data_manager.SelectAccessTokens()

	t := getHtmlTemplates()
	err := t.ExecuteTemplate(w, "access_tokens", s)
	if err != nil {
		panic(err.Error())
	}
}

func accessTokenSettingsPageHandler(w http.ResponseWriter, r *http.Request) {
	accessTokenID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		panic(err.Error())
	}

	var a data_manager.AccessToken
	err = a.SelectFromDBByID(accessTokenID)
	if err != nil {
		panic(err.Error())
	}

	s := strings.Split(a.Value, "")
	switch true {
	case len(a.Value) >= 8:
		begin := s[:3]
		end := s[len(s)-5:]
		a.Value = fmt.Sprintf("%s ****** %s", strings.Join(begin, ""), strings.Join(end, ""))
	case len(a.Value) < 8 && len(a.Value) >= 4:
		end := s[len(s)-4:]
		a.Value = fmt.Sprintf("****** %s", strings.Join(end, ""))
	case len(a.Value) > 0 && len(a.Value) < 4:
		end := s[len(s)-1]
		a.Value = fmt.Sprintf("****** %s", end)
	default:
		a.Value = "******"
	}

	t := getHtmlTemplates()
	err = t.ExecuteTemplate(w, "access_token_settings", a)
	if err != nil {
		panic(err.Error())
	}
}

func accessTokenNewPageHandler(w http.ResponseWriter, _ *http.Request) {
	t := getHtmlTemplates()
	err := t.ExecuteTemplate(w, "access_token_new", nil)
	if err != nil {
		panic(err.Error())
	}
}

func accessTokenCreateNewPageHandler(w http.ResponseWriter, r *http.Request) {
	var a data_manager.AccessToken

	name := r.FormValue("name")
	if len(name) > 0 {
		err := a.SetName(name)
		if err != nil {
			switch true {
			case strings.Contains(strings.ToLower(err.Error()), "string length is zero"):
				fmt.Println("asdasdas") // TODO: обработать ошибку
			case strings.Contains(strings.ToLower(err.Error()), "access token with this name already exists"):
				fmt.Println("asdasdas") // TODO: обработать ошибку
			default:
				tools.WriteToLog(err, debug.Stack())
				panic(err.Error())
			}
		}
	} else {
		return // TODO: обработать ошибку
	}

	value := r.FormValue("value")
	if len(value) > 0 {
		err := a.SetValue(value)
		if err != nil {
			if strings.Contains(strings.ToLower(err.Error()), "string length is zero") {
				fmt.Println("asdasdas") // TODO: обработать ошибку
			} else {
				tools.WriteToLog(err, debug.Stack())
				panic(err.Error())
			}
		}
	} else {
		return // TODO: обработать ошибку
	}

	a.SaveToDB()

	http.Redirect(w, r, "/settings/access_tokens", http.StatusSeeOther)
}

func accessTokenUpdatePageHandler(w http.ResponseWriter, r *http.Request) {
	accessTokenID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		panic(err.Error())
	}

	var a data_manager.AccessToken
	err = a.SelectFromDBByID(accessTokenID)
	if err != nil {
		panic(err.Error())
	}

	var dataIsUpdated bool

	name := r.FormValue("name")
	if len(name) > 0 {
		err = a.SetName(name)
		if err != nil {
			switch true {
			case strings.Contains(strings.ToLower(err.Error()), "string length is zero"):
				fmt.Println("asdasdas") // TODO: обработать ошибку
			case strings.Contains(strings.ToLower(err.Error()), "access token with this name already exists"):
				fmt.Println("asdasdas") // TODO: обработать ошибку
			default:
				tools.WriteToLog(err, debug.Stack())
				panic(err.Error())
			}
		} else {
			dataIsUpdated = true
		}
	}

	value := r.FormValue("value")
	if len(value) > 0 {
		err = a.SetValue(value)
		if err != nil {
			if strings.Contains(strings.ToLower(err.Error()), "string length is zero") {
				fmt.Println("asdasdas") // TODO: обработать ошибку
			} else {
				tools.WriteToLog(err, debug.Stack())
				panic(err.Error())
			}
		} else {
			dataIsUpdated = true
		}
	}

	if dataIsUpdated {
		a.UpdateInDB()
	}

	s := strings.Split(a.Value, "")
	switch true {
	case len(a.Value) >= 8:
		begin := s[:3]
		end := s[len(s)-5:]
		a.Value = fmt.Sprintf("%s ****** %s", strings.Join(begin, ""), strings.Join(end, ""))
	case len(a.Value) < 8 && len(a.Value) >= 4:
		end := s[len(s)-5:]
		a.Value = fmt.Sprintf("****** %s", strings.Join(end, ""))
	case len(a.Value) > 0 && len(a.Value) < 4:
		end := s[len(s)-1]
		a.Value = fmt.Sprintf("****** %s", end)
	default:
		a.Value = "******"
	}

	http.Redirect(w, r, fmt.Sprintf("/settings/access_tokens/%d", a.ID), http.StatusSeeOther)
}

func accessTokenDeletePageHandler(w http.ResponseWriter, r *http.Request) {
	accessTokenID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		panic(err.Error())
	}

	var a data_manager.AccessToken
	err = a.SelectFromDBByID(accessTokenID)
	if err != nil {
		panic(err.Error())
	}

	// TODO: запрос подтверждения на удаления

	a.DeleteFromDB()

	http.Redirect(w, r, "/settings/access_tokens", http.StatusSeeOther)
}

func operatorsPageHandler(w http.ResponseWriter, _ *http.Request) {
	var s settings
	s.Operators = data_manager.SelectOperators()

	t := getHtmlTemplates()
	err := t.ExecuteTemplate(w, "operators", s)
	if err != nil {
		panic(err.Error())
	}
}

func operatorSettingsPageHandler(w http.ResponseWriter, r *http.Request) {
	operatorID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		panic(err.Error())
	}

	var o data_manager.Operator
	err = o.SelectFromDBByID(operatorID)
	if err != nil {
		panic(err.Error())
	}

	t := getHtmlTemplates()
	err = t.ExecuteTemplate(w, "operator_settings", o)
	if err != nil {
		panic(err.Error())
	}
}

func operatorNewPageHandler(w http.ResponseWriter, _ *http.Request) {
	t := getHtmlTemplates()
	err := t.ExecuteTemplate(w, "operator_new", nil)
	if err != nil {
		panic(err.Error())
	}
}

func operatorCreateNewPageHandler(w http.ResponseWriter, r *http.Request) {
	var o data_manager.Operator

	name := r.FormValue("name")
	if len(name) > 0 {
		err := o.SetName(name)
		if err != nil {
			switch true {
			case strings.Contains(strings.ToLower(err.Error()), "string length is zero"):
				fmt.Println("asdasdas") // TODO: обработать ошибку
			case strings.Contains(strings.ToLower(err.Error()), "operator with this name already exists"):
				fmt.Println("asdasdas") // TODO: обработать ошибку
			default:
				tools.WriteToLog(err, debug.Stack())
				panic(err.Error())
			}
		}
	} else {
		return // TODO: обработать ошибку
	}

	value := r.FormValue("vk_id")
	if len(value) > 0 {
		err := o.SetVkID(value)
		if err != nil {
			switch true {
			case strings.Contains(strings.ToLower(err.Error()), "string length is zero"):
				fmt.Println("asdasdas") // TODO: обработать ошибку
			case strings.Contains(strings.ToLower(err.Error()), "vk id starts with zero"):
				fmt.Println("asdasdas") // TODO: обработать ошибку
			case strings.Contains(strings.ToLower(err.Error()), "invalid syntax"):
				fmt.Println("asdasdas") // TODO: обработать ошибку
			default:
				tools.WriteToLog(err, debug.Stack())
				panic(err.Error())
			}
		}
	} else {
		return // TODO: обработать ошибку
	}

	o.SaveToDB()

	http.Redirect(w, r, "/settings/operators", http.StatusSeeOther)
}

func operatorUpdatePageHandler(w http.ResponseWriter, r *http.Request) {
	operatorID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		panic(err.Error())
	}

	var o data_manager.Operator
	err = o.SelectFromDBByID(operatorID)
	if err != nil {
		panic(err.Error())
	}

	var dataIsUpdated bool

	name := r.FormValue("name")
	if len(name) > 0 {
		err = o.SetName(name)
		if err != nil {
			switch true {
			case strings.Contains(strings.ToLower(err.Error()), "string length is zero"):
				fmt.Println("asdasdas") // TODO: обработать ошибку
			case strings.Contains(strings.ToLower(err.Error()), "operator with this name already exists"):
				fmt.Println("asdasdas") // TODO: обработать ошибку
			default:
				tools.WriteToLog(err, debug.Stack())
				panic(err.Error())
			}
		} else {
			dataIsUpdated = true
		}
	}

	value := r.FormValue("vk_id")
	if len(value) > 0 {
		err = o.SetVkID(value)
		if err != nil {
			switch true {
			case strings.Contains(strings.ToLower(err.Error()), "string length is zero"):
				fmt.Println("asdasdas") // TODO: обработать ошибку
			case strings.Contains(strings.ToLower(err.Error()), "vk id starts with zero"):
				fmt.Println("asdasdas") // TODO: обработать ошибку
			case strings.Contains(strings.ToLower(err.Error()), "invalid syntax"):
				fmt.Println("asdasdas") // TODO: обработать ошибку
			default:
				tools.WriteToLog(err, debug.Stack())
				panic(err.Error())
			}
		} else {
			dataIsUpdated = true
		}
	}

	if dataIsUpdated {
		o.UpdateIdDB()
	}

	http.Redirect(w, r, fmt.Sprintf("/settings/operators/%d", o.ID), http.StatusSeeOther)
}

func operatorDeletePageHandler(w http.ResponseWriter, r *http.Request) {
	operatorID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		panic(err.Error())
	}

	var o data_manager.Operator
	err = o.SelectFromDBByID(operatorID)
	if err != nil {
		panic(err.Error())
	}

	// TODO: запрос подтверждения на удаления

	o.DeleteFromDB()

	http.Redirect(w, r, "/settings/operators", http.StatusSeeOther)
}

func wardsPageHandler(w http.ResponseWriter, _ *http.Request) {
	var s settings
	s.Wards = data_manager.SelectWards()

	t := getHtmlTemplates()
	err := t.ExecuteTemplate(w, "wards", s)
	if err != nil {
		panic(err.Error())
	}
}

type wardSettings struct {
	Ward            data_manager.Ward
	AccessTokens    []data_manager.AccessToken
	Operators       []data_manager.Operator
	Observers       []data_manager.Observer
	ObserversTypes  []string
	WallPostTypes   []string
	WallPostTypesRu []string
}

func wardSettingsPageHandler(w http.ResponseWriter, r *http.Request) {
	wardID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		panic(err.Error())
	}

	var ws wardSettings

	err = ws.Ward.SelectFromDBByID(wardID)
	if err != nil {
		panic(err.Error())
	}

	ws.AccessTokens = data_manager.SelectAccessTokens()

	ws.Operators = data_manager.SelectOperators()

	ws.ObserversTypes = []string{"Посты на стене", "Комментарии под постами", "Фото в альбомах",
		"Комментарии под фото", "Видео в альбомах", "Комментарии под видео", "Обсуждения"}

	observersTypes := []string{
		"wall_post", "wall_reply", "photo", "photo_comment", "video", "video_comment", "board_post",
	}
	for _, item := range observersTypes {
		var o data_manager.Observer
		err = o.SelectFromDB(item, ws.Ward.ID)
		if err != nil {
			panic(err.Error())
		}
		ws.Observers = append(ws.Observers, o)
	}

	ws.WallPostTypes = []string{
		"post", "suggest", "postponed",
	}

	ws.WallPostTypesRu = []string{
		"Опубликованные", "Предложенные", "Отложенные",
	}

	t := getHtmlTemplates()
	err = t.ExecuteTemplate(w, "ward_settings", ws)
	if err != nil {
		panic(err.Error())
	}
}

func wardNewPageHandler(w http.ResponseWriter, _ *http.Request) {
	var ws wardSettings

	ws.AccessTokens = data_manager.SelectAccessTokens()

	ws.Operators = data_manager.SelectOperators()

	ws.ObserversTypes = []string{"Посты на стене", "Комментарии под постами", "Фото в альбомах",
		"Комментарии под фото", "Видео в альбомах", "Комментарии под видео", "Обсуждения"}

	observersTypes := []string{
		"wall_post", "wall_reply", "photo", "photo_comment", "video", "video_comment", "board_post",
	}
	for _, item := range observersTypes {
		var o data_manager.Observer
		o.SetName(item)
		ws.Observers = append(ws.Observers, o)
	}

	ws.WallPostTypes = []string{
		"post", "suggest", "postponed",
	}

	ws.WallPostTypesRu = []string{
		"Опубликованные", "Предложенные", "Отложенные",
	}

	t := getHtmlTemplates()
	err := t.ExecuteTemplate(w, "ward_new", ws)
	if err != nil {
		panic(err.Error())
	}
}

func wardCreateNewPageHandler(w http.ResponseWriter, r *http.Request) {
	var ward data_manager.Ward

	name := r.FormValue("name")
	if len(name) > 0 {
		err := ward.SetName(name)
		if err != nil {
			switch true {
			case strings.Contains(strings.ToLower(err.Error()), "string length is zero"):
				return // TODO: обработать ошибку
			case strings.Contains(strings.ToLower(err.Error()), "ward with this name already exists"):
				return // TODO: обработать ошибку
			default:
				tools.WriteToLog(err, debug.Stack())
				panic(err.Error())
			}
		}
	} else {
		return // TODO: обработать ошибку
	}

	vkID := r.FormValue("vk_id")
	if len(vkID) > 0 {
		err := ward.SetVkID(vkID)
		if err != nil {
			switch true {
			case strings.Contains(strings.ToLower(err.Error()), "string length is zero"):
				return // TODO: обработать ошибку
			case strings.Contains(strings.ToLower(err.Error()), "vk id starts with zero"):
				return // TODO: обработать ошибку
			case strings.Contains(strings.ToLower(err.Error()), "invalid syntax"):
				return // TODO: обработать ошибку
			default:
				tools.WriteToLog(err, debug.Stack())
				panic(err.Error())
			}
		}
	} else {
		return // TODO: обработать ошибку
	}

	getAccessToken := r.FormValue("get_access_token")
	if len(getAccessToken) > 0 {
		err := ward.SetAccessToken(getAccessToken)
		if err != nil {
			switch true {
			case strings.Contains(strings.ToLower(err.Error()), "string length is zero"):
				return // TODO: обработать ошибку
			case strings.Contains(strings.ToLower(err.Error()), "no such access token found"):
				return // TODO: обработать ошибку
			default:
				tools.WriteToLog(err, debug.Stack())
				panic(err.Error())
			}
		}
	} else {
		return // TODO: обработать ошибку
	}

	observersTypes := []string{
		"wall_post", "wall_reply", "photo", "photo_comment", "video", "video_comment", "board_post",
	}

	var observers []data_manager.Observer

	for _, item := range observersTypes {
		var observer data_manager.Observer
		observer.SetName(item)

		if item == "wall_post" {
			postType := r.FormValue("post_type")
			if len(postType) > 0 {
				observer.SetAdditionalParams(postType)
			}
		}

		operatorName := r.FormValue(fmt.Sprintf("%s_operator", item))
		if len(operatorName) > 0 {
			err := observer.SetOperator(operatorName)
			if err != nil {
				panic(err.Error()) // TODO: обработать ошибку
			}
		} else {
			return // TODO: обработать ошибку
		}

		sendAccessToken := r.FormValue(fmt.Sprintf("%s_send_access_token", item))
		if len(sendAccessToken) > 0 {
			err := observer.SetAccessToken(sendAccessToken)
			if err != nil {
				panic(err.Error()) // TODO: обработать ошибку
			}
		} else {
			return // TODO: обработать ошибку
		}

		observers = append(observers, observer)
	}

	ward.SaveToDB()

	err := ward.SelectFromDB(ward.Name)
	if err != nil {
		tools.WriteToLog(err, debug.Stack())
		panic(err.Error())
	}

	for _, observer := range observers {
		observer.SetWardID(ward.ID)
		observer.SaveToDB()
	}

	http.Redirect(w, r, "/settings/wards", http.StatusSeeOther)
}

func wardUpdatePageHandler(w http.ResponseWriter, r *http.Request) {
	wardID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		panic(err.Error())
	}

	var ward data_manager.Ward

	err = ward.SelectFromDBByID(wardID)
	if err != nil {
		panic(err.Error())
	}

	var dataIsUpdated bool

	name := r.FormValue("name")
	if len(name) > 0 {
		err = ward.SetName(name)
		if err != nil {
			switch true {
			case strings.Contains(strings.ToLower(err.Error()), "string length is zero"):
				return // TODO: обработать ошибку
			case strings.Contains(strings.ToLower(err.Error()), "ward with this name already exists"):
				return // TODO: обработать ошибку
			default:
				tools.WriteToLog(err, debug.Stack())
				panic(err.Error())
			}
		} else {
			dataIsUpdated = true
		}
	}

	vkID := r.FormValue("vk_id")
	if len(vkID) > 0 {
		err = ward.SetVkID(vkID)
		if err != nil {
			switch true {
			case strings.Contains(strings.ToLower(err.Error()), "string length is zero"):
				return // TODO: обработать ошибку
			case strings.Contains(strings.ToLower(err.Error()), "vk id starts with zero"):
				return // TODO: обработать ошибку
			case strings.Contains(strings.ToLower(err.Error()), "invalid syntax"):
				return // TODO: обработать ошибку
			default:
				tools.WriteToLog(err, debug.Stack())
				panic(err.Error())
			}
		} else {
			dataIsUpdated = true
		}
	}

	getAccessToken := r.FormValue("get_access_token")
	if len(getAccessToken) > 0 {
		err = ward.SetAccessToken(getAccessToken)
		if err != nil {
			switch true {
			case strings.Contains(strings.ToLower(err.Error()), "string length is zero"):
				return // TODO: обработать ошибку
			case strings.Contains(strings.ToLower(err.Error()), "no such access token found"):
				return // TODO: обработать ошибку
			default:
				tools.WriteToLog(err, debug.Stack())
				panic(err.Error())
			}
		} else {
			dataIsUpdated = true
		}
	}

	observersTypes := []string{
		"wall_post", "wall_reply", "photo", "photo_comment", "video", "video_comment", "board_post",
	}

	var observers []data_manager.Observer

	for _, item := range observersTypes {
		var observer data_manager.Observer
		err = observer.SelectFromDB(item, wardID)
		if err != nil {
			panic(err.Error())
		}

		if item == "wall_post" {
			postType := r.FormValue("post_type")
			if len(postType) > 0 {
				observer.SetAdditionalParams(postType)
				dataIsUpdated = true
			}
		}

		operatorName := r.FormValue(fmt.Sprintf("%s_operator", item))
		if len(operatorName) > 0 {
			err = observer.SetOperator(operatorName)
			if err != nil {
				panic(err.Error()) // TODO: обработать ошибку
			} else {
				dataIsUpdated = true
			}
		}

		sendAccessToken := r.FormValue(fmt.Sprintf("%s_send_access_token", item))
		if len(sendAccessToken) > 0 {
			err = observer.SetAccessToken(sendAccessToken)
			if err != nil {
				panic(err.Error()) // TODO: обработать ошибку
			} else {
				dataIsUpdated = true
			}
		}

		observers = append(observers, observer)
	}

	if dataIsUpdated {
		ward.UpdateInDB()
		for _, observer := range observers {
			observer.UpdateInDB()
		}
	}

	http.Redirect(w, r, fmt.Sprintf("/settings/wards/%d", ward.ID), http.StatusSeeOther)
}

func wardDeletePageHandler(w http.ResponseWriter, r *http.Request) {
	wardID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		panic(err.Error())
	}

	var ward data_manager.Ward
	err = ward.SelectFromDBByID(wardID)
	if err != nil {
		panic(err.Error())
	}

	var observers []data_manager.Observer

	observersTypes := []string{
		"wall_post", "wall_reply", "photo", "photo_comment", "video", "video_comment", "board_post",
	}

	for _, item := range observersTypes {
		var o data_manager.Observer
		o.SelectFromDB(item, ward.ID)

		observers = append(observers, o)
	}

	// TODO: запрос подтверждения на удаления

	for _, observer := range observers {
		observer.DeleteFromDB()
	}
	ward.DeleteFromDB()

	http.Redirect(w, r, "/settings/wards", http.StatusSeeOther)
}

func getHtmlTemplates() *template.Template {
	pathToResourcesWebview := tools.GetPath("ui/webview/")

	t, err := template.ParseFiles(
		pathToResourcesWebview+"html/header.html",
		pathToResourcesWebview+"html/general/index.html",
		pathToResourcesWebview+"html/observers/observers.html",
		pathToResourcesWebview+"html/settings/settings.html",
		pathToResourcesWebview+"html/settings/access_tokens.html",
		pathToResourcesWebview+"html/settings/access_token_new.html",
		pathToResourcesWebview+"html/settings/access_token_settings.html",
		pathToResourcesWebview+"html/settings/operators.html",
		pathToResourcesWebview+"html/settings/operator_new.html",
		pathToResourcesWebview+"html/settings/operator_settings.html",
		pathToResourcesWebview+"html/settings/wards.html",
		pathToResourcesWebview+"html/settings/ward_new.html",
		pathToResourcesWebview+"html/settings/ward_settings.html",
		pathToResourcesWebview+"html/footer.html")
	if err != nil {
		panic(err.Error())
	}
	return t
}
