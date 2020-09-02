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
	wndMain := ui.NewWindow("GroupsMonitor", 255, 160, true)
	wndMain.SetMargined(true)
	wndMain.OnClosing(func(*ui.Window) bool {
		ui.Quit()
		return true
	})
	ui.OnShouldQuit(func() bool {
		wndMain.Destroy()
		return true
	})

	// описываем главную коробку для объектов интерфейса
	boxMain := ui.NewVerticalBox()
	boxMain.SetPadded(true)
	// и добавляем ее на главное окно
	wndMain.SetChild(boxMain)

	// описываем верхнюю коробку
	boxUpper := ui.NewHorizontalBox()
	// и добавляем ее на главную коробку
	boxMain.Append(boxUpper, false)

	// описываем левую верхнюю коробку
	boxUpperLeft := ui.NewHorizontalBox()
	// и добавляем ее на верхнюю коробку
	boxUpper.Append(boxUpperLeft, true)

	// описываем коробку для кнопок переключения
	boxSelectButton := ui.NewHorizontalBox()
	// и добавляем ее на верхнюю коробку
	boxUpper.Append(boxSelectButton, false)

	// описываем правую верхнюю коробку
	boxUpperRight := ui.NewHorizontalBox()
	// и добавляем ее на верхнюю коробку
	boxUpper.Append(boxUpperRight, true)

	// описываем кнопки для переключения между коробками
	btnGeneralBox := ui.NewButton("General")
	btnGeneralBox.Disable()
	btnThreadControlBox := ui.NewButton("Threads")
	btnSettings := ui.NewButton("Settings")

	// затем добавляем эти кнопки в коробку для кнопок переключения между коробками
	boxSelectButton.Append(btnGeneralBox, false)
	boxSelectButton.Append(btnThreadControlBox, false)
	boxSelectButton.Append(btnSettings, false)

	// описываем нижнюю коробку (для панели с основными кнопками)
	boxBottom := ui.NewHorizontalBox()
	// описываем группу для нижней панели окна
	groupBottom := ui.NewGroup("Main")
	groupBottom.SetMargined(true)
	// и добавляем коробку в группу
	groupBottom.SetChild(boxBottom)
	// затем добавляем группу на главную коробку
	boxMain.Append(groupBottom, false)

	// описываем коробку для управления потоками
	boxThreadControl := ui.NewVerticalBox()

	// описываем коробку для управления настройками монитора
	boxSettings := ui.NewVerticalBox()

	// затем привязываем к каждой кнопке-переключателе коробок соответствующую процедуру
	btnGeneralBox.OnClicked(func(*ui.Button) {
		groupBottom.SetChild(boxBottom)
		groupBottom.SetTitle("General")
		btnGeneralBox.Disable()
		if !(btnThreadControlBox.Enabled()) {
			btnThreadControlBox.Enable()
		}
		if !(btnSettings.Enabled()) {
			btnSettings.Enable()
		}
	})
	btnThreadControlBox.OnClicked(func(*ui.Button) {
		groupBottom.SetChild(boxThreadControl)
		btnThreadControlBox.Disable()
		groupBottom.SetTitle("Thread control")
		if !(btnGeneralBox.Enabled()) {
			btnGeneralBox.Enable()
		}
		if !(btnSettings.Enabled()) {
			btnSettings.Enable()
		}
	})
	btnSettings.OnClicked(func(*ui.Button) {
		groupBottom.SetChild(boxSettings)
		groupBottom.SetTitle("Settings")
		btnSettings.Disable()
		if !(btnGeneralBox.Enabled()) {
			btnGeneralBox.Enable()
		}
		if !(btnThreadControlBox.Enabled()) {
			btnThreadControlBox.Enable()
		}
	})

	// описываем левую нижнюю коробку
	boxBottomLeft := ui.NewHorizontalBox()
	// и добавляем ее на нижнюю коробку
	boxBottom.Append(boxBottomLeft, true)

	// описываем коробку для основных кнопок
	buttonsBox := ui.NewVerticalBox()
	buttonsBox.SetPadded(true)
	// и добавляем ее на левую коробку
	boxBottom.Append(buttonsBox, false)

	// описываем основные кнопки программы
	btnStart := ui.NewButton("Start")
	btnRestart := ui.NewButton("Restart")
	btnRestart.Disable()
	btnStop := ui.NewButton("Stop")
	btnStop.Disable()

	// и привязываем к каждой соответствующую процедуру
	btnStart.OnClicked(func(*ui.Button) {
		btnStartFunction(threads)
		btnStart.Disable()
		btnRestart.Enable()
		btnStop.Enable()
	})
	btnRestart.OnClicked(func(*ui.Button) {
		btnRestartFunction(threads)
	})
	btnStop.OnClicked(func(*ui.Button) {
		btnStopFunction(threads)
		btnRestart.Disable()
		btnStop.Disable()
	})

	// затем добавляем эти кнопки в коробку для основных кнопок
	buttonsBox.Append(btnStart, false)
	buttonsBox.Append(btnRestart, false)
	buttonsBox.Append(btnStop, false)

	// описываем правую нижняя коробку
	boxBottomRight := ui.NewHorizontalBox()
	// и добавляем ее на нижнюю коробку
	boxBottom.Append(boxBottomRight, true)

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
