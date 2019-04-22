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

	subjects, err := SelectDBSubjects()
	if err != nil {
		return threads, err
	}

	for _, subject := range subjects {
		wallPostMonitorParam, err := SelectDBWallPostMonitorParam(subject.ID)
		if err != nil {
			return threads, err
		}
		if wallPostMonitorParam.NeedMonitoring == 1 {
			var thread Thread
			thread.Name = fmt.Sprintf("%v's wall post monitoring", subject.Name)
			thread.Status = "alive"
			go wallPostMonitoring(&thread, subject, wallPostMonitorParam)
			threads = append(threads, &thread)
		}

		albumPhotoMonitorParam, err := SelectDBAlbumPhotoMonitorParam(subject.ID)
		if err != nil {
			return threads, err
		}
		if albumPhotoMonitorParam.NeedMonitoring == 1 {
			var thread Thread
			thread.Name = fmt.Sprintf("%v's album photo monitoring", subject.Name)
			thread.Status = "alive"
			go albumPhotoMonitoring(&thread, subject, albumPhotoMonitorParam)
			threads = append(threads, &thread)
		}
		// video_monitor
		// photo_comment_monitor
		// video_comment_monitor
		// topic_monitor
		// wall_post_comment_monitor
	}

	//
	// тут нужен поток с функцией проверки жизни остальных потоков
	//

	return threads, nil
}

func wallPostMonitoring(threadData *Thread, subject Subject, wallPostMonitorParam WallPostMonitorParam) {
	sender := threadData.Name
	message := "Started..."
	OutputMessage(sender, message)
	interval := wallPostMonitorParam.Interval
	for true {
		if err := WallPostMonitor(subject); err != nil {
			ErrorHandler(err)
		}
		for i := 0; i < interval; i++ {
			time.Sleep(1 * time.Second)
			if threadData.StopFlag == 1 {
				threadData.Status = "stopped"
				runtime.Goexit()
			}
		}
	}
}

func albumPhotoMonitoring(threadData *Thread, subject Subject, albumPhotoMonitorParam AlbumPhotoMonitorParam) {
	sender := threadData.Name
	message := "Started..."
	OutputMessage(sender, message)
	interval := albumPhotoMonitorParam.Interval
	for true {
		if err := AlbumPhotoMonitor(subject); err != nil {
			ErrorHandler(err)
		}
		for i := 0; i < interval; i++ {
			time.Sleep(1 * time.Second)
			if threadData.StopFlag == 1 {
				threadData.Status = "stopped"
				runtime.Goexit()
			}
		}
	}
}
