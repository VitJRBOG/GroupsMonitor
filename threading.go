package main

import (
	"fmt"
	"runtime/debug"
	"time"
)

// Thread - структура для хранения данных о потоке
type Thread struct {
	Name       string
	ActionFlag int // 0 - nothing, 1 - stopping, 2 - restarting, 3 - waiting
	Status     string
}

// MakeThreads создает и запускает потоки
func MakeThreads() ([]*Thread, error) {
	var threads []*Thread

	// получаем из БД список субъектов
	var dbKit DataBaseKit
	subjects, err := dbKit.selectTableSubject()
	if err != nil {
		return threads, err
	}

	// перебираем список субъектов и запускаем модули мониторинга в отдельных потоках
	for _, subject := range subjects {

		// получаем из БД параметры для модуля мониторинга постов на стене
		var wallPostMonitorParam WallPostMonitorParam
		err := wallPostMonitorParam.selectFromDBBySubjectID(subject.ID)
		if err != nil {
			ToLogFile(err.Error(), string(debug.Stack()))
			panic(err.Error())
		}

		// проверяем параметр и определяем, нужно ли запускать этот модуль
		if wallPostMonitorParam.NeedMonitoring == 1 {

			// создаем структуру с данными о потоке и наполняем ее данными
			var thread Thread
			thread.Name = fmt.Sprintf("%v's wall post monitoring", subject.Name)
			thread.Status = "waiting"

			// запускаем поток
			go wallPostMonitoring(&thread, subject)
			threads = append(threads, &thread)
		}

		// получаем из БД параметры для модуля мониторинга фотографий в альбомах
		var albumPhotoMonitorParam AlbumPhotoMonitorParam
		err = albumPhotoMonitorParam.selectFromDBBySubjectID(subject.ID)
		if err != nil {
			return threads, err
		}

		// проверяем параметр и определяем, нужно ли запускать этот модуль
		if albumPhotoMonitorParam.NeedMonitoring == 1 {

			// создаем структуру с данными о потоке и наполняем ее данными
			var thread Thread
			thread.Name = fmt.Sprintf("%v's album photo monitoring", subject.Name)
			thread.Status = "waiting"

			// запускаем поток
			go albumPhotoMonitoring(&thread, subject)
			threads = append(threads, &thread)
		}

		// получаем из БД параметры для модуля мониторинга видео в альбомах
		var videoMonitorParam VideoMonitorParam
		err = videoMonitorParam.selectFromDBBySubjectID(subject.ID)
		if err != nil {
			return threads, err
		}

		// проверяем параметр и определяем, нужно ли запускать этот модуль
		if videoMonitorParam.NeedMonitoring == 1 {

			// создаем структуру с данными о потоке и наполняем ее данными
			var thread Thread
			thread.Name = fmt.Sprintf("%v's video monitoring", subject.Name)
			thread.Status = "waiting"

			// запускаем поток
			go videoMonitoring(&thread, subject)
			threads = append(threads, &thread)
		}

		// получаем из БД параметры для модуля мониторинга комментариев под фотографиями
		var photoCommentMonitorParam PhotoCommentMonitorParam
		err = photoCommentMonitorParam.selectFromDBBySubjectID(subject.ID)
		if err != nil {
			return threads, err
		}

		// проверяем параметр и определяем, нужно ли запускать этот модуль
		if photoCommentMonitorParam.NeedMonitoring == 1 {

			// создаем структуру с данными о потоке и наполняем ее данными
			var thread Thread
			thread.Name = fmt.Sprintf("%v's photo comment monitoring", subject.Name)
			thread.Status = "waiting"

			// запускаем поток
			go photoCommentMonitoring(&thread, subject)
			threads = append(threads, &thread)
		}

		// получаем из БД параметры для модуля мониторинга комментариев под видеозаписями
		videoCommentMonitorParam, err := SelectDBVideoCommentMonitorParam(subject.ID)
		if err != nil {
			return threads, err
		}

		// проверяем параметр и определяем, нужно ли запускать этот модуль
		if videoCommentMonitorParam.NeedMonitoring == 1 {

			// создаем структуру с данными о потоке и наполняем ее данными
			var thread Thread
			thread.Name = fmt.Sprintf("%v's video comment monitoring", subject.Name)
			thread.Status = "waiting"

			// запускаем поток
			go videoCommentMonitoring(&thread, subject)
			threads = append(threads, &thread)
		}

		// получаем из БД параметры для модуля мониторинга комментариев в топиках обсуждений
		topicMonitorParam, err := SelectDBTopicMonitorParam(subject.ID)
		if err != nil {
			return threads, err
		}

		// проверяем параметр и определяем, нужно ли запускать этот модуль
		if topicMonitorParam.NeedMonitoring == 1 {

			// создаем структуру с данными о потоке и наполняем ее данными
			var thread Thread
			thread.Name = fmt.Sprintf("%v's topic monitoring", subject.Name)
			thread.Status = "waiting"

			// запускаем поток
			go topicMonitoring(&thread, subject)
			threads = append(threads, &thread)
		}

		// получаем из БД параметры для модуля мониторинга комментариев под постами на стене
		wallPostCommentMonitorParam, err := SelectDBWallPostCommentMonitorParam(subject.ID)
		if err != nil {
			return threads, err
		}

		// проверяем параметр и определяем, нужно ли запускать этот модуль
		if wallPostCommentMonitorParam.NeedMonitoring == 1 {

			// создаем структуру с данными о потоке и наполняем ее данными
			var thread Thread
			thread.Name = fmt.Sprintf("%v's wall post comment monitoring", subject.Name)
			thread.Status = "waiting"

			// запускаем поток
			go wallPostCommentMonitoring(&thread, subject)
			threads = append(threads, &thread)
		}

	}

	if len(threads) == 0 {
		sender := "Core"
		message := "WARNING! No thread has been created."
		OutputMessage(sender, message)
	}

	// проверяем количество созданных потоков
	if len(threads) > 0 {

		// если их больше 0, то запускаем функцию поиска потоков, завершивших свою работу из-за ошибки
		go threadsStatusMonitoring(threads)
	}

	return threads, nil
}

