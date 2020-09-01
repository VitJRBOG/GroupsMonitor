package main

import (
	"fmt"
	"log"
	"time"

	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"
)

// RunGui запускает собранный GUI
func RunGui() {
	ui.Main(initGui)
}

func initGui() {
	// запускаем функцию создания потоков с модулями проверки
	threads, err := MakeThreads()
	if err != nil {
		date := UnixTimeStampToDate(int(time.Now().Unix()))
		log.Fatal(fmt.Errorf("> [%v] WARNING! Error: %v", date, err))
	}

	// описываем главное окно GUI
	wndMain := ui.NewWindow("GroupsMonitor", 300, 150, true)
	wndMain.OnClosing(func(*ui.Window) bool {
		ui.Quit()
		return true
	})
	ui.OnShouldQuit(func() bool {
		wndMain.Destroy()
		return true
	})

	// описываем главную коробку для объектов интерфейса
	vBox := ui.NewVerticalBox()
	// и добавляем ее на главное окно
	wndMain.SetChild(vBox)

	// описываем основные кнопки программы
	btnStart := ui.NewButton("Start")
	btnRestart := ui.NewButton("Restart")
	btnStop := ui.NewButton("Stop")

	// и привязываем к каждой соответствующую процедуру
	btnStart.OnClicked(func(*ui.Button) {
		btnStartFunction(threads)
	})
	btnRestart.OnClicked(func(*ui.Button) {
		btnRestartFunction(threads)
	})
	btnStop.OnClicked(func(*ui.Button) {
		btnStopFunction(threads)
	})

	// затем добавляем эти кнопки в главную коробку интерфейса
	vBox.Append(btnStart, false)
	vBox.Append(btnRestart, false)
	vBox.Append(btnStop, false)

	// отображаем главное окно
	wndMain.Show()
}

func btnStartFunction(threads []*Thread) {
	StartThreads(threads)
}

func btnRestartFunction(threads []*Thread) {
	RestartThreads(threads)
}

func btnStopFunction(threads []*Thread) {
	StopThreads(threads)
}
