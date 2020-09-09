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
	boxAccessTokensSettings := ui.NewVerticalBox()

	// запрашиваем список токенов доступа из базы данных
	accessTokens, err := SelectDBAccessTokens()
	if err != nil {
		date := UnixTimeStampToDate(int(time.Now().Unix()))
		log.Fatal(fmt.Errorf("> [%v] WARNING! Error: %v", date, err))
	}

	// описываем коробку для установок токенов доступа
	boxATUpper := ui.NewVerticalBox()

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
		boxATUpper.Append(boxAccessTokenSettings, false)
	}

	// описываем коробку для кнопки добавления нового токена доступа
	boxATBottom := ui.NewHorizontalBox()

	// описываем кнопку для добавления нового токена и декоративную кнопку для выравнивания
	btnAddNewAccessToken := ui.NewButton("＋")
	btnDecorative := ui.NewButton("")
	btnDecorative.Disable()

	// привязываем к кнопке для добавления соответствующую процедуру
	btnAddNewAccessToken.OnClicked(func(*ui.Button) {
		showAccessTokenAdditionWindow()
	})

	// добавляем кнопки на коробку для кнопок
	boxATBottom.Append(btnAddNewAccessToken, false)
	boxATBottom.Append(btnDecorative, true)

	// добавляем обе коробки на основную коробку
	boxAccessTokensSettings.Append(boxATUpper, false)
	boxAccessTokensSettings.Append(boxATBottom, false)

	return boxAccessTokensSettings
}

func showAccessTokenAdditionWindow() {
	// получаем набор для отображения окна для добавления нового токена доступа
	windowTitle := fmt.Sprintf("New access token addition")
	kitWindowAccessTokenAddition := makeSettingWindowKit(windowTitle, 300, 100)

	boxWndATAddition := ui.NewVerticalBox()

	// получаем набор для ввода названия нового токена доступа
	kitATCreationName := makeSettingEntryKit("Name", "")

	// получаем набор для ввода значения нового токена доступа
	kitATCreationValue := makeSettingEntryKit("Value", "")

	// описываем группу, в которой будут размещены элементы
	groupWndATAddition := ui.NewGroup("")
	groupWndATAddition.SetMargined(true)
	boxWndATAddition.Append(kitATCreationName.Box, false)
	boxWndATAddition.Append(kitATCreationValue.Box, false)
	groupWndATAddition.SetChild(boxWndATAddition)

	// добавляем группу в основную коробку окна
	kitWindowAccessTokenAddition.Box.Append(groupWndATAddition, false)

	// получаем набор для кнопок принятия и отмены изменений
	kitButtonsATAddition := makeSettingButtonsKit()

	// привязываем кнопки к соответствующим процедурам
	kitButtonsATAddition.ButtonCancel.OnClicked(func(*ui.Button) {
		// TODO: как-нибудь надо закрывать окно
	})
	kitButtonsATAddition.ButtonApply.OnClicked(func(*ui.Button) {
		var accessToken AccessToken

		accessToken.Name = kitATCreationName.Entry.Text()
		accessToken.Value = kitATCreationValue.Entry.Text()

		err := InsertDBAccessToken(accessToken)
		if err != nil {
			date := UnixTimeStampToDate(int(time.Now().Unix()))
			log.Fatal(fmt.Errorf("> [%v] WARNING! Error: %v", date, err))
		}

		// TODO: как-нибудь надо закрывать окно
	})

	// добавляем коробку с кнопками на основную коробку окна
	kitWindowAccessTokenAddition.Box.Append(kitButtonsATAddition.Box, true)
	// затем еще одну коробку, для выравнивания расположения кнопок при растягивании окна
	boxWndATAdditionBottom := ui.NewHorizontalBox()
	kitWindowAccessTokenAddition.Box.Append(boxWndATAdditionBottom, false)

	kitWindowAccessTokenAddition.Window.Show()
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
		case "Album photo monitor":
			btnSettingsSection.OnClicked(func(*ui.Button) {
				showSubjectAlbumPhotoSettingWindow(subjectData.ID, subjectData.Name, btnName)
			})
		case "Video monitor":
			btnSettingsSection.OnClicked(func(*ui.Button) {
				showSubjectVideoSettingWindow(subjectData.ID, subjectData.Name, btnName)
			})
		case "Photo comment monitor":
			btnSettingsSection.OnClicked(func(*ui.Button) {
				showSubjectPhotoCommentSettingWindow(subjectData.ID, subjectData.Name, btnName)
			})
		case "Video comment monitor":
			btnSettingsSection.OnClicked(func(*ui.Button) {
				showSubjectVideoCommentSettingWindow(subjectData.ID, subjectData.Name, btnName)
			})
		case "Topic monitor":
			btnSettingsSection.OnClicked(func(*ui.Button) {
				showSubjectTopicSettingWindow(subjectData.ID, subjectData.Name, btnName)
			})
		case "Wall post comment monitor":
			btnSettingsSection.OnClicked(func(*ui.Button) {
				showSubjectWallPostCommentSettings(subjectData.ID, subjectData.Name, btnName)
			})
		}

		boxSubjectAdditionalSettingsBox.Append(boxSettingsSection, false)
	}

	return boxSubjectAdditionalSettingsBox
}

