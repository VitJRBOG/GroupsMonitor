package main

import (
	"fmt"
	"time"
)

// Thread - структура для хранения данных о потоке
type Thread struct {
	Name       string
	ActionFlag int // 0 - nothing, 1 - stopping, 2 - restarting
	Status     string
	Subject    Subject
}

func (thread *Thread) initWallPostMonitoring(subject Subject) {
	thread.ActionFlag = 0
	thread.Name = fmt.Sprintf("%v's wall post monitoring", subject.Name)
	thread.Status = "stopped"
	thread.Subject = subject
}

func (thread *Thread) runWallPostMonitoring() {
	thread.ActionFlag = 0
	go wallPostMonitoring(thread, thread.Subject)
}

func (thread *Thread) initAlbumPhotoMonitoring(subject Subject) {
	thread.ActionFlag = 0
	thread.Name = fmt.Sprintf("%v's album photo monitoring", subject.Name)
	thread.Status = "stopped"
	thread.Subject = subject
}

func (thread *Thread) runAlbumPhotoMonitoring() {
	thread.ActionFlag = 0
	go albumPhotoMonitoring(thread, thread.Subject)
}

func (thread *Thread) initVideoMonitoring(subject Subject) {
	thread.ActionFlag = 0
	thread.Name = fmt.Sprintf("%v's video monitoring", subject.Name)
	thread.Status = "stopped"
	thread.Subject = subject
}

func (thread *Thread) runVideoMonitoring() {
	thread.ActionFlag = 0
	go videoMonitoring(thread, thread.Subject)
}

func (thread *Thread) initPhotoCommentMonitoring(subject Subject) {
	thread.ActionFlag = 0
	thread.Name = fmt.Sprintf("%v's photo comment monitoring", subject.Name)
	thread.Status = "stopped"
	thread.Subject = subject
}

func (thread *Thread) runPhotoCommentMonitoring() {
	thread.ActionFlag = 0
	go photoCommentMonitoring(thread, thread.Subject)
}

func (thread *Thread) initVideoCommentMonitoring(subject Subject) {
	thread.ActionFlag = 0
	thread.Name = fmt.Sprintf("%v's video comment monitoring", subject.Name)
	thread.Status = "stopped"
	thread.Subject = subject
}

func (thread *Thread) runVideoCommentMonitoring() {
	thread.ActionFlag = 0
	go videoCommentMonitoring(thread, thread.Subject)
}

func (thread *Thread) initTopicMonitoring(subject Subject) {
	thread.ActionFlag = 0
	thread.Name = fmt.Sprintf("%v's topic monitoring", subject.Name)
	thread.Status = "stopped"
	thread.Subject = subject
}

func (thread *Thread) runTopicMonitoring() {
	thread.ActionFlag = 0
	go topicMonitoring(thread, thread.Subject)
}

func (thread *Thread) initWallPostCommentMonitoring(subject Subject) {
	thread.ActionFlag = 0
	thread.Name = fmt.Sprintf("%v's wall post comment monitoring", subject.Name)
	thread.Status = "stopped"
	thread.Subject = subject
}

func (thread *Thread) runWallPostCommentMonitoring() {
	thread.ActionFlag = 0
	go wallPostCommentMonitoring(thread, thread.Subject)
}

