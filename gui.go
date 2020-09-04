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
	Box    *ui.Box // лишний параметр, нигде не участвует
}

func makeThreadControlBox(threads []*Thread) *ui.Box {
	// описываем коробку для управления потоками
	boxThreadControl := ui.NewHorizontalBox()
	boxThreadControl.SetPadded(true)

	// в этом списке будут хранится данные о коробках, ориентированных на конкретные субъекты (сообщества)
	var listSubjectBoxData []SubjectBoxData

	// описываем коробку для кнопок переключения между субъектами
	boxSubjectsSelection := ui.NewVerticalBox()
	groupSubject := ui.NewGroup("")

	// получаем список названий субъектов
	subjectsNames := getSubjectsNames()

	// проверяем, если в готовом списке названий субъектов пусто
	if len(subjectsNames) == 0 {
		// то отображаем сообщение об этом в окне программы
		labelNone := ui.NewLabel("Subjects not found")
		boxThreadControl.Append(labelNone, true)
	} else {
		// если список не пуст, перечисляем названия
		for _, subjectName := range subjectsNames {
			var subjectBoxData SubjectBoxData

			// для начала создаем кнопку для отображения коробки управления данного конкретного субъекта
			button := ui.NewButton(subjectName)
			// затем добавляем название субъекта и ссылку на кнопку в структуру данных
			subjectBoxData.Title = subjectName
			subjectBoxData.Button = button

			// и помещаем эту структуру в список
			listSubjectBoxData = append(listSubjectBoxData, subjectBoxData)
		}

		// перечисляем структуры из ранее собранного списка
		for i := 0; i < len(listSubjectBoxData); i++ {
			// присваиваем данные в переменные, т.к. в методе OnClicked полученные извне индексы не работают
			buttonTitle := listSubjectBoxData[i].Title
			button := listSubjectBoxData[i].Button

			// получаем получаем коробку, ориентированную на данный конкретный субъект
			boxSubject := makeSubjectBox(buttonTitle, threads)

			listSubjectBoxData[i].Box = boxThreadControl // лишняя строчка, нигде не участвует

			// привязываем к кнопке отображения коробки управления субъекта соответствующую процедуру
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

		// добавляем кнопки переключения между субъектами на коробку для этих кнопок
		for _, subjectBoxData := range listSubjectBoxData {
			boxSubjectsSelection.Append(subjectBoxData.Button, false)
		}

		// добавляем на коробку для управления потоками
		// коробку с кнопками для переключения между субъектами
		boxThreadControl.Append(boxSubjectsSelection, false)
		// и группу с коробками управления субъектами
		boxThreadControl.Append(groupSubject, true)
	}

	return boxThreadControl
}

func getSubjectsNames() []string {
	var subjectsNames []string

	// запрашиваем список субъектов из базы данных
	subjects, err := SelectDBSubjects()
	if err != nil {
		date := UnixTimeStampToDate(int(time.Now().Unix()))
		log.Fatal(fmt.Errorf("> [%v] WARNING! Error: %v", date, err))
	}
	// и добавляем их названия в список
	for _, subject := range subjects {
		subjectsNames = append(subjectsNames, subject.Name)
	}

	return subjectsNames
}

func makeSubjectBox(subjectName string, threads []*Thread) *ui.Box {
	// описываем коробку управления конкретным субъектом
	boxSubject := ui.NewVerticalBox()

	// создаем список с названиями модулей мониторинга (для названий кнопок управления потоками)
	monitorsNames := []string{"Wall post monitoring", "Album photo monitoring", "Video monitoring",
		"Photo comment monitoring", "Video comment monitoring", "Topic monitoring",
		"Wall post comment monitoring"}

	// в этом списке будут храниться кнопки для управления потоками
	var btnsMonitorControl []*ui.Button
	// а в этом - метки для отображения статусов соответствующих потоков
	var lblsMonitorControl []*ui.Label

	// перечисляем названия модулей мониторинга
	for _, monitorName := range monitorsNames {
		// описываем коробку для отображения кнопки управления и метки для статуса
		boxMonitorControl := ui.NewHorizontalBox()

		// описываем коробку для кнопки управления потоком
		boxBtnMonitorControl := ui.NewVerticalBox()
		// и описываем саму кнопку
		btnMonitorControl := ui.NewButton(monitorName)
		// затем добавляем на коробку для кнопки
		boxBtnMonitorControl.Append(btnMonitorControl, false)
		// по умолчанию делаем ее неактивной (делать активной будем потом)
		btnMonitorControl.Disable()
		// добавляем кнопку в список для кнопок
		btnsMonitorControl = append(btnsMonitorControl, btnMonitorControl)
		// а затем на коробку (сюда потом еще добавим метку для статуса)
		boxMonitorControl.Append(boxBtnMonitorControl, true)

		// описываем коробку для метки для статуса
		boxLblMonitorControl := ui.NewVerticalBox()
		// по умолчанию ставим статус stopped (менять статус будем потом)
		lblMonitorControl := ui.NewLabel("stopped")
		// затем добавляем на коробку для метки
		boxLblMonitorControl.Append(lblMonitorControl, false)
		// добавляем метку в списко для меток
		lblsMonitorControl = append(lblsMonitorControl, lblMonitorControl)
		// а затем на коробку, где ранее разместили кнопку управления потоком
		boxMonitorControl.Append(boxLblMonitorControl, true)

		// размещаем коробку с кнопкой и меткой на коробку управления потоками субъекта
		boxSubject.Append(boxMonitorControl, false)
	}

	// перечисляем описанные кнопки управления потоками
	for i, btnMonitorControl := range btnsMonitorControl {
		btnTitle := btnMonitorControl.Text()

		// перечисляем данные о созданных потоках
		for _, threadData := range threads {

			// проверяем, содержится ли название субъекта в названии потока
			if strings.Contains(strings.ToLower(threadData.Name), strings.ToLower(subjectName)) {

				// сравниваем название потока с конкатенацией названия субъекта и названия кнопки
				if strings.ToLower(threadData.Name) == strings.ToLower(subjectName+"'s "+btnTitle) {
					// привязываем к кнопке управления потоком соответствующую процедуру
					btnMonitorControl.OnClicked(func(*ui.Button) {
						if threadData.Status == "waiting" {
							threadData.Status = "alive"
						} else {
							if threadData.Status == "alive" {
								threadData.Status = "waiting"
							}
						}
					})
					// делаем кнопку активной, так как поток, которым она управляет, существует
					btnMonitorControl.Enable()

					// проверку статуса данного потока мониторинга запускаем в отдельном потоке
					go threadStatusChecking(lblsMonitorControl[i], threadData)

					break
				}
			}
		}
	}

	return boxSubject
}

func threadStatusChecking(lblMonitorControl *ui.Label, threadData *Thread) {
	// запускаем вечный цикл
	for true {
		// если текст в метке для отображения статуса не совпадает с названием статуса потока
		if lblMonitorControl.Text() != threadData.Status {
			// то меняем текст в метке на соответствующий названию статуса потока
			lblMonitorControl.SetText(threadData.Status)
		}
		// и ждем 1 секунду, затем повторяем цикл проверки
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
