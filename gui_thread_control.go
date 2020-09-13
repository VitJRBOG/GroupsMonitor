package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/andlabs/ui"
)

// SubjectBoxData хранит данные для коробок с управленем потоками
type SubjectBoxData struct {
	Title  string
	Button *ui.Button
	Box    *ui.Box
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

			listSubjectBoxData[i].Box = boxSubject

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

		// для выравнивания расположения кнопок относительно коробки управления потоками
		// создаем группу для кнопок переключения между субъектами
		groupSubjectsSelection := ui.NewGroup("")
		// и добавляем на нее коробку с кнопками
		groupSubjectsSelection.SetChild(boxSubjectsSelection)

		// добавляем на коробку для управления потоками
		// группу с кнопками для переключения между субъектами
		boxThreadControl.Append(groupSubjectsSelection, false)
		// и группу с коробками управления субъектами
		boxThreadControl.Append(groupSubject, true)

		// по умолчанию отображаем коробку первого в списке субъекта
		groupSubject.SetChild(listSubjectBoxData[0].Box)
		groupSubject.SetTitle(listSubjectBoxData[0].Title)
		// и делаем неактивной кнопку отображения коробки этого субъекта
		listSubjectBoxData[0].Button.Disable()
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
		// по умолчанию ставим статус inactive (менять статус будем потом)
		lblMonitorControl := ui.NewLabel("inactive")
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
