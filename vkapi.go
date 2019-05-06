package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"runtime"
	"strings"
	"time"
)

// SendVKAPIQuery отправляет post-запрос к VK API
// возвращает ответ сервера и ошибку
func SendVKAPIQuery(sender string, methodName string,
	valuesBytes []byte, subject Subject) (map[string]interface{}, error) {

	// получаем имя текущего модуля мониторинга для запроса в БД
	monitorName := strings.Replace(sender, subject.Name+"'s ", "", 1)
	monitorName = strings.Replace(monitorName, "monitoring", "monitor", 1)
	monitorName = strings.Replace(monitorName, " ", "_", -1)

	// получаем данные о текущем модуле мониторинга
	monitor, err := SelectDBMonitor(monitorName, subject.ID)
	if err != nil {
		return nil, err
	}

	// получаем токен доступа для данного метода vk api, субъекта и модуля мониторинга
	accessToken, err := GetAccessToken(methodName, subject, monitor.ID)
	if err != nil {
		return nil, err
	}

	// формируем url для запроса к vk api
	query := "https://api.vk.com/method/"
	query += methodName
	query += "?access_token=" + accessToken.Value
	query += "&lang=0" // чтобы vk api отвечал на запросы по-русски

	// преобразуем массив байт в словарь
	var f interface{}
	err = json.Unmarshal(valuesBytes, &f)
	if err != nil {
		return nil, err
	}
	values := f.(map[string]interface{})

	// извлекаем из словаря ключи
	keys := make([]string, 0, len(values))
	for k := range values {
		keys = append(keys, k)
	}

	// собираем url запроса из ключей и значений
	for _, key := range keys {
		query += "&" + key + "=" + url.QueryEscape(values[key].(string))
	}

	// отправляем запрос
	requestResult, err := http.Post(query, "", nil)
	if err != nil {

		// если словил ошибку, то преобразуем ее в текст
		errorMessage := fmt.Sprintf("%v", err)

		// и отправляем в обработчик ошибок
		typeError, causeError := RequestError(errorMessage)

		// потом проверяем тип полученной ошибки
		switch typeError {

		// задержка, если получена ошибка, которую можно решить таким образом
		case "timeout error":
			interval := 60 * time.Second
			message := fmt.Sprintf("ERROR: %v. Timeout for %v...", causeError, interval)
			OutputMessage(sender, message)
			time.Sleep(interval)
			return SendVKAPIQuery(sender, methodName, valuesBytes, subject)
		}

		// если пойманная ошибка не обрабатывается, то возвращаем ее стандартным путем
		return nil, err
	}

	// извлекаем тело ответа и преобразуем его полученный ответ сервера в массив байт
	body, err := ioutil.ReadAll(requestResult.Body)
	defer requestResult.Body.Close()
	if err != nil {
		return nil, err
	}

	// преобразуем тело ответа в карту
	err = json.Unmarshal(body, &f)
	if err != nil {
		return nil, err
	}
	response := f.(map[string]interface{})

	if _, exist := response["response"]; exist == true {
		return response, nil
	}

	// проверяем ошибки, которые вернул сервер vk api
	if _, exist := response["error"]; exist == true {
		errorItem := response["error"]
		errorMessage, _ := errorItem.(map[string]interface{})["error_msg"].(string)
		typeError, causeError := VkAPIError(errorMessage)

		// смотрим тип ошибки
		switch typeError {

		// задержка, если слишком часто отправляются запросы
		case "timeout error":
			if causeError != "many requests per second" {

				// однотипно обработаем все ошибки, которые не "many requests per second"
				interval := 1 * time.Second
				message := fmt.Sprintf("ERROR: %v. Timeout for %v...", causeError, interval)
				OutputMessage(sender, message)
				time.Sleep(interval)
			} else {

				// индивидуально обработаем ошибку "many requests per second"
				// потому что они выскакивают регулярно
				interval := 500 * time.Millisecond
				time.Sleep(interval)
			}
			return SendVKAPIQuery(sender, methodName, valuesBytes, subject)

		// задержка, если токен доступа невалидный
		case "access token error":
			interval := 60 * time.Second
			message := fmt.Sprintf("ERROR: %v. Need new %v's access token. Waiting for %v...",
				causeError, accessToken.Name, interval)
			OutputMessage(sender, message)
			time.Sleep(interval)
			return SendVKAPIQuery(sender, methodName, valuesBytes, subject)

		// выход из обработчика ошибок, если эта ошибка обрабатывается в другом месте
		case "skip error":
			return response, nil

		// выход, если неизвестная ошибка
		default:
			message := fmt.Sprintf("UNKNOWN ERROR: %v. Exit...",
				errorItem.(map[string]interface{})["error_msg"].(string))
			OutputMessage(sender, message)
			runtime.Goexit()
		}

	}

	return response, nil
}