// InitThreads инициализирует потоки и заполняет данными о них список для них
func InitThreads() (*[]*Thread, error) {
	var threads []*Thread

	var dbKit DataBaseKit
	subjects, err := dbKit.selectTableSubject()
	if err != nil {
		return nil, err
	}

	for _, subject := range subjects {

		var wallPostMonitorParam WallPostMonitorParam
		err := wallPostMonitorParam.selectFromDBBySubjectID(subject.ID)
		if err != nil {
			return nil, err
		}
		var threadWPM Thread
		threadWPM.initWallPostMonitoring(subject)
		if wallPostMonitorParam.NeedMonitoring == 0 {
			threadWPM.Status = "inactive"
		}
		threads = append(threads, &threadWPM)

		var albumPhotoMonitorParam AlbumPhotoMonitorParam
		err = albumPhotoMonitorParam.selectFromDBBySubjectID(subject.ID)
		if err != nil {
			return nil, err
		}
		var threadAPM Thread
		threadAPM.initAlbumPhotoMonitoring(subject)
		if albumPhotoMonitorParam.NeedMonitoring == 0 {
			threadAPM.Status = "inactive"
		}
		threads = append(threads, &threadAPM)

		var videoMonitorParam VideoMonitorParam
		err = videoMonitorParam.selectFromDBBySubjectID(subject.ID)
		if err != nil {
			return nil, err
		}
		var threadVM Thread
		threadVM.initVideoMonitoring(subject)
		if videoMonitorParam.NeedMonitoring == 0 {
			threadVM.Status = "inactive"
		}
		threads = append(threads, &threadVM)

		var photoCommentMonitorParam PhotoCommentMonitorParam
		err = photoCommentMonitorParam.selectFromDBBySubjectID(subject.ID)
		if err != nil {
			return nil, err
		}
		var threadPCM Thread
		threadPCM.initPhotoCommentMonitoring(subject)
		if photoCommentMonitorParam.NeedMonitoring == 0 {
			threadPCM.Status = "inactive"
		}
		threads = append(threads, &threadPCM)

		var videoCommentMonitorParam VideoCommentMonitorParam
		err = videoCommentMonitorParam.selectFromDBBySubjectID(subject.ID)
		if err != nil {
			return nil, err
		}
		var threadVCM Thread
		threadVCM.initVideoCommentMonitoring(subject)
		if videoCommentMonitorParam.NeedMonitoring == 0 {
			threadVCM.Status = "inactive"
		}
		threads = append(threads, &threadVCM)

		var topicMonitorParam TopicMonitorParam
		err = topicMonitorParam.selectFromDBBySubjectID(subject.ID)
		if err != nil {
			return nil, err
		}
		var threadTM Thread
		threadTM.initTopicMonitoring(subject)
		if topicMonitorParam.NeedMonitoring == 0 {
			threadTM.Status = "inactive"
		}
		threads = append(threads, &threadTM)

		var wallPostCommentMonitorParam WallPostCommentMonitorParam
		err = wallPostCommentMonitorParam.selectFromDBBySubjectID(subject.ID)
		if err != nil {
			return nil, err
		}
		var threadWPCM Thread
		threadWPCM.initWallPostCommentMonitoring(subject)
		if wallPostCommentMonitorParam.NeedMonitoring == 0 {
			threadWPCM.Status = "inactive"
		}
		threads = append(threads, &threadWPCM)

	}

	if len(threads) == 0 {
		sender := "Core"
		message := "WARNING! No thread has been created."
		OutputMessage(sender, message)
	}

	// проверяем количество созданных потоков
	if len(threads) > 0 {

		// если их больше 0, то запускаем функцию поиска потоков, завершивших свою работу из-за ошибки
		go threadsStatusMonitoring(&threads)
	}

	return &threads, nil
}

// threadsStatusMonitoring ищет потоки, завершившие свою работу из-за ошибки
func threadsStatusMonitoring(threads *[]*Thread) {

	// перебираем структуры с данными о потоках
	for j, thread := range *threads {

		// если статус потока "error", то сообщаем об этом пользователю
		if thread.Status == "error" {
			message := "WARNING! Thread is stopped with error!"
			OutputMessage(thread.Name, message)
			(*threads)[j] = nil
		}
	}

	// после завершения перебора включаем режим ожидания
	time.Sleep(10 * time.Second)
}

func wallPostMonitoring(threadData *Thread, subject Subject) {

	// сообщаем пользователю о запуске модуля
	sender := threadData.Name
	message := "Started..."
	OutputMessage(sender, message)

	// создаем счетчик ошибок
	errorsCounter := 0

	// запускаем бесконечный цикл
	for true {
		// меняем статус потока
		threadData.Status = "alive"
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
				// и ставим потоку статус error
				threadData.Status = "error"
			}
		}

		// после успешного завершения работы функции мониторинга получаем значение интервала
		if wallPostMonitorParam != nil {
			interval = wallPostMonitorParam.Interval
		}

		// и включаем режим ожидания
		for i := 0; i < interval; i++ {
			// если статус потока error
			if threadData.Status == "error" {
				// то каждый раз обнуляем i, и тем самым вводим поток в вечное ожидание
				i = 0
			} else { // если нет, то ставим статус waiting и указываем оставшееся время ожидания
				threadData.Status = fmt.Sprintf("waiting for %d sec", interval-i)
			}
			time.Sleep(1 * time.Second)

			// периодически проверяем, был ли выставлен флаг остановки
			if threadData.ActionFlag == 1 {
				// если был, то меняем статус потока на "stopped"
				threadData.Status = "stopped"
				// и завершаем работу потока
				return
			}

			// если выставлен флаг рестарта
			if threadData.ActionFlag == 2 {
				// то обновляем статус потока
				threadData.Status = "alive"
				// и перезапускаем функцию
				threadData.ActionFlag = 0
				wallPostMonitoring(threadData, subject)
			}
		}
	}
}

