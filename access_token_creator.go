package main

import "fmt"

// CreateAccessToken создает субъект и наполняет его параметрами
func CreateAccessToken() error {
	var accessToken AccessToken

	// получаем имя токена доступа
	name, err := getAccessTokenName()
	if err != nil {
		return err
	}
	accessToken.Name = name

	// получаем значение токена доступа
	value, err := getValue()
	if err != nil {
		return err
	}
	accessToken.Value = value

	// создаем новое поле в таблице
	err = InsertDBAccessToken(accessToken)
	if err != nil {
		return err
	}

	// сообщаем пользователю о том, что токен доступа создан
	sender := "> [Create access token]: "
	message := "Access token has been successfully created."
	OutputMessage(sender, message)

	return nil
}

// getAccessTokenName запрашивает у пользователя название токена для своей базы данных
func getAccessTokenName() (string, error) {
	fmt.Print("> [Create access token -> Get name]: ")
	var userAnswer string
	_, err := fmt.Scan(&userAnswer)
	if err != nil {
		return "", err
	}
	return userAnswer, nil
}

// getValue запрашивает у пользователя значение токена доступа
func getValue() (string, error) {
	fmt.Print("> [Create access token -> Get value]: ")
	var userAnswer string
	_, err := fmt.Scan(&userAnswer)
	if err != nil {
		return "", err
	}
	return userAnswer, nil
}