// SendMessage посылает запрос к vk api на отправку сообщения
func SendMessage(sender string, jsonDump string, subject Subject) error {
	// преобразуем строку с данными для карты в массив байт из этой карты
	values, err := MakeJSON(jsonDump)
	if err != nil {
		return err
	}

	// отправляем к vk api запрос на отправку сообщения с этими данными
	_, err = SendVKAPIQuery(sender, "messages.send", values, subject)
	if err != nil {
		return err
	}

	return nil
}

// Attachment - структура для данных о прикрепленном контенте
type Attachment struct {
	Type      string `json:"type"`
	OwnerID   int    `json:"owner_id"`
	ID        int    `json:"id"`
	AccessKey string `json:"access_key"`
	URL       string `json:"url"`
}

// ParseAttachments извлекает данные о прикреплениях из карты vk api
func ParseAttachments(mediaContentMaps []interface{}) []Attachment {
	var attachments []Attachment

	// перебираем элементы с данными о прикреплениях
	for _, mediaItemMap := range mediaContentMaps {
		mediaItem := mediaItemMap.(map[string]interface{})

		var attachment Attachment

		// получаем тип прикрепления
		typeMediaItem := mediaItem["type"].(string)

		// проверяем тип прикрепления на соответствие обрабатываемым типам
		match := false
		switch typeMediaItem {
		case "photo":
			match = true
		case "video":
			match = true
		case "audio":
			match = true
		case "doc":
			match = true
		case "poll":
			match = true
		case "link":
			match = true
		}

		// если соответствует, то парсим данные в структуру
		if match {
			attachment.Type = typeMediaItem
			if typeMediaItem == "link" {
				attachment.URL = mediaItem["link"].(map[string]interface{})["url"].(string)
			} else {
				attachment.OwnerID = int(mediaItem[typeMediaItem].(map[string]interface{})["owner_id"].(float64))
				attachment.ID = int(mediaItem[typeMediaItem].(map[string]interface{})["id"].(float64))
				if accessKey, exist := mediaItem[typeMediaItem].(map[string]interface{})["access_key"]; exist {
					attachment.AccessKey = accessKey.(string)
				}
			}
			attachments = append(attachments, attachment)
		}
	}
	return attachments
}

// VKCommunity - структура данных о сообществе VK
type VKCommunity struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	ScreenName string `json:"screen_name"`
}

// GetCommunityInfo получает название сообщества по его id
func GetCommunityInfo(sender string, subject Subject, groupID int) (VKCommunity, error) {
	var vkCommunity VKCommunity

	// формируем json запроса к vk api
	jsonDump := fmt.Sprintf(`{
		"group_ids": "%v",
		"v": "%v"
	}`, groupID*-1, // умножаем на -1, потому что id группы в этом методе должен быть без минуса
		"5.95")
	values, err := MakeJSON(jsonDump)
	if err != nil {
		return vkCommunity, err
	}

	// отправляем запрос
	response, err := SendVKAPIQuery(sender, "groups.getById", values, subject)
	if err != nil {
		return vkCommunity, err
	}

	// заранее извлекаем из общей карты карту с данными о сообществе
	groupInfoMap := response["response"].([]interface{})[0]
	groupInfo := groupInfoMap.(map[string]interface{})

	// собираем данные о сообществе из ответа сервера
	vkCommunity.ID = int(groupInfo["id"].(float64))
	vkCommunity.Name = groupInfo["name"].(string)
	vkCommunity.ScreenName = groupInfo["screen_name"].(string)

	return vkCommunity, nil
}

// VKUser - структура данных о пользователе VK
type VKUser struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// GetUserInfo получает имя и фамилию пользователя по его id
func GetUserInfo(sender string, subject Subject, userID int) (VKUser, error) {
	var vkUser VKUser

	// формируем json запроса к vk api
	jsonDump := fmt.Sprintf(`{
		"user_ids": "%v",
		"v": "%v"
	}`, userID, "5.95")
	values, err := MakeJSON(jsonDump)
	if err != nil {
		return vkUser, err
	}

	// отправляем запрос
	response, err := SendVKAPIQuery(sender, "users.get", values, subject)
	if err != nil {
		return vkUser, err
	}

	// заранее извлекаем из общей карты карту с данными о сообществе
	userInfoMap := response["response"].([]interface{})[0]
	userInfo := userInfoMap.(map[string]interface{})

	// собираем данные о сообществе из ответа сервера
	vkUser.ID = int(userInfo["id"].(float64))
	vkUser.FirstName = userInfo["first_name"].(string)
	vkUser.LastName = userInfo["last_name"].(string)

	return vkUser, nil
}