// threadsStatusMonitoring ищет потоки, завершившие свою работу из-за ошибки
func threadsStatusMonitoring(threads []*Thread) {

	// перебираем структуры с данными о потоках
	for j, thread := range threads {

		// если статус потока "error", то сообщаем об этом пользователю
		if thread.Status == "error" {
			message := "WARNING! Thread is stopped with error!"
			OutputMessage(thread.Name, message)
			threads[j] = nil
		}
	}

	// после завершения перебора включаем режим ожидания
	time.Sleep(10 * time.Second)
}

func wallPostMonitoring(threadData *Thread, subject Subject) {

	// включаем режим ожидания команды запуска
	needWaiting := true
	for needWaiting {
		// если статус - ожидание, то ждем
		if threadData.Status == "waiting" {
			time.Sleep(1 * time.Second)
		} else {
			// если нет, то выходим из цикла
			needWaiting = false
		}
	}

	// сообщаем пользователю о запуске модуля
	sender := threadData.Name
	message := "Started..."
	OutputMessage(sender, message)

	// создаем счетчик ошибок
	errorsCounter := 0

	// запускаем бесконечный цикл
	for true {
		// заранее присваиваем значение интервала
		interval := 20

		// запускаем функцию мониторинга
		wallPostMonitorParam, err := WallPostMonitor(subject)
		if err != nil {
			// если функция вернула ошибку, то увеличиваем счетчик на 1
			errorsCounter++
			// если в результате счетчик не стал равен 4, то продолжаем
			if errorsCounter < 4 {
				// сообщаем пользователю об ошибке
				sender := fmt.Sprintf("%v -> Thread control", threadData.Name)
				// 20-секундный таймер умножаем на количество ошибок
				interval *= errorsCounter
				message := fmt.Sprintf("ERROR: %v. Time out for %ds", err, interval)
				OutputMessage(sender, message)
			} else {
				// если стал, то сообщаем об этом пользователю
				message := fmt.Sprintf("ERROR: %v. Thread is paused. Type \"restart\" for turn on again...", err)
				OutputMessage(sender, message)
				// и ставим потоку статус waiting
				threadData.Status = "waiting"
			}
		}

		// после успешного завершения работы функции мониторинга получаем значение интервала
		if wallPostMonitorParam != nil {
			interval = wallPostMonitorParam.Interval
		}

		// и включаем режим ожидания
		for i := 0; i < interval; i++ {
			// если статус потока waiting
			if threadData.Status == "waiting" {
				// то каждый раз обнуляем i, и тем самым вводим поток в вечное ожидание
				i = 0
			}
			time.Sleep(1 * time.Second)

			// периодически проверяем, был ли выставлен флаг остановки
			if threadData.ActionFlag == 1 {
				// если был, то меняем статус потока на "stopped"
				threadData.Status = "stopped"
			}

			// если выставлен флаг рестарта
			if threadData.ActionFlag == 2 {
				// то обновляем статус потока
				threadData.Status = "alive"
				// и перезапускаем функцию
				threadData.ActionFlag = 0
				wallPostMonitoring(threadData, subject)
			}

			// если выставлен флаг ожидания
			if threadData.ActionFlag == 3 {
				// то обновляем статус потока
				threadData.Status = "waiting"
			}
		}

		// если статус потока "stopped", то завершаем его работу
		if threadData.Status == "stopped" {
			return
		}
	}
}

