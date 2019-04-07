package tools

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// Method отправляет post-запрос к VK API
// возвращает ответ сервера и ошибку
func Method(methodName string, valuesBytes []byte, accessToken string) ([]byte, error) {
	url := "https://api.vk.com/method/"
	url += methodName
	url += "?access_token=" + accessToken

	// преобразуем массив байт в словарь
	var f interface{}
	err := json.Unmarshal(valuesBytes, &f)
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
	response, err := http.Post(url, "", nil)
	if err != nil {
		return nil, err
	}

	// преобразуем полученный ответ сервера в массив байт
	body, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		return nil, err
	}

	return body, nil
}
