package main

import (
	"errors"
	"fmt"
	"time"
)

// Thread - структура для хранения данных о потоке
type Thread struct {
	Name       string
	ActionFlag int // 0 - nothing, 1 - stopping, 2 - restarting
	Status     string
	Subject    Subject
	Errs       *Errors
}

func (thread *Thread) initWallPostMonitoring(subject Subject, errs *Errors) {
	thread.ActionFlag = 0
	thread.Name = fmt.Sprintf("%v: посты на стене", subject.Name)
	thread.Status = "остановлен"
	thread.Subject = subject
	thread.Errs = errs
}

func (thread *Thread) runWallPostMonitoring() {
	thread.ActionFlag = 0
	go monitoringCycle("wall_post_monitor", thread, thread.Subject)
}

func (thread *Thread) initAlbumPhotoMonitoring(subject Subject, errs *Errors) {
	thread.ActionFlag = 0
	thread.Name = fmt.Sprintf("%v: фото в альбомах", subject.Name)
	thread.Status = "остановлен"
	thread.Subject = subject
	thread.Errs = errs
}

func (thread *Thread) runAlbumPhotoMonitoring() {
	thread.ActionFlag = 0
	go monitoringCycle("album_photo_monitor", thread, thread.Subject)
}

func (thread *Thread) initVideoMonitoring(subject Subject, errs *Errors) {
	thread.ActionFlag = 0
	thread.Name = fmt.Sprintf("%v: видео в альбомах", subject.Name)
	thread.Status = "остановлен"
	thread.Subject = subject
	thread.Errs = errs
}

func (thread *Thread) runVideoMonitoring() {
	thread.ActionFlag = 0
	go monitoringCycle("video_monitor", thread, thread.Subject)
}

func (thread *Thread) initPhotoCommentMonitoring(subject Subject, errs *Errors) {
	thread.ActionFlag = 0
	thread.Name = fmt.Sprintf("%v: комментарии под фото", subject.Name)
	thread.Status = "остановлен"
	thread.Subject = subject
	thread.Errs = errs
}

func (thread *Thread) runPhotoCommentMonitoring() {
	thread.ActionFlag = 0
	go monitoringCycle("photo_comment_monitor", thread, thread.Subject)
}

func (thread *Thread) initVideoCommentMonitoring(subject Subject, errs *Errors) {
	thread.ActionFlag = 0
	thread.Name = fmt.Sprintf("%v: комментарии под видео", subject.Name)
	thread.Status = "остановлен"
	thread.Subject = subject
	thread.Errs = errs
}

func (thread *Thread) runVideoCommentMonitoring() {
	thread.ActionFlag = 0
	go monitoringCycle("video_comment_monitor", thread, thread.Subject)
}

func (thread *Thread) initTopicMonitoring(subject Subject, errs *Errors) {
	thread.ActionFlag = 0
	thread.Name = fmt.Sprintf("%v: комментарии в обсуждениях", subject.Name)
	thread.Status = "остановлен"
	thread.Subject = subject
	thread.Errs = errs
}

func (thread *Thread) runTopicMonitoring() {
	thread.ActionFlag = 0
	go monitoringCycle("topic_monitor", thread, thread.Subject)
}

func (thread *Thread) initWallPostCommentMonitoring(subject Subject, errs *Errors) {
	thread.ActionFlag = 0
	thread.Name = fmt.Sprintf("%v: комментарии под постами", subject.Name)
	thread.Status = "остановлен"
	thread.Subject = subject
	thread.Errs = errs
}

func (thread *Thread) runWallPostCommentMonitoring() {
	thread.ActionFlag = 0
	go monitoringCycle("wall_post_comment_monitor", thread, thread.Subject)
}