func albumPhotoMonitoring(threadData *Thread, subject Subject) {

	// сообщаем пользователю о запуске модуля
	sender := threadData.Name
	message := "Started..."
	OutputMessage(sender, message)

	// создаем счетчик ошибок
	errorsCounter := 0

	// запускаем бесконечный цикл
	for true {
		// меняем статус потока
		threadData.Status = "alive"
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
				// и ставим потоку статус error
				threadData.Status = "error"
			}
		}

		// после успешного завершения работы функции мониторинга получаем значение интервала
		if albumPhotoMonitorParam != nil {
			interval = albumPhotoMonitorParam.Interval
		}

		// и включаем режим ожидания
		for i := 0; i < interval; i++ {
			// если статус потока error
			if threadData.Status == "error" {
				// то каждый раз обнуляем i, и тем самым вводим поток в вечное ожидание
				i = 0
			} else { // если нет, то ставим статус waiting и указываем оставшееся время ожидания
				threadData.Status = fmt.Sprintf("waiting for %d sec", interval-i)
			}
			time.Sleep(1 * time.Second)

			// периодически проверяем, был ли выставлен флаг остановки
			if threadData.ActionFlag == 1 {
				// если был, то меняем статус потока на "stopped"
				threadData.Status = "stopped"
				// и завершаем работу потока
				return
			}

			// если выставлен флаг рестарта
			if threadData.ActionFlag == 2 {
				// то обновляем статус потока
				threadData.Status = "alive"
				// и перезапускаем функцию
				threadData.ActionFlag = 0
				albumPhotoMonitoring(threadData, subject)
			}
		}
	}
}

func videoMonitoring(threadData *Thread, subject Subject) {

	// сообщаем пользователю о запуске модуля
	sender := threadData.Name
	message := "Started..."
	OutputMessage(sender, message)

	// создаем счетчик ошибок
	errorsCounter := 0

	// запускаем бесконечный цикл
	for true {
		// меняем статус потока
		threadData.Status = "alive"
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
				// и ставим потоку статус error
				threadData.Status = "error"
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
			// если статус потока error
			if threadData.Status == "error" {
				// то каждый раз обнуляем i, и тем самым вводим поток в вечное ожидание
				i = 0
			} else { // если нет, то ставим статус waiting и указываем оставшееся время ожидания
				threadData.Status = fmt.Sprintf("waiting for %d sec", interval-i)
			}
			time.Sleep(1 * time.Second)

			// периодически проверяем, был ли выставлен флаг остановки
			if threadData.ActionFlag == 1 {
				// если был, то меняем статус потока на "stopped"
				threadData.Status = "stopped"
				// и завершаем работу потока
				return
			}

			// если выставлен флаг рестарта
			if threadData.ActionFlag == 2 {
				// то обновляем статус потока
				threadData.Status = "alive"
				// и перезапускаем функцию
				threadData.ActionFlag = 0
				videoMonitoring(threadData, subject)
			}
		}
	}
}

func photoCommentMonitoring(threadData *Thread, subject Subject) {

	// сообщаем пользователю о запуске модуля
	sender := threadData.Name
	message := "Started..."
	OutputMessage(sender, message)

	// создаем счетчик ошибок
	errorsCounter := 0

	// запускаем бесконечный цикл
	for true {
		// меняем статус потока
		threadData.Status = "alive"
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
				// и ставим потоку статус error
				threadData.Status = "error"
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
			// если статус потока error
			if threadData.Status == "error" {
				// то каждый раз обнуляем i, и тем самым вводим поток в вечное ожидание
				i = 0
			} else { // если нет, то ставим статус waiting и указываем оставшееся время ожидания
				threadData.Status = fmt.Sprintf("waiting for %d sec", interval-i)
			}
			time.Sleep(1 * time.Second)

			// периодически проверяем, был ли выставлен флаг остановки
			if threadData.ActionFlag == 1 {
				// если был, то меняем статус потока на "stopped"
				threadData.Status = "stopped"
				// и завершаем работу потока
				return
			}

			// если выставлен флаг рестарта
			if threadData.ActionFlag == 2 {
				// то обновляем статус потока
				threadData.Status = "alive"
				// и перезапускаем функцию
				threadData.ActionFlag = 0
				photoCommentMonitoring(threadData, subject)
			}
		}
	}
}

