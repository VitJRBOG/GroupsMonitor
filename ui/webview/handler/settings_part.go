package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/VitJRBOG/GroupsObserver/data_manager"
	"github.com/VitJRBOG/GroupsObserver/tools"
	"github.com/gorilla/mux"
)

func settingsPageHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/settings/access_tokens", http.StatusSeeOther)
}

type accessTokenPageData struct {
	AccessTokenID int                        `json:"access_token_id"`
	AccessToken   data_manager.AccessToken   `json:"access_token"`
	AccessTokens  []data_manager.AccessToken `json:"access_tokens"`
	Error         string                     `json:"error"`
}

func (a *accessTokenPageData) hideAccessTokenValues() {
	for i := 0; i < len(a.AccessTokens); i++ {
		a.AccessTokens[i].Value = a.hideAccessTokenValue(a.AccessTokens[i].Value)
	}
}

func (a *accessTokenPageData) hideAccessTokenValue(value string) string {
	v := strings.Split(value, "")
	switch true {
	case len(v) >= 8:
		begin := v[:3]
		end := v[len(v)-5:]
		value = fmt.Sprintf("%s ****** %s", strings.Join(begin, ""), strings.Join(end, ""))
	case len(v) < 8 && len(v) >= 4:
		end := v[len(v)-4:]
		value = fmt.Sprintf("****** %s", strings.Join(end, ""))
	case len(v) > 0 && len(v) < 4:
		end := v[len(v)-1]
		value = fmt.Sprintf("****** %s", end)
	default:
		value = "******"
	}
	return value
}

func accessTokensPageHandler(w http.ResponseWriter, _ *http.Request) {
	var data accessTokenPageData
	data.AccessTokens = data_manager.SelectAccessTokens()

	data.hideAccessTokenValues()

	t := getHtmlTemplates()
	err := t.ExecuteTemplate(w, "access_tokens", data)
	if err != nil {
		tools.WriteToLog(err, debug.Stack())
		panic(err.Error())
	}
}