func albumPhotoMonitoring(threadData *Thread, subject Subject) {

	// включаем режим ожидания команды запуска
	needWaiting := true
	for needWaiting {
		// если статус - ожидание, то ждем
		if threadData.Status == "waiting" {
			time.Sleep(1 * time.Second)
		} else {
			// если нет, то выходим из цикла
			needWaiting = false
		}
	}

	// сообщаем пользователю о запуске модуля
	sender := threadData.Name
	message := "Started..."
	OutputMessage(sender, message)

	// создаем счетчик ошибок
	errorsCounter := 0

	// запускаем бесконечный цикл
	for true {
		// заранее присваиваем значение интервала
		interval := 20

		// запускаем функцию мониторинга
		albumPhotoMonitorParam, err := AlbumPhotoMonitor(subject)
		if err != nil {
			// если функция вернула ошибку, то увеличиваем счетчик на 1
			errorsCounter++
			// если в результате счетчик не стал равен 4, то продолжаем
			if errorsCounter < 4 {
				// сообщаем пользователю об ошибке
				sender := fmt.Sprintf("%v -> Thread control", threadData.Name)
				// 20-секундный таймер умножаем на количество ошибок
				interval *= errorsCounter
				message := fmt.Sprintf("ERROR: %v. Time out for %ds", err, interval)
				OutputMessage(sender, message)
			} else {
				// если стал, то сообщаем об этом пользователю
				message := fmt.Sprintf("ERROR: %v. Thread is paused. Type \"restart\" for turn on again...", err)
				OutputMessage(sender, message)
				// и ставим потоку статус waiting
				threadData.Status = "waiting"
			}
		}

		// после успешного завершения работы функции мониторинга получаем значение интервала
		if albumPhotoMonitorParam != nil {
			interval = albumPhotoMonitorParam.Interval
		}

		// и включаем режим ожидания
		for i := 0; i < interval; i++ {
			// если статус потока waiting
			if threadData.Status == "waiting" {
				// то каждый раз обнуляем i, и тем самым вводим поток в вечное ожидание
				i = 0
			}
			time.Sleep(1 * time.Second)

			// периодически проверяем, был ли выставлен флаг остановки
			if threadData.ActionFlag == 1 {
				// если был, то меняем статус потока на "stopped"
				threadData.Status = "stopped"
			}

			// если выставлен флаг рестарта
			if threadData.ActionFlag == 2 {
				// то обновляем статус потока
				threadData.Status = "alive"
				// и перезапускаем функцию
				threadData.ActionFlag = 0
				albumPhotoMonitoring(threadData, subject)
			}

			// если выставлен флаг ожидания
			if threadData.ActionFlag == 3 {
				// то обновляем статус потока
				threadData.Status = "waiting"
			}
		}

		// если статус потока "stopped", то завершаем его работу
		if threadData.Status == "stopped" {
			return
		}
	}
}