func videoCommentMonitoring(threadData *Thread, subject Subject) {

	// сообщаем пользователю о запуске модуля
	sender := threadData.Name
	message := "Started..."
	OutputMessage(sender, message)

	// создаем счетчик ошибок
	errorsCounter := 0

	// запускаем бесконечный цикл
	for true {
		// меняем статус потока
		threadData.Status = "alive"
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
				// и ставим потоку статус error
				threadData.Status = "error"
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
			// если статус потока error
			if threadData.Status == "error" {
				// то каждый раз обнуляем i, и тем самым вводим поток в вечное ожидание
				i = 0
			} else { // если нет, то ставим статус waiting и указываем оставшееся время ожидания
				threadData.Status = fmt.Sprintf("waiting for %d sec", interval-i)
			}
			time.Sleep(1 * time.Second)

			// периодически проверяем, был ли выставлен флаг остановки
			if threadData.ActionFlag == 1 {
				// если был, то меняем статус потока на "stopped"
				threadData.Status = "stopped"
				// и завершаем работу потока
				return
			}

			// если выставлен флаг рестарта
			if threadData.ActionFlag == 2 {
				// то обновляем статус потока
				threadData.Status = "alive"
				// и перезапускаем функцию
				threadData.ActionFlag = 0
				videoCommentMonitoring(threadData, subject)
			}
		}
	}
}

func topicMonitoring(threadData *Thread, subject Subject) {

	// сообщаем пользователю о запуске модуля
	sender := threadData.Name
	message := "Started..."
	OutputMessage(sender, message)

	// создаем счетчик ошибок
	errorsCounter := 0

	// запускаем бесконечный цикл
	for true {
		// меняем статус потока
		threadData.Status = "alive"
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
				// и ставим потоку статус error
				threadData.Status = "error"
			}
		}

		// после успешного завершения работы функции мониторинга получаем значение интервала
		if topicMonitorParam != nil {
			interval = topicMonitorParam.Interval
		}

		// и включаем режим ожидания
		for i := 0; i < interval; i++ {
			// если статус потока error
			if threadData.Status == "error" {
				// то каждый раз обнуляем i, и тем самым вводим поток в вечное ожидание
				i = 0
			} else { // если нет, то ставим статус waiting и указываем оставшееся время ожидания
				threadData.Status = fmt.Sprintf("waiting for %d sec", interval-i)
			}
			time.Sleep(1 * time.Second)

			// периодически проверяем, был ли выставлен флаг остановки
			if threadData.ActionFlag == 1 {
				// если был, то меняем статус потока на "stopped"
				threadData.Status = "stopped"
				// и завершаем работу потока
				return
			}

			// если выставлен флаг рестарта
			if threadData.ActionFlag == 2 {
				// то обновляем статус потока
				threadData.Status = "alive"
				// и перезапускаем функцию
				threadData.ActionFlag = 0
				topicMonitoring(threadData, subject)
			}
		}
	}
}

func wallPostCommentMonitoring(threadData *Thread, subject Subject) {

	// сообщаем пользователю о запуске модуля
	sender := threadData.Name
	message := "Started..."
	OutputMessage(sender, message)

	// создаем счетчик ошибок
	errorsCounter := 0

	// запускаем бесконечный цикл
	for true {
		// меняем статус потока
		threadData.Status = "alive"
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
				// и ставим потоку статус error
				threadData.Status = "error"
			}
		}

		// после успешного завершения работы функции мониторинга получаем значение интервала
		if wallPostCommentMonitorParam != nil {
			interval = wallPostCommentMonitorParam.Interval
		}

		// и включаем режим ожидания
		for i := 0; i < interval; i++ {
			// если статус потока error
			if threadData.Status == "error" {
				// то каждый раз обнуляем i, и тем самым вводим поток в вечное ожидание
				i = 0
			} else { // если нет, то ставим статус waiting и указываем оставшееся время ожидания
				threadData.Status = fmt.Sprintf("waiting for %d sec", interval-i)
			}
			time.Sleep(1 * time.Second)

			// периодически проверяем, был ли выставлен флаг остановки
			if threadData.ActionFlag == 1 {
				// если был, то меняем статус потока на "stopped"
				threadData.Status = "stopped"
				// и завершаем работу потока
				return
			}

			// если выставлен флаг рестарта
			if threadData.ActionFlag == 2 {
				// то обновляем статус потока
				threadData.Status = "alive"
				// и перезапускаем функцию
				threadData.ActionFlag = 0
				wallPostCommentMonitoring(threadData, subject)
			}
		}
	}
}