// InitThreads инициализирует потоки и заполняет данными о них список для них
func InitThreads() (*[]*Thread, *Errors, error) {
	var threads []*Thread
	var errs Errors

	var dbKit DataBaseKit
	subjects, err := dbKit.selectTableSubject()
	if err != nil {
		return nil, nil, err
	}

	for _, subject := range subjects {

		var wallPostMonitorParam WallPostMonitorParam
		err := wallPostMonitorParam.selectFromDBBySubjectID(subject.ID)
		if err != nil {
			return nil, nil, err
		}
		var threadWPM Thread
		threadWPM.initWallPostMonitoring(subject, &errs)
		if wallPostMonitorParam.NeedMonitoring == 0 {
			threadWPM.Status = "неактивен"
		}
		threads = append(threads, &threadWPM)

		var albumPhotoMonitorParam AlbumPhotoMonitorParam
		err = albumPhotoMonitorParam.selectFromDBBySubjectID(subject.ID)
		if err != nil {
			return nil, nil, err
		}
		var threadAPM Thread
		threadAPM.initAlbumPhotoMonitoring(subject, &errs)
		if albumPhotoMonitorParam.NeedMonitoring == 0 {
			threadAPM.Status = "неактивен"
		}
		threads = append(threads, &threadAPM)

		var videoMonitorParam VideoMonitorParam
		err = videoMonitorParam.selectFromDBBySubjectID(subject.ID)
		if err != nil {
			return nil, nil, err
		}
		var threadVM Thread
		threadVM.initVideoMonitoring(subject, &errs)
		if videoMonitorParam.NeedMonitoring == 0 {
			threadVM.Status = "неактивен"
		}
		threads = append(threads, &threadVM)

		var photoCommentMonitorParam PhotoCommentMonitorParam
		err = photoCommentMonitorParam.selectFromDBBySubjectID(subject.ID)
		if err != nil {
			return nil, nil, err
		}
		var threadPCM Thread
		threadPCM.initPhotoCommentMonitoring(subject, &errs)
		if photoCommentMonitorParam.NeedMonitoring == 0 {
			threadPCM.Status = "неактивен"
		}
		threads = append(threads, &threadPCM)

		var videoCommentMonitorParam VideoCommentMonitorParam
		err = videoCommentMonitorParam.selectFromDBBySubjectID(subject.ID)
		if err != nil {
			return nil, nil, err
		}
		var threadVCM Thread
		threadVCM.initVideoCommentMonitoring(subject, &errs)
		if videoCommentMonitorParam.NeedMonitoring == 0 {
			threadVCM.Status = "неактивен"
		}
		threads = append(threads, &threadVCM)

		var topicMonitorParam TopicMonitorParam
		err = topicMonitorParam.selectFromDBBySubjectID(subject.ID)
		if err != nil {
			return nil, nil, err
		}
		var threadTM Thread
		threadTM.initTopicMonitoring(subject, &errs)
		if topicMonitorParam.NeedMonitoring == 0 {
			threadTM.Status = "неактивен"
		}
		threads = append(threads, &threadTM)

		var wallPostCommentMonitorParam WallPostCommentMonitorParam
		err = wallPostCommentMonitorParam.selectFromDBBySubjectID(subject.ID)
		if err != nil {
			return nil, nil, err
		}
		var threadWPCM Thread
		threadWPCM.initWallPostCommentMonitoring(subject, &errs)
		if wallPostCommentMonitorParam.NeedMonitoring == 0 {
			threadWPCM.Status = "неактивен"
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

	return &threads, &errs, nil
}

// threadsStatusMonitoring ищет потоки, завершившие свою работу из-за ошибки
func threadsStatusMonitoring(threads *[]*Thread) {

	// перебираем структуры с данными о потоках
	for j, thread := range *threads {

		// если статус потока "error", то сообщаем об этом пользователю
		if thread.Status == "ошибка" {
			message := "WARNING! Thread is stopped with error!"
			OutputMessage(thread.Name, message)
			(*threads)[j] = nil
		}
	}

	// после завершения перебора включаем режим ожидания
	time.Sleep(10 * time.Second)
}

func monitoringCycle(moduleName string, threadData *Thread, subject Subject) {
	// сообщаем пользователю о запуске модуля
	sender := threadData.Name
	message := "Started..."
	OutputMessage(sender, message)

	// создаем счетчик ошибок
	errorsCounter := 0

	// запускаем бесконечный цикл
	for true {
		// меняем статус потока
		threadData.Status = "работает"
		// заранее присваиваем значение интервала
		interval := 20

		// запускаем функцию мониторинга
		monitInterval, err := startMonitoringModule(moduleName, subject)
		if err != nil {
			threadData.Errs.AddNewError(err.Error())
			// если функция вернула ошибку, то увеличиваем счетчик на 1
			errorsCounter++
			// если в результате счетчик не стал равен 4, то продолжаем
			if errorsCounter < 4 {
				// сообщаем пользователю об ошибке
				sender := fmt.Sprintf("%v -> Thread control", threadData.Name)
				// 20-секундный таймер умножаем на количество ошибок
				interval *= errorsCounter
				message := fmt.Sprintf("Error: %v. Time out for %ds", err, interval)
				OutputMessage(sender, message)
			} else {
				// если стал, то сообщаем об этом пользователю
				message := fmt.Sprintf("Error: %v. Thread is paused. Type \"restart\" for turn on again...", err)
				OutputMessage(sender, message)
				// и ставим потоку статус ошибка
				threadData.Status = "ошибка"
			}
		}

		// после успешного завершения работы функции мониторинга получаем значение интервала
		if monitInterval > 0 {
			interval = monitInterval
		}

		// и включаем режим ожидания
		for i := 0; i < interval; i++ {
			// если статус потока ошибка
			if threadData.Status == "ошибка" {
				// то каждый раз обнуляем i, и тем самым вводим поток в вечное ожидание
				i = 0
			} else { // если нет, то ставим статус waiting и указываем оставшееся время ожидания
				threadData.Status = fmt.Sprintf("ожидание еще %d сек.", interval-i)
			}
			time.Sleep(1 * time.Second)

			// периодически проверяем, был ли выставлен флаг остановки
			if threadData.ActionFlag == 1 {
				// если был, то меняем статус потока на "остановлен"
				threadData.Status = "остановлен"
				// и завершаем работу потока
				return
			}

			// если выставлен флаг рестарта
			if threadData.ActionFlag == 2 {
				// то обновляем статус потока
				threadData.Status = "работает"
				// и перезапускаем функцию
				threadData.ActionFlag = 0
				monitoringCycle(moduleName, threadData, subject)
			}
		}
	}
}

func startMonitoringModule(moduleName string, subject Subject) (int, error) {
	switch moduleName {
	case "wall_post_monitor":
		var monitInterval int
		monitorParams, err := WallPostMonitor(subject)
		if monitorParams != nil {
			monitInterval = monitorParams.Interval
		}
		return monitInterval, err
	case "album_photo_monitor":
		var monitInterval int
		monitorParams, err := AlbumPhotoMonitor(subject)
		if monitorParams != nil {
			monitInterval = monitorParams.Interval
		}
		return monitInterval, err
	case "video_monitor":
		var monitInterval int
		monitorParams, err := VideoMonitor(subject)
		if monitorParams != nil {
			monitInterval = monitorParams.Interval
		}
		return monitInterval, err
	case "photo_comment_monitor":
		var monitInterval int
		monitorParams, err := PhotoCommentMonitor(subject)
		if monitorParams != nil {
			monitInterval = monitorParams.Interval
		}
		return monitInterval, err
	case "video_comment_monitor":
		var monitInterval int
		monitorParams, err := VideoCommentMonitor(subject)
		if monitorParams != nil {
			monitInterval = monitorParams.Interval
		}
		return monitInterval, err
	case "topic_monitor":
		var monitInterval int
		monitorParams, err := TopicMonitor(subject)
		if monitorParams != nil {
			monitInterval = monitorParams.Interval
		}
		return monitInterval, err
	case "wall_post_comment_monitor":
		var monitInterval int
		monitorParams, err := WallPostCommentMonitor(subject)
		if monitorParams != nil {
			monitInterval = monitorParams.Interval
		}
		return monitInterval, err
	default:
		err := errors.New(moduleName + " - is unknown monit module name")
		panic(err.Error())
	}
}
