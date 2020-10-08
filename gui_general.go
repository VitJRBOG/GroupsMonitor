package main

import (
	"strings"

	"github.com/andlabs/ui"
)

// boxGeneral хранит данные о боксе с кнопками запуска, перезапуска и остановки модулей мониторинга
type boxGeneral struct {
	box *ui.Box
}

func (bg *boxGeneral) init() {
	bg.box = ui.NewHorizontalBox()
}

func (bg *boxGeneral) setBtnsBox(btnsBox generalButtonsBox) {
	bg.box.Append(btnsBox.box, false)
}

func (bg *boxGeneral) initFlexibleSpaceBox() {
	box := ui.NewHorizontalBox()
	bg.box.Append(box, true)
}

// generalButtonsBox хранит данные о кнопках запуска, перезапуска и остановки модулей мониторинга
type generalButtonsBox struct {
	box        *ui.Box
	btnStart   *ui.Button
	btnRestart *ui.Button
	btnStop    *ui.Button
}

func (gbb *generalButtonsBox) init() {
	gbb.box = ui.NewVerticalBox()
	gbb.box.SetPadded(true)
}

func (gbb *generalButtonsBox) initBtnStart(threads *[]*Thread) {
	gbb.btnStart = ui.NewButton("Запуск")

	gbb.btnStart.OnClicked(func(*ui.Button) {
		go StartThreads(threads)
		gbb.btnRestart.Enable()
		gbb.btnStop.Enable()
	})

	gbb.box.Append(gbb.btnStart, false)
}

func (gbb *generalButtonsBox) initBtnRestart(threads *[]*Thread) {
	gbb.btnRestart = ui.NewButton("Перезапуск")
	gbb.btnRestart.Disable()

	gbb.btnRestart.OnClicked(func(*ui.Button) {
		go RestartThreads(threads)
	})

	gbb.box.Append(gbb.btnRestart, false)
}

func (gbb *generalButtonsBox) initBtnStop(threads *[]*Thread) {
	gbb.btnStop = ui.NewButton("Остановка")
	gbb.btnStop.Disable()

	gbb.btnStop.OnClicked(func(*ui.Button) {
		go StopThreads(threads)
		gbb.btnRestart.Disable()
		gbb.btnStop.Disable()
	})

	gbb.box.Append(gbb.btnStop, false)
}

// makeGeneralBox собирает бокс с кнопками запуска, перезапуска и остановки модулей мониторинга
func makeGeneralBox(threads *[]*Thread) *ui.Box {

	var gbb generalButtonsBox

	gbb.init()
	gbb.initBtnStart(threads)
	gbb.initBtnRestart(threads)
	gbb.initBtnStop(threads)

	var bg boxGeneral
	bg.init()
	bg.initFlexibleSpaceBox()
	bg.setBtnsBox(gbb)
	bg.initFlexibleSpaceBox()

	return bg.box
}

// StartThreads запускает потоки, находящиеся в режиме ожидания после запуска программы
func StartThreads(threads *[]*Thread) {
	// сообщаем пользователю о начале операции запуска потоков
	sender := "Core"
	message := "Starting threads. Please stand by..."
	OutputMessage(sender, message)

	for i := 0; i < len(*threads); i++ {
		monitorName := strings.ReplaceAll((*threads)[i].Name, (*threads)[i].Subject.Name+": ", "")
		switch monitorName {
		case "посты на стене":
			if (*threads)[i].Status != "неактивен" {
				(*threads)[i].runWallPostMonitoring()
			}
		case "фото в альбомах":
			if (*threads)[i].Status != "неактивен" {
				(*threads)[i].runAlbumPhotoMonitoring()
			}
		case "видео в альбомах":
			if (*threads)[i].Status != "неактивен" {
				(*threads)[i].runVideoMonitoring()
			}
		case "комментарии под фото":
			if (*threads)[i].Status != "неактивен" {
				(*threads)[i].runPhotoCommentMonitoring()
			}
		case "комментарии под видео":
			if (*threads)[i].Status != "неактивен" {
				(*threads)[i].runVideoCommentMonitoring()
			}
		case "комментарии в обсуждениях":
			if (*threads)[i].Status != "неактивен" {
				(*threads)[i].runTopicMonitoring()
			}
		case "комментарии под постами":
			if (*threads)[i].Status != "неактивен" {
				(*threads)[i].runWallPostCommentMonitoring()
			}
		}
	}

	return
}

// RestartThreads перезапускает запущенные потоки
func RestartThreads(threads *[]*Thread) {

	// сообщаем пользователю о перезапуске алгоритмов мониторинга
	sender := "Core"
	message := "Restarting monitors. Please stand by..."
	OutputMessage(sender, message)

	// пробегаем по всем потокам и выставляем флаги на перезапуск потоков
	for _, thread := range *threads {
		if thread != nil {
			thread.ActionFlag = 2
		}
	}

	return
}

// StopThreads останавливает все потоки
func StopThreads(threads *[]*Thread) {

	// сообщаем пользователю о начале операции по остановке потоков
	sender := "Core"
	message := "Stopping threads..."
	OutputMessage(sender, message)

	// пробегаем по всем потокам и выставляем флаги на остановку
	for _, thread := range *threads {
		if thread != nil {
			thread.ActionFlag = 1
		}
	}

	return
}
