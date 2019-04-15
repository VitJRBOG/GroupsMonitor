package main

import (
	"fmt"
)

// InputData принимает ввод данных пользователем в консоль и возвращает их
func InputData(sender string, message string) (string, error) {
	fmt.Print("USER [" + sender + "] " + message + ": ")
	var userAnswer string
	_, err := fmt.Scan(&userAnswer)
	if err != nil {
		return "", err
	}
	return userAnswer, nil
}
