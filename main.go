package main

import "log"

func main() {
	// запускаем функцию проверки токенов доступа
	err := CheckAccessTokens()
	if err != nil {
		log.Fatalln(err)
	}
	// запускаем функцию старта модулей мониторинга
	// ...
}