func videoMonitoring(threadData *Thread, subject Subject) {

	// включаем режим ожидания команды запуска
	needWaiting := true
	for needWaiting {
		// если статус - ожидание, то ждем
		if threadData.Status == "waiting" {
			time.Sleep(1 * time.Second)
		} else {
			// если нет, то выходим из цикла
			needWaiting = false
		}
	}

	// сообщаем пользователю о запуске модуля
	sender := threadData.Name
	message := "Started..."
	OutputMessage(sender, message)

	// создаем счетчик ошибок
	errorsCounter := 0

	// запускаем бесконечный цикл
	for true {
		// заранее присваиваем значение интервала
		interval := 20

		// запускаем функцию мониторинга
		videoMonitorParam, err := VideoMonitor(subject)
		if err != nil {
			// если функция вернула ошибку, то увеличиваем счетчик на 1
			errorsCounter++
			// если в результате счетчик не стал равен 4, то продолжаем
			if errorsCounter < 4 {
				// сообщаем пользователю об ошибке
				sender := fmt.Sprintf("%v -> Thread control", threadData.Name)
				// 20-секундный таймер умножаем на количество ошибок
				interval *= errorsCounter
				message := fmt.Sprintf("ERROR: %v. Time out for %ds", err, interval)
				OutputMessage(sender, message)
			} else {
				// если стал, то сообщаем об этом пользователю
				message := fmt.Sprintf("ERROR: %v. Thread is paused. Type \"restart\" for turn on again...", err)
				OutputMessage(sender, message)
				// и ставим потоку статус waiting
				threadData.Status = "waiting"
			}
		}

		// после успешного завершения работы функции мониторинга получаем значение интервала
		if videoMonitorParam != nil {
			interval = videoMonitorParam.Interval
			// и обнуляем счетчик ошибок
			errorsCounter = 0
		}

		// и включаем режим ожидания
		for i := 0; i < interval; i++ {
			// если статус потока waiting
			if threadData.Status == "waiting" {
				// то каждый раз обнуляем i, и тем самым вводим поток в вечное ожидание
				i = 0
			}
			time.Sleep(1 * time.Second)

			// периодически проверяем, был ли выставлен флаг остановки
			if threadData.ActionFlag == 1 {
				// если был, то меняем статус потока на "stopped"
				threadData.Status = "stopped"
			}

			// если выставлен флаг рестарта
			if threadData.ActionFlag == 2 {
				// то обновляем статус потока
				threadData.Status = "alive"
				// и перезапускаем функцию
				threadData.ActionFlag = 0
				videoMonitoring(threadData, subject)
			}

			// если выставлен флаг ожидания
			if threadData.ActionFlag == 3 {
				// то обновляем статус потока
				threadData.Status = "waiting"
			}
		}

		// если статус потока "stopped", то завершаем его работу
		if threadData.Status == "stopped" {
			return
		}
	}
}

func photoCommentMonitoring(threadData *Thread, subject Subject) {

	// включаем режим ожидания команды запуска
	needWaiting := true
	for needWaiting {
		// если статус - ожидание, то ждем
		if threadData.Status == "waiting" {
			time.Sleep(1 * time.Second)
		} else {
			// если нет, то выходим из цикла
			needWaiting = false
		}
	}

	// сообщаем пользователю о запуске модуля
	sender := threadData.Name
	message := "Started..."
	OutputMessage(sender, message)

	// создаем счетчик ошибок
	errorsCounter := 0

	// запускаем бесконечный цикл
	for true {
		// заранее присваиваем значение интервала
		interval := 20

		// запускаем функцию мониторинга
		photoCommentMonitorParam, err := PhotoCommentMonitor(subject)
		if err != nil {
			// если функция вернула ошибку, то увеличиваем счетчик на 1
			errorsCounter++
			// если в результате счетчик не стал равен 4, то продолжаем
			if errorsCounter < 4 {
				// сообщаем пользователю об ошибке
				sender := fmt.Sprintf("%v -> Thread control", threadData.Name)
				// 20-секундный таймер умножаем на количество ошибок
				interval *= errorsCounter
				message := fmt.Sprintf("ERROR: %v. Time out for %ds", err, interval)
				OutputMessage(sender, message)
			} else {
				// если стал, то сообщаем об этом пользователю
				message := fmt.Sprintf("ERROR: %v. Thread is paused. Type \"restart\" for turn on again...", err)
				OutputMessage(sender, message)
				// и ставим потоку статус waiting
				threadData.Status = "waiting"
			}
		}

		// после успешного завершения работы функции мониторинга получаем значение интервала
		if photoCommentMonitorParam != nil {
			interval = photoCommentMonitorParam.Interval
			// и обнуляем счетчик ошибок
			errorsCounter = 0
		}

		// и включаем режим ожидания
		for i := 0; i < interval; i++ {
			// если статус потока waiting
			if threadData.Status == "waiting" {
				// то каждый раз обнуляем i, и тем самым вводим поток в вечное ожидание
				i = 0
			}
			time.Sleep(1 * time.Second)

			// периодически проверяем, был ли выставлен флаг остановки
			if threadData.ActionFlag == 1 {
				// если был, то меняем статус потока на "stopped"
				threadData.Status = "stopped"
			}

			// если выставлен флаг рестарта
			if threadData.ActionFlag == 2 {
				// то обновляем статус потока
				threadData.Status = "alive"
				// и перезапускаем функцию
				threadData.ActionFlag = 0
				photoCommentMonitoring(threadData, subject)
			}

			// если выставлен флаг ожидания
			if threadData.ActionFlag == 3 {
				// то обновляем статус потока
				threadData.Status = "waiting"
			}
		}

		// если статус потока "stopped", то завершаем его работу
		if threadData.Status == "stopped" {
			return
		}
	}
}

