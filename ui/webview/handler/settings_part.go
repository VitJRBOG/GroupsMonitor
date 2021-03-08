package handler

import (
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
	AccessTokenID int
	AccessToken   data_manager.AccessToken
	AccessTokens  []data_manager.AccessToken
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

func accessTokenSettingsPageHandler(w http.ResponseWriter, r *http.Request) {
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

	t := getHtmlTemplates()
	err = t.ExecuteTemplate(w, "access_token_settings", data)
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
	if len(name) > 0 {
		err := data.AccessToken.SetName(name)
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
		err := data.AccessToken.SetValue(value)
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

	data.AccessToken.SaveToDB()

	http.Redirect(w, r, "/settings/access_tokens", http.StatusSeeOther)
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

	var dataIsUpdated bool

	name := r.FormValue("name")
	if len(name) > 0 {
		err = data.AccessToken.SetName(name)
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
		err = data.AccessToken.SetValue(value)
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
		data.AccessToken.UpdateInDB()
	}

	http.Redirect(w, r, fmt.Sprintf("/settings/access_tokens/%d",
		data.AccessTokenID), http.StatusSeeOther)
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

	// TODO: запрос подтверждения на удаления

	data.AccessToken.DeleteFromDB()

	http.Redirect(w, r, "/settings/access_tokens", http.StatusSeeOther)
}

type operatorPagesData struct {
	OperatorID int
	Operator   data_manager.Operator
	Operators  []data_manager.Operator
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

func operatorSettingsPageHandler(w http.ResponseWriter, r *http.Request) {
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

	t := getHtmlTemplates()
	err = t.ExecuteTemplate(w, "operator_settings", data)
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
	if len(name) > 0 {
		err := data.Operator.SetName(name)
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
		err := data.Operator.SetVkID(value)
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

	data.Operator.SaveToDB()

	http.Redirect(w, r, "/settings/operators", http.StatusSeeOther)
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

	var dataIsUpdated bool

	name := r.FormValue("name")
	if len(name) > 0 {
		err = data.Operator.SetName(name)
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
		err = data.Operator.SetVkID(value)
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
		data.Operator.UpdateIdDB()
	}

	http.Redirect(w, r, fmt.Sprintf("/settings/operators/%d",
		data.OperatorID), http.StatusSeeOther)
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

	// TODO: запрос подтверждения на удаления

	data.Operator.DeleteFromDB()

	http.Redirect(w, r, "/settings/operators", http.StatusSeeOther)
}

type wardPageData struct {
	WardID          int
	Ward            data_manager.Ward
	Wards           []data_manager.Ward
	AccessTokens    []data_manager.AccessToken
	Operators       []data_manager.Operator
	Observers       []data_manager.Observer
	ObserverTypes   []string
	WallPostTypes   []string
	WallPostTypesRu []string
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

func wardSettingsPageHandler(w http.ResponseWriter, r *http.Request) {
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

	t := getHtmlTemplates()
	err = t.ExecuteTemplate(w, "ward_settings", data)
	if err != nil {
		tools.WriteToLog(err, debug.Stack())
		panic(err.Error())
	}
}

func wardNewPageHandler(w http.ResponseWriter, _ *http.Request) {
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

	t := getHtmlTemplates()
	err := t.ExecuteTemplate(w, "ward_new", data)
	if err != nil {
		panic(err.Error())
	}
}

func wardCreateNewPageHandler(w http.ResponseWriter, r *http.Request) {
	var data wardPageData

	name := r.FormValue("name")
	if len(name) > 0 {
		err := data.Ward.SetName(name)
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
		err := data.Ward.SetVkID(vkID)
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
		err := data.Ward.SetAccessToken(getAccessToken)
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
		if len(operatorName) > 0 {
			err := o.SetOperator(operatorName)
			if err != nil {
				panic(err.Error()) // TODO: обработать ошибку
			}
		} else {
			return // TODO: обработать ошибку
		}

		sendAccessToken := r.FormValue(fmt.Sprintf("%s_send_access_token", item))
		if len(sendAccessToken) > 0 {
			err := o.SetAccessToken(sendAccessToken)
			if err != nil {
				panic(err.Error()) // TODO: обработать ошибку
			}
		} else {
			return // TODO: обработать ошибку
		}

		data.Observers = append(data.Observers, o)
	}

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

	http.Redirect(w, r, "/settings/wards", http.StatusSeeOther)
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

	var dataIsUpdated bool

	name := r.FormValue("name")
	if len(name) > 0 {
		err = data.Ward.SetName(name)
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
		err = data.Ward.SetVkID(vkID)
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
		err = data.Ward.SetAccessToken(getAccessToken)
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
			if len(postType) > 0 {
				observer.SetAdditionalParams(postType)
				dataIsUpdated = true
			}
		}

		operatorName := r.FormValue(fmt.Sprintf("%s_operator", item))
		if len(operatorName) > 0 {
			err = observer.SetOperator(operatorName)
			if err != nil {
				tools.WriteToLog(err, debug.Stack())
				panic(err.Error()) // TODO: обработать ошибку
			} else {
				dataIsUpdated = true
			}
		}

		sendAccessToken := r.FormValue(fmt.Sprintf("%s_send_access_token", item))
		if len(sendAccessToken) > 0 {
			err = observer.SetAccessToken(sendAccessToken)
			if err != nil {
				tools.WriteToLog(err, debug.Stack())
				panic(err.Error()) // TODO: обработать ошибку
			} else {
				dataIsUpdated = true
			}
		}

		data.Observers = append(data.Observers, observer)
	}

	if dataIsUpdated {
		data.Ward.UpdateInDB()
		for _, observer := range data.Observers {
			observer.UpdateInDB()
		}
	}

	http.Redirect(w, r, fmt.Sprintf("/settings/wards/%d", data.Ward.ID), http.StatusSeeOther)
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

	// TODO: запрос подтверждения на удаления

	for _, observer := range data.Observers {
		observer.DeleteFromDB()
	}
	data.Ward.DeleteFromDB()

	http.Redirect(w, r, "/settings/wards", http.StatusSeeOther)
}
