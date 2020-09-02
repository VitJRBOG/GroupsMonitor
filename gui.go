package main

import (
	"fmt"
	"log"
	"strings"
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
	boxThreadControl := makeThreadControlBox(threads)

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

// SubjectBoxData хранит данные для коробок с управленем потоками
type SubjectBoxData struct {
	Title  string
	Button *ui.Button
	Box    *ui.Box
}

func makeThreadControlBox(threads []*Thread) *ui.Box {
	boxThreadControl := ui.NewHorizontalBox()
	boxThreadControl.SetPadded(true)

	var listSubjectBoxData []SubjectBoxData

	boxSubjectsSelection := ui.NewVerticalBox()
	groupSubject := ui.NewGroup("")

	subjectsNames := getSubjectsNames()

	if len(subjectsNames) == 0 {
		labelNone := ui.NewLabel("Subjects not found")
		boxThreadControl.Append(labelNone, true)
	} else {
		for _, subjectName := range subjectsNames {
			var subjectBoxData SubjectBoxData

			button := ui.NewButton(subjectName)
			subjectBoxData.Title = subjectName
			subjectBoxData.Button = button

			listSubjectBoxData = append(listSubjectBoxData, subjectBoxData)
		}

		for i := 0; i < len(listSubjectBoxData); i++ {
			buttonTitle := listSubjectBoxData[i].Title
			button := listSubjectBoxData[i].Button
			boxSubject := makeSubjectBox(buttonTitle, threads)

			listSubjectBoxData[i].Box = boxThreadControl

			listSubjectBoxData[i].Button.OnClicked(func(*ui.Button) {
				groupSubject.SetChild(boxSubject)
				groupSubject.SetTitle(buttonTitle)

				for n := 0; n < len(listSubjectBoxData); n++ {
					if !(listSubjectBoxData[n].Button.Enabled()) {
						listSubjectBoxData[n].Button.Enable()
					}
				}

				button.Disable()
			})
		}

		for _, subjectBoxData := range listSubjectBoxData {
			boxSubjectsSelection.Append(subjectBoxData.Button, false)
		}

		boxThreadControl.Append(boxSubjectsSelection, false)
		boxThreadControl.Append(groupSubject, true)
	}

	return boxThreadControl
}

func getSubjectsNames() []string {
	var subjectsNames []string

	subjects, err := SelectDBSubjects()
	if err != nil {
		date := UnixTimeStampToDate(int(time.Now().Unix()))
		log.Fatal(fmt.Errorf("> [%v] WARNING! Error: %v", date, err))
	}

	for _, subject := range subjects {
		subjectsNames = append(subjectsNames, subject.Name)
	}

	return subjectsNames
}

func makeSubjectBox(subjectName string, threads []*Thread) *ui.Box {
	boxSubject := ui.NewVerticalBox()

	monitorsNames := []string{"Wall post monitoring", "Album photo monitoring", "Video monitoring",
		"Photo comment monitoring", "Video comment monitoring", "Topic monitoring",
		"Wall post comment monitoring"}

	var btnsMonitorControl []*ui.Button
	var lblsMonitorControl []*ui.Label

	for _, monitorName := range monitorsNames {
		boxMonitorControl := ui.NewHorizontalBox()

		boxBtnMonitorControl := ui.NewVerticalBox()
		btnMonitorControl := ui.NewButton(monitorName)
		boxBtnMonitorControl.Append(btnMonitorControl, false)
		btnMonitorControl.Disable()
		btnsMonitorControl = append(btnsMonitorControl, btnMonitorControl)
		boxMonitorControl.Append(boxBtnMonitorControl, true)

		boxLblMonitorControl := ui.NewVerticalBox()
		lblMonitorControl := ui.NewLabel("stopped")
		boxLblMonitorControl.Append(lblMonitorControl, false)
		lblsMonitorControl = append(lblsMonitorControl, lblMonitorControl)
		boxMonitorControl.Append(boxLblMonitorControl, true)

		boxSubject.Append(boxMonitorControl, false)
	}

	for i, btnMonitorControl := range btnsMonitorControl {
		btnTitle := btnMonitorControl.Text()

		for _, threadData := range threads {

			if strings.Contains(strings.ToLower(threadData.Name), strings.ToLower(subjectName)) {

				if strings.ToLower(threadData.Name) == strings.ToLower(subjectName+"'s "+btnTitle) {
					btnMonitorControl.OnClicked(func(*ui.Button) {
						if threadData.Status == "waiting" {
							threadData.Status = "alive"
						} else {
							if threadData.Status == "alive" {
								threadData.Status = "waiting"
							}
						}
					})
					btnMonitorControl.Enable()

					go threadStatusChecking(lblsMonitorControl[i], threadData)

					break
				}
			}
		}
	}

	return boxSubject
}

func threadStatusChecking(lblMonitorControl *ui.Label, threadData *Thread) {
	for true {
		if lblMonitorControl.Text() != threadData.Status {
			lblMonitorControl.SetText(threadData.Status)
		}
		time.Sleep(1 * time.Second)
	}
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
