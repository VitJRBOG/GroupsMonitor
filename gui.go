package main

import (
	"fmt"
	"log"
	"strconv"
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

// GeneralBoxesData хранит коробки для основных установок
type GeneralBoxesData struct {
	AccessTokens *ui.Box
	Subjects     *ui.Box
}

// GroupsSettingsData хранит группы для установок
type GroupsSettingsData struct {
	Primary    *ui.Group
	General    *ui.Group
	Additional *ui.Group
}

func makeSettingsBox() *ui.Box {
	boxSettings := ui.NewHorizontalBox()

	// описываем три группы для отображения установок:
	var groupsSettingsData GroupsSettingsData
	// первичная
	groupPrimarySettings := ui.NewGroup("")
	groupsSettingsData.Primary = groupPrimarySettings
	// общие
	groupGeneralSettings := ui.NewGroup("")
	groupsSettingsData.General = groupGeneralSettings
	// дополнительная
	groupAdditionalSettings := ui.NewGroup("")
	groupsSettingsData.Additional = groupAdditionalSettings

	var generalBoxesData GeneralBoxesData

	// получаем коробку для установок токенов доступа
	boxAccessTokensSettings := makeAccessTokensSettingsBox()
	generalBoxesData.AccessTokens = boxAccessTokensSettings
	// по умолчанию отображаем ее в группе общих настроек
	groupGeneralSettings.SetTitle("Access tokens")
	groupGeneralSettings.SetChild(boxAccessTokensSettings)
	groupAdditionalSettings.SetChild(ui.NewLabel("Nothing to show here..."))

	// получаем коробку для установок субъектов
	boxSubjectsSettings := makeSubjectsSettingsBox(groupsSettingsData)
	generalBoxesData.Subjects = boxSubjectsSettings

	// получаем коробку для первичных установок
	boxPrimarySettings := makePrimarySettingsBox(generalBoxesData, groupsSettingsData)
	// и добавляем ее в соответствующую группу
	groupPrimarySettings.SetChild(boxPrimarySettings)

	// добавляем группы на коробку с настройками
	boxSettings.Append(groupPrimarySettings, false)
	boxSettings.Append(groupGeneralSettings, false)
	boxSettings.Append(groupAdditionalSettings, false)

	return boxSettings
}

func makeAccessTokensSettingsBox() *ui.Box {
	// описываем коробку для установок токенов доступа
	boxAccessTokensSettings := ui.NewVerticalBox()

	// запрашиваем список токенов доступа из базы данных
	accessTokens, err := SelectDBAccessTokens()
	if err != nil {
		date := UnixTimeStampToDate(int(time.Now().Unix()))
		log.Fatal(fmt.Errorf("> [%v] WARNING! Error: %v", date, err))
	}

	// перечисляем токены доступа
	for i := 0; i < len(accessTokens); i++ {
		accessTokenData := accessTokens[i]
		// описываем коробку для отображения названия токена доступа и кнопки вызова настроек
		boxAccessTokenSettings := ui.NewHorizontalBox()

		// описываем метку для отображения названия токена доступа
		lblAccessTokenName := ui.NewLabel(accessTokenData.Name)
		// и добавляем ее в коробку (потом еще добавим кнопку)
		boxAccessTokenSettings.Append(lblAccessTokenName, true)

		// описываем кнопку для вызова настроек соответствующего токена доступа
		btnAccessTokenSettings := ui.NewButton("Set...")
		// и добавляем ее в коробку (сюда ранее добавили кнопку)
		boxAccessTokenSettings.Append(btnAccessTokenSettings, true)

		// привязываем кнопку к процедуре отображения окна с параметрами
		btnAccessTokenSettings.OnClicked(func(*ui.Button) {
			showAccessTokenSettingWindow(accessTokenData.ID)
		})

		// размещаем коробку с меткой и кнопкой на коробке для установок токенов доступа
		boxAccessTokensSettings.Append(boxAccessTokenSettings, false)
	}

	return boxAccessTokensSettings
}

func showAccessTokenSettingWindow(IDAccessToken int) {
	// описываем окно для отображения установок токена доступа
	wndAccessTokenSettings := ui.NewWindow("", 300, 100, true)
	wndAccessTokenSettings.OnClosing(func(*ui.Window) bool {
		wndAccessTokenSettings.Disable()
		return true
	})
	wndAccessTokenSettings.SetMargined(true)
	boxWndMain := ui.NewVerticalBox()
	boxWndMain.SetPadded(true)
	wndAccessTokenSettings.SetChild(boxWndMain)

	// запрашиваем список токенов доступа из базы данных
	accessTokens, err := SelectDBAccessTokens()
	if err != nil {
		date := UnixTimeStampToDate(int(time.Now().Unix()))
		log.Fatal(fmt.Errorf("> [%v] WARNING! Error: %v", date, err))
	}

	// перечисляем токены доступа
	for _, accessToken := range accessTokens {
		// и ищем токен с подходящим идентификатором
		if accessToken.ID == IDAccessToken {
			// устанавливаем заголовок окна в соответствии с названием токена доступа
			wndAccessTokenSettings.SetTitle("Settings of " + accessToken.Name + "'s access token")

			boxWndAT := ui.NewVerticalBox()

			// описываем коробку с меткой и полем для названия токена доступа
			boxWndATName := ui.NewHorizontalBox()
			boxWndATName.SetPadded(true)
			lblWndATName := ui.NewLabel("Name")
			boxWndATName.Append(lblWndATName, false)
			entryWndATName := ui.NewEntry()
			entryWndATName.SetText(accessToken.Name)
			boxWndATName.Append(entryWndATName, true)

			// описываем коробку с меткой и полем для значения токена доступа
			boxWndATValue := ui.NewHorizontalBox()
			boxWndATValue.SetPadded(true)
			lblWndATValue := ui.NewLabel("Value")
			boxWndATValue.Append(lblWndATValue, false)
			entryWndATValue := ui.NewEntry()
			entryWndATValue.SetText(accessToken.Value)
			boxWndATValue.Append(entryWndATValue, true)

			// описываем группу, в которой будут размещены элементы
			groupWndAT := ui.NewGroup("")
			groupWndAT.SetMargined(true)
			boxWndAT.Append(boxWndATName, false)
			boxWndAT.Append(boxWndATValue, false)
			groupWndAT.SetChild(boxWndAT)

			// добавляем группу в основную коробку окна
			boxWndMain.Append(groupWndAT, false)

			// описываем коробку для кнопок
			boxWndATBtns := ui.NewHorizontalBox()
			boxWndATBtns.SetPadded(true)
			// и несколько коробок для выравнивания кнопок
			btnWndATBtnsLeft := ui.NewHorizontalBox()
			btnWndATBtnsCenter := ui.NewHorizontalBox()
			btnWndATBtnsRight := ui.NewHorizontalBox()
			btnWndATBtnsRight.SetPadded(true)
			// а затем сами кнопки
			btnwndATCancel := ui.NewButton("Cancel")
			btnWndATBtnsRight.Append(btnwndATCancel, false)
			btnWndATApplyChanges := ui.NewButton("Apply")
			btnWndATBtnsRight.Append(btnWndATApplyChanges, false)
			// и добавляем их в коробку для кнопок
			boxWndATBtns.Append(btnWndATBtnsLeft, false)
			boxWndATBtns.Append(btnWndATBtnsCenter, false)
			boxWndATBtns.Append(btnWndATBtnsRight, false)

			btnwndATCancel.OnClicked(func(*ui.Button) {
				// TODO: как-нибудь надо закрывать окно
			})
			// привязываем кнопки к соответствующим процедурам
			btnWndATApplyChanges.OnClicked(func(*ui.Button) {
				var updatedAccessToken AccessToken
				updatedAccessToken.ID = accessToken.ID
				updatedAccessToken.Name = entryWndATName.Text()
				updatedAccessToken.Value = entryWndATValue.Text()

				err := UpdateDBAccessToken(updatedAccessToken)
				if err != nil {
					date := UnixTimeStampToDate(int(time.Now().Unix()))
					log.Fatal(fmt.Errorf("> [%v] WARNING! Error: %v", date, err))
				}

				// TODO: как-нибудь надо закрывать окно
			})

			// добавляем коробку с кнопками на основную коробку окна
			boxWndMain.Append(boxWndATBtns, true)
			// затем еще одну коробку, для выравнивания расположения кнопок при растягивании окна
			boxWndATBottom := ui.NewHorizontalBox()
			boxWndMain.Append(boxWndATBottom, false)
			break
		}
	}

	wndAccessTokenSettings.Show()
}

func makeSubjectsSettingsBox(groupsSettingsData GroupsSettingsData) *ui.Box {
	// описываем коробку для установок субъектов
	boxSubjectsSettings := ui.NewVerticalBox()

	// запрашиваем список субъектов из базы данных
	subjects, err := SelectDBSubjects()
	if err != nil {
		date := UnixTimeStampToDate(int(time.Now().Unix()))
		log.Fatal(fmt.Errorf("> [%v] WARNING! Error: %v", date, err))
	}

	// в этом списке будут храниться ссылки на кнопки для отображения доп. настроек
	var listBtnsSubjectSettings []*ui.Button

	// перечисляем субъекты
	for _, subjectData := range subjects {
		// описываем кнопку для отображения доп. настроек соответствующего субъекта
		btnSubjectSettings := ui.NewButton(subjectData.Name)
		// и добавляем ее в коробку
		boxSubjectsSettings.Append(btnSubjectSettings, false)
		// добавляем кнопку в список
		listBtnsSubjectSettings = append(listBtnsSubjectSettings, btnSubjectSettings)
	}

	// перечисляем кнопки для отображения доп. настроек
	for i := 0; i < len(listBtnsSubjectSettings); i++ {
		btnSubjectSettings := listBtnsSubjectSettings[i]
		// еще раз перечисляем субъекты
		for _, subjectData := range subjects {
			// если название субъекта совпало с названием кнопки
			if subjectData.Name == btnSubjectSettings.Text() {
				// то получаем коробку для отображения кнопок для вызова доп. настроек
				boxSubjectAdditionalSettingsBox := makeSubjectAdditionalSettingsBox(subjectData)
				// и привязываем кнопку к процедуре отображения соответствующих доп. настроек
				btnSubjectSettings.OnClicked(func(*ui.Button) {
					groupsSettingsData.Additional.SetChild(boxSubjectAdditionalSettingsBox)
					groupsSettingsData.Additional.SetTitle(subjectData.Name)

					for n := 0; n < len(listBtnsSubjectSettings); n++ {
						if !(listBtnsSubjectSettings[n].Enabled()) {
							listBtnsSubjectSettings[n].Enable()
						}
					}

					btnSubjectSettings.Disable()
				})

				break
			}
		}
	}

	return boxSubjectsSettings
}

func makeSubjectAdditionalSettingsBox(subjectData Subject) *ui.Box {
	boxSubjectAdditionalSettingsBox := ui.NewVerticalBox()

	// создаем список с названиями кнопок для вызова окна доп. с установками
	btnsNames := []string{"General", "Wall post monitor", "Album photo monitor", "Video monitor",
		"Photo comment monitor", "Video comment monitor", "Topic monitor",
		"Wall post comment monitor"}

	// перечисляем названия кнопок
	for i := 0; i < len(btnsNames); i++ {
		btnName := btnsNames[i]
		// описываем коробку для отображения метки с названием доп. установок и кнопкой
		boxSettingsSection := ui.NewHorizontalBox()

		// описываем коробку для метки с названием доп. установок
		boxLblSettingsSection := ui.NewVerticalBox()
		lblSettingsSection := ui.NewLabel(btnName)
		boxLblSettingsSection.Append(lblSettingsSection, false)
		boxSettingsSection.Append(boxLblSettingsSection, true)

		// описываем коробку для кнопки вызова окна с доп. установками
		boxBtnSettingsSection := ui.NewVerticalBox()
		btnSettingsSection := ui.NewButton("Set...")
		boxBtnSettingsSection.Append(btnSettingsSection, false)
		boxSettingsSection.Append(boxBtnSettingsSection, true)

		// привязываем к кнопке отображения окна с доп. установками соответствующую процедуру
		switch btnName {
		case "General":
			btnSettingsSection.OnClicked(func(*ui.Button) {
				showSubjectGeneralSettingWindow(subjectData.ID, btnName)
			})
		case "Wall post monitor":
			btnSettingsSection.OnClicked(func(*ui.Button) {
				showSubjectWallPostSettingWindow(subjectData.ID, subjectData.Name, btnName)
			})
			// case "Album photo monitor":
			// case "Video monitor":
			// case "Photo comment monitor":
			// case "Video comment monitor":
			// case "Topic monitor":
			// case "Wall post comment monitor":
		}

		boxSubjectAdditionalSettingsBox.Append(boxSettingsSection, false)
	}

	return boxSubjectAdditionalSettingsBox
}

func showSubjectGeneralSettingWindow(IDSubject int, btnName string) {
	// описываем окно для отображения общих установок субъекта
	wndSubjectGeneralSettings := ui.NewWindow("", 300, 100, true)
	wndSubjectGeneralSettings.OnClosing(func(*ui.Window) bool {
		wndSubjectGeneralSettings.Disable()
		return true
	})
	wndSubjectGeneralSettings.SetMargined(true)
	boxWndMain := ui.NewVerticalBox()
	boxWndMain.SetPadded(true)
	wndSubjectGeneralSettings.SetChild(boxWndMain)

	// запрашиваем список субъектов из базы данных
	subjects, err := SelectDBSubjects()
	if err != nil {
		date := UnixTimeStampToDate(int(time.Now().Unix()))
		log.Fatal(fmt.Errorf("> [%v] WARNING! Error: %v", date, err))
	}

	// перечисляем субъекты
	for _, subject := range subjects {
		// ищем субъект с подходящим идентификатором
		if subject.ID == IDSubject {
			// устанавливаем заголовок окна в соответствии с названием субъекта и назначением установок
			wndTitle := fmt.Sprintf("%v settings for %v", btnName, subject.Name)
			wndSubjectGeneralSettings.SetTitle(wndTitle)

			boxWndS := ui.NewVerticalBox()

			// описываем коробку с меткой и полем для названия субъекта
			boxWndSName := ui.NewHorizontalBox()
			boxWndSName.SetPadded(true)
			lblWndSName := ui.NewLabel("Name")
			boxWndSName.Append(lblWndSName, true)
			entryWndSName := ui.NewEntry()
			entryWndSName.SetText(subject.Name)
			boxWndSName.Append(entryWndSName, true)

			// описываем коробку с меткой и полем для идентификатора субъекта в базе ВК
			boxWndSSubjectID := ui.NewHorizontalBox()
			boxWndSSubjectID.SetPadded(true)
			lblWndSSubjectID := ui.NewLabel("Subject ID")
			boxWndSSubjectID.Append(lblWndSSubjectID, true)
			entryWndSSubjectID := ui.NewEntry()
			entryWndSSubjectID.SetText(strconv.Itoa(subject.SubjectID))
			boxWndSSubjectID.Append(entryWndSSubjectID, true)

			// описываем группу, в которой будут размещены элементы
			groupWndS := ui.NewGroup("")
			groupWndS.SetMargined(true)
			boxWndS.Append(boxWndSName, false)
			boxWndS.Append(boxWndSSubjectID, false)
			groupWndS.SetChild(boxWndS)

			// добавляем группу в основную коробку окна
			boxWndMain.Append(groupWndS, false)

			// описываем коробку для кнопок
			boxWndSBtns := ui.NewHorizontalBox()
			boxWndSBtns.SetPadded(true)
			// и несколько коробок для выравнивания кнопок
			btnWndSBtnsLeft := ui.NewHorizontalBox()
			btnWndSBtnsCenter := ui.NewHorizontalBox()
			btnWndSBtnsRight := ui.NewHorizontalBox()
			btnWndSBtnsRight.SetPadded(true)
			// а затем сами кнопки
			btnWndSCancel := ui.NewButton("Cancel")
			btnWndSBtnsRight.Append(btnWndSCancel, false)
			btnWndSApplyChanges := ui.NewButton("Apply")
			btnWndSBtnsRight.Append(btnWndSApplyChanges, false)
			// и добавляем их в коробку для кнопок
			boxWndSBtns.Append(btnWndSBtnsLeft, false)
			boxWndSBtns.Append(btnWndSBtnsCenter, false)
			boxWndSBtns.Append(btnWndSBtnsRight, false)

			// привязываем к кнопкам соответствующие процедуры
			btnWndSCancel.OnClicked(func(*ui.Button) {
				// TODO: как-нибудь надо закрывать окно
			})
			// привязываем кнопки к соответствующим процедурам
			btnWndSApplyChanges.OnClicked(func(*ui.Button) {
				var updatedSubject Subject
				updatedSubject.ID = subject.ID
				updatedSubject.SubjectID, err = strconv.Atoi(entryWndSSubjectID.Text())
				if err != nil {
					date := UnixTimeStampToDate(int(time.Now().Unix()))
					log.Fatal(fmt.Errorf("> [%v] WARNING! Error: %v", date, err))
				}
				updatedSubject.Name = entryWndSName.Text()
				updatedSubject.BackupWikipage = subject.BackupWikipage
				updatedSubject.LastBackup = subject.LastBackup

				err := UpdateDBSubject(updatedSubject)
				if err != nil {
					date := UnixTimeStampToDate(int(time.Now().Unix()))
					log.Fatal(fmt.Errorf("> [%v] WARNING! Error: %v", date, err))
				}

				// TODO: как-нибудь надо закрывать окно
			})

			// добавляем коробку с кнопками на основную коробку окна
			boxWndMain.Append(boxWndSBtns, true)
			// затем еще одну коробку, для выравнивания расположения кнопок при растягивании окна
			boxWndSBottom := ui.NewHorizontalBox()
			boxWndMain.Append(boxWndSBottom, false)

			break
		}
	}

	wndSubjectGeneralSettings.Show()
}

func showSubjectWallPostSettingWindow(IDSubject int, nameSubject, btnName string) {
	// описываем окно для отображения общих установок субъекта
	wndSubjectWallPostSettings := ui.NewWindow("", 300, 100, true)
	wndSubjectWallPostSettings.OnClosing(func(*ui.Window) bool {
		wndSubjectWallPostSettings.Disable()
		return true
	})
	wndSubjectWallPostSettings.SetMargined(true)
	boxWndMain := ui.NewVerticalBox()
	boxWndMain.SetPadded(true)
	wndSubjectWallPostSettings.SetChild(boxWndMain)

	// запрашиваем параметры мониторинга из базы данных
	wallPostMonitorParam, err := SelectDBWallPostMonitorParam(IDSubject)
	if err != nil {
		date := UnixTimeStampToDate(int(time.Now().Unix()))
		log.Fatal(fmt.Errorf("> [%v] WARNING! Error: %v", date, err))
	}

	// устанавливаем заголовок окна в соответствии с названием субъекта и назначением установок
	wndTitle := fmt.Sprintf("%v settings for %v", btnName, nameSubject)
	wndSubjectWallPostSettings.SetTitle(wndTitle)

	boxWndWP := ui.NewVerticalBox()

	// описываем коробку с меткой и чекбоксом для флага необходимости активировать модуль мониторинга
	boxWndWPMonitoring := ui.NewHorizontalBox()
	boxWndWPMonitoring.SetPadded(true)
	lblWndWPMonitoring := ui.NewLabel("Need monitoring")
	boxWndWPMonitoring.Append(lblWndWPMonitoring, true)
	cboxWndWPNeedMonitoring := ui.NewCheckbox("")
	if wallPostMonitorParam.NeedMonitoring == 1 {
		cboxWndWPNeedMonitoring.SetChecked(true)
	} else {
		cboxWndWPNeedMonitoring.SetChecked(false)
	}
	boxWndWPMonitoring.Append(cboxWndWPNeedMonitoring, true)

	// описываем коробку с меткой и спинбоксом для интервала между запусками функции мониторинга
	boxWndWPInterval := ui.NewHorizontalBox()
	boxWndWPInterval.SetPadded(true)
	lblWndWPInterval := ui.NewLabel("Interval")
	boxWndWPInterval.Append(lblWndWPInterval, true)
	sboxWndWPInterval := ui.NewSpinbox(5, 21600)
	sboxWndWPInterval.SetValue(wallPostMonitorParam.Interval)
	boxWndWPInterval.Append(sboxWndWPInterval, true)

	// описываем коробку с меткой и полем для идентификатора получателя сообщений
	boxWndWPSendTo := ui.NewHorizontalBox()
	boxWndWPSendTo.SetPadded(true)
	lblWndWPSendTo := ui.NewLabel("Send to")
	boxWndWPSendTo.Append(lblWndWPSendTo, true)
	entryWndWPSendTo := ui.NewEntry()
	entryWndWPSendTo.SetText(strconv.Itoa(wallPostMonitorParam.SendTo))
	boxWndWPSendTo.Append(entryWndWPSendTo, true)

	// описываем коробку с меткой и выпадающим списком для названия фильтра постов
	boxWndWPFilter := ui.NewHorizontalBox()
	boxWndWPFilter.SetPadded(true)
	lblWndWPFilter := ui.NewLabel("Filter")
	boxWndWPFilter.Append(lblWndWPFilter, true)
	comboboxWndWPFilter := ui.NewCombobox()
	listPostsFilters := []string{"all", "others", "owner", "suggests"}
	var slctd int
	for i, postFilter := range listPostsFilters {
		comboboxWndWPFilter.Append(postFilter)
		if wallPostMonitorParam.Filter == postFilter {
			slctd = i
		}
	}
	comboboxWndWPFilter.SetSelected(slctd)
	boxWndWPFilter.Append(comboboxWndWPFilter, true)

	// описываем коробку с меткой и спинбоксом для количества проверяемых постов
	boxWndWPPostsCount := ui.NewHorizontalBox()
	boxWndWPPostsCount.SetPadded(true)
	lblWndWPPostsCount := ui.NewLabel("Posts cound")
	boxWndWPPostsCount.Append(lblWndWPPostsCount, true)
	sboxWndWPPostsCount := ui.NewSpinbox(1, 50)
	sboxWndWPPostsCount.SetValue(wallPostMonitorParam.PostsCount)
	boxWndWPPostsCount.Append(sboxWndWPPostsCount, true)

	// описываем коробку с меткой и полем для списка ключевых слов для поиска постов
	boxWndWPKwrdsForMntrng := ui.NewHorizontalBox()
	boxWndWPKwrdsForMntrng.SetPadded(true)
	lblWndWPKwrdsForMntrng := ui.NewLabel("Keywords")
	boxWndWPKwrdsForMntrng.Append(lblWndWPKwrdsForMntrng, true)
	entryWndWPKwrdsForMntrng := ui.NewEntry()
	listKwrdsForMntrng, err := MakeParamList(wallPostMonitorParam.KeywordsForMonitoring)
	if err != nil {
		date := UnixTimeStampToDate(int(time.Now().Unix()))
		log.Fatal(fmt.Errorf("> [%v] WARNING! Error: %v", date, err))
	}
	if len(listKwrdsForMntrng.List) > 0 {
		var kwrdsForMntrng string
		for i, keyword := range listKwrdsForMntrng.List {
			if i > 0 {
				kwrdsForMntrng += ", "
			}
			kwrdsForMntrng += fmt.Sprintf("\"%v\"", keyword)
		}
		entryWndWPKwrdsForMntrng.SetText(kwrdsForMntrng)
	}
	boxWndWPKwrdsForMntrng.Append(entryWndWPKwrdsForMntrng, true)

	// описываем группу, в которой будут размещены элементы
	groupWndWP := ui.NewGroup("")
	groupWndWP.SetMargined(true)
	boxWndWP.Append(boxWndWPMonitoring, false)
	boxWndWP.Append(boxWndWPInterval, false)
	boxWndWP.Append(boxWndWPSendTo, false)
	boxWndWP.Append(boxWndWPFilter, false)
	boxWndWP.Append(boxWndWPPostsCount, false)
	boxWndWP.Append(boxWndWPKwrdsForMntrng, false)
	groupWndWP.SetChild(boxWndWP)

	// добавляем группу в основную коробку окна
	boxWndMain.Append(groupWndWP, false)

	// описываем коробку для кнопок
	boxWndWPBtns := ui.NewHorizontalBox()
	boxWndWPBtns.SetPadded(true)
	// и несколько коробок для выравнивания кнопок
	btnWndWPBtnsLeft := ui.NewHorizontalBox()
	btnWndWPBtnsCenter := ui.NewHorizontalBox()
	btnWndWPBtnsRight := ui.NewHorizontalBox()
	btnWndWPBtnsRight.SetPadded(true)
	// а затем сами кнопки
	btnWndWPCancel := ui.NewButton("Cancel")
	btnWndWPBtnsRight.Append(btnWndWPCancel, false)
	btnWndWPApplyChanges := ui.NewButton("Apply")
	btnWndWPBtnsRight.Append(btnWndWPApplyChanges, false)
	// и добавляем их в коробку для кнопок
	boxWndWPBtns.Append(btnWndWPBtnsLeft, false)
	boxWndWPBtns.Append(btnWndWPBtnsCenter, false)
	boxWndWPBtns.Append(btnWndWPBtnsRight, false)

	// привязываем к кнопкам соответствующие процедуры
	btnWndWPCancel.OnClicked(func(*ui.Button) {
		// TODO: как-нибудь надо закрывать окно
	})
	// привязываем кнопки к соответствующим процедурам
	btnWndWPApplyChanges.OnClicked(func(*ui.Button) {
		var updatedWallPostMonitorParam WallPostMonitorParam
		updatedWallPostMonitorParam.ID = wallPostMonitorParam.ID
		updatedWallPostMonitorParam.SubjectID = wallPostMonitorParam.SubjectID
		if cboxWndWPNeedMonitoring.Checked() {
			updatedWallPostMonitorParam.NeedMonitoring = 1
		} else {
			updatedWallPostMonitorParam.NeedMonitoring = 0
		}
		updatedWallPostMonitorParam.Interval = sboxWndWPInterval.Value()
		updatedWallPostMonitorParam.SendTo, err = strconv.Atoi(entryWndWPSendTo.Text())
		if err != nil {
			date := UnixTimeStampToDate(int(time.Now().Unix()))
			log.Fatal(fmt.Errorf("> [%v] WARNING! Error: %v", date, err))
		}
		listPostsFilters := []string{"all", "others", "owner", "suggests"}
		updatedWallPostMonitorParam.Filter = listPostsFilters[comboboxWndWPFilter.Selected()]
		updatedWallPostMonitorParam.LastDate = wallPostMonitorParam.LastDate
		updatedWallPostMonitorParam.PostsCount = sboxWndWPPostsCount.Value()
		jsonDump := fmt.Sprintf("{\"list\":[%v]}", entryWndWPKwrdsForMntrng.Text())
		updatedWallPostMonitorParam.KeywordsForMonitoring = jsonDump
		updatedWallPostMonitorParam.UsersIDsForIgnore = wallPostMonitorParam.UsersIDsForIgnore

		err = UpdateDBWallPostMonitor(updatedWallPostMonitorParam)
		if err != nil {
			date := UnixTimeStampToDate(int(time.Now().Unix()))
			log.Fatal(fmt.Errorf("> [%v] WARNING! Error: %v", date, err))
		}

		// TODO: как-нибудь надо закрывать окно
	})

	// добавляем коробку с кнопками на основную коробку окна
	boxWndMain.Append(boxWndWPBtns, true)
	// затем еще одну коробку, для выравнивания расположения кнопок при растягивании окна
	boxWndWPBottom := ui.NewHorizontalBox()
	boxWndMain.Append(boxWndWPBottom, false)

	wndSubjectWallPostSettings.Show()
}

func makePrimarySettingsBox(generalBoxesData GeneralBoxesData, groupsSettingsData GroupsSettingsData) *ui.Box {
	// описываем коробку для первичных установок
	boxPrimarySettings := ui.NewVerticalBox()

	// описываем кнопку для отображения установок токенов доступа
	btnAccessTokensSettings := ui.NewButton("Access tokens")
	// по умолчанию делаем ее неактивной
	btnAccessTokensSettings.Disable()
	// описываем кнопку для отображения установок субъектов
	btnSubjectsSettings := ui.NewButton("Subjects")

	// привязываем кнопки к процедурам отображения соответствующих блоков настроек
	btnAccessTokensSettings.OnClicked(func(*ui.Button) {
		groupsSettingsData.Additional.SetChild(ui.NewLabel("Nothing to show here..."))
		groupsSettingsData.General.SetTitle("Access tokens")
		groupsSettingsData.General.SetChild(generalBoxesData.AccessTokens)
		btnAccessTokensSettings.Disable()
		if !(btnSubjectsSettings.Enabled()) {
			btnSubjectsSettings.Enable()
		}
	})
	btnSubjectsSettings.OnClicked(func(*ui.Button) {
		groupsSettingsData.General.SetTitle("Subjects")
		groupsSettingsData.General.SetChild(generalBoxesData.Subjects)
		btnSubjectsSettings.Disable()
		if !(btnAccessTokensSettings.Enabled()) {
			btnAccessTokensSettings.Enable()
		}
	})

	// добавляем кнопки на коробку для первичных установок
	boxPrimarySettings.Append(btnAccessTokensSettings, false)
	boxPrimarySettings.Append(btnSubjectsSettings, false)

	return boxPrimarySettings
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
