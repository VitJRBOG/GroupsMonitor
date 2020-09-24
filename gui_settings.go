package main

import (
	"fmt"
	"log"
	"runtime/debug"
	"strconv"
	"time"

	"github.com/andlabs/ui"
)

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

func makeAccessTokensSettingsBox() *ui.Box {
	boxAccessTokensSettings := ui.NewVerticalBox()

	// запрашиваем список токенов доступа из базы данных
	accessTokens, err := SelectDBAccessTokens()
	if err != nil {
		ToLogFile(err.Error(), string(debug.Stack()))
		panic(err.Error())
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
	kitWindowAccessTokenAddition := MakeSettingWindowKit(windowTitle, 300, 100)

	boxWndATAddition := ui.NewVerticalBox()

	// получаем набор для ввода названия нового токена доступа
	kitATCreationName := MakeSettingEntryKit("Name", "")

	// получаем набор для ввода значения нового токена доступа
	kitATCreationValue := MakeSettingEntryKit("Value", "")

	// описываем группу, в которой будут размещены элементы
	groupWndATAddition := ui.NewGroup("")
	groupWndATAddition.SetMargined(true)
	boxWndATAddition.Append(kitATCreationName.Box, false)
	boxWndATAddition.Append(kitATCreationValue.Box, false)
	groupWndATAddition.SetChild(boxWndATAddition)

	// добавляем группу в основную коробку окна
	kitWindowAccessTokenAddition.Box.Append(groupWndATAddition, false)

	// получаем набор для кнопок принятия и отмены изменений
	kitButtonsATAddition := MakeSettingButtonsKit()

	// привязываем кнопки к соответствующим процедурам
	kitButtonsATAddition.ButtonCancel.OnClicked(func(*ui.Button) {
		kitWindowAccessTokenAddition.Window.Disable()
		kitWindowAccessTokenAddition.Window.Hide()
	})
	kitButtonsATAddition.ButtonApply.OnClicked(func(*ui.Button) {
		var accessToken AccessToken

		accessToken.Name = kitATCreationName.Entry.Text()
		accessToken.Value = kitATCreationValue.Entry.Text()

		err := InsertDBAccessToken(accessToken)
		if err != nil {
			ToLogFile(err.Error(), string(debug.Stack()))
			panic(err.Error())
		}

		kitWindowAccessTokenAddition.Window.Disable()
		kitWindowAccessTokenAddition.Window.Hide()
	})

	// добавляем коробку с кнопками на основную коробку окна
	kitWindowAccessTokenAddition.Box.Append(kitButtonsATAddition.Box, true)
	// затем еще одну коробку, для выравнивания расположения кнопок при растягивании окна
	boxWndATAdditionBottom := ui.NewHorizontalBox()
	kitWindowAccessTokenAddition.Box.Append(boxWndATAdditionBottom, false)

	kitWindowAccessTokenAddition.Window.Show()
}

func showAccessTokenSettingWindow(IDAccessToken int) {
	// запрашиваем список токенов доступа из базы данных
	accessTokens, err := SelectDBAccessTokens()
	if err != nil {
		ToLogFile(err.Error(), string(debug.Stack()))
		panic(err.Error())
	}

	// получаем набор для окна отображения установок токена доступа
	kitWindowAccessTokenSettings := MakeSettingWindowKit("", 300, 100)

	// перечисляем токены доступа
	for _, accessToken := range accessTokens {
		// и ищем токен с подходящим идентификатором
		if accessToken.ID == IDAccessToken {
			// устанавливаем заголовок окна в соответствии с названием токена доступа
			windowTitle := "Settings of " + accessToken.Name + "'s access token"
			kitWindowAccessTokenSettings.Window.SetTitle(windowTitle)

			boxWndAT := ui.NewVerticalBox()

			// получаем набор для названия токена доступа
			kitWndATName := MakeSettingEntryKit("Name", accessToken.Name)

			// получаем набор для значения токена доступа
			kitWndATValue := MakeSettingEntryKit("Value", accessToken.Value)

			// описываем группу, в которой будут размещены элементы
			groupWndAT := ui.NewGroup("")
			groupWndAT.SetMargined(true)
			boxWndAT.Append(kitWndATName.Box, false)
			boxWndAT.Append(kitWndATValue.Box, false)
			groupWndAT.SetChild(boxWndAT)

			// добавляем группу в основную коробку окна
			kitWindowAccessTokenSettings.Box.Append(groupWndAT, false)

			// получаем набор для кнопок принятия и отмены изменений
			kitButtonsAT := MakeSettingButtonsKit()

			kitButtonsAT.ButtonCancel.OnClicked(func(*ui.Button) {
				kitWindowAccessTokenSettings.Window.Disable()
				kitWindowAccessTokenSettings.Window.Hide()
			})
			// привязываем кнопки к соответствующим процедурам
			kitButtonsAT.ButtonApply.OnClicked(func(*ui.Button) {
				var updatedAccessToken AccessToken
				updatedAccessToken.ID = accessToken.ID
				updatedAccessToken.Name = kitWndATName.Entry.Text()
				updatedAccessToken.Value = kitWndATValue.Entry.Text()

				err := UpdateDBAccessToken(updatedAccessToken)
				if err != nil {
					ToLogFile(err.Error(), string(debug.Stack()))
					panic(err.Error())
				}

				kitWindowAccessTokenSettings.Window.Disable()
				kitWindowAccessTokenSettings.Window.Hide()
			})

			// добавляем коробку с кнопками на основную коробку окна
			kitWindowAccessTokenSettings.Box.Append(kitButtonsAT.Box, true)
			// затем еще одну коробку, для выравнивания расположения кнопок при растягивании окна
			boxWndATBottom := ui.NewHorizontalBox()
			kitWindowAccessTokenSettings.Box.Append(boxWndATBottom, false)
			break
		}
	}

	kitWindowAccessTokenSettings.Window.Show()
}

func makeSubjectsSettingsBox(groupsSettingsData GroupsSettingsData) *ui.Box {
	// описываем коробку для установок субъектов
	boxSubjectsSettings := ui.NewVerticalBox()

	// запрашиваем список субъектов из базы данных
	subjects, err := SelectDBSubjects()
	if err != nil {
		ToLogFile(err.Error(), string(debug.Stack()))
		panic(err.Error())
	}

	boxSUpper := ui.NewVerticalBox()

	// в этом списке будут храниться ссылки на кнопки для отображения доп. настроек
	var listBtnsSubjectSettings []*ui.Button

	// перечисляем субъекты
	for _, subjectData := range subjects {
		// описываем кнопку для отображения доп. настроек соответствующего субъекта
		btnSubjectSettings := ui.NewButton(subjectData.Name)
		// и добавляем ее в коробку
		boxSUpper.Append(btnSubjectSettings, false)
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

	// описываем коробку для кнопки добавления нового токена доступа
	boxSBottom := ui.NewHorizontalBox()

	// описываем кнопку для добавления нового токена и декоративную кнопку для выравнивания
	btnAddNewSubject := ui.NewButton("＋")
	btnDecorative := ui.NewButton("")
	btnDecorative.Disable()

	// привязываем к кнопке для добавления соответствующую процедуру
	btnAddNewSubject.OnClicked(func(*ui.Button) {
		showSubjectAdditionWindow()
	})

	// добавляем кнопки на коробку для кнопок
	boxSBottom.Append(btnAddNewSubject, false)
	boxSBottom.Append(btnDecorative, true)

	// добавляем обе коробки на основную коробку
	boxSubjectsSettings.Append(boxSUpper, false)
	boxSubjectsSettings.Append(boxSBottom, false)

	return boxSubjectsSettings
}

// NewSubjectData хранит данные для создаваемого субъекта
type NewSubjectData struct {
	ID   int
	Name string
}

// NewMonitorData хранит данные для создаваемого монитора
type NewMonitorData struct {
	Name      string
	SubjectID int
}

// NewMethodData хранит данные для создаваемого метода
type NewMethodData struct {
	Name          string
	SubjectID     int
	AccessTokenID int
	MonitorID     int
}

// ListNewMethodData хранит списки со структурами с данными для создаваемых методов
type ListNewMethodData struct {
	WPM  []NewMethodData // wall_post_monitor
	APM  []NewMethodData // album_photo_monitor
	VM   []NewMethodData // video_monitor
	PCM  []NewMethodData // album_photo_comment_monitor
	VCM  []NewMethodData // video_monitor
	TM   []NewMethodData // topic_monitor
	WPCM []NewMethodData // wall_post_comment_monitor
}

// NewMonitorModuleData хранит данные для создаваемого модуля мониторинга
type NewMonitorModuleData struct {
	SubjectID int
	SendTo    int
}

// ListNewMonitorModuleData хранит структуры с данными для создаваемых модулей мониторинга
type ListNewMonitorModuleData struct {
	WPM  NewMonitorModuleData // wall_post_monitor
	APM  NewMonitorModuleData // album_photo_monitor
	VM   NewMonitorModuleData // video_monitor
	PCM  NewMonitorModuleData // album_photo_comment_monitor
	VCM  NewMonitorModuleData // video_monitor
	TM   NewMonitorModuleData // topic_monitor
	WPCM NewMonitorModuleData // wall_post_comment_monitor
}

func showSubjectAdditionWindow() {
	// запрашиваем список токенов доступа из базы данных
	accessTokens, err := SelectDBAccessTokens()
	if err != nil {
		ToLogFile(err.Error(), string(debug.Stack()))
		panic(err.Error())
	}
	// формируем список с названиями токенов доступа
	var accessTokensNames []string
	for _, accessToken := range accessTokens {
		accessTokensNames = append(accessTokensNames, accessToken.Name)
	}

	// получаем набор для отображения окна для добавления нового субъекта мониторинга
	windowTitle := fmt.Sprintf("New subject addition")
	kitWindowSubjectAddition := MakeSettingWindowKit(windowTitle, 300, 100)

	boxWndSAddition := ui.NewHorizontalBox()
	boxWndSAddition.SetPadded(true)

	boxWndSAdditionLeft := ui.NewVerticalBox()
	boxWndSAdditionLeft.SetPadded(true)
	boxWndSAdditionCenter := ui.NewVerticalBox()
	boxWndSAdditionCenter.SetPadded(true)
	boxWndSAdditionRight := ui.NewVerticalBox()
	boxWndSAdditionRight.SetPadded(true)

	// получаем набор для ввода названия нового субъекта
	kitSAdditionName := MakeSettingEntryKit("Name", "")

	// получаем набор для ввода идентификатора субъекта в базе данных ВК
	kitSAdditionSubjectID := MakeSettingEntryKit("Subject ID", "")
	// и привязываем его к обработке ввода, чтобы кроме чисел ничего не вводилось
	kitSAdditionSubjectID.Entry.OnChanged(func(*ui.Entry) {
		NumericEntriesHandler(kitSAdditionSubjectID.Entry)
	})

	// описываем группу для общих установок субъекта
	groupSAdditionGeneral := ui.NewGroup("General")
	groupSAdditionGeneral.SetMargined(true)
	boxSAdditionGeneral := ui.NewVerticalBox()
	boxSAdditionGeneral.Append(kitSAdditionName.Box, false)
	boxSAdditionGeneral.Append(kitSAdditionSubjectID.Box, false)
	groupSAdditionGeneral.SetChild(boxSAdditionGeneral)

	// получаем набор для ввода идентификатора получателя сообщений в модуле wall_post_monitor
	kitSAdditionSendToinWPM := MakeSettingEntryKit("Send to", "")
	// и привязываем его к обработке ввода, чтобы кроме чисел ничего не вводилось
	kitSAdditionSendToinWPM.Entry.OnChanged(func(*ui.Entry) {
		NumericEntriesHandler(kitSAdditionSendToinWPM.Entry)
	})

	// получаем набор для выбора токена доступа для метода wall.get в модуле wall_post_monitor
	kitSAdditionWGinWPM := MakeSettingComboboxKit("Access token for \"wall.get\"", accessTokensNames, "")

	// получаем набор для выбора токена доступа для метода users.get в модуле wall_post_monitor
	kitSAdditionUGinWPM := MakeSettingComboboxKit("Access token for \"users.get\"", accessTokensNames, "")

	// получаем набор для выбора токена доступа для метода groups.getById в модуле wall_post_monitor
	kitSAdditionGGBIinWPM := MakeSettingComboboxKit("Access token for \"groups.getById\"", accessTokensNames, "")

	// получаем набор для выбора токена доступа для метода messages.send в модуле wall_post_monitor
	kitSAdditionMSinWPM := MakeSettingComboboxKit("Access token for \"messages.send\"", accessTokensNames, "")

	// описываем группу для установок модуля wall_post_monitor субъекта
	groupSAdditionWPM := ui.NewGroup("Wall post monitor")
	groupSAdditionWPM.SetMargined(true)
	boxSAdditionWPM := ui.NewVerticalBox()
	boxSAdditionWPM.Append(kitSAdditionSendToinWPM.Box, false)
	boxSAdditionWPM.Append(kitSAdditionWGinWPM.Box, false)
	boxSAdditionWPM.Append(kitSAdditionUGinWPM.Box, false)
	boxSAdditionWPM.Append(kitSAdditionGGBIinWPM.Box, false)
	boxSAdditionWPM.Append(kitSAdditionMSinWPM.Box, false)
	groupSAdditionWPM.SetChild(boxSAdditionWPM)

	// получаем набор для ввода идентификатора получателя сообщений в модуле album_photo_monitor
	kitSAdditionSendToinAPM := MakeSettingEntryKit("Send to", "")
	// и привязываем его к обработке ввода, чтобы кроме чисел ничего не вводилось
	kitSAdditionSendToinAPM.Entry.OnChanged(func(*ui.Entry) {
		NumericEntriesHandler(kitSAdditionSendToinAPM.Entry)
	})

	// получаем набор для выбора токена доступа для метода photos.get в модуле album_photo_monitor
	kitSAdditionPGinAPM := MakeSettingComboboxKit("Access token for \"photos.get\"", accessTokensNames, "")

	// получаем набор для выбора токена доступа для метода photos.getAlbums в модуле album_photo_monitor
	kitSAdditionPGAinAPM := MakeSettingComboboxKit("Access token for \"photos.getAlbums\"", accessTokensNames, "")

	// получаем набор для выбора токена доступа для метода users.get в модуле album_photo_monitor
	kitSAdditionUGinAPM := MakeSettingComboboxKit("Access token for \"users.get\"", accessTokensNames, "")

	// получаем набор для выбора токена доступа для метода groups.getById в модуле album_photo_monitor
	kitSAdditionGGBIinAPM := MakeSettingComboboxKit("Access token for \"groups.getById\"", accessTokensNames, "")

	// получаем набор для выбора токена доступа для метода messages.send в модуле album_photo_monitor
	kitSAdditionMSinAPM := MakeSettingComboboxKit("Access token for \"messages.send\"", accessTokensNames, "")

	// описываем группу для установок модуля album_photo_monitor субъекта
	groupSAdditionAPM := ui.NewGroup("Album photo monitor")
	groupSAdditionAPM.SetMargined(true)
	boxSAdditionAPM := ui.NewVerticalBox()
	boxSAdditionAPM.Append(kitSAdditionSendToinAPM.Box, false)
	boxSAdditionAPM.Append(kitSAdditionPGinAPM.Box, false)
	boxSAdditionAPM.Append(kitSAdditionPGAinAPM.Box, false)
	boxSAdditionAPM.Append(kitSAdditionUGinAPM.Box, false)
	boxSAdditionAPM.Append(kitSAdditionGGBIinAPM.Box, false)
	boxSAdditionAPM.Append(kitSAdditionMSinAPM.Box, false)
	groupSAdditionAPM.SetChild(boxSAdditionAPM)

	// получаем набор для ввода идентификатора получателя сообщений в модуле video_monitor
	kitSAdditionSendToinVM := MakeSettingEntryKit("Send to", "")
	// и привязываем его к обработке ввода, чтобы кроме чисел ничего не вводилось
	kitSAdditionSendToinVM.Entry.OnChanged(func(*ui.Entry) {
		NumericEntriesHandler(kitSAdditionSendToinVM.Entry)
	})

	// получаем набор для выбора токена доступа для метода video.get в модуле video_monitor
	kitSAdditionVGinVM := MakeSettingComboboxKit("Access token for \"video.get\"", accessTokensNames, "")

	// получаем набор для выбора токена доступа для метода users.get в модуле video_monitor
	kitSAdditionUGinVM := MakeSettingComboboxKit("Access token for \"users.get\"", accessTokensNames, "")

	// получаем набор для выбора токена доступа для метода groups.getById в модуле video_monitor
	kitSAdditionGGBIinVM := MakeSettingComboboxKit("Access token for \"groups.getById\"", accessTokensNames, "")

	// получаем набор для выбора токена доступа для метода messages.send в модуле video_monitor
	kitSAdditionMSinVM := MakeSettingComboboxKit("Access token for \"messages.send\"", accessTokensNames, "")

	// описываем группу для установок модуля video_monitor субъекта
	groupSAdditionVM := ui.NewGroup("Video monitor")
	groupSAdditionVM.SetMargined(true)
	boxSAdditionVM := ui.NewVerticalBox()
	boxSAdditionVM.Append(kitSAdditionSendToinVM.Box, false)
	boxSAdditionVM.Append(kitSAdditionVGinVM.Box, false)
	boxSAdditionVM.Append(kitSAdditionUGinVM.Box, false)
	boxSAdditionVM.Append(kitSAdditionGGBIinVM.Box, false)
	boxSAdditionVM.Append(kitSAdditionMSinVM.Box, false)
	groupSAdditionVM.SetChild(boxSAdditionVM)

	// получаем набор для ввода идентификатора получателя сообщений в модуле photo_comment_monitor
	kitSAdditionSendToinPCM := MakeSettingEntryKit("Send to", "")
	// и привязываем его к обработке ввода, чтобы кроме чисел ничего не вводилось
	kitSAdditionSendToinPCM.Entry.OnChanged(func(*ui.Entry) {
		NumericEntriesHandler(kitSAdditionSendToinPCM.Entry)
	})

	// получаем набор для выбора токена доступа для метода photos.getAllComments в модуле photo_comment_monitor
	kitSAdditionPGACinPCM := MakeSettingComboboxKit("Access token for \"photos.getAllComments\"", accessTokensNames, "")

	// получаем набор для выбора токена доступа для метода users.get в модуле photo_comment_monitor
	kitSAdditionUGinPCM := MakeSettingComboboxKit("Access token for \"users.get\"", accessTokensNames, "")

	// получаем набор для выбора токена доступа для метода groups.getById в модуле photo_comment_monitor
	kitSAdditionGGBIinPCM := MakeSettingComboboxKit("Access token for \"groups.getById\"", accessTokensNames, "")

	// получаем набор для выбора токена доступа для метода messages.send в модуле photo_comment_monitor
	kitSAdditionMSinPCM := MakeSettingComboboxKit("Access token for \"messages.send\"", accessTokensNames, "")

	// описываем группу для установок модуля photo_comment_monitor субъекта
	groupSAdditionPCM := ui.NewGroup("Photo comment monitor")
	groupSAdditionPCM.SetMargined(true)
	boxSAdditionPCM := ui.NewVerticalBox()
	boxSAdditionPCM.Append(kitSAdditionSendToinPCM.Box, false)
	boxSAdditionPCM.Append(kitSAdditionPGACinPCM.Box, false)
	boxSAdditionPCM.Append(kitSAdditionUGinPCM.Box, false)
	boxSAdditionPCM.Append(kitSAdditionGGBIinPCM.Box, false)
	boxSAdditionPCM.Append(kitSAdditionMSinPCM.Box, false)
	groupSAdditionPCM.SetChild(boxSAdditionPCM)

	// получаем набор для ввода идентификатора получателя сообщений в модуле video_comment_monitor
	kitSAdditionSendToinVCM := MakeSettingEntryKit("Send to", "")
	// и привязываем его к обработке ввода, чтобы кроме чисел ничего не вводилось
	kitSAdditionSendToinVCM.Entry.OnChanged(func(*ui.Entry) {
		NumericEntriesHandler(kitSAdditionSendToinVCM.Entry)
	})

	// получаем набор для выбора токена доступа для метода video.getComments в модуле video_comment_monitor
	kitSAdditionVGCinVCM := MakeSettingComboboxKit("Access token for \"video.getComments\"", accessTokensNames, "")

	// получаем набор для выбора токена доступа для метода users.get в модуле video_comment_monitor
	kitSAdditionUGinVCM := MakeSettingComboboxKit("Access token for \"users.get\"", accessTokensNames, "")

	// получаем набор для выбора токена доступа для метода groups.getById в модуле video_comment_monitor
	kitSAdditionGGBIinVCM := MakeSettingComboboxKit("Access token for \"groups.getById\"", accessTokensNames, "")

	// получаем набор для выбора токена доступа для метода video.get в модуле video_comment_monitor
	kitSAdditionVGinVCM := MakeSettingComboboxKit("Access token for \"video.get\"", accessTokensNames, "")

	// получаем набор для выбора токена доступа для метода messages.send в модуле video_comment_monitor
	kitSAdditionMSinVCM := MakeSettingComboboxKit("Access token for \"messages.send\"", accessTokensNames, "")

	// описываем группу для установок модуля video_comment_monitor субъекта
	groupSAdditionVCM := ui.NewGroup("Video comment monitor")
	groupSAdditionVCM.SetMargined(true)
	boxSAdditionVCM := ui.NewVerticalBox()
	boxSAdditionVCM.Append(kitSAdditionSendToinVCM.Box, false)
	boxSAdditionVCM.Append(kitSAdditionVGCinVCM.Box, false)
	boxSAdditionVCM.Append(kitSAdditionUGinVCM.Box, false)
	boxSAdditionVCM.Append(kitSAdditionGGBIinVCM.Box, false)
	boxSAdditionVCM.Append(kitSAdditionVGinVCM.Box, false)
	boxSAdditionVCM.Append(kitSAdditionMSinVCM.Box, false)
	groupSAdditionVCM.SetChild(boxSAdditionVCM)

	// получаем набор для ввода идентификатора получателя сообщений в модуле topic_monitor
	kitSAdditionSendToinTM := MakeSettingEntryKit("Send to", "")
	// и привязываем его к обработке ввода, чтобы кроме чисел ничего не вводилось
	kitSAdditionSendToinTM.Entry.OnChanged(func(*ui.Entry) {
		NumericEntriesHandler(kitSAdditionSendToinTM.Entry)
	})

	// получаем набор для выбора токена доступа для метода board.getComments в модуле topic_monitor
	kitSAdditionBGCinTM := MakeSettingComboboxKit("Access token for \"board.getComments\"", accessTokensNames, "")

	// получаем набор для выбора токена доступа для метода board.getTopics в модуле topic_monitor
	kitSAdditionBGTinTM := MakeSettingComboboxKit("Access token for \"board.getTopics\"", accessTokensNames, "")

	// получаем набор для выбора токена доступа для метода users.get в модуле topic_monitor
	kitSAdditionUGinTM := MakeSettingComboboxKit("Access token for \"users.get\"", accessTokensNames, "")

	// получаем набор для выбора токена доступа для метода groups.getById в модуле topic_monitor
	kitSAdditionGGBIinTM := MakeSettingComboboxKit("Access token for \"groups.getById\"", accessTokensNames, "")

	// получаем набор для выбора токена доступа для метода messages.send в модуле topic_monitor
	kitSAdditionMSinTM := MakeSettingComboboxKit("Access token for \"messages.send\"", accessTokensNames, "")

	// описываем группу для установок модуля topic_monitor субъекта
	groupSAdditionTM := ui.NewGroup("Topic monitor")
	groupSAdditionTM.SetMargined(true)
	boxSAdditionTM := ui.NewVerticalBox()
	boxSAdditionTM.Append(kitSAdditionSendToinTM.Box, false)
	boxSAdditionTM.Append(kitSAdditionBGCinTM.Box, false)
	boxSAdditionTM.Append(kitSAdditionBGTinTM.Box, false)
	boxSAdditionTM.Append(kitSAdditionUGinTM.Box, false)
	boxSAdditionTM.Append(kitSAdditionGGBIinTM.Box, false)
	boxSAdditionTM.Append(kitSAdditionMSinTM.Box, false)
	groupSAdditionTM.SetChild(boxSAdditionTM)

	// получаем набор для ввода идентификатора получателя сообщений в модуле wall_post_comment_monitor
	kitSAdditionSendToinWPCM := MakeSettingEntryKit("Send to", "")
	// и привязываем его к обработке ввода, чтобы кроме чисел ничего не вводилось
	kitSAdditionSendToinWPCM.Entry.OnChanged(func(*ui.Entry) {
		NumericEntriesHandler(kitSAdditionSendToinWPCM.Entry)
	})

	// получаем набор для выбора токена доступа для метода wall.getComments в модуле wall_post_comment_monitor
	kitSAdditionWGCsinWPCM := MakeSettingComboboxKit("Access token for \"wall.getComments\"", accessTokensNames, "")

	// получаем набор для выбора токена доступа для метода users.get в модуле wall_post_comment_monitor
	kitSAdditionUGinWPCM := MakeSettingComboboxKit("Access token for \"users.get\"", accessTokensNames, "")

	// получаем набор для выбора токена доступа для метода groups.getById в модуле wall_post_comment_monitor
	kitSAdditionGGBIinWPCM := MakeSettingComboboxKit("Access token for \"groups.getById\"", accessTokensNames, "")

	// получаем набор для выбора токена доступа для метода wall.get в модуле wall_post_comment_monitor
	kitSAdditionWGinWPCM := MakeSettingComboboxKit("Access token for \"wall.get\"", accessTokensNames, "")

	// получаем набор для выбора токена доступа для метода wall.getComment в модуле wall_post_comment_monitor
	kitSAdditionWGCinWPCM := MakeSettingComboboxKit("Access token for \"wall.getComment\"", accessTokensNames, "")

	// получаем набор для выбора токена доступа для метода messages.send в модуле wall_post_comment_monitor
	kitSAdditionMSinWPCM := MakeSettingComboboxKit("Access token for \"messages.send\"", accessTokensNames, "")

	// описываем группу для установок модуля wall_post_comment_monitor субъекта
	groupSAdditionWPCM := ui.NewGroup("Wall post comment monitor")
	groupSAdditionWPCM.SetMargined(true)
	boxSAdditionWPCM := ui.NewVerticalBox()
	boxSAdditionWPCM.Append(kitSAdditionSendToinWPCM.Box, false)
	boxSAdditionWPCM.Append(kitSAdditionWGCsinWPCM.Box, false)
	boxSAdditionWPCM.Append(kitSAdditionUGinWPCM.Box, false)
	boxSAdditionWPCM.Append(kitSAdditionGGBIinWPCM.Box, false)
	boxSAdditionWPCM.Append(kitSAdditionWGinWPCM.Box, false)
	boxSAdditionWPCM.Append(kitSAdditionWGCinWPCM.Box, false)
	boxSAdditionWPCM.Append(kitSAdditionMSinWPCM.Box, false)
	groupSAdditionWPCM.SetChild(boxSAdditionWPCM)

	// добавляем все заполненные группы на
	// левую коробку
	boxWndSAdditionLeft.Append(groupSAdditionGeneral, false)
	boxWndSAdditionLeft.Append(groupSAdditionWPM, false)
	boxWndSAdditionLeft.Append(groupSAdditionAPM, false)

	// коробку посередине
	boxWndSAdditionCenter.Append(groupSAdditionVM, false)
	boxWndSAdditionCenter.Append(groupSAdditionPCM, false)
	boxWndSAdditionCenter.Append(groupSAdditionVCM, false)

	// и правую коробку (для экономии места из-за отсутствия в данной библиотеке для gui скроллинга)
	boxWndSAdditionRight.Append(groupSAdditionTM, false)
	boxWndSAdditionRight.Append(groupSAdditionWPCM, false)

	// затем добавляем левую и правую коробки на одну общую
	boxWndSAddition.Append(boxWndSAdditionLeft, false)
	boxWndSAddition.Append(boxWndSAdditionCenter, false)
	boxWndSAddition.Append(boxWndSAdditionRight, false)

	// добавляем коробку в основную коробку окна
	kitWindowSubjectAddition.Box.Append(boxWndSAddition, false)

	// получаем набор для кнопок принятия и отмены изменений
	kitButtonsSAddition := MakeSettingButtonsKit()

	// привязываем кнопки к соответствующим процедурам
	kitButtonsSAddition.ButtonCancel.OnClicked(func(*ui.Button) {
		kitWindowSubjectAddition.Window.Disable()
		kitWindowSubjectAddition.Window.Hide()
	})
	kitButtonsSAddition.ButtonApply.OnClicked(func(*ui.Button) {
		var newSubjectData NewSubjectData

		if len(kitSAdditionSubjectID.Entry.Text()) == 0 {
			warningTitle := "Field \"Subject ID\" must not be empty."
			ShowWarningWindow(warningTitle)
			return
		}
		subjectID, err := strconv.Atoi(kitSAdditionSubjectID.Entry.Text())
		if err != nil {
			ToLogFile(err.Error(), string(debug.Stack()))
			panic(err.Error())
		}
		newSubjectData.ID = subjectID
		if len(kitSAdditionName.Entry.Text()) == 0 {
			warningTitle := "Field \"Name\" must not be empty."
			ShowWarningWindow(warningTitle)
			return
		}
		newSubjectData.Name = kitSAdditionName.Entry.Text()

		var listNewMethodData ListNewMethodData
		var listNewMonitorModuleData ListNewMonitorModuleData

		monitorsNames := []string{"wall_post_monitor", "album_photo_monitor", "video_monitor",
			"photo_comment_monitor", "video_comment_monitor", "topic_monitor", "wall_post_comment_monitor"}

		for _, monitorName := range monitorsNames {
			switch monitorName {
			case "wall_post_monitor":

				if kitSAdditionWGinWPM.Combobox.Selected() == -1 {
					warningTitle := "You must select an item in the combobox " +
						"\"Access token for \"wall.get\"\""
					ShowWarningWindow(warningTitle)
					return
				}
				accessTokenName := accessTokensNames[kitSAdditionWGinWPM.Combobox.Selected()]
				for _, accessToken := range accessTokens {
					if accessTokenName == accessToken.Name {
						var newMethodData NewMethodData
						newMethodData.AccessTokenID = accessToken.ID
						newMethodData.Name = "wall.get"
						listNewMethodData.WPM = append(listNewMethodData.WPM, newMethodData)
					}
				}

				if kitSAdditionUGinWPM.Combobox.Selected() == -1 {
					warningTitle := "You must select an item in the combobox " +
						"\"Access token for \"Users.get\"\""
					ShowWarningWindow(warningTitle)
					return
				}
				accessTokenName = accessTokensNames[kitSAdditionUGinWPM.Combobox.Selected()]
				for _, accessToken := range accessTokens {
					if accessTokenName == accessToken.Name {
						var newMethodData NewMethodData
						newMethodData.AccessTokenID = accessToken.ID
						newMethodData.Name = "users.get"
						listNewMethodData.WPM = append(listNewMethodData.WPM, newMethodData)
					}
				}

				if kitSAdditionGGBIinWPM.Combobox.Selected() == -1 {
					warningTitle := "You must select an item in the combobox " +
						"\"Access token for \"Groups.getById\"\""
					ShowWarningWindow(warningTitle)
					return
				}
				accessTokenName = accessTokensNames[kitSAdditionGGBIinWPM.Combobox.Selected()]
				for _, accessToken := range accessTokens {
					if accessTokenName == accessToken.Name {
						var newMethodData NewMethodData
						newMethodData.AccessTokenID = accessToken.ID
						newMethodData.Name = "groups.getById"
						listNewMethodData.WPM = append(listNewMethodData.WPM, newMethodData)
					}
				}

				if kitSAdditionMSinWPM.Combobox.Selected() == -1 {
					warningTitle := "You must select an item in the combobox " +
						"\"Access token for \"messages.send\"\""
					ShowWarningWindow(warningTitle)
					return
				}
				accessTokenName = accessTokensNames[kitSAdditionMSinWPM.Combobox.Selected()]
				for _, accessToken := range accessTokens {
					if accessTokenName == accessToken.Name {
						var newMethodData NewMethodData
						newMethodData.AccessTokenID = accessToken.ID
						newMethodData.Name = "messages.send"
						listNewMethodData.WPM = append(listNewMethodData.WPM, newMethodData)
					}
				}

				if len(kitSAdditionSendToinWPM.Entry.Text()) == 0 {
					warningTitle := "Field \"Send to\" must not be empty."
					ShowWarningWindow(warningTitle)
					return
				}
				sendTo, err := strconv.Atoi(kitSAdditionSendToinWPM.Entry.Text())
				if err != nil {
					ToLogFile(err.Error(), string(debug.Stack()))
					panic(err.Error())
				}
				var newMonitorModuleData NewMonitorModuleData
				newMonitorModuleData.SendTo = sendTo
				listNewMonitorModuleData.WPM = newMonitorModuleData

			case "album_photo_monitor":

				if kitSAdditionPGinAPM.Combobox.Selected() == -1 {
					warningTitle := "You must select an item in the combobox " +
						"\"Access token for \"photos.get\"\""
					ShowWarningWindow(warningTitle)
					return
				}
				accessTokenName := accessTokensNames[kitSAdditionPGinAPM.Combobox.Selected()]
				for _, accessToken := range accessTokens {
					if accessTokenName == accessToken.Name {
						var newMethodData NewMethodData
						newMethodData.AccessTokenID = accessToken.ID
						newMethodData.Name = "photos.get"
						listNewMethodData.APM = append(listNewMethodData.APM, newMethodData)
					}
				}

				if kitSAdditionPGAinAPM.Combobox.Selected() == -1 {
					warningTitle := "You must select an item in the combobox " +
						"\"Access token for \"photos.getAlbums\"\""
					ShowWarningWindow(warningTitle)
					return
				}
				accessTokenName = accessTokensNames[kitSAdditionPGAinAPM.Combobox.Selected()]
				for _, accessToken := range accessTokens {
					if accessTokenName == accessToken.Name {
						var newMethodData NewMethodData
						newMethodData.AccessTokenID = accessToken.ID
						newMethodData.Name = "photos.getAlbums"
						listNewMethodData.APM = append(listNewMethodData.APM, newMethodData)
					}
				}

				if kitSAdditionUGinAPM.Combobox.Selected() == -1 {
					warningTitle := "You must select an item in the combobox " +
						"\"Access token for \"users.get\"\""
					ShowWarningWindow(warningTitle)
					return
				}
				accessTokenName = accessTokensNames[kitSAdditionUGinAPM.Combobox.Selected()]
				for _, accessToken := range accessTokens {
					if accessTokenName == accessToken.Name {
						var newMethodData NewMethodData
						newMethodData.AccessTokenID = accessToken.ID
						newMethodData.Name = "users.get"
						listNewMethodData.APM = append(listNewMethodData.APM, newMethodData)
					}
				}

				if kitSAdditionGGBIinAPM.Combobox.Selected() == -1 {
					warningTitle := "You must select an item in the combobox " +
						"\"Access token for \"groups.getById\"\""
					ShowWarningWindow(warningTitle)
					return
				}
				accessTokenName = accessTokensNames[kitSAdditionGGBIinAPM.Combobox.Selected()]
				for _, accessToken := range accessTokens {
					if accessTokenName == accessToken.Name {
						var newMethodData NewMethodData
						newMethodData.AccessTokenID = accessToken.ID
						newMethodData.Name = "groups.getById"
						listNewMethodData.APM = append(listNewMethodData.APM, newMethodData)
					}
				}

				if kitSAdditionMSinAPM.Combobox.Selected() == -1 {
					warningTitle := "You must select an item in the combobox " +
						"\"Access token for \"messages.send\"\""
					ShowWarningWindow(warningTitle)
					return
				}
				accessTokenName = accessTokensNames[kitSAdditionMSinAPM.Combobox.Selected()]
				for _, accessToken := range accessTokens {
					if accessTokenName == accessToken.Name {
						var newMethodData NewMethodData
						newMethodData.AccessTokenID = accessToken.ID
						newMethodData.Name = "messages.send"
						listNewMethodData.APM = append(listNewMethodData.APM, newMethodData)
					}
				}

				if len(kitSAdditionSendToinAPM.Entry.Text()) == 0 {
					warningTitle := "Field \"Send to\" must not be empty."
					ShowWarningWindow(warningTitle)
					return
				}
				sendTo, err := strconv.Atoi(kitSAdditionSendToinAPM.Entry.Text())
				if err != nil {
					ToLogFile(err.Error(), string(debug.Stack()))
					panic(err.Error())
				}
				var newMonitorModuleData NewMonitorModuleData
				newMonitorModuleData.SendTo = sendTo
				listNewMonitorModuleData.APM = newMonitorModuleData

			case "video_monitor":

				if kitSAdditionVGinVM.Combobox.Selected() == -1 {
					warningTitle := "You must select an item in the combobox " +
						"\"Access token for \"video.get\"\""
					ShowWarningWindow(warningTitle)
					return
				}
				accessTokenName := accessTokensNames[kitSAdditionVGinVM.Combobox.Selected()]
				for _, accessToken := range accessTokens {
					if accessTokenName == accessToken.Name {
						var newMethodData NewMethodData
						newMethodData.AccessTokenID = accessToken.ID
						newMethodData.Name = "video.get"
						listNewMethodData.VM = append(listNewMethodData.VM, newMethodData)
					}
				}

				if kitSAdditionUGinVM.Combobox.Selected() == -1 {
					warningTitle := "You must select an item in the combobox " +
						"\"Access token for \"users.get\"\""
					ShowWarningWindow(warningTitle)
					return
				}
				accessTokenName = accessTokensNames[kitSAdditionUGinVM.Combobox.Selected()]
				for _, accessToken := range accessTokens {
					if accessTokenName == accessToken.Name {
						var newMethodData NewMethodData
						newMethodData.AccessTokenID = accessToken.ID
						newMethodData.Name = "users.get"
						listNewMethodData.VM = append(listNewMethodData.VM, newMethodData)
					}
				}

				if kitSAdditionGGBIinVM.Combobox.Selected() == -1 {
					warningTitle := "You must select an item in the combobox " +
						"\"Access token for \"groups.getById\"\""
					ShowWarningWindow(warningTitle)
					return
				}
				accessTokenName = accessTokensNames[kitSAdditionGGBIinVM.Combobox.Selected()]
				for _, accessToken := range accessTokens {
					if accessTokenName == accessToken.Name {
						var newMethodData NewMethodData
						newMethodData.AccessTokenID = accessToken.ID
						newMethodData.Name = "groups.getById"
						listNewMethodData.VM = append(listNewMethodData.VM, newMethodData)
					}
				}

				if kitSAdditionMSinVM.Combobox.Selected() == -1 {
					warningTitle := "You must select an item in the combobox " +
						"\"Access token for \"messages.send\"\""
					ShowWarningWindow(warningTitle)
					return
				}
				accessTokenName = accessTokensNames[kitSAdditionMSinVM.Combobox.Selected()]
				for _, accessToken := range accessTokens {
					if accessTokenName == accessToken.Name {
						var newMethodData NewMethodData
						newMethodData.AccessTokenID = accessToken.ID
						newMethodData.Name = "messages.send"
						listNewMethodData.VM = append(listNewMethodData.VM, newMethodData)
					}
				}

				if len(kitSAdditionSendToinVM.Entry.Text()) == 0 {
					warningTitle := "Field \"Send to\" must not be empty."
					ShowWarningWindow(warningTitle)
					return
				}
				sendTo, err := strconv.Atoi(kitSAdditionSendToinVM.Entry.Text())
				if err != nil {
					ToLogFile(err.Error(), string(debug.Stack()))
					panic(err.Error())
				}
				var newMonitorModuleData NewMonitorModuleData
				newMonitorModuleData.SendTo = sendTo
				listNewMonitorModuleData.VM = newMonitorModuleData

			case "photo_comment_monitor":
				if kitSAdditionPGACinPCM.Combobox.Selected() == -1 {
					warningTitle := "You must select an item in the combobox " +
						"\"Access token for \"photos.getAllComments\"\""
					ShowWarningWindow(warningTitle)
					return
				}
				accessTokenName := accessTokensNames[kitSAdditionPGACinPCM.Combobox.Selected()]
				for _, accessToken := range accessTokens {
					if accessTokenName == accessToken.Name {
						var newMethodData NewMethodData
						newMethodData.AccessTokenID = accessToken.ID
						newMethodData.Name = "photos.getAllComments"
						listNewMethodData.PCM = append(listNewMethodData.PCM, newMethodData)
					}
				}

				if kitSAdditionUGinPCM.Combobox.Selected() == -1 {
					warningTitle := "You must select an item in the combobox " +
						"\"Access token for \"users.get\"\""
					ShowWarningWindow(warningTitle)
					return
				}
				accessTokenName = accessTokensNames[kitSAdditionUGinPCM.Combobox.Selected()]
				for _, accessToken := range accessTokens {
					if accessTokenName == accessToken.Name {
						var newMethodData NewMethodData
						newMethodData.AccessTokenID = accessToken.ID
						newMethodData.Name = "users.get"
						listNewMethodData.PCM = append(listNewMethodData.PCM, newMethodData)
					}
				}

				if kitSAdditionGGBIinPCM.Combobox.Selected() == -1 {
					warningTitle := "You must select an item in the combobox " +
						"\"Access token for \"groups.getById\"\""
					ShowWarningWindow(warningTitle)
					return
				}
				accessTokenName = accessTokensNames[kitSAdditionGGBIinPCM.Combobox.Selected()]
				for _, accessToken := range accessTokens {
					if accessTokenName == accessToken.Name {
						var newMethodData NewMethodData
						newMethodData.AccessTokenID = accessToken.ID
						newMethodData.Name = "groups.getById"
						listNewMethodData.PCM = append(listNewMethodData.PCM, newMethodData)
					}
				}

				if kitSAdditionMSinPCM.Combobox.Selected() == -1 {
					warningTitle := "You must select an item in the combobox " +
						"\"Access token for \"messages.send\"\""
					ShowWarningWindow(warningTitle)
					return
				}
				accessTokenName = accessTokensNames[kitSAdditionMSinPCM.Combobox.Selected()]
				for _, accessToken := range accessTokens {
					if accessTokenName == accessToken.Name {
						var newMethodData NewMethodData
						newMethodData.AccessTokenID = accessToken.ID
						newMethodData.Name = "messages.send"
						listNewMethodData.PCM = append(listNewMethodData.PCM, newMethodData)
					}
				}

				if len(kitSAdditionSendToinPCM.Entry.Text()) == 0 {
					warningTitle := "Field \"Send to\" must not be empty."
					ShowWarningWindow(warningTitle)
					return
				}
				sendTo, err := strconv.Atoi(kitSAdditionSendToinPCM.Entry.Text())
				if err != nil {
					ToLogFile(err.Error(), string(debug.Stack()))
					panic(err.Error())
				}
				var newMonitorModuleData NewMonitorModuleData
				newMonitorModuleData.SendTo = sendTo
				listNewMonitorModuleData.PCM = newMonitorModuleData

			case "video_comment_monitor":
				if kitSAdditionVGCinVCM.Combobox.Selected() == -1 {
					warningTitle := "You must select an item in the combobox " +
						"\"Access token for \"video.getComments\"\""
					ShowWarningWindow(warningTitle)
					return
				}
				accessTokenName := accessTokensNames[kitSAdditionVGCinVCM.Combobox.Selected()]
				for _, accessToken := range accessTokens {
					if accessTokenName == accessToken.Name {
						var newMethodData NewMethodData
						newMethodData.AccessTokenID = accessToken.ID
						newMethodData.Name = "video.getComments"
						listNewMethodData.VCM = append(listNewMethodData.VCM, newMethodData)
					}
				}

				if kitSAdditionUGinVCM.Combobox.Selected() == -1 {
					warningTitle := "You must select an item in the combobox " +
						"\"Access token for \"users.get\"\""
					ShowWarningWindow(warningTitle)
					return
				}
				accessTokenName = accessTokensNames[kitSAdditionUGinVCM.Combobox.Selected()]
				for _, accessToken := range accessTokens {
					if accessTokenName == accessToken.Name {
						var newMethodData NewMethodData
						newMethodData.AccessTokenID = accessToken.ID
						newMethodData.Name = "users.get"
						listNewMethodData.VCM = append(listNewMethodData.VCM, newMethodData)
					}
				}

				if kitSAdditionGGBIinVCM.Combobox.Selected() == -1 {
					warningTitle := "You must select an item in the combobox " +
						"\"Access token for \"groups.getById\"\""
					ShowWarningWindow(warningTitle)
					return
				}
				accessTokenName = accessTokensNames[kitSAdditionGGBIinVCM.Combobox.Selected()]
				for _, accessToken := range accessTokens {
					if accessTokenName == accessToken.Name {
						var newMethodData NewMethodData
						newMethodData.AccessTokenID = accessToken.ID
						newMethodData.Name = "groups.getById"
						listNewMethodData.VCM = append(listNewMethodData.VCM, newMethodData)
					}
				}

				if kitSAdditionVGinVCM.Combobox.Selected() == -1 {
					warningTitle := "You must select an item in the combobox " +
						"\"Access token for \"video.get\"\""
					ShowWarningWindow(warningTitle)
					return
				}
				accessTokenName = accessTokensNames[kitSAdditionVGinVCM.Combobox.Selected()]
				for _, accessToken := range accessTokens {
					if accessTokenName == accessToken.Name {
						var newMethodData NewMethodData
						newMethodData.AccessTokenID = accessToken.ID
						newMethodData.Name = "video.get"
						listNewMethodData.VCM = append(listNewMethodData.VCM, newMethodData)
					}
				}

				if kitSAdditionMSinVCM.Combobox.Selected() == -1 {
					warningTitle := "You must select an item in the combobox " +
						"\"Access token for \"messages.send\"\""
					ShowWarningWindow(warningTitle)
					return
				}
				accessTokenName = accessTokensNames[kitSAdditionMSinVCM.Combobox.Selected()]
				for _, accessToken := range accessTokens {
					if accessTokenName == accessToken.Name {
						var newMethodData NewMethodData
						newMethodData.AccessTokenID = accessToken.ID
						newMethodData.Name = "messages.send"
						listNewMethodData.VCM = append(listNewMethodData.VCM, newMethodData)
					}
				}

				if len(kitSAdditionSendToinVCM.Entry.Text()) == 0 {
					warningTitle := "Field \"Send to\" must not be empty."
					ShowWarningWindow(warningTitle)
					return
				}
				sendTo, err := strconv.Atoi(kitSAdditionSendToinVCM.Entry.Text())
				if err != nil {
					ToLogFile(err.Error(), string(debug.Stack()))
					panic(err.Error())
				}
				var newMonitorModuleData NewMonitorModuleData
				newMonitorModuleData.SendTo = sendTo
				listNewMonitorModuleData.VCM = newMonitorModuleData

			case "topic_monitor":
				if kitSAdditionBGCinTM.Combobox.Selected() == -1 {
					warningTitle := "You must select an item in the combobox " +
						"\"Access token for \"board.getComments\"\""
					ShowWarningWindow(warningTitle)
					return
				}
				accessTokenName := accessTokensNames[kitSAdditionBGCinTM.Combobox.Selected()]
				for _, accessToken := range accessTokens {
					if accessTokenName == accessToken.Name {
						var newMethodData NewMethodData
						newMethodData.AccessTokenID = accessToken.ID
						newMethodData.Name = "board.getComments"
						listNewMethodData.TM = append(listNewMethodData.TM, newMethodData)
					}
				}

				if kitSAdditionBGTinTM.Combobox.Selected() == -1 {
					warningTitle := "You must select an item in the combobox " +
						"\"Access token for \"board.getTopicsgetComments\"\""
					ShowWarningWindow(warningTitle)
					return
				}
				accessTokenName = accessTokensNames[kitSAdditionBGTinTM.Combobox.Selected()]
				for _, accessToken := range accessTokens {
					if accessTokenName == accessToken.Name {
						var newMethodData NewMethodData
						newMethodData.AccessTokenID = accessToken.ID
						newMethodData.Name = "board.getTopics"
						listNewMethodData.TM = append(listNewMethodData.TM, newMethodData)
					}
				}

				if kitSAdditionUGinTM.Combobox.Selected() == -1 {
					warningTitle := "You must select an item in the combobox " +
						"\"Access token for \"users.get\"\""
					ShowWarningWindow(warningTitle)
					return
				}
				accessTokenName = accessTokensNames[kitSAdditionUGinTM.Combobox.Selected()]
				for _, accessToken := range accessTokens {
					if accessTokenName == accessToken.Name {
						var newMethodData NewMethodData
						newMethodData.AccessTokenID = accessToken.ID
						newMethodData.Name = "users.get"
						listNewMethodData.TM = append(listNewMethodData.TM, newMethodData)
					}
				}

				if kitSAdditionGGBIinTM.Combobox.Selected() == -1 {
					warningTitle := "You must select an item in the combobox " +
						"\"Access token for \"groups.getById\"\""
					ShowWarningWindow(warningTitle)
					return
				}
				accessTokenName = accessTokensNames[kitSAdditionGGBIinTM.Combobox.Selected()]
				for _, accessToken := range accessTokens {
					if accessTokenName == accessToken.Name {
						var newMethodData NewMethodData
						newMethodData.AccessTokenID = accessToken.ID
						newMethodData.Name = "groups.getById"
						listNewMethodData.TM = append(listNewMethodData.TM, newMethodData)
					}
				}

				if kitSAdditionMSinTM.Combobox.Selected() == -1 {
					warningTitle := "You must select an item in the combobox " +
						"\"Access token for \"messages.send\"\""
					ShowWarningWindow(warningTitle)
					return
				}
				accessTokenName = accessTokensNames[kitSAdditionMSinTM.Combobox.Selected()]
				for _, accessToken := range accessTokens {
					if accessTokenName == accessToken.Name {
						var newMethodData NewMethodData
						newMethodData.AccessTokenID = accessToken.ID
						newMethodData.Name = "messages.send"
						listNewMethodData.TM = append(listNewMethodData.TM, newMethodData)
					}
				}

				if len(kitSAdditionSendToinTM.Entry.Text()) == 0 {
					warningTitle := "Field \"Send to\" must not be empty."
					ShowWarningWindow(warningTitle)
					return
				}
				sendTo, err := strconv.Atoi(kitSAdditionSendToinTM.Entry.Text())
				if err != nil {
					ToLogFile(err.Error(), string(debug.Stack()))
					panic(err.Error())
				}
				var newMonitorModuleData NewMonitorModuleData
				newMonitorModuleData.SendTo = sendTo
				listNewMonitorModuleData.TM = newMonitorModuleData

			case "wall_post_comment_monitor":
				if kitSAdditionWGCsinWPCM.Combobox.Selected() == -1 {
					warningTitle := "You must select an item in the combobox " +
						"\"Access token for \"wall.getComments\"\""
					ShowWarningWindow(warningTitle)
					return
				}
				accessTokenName := accessTokensNames[kitSAdditionWGCsinWPCM.Combobox.Selected()]
				for _, accessToken := range accessTokens {
					if accessTokenName == accessToken.Name {
						var newMethodData NewMethodData
						newMethodData.AccessTokenID = accessToken.ID
						newMethodData.Name = "wall.getComments"
						listNewMethodData.WPCM = append(listNewMethodData.WPCM, newMethodData)
					}
				}

				if kitSAdditionUGinWPCM.Combobox.Selected() == -1 {
					warningTitle := "You must select an item in the combobox " +
						"\"Access token for \"users.get\"\""
					ShowWarningWindow(warningTitle)
					return
				}
				accessTokenName = accessTokensNames[kitSAdditionUGinWPCM.Combobox.Selected()]
				for _, accessToken := range accessTokens {
					if accessTokenName == accessToken.Name {
						var newMethodData NewMethodData
						newMethodData.AccessTokenID = accessToken.ID
						newMethodData.Name = "users.get"
						listNewMethodData.WPCM = append(listNewMethodData.WPCM, newMethodData)
					}
				}

				if kitSAdditionGGBIinWPCM.Combobox.Selected() == -1 {
					warningTitle := "You must select an item in the combobox " +
						"\"Access token for \"groups.getById\"\""
					ShowWarningWindow(warningTitle)
					return
				}
				accessTokenName = accessTokensNames[kitSAdditionGGBIinWPCM.Combobox.Selected()]
				for _, accessToken := range accessTokens {
					if accessTokenName == accessToken.Name {
						var newMethodData NewMethodData
						newMethodData.AccessTokenID = accessToken.ID
						newMethodData.Name = "groups.getById"
						listNewMethodData.WPCM = append(listNewMethodData.WPCM, newMethodData)
					}
				}

				if kitSAdditionWGinWPCM.Combobox.Selected() == -1 {
					warningTitle := "You must select an item in the combobox " +
						"\"Access token for \"wall.get\"\""
					ShowWarningWindow(warningTitle)
					return
				}
				accessTokenName = accessTokensNames[kitSAdditionWGinWPCM.Combobox.Selected()]
				for _, accessToken := range accessTokens {
					if accessTokenName == accessToken.Name {
						var newMethodData NewMethodData
						newMethodData.AccessTokenID = accessToken.ID
						newMethodData.Name = "wall.get"
						listNewMethodData.WPCM = append(listNewMethodData.WPCM, newMethodData)
					}
				}

				if kitSAdditionWGCinWPCM.Combobox.Selected() == -1 {
					warningTitle := "You must select an item in the combobox " +
						"\"Access token for \"wall.getComment\"\""
					ShowWarningWindow(warningTitle)
					return
				}
				accessTokenName = accessTokensNames[kitSAdditionWGCinWPCM.Combobox.Selected()]
				for _, accessToken := range accessTokens {
					if accessTokenName == accessToken.Name {
						var newMethodData NewMethodData
						newMethodData.AccessTokenID = accessToken.ID
						newMethodData.Name = "wall.getComment"
						listNewMethodData.WPCM = append(listNewMethodData.WPCM, newMethodData)
					}
				}

				if kitSAdditionMSinWPCM.Combobox.Selected() == -1 {
					warningTitle := "You must select an item in the combobox " +
						"\"Access token for \"messages.send\"\""
					ShowWarningWindow(warningTitle)
					return
				}
				accessTokenName = accessTokensNames[kitSAdditionMSinWPCM.Combobox.Selected()]
				for _, accessToken := range accessTokens {
					if accessTokenName == accessToken.Name {
						var newMethodData NewMethodData
						newMethodData.AccessTokenID = accessToken.ID
						newMethodData.Name = "messages.send"
						listNewMethodData.WPCM = append(listNewMethodData.WPCM, newMethodData)
					}
				}

				if len(kitSAdditionSendToinWPCM.Entry.Text()) == 0 {
					warningTitle := "Field \"Send to\" must not be empty."
					ShowWarningWindow(warningTitle)
					return
				}
				sendTo, err := strconv.Atoi(kitSAdditionSendToinWPCM.Entry.Text())
				if err != nil {
					ToLogFile(err.Error(), string(debug.Stack()))
					panic(err.Error())
				}
				var newMonitorModuleData NewMonitorModuleData
				newMonitorModuleData.SendTo = sendTo
				listNewMonitorModuleData.WPCM = newMonitorModuleData
			}
		}

		additionNewSubject(newSubjectData)

		subjects, err := SelectDBSubjects()
		if err != nil {
			ToLogFile(err.Error(), string(debug.Stack()))
			panic(err.Error())
		}
		id := subjects[len(subjects)-1].ID

		for _, monitorName := range monitorsNames {
			switch monitorName {
			case "wall_post_monitor":
				additionNewMonitor(monitorName, id)

				monitor, err := SelectDBMonitor(monitorName, id)
				if err != nil {
					ToLogFile(err.Error(), string(debug.Stack()))
					panic(err.Error())
				}

				for i := 0; i < len(listNewMethodData.WPM); i++ {
					listNewMethodData.WPM[i].SubjectID = id
					listNewMethodData.WPM[i].MonitorID = monitor.ID
					additionNewMethod(listNewMethodData.WPM[i])
				}

				listNewMonitorModuleData.WPM.SubjectID = id
				additionNewWallPostMonitor(listNewMonitorModuleData.WPM)

			case "album_photo_monitor":
				additionNewMonitor(monitorName, id)

				monitor, err := SelectDBMonitor(monitorName, id)
				if err != nil {
					ToLogFile(err.Error(), string(debug.Stack()))
					panic(err.Error())
				}

				for i := 0; i < len(listNewMethodData.APM); i++ {
					listNewMethodData.APM[i].SubjectID = id
					listNewMethodData.APM[i].MonitorID = monitor.ID
					additionNewMethod(listNewMethodData.APM[i])
				}

				listNewMonitorModuleData.APM.SubjectID = id
				additionNewAlbumPhotoMonitor(listNewMonitorModuleData.APM)

			case "video_monitor":
				additionNewMonitor(monitorName, id)

				monitor, err := SelectDBMonitor(monitorName, id)
				if err != nil {
					ToLogFile(err.Error(), string(debug.Stack()))
					panic(err.Error())
				}

				for i := 0; i < len(listNewMethodData.VM); i++ {
					listNewMethodData.VM[i].SubjectID = id
					listNewMethodData.VM[i].MonitorID = monitor.ID
					additionNewMethod(listNewMethodData.VM[i])
				}

				listNewMonitorModuleData.VM.SubjectID = id
				additionNewVideoMonitor(listNewMonitorModuleData.VM)

			case "photo_comment_monitor":
				additionNewMonitor(monitorName, id)

				monitor, err := SelectDBMonitor(monitorName, id)
				if err != nil {
					ToLogFile(err.Error(), string(debug.Stack()))
					panic(err.Error())
				}

				for i := 0; i < len(listNewMethodData.PCM); i++ {
					listNewMethodData.PCM[i].SubjectID = id
					listNewMethodData.PCM[i].MonitorID = monitor.ID
					additionNewMethod(listNewMethodData.PCM[i])
				}

				listNewMonitorModuleData.PCM.SubjectID = id
				additionPhotoCommentMonitor(listNewMonitorModuleData.PCM)

			case "video_comment_monitor":
				additionNewMonitor(monitorName, id)

				monitor, err := SelectDBMonitor(monitorName, id)
				if err != nil {
					ToLogFile(err.Error(), string(debug.Stack()))
					panic(err.Error())
				}

				for i := 0; i < len(listNewMethodData.VCM); i++ {
					listNewMethodData.VCM[i].SubjectID = id
					listNewMethodData.VCM[i].MonitorID = monitor.ID
					additionNewMethod(listNewMethodData.VCM[i])
				}

				listNewMonitorModuleData.VCM.SubjectID = id
				additionVideoCommentMonitor(listNewMonitorModuleData.VCM)

			case "topic_monitor":
				additionNewMonitor(monitorName, id)

				monitor, err := SelectDBMonitor(monitorName, id)
				if err != nil {
					ToLogFile(err.Error(), string(debug.Stack()))
					panic(err.Error())
				}

				for i := 0; i < len(listNewMethodData.TM); i++ {
					listNewMethodData.TM[i].SubjectID = id
					listNewMethodData.TM[i].MonitorID = monitor.ID
					additionNewMethod(listNewMethodData.TM[i])
				}

				listNewMonitorModuleData.TM.SubjectID = id
				additionTopicMonitor(listNewMonitorModuleData.TM)

			case "wall_post_comment_monitor":
				additionNewMonitor(monitorName, id)

				monitor, err := SelectDBMonitor(monitorName, id)
				if err != nil {
					ToLogFile(err.Error(), string(debug.Stack()))
					panic(err.Error())
				}

				for i := 0; i < len(listNewMethodData.WPCM); i++ {
					listNewMethodData.WPCM[i].SubjectID = id
					listNewMethodData.WPCM[i].MonitorID = monitor.ID
					additionNewMethod(listNewMethodData.WPCM[i])
				}

				listNewMonitorModuleData.WPCM.SubjectID = id
				additionWallPostCommentMonitor(listNewMonitorModuleData.WPCM)
			}
		}

		kitWindowSubjectAddition.Window.Disable()
		kitWindowSubjectAddition.Window.Hide()
	})

	// добавляем коробку с кнопками на основную коробку окна
	kitWindowSubjectAddition.Box.Append(kitButtonsSAddition.Box, true)
	// затем еще одну коробку, для выравнивания расположения кнопок при растягивании окна
	boxWndSAdditionBottom := ui.NewHorizontalBox()
	kitWindowSubjectAddition.Box.Append(boxWndSAdditionBottom, false)

	kitWindowSubjectAddition.Window.Show()
}

func additionNewSubject(newSubjectData NewSubjectData) {
	var subject Subject

	subject.Name = newSubjectData.Name
	subject.SubjectID = newSubjectData.ID
	subject.BackupWikipage = "-0_0" // этот параметр нигде не используется
	subject.LastBackup = 0          // этот тоже

	err := InsertDBSubject(subject)
	if err != nil {
		ToLogFile(err.Error(), string(debug.Stack()))
		panic(err.Error())
	}
}

func additionNewMonitor(monitorName string, subjectID int) {
	var monitor Monitor

	monitor.Name = monitorName
	monitor.SubjectID = subjectID

	err := InsertDBMonitor(monitor)
	if err != nil {
		ToLogFile(err.Error(), string(debug.Stack()))
		panic(err.Error())
	}
}

func additionNewMethod(newMethodData NewMethodData) {
	var method Method

	method.Name = newMethodData.Name
	method.SubjectID = newMethodData.SubjectID
	method.AccessTokenID = newMethodData.AccessTokenID
	method.MonitorID = newMethodData.MonitorID

	err := InsertDBMethod(method)
	if err != nil {
		ToLogFile(err.Error(), string(debug.Stack()))
		panic(err.Error())
	}
}

func additionNewWallPostMonitor(newMonitorModuleData NewMonitorModuleData) {
	var wallPostMonitorParam WallPostMonitorParam

	wallPostMonitorParam.SubjectID = newMonitorModuleData.SubjectID
	wallPostMonitorParam.NeedMonitoring = 0
	wallPostMonitorParam.Interval = 60
	wallPostMonitorParam.SendTo = newMonitorModuleData.SendTo
	wallPostMonitorParam.Filter = "all"
	wallPostMonitorParam.LastDate = 0
	wallPostMonitorParam.PostsCount = 5
	wallPostMonitorParam.KeywordsForMonitoring = "{\"list\":[]}"
	wallPostMonitorParam.UsersIDsForIgnore = "{\"list\":[]}"

	err := InsertDBWallPostMonitor(wallPostMonitorParam)
	if err != nil {
		ToLogFile(err.Error(), string(debug.Stack()))
		panic(err.Error())
	}
}

func additionNewAlbumPhotoMonitor(newMonitorModuleData NewMonitorModuleData) {
	var albumPhotoMonitorParam AlbumPhotoMonitorParam

	albumPhotoMonitorParam.SubjectID = newMonitorModuleData.SubjectID
	albumPhotoMonitorParam.NeedMonitoring = 0
	albumPhotoMonitorParam.SendTo = newMonitorModuleData.SendTo
	albumPhotoMonitorParam.Interval = 60
	albumPhotoMonitorParam.LastDate = 0
	albumPhotoMonitorParam.PhotosCount = 5

	err := InsertDBAlbumPhotoMonitor(albumPhotoMonitorParam)
	if err != nil {
		ToLogFile(err.Error(), string(debug.Stack()))
		panic(err.Error())
	}
}

func additionNewVideoMonitor(newMonitorModuleData NewMonitorModuleData) {
	var videoMonitorParam VideoMonitorParam

	videoMonitorParam.SubjectID = newMonitorModuleData.SubjectID
	videoMonitorParam.NeedMonitoring = 0
	videoMonitorParam.SendTo = newMonitorModuleData.SendTo
	videoMonitorParam.VideoCount = 5
	videoMonitorParam.LastDate = 0
	videoMonitorParam.Interval = 60

	err := InsertDBVideoMonitor(videoMonitorParam)
	if err != nil {
		ToLogFile(err.Error(), string(debug.Stack()))
		panic(err.Error())
	}
}

func additionPhotoCommentMonitor(newMonitorModuleData NewMonitorModuleData) {
	var photoCommentMonitorParam PhotoCommentMonitorParam

	photoCommentMonitorParam.SubjectID = newMonitorModuleData.SubjectID
	photoCommentMonitorParam.NeedMonitoring = 0
	photoCommentMonitorParam.CommentsCount = 5
	photoCommentMonitorParam.LastDate = 0
	photoCommentMonitorParam.Interval = 60
	photoCommentMonitorParam.SendTo = newMonitorModuleData.SendTo

	err := InsertDBPhotoCommentMonitor(photoCommentMonitorParam)
	if err != nil {
		ToLogFile(err.Error(), string(debug.Stack()))
		panic(err.Error())
	}
}

func additionVideoCommentMonitor(newMonitorModuleData NewMonitorModuleData) {
	var videoCommentMonitorParam VideoCommentMonitorParam

	videoCommentMonitorParam.SubjectID = newMonitorModuleData.SubjectID
	videoCommentMonitorParam.NeedMonitoring = 0
	videoCommentMonitorParam.VideosCount = 5
	videoCommentMonitorParam.Interval = 60
	videoCommentMonitorParam.CommentsCount = 5
	videoCommentMonitorParam.SendTo = newMonitorModuleData.SendTo
	videoCommentMonitorParam.LastDate = 0

	err := InsertDBVideoCommentMonitor(videoCommentMonitorParam)
	if err != nil {
		ToLogFile(err.Error(), string(debug.Stack()))
		panic(err.Error())
	}
}

func additionTopicMonitor(newMonitorModuleData NewMonitorModuleData) {
	var topicMonitorParam TopicMonitorParam

	topicMonitorParam.SubjectID = newMonitorModuleData.SubjectID
	topicMonitorParam.NeedMonitoring = 0
	topicMonitorParam.TopicsCount = 5
	topicMonitorParam.CommentsCount = 5
	topicMonitorParam.Interval = 60
	topicMonitorParam.SendTo = newMonitorModuleData.SendTo
	topicMonitorParam.LastDate = 0

	err := InsertDBTopicMonitor(topicMonitorParam)
	if err != nil {
		ToLogFile(err.Error(), string(debug.Stack()))
		panic(err.Error())
	}
}

func additionWallPostCommentMonitor(newMonitorModuleData NewMonitorModuleData) {
	var wallPostCommentMonitorParam WallPostCommentMonitorParam

	wallPostCommentMonitorParam.SubjectID = newMonitorModuleData.SubjectID
	wallPostCommentMonitorParam.NeedMonitoring = 0
	wallPostCommentMonitorParam.PostsCount = 5
	wallPostCommentMonitorParam.CommentsCount = 5
	wallPostCommentMonitorParam.MonitoringAll = 1
	wallPostCommentMonitorParam.UsersIDsForMonitoring = "{\"list\":[]}"
	wallPostCommentMonitorParam.UsersNamesForMonitoring = "{\"list\":[]}"
	wallPostCommentMonitorParam.AttachmentsTypesForMonitoring = "{\"list\":[\"photo\", \"video\", \"audio\", \"doc\", \"poll\", \"link\"]}"
	wallPostCommentMonitorParam.UsersIDsForIgnore = "{\"list\":[]}"
	wallPostCommentMonitorParam.Interval = 60
	wallPostCommentMonitorParam.SendTo = newMonitorModuleData.SendTo
	wallPostCommentMonitorParam.Filter = "all"
	wallPostCommentMonitorParam.LastDate = 0
	wallPostCommentMonitorParam.KeywordsForMonitoring = "{\"list\":[]}"
	wallPostCommentMonitorParam.SmallCommentsForMonitoring = "{\"list\":[]}"
	wallPostCommentMonitorParam.DigitsCountForCardNumberMonitoring = "{\"list\":[\"16\"]}"
	wallPostCommentMonitorParam.DigitsCountForPhoneNumberMonitoring = "{\"list\":[\"6\",\"11\"]}"
	wallPostCommentMonitorParam.MonitorByCommunity = 1

	err := InsertDBWallPostCommentMonitor(wallPostCommentMonitorParam)
	if err != nil {
		ToLogFile(err.Error(), string(debug.Stack()))
		panic(err.Error())
	}
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
		ToLogFile(err.Error(), string(debug.Stack()))
		panic(err.Error())
	}

	// получаем набор для отображения общих установок субъекта мониторинга
	kitWindowGeneralSettings := MakeSettingWindowKit("", 300, 100)

	// перечисляем субъекты
	for _, subject := range subjects {
		// ищем субъект с подходящим идентификатором
		if subject.ID == IDSubject {
			// устанавливаем заголовок окна в соответствии с названием субъекта и назначением установок
			windowTitle := fmt.Sprintf("%v settings for %v", btnName, subject.Name)
			kitWindowGeneralSettings.Window.SetTitle(windowTitle)

			boxWndS := ui.NewVerticalBox()

			// получаем набор для названия субъекта мониторинга
			kitWndSName := MakeSettingEntryKit("Name", subject.Name)

			// получаем набор для идентификатора субъекта мониторинга в базе ВК
			kitWndSSubjectID := MakeSettingEntryKit("Subject ID", strconv.Itoa(subject.SubjectID))
			// и привязываем его к обработке ввода, чтобы кроме чисел ничего не вводилось
			kitWndSSubjectID.Entry.OnChanged(func(*ui.Entry) {
				NumericEntriesHandler(kitWndSSubjectID.Entry)
			})

			// описываем группу, в которой будут размещены элементы
			groupWndS := ui.NewGroup("")
			groupWndS.SetMargined(true)
			boxWndS.Append(kitWndSName.Box, false)
			boxWndS.Append(kitWndSSubjectID.Box, false)
			groupWndS.SetChild(boxWndS)

			// добавляем группу в основную коробку окна
			kitWindowGeneralSettings.Box.Append(groupWndS, false)

			// получаем набор для кнопок принятия и отмены изменений
			kitButtonsS := MakeSettingButtonsKit()

			// привязываем к кнопкам соответствующие процедуры
			kitButtonsS.ButtonCancel.OnClicked(func(*ui.Button) {
				kitWindowGeneralSettings.Window.Disable()
				kitWindowGeneralSettings.Window.Hide()
			})
			// привязываем кнопки к соответствующим процедурам
			kitButtonsS.ButtonApply.OnClicked(func(*ui.Button) {
				var updatedSubject Subject
				updatedSubject.ID = subject.ID
				if len(kitWndSSubjectID.Entry.Text()) == 0 {
					warningTitle := "Field \"Subject ID\" must not be empty."
					ShowWarningWindow(warningTitle)
					return
				}
				updatedSubject.SubjectID, err = strconv.Atoi(kitWndSSubjectID.Entry.Text())
				if err != nil {
					ToLogFile(err.Error(), string(debug.Stack()))
					panic(err.Error())
				}
				if len(kitWndSName.Entry.Text()) == 0 {
					warningTitle := "Field \"Name\" must not be empty."
					ShowWarningWindow(warningTitle)
					return
				}
				updatedSubject.Name = kitWndSName.Entry.Text()
				updatedSubject.BackupWikipage = subject.BackupWikipage
				updatedSubject.LastBackup = subject.LastBackup

				err := UpdateDBSubject(updatedSubject)
				if err != nil {
					ToLogFile(err.Error(), string(debug.Stack()))
					panic(err.Error())
				}

				kitWindowGeneralSettings.Window.Disable()
				kitWindowGeneralSettings.Window.Hide()
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
		ToLogFile(err.Error(), string(debug.Stack()))
		panic(err.Error())
	}

	// получаем набор для отображения установок модуля мониторинга постов на стене
	windowTitle := fmt.Sprintf("%v settings for %v", btnName, nameSubject)
	kitWindowWallPostSettings := MakeSettingWindowKit(windowTitle, 300, 100)

	boxWndWP := ui.NewVerticalBox()

	// получаем набор для флага необходимости активировать модуль мониторинга
	kitWndWPMonitoring := MakeSettingCheckboxKit("Need monitoring", wallPostMonitorParam.NeedMonitoring)

	// получаем набор для интервала между запусками функции мониторинга
	kitWndWPInterval := MakeSettingSpinboxKit("Interval", 5, 21600, wallPostMonitorParam.Interval)

	// получаем набор для количества проверяемых постов
	kitWndWPSendTo := MakeSettingEntryKit("Send to", strconv.Itoa(wallPostMonitorParam.SendTo))
	// и привязываем его к обработке ввода, чтобы кроме чисел ничего не вводилось
	kitWndWPSendTo.Entry.OnChanged(func(*ui.Entry) {
		NumericEntriesHandler(kitWndWPSendTo.Entry)
	})

	// получаем набор для фильтра получаемых для проверки постов
	listPostsFilters := []string{"all", "others", "owner", "suggests"}
	kitWndWPFilter := MakeSettingComboboxKit("Filter", listPostsFilters, wallPostMonitorParam.Filter)

	// получаем набор для количества проверяемых постов
	kitWndWPPostsCount := MakeSettingSpinboxKit("Posts count", 1, 100, wallPostMonitorParam.PostsCount)

	// получаем набор для списка ключевых слов для отбора постов
	kitWndWPKeywordsForMonitoring := MakeSettingEntryListKit("Keywords", wallPostMonitorParam.KeywordsForMonitoring)

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
	kitButtonsWP := MakeSettingButtonsKit()

	// привязываем к кнопкам соответствующие процедуры
	kitButtonsWP.ButtonCancel.OnClicked(func(*ui.Button) {
		kitWindowWallPostSettings.Window.Disable()
		kitWindowWallPostSettings.Window.Hide()
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
		if len(kitWndWPSendTo.Entry.Text()) == 0 {
			warningTitle := "Field \"Send to\" must not be empty."
			ShowWarningWindow(warningTitle)
			return
		}
		updatedWallPostMonitorParam.SendTo, err = strconv.Atoi(kitWndWPSendTo.Entry.Text())
		if err != nil {
			ToLogFile(err.Error(), string(debug.Stack()))
			panic(err.Error())
		}
		listPostsFilters := []string{"all", "others", "owner", "suggests"}
		if kitWndWPFilter.Combobox.Selected() == -1 {
			warningTitle := "You must select an item in the combobox \"Filter\""
			ShowWarningWindow(warningTitle)
			return
		}
		updatedWallPostMonitorParam.Filter = listPostsFilters[kitWndWPFilter.Combobox.Selected()]
		updatedWallPostMonitorParam.LastDate = wallPostMonitorParam.LastDate
		updatedWallPostMonitorParam.PostsCount = kitWndWPPostsCount.Spinbox.Value()
		// TODO: проверка соответствия оформления требованиям json
		jsonDump := fmt.Sprintf("{\"list\":[%v]}", kitWndWPKeywordsForMonitoring.Entry.Text())
		updatedWallPostMonitorParam.KeywordsForMonitoring = jsonDump
		updatedWallPostMonitorParam.UsersIDsForIgnore = wallPostMonitorParam.UsersIDsForIgnore

		err = UpdateDBWallPostMonitor(updatedWallPostMonitorParam)
		if err != nil {
			ToLogFile(err.Error(), string(debug.Stack()))
			panic(err.Error())
		}

		kitWindowWallPostSettings.Window.Disable()
		kitWindowWallPostSettings.Window.Hide()
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
		ToLogFile(err.Error(), string(debug.Stack()))
		panic(err.Error())
	}

	// получаем набор для отображения установок модуля мониторинга фотографий в альбомах
	windowTitle := fmt.Sprintf("%v settings for %v", btnName, nameSubject)
	kitWindowAlbumPhotoSettings := MakeSettingWindowKit(windowTitle, 300, 100)

	boxWndAP := ui.NewVerticalBox()

	// получаем набор для флага необходимости активировать модуль мониторинга
	kitWndAPMonitoring := MakeSettingCheckboxKit("Need monitoring", albumPhotoMonitorParam.NeedMonitoring)

	// получаем набор для интервала между запусками функции мониторинга
	kitWndAPInterval := MakeSettingSpinboxKit("Interval", 5, 21600, albumPhotoMonitorParam.Interval)

	// получаем набор для идентификатора получателя сообщений
	kitWndAPSendTo := MakeSettingEntryKit("Send to", strconv.Itoa(albumPhotoMonitorParam.SendTo))
	// и привязываем его к обработке ввода, чтобы кроме чисел ничего не вводилось
	kitWndAPSendTo.Entry.OnChanged(func(*ui.Entry) {
		NumericEntriesHandler(kitWndAPSendTo.Entry)
	})

	// получаем набор для количества проверяемых фото
	kitWndApPhotosCount := MakeSettingSpinboxKit("Photos count", 1, 1000, albumPhotoMonitorParam.PhotosCount)

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
	kitButtonsAP := MakeSettingButtonsKit()

	// привязываем к кнопкам соответствующие процедуры
	kitButtonsAP.ButtonCancel.OnClicked(func(*ui.Button) {
		kitWindowAlbumPhotoSettings.Window.Disable()
		kitWindowAlbumPhotoSettings.Window.Hide()
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
		if len(kitWndAPSendTo.Entry.Text()) == 0 {
			warningTitle := "Field \"Send to\" must not be empty."
			ShowWarningWindow(warningTitle)
			return
		}
		updatedAlbumPhotoMonitorParam.SendTo, err = strconv.Atoi(kitWndAPSendTo.Entry.Text())
		if err != nil {
			ToLogFile(err.Error(), string(debug.Stack()))
			panic(err.Error())
		}
		updatedAlbumPhotoMonitorParam.Interval = kitWndAPInterval.Spinbox.Value()
		updatedAlbumPhotoMonitorParam.LastDate = albumPhotoMonitorParam.LastDate
		updatedAlbumPhotoMonitorParam.PhotosCount = kitWndApPhotosCount.Spinbox.Value()

		err = UpdateDBAlbumPhotoMonitor(updatedAlbumPhotoMonitorParam)
		if err != nil {
			ToLogFile(err.Error(), string(debug.Stack()))
			panic(err.Error())
		}

		kitWindowAlbumPhotoSettings.Window.Disable()
		kitWindowAlbumPhotoSettings.Window.Hide()
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
		ToLogFile(err.Error(), string(debug.Stack()))
		panic(err.Error())
	}

	// получаем набор для отображения установок модуля мониторинга видео
	windowTitle := fmt.Sprintf("%v settings for %v", btnName, nameSubject)
	kitWindowVideoSettings := MakeSettingWindowKit(windowTitle, 300, 100)

	boxWndV := ui.NewVerticalBox()

	// получаем набор для флага необходимости активировать модуль мониторинга
	kitWndVMonitoring := MakeSettingCheckboxKit("Need monitoring", videoMonitorParam.NeedMonitoring)

	// получаем набор для интервала между запусками функции мониторинга
	kitWndVInterval := MakeSettingSpinboxKit("Interval", 5, 21600, videoMonitorParam.Interval)

	// получаем набор для идентификатора получателя сообщений
	kitWndVSendTo := MakeSettingEntryKit("Send to", strconv.Itoa(videoMonitorParam.SendTo))
	// и привязываем его к обработке ввода, чтобы кроме чисел ничего не вводилось
	kitWndVSendTo.Entry.OnChanged(func(*ui.Entry) {
		NumericEntriesHandler(kitWndVSendTo.Entry)
	})

	// получаем набор для количества проверяемых видео
	kitWndVVideoCount := MakeSettingSpinboxKit("Video count", 1, 1000, videoMonitorParam.VideoCount)

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
	kitButtonsV := MakeSettingButtonsKit()

	// привязываем к кнопкам соответствующие процедуры
	kitButtonsV.ButtonCancel.OnClicked(func(*ui.Button) {
		kitWindowVideoSettings.Window.Disable()
		kitWindowVideoSettings.Window.Hide()
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
		if len(kitWndVSendTo.Entry.Text()) == 0 {
			warningTitle := "Field \"Send to\" must not be empty."
			ShowWarningWindow(warningTitle)
			return
		}
		updatedVideoMonitorParam.SendTo, err = strconv.Atoi(kitWndVSendTo.Entry.Text())
		if err != nil {
			ToLogFile(err.Error(), string(debug.Stack()))
			panic(err.Error())
		}
		updatedVideoMonitorParam.Interval = kitWndVInterval.Spinbox.Value()
		updatedVideoMonitorParam.LastDate = videoMonitorParam.LastDate
		updatedVideoMonitorParam.VideoCount = kitWndVVideoCount.Spinbox.Value()

		err = UpdateDBVideoMonitor(updatedVideoMonitorParam)
		if err != nil {
			ToLogFile(err.Error(), string(debug.Stack()))
			panic(err.Error())
		}

		kitWindowVideoSettings.Window.Disable()
		kitWindowVideoSettings.Window.Hide()
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
		ToLogFile(err.Error(), string(debug.Stack()))
		panic(err.Error())
	}

	// получаем набор для отображения установок модуля мониторинга комментариев под фотками
	windowTitle := fmt.Sprintf("%v settings for %v", btnName, nameSubject)
	kitWindowPhotoCommentSettings := MakeSettingWindowKit(windowTitle, 300, 100)

	boxWndPC := ui.NewVerticalBox()

	// получаем набор для флага необходимости активировать модуль мониторинга
	kitWndPCMonitoring := MakeSettingCheckboxKit("Need monitoring", photoCommentMonitorParam.NeedMonitoring)

	// получаем набор для интервала между запусками функции мониторинга
	kitWndPCInterval := MakeSettingSpinboxKit("Interval", 5, 21600, photoCommentMonitorParam.Interval)

	// получаем набор для идентификатора получателя сообщений
	kitWndPCSendTo := MakeSettingEntryKit("Send to", strconv.Itoa(photoCommentMonitorParam.SendTo))
	// и привязываем его к обработке ввода, чтобы кроме чисел ничего не вводилось
	kitWndPCSendTo.Entry.OnChanged(func(*ui.Entry) {
		NumericEntriesHandler(kitWndPCSendTo.Entry)
	})

	// получаем набор для количества проверяемых комментариев
	kitWndPCCommentsCount := MakeSettingSpinboxKit("Comments count", 1, 1000, photoCommentMonitorParam.CommentsCount)

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
	kitButtonsPC := MakeSettingButtonsKit()

	// привязываем к кнопкам соответствующие процедуры
	kitButtonsPC.ButtonCancel.OnClicked(func(*ui.Button) {
		kitWindowPhotoCommentSettings.Window.Disable()
		kitWindowPhotoCommentSettings.Window.Hide()
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
		if len(kitWndPCSendTo.Entry.Text()) == 0 {
			warningTitle := "Field \"Send to\" must not be empty."
			ShowWarningWindow(warningTitle)
			return
		}
		updatedPhotoCommentMonitorParam.SendTo, err = strconv.Atoi(kitWndPCSendTo.Entry.Text())
		if err != nil {
			ToLogFile(err.Error(), string(debug.Stack()))
			panic(err.Error())
		}
		updatedPhotoCommentMonitorParam.Interval = kitWndPCInterval.Spinbox.Value()
		updatedPhotoCommentMonitorParam.LastDate = photoCommentMonitorParam.LastDate
		updatedPhotoCommentMonitorParam.CommentsCount = kitWndPCCommentsCount.Spinbox.Value()

		err = UpdateDBPhotoCommentMonitor(updatedPhotoCommentMonitorParam)
		if err != nil {
			ToLogFile(err.Error(), string(debug.Stack()))
			panic(err.Error())
		}

		kitWindowPhotoCommentSettings.Window.Disable()
		kitWindowPhotoCommentSettings.Window.Hide()
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
		ToLogFile(err.Error(), string(debug.Stack()))
		panic(err.Error())
	}

	// получаем набор для отображения установок модуля мониторинга комментариев в обсуждениях
	windowTitle := fmt.Sprintf("%v settings for %v", btnName, nameSubject)
	kitWindowVideoCommentSettings := MakeSettingWindowKit(windowTitle, 300, 100)

	boxWndVC := ui.NewVerticalBox()

	// получаем набор для флага необходимости активировать модуль мониторинга
	kitWndVCMonitoring := MakeSettingCheckboxKit("Need monitoring", videoCommentMonitorParam.NeedMonitoring)

	// получаем набор для интервала между запусками функции мониторинга
	kitWndVCInterval := MakeSettingSpinboxKit("Interval", 5, 21600, videoCommentMonitorParam.Interval)

	// получаем набор для идентификатора получателя сообщений
	kitWndVCSendTo := MakeSettingEntryKit("Send to", strconv.Itoa(videoCommentMonitorParam.SendTo))
	// и привязываем его к обработке ввода, чтобы кроме чисел ничего не вводилось
	kitWndVCSendTo.Entry.OnChanged(func(*ui.Entry) {
		NumericEntriesHandler(kitWndVCSendTo.Entry)
	})

	// получаем набор для количества проверяемых видео
	kitWndVCVideosCount := MakeSettingSpinboxKit("Videos count", 1, 200, videoCommentMonitorParam.VideosCount)

	// получаем набор для количества проверяемых комментариев
	kitWndVCCommentsCount := MakeSettingSpinboxKit("Comments count", 1, 100, videoCommentMonitorParam.CommentsCount)

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
	kitButtonsVC := MakeSettingButtonsKit()

	// привязываем к кнопкам соответствующие процедуры
	kitButtonsVC.ButtonCancel.OnClicked(func(*ui.Button) {
		kitWindowVideoCommentSettings.Window.Disable()
		kitWindowVideoCommentSettings.Window.Hide()
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
		if len(kitWndVCSendTo.Entry.Text()) == 0 {
			warningTitle := "Field \"Send to\" must not be empty."
			ShowWarningWindow(warningTitle)
			return
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
			ToLogFile(err.Error(), string(debug.Stack()))
			panic(err.Error())
		}

		kitWindowVideoCommentSettings.Window.Disable()
		kitWindowVideoCommentSettings.Window.Hide()
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
		ToLogFile(err.Error(), string(debug.Stack()))
		panic(err.Error())
	}

	// получаем набор для отображения установок модуля мониторинга комментариев в обсуждениях
	windowTitle := fmt.Sprintf("%v settings for %v", btnName, nameSubject)
	kitWindowTopicSettings := MakeSettingWindowKit(windowTitle, 300, 100)

	boxWndT := ui.NewVerticalBox()

	// получаем набор для флага необходимости активировать модуль мониторинга
	kitWndTMonitoring := MakeSettingCheckboxKit("Need monitoring", topicMonitorParam.NeedMonitoring)

	// получаем набор для интервала между запусками функции мониторинга
	kitWndTInterval := MakeSettingSpinboxKit("Interval", 5, 21600, topicMonitorParam.Interval)

	// получаем набор для идентификатора получателя сообщений
	kitWndTSendTo := MakeSettingEntryKit("Send to", strconv.Itoa(topicMonitorParam.SendTo))
	// и привязываем его к обработке ввода, чтобы кроме чисел ничего не вводилось
	kitWndTSendTo.Entry.OnChanged(func(*ui.Entry) {
		NumericEntriesHandler(kitWndTSendTo.Entry)
	})

	// получаем набор для количества проверяемых топиков обсуждений
	kitWndTTopicsCount := MakeSettingSpinboxKit("Topics count", 1, 100, topicMonitorParam.TopicsCount)

	// получаем набор для количества проверяемых комментариев
	kitWndTCommentsCount := MakeSettingSpinboxKit("Comments count", 1, 100, topicMonitorParam.TopicsCount)

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
	kitButtonsT := MakeSettingButtonsKit()

	// привязываем к кнопкам соответствующие процедуры
	kitButtonsT.ButtonCancel.OnClicked(func(*ui.Button) {
		kitWindowTopicSettings.Window.Disable()
		kitWindowTopicSettings.Window.Hide()
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
		if len(kitWndTSendTo.Entry.Text()) == 0 {
			warningTitle := "Field \"Send to\" must not be empty."
			ShowWarningWindow(warningTitle)
			return
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
			ToLogFile(err.Error(), string(debug.Stack()))
			panic(err.Error())
		}

		kitWindowTopicSettings.Window.Disable()
		kitWindowTopicSettings.Window.Hide()
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
		ToLogFile(err.Error(), string(debug.Stack()))
		panic(err.Error())
	}

	// получаем набор для отображения установок модуля мониторинга комментариев под постами
	windowTitle := fmt.Sprintf("%v settings for %v", btnName, nameSubject)
	kitWindowWallPostCommentSettings := MakeSettingWindowKit(windowTitle, 300, 100)

	boxWndWPC := ui.NewVerticalBox()

	// получаем набор для флага необходимости активировать модуль мониторинга
	kitWndWPCMonitoring := MakeSettingCheckboxKit("Need monitoring", wallPostCommentMonitorParam.NeedMonitoring)

	// получаем набор для интервала между запусками функции мониторинга
	kitWndWPCInterval := MakeSettingSpinboxKit("Interval", 5, 21600, wallPostCommentMonitorParam.Interval)

	// получаем набор для идентификатора получателя сообщений
	kitWndWPCSendTo := MakeSettingEntryKit("Send to", strconv.Itoa(wallPostCommentMonitorParam.SendTo))
	// и привязываем его к обработке ввода, чтобы кроме чисел ничего не вводилось
	kitWndWPCSendTo.Entry.OnChanged(func(*ui.Entry) {
		NumericEntriesHandler(kitWndWPCSendTo.Entry)
	})

	// получаем набор для количества проверяемых постов
	kitWndWPCPostsCount := MakeSettingSpinboxKit("Posts count", 1, 100, wallPostCommentMonitorParam.PostsCount)

	// получаем набор для количества проверяемых комментариев
	kitWndWPCCommentsCount := MakeSettingSpinboxKit("Comments count", 1, 100, wallPostCommentMonitorParam.CommentsCount)

	// получаем набор для фильтров постов для проверки комментариев
	listPostsFilters := []string{"all", "others", "owner"}
	kitWndWPCFilter := MakeSettingComboboxKit("Filter", listPostsFilters, wallPostCommentMonitorParam.Filter)

	// получаем набор для флага необходимости проверять все комментарии без исключения
	kitWndWPCMonitoringAll := MakeSettingCheckboxKit("Monitoring all", wallPostCommentMonitorParam.MonitoringAll)

	// получаем набор для флага необходимости проверять комментарии от сообществ
	kitWndWPCMonitorByCommunity := MakeSettingCheckboxKit("Monitor by community", wallPostCommentMonitorParam.MonitorByCommunity)

	// получаем набор для списка ключевых слов для поиска комментариев
	kitWndWPCKeywordsForMonitoring := MakeSettingEntryListKit("Keywords for monitoring", wallPostCommentMonitorParam.KeywordsForMonitoring)

	// получаем набор для списка комментариев для поиска
	kitWndWPCSmallCommentsForMonitoring := MakeSettingEntryListKit("Small comments for monitoring", wallPostCommentMonitorParam.SmallCommentsForMonitoring)

	// получаем набор для списка имен и фамилий авторов комментариев для поиска комментариев
	kitWndWPCUsersNamesForMonitoring := MakeSettingEntryListKit("Users names for monitoring", wallPostCommentMonitorParam.UsersNamesForMonitoring)

	// получаем набор для списка идентификаторов авторов комментариев для поиска комментариев
	kitWndWPCUsersIdsForMonitoring := MakeSettingEntryListKit("Users IDs for monitoring", wallPostCommentMonitorParam.UsersIDsForMonitoring)

	// получаем набор для списка идентификаторов авторов комментариев для их игнорирования при проверке комментариев
	kitWndWPCUsersIdsForIgnore := MakeSettingEntryListKit("Users IDs for ignore", wallPostCommentMonitorParam.UsersIDsForIgnore)

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
	kitButtonsWPC := MakeSettingButtonsKit()

	// привязываем к кнопкам соответствующие процедуры
	kitButtonsWPC.ButtonCancel.OnClicked(func(*ui.Button) {
		kitWindowWallPostCommentSettings.Window.Disable()
		kitWindowWallPostCommentSettings.Window.Hide()
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
		updatedWallPostCommentMonitorParam.Interval = kitWndWPCInterval.Spinbox.Value()
		if len(kitWndWPCSendTo.Entry.Text()) == 0 {
			warningTitle := "Field \"Send to\" must not be empty."
			ShowWarningWindow(warningTitle)
			return
		}
		updatedWallPostCommentMonitorParam.SendTo, err = strconv.Atoi(kitWndWPCSendTo.Entry.Text())
		if err != nil {
			ToLogFile(err.Error(), string(debug.Stack()))
			panic(err.Error())
		}
		listPostsFilters := []string{"all", "others", "owner"}
		if kitWndWPCFilter.Combobox.Selected() == -1 {
			warningTitle := "You must select an item in the combobox \"Filter\""
			ShowWarningWindow(warningTitle)
			return
		}
		updatedWallPostCommentMonitorParam.Filter = listPostsFilters[kitWndWPCFilter.Combobox.Selected()]
		updatedWallPostCommentMonitorParam.LastDate = wallPostCommentMonitorParam.LastDate
		// TODO: проверка соответствия оформления требованиям json
		jsonDump = fmt.Sprintf("{\"list\":[%v]}", kitWndWPCKeywordsForMonitoring.Entry.Text())
		updatedWallPostCommentMonitorParam.KeywordsForMonitoring = jsonDump
		// TODO: проверка соответствия оформления требованиям json
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
			ToLogFile(err.Error(), string(debug.Stack()))
			panic(err.Error())
		}
		kitWindowWallPostCommentSettings.Window.Disable()
		kitWindowWallPostCommentSettings.Window.Hide()
	})

	// добавляем коробку с кнопками на основную коробку окна
	kitWindowWallPostCommentSettings.Box.Append(kitButtonsWPC.Box, true)
	// затем еще одну коробку, для выравнивания расположения кнопок при растягивании окна
	boxWndWPCBottom := ui.NewHorizontalBox()
	kitWindowWallPostCommentSettings.Box.Append(boxWndWPCBottom, false)

	kitWindowWallPostCommentSettings.Window.Show()
}