func videoCommentMonitoring(threadData *Thread, subject Subject) {

	// включаем режим ожидания команды запуска
	needWaiting := true
	for needWaiting {
		// если статус - ожидание, то ждем
		if threadData.Status == "waiting" {
			time.Sleep(1 * time.Second)
		} else {
			// если нет, то выходим из цикла
			needWaiting = false
		}
	}

	// сообщаем пользователю о запуске модуля
	sender := threadData.Name
	message := "Started..."
	OutputMessage(sender, message)

	// создаем счетчик ошибок
	errorsCounter := 0

	// запускаем бесконечный цикл
	for true {
		// заранее присваиваем значение интервала
		interval := 20

		// запускаем функцию мониторинга
		videoCommentMonitorParam, err := VideoCommentMonitor(subject)
		if err != nil {
			// если функция вернула ошибку, то увеличиваем счетчик на 1
			errorsCounter++
			// если в результате счетчик не стал равен 4, то продолжаем
			if errorsCounter < 4 {
				// сообщаем пользователю об ошибке
				sender := fmt.Sprintf("%v -> Thread control", threadData.Name)
				// 20-секундный таймер умножаем на количество ошибок
				interval *= errorsCounter
				message := fmt.Sprintf("ERROR: %v. Time out for %ds", err, interval)
				OutputMessage(sender, message)
			} else {
				// если стал, то сообщаем об этом пользователю
				message := fmt.Sprintf("ERROR: %v. Thread is paused. Type \"restart\" for turn on again...", err)
				OutputMessage(sender, message)
				// и ставим потоку статус waiting
				threadData.Status = "waiting"
			}
		}

		// после успешного завершения работы функции мониторинга получаем значение интервала
		if videoCommentMonitorParam != nil {
			interval = videoCommentMonitorParam.Interval
			// и обнуляем счетчик ошибок
			errorsCounter = 0
		}

		// и включаем режим ожидания
		for i := 0; i < interval; i++ {
			// если статус потока waiting
			if threadData.Status == "waiting" {
				// то каждый раз обнуляем i, и тем самым вводим поток в вечное ожидание
				i = 0
			}
			time.Sleep(1 * time.Second)

			// периодически проверяем, был ли выставлен флаг остановки
			if threadData.ActionFlag == 1 {
				// если был, то меняем статус потока на "stopped"
				threadData.Status = "stopped"
			}

			// если выставлен флаг рестарта
			if threadData.ActionFlag == 2 {
				// то обновляем статус потока
				threadData.Status = "alive"
				// и перезапускаем функцию
				threadData.ActionFlag = 0
				videoCommentMonitoring(threadData, subject)
			}

			// если выставлен флаг ожидания
			if threadData.ActionFlag == 3 {
				// то обновляем статус потока
				threadData.Status = "waiting"
			}
		}

		// если статус потока "stopped", то завершаем его работу
		if threadData.Status == "stopped" {
			return
		}
	}
}

