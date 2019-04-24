package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// UnixTimeStampToDate преобразовывает дату из unix time stamp в читабельный вид
func UnixTimeStampToDate(timestampDate int) string {
	tm := time.Unix(int64(timestampDate), 0)
	timeFormat := "02.01.2006 15:04:05"
	readableTime := tm.Format(timeFormat)
	return readableTime
}

// MakeJSON формирует json-словарь
func MakeJSON(jsonDump string) ([]byte, error) {

	// сначала собираем карту из полученной json-строки
	var values interface{}
	valuesBytes := []byte(jsonDump)
	err := json.Unmarshal(valuesBytes, &values)
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

// ListFromDB - структура для списка из параметра модуля мониторинга
type ListFromDB struct {
	List []string `json:"list"`
}

// MakeParamList формирует список из массива в json
func MakeParamList(jsonDump string) (ListFromDB, error) {
	var values ListFromDB

	// собираем структуру из полученной json-строки
	valuesBytes := []byte(jsonDump)
	err := json.Unmarshal(valuesBytes, &values)
	if err != nil {
		return values, err
	}

	return values, nil
}

// CharChangeChars - структура для хранение символов и их аналогий в другом языке
type CharChangeChars struct {
	CyrChars []string
	LatChars []string
}

// CharChange выполняет замену символов с кириллических на латинские и наоборот
func CharChange(text string, changeType string) string {
	var charChangeChars CharChangeChars
	charChangeChars.CyrChars = []string{"С", "Р", "у", "Х", "о", "р", "а", "с", "х", "Е", "е", "А", "О"}
	charChangeChars.LatChars = []string{"C", "P", "y", "X", "o", "p", "a", "c", "x", "E", "e", "A", "O"}

	// замена символов не происходит

	switch changeType {
	case "lat_to_cyr":
		for i, symb := range charChangeChars.LatChars {
			text = strings.Replace(text, symb, charChangeChars.CyrChars[i], -1)
		}
	case "cyr_to_lat":
		for i, symb := range charChangeChars.CyrChars {
			text = strings.Replace(text, symb, charChangeChars.LatChars[i], -1)
		}
	}
	return text
}

// GetAccessToken получает токен доступа из БД по названию метода и id субъекта
func GetAccessToken(methodName string, subject Subject, monitorID int) (AccessToken, error) {
	var accessToken AccessToken

	// запрашиваем из БД данные по методу vk api
	method, err := SelectDBMethod(methodName, subject.ID, monitorID)
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

	// запрашиваем все данные по указанному токену
	accessToken, err := SelectDBAccessTokenByName(nameAccessToken)
	if err != nil {
		return err
	}

	// проверяем существование токена с таким именем в БД
	if accessToken.ID == 0 {
		if accessToken.Name == "" {
			if accessToken.Value == "" {
				message := fmt.Sprintf(`Error: access token with name "%v" is not exist.`, nameAccessToken)
				OutputMessage(sender, message)
				return nil
			}
		}
	}

	// сообщаем пользователю о запуске алгоритма обновления токена
	message := fmt.Sprintf("Request new access token for %v.", nameAccessToken)
	OutputMessage(sender, message)

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
