package main

import (
	"fmt"
	"time"

	"github.com/andlabs/ui"
)

func makeGeneralBox(threads []*Thread) *ui.Box {
	boxGeneral := ui.NewHorizontalBox()

	// описываем левую нижнюю коробку
	boxBottomLeft := ui.NewHorizontalBox()
	// и добавляем ее на основную коробку
	boxGeneral.Append(boxBottomLeft, true)

	// описываем коробку для основных кнопок
	buttonsBox := ui.NewVerticalBox()
	buttonsBox.SetPadded(true)
	// и добавляем ее на основную коробку
	boxGeneral.Append(buttonsBox, false)

	// описываем основные кнопки программы
	btnStart := ui.NewButton("Start")
	btnRestart := ui.NewButton("Restart")
	btnRestart.Disable()
	btnStop := ui.NewButton("Stop")
	btnStop.Disable()

	// и привязываем к каждой соответствующую процедуру
	btnStart.OnClicked(func(*ui.Button) {
		go StartThreads(threads)
		btnStart.Disable()
		btnRestart.Enable()
		btnStop.Enable()
	})
	btnRestart.OnClicked(func(*ui.Button) {
		go RestartThreads(threads)
	})
	btnStop.OnClicked(func(*ui.Button) {
		go StopThreads(threads)
		btnRestart.Disable()
		btnStop.Disable()
	})

	// затем добавляем эти кнопки в коробку для основных кнопок
	buttonsBox.Append(btnStart, false)
	buttonsBox.Append(btnRestart, false)
	buttonsBox.Append(btnStop, false)

	// описываем правую нижнюю коробку
	boxBottomRight := ui.NewHorizontalBox()
	// и добавляем ее на нижнюю коробку
	boxGeneral.Append(boxBottomRight, true)

	return boxGeneral
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
			thread.StopFlag = 2
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
			thread.StopFlag = 1
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
