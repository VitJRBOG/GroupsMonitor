package main

import (
	"fmt"
	"runtime"
	"time"
)

// Thread - структура для хранения данных о потоке
type Thread struct {
	Name     string
	StopFlag int
	Status   string
}

// MakeThreads создает и запускает потоки
func MakeThreads() ([]*Thread, error) {
	var threads []*Thread

	// получаем из БД список субъектов
	subjects, err := SelectDBSubjects()
	if err != nil {
		return threads, err
	}

	// перебираем список субъектов и запускаем модули мониторинга в отдельных потоках
	for _, subject := range subjects {

		// получаем из БД параметры для модуля мониторинга постов на стене
		wallPostMonitorParam, err := SelectDBWallPostMonitorParam(subject.ID)
		if err != nil {
			return threads, err
		}

		// проверяем параметр и определяем, нужно ли запускать этот модуль
		if wallPostMonitorParam.NeedMonitoring == 1 {

			// создаем структуру с данными о потоке и наполняем ее данными
			var thread Thread
			thread.Name = fmt.Sprintf("%v's wall post monitoring", subject.Name)
			thread.Status = "alive"

			// запускаем поток
			go wallPostMonitoring(&thread, subject, wallPostMonitorParam)
			threads = append(threads, &thread)
		}

		// получаем из БД параметры для модуля мониторинга фотографий в альбомах
		albumPhotoMonitorParam, err := SelectDBAlbumPhotoMonitorParam(subject.ID)
		if err != nil {
			return threads, err
		}

		// проверяем параметр и определяем, нужно ли запускать этот модуль
		if albumPhotoMonitorParam.NeedMonitoring == 1 {

			// создаем структуру с данными о потоке и наполняем ее данными
			var thread Thread
			thread.Name = fmt.Sprintf("%v's album photo monitoring", subject.Name)
			thread.Status = "alive"

			// запускаем поток
			go albumPhotoMonitoring(&thread, subject, albumPhotoMonitorParam)
			threads = append(threads, &thread)
		}

		// получаем из БД параметры для модуля мониторинга видео в альбомах
		videoMonitorParam, err := SelectDBVideoMonitorParam(subject.ID)
		if err != nil {
			return threads, err
		}

		// проверяем параметр и определяем, нужно ли запускать этот модуль
		if videoMonitorParam.NeedMonitoring == 1 {

			// создаем структуру с данными о потоке и наполняем ее данными
			var thread Thread
			thread.Name = fmt.Sprintf("%v's video monitoring", subject.Name)
			thread.Status = "alive"

			// запускаем поток
			go videoMonitoring(&thread, subject, videoMonitorParam)
			threads = append(threads, &thread)
		}

		// получаем из БД параметры для модуля мониторинга комментариев под фотографиями
		photoCommentMonitorParam, err := SelectDBPhotoCommentMonitorParam(subject.ID)
		if err != nil {
			return threads, err
		}

		// проверяем параметр и определяем, нужно ли запускать этот модуль
		if photoCommentMonitorParam.NeedMonitoring == 1 {

			// создаем структуру с данными о потоке и наполняем ее данными
			var thread Thread
			thread.Name = fmt.Sprintf("%v's photo comment monitoring", subject.Name)
			thread.Status = "alive"

			// запускаем поток
			go photoCommentMonitoring(&thread, subject, photoCommentMonitorParam)
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
			thread.Status = "alive"

			// запускаем поток
			go videoCommentMonitoring(&thread, subject, videoCommentMonitorParam)
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
			thread.Status = "alive"

			// запускаем поток
			go topicMonitoring(&thread, subject, topicMonitorParam)
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
			thread.Status = "alive"

			// запускаем поток
			go wallPostCommentMonitoring(&thread, subject, wallPostCommentMonitorParam)
			threads = append(threads, &thread)
		}

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

func wallPostMonitoring(threadData *Thread, subject Subject, wallPostMonitorParam WallPostMonitorParam) error {

	// сообщаем пользователю о запуске модуля
	sender := threadData.Name
	message := "Started..."
	OutputMessage(sender, message)

	// получаем значение интервала между итерациями и запускаем бесконечный цикл
	interval := wallPostMonitorParam.Interval
	for true {

		// запускаем функцию мониторинга
		if err := WallPostMonitor(subject); err != nil {

			// если функция вернула ошибку, то сообщаем об этом пользователю
			message := fmt.Sprintf("Error: %v", err)
			OutputMessage(threadData.Name, message)

			// и меняем статус на "error"
			threadData.Status = "error"

			return err
		}

		// после успешного завершения работы функции мониторинга включаем режим ожидания
		for i := 0; i < interval; i++ {
			time.Sleep(1 * time.Second)

			// периодически проверяем, был ли выставлен флаг остановки
			if threadData.StopFlag == 1 {
				// если был, то меняем статус потока на "stopped" и завершаем его работу
				threadData.Status = "stopped"
				runtime.Goexit()
			}
		}
	}
	return nil
}

func albumPhotoMonitoring(threadData *Thread, subject Subject, albumPhotoMonitorParam AlbumPhotoMonitorParam) error {

	// сообщаем пользователю о запуске модуля
	sender := threadData.Name
	message := "Started..."
	OutputMessage(sender, message)

	// получаем значение интервала между итерациями и запускаем бесконечный цикл
	interval := albumPhotoMonitorParam.Interval
	for true {

		// запускаем функцию мониторинга
		if err := AlbumPhotoMonitor(subject); err != nil {

			// если функция вернула ошибку, то сообщаем об этом пользователю
			message := fmt.Sprintf("Error: %v", err)
			OutputMessage(threadData.Name, message)

			// и меняем статус на "error"
			threadData.Status = "error"
			return err
		}

		// после успешного завершения работы функции мониторинга включаем режим ожидания
		for i := 0; i < interval; i++ {
			time.Sleep(1 * time.Second)

			// периодически проверяем, был ли выставлен флаг остановки
			if threadData.StopFlag == 1 {
				// если был, то меняем статус потока на "stopped" и завершаем его работу
				threadData.Status = "stopped"
				runtime.Goexit()
			}
		}
	}
	return nil
}

func videoMonitoring(threadData *Thread, subject Subject, videoMonitorParam VideoMonitorParam) error {

	// сообщаем пользователю о запуске модуля
	sender := threadData.Name
	message := "Started..."
	OutputMessage(sender, message)

	// получаем значение интервала между итерациями и запускаем бесконечный цикл
	interval := videoMonitorParam.Interval
	for true {

		// запускаем функцию мониторинга
		if err := VideoMonitor(subject); err != nil {

			// если функция вернула ошибку, то сообщаем об этом пользователю
			message := fmt.Sprintf("Error: %v", err)
			OutputMessage(threadData.Name, message)

			// и меняем статус на "error"
			threadData.Status = "error"
			return err
		}

		// после успешного завершения работы функции мониторинга включаем режим ожидания
		for i := 0; i < interval; i++ {
			time.Sleep(1 * time.Second)

			// периодически проверяем, был ли выставлен флаг остановки
			if threadData.StopFlag == 1 {
				// если был, то меняем статус потока на "stopped" и завершаем его работу
				threadData.Status = "stopped"
				runtime.Goexit()
			}
		}
	}
	return nil
}

func photoCommentMonitoring(threadData *Thread, subject Subject,
	photoCommentMonitorParam PhotoCommentMonitorParam) error {

	// сообщаем пользователю о запуске модуля
	sender := threadData.Name
	message := "Started..."
	OutputMessage(sender, message)

	// получаем значение интервала между итерациями и запускаем бесконечный цикл
	interval := photoCommentMonitorParam.Interval
	for true {

		// запускаем функцию мониторинга
		if err := PhotoCommentMonitor(subject); err != nil {

			// если функция вернула ошибку, то сообщаем об этом пользователю
			message := fmt.Sprintf("Error: %v", err)
			OutputMessage(threadData.Name, message)

			// и меняем статус на "error"
			threadData.Status = "error"
			return err
		}

		// после успешного завершения работы функции мониторинга включаем режим ожидания
		for i := 0; i < interval; i++ {
			time.Sleep(1 * time.Second)

			// периодически проверяем, был ли выставлен флаг остановки
			if threadData.StopFlag == 1 {
				// если был, то меняем статус потока на "stopped" и завершаем его работу
				threadData.Status = "stopped"
				runtime.Goexit()
			}
		}
	}
	return nil
}

func videoCommentMonitoring(threadData *Thread, subject Subject,
	videoCommentMonitorParam VideoCommentMonitorParam) error {

	// сообщаем пользователю о запуске модуля
	sender := threadData.Name
	message := "Started..."
	OutputMessage(sender, message)

	// получаем значение интервала между итерациями и запускаем бесконечный цикл
	interval := videoCommentMonitorParam.Interval
	for true {

		// запускаем функцию мониторинга
		if err := VideoCommentMonitor(subject); err != nil {

			// если функция вернула ошибку, то сообщаем об этом пользователю
			message := fmt.Sprintf("Error: %v", err)
			OutputMessage(threadData.Name, message)

			// и меняем статус на "error"
			threadData.Status = "error"
			return err
		}

		// после успешного завершения работы функции мониторинга включаем режим ожидания
		for i := 0; i < interval; i++ {
			time.Sleep(1 * time.Second)

			// периодически проверяем, был ли выставлен флаг остановки
			if threadData.StopFlag == 1 {
				// если был, то меняем статус потока на "stopped" и завершаем его работу
				threadData.Status = "stopped"
				runtime.Goexit()
			}
		}
	}
	return nil
}

func topicMonitoring(threadData *Thread, subject Subject,
	topicMonitorParam TopicMonitorParam) error {

	// сообщаем пользователю о запуске модуля
	sender := threadData.Name
	message := "Started..."
	OutputMessage(sender, message)

	// получаем значение интервала между итерациями и запускаем бесконечный цикл
	interval := topicMonitorParam.Interval
	for true {

		// запускаем функцию мониторинга
		if err := TopicMonitor(subject); err != nil {

			// если функция вернула ошибку, то сообщаем об этом пользователю
			message := fmt.Sprintf("Error: %v", err)
			OutputMessage(threadData.Name, message)

			// и меняем статус на "error"
			threadData.Status = "error"
			return err
		}

		// после успешного завершения работы функции мониторинга включаем режим ожидания
		for i := 0; i < interval; i++ {
			time.Sleep(1 * time.Second)

			// периодически проверяем, был ли выставлен флаг остановки
			if threadData.StopFlag == 1 {
				// если был, то меняем статус потока на "stopped" и завершаем его работу
				threadData.Status = "stopped"
				runtime.Goexit()
			}
		}
	}
	return nil
}

func wallPostCommentMonitoring(threadData *Thread, subject Subject,
	wallPostCommentMonitorParam WallPostCommentMonitorParam) error {

	// сообщаем пользователю о запуске модуля
	sender := threadData.Name
	message := "Started..."
	OutputMessage(sender, message)

	// получаем значение интервала между итерациями и запускаем бесконечный цикл
	interval := wallPostCommentMonitorParam.Interval
	for true {

		// запускаем функцию мониторинга
		if err := WallPostCommentMonitor(subject); err != nil {

			// если функция вернула ошибку, то сообщаем об этом пользователю
			message := fmt.Sprintf("Error: %v", err)
			OutputMessage(threadData.Name, message)

			// и меняем статус на "error"
			threadData.Status = "error"
			return err
		}

		// после успешного завершения работы функции мониторинга включаем режим ожидания
		for i := 0; i < interval; i++ {
			time.Sleep(1 * time.Second)

			// периодически проверяем, был ли выставлен флаг остановки
			if threadData.StopFlag == 1 {
				// если был, то меняем статус потока на "stopped" и завершаем его работу
				threadData.Status = "stopped"
				runtime.Goexit()
			}
		}
	}
	return nil
}
