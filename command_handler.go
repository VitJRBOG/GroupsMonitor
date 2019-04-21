package main

import (
	"fmt"
	"os"
	"time"
)

// ListenUserCommands принимает и обрабатывает консольные команды от пользователя
func ListenUserCommands(threads []*Thread) error {

	// в бесконечном цикле будем проверять ввод в консоль
	for true {
		var userAnswer string
		if _, err := fmt.Scan(&userAnswer); err != nil {
			return err
		}

		// после ввода проверяем содержимое строки
		switch userAnswer {

		// команда на обновление токена доступа
		case "upd_at":
			if err := updateAccessToken(); err != nil {
				return err
			}

		// команда на нормальную остановку потоков
		case "stop":
			stopThreads(threads)

		// команда на принудительное завершение работы
		case "quit":
			forceQuit()
		}
	}

	return nil
}

func updateAccessToken() error {
	sender := "Update access token"

	// запрашиваем у пользователя имя обновляемого токена
	nameAccessToken, err := InputNameAccessToken()
	if err != nil {
		return err
	}

	// проверяем ввод на команду отмены
	if nameAccessToken == "cancel" {
		message := "Abort."
		OutputMessage(sender, message)
		return nil
	}

	// вызываем функцию получения нового токена
	if err = GetNewAccessToken(sender, nameAccessToken); err != nil {
		return err
	}

	return nil
}

func stopThreads(threads []*Thread) {

	// пробегаем по всем потокам и выставляем флаги на остановку
	for _, thread := range threads {
		thread.StopFlag = 1
	}

	// проверяем успешность остановки потоков
	repeats := 60
	cantStop := len(threads)
	for i := 0; i < repeats; i++ {
		for _, thread := range threads {
			// если поток имеет статус stopped, то обнуляем ссылку на него
			if thread != nil {
				if thread.Status == "stopped" {
					sender := thread.Name
					message := "OK! Monitoring is stopped!"
					OutputMessage(sender, message)
					thread = nil
					cantStop--
				}
			}
		}
		// если остались работающие потоки, то вызываем задержку
		// если все потоки завершились, то сообщаем об этом пользователю и завершаем работу
		if cantStop > 0 {
			interval := 1
			time.Sleep(time.Duration(interval) * time.Second)
		} else {
			sender := "Core"
			message := "All threads is stopped. Quit..."
			OutputMessage(sender, message)
			os.Exit(0)
		}
	}

	// если после максимального количества повторов проверки еще остались работающие потоки,
	// то сообщаем об этом пользователю
	for _, thread := range threads {
		if thread != nil {
			sender := thread.Name
			message := fmt.Sprintf("WARNING! Can't stop thread.")
			OutputMessage(sender, message)
		}
	}

	// принудительное завершение работы, если дело дошло до этой строки
	forceQuit()
}

func forceQuit() {
	// сообщаем о принудительном завершении работы и выходим
	sender := "Core"
	message := fmt.Sprintf("Force quit...")
	OutputMessage(sender, message)
	os.Exit(0)
}
