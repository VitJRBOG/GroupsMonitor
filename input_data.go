package main

import (
	"fmt"
)

// InputAccessToken принимает ввод токена доступа пользователем в консоль и возвращает его
func InputAccessToken(name string) (string, error) {
	fmt.Printf("USER [New access token for %v]: ", name)
	var userAnswer string
	_, err := fmt.Scan(&userAnswer)
	if err != nil {
		return "", err
	}
	return userAnswer, nil
}
