package main

import (
	"encoding/json"
	"fmt"
)

// MakeJSON формирует json-словарь
func MakeJSON(jsonDump string) ([]byte, error) {

	// сначала собираем карту из полученной json-строки
	var values interface{}
	valuesBytes := []byte(jsonDump)
	err := json.Unmarshal(valuesBytes, values)
	if err != nil {
		return nil, err
	}

	// затем преобразуем эту карту в массив байт
	valuesBytes, err = json.Marshal(values)
	if err != nil {
		return nil, err
	}

	return valuesBytes, nil
}

// GetAccessToken получает токен доступа из БД по названию метода и id субъекта
func GetAccessToken(methodName string, subject Subject) (AccessToken, error) {
	var accessToken AccessToken

	// запрашиваем из БД данные по методу vk api
	method, err := SelectDBMethod(methodName, subject.ID)
	if err != nil {
		return accessToken, err
	}

	// запрашиваем из БД токен доступа, связанный с этим методом vk api
	accessToken, err = SelectDBAccessTokenByID(method.AccessTokenID)
	if err != nil {
		return accessToken, err
	}

	return accessToken, nil
}

// GetNewAccessToken запрашивает у пользователя новый токен из консоли и обновляет значение в БД
func GetNewAccessToken(sender string, nameAccessToken string) error {

	// сообщаем пользователю о запуске алгоритма обновления токена
	message := fmt.Sprintf("Request new access token for %v.", nameAccessToken)
	OutputMessage(sender, message)

	// запрашиваем все данные по указанному токену
	accessToken, err := SelectDBAccessTokenByName(nameAccessToken)
	if err != nil {
		return err
	}

	// запрашиваем у пользователя новый токен
	accessToken.Value, err = InputAccessToken(nameAccessToken)
	if err != nil {
		return err
	}

	// сохраняем новые данные в БД
	if err = UpdateDBAccessToken(accessToken); err != nil {
		return err
	}

	// сообщаем пользователю об успехе
	message = fmt.Sprintf("Access token for %v has been successfully updated!", nameAccessToken)
	OutputMessage(sender, message)

	return nil
}

// MakeCommunityHyperlink собирает гиперссылку на сообщество
func MakeCommunityHyperlink(vkCommunity VKCommunity) string {
	hyperlink := fmt.Sprintf("*%v (%v)", vkCommunity.ScreenName, vkCommunity.Name)
	return hyperlink
}

// MakeUserHyperlink собирает гиперссылку на пользователя
func MakeUserHyperlink(vkUser VKUser) string {
	hyperlink := fmt.Sprintf("*id%d (%v %v)", vkUser.ID, vkUser.FirstName, vkUser.LastName)
	return hyperlink
}
