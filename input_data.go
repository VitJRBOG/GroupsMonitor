package main

import (
	"fmt"
)

// InputNameAccessToken принимает ввод имени токена доступа пользователем в консоль и возвращает его
func InputNameAccessToken() (string, error) {
	fmt.Print("\n> [Name of access token for update]: ")
	var userAnswer string
	_, err := fmt.Scan(&userAnswer)
	if err != nil {
		return "", err
	}
	return userAnswer, nil
}

// InputAccessToken принимает ввод токена доступа пользователем в консоль и возвращает его
func InputAccessToken(name string) (string, error) {
	fmt.Printf("\n> [New access token for %v]: ", name)
	var userAnswer string
	_, err := fmt.Scan(&userAnswer)
	if err != nil {
		return "", err
	}
	return userAnswer, nil
}