func accessTokensGet(w http.ResponseWriter, _ *http.Request) {
	var data accessTokenPageData
	data.AccessTokens = data_manager.SelectAccessTokens()

	data.hideAccessTokenValues()

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

func accessTokenSettingsPageHandler(w http.ResponseWriter, r *http.Request) {
	accessTokenID := mux.Vars(r)["id"]

	t := getHtmlTemplates()
	err := t.ExecuteTemplate(w, "access_token_settings", accessTokenID)
	if err != nil {
		tools.WriteToLog(err, debug.Stack())
		panic(err.Error())
	}
}

func accessTokenSettingsGet(w http.ResponseWriter, r *http.Request) {
	var data accessTokenPageData
	var err error
	data.AccessTokenID, err = strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		tools.WriteToLog(err, debug.Stack())
		panic(err.Error())
	}

	err = data.AccessToken.SelectFromDBByID(data.AccessTokenID)
	if err != nil {
		tools.WriteToLog(err, debug.Stack())
		panic(err.Error())
	}

	data.AccessToken.Value = data.hideAccessTokenValue(data.AccessToken.Value)

	data.AccessTokens = data_manager.SelectAccessTokens()

	data.hideAccessTokenValues()

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

func accessTokenNewPageHandler(w http.ResponseWriter, _ *http.Request) {
	t := getHtmlTemplates()
	err := t.ExecuteTemplate(w, "access_token_new", nil)
	if err != nil {
		tools.WriteToLog(err, debug.Stack())
		panic(err.Error())
	}
}

func accessTokenCreateNewPageHandler(w http.ResponseWriter, r *http.Request) {
	var data accessTokenPageData

	name := r.FormValue("name")
	err := data.AccessToken.SetName(name)
	if err != nil {
		switch true {
		case strings.Contains(strings.ToLower(err.Error()), "string length is zero"):
			data.Error = "Поле «Название» не должно быть пустым."
		case strings.Contains(strings.ToLower(err.Error()), "access token with this name already exists"):
			data.Error = "Ключ с таким названием уже существует."
		default:
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
	}

	value := r.FormValue("value")
	err = data.AccessToken.SetValue(value)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "string length is zero") {
			data.Error = "Поле «Значение» не должно быть пустым."
		} else {
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
	}

	if data.Error == "" {
		data.AccessToken.SaveToDB()
	}

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

func accessTokenUpdatePageHandler(w http.ResponseWriter, r *http.Request) {
	var data accessTokenPageData
	var err error

	data.AccessTokenID, err = strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		tools.WriteToLog(err, debug.Stack())
		panic(err.Error())
	}

	err = data.AccessToken.SelectFromDBByID(data.AccessTokenID)
	if err != nil {
		tools.WriteToLog(err, debug.Stack())
		panic(err.Error())
	}

	name := r.FormValue("name")
	err = data.AccessToken.SetName(name)
	if err != nil {
		switch true {
		case strings.Contains(strings.ToLower(err.Error()), "string length is zero"):
			data.Error = "Поле «Название» не должно быть пустым."
		case strings.Contains(strings.ToLower(err.Error()), "access token with this name already exists"):
			data.Error = "Ключ с таким названием уже существует."
		default:
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
	}

	value := r.FormValue("value")
	err = data.AccessToken.SetValue(value)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "string length is zero") {
			data.Error = "Поле «Значение» не должно быть пустым."
		} else {
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
	}

	if data.Error == "" {
		data.AccessToken.UpdateInDB()
	}

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

func accessTokenDeletePageHandler(w http.ResponseWriter, r *http.Request) {
	var data accessTokenPageData
	var err error

	data.AccessTokenID, err = strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		tools.WriteToLog(err, debug.Stack())
		panic(err.Error())
	}

	err = data.AccessToken.SelectFromDBByID(data.AccessTokenID)
	if err != nil {
		tools.WriteToLog(err, debug.Stack())
		panic(err.Error())
	}

	data.AccessToken.DeleteFromDB()

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

type operatorPagesData struct {
	OperatorID int                     `json:"operator_id"`
	Operator   data_manager.Operator   `json:"operator"`
	Operators  []data_manager.Operator `json:"operators"`
	Error      string                  `json:"error"`
}

func operatorsPageHandler(w http.ResponseWriter, _ *http.Request) {
	var data operatorPagesData
	data.Operators = data_manager.SelectOperators()

	t := getHtmlTemplates()
	err := t.ExecuteTemplate(w, "operators", data)
	if err != nil {
		tools.WriteToLog(err, debug.Stack())
		panic(err.Error())
	}
}

func operatorsGet(w http.ResponseWriter, r *http.Request) {
	var data operatorPagesData
	data.Operators = data_manager.SelectOperators()

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

func operatorSettingsPageHandler(w http.ResponseWriter, r *http.Request) {
	operatorID := mux.Vars(r)["id"]

	t := getHtmlTemplates()
	err := t.ExecuteTemplate(w, "operator_settings", operatorID)
	if err != nil {
		tools.WriteToLog(err, debug.Stack())
		panic(err.Error())
	}
}

func operatorSettingsGet(w http.ResponseWriter, r *http.Request) {
	var data operatorPagesData
	var err error
	data.OperatorID, err = strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		tools.WriteToLog(err, debug.Stack())
		panic(err.Error())
	}

	err = data.Operator.SelectFromDBByID(data.OperatorID)
	if err != nil {
		tools.WriteToLog(err, debug.Stack())
		panic(err.Error())
	}

	data.Operators = data_manager.SelectOperators()

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

func operatorNewPageHandler(w http.ResponseWriter, _ *http.Request) {
	t := getHtmlTemplates()
	err := t.ExecuteTemplate(w, "operator_new", nil)
	if err != nil {
		tools.WriteToLog(err, debug.Stack())
		panic(err.Error())
	}
}

func operatorCreateNewPageHandler(w http.ResponseWriter, r *http.Request) {
	var data operatorPagesData

	name := r.FormValue("name")
	err := data.Operator.SetName(name)
	if err != nil {
		switch true {
		case strings.Contains(strings.ToLower(err.Error()), "string length is zero"):
			data.Error = "Поле «Название» не должно быть пустым."
		case strings.Contains(strings.ToLower(err.Error()), "operator with this name already exists"):
			data.Error = "Оператор с таким названием уже существует."
		default:
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
	}

	value := r.FormValue("vk_id")
	err = data.Operator.SetVkID(value)
	if err != nil {
		switch true {
		case strings.Contains(strings.ToLower(err.Error()), "string length is zero"):
			data.Error = "Поле «ID в ВК» не должно быть пустым."
		case strings.Contains(strings.ToLower(err.Error()), "vk id starts with zero"):
			data.Error = "Идентификатор ВК не должен начинаться с нуля."
		case strings.Contains(strings.ToLower(err.Error()), "invalid syntax"):
			data.Error = "Идентификатор ВК должен быть числом."
		default:
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
	}

	if data.Error == "" {
		data.Operator.SaveToDB()
	}

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

func operatorUpdatePageHandler(w http.ResponseWriter, r *http.Request) {
	var data operatorPagesData
	var err error

	data.OperatorID, err = strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		tools.WriteToLog(err, debug.Stack())
		panic(err.Error())
	}

	err = data.Operator.SelectFromDBByID(data.OperatorID)
	if err != nil {
		tools.WriteToLog(err, debug.Stack())
		panic(err.Error())
	}

	name := r.FormValue("name")
	err = data.Operator.SetName(name)
	if err != nil {
		switch true {
		case strings.Contains(strings.ToLower(err.Error()), "string length is zero"):
			data.Error = "Поле «Название» не должно быть пустым."
		case strings.Contains(strings.ToLower(err.Error()), "operator with this name already exists"):
			data.Error = "Оператор с таким названием уже существует."
		default:
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
	}

	value := r.FormValue("vk_id")
	err = data.Operator.SetVkID(value)
	if err != nil {
		switch true {
		case strings.Contains(strings.ToLower(err.Error()), "string length is zero"):
			data.Error = "Поле «ID в ВК» не должно быть пустым."
		case strings.Contains(strings.ToLower(err.Error()), "vk id starts with zero"):
			data.Error = "Идентификатор ВК не должен начинаться с нуля."
		case strings.Contains(strings.ToLower(err.Error()), "invalid syntax"):
			data.Error = "Идентификатор ВК должен быть числом."
		default:
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
	}

	if data.Error == "" {
		data.Operator.UpdateIdDB()
	}

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

func operatorDeletePageHandler(w http.ResponseWriter, r *http.Request) {
	var data operatorPagesData
	var err error

	data.OperatorID, err = strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		tools.WriteToLog(err, debug.Stack())
		panic(err.Error())
	}

	err = data.Operator.SelectFromDBByID(data.OperatorID)
	if err != nil {
		tools.WriteToLog(err, debug.Stack())
		panic(err.Error())
	}

	data.Operator.DeleteFromDB()

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

type wardPageData struct {
	WardID          int                        `json:"ward_id"`
	Ward            data_manager.Ward          `json:"ward"`
	Wards           []data_manager.Ward        `json:"wards"`
	AccessTokens    []data_manager.AccessToken `json:"access_tokens"`
	Operators       []data_manager.Operator    `json:"operators"`
	Observers       []data_manager.Observer    `json:"observers"`
	ObserverTypes   []string                   `json:"observer_types"`
	WallPostTypes   []string                   `json:"wall_post_types"`
	WallPostTypesRu []string                   `json:"wall_post_types_ru"`
	Error           string                     `json:"error"`
}

func wardsPageHandler(w http.ResponseWriter, _ *http.Request) {
	var data wardPageData
	data.Wards = data_manager.SelectWards()

	t := getHtmlTemplates()
	err := t.ExecuteTemplate(w, "wards", data)
	if err != nil {
		tools.WriteToLog(err, debug.Stack())
		panic(err.Error())
	}
}

func wardsGet(w http.ResponseWriter, r *http.Request) {
	var data wardPageData
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

func wardSettingsPageHandler(w http.ResponseWriter, r *http.Request) {
	wardID := mux.Vars(r)["id"]

	t := getHtmlTemplates()
	err := t.ExecuteTemplate(w, "ward_settings", wardID)
	if err != nil {
		tools.WriteToLog(err, debug.Stack())
		panic(err.Error())
	}
}

func wardSettingsGet(w http.ResponseWriter, r *http.Request) {
	var data wardPageData
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

	data.AccessTokens = data_manager.SelectAccessTokens()

	data.Operators = data_manager.SelectOperators()

	data.ObserverTypes = []string{"Посты на стене", "Комментарии на стене", "Фото в альбомах",
		"Комментарии под фото", "Видео в альбомах", "Комментарии под видео", "Обсуждения"}

	observerTypes := []string{
		"wall_post", "wall_reply", "photo", "photo_comment", "video", "video_comment", "board_post",
	}
	for _, item := range observerTypes {
		var o data_manager.Observer
		err = o.SelectFromDB(item, data.Ward.ID)
		if err != nil {
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
		data.Observers = append(data.Observers, o)
	}

	data.WallPostTypes = []string{
		"post", "suggest", "postponed",
	}

	data.WallPostTypesRu = []string{
		"Опубликованные", "Предложенные", "Отложенные",
	}

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

func wardNewPageHandler(w http.ResponseWriter, _ *http.Request) {
	t := getHtmlTemplates()
	err := t.ExecuteTemplate(w, "ward_new", nil)
	if err != nil {
		panic(err.Error())
	}
}

func wardGetNew(w http.ResponseWriter, _ *http.Request) {
	var data wardPageData

	data.AccessTokens = data_manager.SelectAccessTokens()

	data.Operators = data_manager.SelectOperators()

	data.ObserverTypes = []string{"Посты на стене", "Комментарии на стене", "Фото в альбомах",
		"Комментарии под фото", "Видео в альбомах", "Комментарии под видео", "Обсуждения"}

	observerTypes := []string{
		"wall_post", "wall_reply", "photo", "photo_comment", "video", "video_comment", "board_post",
	}
	for _, item := range observerTypes {
		var o data_manager.Observer
		o.SetName(item)
		data.Observers = append(data.Observers, o)
	}

	data.WallPostTypes = []string{
		"post", "suggest", "postponed",
	}

	data.WallPostTypesRu = []string{
		"Опубликованные", "Предложенные", "Отложенные",
	}

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

func wardCreateNewPageHandler(w http.ResponseWriter, r *http.Request) {
	var data wardPageData

	name := r.FormValue("name")
	err := data.Ward.SetName(name)
	if err != nil {
		switch true {
		case strings.Contains(strings.ToLower(err.Error()), "string length is zero"):
			data.Error = "Поле «Название» не должно быть пустым."
		case strings.Contains(strings.ToLower(err.Error()), "ward with this name already exists"):
			data.Error = "Подопечный с таким названием уже существует."
		default:
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
	}

	vkID := r.FormValue("vk_id")
	err = data.Ward.SetVkID(vkID)
	if err != nil {
		switch true {
		case strings.Contains(strings.ToLower(err.Error()), "string length is zero"):
			data.Error = "Поле «ID в ВК» не должно быть пустым."
		case strings.Contains(strings.ToLower(err.Error()), "vk id starts with zero"):
			data.Error = "Идентификатор ВК не должен начинаться с нуля."
		case strings.Contains(strings.ToLower(err.Error()), "invalid syntax"):
			data.Error = "Идентификатор ВК должен быть числом."
		default:
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
	}

	getAccessToken := r.FormValue("get_access_token")
	err = data.Ward.SetAccessToken(getAccessToken)
	if err != nil {
		switch true {
		case strings.Contains(strings.ToLower(err.Error()), "string length is zero"):
			data.Error = "Значение «Get-ключ доступа» не должно быть пустым."
		case strings.Contains(strings.ToLower(err.Error()), "no such access token found"):
			data.Error = "Указанный ключ доступ не найден."
		default:
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
	}

	observerTypes := []string{
		"wall_post", "wall_reply", "photo", "photo_comment", "video", "video_comment", "board_post",
	}

	for _, item := range observerTypes {
		var o data_manager.Observer
		o.SetName(item)

		if item == "wall_post" {
			postType := r.FormValue("post_type")
			if len(postType) > 0 {
				o.SetAdditionalParams(postType)
			}
		}

		operatorName := r.FormValue(fmt.Sprintf("%s_operator", item))
		err = o.SetOperator(operatorName)
		if err != nil {
			switch true {
			case strings.Contains(strings.ToLower(err.Error()), "string length is zero"):
				data.Error = "Значение «Оператор» не должно быть пустым."
			case strings.Contains(strings.ToLower(err.Error()), "no such operator found"):
				data.Error = "Оператор с таким названием не найден."
			default:
				tools.WriteToLog(err, debug.Stack())
				panic(err.Error())
			}
		}

		sendAccessToken := r.FormValue(fmt.Sprintf("%s_send_access_token", item))
		err = o.SetAccessToken(sendAccessToken)
		if err != nil {
			switch true {
			case strings.Contains(strings.ToLower(err.Error()), "string length is zero"):
				data.Error = "Значение «Get-ключ доступа» не должно быть пустым."
			case strings.Contains(strings.ToLower(err.Error()), "no such access token found"):
				data.Error = "Указанный ключ доступ не найден."
			default:
				tools.WriteToLog(err, debug.Stack())
				panic(err.Error())
			}
		}

		data.Observers = append(data.Observers, o)
	}

	if data.Error == "" {
		data.Ward.SaveToDB()

		err := data.Ward.SelectFromDB(data.Ward.Name)
		if err != nil {
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}

		for _, observer := range data.Observers {
			observer.SetWardID(data.Ward.ID)
			observer.SaveToDB()
		}
	}

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

func wardUpdatePageHandler(w http.ResponseWriter, r *http.Request) {
	var data wardPageData
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

	name := r.FormValue("name")
	err = data.Ward.SetName(name)
	if err != nil {
		switch true {
		case strings.Contains(strings.ToLower(err.Error()), "string length is zero"):
			data.Error = "Поле «Название» не должно быть пустым."
		case strings.Contains(strings.ToLower(err.Error()), "ward with this name already exists"):
			data.Error = "Подопечный с таким названием уже существует."
		default:
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
	}

	vkID := r.FormValue("vk_id")
	err = data.Ward.SetVkID(vkID)
	if err != nil {
		switch true {
		case strings.Contains(strings.ToLower(err.Error()), "string length is zero"):
			data.Error = "Поле «ID в ВК» не должно быть пустым."
		case strings.Contains(strings.ToLower(err.Error()), "vk id starts with zero"):
			data.Error = "Идентификатор ВК не должен начинаться с нуля."
		case strings.Contains(strings.ToLower(err.Error()), "invalid syntax"):
			data.Error = "Идентификатор ВК должен быть числом."
		default:
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
	}

	getAccessToken := r.FormValue("get_access_token")
	err = data.Ward.SetAccessToken(getAccessToken)
	if err != nil {
		switch true {
		case strings.Contains(strings.ToLower(err.Error()), "string length is zero"):
			data.Error = "Значение «Get-ключ доступа» не должно быть пустым."
		case strings.Contains(strings.ToLower(err.Error()), "no such access token found"):
			data.Error = "Указанный ключ доступ не найден."
		default:
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
	}

	observerTypes := []string{
		"wall_post", "wall_reply", "photo", "photo_comment", "video", "video_comment", "board_post",
	}

	for _, item := range observerTypes {
		var observer data_manager.Observer
		err = observer.SelectFromDB(item, data.WardID)
		if err != nil {
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}

		if item == "wall_post" {
			postType := r.FormValue("post_type")
			observer.SetAdditionalParams(postType)
		}

		operatorName := r.FormValue(fmt.Sprintf("%s_operator", item))
		err = observer.SetOperator(operatorName)
		if err != nil {
			switch true {
			case strings.Contains(strings.ToLower(err.Error()), "string length is zero"):
				data.Error = "Значение «Оператор» не должно быть пустым."
			case strings.Contains(strings.ToLower(err.Error()), "no such operator found"):
				data.Error = "Оператор с таким названием не найден."
			default:
				tools.WriteToLog(err, debug.Stack())
				panic(err.Error())
			}
		}

		sendAccessToken := r.FormValue(fmt.Sprintf("%s_send_access_token", item))
		err = observer.SetAccessToken(sendAccessToken)
		if err != nil {
			switch true {
			case strings.Contains(strings.ToLower(err.Error()), "string length is zero"):
				data.Error = "Значение «Get-ключ доступа» не должно быть пустым."
			case strings.Contains(strings.ToLower(err.Error()), "no such access token found"):
				data.Error = "Указанный ключ доступ не найден."
			default:
				tools.WriteToLog(err, debug.Stack())
				panic(err.Error())
			}
		}

		data.Observers = append(data.Observers, observer)
	}

	if data.Error == "" {
		data.Ward.UpdateInDB()
		for _, observer := range data.Observers {
			observer.UpdateInDB()
		}
	}

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

func wardDeletePageHandler(w http.ResponseWriter, r *http.Request) {
	var data wardPageData
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

	observerTypes := []string{
		"wall_post", "wall_reply", "photo", "photo_comment", "video", "video_comment", "board_post",
	}

	for _, item := range observerTypes {
		var o data_manager.Observer
		o.SelectFromDB(item, data.Ward.ID)

		data.Observers = append(data.Observers, o)
	}

	for _, observer := range data.Observers {
		observer.DeleteFromDB()
	}
	data.Ward.DeleteFromDB()

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