func showSubjectGeneralSettingWindow(IDSubject int, btnName string) {
	// запрашиваем список субъектов из базы данных
	subjects, err := SelectDBSubjects()
	if err != nil {
		date := UnixTimeStampToDate(int(time.Now().Unix()))
		log.Fatal(fmt.Errorf("> [%v] WARNING! Error: %v", date, err))
	}

	// получаем набор для отображения общих установок субъекта мониторинга
	kitWindowGeneralSettings := makeSettingWindowKit("", 300, 100)

	// перечисляем субъекты
	for _, subject := range subjects {
		// ищем субъект с подходящим идентификатором
		if subject.ID == IDSubject {
			// устанавливаем заголовок окна в соответствии с названием субъекта и назначением установок
			windowTitle := fmt.Sprintf("%v settings for %v", btnName, subject.Name)
			kitWindowGeneralSettings.Window.SetTitle(windowTitle)

			boxWndS := ui.NewVerticalBox()

			// получаем набор для названия субъекта мониторинга
			kitWndSName := makeSettingEntryKit("Name", subject.Name)

			// получаем набор для идентификатора субъекта мониторинга в базе ВК
			kitWndSSubjectID := makeSettingEntryKit("Subject ID", strconv.Itoa(subject.SubjectID))

			// описываем группу, в которой будут размещены элементы
			groupWndS := ui.NewGroup("")
			groupWndS.SetMargined(true)
			boxWndS.Append(kitWndSName.Box, false)
			boxWndS.Append(kitWndSSubjectID.Box, false)
			groupWndS.SetChild(boxWndS)

			// добавляем группу в основную коробку окна
			kitWindowGeneralSettings.Box.Append(groupWndS, false)

			// получаем набор для кнопок принятия и отмены изменений
			kitButtonsS := makeSettingButtonsKit()

			// привязываем к кнопкам соответствующие процедуры
			kitButtonsS.ButtonCancel.OnClicked(func(*ui.Button) {
				// TODO: как-нибудь надо закрывать окно
			})
			// привязываем кнопки к соответствующим процедурам
			kitButtonsS.ButtonApply.OnClicked(func(*ui.Button) {
				var updatedSubject Subject
				updatedSubject.ID = subject.ID
				updatedSubject.SubjectID, err = strconv.Atoi(kitWndSSubjectID.Entry.Text())
				if err != nil {
					date := UnixTimeStampToDate(int(time.Now().Unix()))
					log.Fatal(fmt.Errorf("> [%v] WARNING! Error: %v", date, err))
				}
				updatedSubject.Name = kitWndSName.Entry.Text()
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
			kitWindowGeneralSettings.Box.Append(kitButtonsS.Box, true)
			// затем еще одну коробку, для выравнивания расположения кнопок при растягивании окна
			boxWndSBottom := ui.NewHorizontalBox()
			kitWindowGeneralSettings.Box.Append(boxWndSBottom, false)

			break
		}
	}

	kitWindowGeneralSettings.Window.Show()
}

func showSubjectWallPostSettingWindow(IDSubject int, nameSubject, btnName string) {
	// запрашиваем параметры мониторинга из базы данных
	wallPostMonitorParam, err := SelectDBWallPostMonitorParam(IDSubject)
	if err != nil {
		date := UnixTimeStampToDate(int(time.Now().Unix()))
		log.Fatal(fmt.Errorf("> [%v] WARNING! Error: %v", date, err))
	}

	// получаем набор для отображения установок модуля мониторинга постов на стене
	windowTitle := fmt.Sprintf("%v settings for %v", btnName, nameSubject)
	kitWindowWallPostSettings := makeSettingWindowKit(windowTitle, 300, 100)

	boxWndWP := ui.NewVerticalBox()

	// получаем набор для флага необходимости активировать модуль мониторинга
	kitWndWPMonitoring := makeSettingCheckboxKit("Need monitoring", wallPostMonitorParam.NeedMonitoring)

	// получаем набор для интервала между запусками функции мониторинга
	kitWndWPInterval := makeSettingSpinboxKit("Interval", 5, 21600, wallPostMonitorParam.Interval)

	// получаем набор для количества проверяемых постов
	kitWndWPSendTo := makeSettingEntryKit("Send to", strconv.Itoa(wallPostMonitorParam.SendTo))

	// получаем набор для фильтра получаемых для проверки постов
	listPostsFilters := []string{"all", "others", "owner", "suggests"}
	kitWndWPFilter := makeSettingComboboxKit("Filter", listPostsFilters, wallPostMonitorParam.Filter)

	// получаем набор для количества проверяемых постов
	kitWndWPPostsCount := makeSettingSpinboxKit("Posts count", 1, 100, wallPostMonitorParam.PostsCount)

	// получаем набор для списка ключевых слов для отбора постов
	kitWndWPKeywordsForMonitoring := makeSettingEntryListKit("Keywords", wallPostMonitorParam.KeywordsForMonitoring)

	// описываем группу, в которой будут размещены элементы
	groupWndWP := ui.NewGroup("")
	groupWndWP.SetMargined(true)
	boxWndWP.Append(kitWndWPMonitoring.Box, false)
	boxWndWP.Append(kitWndWPInterval.Box, false)
	boxWndWP.Append(kitWndWPSendTo.Box, false)
	boxWndWP.Append(kitWndWPFilter.Box, false)
	boxWndWP.Append(kitWndWPPostsCount.Box, false)
	boxWndWP.Append(kitWndWPKeywordsForMonitoring.Box, false)
	groupWndWP.SetChild(boxWndWP)

	// добавляем группу в основную коробку окна
	kitWindowWallPostSettings.Box.Append(groupWndWP, false)

	// получаем набор для кнопок принятия и отмены изменений
	kitButtonsWP := makeSettingButtonsKit()

	// привязываем к кнопкам соответствующие процедуры
	kitButtonsWP.ButtonCancel.OnClicked(func(*ui.Button) {
		// TODO: как-нибудь надо закрывать окно
	})
	// привязываем кнопки к соответствующим процедурам
	kitButtonsWP.ButtonApply.OnClicked(func(*ui.Button) {
		var updatedWallPostMonitorParam WallPostMonitorParam
		updatedWallPostMonitorParam.ID = wallPostMonitorParam.ID
		updatedWallPostMonitorParam.SubjectID = wallPostMonitorParam.SubjectID
		if kitWndWPMonitoring.CheckBox.Checked() {
			updatedWallPostMonitorParam.NeedMonitoring = 1
		} else {
			updatedWallPostMonitorParam.NeedMonitoring = 0
		}
		updatedWallPostMonitorParam.Interval = kitWndWPInterval.Spinbox.Value()
		updatedWallPostMonitorParam.SendTo, err = strconv.Atoi(kitWndWPSendTo.Entry.Text())
		if err != nil {
			date := UnixTimeStampToDate(int(time.Now().Unix()))
			log.Fatal(fmt.Errorf("> [%v] WARNING! Error: %v", date, err))
		}
		listPostsFilters := []string{"all", "others", "owner", "suggests"}
		updatedWallPostMonitorParam.Filter = listPostsFilters[kitWndWPFilter.Combobox.Selected()]
		updatedWallPostMonitorParam.LastDate = wallPostMonitorParam.LastDate
		updatedWallPostMonitorParam.PostsCount = kitWndWPPostsCount.Spinbox.Value()
		jsonDump := fmt.Sprintf("{\"list\":[%v]}", kitWndWPKeywordsForMonitoring.Entry.Text())
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
	kitWindowWallPostSettings.Box.Append(kitButtonsWP.Box, true)
	// затем еще одну коробку, для выравнивания расположения кнопок при растягивании окна
	boxWndWPBottom := ui.NewHorizontalBox()
	kitWindowWallPostSettings.Box.Append(boxWndWPBottom, false)

	kitWindowWallPostSettings.Window.Show()
}

func showSubjectAlbumPhotoSettingWindow(IDSubject int, nameSubject, btnName string) {
	// запрашиваем параметры мониторинга из базы данных
	albumPhotoMonitorParam, err := SelectDBAlbumPhotoMonitorParam(IDSubject)
	if err != nil {
		date := UnixTimeStampToDate(int(time.Now().Unix()))
		log.Fatal(fmt.Errorf("> [%v] WARNING! Error: %v", date, err))
	}

	// получаем набор для отображения установок модуля мониторинга фотографий в альбомах
	windowTitle := fmt.Sprintf("%v settings for %v", btnName, nameSubject)
	kitWindowAlbumPhotoSettings := makeSettingWindowKit(windowTitle, 300, 100)

	boxWndAP := ui.NewVerticalBox()

	// получаем набор для флага необходимости активировать модуль мониторинга
	kitWndAPMonitoring := makeSettingCheckboxKit("Need monitoring", albumPhotoMonitorParam.NeedMonitoring)

	// получаем набор для интервала между запусками функции мониторинга
	kitWndAPInterval := makeSettingSpinboxKit("Interval", 5, 21600, albumPhotoMonitorParam.Interval)

	// получаем набор для идентификатора получателя сообщений
	kitWndAPSendTo := makeSettingEntryKit("Send to", strconv.Itoa(albumPhotoMonitorParam.SendTo))

	// получаем набор для количества проверяемых фото
	kitWndApPhotosCount := makeSettingSpinboxKit("Photos count", 1, 1000, albumPhotoMonitorParam.PhotosCount)

	// описываем группу, в которой будут размещены элементы
	groupWndAP := ui.NewGroup("")
	groupWndAP.SetMargined(true)
	boxWndAP.Append(kitWndAPMonitoring.Box, false)
	boxWndAP.Append(kitWndAPInterval.Box, false)
	boxWndAP.Append(kitWndAPSendTo.Box, false)
	boxWndAP.Append(kitWndApPhotosCount.Box, false)
	groupWndAP.SetChild(boxWndAP)

	// добавляем группу в основную коробку окна
	kitWindowAlbumPhotoSettings.Box.Append(groupWndAP, false)

	// получаем набор для кнопок принятия и отмены изменений
	kitButtonsAP := makeSettingButtonsKit()

	// привязываем к кнопкам соответствующие процедуры
	kitButtonsAP.ButtonCancel.OnClicked(func(*ui.Button) {
		// TODO: как-нибудь надо закрывать окно
	})
	// привязываем кнопки к соответствующим процедурам
	kitButtonsAP.ButtonApply.OnClicked(func(*ui.Button) {
		var updatedAlbumPhotoMonitorParam AlbumPhotoMonitorParam
		updatedAlbumPhotoMonitorParam.ID = albumPhotoMonitorParam.ID
		updatedAlbumPhotoMonitorParam.SubjectID = albumPhotoMonitorParam.SubjectID
		if kitWndAPMonitoring.CheckBox.Checked() {
			updatedAlbumPhotoMonitorParam.NeedMonitoring = 1
		} else {
			updatedAlbumPhotoMonitorParam.NeedMonitoring = 0
		}
		updatedAlbumPhotoMonitorParam.SendTo, err = strconv.Atoi(kitWndAPSendTo.Entry.Text())
		if err != nil {
			date := UnixTimeStampToDate(int(time.Now().Unix()))
			log.Fatal(fmt.Errorf("> [%v] WARNING! Error: %v", date, err))
		}
		updatedAlbumPhotoMonitorParam.Interval = kitWndAPInterval.Spinbox.Value()
		updatedAlbumPhotoMonitorParam.LastDate = albumPhotoMonitorParam.LastDate
		updatedAlbumPhotoMonitorParam.PhotosCount = kitWndApPhotosCount.Spinbox.Value()

		err = UpdateDBAlbumPhotoMonitor(updatedAlbumPhotoMonitorParam)
		if err != nil {
			date := UnixTimeStampToDate(int(time.Now().Unix()))
			log.Fatal(fmt.Errorf("> [%v] WARNING! Error: %v", date, err))
		}

		// TODO: как-нибудь надо закрывать окно
	})

	// добавляем коробку с кнопками на основную коробку окна
	kitWindowAlbumPhotoSettings.Box.Append(kitButtonsAP.Box, true)
	// затем еще одну коробку, для выравнивания расположения кнопок при растягивании окна
	boxWndAPBottom := ui.NewHorizontalBox()
	kitWindowAlbumPhotoSettings.Box.Append(boxWndAPBottom, false)

	kitWindowAlbumPhotoSettings.Window.Show()
}

func showSubjectVideoSettingWindow(IDSubject int, nameSubject, btnName string) {
	// запрашиваем параметры мониторинга из базы данных
	videoMonitorParam, err := SelectDBVideoMonitorParam(IDSubject)
	if err != nil {
		date := UnixTimeStampToDate(int(time.Now().Unix()))
		log.Fatal(fmt.Errorf("> [%v] WARNING! Error: %v", date, err))
	}

	// получаем набор для отображения установок модуля мониторинга видео
	windowTitle := fmt.Sprintf("%v settings for %v", btnName, nameSubject)
	kitWindowVideoSettings := makeSettingWindowKit(windowTitle, 300, 100)

	boxWndV := ui.NewVerticalBox()

	// получаем набор для флага необходимости активировать модуль мониторинга
	kitWndVMonitoring := makeSettingCheckboxKit("Need monitoring", videoMonitorParam.NeedMonitoring)

	// получаем набор для интервала между запусками функции мониторинга
	kitWndVInterval := makeSettingSpinboxKit("Interval", 5, 21600, videoMonitorParam.Interval)

	// получаем набор для идентификатора получателя сообщений
	kitWndVSendTo := makeSettingEntryKit("Send to", strconv.Itoa(videoMonitorParam.SendTo))

	// получаем набор для количества проверяемых видео
	kitWndVVideoCount := makeSettingSpinboxKit("Video count", 1, 1000, videoMonitorParam.VideoCount)

	// описываем группу, в которой будут размещены элементы
	groupWndV := ui.NewGroup("")
	groupWndV.SetMargined(true)
	boxWndV.Append(kitWndVMonitoring.Box, false)
	boxWndV.Append(kitWndVInterval.Box, false)
	boxWndV.Append(kitWndVSendTo.Box, false)
	boxWndV.Append(kitWndVVideoCount.Box, false)
	groupWndV.SetChild(boxWndV)

	// добавляем группу в основную коробку окна
	kitWindowVideoSettings.Box.Append(groupWndV, false)

	// получаем набор для кнопок принятия и отмены изменений
	kitButtonsV := makeSettingButtonsKit()

	// привязываем к кнопкам соответствующие процедуры
	kitButtonsV.ButtonCancel.OnClicked(func(*ui.Button) {
		// TODO: как-нибудь надо закрывать окно
	})
	// привязываем кнопки к соответствующим процедурам
	kitButtonsV.ButtonApply.OnClicked(func(*ui.Button) {
		var updatedVideoMonitorParam VideoMonitorParam
		updatedVideoMonitorParam.ID = videoMonitorParam.ID
		updatedVideoMonitorParam.SubjectID = videoMonitorParam.SubjectID
		if kitWndVMonitoring.CheckBox.Checked() {
			updatedVideoMonitorParam.NeedMonitoring = 1
		} else {
			updatedVideoMonitorParam.NeedMonitoring = 0
		}
		updatedVideoMonitorParam.SendTo, err = strconv.Atoi(kitWndVSendTo.Entry.Text())
		if err != nil {
			date := UnixTimeStampToDate(int(time.Now().Unix()))
			log.Fatal(fmt.Errorf("> [%v] WARNING! Error: %v", date, err))
		}
		updatedVideoMonitorParam.Interval = kitWndVInterval.Spinbox.Value()
		updatedVideoMonitorParam.LastDate = videoMonitorParam.LastDate
		updatedVideoMonitorParam.VideoCount = kitWndVVideoCount.Spinbox.Value()

		err = UpdateDBVideoMonitor(updatedVideoMonitorParam)
		if err != nil {
			date := UnixTimeStampToDate(int(time.Now().Unix()))
			log.Fatal(fmt.Errorf("> [%v] WARNING! Error: %v", date, err))
		}

		// TODO: как-нибудь надо закрывать окно
	})

	// добавляем коробку с кнопками на основную коробку окна
	kitWindowVideoSettings.Box.Append(kitButtonsV.Box, true)
	// затем еще одну коробку, для выравнивания расположения кнопок при растягивании окна
	boxWndVBottom := ui.NewHorizontalBox()
	kitWindowVideoSettings.Box.Append(boxWndVBottom, false)

	kitWindowVideoSettings.Window.Show()
}

func showSubjectPhotoCommentSettingWindow(IDSubject int, nameSubject, btnName string) {
	// запрашиваем параметры мониторинга из базы данных
	photoCommentMonitorParam, err := SelectDBPhotoCommentMonitorParam(IDSubject)
	if err != nil {
		date := UnixTimeStampToDate(int(time.Now().Unix()))
		log.Fatal(fmt.Errorf("> [%v] WARNING! Error: %v", date, err))
	}

	// получаем набор для отображения установок модуля мониторинга комментариев под фотками
	windowTitle := fmt.Sprintf("%v settings for %v", btnName, nameSubject)
	kitWindowPhotoCommentSettings := makeSettingWindowKit(windowTitle, 300, 100)

	boxWndPC := ui.NewVerticalBox()

	// получаем набор для флага необходимости активировать модуль мониторинга
	kitWndPCMonitoring := makeSettingCheckboxKit("Need monitoring", photoCommentMonitorParam.NeedMonitoring)

	// получаем набор для интервала между запусками функции мониторинга
	kitWndPCInterval := makeSettingSpinboxKit("Interval", 5, 21600, photoCommentMonitorParam.Interval)

	// получаем набор для идентификатора получателя сообщений
	kitWndPCSendTo := makeSettingEntryKit("Send to", strconv.Itoa(photoCommentMonitorParam.SendTo))

	// получаем набор для количества проверяемых комментариев
	kitWndPCCommentsCount := makeSettingSpinboxKit("Comments count", 1, 1000, photoCommentMonitorParam.CommentsCount)

	// описываем группу, в которой будут размещены элементы
	groupWndPC := ui.NewGroup("")
	groupWndPC.SetMargined(true)
	boxWndPC.Append(kitWndPCMonitoring.Box, false)
	boxWndPC.Append(kitWndPCInterval.Box, false)
	boxWndPC.Append(kitWndPCSendTo.Box, false)
	boxWndPC.Append(kitWndPCCommentsCount.Box, false)
	groupWndPC.SetChild(boxWndPC)

	// добавляем группу в основную коробку окна
	kitWindowPhotoCommentSettings.Box.Append(groupWndPC, false)

	// получаем набор для кнопок принятия и отмены изменений
	kitButtonsPC := makeSettingButtonsKit()

	// привязываем к кнопкам соответствующие процедуры
	kitButtonsPC.ButtonCancel.OnClicked(func(*ui.Button) {
		// TODO: как-нибудь надо закрывать окно
	})
	// привязываем кнопки к соответствующим процедурам
	kitButtonsPC.ButtonApply.OnClicked(func(*ui.Button) {
		var updatedPhotoCommentMonitorParam PhotoCommentMonitorParam
		updatedPhotoCommentMonitorParam.ID = photoCommentMonitorParam.ID
		updatedPhotoCommentMonitorParam.SubjectID = photoCommentMonitorParam.SubjectID
		if kitWndPCMonitoring.CheckBox.Checked() {
			updatedPhotoCommentMonitorParam.NeedMonitoring = 1
		} else {
			updatedPhotoCommentMonitorParam.NeedMonitoring = 0
		}
		updatedPhotoCommentMonitorParam.SendTo, err = strconv.Atoi(kitWndPCSendTo.Entry.Text())
		if err != nil {
			date := UnixTimeStampToDate(int(time.Now().Unix()))
			log.Fatal(fmt.Errorf("> [%v] WARNING! Error: %v", date, err))
		}
		updatedPhotoCommentMonitorParam.Interval = kitWndPCInterval.Spinbox.Value()
		updatedPhotoCommentMonitorParam.LastDate = photoCommentMonitorParam.LastDate
		updatedPhotoCommentMonitorParam.CommentsCount = kitWndPCCommentsCount.Spinbox.Value()

		err = UpdateDBPhotoCommentMonitor(updatedPhotoCommentMonitorParam)
		if err != nil {
			date := UnixTimeStampToDate(int(time.Now().Unix()))
			log.Fatal(fmt.Errorf("> [%v] WARNING! Error: %v", date, err))
		}

		// TODO: как-нибудь надо закрывать окно
	})

	// добавляем коробку с кнопками на основную коробку окна
	kitWindowPhotoCommentSettings.Box.Append(kitButtonsPC.Box, true)
	// затем еще одну коробку, для выравнивания расположения кнопок при растягивании окна
	boxWndPCBottom := ui.NewHorizontalBox()
	kitWindowPhotoCommentSettings.Box.Append(boxWndPCBottom, false)

	kitWindowPhotoCommentSettings.Window.Show()
}

func showSubjectVideoCommentSettingWindow(IDSubject int, nameSubject, btnName string) {
	// запрашиваем параметры мониторинга из базы данных
	videoCommentMonitorParam, err := SelectDBVideoCommentMonitorParam(IDSubject)
	if err != nil {
		date := UnixTimeStampToDate(int(time.Now().Unix()))
		log.Fatal(fmt.Errorf("> [%v] WARNING! Error: %v", date, err))
	}

	// получаем набор для отображения установок модуля мониторинга комментариев в обсуждениях
	windowTitle := fmt.Sprintf("%v settings for %v", btnName, nameSubject)
	kitWindowVideoCommentSettings := makeSettingWindowKit(windowTitle, 300, 100)

	boxWndVC := ui.NewVerticalBox()

	// получаем набор для флага необходимости активировать модуль мониторинга
	kitWndVCMonitoring := makeSettingCheckboxKit("Need monitoring", videoCommentMonitorParam.NeedMonitoring)

	// получаем набор для интервала между запусками функции мониторинга
	kitWndVCInterval := makeSettingSpinboxKit("Interval", 5, 21600, videoCommentMonitorParam.Interval)

	// получаем набор для идентификатора получателя сообщений
	kitWndVCSendTo := makeSettingEntryKit("Send to", strconv.Itoa(videoCommentMonitorParam.SendTo))

	// получаем набор для количества проверяемых видео
	kitWndVCVideosCount := makeSettingSpinboxKit("Videos count", 1, 200, videoCommentMonitorParam.VideosCount)

	// получаем набор для количества проверяемых комментариев
	kitWndVCCommentsCount := makeSettingSpinboxKit("Comments count", 1, 100, videoCommentMonitorParam.CommentsCount)

	// описываем группу, в которой будут размещены элементы
	groupWndVC := ui.NewGroup("")
	groupWndVC.SetMargined(true)
	boxWndVC.Append(kitWndVCMonitoring.Box, false)
	boxWndVC.Append(kitWndVCSendTo.Box, false)
	boxWndVC.Append(kitWndVCInterval.Box, false)
	boxWndVC.Append(kitWndVCVideosCount.Box, false)
	boxWndVC.Append(kitWndVCCommentsCount.Box, false)
	groupWndVC.SetChild(boxWndVC)

	// добавляем группу в основную коробку окна
	kitWindowVideoCommentSettings.Box.Append(groupWndVC, false)

	// получаем набор для кнопок принятия и отмены изменений
	kitButtonsVC := makeSettingButtonsKit()

	// привязываем к кнопкам соответствующие процедуры
	kitButtonsVC.ButtonCancel.OnClicked(func(*ui.Button) {
		// TODO: как-нибудь надо закрывать окно
	})
	// привязываем кнопки к соответствующим процедурам
	kitButtonsVC.ButtonApply.OnClicked(func(*ui.Button) {
		var updatedVideoCommentMonitorParam VideoCommentMonitorParam
		updatedVideoCommentMonitorParam.ID = videoCommentMonitorParam.ID
		updatedVideoCommentMonitorParam.SubjectID = videoCommentMonitorParam.SubjectID
		if kitWndVCMonitoring.CheckBox.Checked() {
			updatedVideoCommentMonitorParam.NeedMonitoring = 1
		} else {
			updatedVideoCommentMonitorParam.NeedMonitoring = 0
		}
		updatedVideoCommentMonitorParam.SendTo, err = strconv.Atoi(kitWndVCSendTo.Entry.Text())
		if err != nil {
			date := UnixTimeStampToDate(int(time.Now().Unix()))
			log.Fatal(fmt.Errorf("> [%v] WARNING! Error: %v", date, err))
		}
		updatedVideoCommentMonitorParam.Interval = kitWndVCInterval.Spinbox.Value()
		updatedVideoCommentMonitorParam.LastDate = videoCommentMonitorParam.LastDate
		updatedVideoCommentMonitorParam.CommentsCount = kitWndVCCommentsCount.Spinbox.Value()
		updatedVideoCommentMonitorParam.VideosCount = kitWndVCVideosCount.Spinbox.Value()

		err = UpdateDBVideoCommentMonitor(updatedVideoCommentMonitorParam)
		if err != nil {
			date := UnixTimeStampToDate(int(time.Now().Unix()))
			log.Fatal(fmt.Errorf("> [%v] WARNING! Error: %v", date, err))
		}

		// TODO: как-нибудь надо закрывать окно
	})

	// добавляем коробку с кнопками на основную коробку окна
	kitWindowVideoCommentSettings.Box.Append(kitButtonsVC.Box, true)
	// затем еще одну коробку, для выравнивания расположения кнопок при растягивании окна
	boxWndVCBottom := ui.NewHorizontalBox()
	kitWindowVideoCommentSettings.Box.Append(boxWndVCBottom, false)

	kitWindowVideoCommentSettings.Window.Show()
}

func showSubjectTopicSettingWindow(IDSubject int, nameSubject, btnName string) {
	// запрашиваем параметры мониторинга из базы данных
	topicMonitorParam, err := SelectDBTopicMonitorParam(IDSubject)
	if err != nil {
		date := UnixTimeStampToDate(int(time.Now().Unix()))
		log.Fatal(fmt.Errorf("> [%v] WARNING! Error: %v", date, err))
	}

	// получаем набор для отображения установок модуля мониторинга комментариев в обсуждениях
	windowTitle := fmt.Sprintf("%v settings for %v", btnName, nameSubject)
	kitWindowTopicSettings := makeSettingWindowKit(windowTitle, 300, 100)

	boxWndT := ui.NewVerticalBox()

	// получаем набор для флага необходимости активировать модуль мониторинга
	kitWndTMonitoring := makeSettingCheckboxKit("Need monitoring", topicMonitorParam.NeedMonitoring)

	// получаем набор для интервала между запусками функции мониторинга
	kitWndTInterval := makeSettingSpinboxKit("Interval", 5, 21600, topicMonitorParam.Interval)

	// получаем набор для идентификатора получателя сообщений
	kitWndTSendTo := makeSettingEntryKit("Send to", strconv.Itoa(topicMonitorParam.SendTo))

	// получаем набор для количества проверяемых топиков обсуждений
	kitWndTTopicsCount := makeSettingSpinboxKit("Topics count", 1, 100, topicMonitorParam.TopicsCount)

	// получаем набор для количества проверяемых комментариев
	kitWndTCommentsCount := makeSettingSpinboxKit("Comments count", 1, 100, topicMonitorParam.TopicsCount)

	// описываем группу, в которой будут размещены элементы
	groupWndT := ui.NewGroup("")
	groupWndT.SetMargined(true)
	boxWndT.Append(kitWndTMonitoring.Box, false)
	boxWndT.Append(kitWndTSendTo.Box, false)
	boxWndT.Append(kitWndTInterval.Box, false)
	boxWndT.Append(kitWndTTopicsCount.Box, false)
	boxWndT.Append(kitWndTCommentsCount.Box, false)
	groupWndT.SetChild(boxWndT)

	// добавляем группу в основную коробку окна
	kitWindowTopicSettings.Box.Append(groupWndT, false)

	// получаем набор для кнопок принятия и отмены изменений
	kitButtonsT := makeSettingButtonsKit()

	// привязываем к кнопкам соответствующие процедуры
	kitButtonsT.ButtonCancel.OnClicked(func(*ui.Button) {
		// TODO: как-нибудь надо закрывать окно
	})
	// привязываем кнопки к соответствующим процедурам
	kitButtonsT.ButtonApply.OnClicked(func(*ui.Button) {
		var updatedTopicMonitorParam TopicMonitorParam
		updatedTopicMonitorParam.ID = topicMonitorParam.ID
		updatedTopicMonitorParam.SubjectID = topicMonitorParam.SubjectID
		if kitWndTMonitoring.CheckBox.Checked() {
			updatedTopicMonitorParam.NeedMonitoring = 1
		} else {
			updatedTopicMonitorParam.NeedMonitoring = 0
		}
		updatedTopicMonitorParam.SendTo, err = strconv.Atoi(kitWndTSendTo.Entry.Text())
		if err != nil {
			date := UnixTimeStampToDate(int(time.Now().Unix()))
			log.Fatal(fmt.Errorf("> [%v] WARNING! Error: %v", date, err))
		}
		updatedTopicMonitorParam.Interval = kitWndTInterval.Spinbox.Value()
		updatedTopicMonitorParam.LastDate = topicMonitorParam.LastDate
		updatedTopicMonitorParam.CommentsCount = kitWndTCommentsCount.Spinbox.Value()
		updatedTopicMonitorParam.TopicsCount = kitWndTTopicsCount.Spinbox.Value()

		err = UpdateDBTopicMonitor(updatedTopicMonitorParam)
		if err != nil {
			date := UnixTimeStampToDate(int(time.Now().Unix()))
			log.Fatal(fmt.Errorf("> [%v] WARNING! Error: %v", date, err))
		}

		// TODO: как-нибудь надо закрывать окно
	})

	// добавляем коробку с кнопками на основную коробку окна
	kitWindowTopicSettings.Box.Append(kitButtonsT.Box, true)
	// затем еще одну коробку, для выравнивания расположения кнопок при растягивании окна
	boxWndTBottom := ui.NewHorizontalBox()
	kitWindowTopicSettings.Box.Append(boxWndTBottom, false)

	kitWindowTopicSettings.Window.Show()
}

func showSubjectWallPostCommentSettings(IDSubject int, nameSubject, btnName string) {
	// запрашиваем параметры мониторинга из базы данных
	wallPostCommentMonitorParam, err := SelectDBWallPostCommentMonitorParam(IDSubject)
	if err != nil {
		date := UnixTimeStampToDate(int(time.Now().Unix()))
		log.Fatal(fmt.Errorf("> [%v] WARNING! Error: %v", date, err))
	}

	// получаем набор для отображения установок модуля мониторинга комментариев под постами
	windowTitle := fmt.Sprintf("%v settings for %v", btnName, nameSubject)
	kitWindowWallPostCommentSettings := makeSettingWindowKit(windowTitle, 300, 100)

	boxWndWPC := ui.NewVerticalBox()

	// получаем набор для флага необходимости активировать модуль мониторинга
	kitWndWPCMonitoring := makeSettingCheckboxKit("Need monitoring", wallPostCommentMonitorParam.NeedMonitoring)

	// получаем набор для интервала между запусками функции мониторинга
	kitWndWPCInterval := makeSettingSpinboxKit("Interval", 5, 21600, wallPostCommentMonitorParam.Interval)

	// получаем набор для идентификатора получателя сообщений
	kitWndWPCSendTo := makeSettingEntryKit("Send to", strconv.Itoa(wallPostCommentMonitorParam.SendTo))

	// получаем набор для количества проверяемых постов
	kitWndWPCPostsCount := makeSettingSpinboxKit("Posts count", 1, 100, wallPostCommentMonitorParam.PostsCount)

	// получаем набор для количества проверяемых комментариев
	kitWndWPCCommentsCount := makeSettingSpinboxKit("Comments count", 1, 100, wallPostCommentMonitorParam.CommentsCount)

	// получаем набор для фильтров постов для проверки комментариев
	listPostsFilters := []string{"all", "others", "owner"}
	kitWndWPCFilter := makeSettingComboboxKit("Filter", listPostsFilters, wallPostCommentMonitorParam.Filter)

	// получаем набор для флага необходимости проверять все комментарии без исключения
	kitWndWPCMonitoringAll := makeSettingCheckboxKit("Monitoring all", wallPostCommentMonitorParam.MonitoringAll)

	// получаем набор для флага необходимости проверять комментарии от сообществ
	kitWndWPCMonitorByCommunity := makeSettingCheckboxKit("Monitor by community", wallPostCommentMonitorParam.MonitorByCommunity)

	// получаем набор для списка ключевых слов для поиска комментариев
	kitWndWPCKeywordsForMonitoring := makeSettingEntryListKit("Keywords for monitoring", wallPostCommentMonitorParam.KeywordsForMonitoring)

	// получаем набор для списка комментариев для поиска
	kitWndWPCSmallCommentsForMonitoring := makeSettingEntryListKit("Small comments for monitoring", wallPostCommentMonitorParam.SmallCommentsForMonitoring)

	// получаем набор для списка имен и фамилий авторов комментариев для поиска комментариев
	kitWndWPCUsersNamesForMonitoring := makeSettingEntryListKit("Users names for monitoring", wallPostCommentMonitorParam.UsersNamesForMonitoring)

	// получаем набор для списка идентификаторов авторов комментариев для поиска комментариев
	kitWndWPCUsersIdsForMonitoring := makeSettingEntryListKit("Users IDs for monitoring", wallPostCommentMonitorParam.UsersIDsForMonitoring)

	// получаем набор для списка идентификаторов авторов комментариев для их игнорирования при проверке комментариев
	kitWndWPCUsersIdsForIgnore := makeSettingEntryListKit("Users IDs for ignore", wallPostCommentMonitorParam.UsersIDsForIgnore)

	// описываем группу, в которой будут размещены элементы
	groupWndWPC := ui.NewGroup("")
	groupWndWPC.SetMargined(true)
	boxWndWPC.Append(kitWndWPCMonitoring.Box, false)
	boxWndWPC.Append(kitWndWPCInterval.Box, false)
	boxWndWPC.Append(kitWndWPCSendTo.Box, false)
	boxWndWPC.Append(kitWndWPCPostsCount.Box, false)
	boxWndWPC.Append(kitWndWPCCommentsCount.Box, false)
	boxWndWPC.Append(kitWndWPCFilter.Box, false)
	boxWndWPC.Append(kitWndWPCMonitoringAll.Box, false)
	boxWndWPC.Append(kitWndWPCMonitorByCommunity.Box, false)
	boxWndWPC.Append(kitWndWPCKeywordsForMonitoring.Box, false)
	boxWndWPC.Append(kitWndWPCSmallCommentsForMonitoring.Box, false)
	boxWndWPC.Append(kitWndWPCUsersNamesForMonitoring.Box, false)
	boxWndWPC.Append(kitWndWPCUsersIdsForMonitoring.Box, false)
	boxWndWPC.Append(kitWndWPCUsersIdsForIgnore.Box, false)
	groupWndWPC.SetChild(boxWndWPC)

	// добавляем группу в основную коробку окна
	kitWindowWallPostCommentSettings.Box.Append(groupWndWPC, false)

	// получаем набор для кнопок принятия и отмены изменений
	kitButtonsWPC := makeSettingButtonsKit()

	// привязываем к кнопкам соответствующие процедуры
	kitButtonsWPC.ButtonCancel.OnClicked(func(*ui.Button) {
		// TODO: как-нибудь надо закрывать окно
	})
	// привязываем кнопки к соответствующим процедурам
	kitButtonsWPC.ButtonApply.OnClicked(func(*ui.Button) {
		var updatedWallPostCommentMonitorParam WallPostCommentMonitorParam

		updatedWallPostCommentMonitorParam.ID = wallPostCommentMonitorParam.ID
		updatedWallPostCommentMonitorParam.SubjectID = wallPostCommentMonitorParam.SubjectID
		if kitWndWPCMonitoring.CheckBox.Checked() {
			updatedWallPostCommentMonitorParam.NeedMonitoring = 1
		} else {
			updatedWallPostCommentMonitorParam.NeedMonitoring = 0
		}
		updatedWallPostCommentMonitorParam.PostsCount = kitWndWPCPostsCount.Spinbox.Value()
		updatedWallPostCommentMonitorParam.CommentsCount = kitWndWPCCommentsCount.Spinbox.Value()
		if kitWndWPCMonitoringAll.CheckBox.Checked() {
			updatedWallPostCommentMonitorParam.MonitoringAll = 1
		} else {
			updatedWallPostCommentMonitorParam.MonitoringAll = 0
		}
		jsonDump := fmt.Sprintf("{\"list\":[%v]}", kitWndWPCUsersIdsForMonitoring.Entry.Text())
		updatedWallPostCommentMonitorParam.UsersIDsForMonitoring = jsonDump
		jsonDump = fmt.Sprintf("{\"list\":[%v]}", kitWndWPCUsersNamesForMonitoring.Entry.Text())
		updatedWallPostCommentMonitorParam.UsersNamesForMonitoring = jsonDump
		updatedWallPostCommentMonitorParam.AttachmentsTypesForMonitoring = wallPostCommentMonitorParam.AttachmentsTypesForMonitoring
		jsonDump = fmt.Sprintf("{\"list\":[%v]}", kitWndWPCUsersIdsForIgnore.Entry.Text())
		updatedWallPostCommentMonitorParam.UsersIDsForIgnore = jsonDump
		updatedWallPostCommentMonitorParam.PostsCount = kitWndWPCInterval.Spinbox.Value()
		updatedWallPostCommentMonitorParam.SendTo, err = strconv.Atoi(kitWndWPCSendTo.Entry.Text())
		if err != nil {
			date := UnixTimeStampToDate(int(time.Now().Unix()))
			log.Fatal(fmt.Errorf("> [%v] WARNING! Error: %v", date, err))
		}
		listPostsFilters := []string{"all", "others", "owner"}
		updatedWallPostCommentMonitorParam.Filter = listPostsFilters[kitWndWPCFilter.Combobox.Selected()]
		updatedWallPostCommentMonitorParam.LastDate = wallPostCommentMonitorParam.LastDate
		jsonDump = fmt.Sprintf("{\"list\":[%v]}", kitWndWPCKeywordsForMonitoring.Entry.Text())
		updatedWallPostCommentMonitorParam.KeywordsForMonitoring = jsonDump
		jsonDump = fmt.Sprintf("{\"list\":[%v]}", kitWndWPCSmallCommentsForMonitoring.Entry.Text())
		updatedWallPostCommentMonitorParam.SmallCommentsForMonitoring = jsonDump
		updatedWallPostCommentMonitorParam.DigitsCountForCardNumberMonitoring = wallPostCommentMonitorParam.DigitsCountForCardNumberMonitoring
		updatedWallPostCommentMonitorParam.DigitsCountForPhoneNumberMonitoring = wallPostCommentMonitorParam.DigitsCountForPhoneNumberMonitoring
		if kitWndWPCMonitorByCommunity.CheckBox.Checked() {
			updatedWallPostCommentMonitorParam.MonitorByCommunity = 1
		} else {
			updatedWallPostCommentMonitorParam.MonitorByCommunity = 0
		}

		err = UpdateDBWallPostCommentMonitor(updatedWallPostCommentMonitorParam)
		if err != nil {
			date := UnixTimeStampToDate(int(time.Now().Unix()))
			log.Fatal(fmt.Errorf("> [%v] WARNING! Error: %v", date, err))
		}
		// TODO: как-нибудь надо закрывать окно
	})

	// добавляем коробку с кнопками на основную коробку окна
	kitWindowWallPostCommentSettings.Box.Append(kitButtonsWPC.Box, true)
	// затем еще одну коробку, для выравнивания расположения кнопок при растягивании окна
	boxWndWPCBottom := ui.NewHorizontalBox()
	kitWindowWallPostCommentSettings.Box.Append(boxWndWPCBottom, false)

	kitWindowWallPostCommentSettings.Window.Show()
}

// WindowSettingsKit хранит ссылки на объекты окна с установками модулей мониторинга
type WindowSettingsKit struct {
	Window *ui.Window
	Box    *ui.Box
}

func makeSettingWindowKit(windowTitle string, width, height int) WindowSettingsKit {
	var windowSettingsKit WindowSettingsKit

	windowSettingsKit.Window = ui.NewWindow(windowTitle, width, height, true)
	windowSettingsKit.Window.OnClosing(func(*ui.Window) bool {
		windowSettingsKit.Window.Disable()
		return true
	})
	windowSettingsKit.Window.SetMargined(true)
	windowSettingsKit.Box = ui.NewVerticalBox()
	windowSettingsKit.Box.SetPadded(true)
	windowSettingsKit.Window.SetChild(windowSettingsKit.Box)

	return windowSettingsKit
}

// CheckboxKit хранит ссылки на объекты для параметров с переключателями
type CheckboxKit struct {
	Box      *ui.Box
	CheckBox *ui.Checkbox
}

func makeSettingCheckboxKit(labelTitle string, needMonitoringFlag int) CheckboxKit {
	var checkboxKit CheckboxKit

	checkboxKit.Box = ui.NewHorizontalBox()
	checkboxKit.Box.SetPadded(true)
	labelObj := ui.NewLabel(labelTitle)
	checkboxKit.Box.Append(labelObj, true)
	checkboxKit.CheckBox = ui.NewCheckbox("")
	if needMonitoringFlag == 1 {
		checkboxKit.CheckBox.SetChecked(true)
	} else {
		checkboxKit.CheckBox.SetChecked(false)
	}
	checkboxKit.Box.Append(checkboxKit.CheckBox, true)

	return checkboxKit
}

// SpinboxKit хранит ссылки на объекты для параметров с спинбоксом
type SpinboxKit struct {
	Box     *ui.Box
	Spinbox *ui.Spinbox
}

func makeSettingSpinboxKit(labelTitle string, minValue, maxValue, currentValue int) SpinboxKit {
	var spinboxKit SpinboxKit

	spinboxKit.Box = ui.NewHorizontalBox()
	spinboxKit.Box.SetPadded(true)
	labelObj := ui.NewLabel(labelTitle)
	spinboxKit.Box.Append(labelObj, true)
	spinboxKit.Spinbox = ui.NewSpinbox(minValue, maxValue)
	spinboxKit.Spinbox.SetValue(currentValue)
	spinboxKit.Box.Append(spinboxKit.Spinbox, true)

	return spinboxKit
}

// EntryKit хранит ссылки на объекты для параметров с полями для ввода текста
type EntryKit struct {
	Box   *ui.Box
	Entry *ui.Entry
}

func makeSettingEntryKit(labelTitle string, entryValue string) EntryKit {
	var entryKit EntryKit

	entryKit.Box = ui.NewHorizontalBox()
	entryKit.Box.SetPadded(true)
	labelObj := ui.NewLabel(labelTitle)
	entryKit.Box.Append(labelObj, true)
	entryKit.Entry = ui.NewEntry()
	entryKit.Entry.SetText(entryValue)
	entryKit.Box.Append(entryKit.Entry, true)

	return entryKit
}

// EntryListKit хранит ссылки на объекты для параметров со списком в поле для ввода текста
type EntryListKit struct {
	Box   *ui.Box
	Entry *ui.Entry
}

func makeSettingEntryListKit(labelTitle, jsonDump string) EntryListKit {
	var entryListKit EntryListKit

	entryListKit.Box = ui.NewHorizontalBox()
	entryListKit.Box.SetPadded(true)
	labelObj := ui.NewLabel(labelTitle)
	entryListKit.Box.Append(labelObj, true)
	entryListKit.Entry = ui.NewEntry()
	structFromDump, err := MakeParamList(jsonDump)
	if err != nil {
		date := UnixTimeStampToDate(int(time.Now().Unix()))
		log.Fatal(fmt.Errorf("> [%v] WARNING! Error: %v", date, err))
	}
	if len(structFromDump.List) > 0 {
		var list string
		for i, item := range structFromDump.List {
			if i > 0 {
				list += ", "
			}
			list += fmt.Sprintf("\"%v\"", item)
		}
		entryListKit.Entry.SetText(list)
	}
	entryListKit.Box.Append(entryListKit.Entry, true)

	return entryListKit
}

// ComboboxKit хранит ссылки на объекты для параметров с выпадающим списком
type ComboboxKit struct {
	Box      *ui.Box
	Combobox *ui.Combobox
}

func makeSettingComboboxKit(labelTitle string, comboboxValues []string, currentValue string) ComboboxKit {
	var comboboxKit ComboboxKit

	comboboxKit.Box = ui.NewHorizontalBox()
	comboboxKit.Box.SetPadded(true)
	labelObj := ui.NewLabel(labelTitle)
	comboboxKit.Box.Append(labelObj, true)
	comboboxKit.Combobox = ui.NewCombobox()
	var slctd int
	for i, item := range comboboxValues {
		comboboxKit.Combobox.Append(item)
		if currentValue == item {
			slctd = i
		}
	}
	comboboxKit.Combobox.SetSelected(slctd)
	comboboxKit.Box.Append(comboboxKit.Combobox, true)

	return comboboxKit
}

// ButtonsKit хранит ссылки на объекты для кнопок принятия и отмены изменений в установках
type ButtonsKit struct {
	Box          *ui.Box
	ButtonApply  *ui.Button
	ButtonCancel *ui.Button
}

func makeSettingButtonsKit() ButtonsKit {
	var buttonsKit ButtonsKit

	buttonsKit.Box = ui.NewHorizontalBox()
	buttonsKit.Box.SetPadded(true)

	boxButtons := ui.NewHorizontalBox()
	boxButtons.SetPadded(true)
	buttonsKit.ButtonCancel = ui.NewButton("Cancel")
	boxButtons.Append(buttonsKit.ButtonCancel, false)
	buttonsKit.ButtonApply = ui.NewButton("Apply")
	boxButtons.Append(buttonsKit.ButtonApply, false)

	// для выравнивания кнопок
	boxEmptyLeft := ui.NewHorizontalBox()
	boxEmptyCenter := ui.NewHorizontalBox()

	buttonsKit.Box.Append(boxEmptyLeft, false)
	buttonsKit.Box.Append(boxEmptyCenter, false)
	buttonsKit.Box.Append(boxButtons, false)

	return buttonsKit
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
