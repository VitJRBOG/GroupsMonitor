package main

import (
	"fmt"
	"time"

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

func (gbb *generalButtonsBox) initBtnStart(threads []*Thread) {
	gbb.btnStart = ui.NewButton("Start")

	gbb.btnStart.OnClicked(func(*ui.Button) {
		go StartThreads(threads)
		gbb.btnStart.Disable()
		gbb.btnRestart.Enable()
		gbb.btnStop.Enable()
	})

	gbb.box.Append(gbb.btnStart, false)
}

func (gbb *generalButtonsBox) initBtnRestart(threads []*Thread) {
	gbb.btnRestart = ui.NewButton("Restart")
	gbb.btnRestart.Disable()

	gbb.btnRestart.OnClicked(func(*ui.Button) {
		go RestartThreads(threads)
	})

	gbb.box.Append(gbb.btnRestart, false)
}

func (gbb *generalButtonsBox) initBtnStop(threads []*Thread) {
	gbb.btnStop = ui.NewButton("Stop")
	gbb.btnStop.Disable()

	gbb.btnStop.OnClicked(func(*ui.Button) {
		go StopThreads(threads)
		gbb.btnRestart.Disable()
		gbb.btnStop.Disable()
	})

	gbb.box.Append(gbb.btnStop, false)
}

// makeGeneralBox собирает бокс с кнопками запуска, перезапуска и остановки модулей мониторинга
func makeGeneralBox(threads []*Thread) *ui.Box {

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
func StartThreads(threads []*Thread) {
	// сообщаем пользователю о начале операции запуска потоков
	sender := "Core"
	message := "Starting threads. Please stand by..."
	OutputMessage(sender, message)

	// пробегаем по всем потокам и снимаем статус ожидания
	for _, thread := range threads {
		if thread != nil {
			thread.Status = "alive"
		}
	}

	return
}

// RestartThreads перезапускает потоки, которые были запущены при старте программы
func RestartThreads(threads []*Thread) {

	// сообщаем пользователю о перезапуске алгоритмов мониторинга
	sender := "Core"
	message := "Restarting monitors. Please stand by..."
	OutputMessage(sender, message)

	// пробегаем по всем потокам и выставляем флаги на перезапуск потоков
	for _, thread := range threads {
		if thread != nil {
			thread.ActionFlag = 2
		}
	}

	return
}

// StopThreads останавливает все потоки
func StopThreads(threads []*Thread) {

	// сообщаем пользователю о начале операции по остановке потоков
	sender := "Core"
	message := "Stopping threads..."
	OutputMessage(sender, message)

	// пробегаем по всем потокам и выставляем флаги на остановку
	for _, thread := range threads {
		if thread != nil {
			thread.ActionFlag = 1
		}
	}

	repeats := 90

	// проверяем успешность остановки потоков
	for i := 0; i < repeats; i++ {

		// определяем количество живых потоков
		var alive int
		for _, thread := range threads {
			if thread != nil {
				alive++
			}
		}

		// если цикл повторился, то проверяем успешность завершения работы потоков
		if i > 0 {
			// если остались работающие потоки, то вызываем задержку
			if alive > 0 {
				interval := 1
				time.Sleep(time.Duration(interval) * time.Second)
			} else {
				// если работающих потоков не осталось, то сообщаем об этом пользователю и завершаем работу
				sender := "Core"
				message := "All threads is stopped. Quit..."
				OutputMessage(sender, message)
				return
			}
		}

		// если имеются работающие потоки, то проверяем результат их остановки
		if alive > 0 {
			for j, thread := range threads {

				// если поток имеет статус stopped, то обнуляем ссылку на него и сообщаем пользователю об успехе
				if thread != nil {
					if thread.Status == "stopped" {
						sender := thread.Name
						message := "OK! Monitoring is stopped!"
						OutputMessage(sender, message)
						threads[j] = nil
					}
				}
			}
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

	return
}
