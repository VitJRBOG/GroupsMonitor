package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

// CheckAccessTokens проверяет валидность токенов доступа из базы данных
func CheckAccessTokens() error {
	// получаем слайс со структурами AccessToken
	accessTokens, err := SelectDBAccessTokens()
	if err != nil {
		return err
	}
	// запускаем цикл перебора структур в слайсе
	for _, accessToken := range accessTokens {
		// передаем структуру на проверку
		updated, err := checkAccessToken(&accessToken, false)
		if err != nil {
			return err
		}
		// если у пользователя был запрошен новый токен, то обновляем его в базе данных
		if updated {
			err = UpdateDBAccessToken(accessToken)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func checkAccessToken(accessToken *AccessToken, updated bool) (bool, error) {
	// формируем json проверочного запроса к vk api
	type Values struct {
		UserIds string `json:"user_ids"`
		Version string `json:"v"`
	}
	jsonDump := `{
		"user_ids": "1",
		"v": "5.95"
	}`
	values := Values{}
	valuesBytes := []byte(jsonDump)
	err := json.Unmarshal(valuesBytes, &values)
	if err != nil {
		return updated, err
	}
	valuesBytes, err = json.Marshal(values)
	if err != nil {
		return updated, err
	}
	// отправляем проверочный запрос
	responseBytes, err := SendVKAPIQuery("users.get", valuesBytes, accessToken.Value)
	if err != nil {
		return updated, err
	}
	// преобразуем полученный ответ сервера vk api в карту
	var f interface{}
	err = json.Unmarshal(responseBytes, &f)
	if err != nil {
		return updated, err
	}
	response := f.(map[string]interface{})

	// проверяем наличие ключа
	// будет присутствовать только если запрос выполнился корректно
	_, exist := response["response"]
	if exist {
		return updated, nil
	}

	// если ключ response отсутствует, то проверяем error
	errorItem, exist := response["error"]
	if exist {
		// снова преобразуем интерфейс в карту
		errorMessage, _ := errorItem.(map[string]interface{})["error_msg"].(string)
		accessTokenErrors := []string{
			"invalid access_token",
			"access_token was given to another ip address",
			"access_token has expired",
		}
		// проверяем ошибку по списку
		// если хоть одна совпадает, то запрашиваем у пользователя новый токен
		for _, item := range accessTokenErrors {
			if strings.Contains(errorMessage, item) {
				accessToken.Value, err = getNewAccessToken(accessToken.Name, item)
				if err != nil {
					return updated, err
				}
				updated = true
				return checkAccessToken(accessToken, updated)
			}
		}
		// если ни одна не совпала, то выводим ошибку, полученную от сервера
		sender := "Access token checker"
		message := "Fatal error: " + errorMessage + "."
		OutputMessage(sender, message)
	} else {
		fmt.Println(response) // на случай, если и ключ error будет отсутствовать
	}
	return updated, nil
}

// getNewAccessToken запрашивает у пользователя новый токен из консоли и возвращает его
func getNewAccessToken(nameAccessToken string, errorCause string) (string, error) {
	// сначала сообщаем пользователю, что случилось с токеном
	sender := "Access token checker"
	message := "Access token of " + nameAccessToken + " is invalid: " + errorCause + "."
	OutputMessage(sender, message)

	// запрашиваем новый токен
	message = "New access token for " + nameAccessToken
	userAnswer, err := InputData(sender, message)
	if err != nil {
		return "", err
	}
	return userAnswer, nil
}
