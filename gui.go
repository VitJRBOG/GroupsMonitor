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
	// получаем список со ссылками на потоки
	threads := createThreads()

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

	// описываем группу для нижней панели окна
	groupBottom := ui.NewGroup("")
	groupBottom.SetMargined(true)

	// получаем коробку с основными кнопками
	boxGeneral := makeGeneralBox(threads)
	// и сразу добавляем ее в группу
	groupBottom.SetChild(boxGeneral)
	groupBottom.SetTitle("General")

	// получаем коробку для управления потоками
	boxThreadControl := makeThreadControlBox()

	// получаем коробку с настройками мониторинга
	boxSettings := makeSettingsBox()

	// получаем коробку с переключателями панелей
	boxSelection := makeSelectionBox(groupBottom, boxGeneral, boxThreadControl, boxSettings)

	// в конце добавляем на главную коробку коробку с кнопками-переключателями
	boxMain.Append(boxSelection, false)
	// а затем группу
	boxMain.Append(groupBottom, false)

	// отображаем главное окно
	wndMain.Show()
}

func makeSelectionBox(groupBottom *ui.Group, boxGeneral *ui.Box, boxThreadControl *ui.Box, boxSettings *ui.Box) *ui.Box {
	// описываем верхнюю коробку
	boxSelection := ui.NewHorizontalBox()

	// описываем левую верхнюю коробку
	boxSelectionLeft := ui.NewHorizontalBox()
	// и добавляем ее на верхнюю коробку
	boxSelection.Append(boxSelectionLeft, true)

	// описываем коробку для кнопок переключения
	boxSelectButton := ui.NewHorizontalBox()
	// и добавляем ее на верхнюю коробку
	boxSelection.Append(boxSelectButton, false)

	// описываем правую верхнюю коробку
	boxSelectionRight := ui.NewHorizontalBox()
	// и добавляем ее на верхнюю коробку
	boxSelection.Append(boxSelectionRight, true)

	// описываем кнопки для переключения между коробками
	btnGeneralBox := ui.NewButton("General")
	btnGeneralBox.Disable()
	btnThreadControlBox := ui.NewButton("Threads")
	btnSettings := ui.NewButton("Settings")

	// затем добавляем эти кнопки в коробку для кнопок переключения между коробками
	boxSelectButton.Append(btnGeneralBox, false)
	boxSelectButton.Append(btnThreadControlBox, false)
	boxSelectButton.Append(btnSettings, false)

	// затем привязываем к каждой кнопке-переключателе коробок соответствующую процедуру
	btnGeneralBox.OnClicked(func(*ui.Button) {
		groupBottom.SetChild(boxGeneral)
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

	return boxSelection
}

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
		StartThreads(threads)
		btnStart.Disable()
		btnRestart.Enable()
		btnStop.Enable()
	})
	btnRestart.OnClicked(func(*ui.Button) {
		RestartThreads(threads)
	})
	btnStop.OnClicked(func(*ui.Button) {
		StopThreads(threads)
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

func makeThreadControlBox() *ui.Box {
	boxThreadControl := ui.NewVerticalBox()

	// TODO

	return boxThreadControl
}

func makeSettingsBox() *ui.Box {
	boxSettings := ui.NewVerticalBox()

	// TODO

	return boxSettings
}

func createThreads() []*Thread {
	// запускаем функцию создания потоков с модулями проверки
	threads, err := MakeThreads()
	if err != nil {
		date := UnixTimeStampToDate(int(time.Now().Unix()))
		log.Fatal(fmt.Errorf("> [%v] WARNING! Error: %v", date, err))
	}

	return threads
}