func topicMonitoring(threadData *Thread, subject Subject) {

	// включаем режим ожидания команды запуска
	needWaiting := true
	for needWaiting {
		// если статус - ожидание, то ждем
		if threadData.Status == "waiting" {
			time.Sleep(1 * time.Second)
		} else {
			// если нет, то выходим из цикла
			needWaiting = false
		}
	}

	// сообщаем пользователю о запуске модуля
	sender := threadData.Name
	message := "Started..."
	OutputMessage(sender, message)

	// создаем счетчик ошибок
	errorsCounter := 0

	// запускаем бесконечный цикл
	for true {
		// заранее присваиваем значение интервала
		interval := 20

		// запускаем функцию мониторинга
		topicMonitorParam, err := TopicMonitor(subject)
		if err != nil {
			// если функция вернула ошибку, то увеличиваем счетчик на 1
			errorsCounter++
			// если в результате счетчик не стал равен 4, то продолжаем
			if errorsCounter < 4 {
				// сообщаем пользователю об ошибке
				sender := fmt.Sprintf("%v -> Thread control", threadData.Name)
				// 20-секундный таймер умножаем на количество ошибок
				interval *= errorsCounter
				message := fmt.Sprintf("ERROR: %v. Time out for %ds", err, interval)
				OutputMessage(sender, message)
			} else {
				// если стал, то сообщаем об этом пользователю
				message := fmt.Sprintf("ERROR: %v. Thread is paused. Type \"restart\" for turn on again...", err)
				OutputMessage(sender, message)
				// и ставим потоку статус waiting
				threadData.Status = "waiting"
			}
		}

		// после успешного завершения работы функции мониторинга получаем значение интервала
		if topicMonitorParam != nil {
			interval = topicMonitorParam.Interval
		}

		// и включаем режим ожидания
		for i := 0; i < interval; i++ {
			// если статус потока waiting
			if threadData.Status == "waiting" {
				// то каждый раз обнуляем i, и тем самым вводим поток в вечное ожидание
				i = 0
			}
			time.Sleep(1 * time.Second)

			// периодически проверяем, был ли выставлен флаг остановки
			if threadData.ActionFlag == 1 {
				// если был, то меняем статус потока на "stopped"
				threadData.Status = "stopped"
			}

			// если выставлен флаг рестарта
			if threadData.ActionFlag == 2 {
				// то обновляем статус потока
				threadData.Status = "alive"
				// и перезапускаем функцию
				threadData.ActionFlag = 0
				topicMonitoring(threadData, subject)
			}

			// если выставлен флаг ожидания
			if threadData.ActionFlag == 3 {
				// то обновляем статус потока
				threadData.Status = "waiting"
			}
		}

		// если статус потока "stopped", то завершаем его работу
		if threadData.Status == "stopped" {
			return
		}
	}
}

func wallPostCommentMonitoring(threadData *Thread, subject Subject) {

	// включаем режим ожидания команды запуска
	needWaiting := true
	for needWaiting {
		// если статус - ожидание, то ждем
		if threadData.Status == "waiting" {
			time.Sleep(1 * time.Second)
		} else {
			// если нет, то выходим из цикла
			needWaiting = false
		}
	}

	// сообщаем пользователю о запуске модуля
	sender := threadData.Name
	message := "Started..."
	OutputMessage(sender, message)

	// создаем счетчик ошибок
	errorsCounter := 0

	// запускаем бесконечный цикл
	for true {
		// заранее присваиваем значение интервала
		interval := 20

		// запускаем функцию мониторинга
		wallPostCommentMonitorParam, err := WallPostCommentMonitor(subject)
		if err != nil {
			// если функция вернула ошибку, то увеличиваем счетчик на 1
			errorsCounter++
			// если в результате счетчик не стал равен 4, то продолжаем
			if errorsCounter < 4 {
				// сообщаем пользователю об ошибке
				sender := fmt.Sprintf("%v -> Thread control", threadData.Name)
				// 20-секундный таймер умножаем на количество ошибок
				interval *= errorsCounter
				message := fmt.Sprintf("ERROR: %v. Time out for %ds", err, interval)
				OutputMessage(sender, message)
			} else {
				// если стал, то сообщаем об этом пользователю
				message := fmt.Sprintf("ERROR: %v. Thread is paused. Type \"restart\" for turn on again...", err)
				OutputMessage(sender, message)
				// и ставим потоку статус waiting
				threadData.Status = "waiting"
			}
		}

		// после успешного завершения работы функции мониторинга получаем значение интервала
		if wallPostCommentMonitorParam != nil {
			interval = wallPostCommentMonitorParam.Interval
		}

		// и включаем режим ожидания
		for i := 0; i < interval; i++ {
			// если статус потока waiting
			if threadData.Status == "waiting" {
				// то каждый раз обнуляем i, и тем самым вводим поток в вечное ожидание
				i = 0
			}
			time.Sleep(1 * time.Second)

			// периодически проверяем, был ли выставлен флаг остановки
			if threadData.ActionFlag == 1 {
				// если был, то меняем статус потока на "stopped"
				threadData.Status = "stopped"
			}

			// если выставлен флаг рестарта
			if threadData.ActionFlag == 2 {
				// то обновляем статус потока
				threadData.Status = "alive"
				// и перезапускаем функцию
				threadData.ActionFlag = 0
				wallPostCommentMonitoring(threadData, subject)
			}

			// если выставлен флаг ожидания
			if threadData.ActionFlag == 3 {
				// то обновляем статус потока
				threadData.Status = "waiting"
			}
		}

		// если статус потока "stopped", то завершаем его работу
		if threadData.Status == "stopped" {
			return
		}
	}
}
