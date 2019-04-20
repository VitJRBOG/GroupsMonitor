package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime"
	"strings"
	"time"
)

// SendVKAPIQuery отправляет post-запрос к VK API
// возвращает ответ сервера и ошибку
func SendVKAPIQuery(sender string, methodName string,
	valuesBytes []byte, subject Subject) (map[string]interface{}, error) {

	// получаем токен доступа для данного метода vk api и субъекта
	accessToken, err := GetAccessToken(methodName, subject)

	// формируем url для запроса к vk api
	url := "https://api.vk.com/method/"
	url += methodName
	url += "?access_token=" + accessToken.Value

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
		url += "&" + key + "=" + values[key].(string)
	}

	// отправляем запрос
	requestResult, err := http.Post(url, "", nil)
	if err != nil {
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

	if exist := response["response"]; exist == true {
		return response, nil
	}

	// проверяем ошибки, которые вернул сервер vk api
	if _, exist := response["error"]; exist == true {
		errorItem := response["error"]
		typeError, causeError := VkAPIErrorHandler(errorItem)

		// смотрим тип ошибки
		switch typeError {

		// задержка, если слишком часто отправляются запросы
		case "timeout error":
			interval := 2
			message := fmt.Sprintf("Error: %v. Timeout for %d seconds...", causeError, interval)
			OutputMessage(sender, message)
			time.Sleep(time.Duration(interval) * time.Second)
			return SendVKAPIQuery(sender, methodName, valuesBytes, subject)

		// задержка, если токен доступа невалидный
		case "access token error":
			interval := 60
			message := fmt.Sprintf("Error: %v. Need new %v's access token. Waiting for %d seconds...",
				causeError, accessToken.Name, interval)
			OutputMessage(sender, message)
			time.Sleep(time.Duration(interval) * time.Second)
			return SendVKAPIQuery(sender, methodName, valuesBytes, subject)

		// выход, если критическая ошибка
		case "fatal error":
			message := fmt.Sprintf("Error: %v. Exit...", causeError)
			OutputMessage(sender, message)
			runtime.Goexit()
		}
	}

	return response, nil
}

// VkAPIErrorHandler проверяет ошибки от VK API
func VkAPIErrorHandler(responseError interface{}) (string, string) {

	// список ошибок для вызова таймаута
	timeoutErrors := []string{
		"captcha needed", "failed to establish a new connection",
		"connection aborted", "internal server error", "response code 504",
		"response code 502", "many requests per second",
	}

	// список ошибок для запроса нового токена доступа
	accessTokenErrors := []string{
		"invalid access_token",
		"access_token was given to another ip address",
		"access_token has expired",
		"no access_token passed",
	}

	// список ошибок для выхода из потока
	fatalErrors := []string{
		"access denied",
	}

	// извлекаем текст ошибки из словаря от vk api
	errorMessage, _ := responseError.(map[string]interface{})["error_msg"].(string)

	// далее проверяем текст ошибки на наличие похожих в трех списках
	for _, item := range timeoutErrors {
		if strings.Contains(errorMessage, item) {
			return "timeout error", item
		}
	}
	for _, item := range accessTokenErrors {
		if strings.Contains(errorMessage, item) {
			return "access token error", item
		}
	}
	for _, item := range fatalErrors {
		if strings.Contains(errorMessage, item) {
			return "fatal error", item
		}
	}

	// если в списках этой ошибки нет, то так и пишем
	return "unknown error", ""
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
		"group_ids": "%d",
		"v": "5.95"
	}`, groupID)
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
	groupInfo := response["response"].([]map[string]interface{})[0]

	// собираем данные о сообществе из ответа сервера
	vkCommunity.ID = groupInfo["id"].(int)
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
		"user_ids": "%d",
		"v": "5.95"
	}`, userID)
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
	userInfo := response["response"].([]map[string]interface{})[0]

	// собираем данные о сообществе из ответа сервера
	vkUser.ID = userInfo["id"].(int)
	vkUser.FirstName = userInfo["first_name"].(string)
	vkUser.LastName = userInfo["last_name"].(string)

	return vkUser, nil
}
