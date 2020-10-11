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
	boxSettings.SetPadded(true)

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
	groupGeneralSettings.SetTitle("Ключи доступа")
	groupGeneralSettings.SetChild(boxAccessTokensSettings)
	groupAdditionalSettings.SetChild(ui.NewLabel("Тут нечего отображать..."))

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
	btnAccessTokensSettings := ui.NewButton("Ключи доступа")
	// по умолчанию делаем ее неактивной
	btnAccessTokensSettings.Disable()
	// описываем кнопку для отображения установок субъектов
	btnSubjectsSettings := ui.NewButton("Субъекты")

	// привязываем кнопки к процедурам отображения соответствующих блоков настроек
	btnAccessTokensSettings.OnClicked(func(*ui.Button) {
		groupsSettingsData.Additional.SetChild(ui.NewLabel("Тут нечего показывать..."))
		groupsSettingsData.Additional.SetTitle("")
		groupsSettingsData.General.SetTitle("Ключи доступа")
		groupsSettingsData.General.SetChild(generalBoxesData.AccessTokens)
		btnAccessTokensSettings.Disable()
		if !(btnSubjectsSettings.Enabled()) {
			btnSubjectsSettings.Enable()
		}
	})
	btnSubjectsSettings.OnClicked(func(*ui.Button) {
		groupsSettingsData.General.SetTitle("Субъекты")
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
	var dbKit DataBaseKit
	accessTokens, err := dbKit.selectTableAccessToken()
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
		btnAccessTokenSettings := ui.NewButton("Настройки...")
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
	windowTitle := fmt.Sprintf("Добавление нового ключа доступа")
	kitWindowAccessTokenAddition := MakeSettingWindowKit(windowTitle, 300, 100)

	boxWndATAddition := ui.NewVerticalBox()

	// получаем набор для ввода названия нового токена доступа
	kitATCreationName := MakeSettingEntryKit("Название", "")

	// получаем набор для ввода значения нового токена доступа
	kitATCreationValue := MakeSettingEntryKit("Значение", "")

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

		err := accessToken.insertIntoDB()
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
	var dbKit DataBaseKit
	accessTokens, err := dbKit.selectTableAccessToken()
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
			windowTitle := "Настройки для ключа доступа " + accessToken.Name + ""
			kitWindowAccessTokenSettings.Window.SetTitle(windowTitle)

			boxWndAT := ui.NewVerticalBox()

			// получаем набор для названия токена доступа
			kitWndATName := MakeSettingEntryKit("Название", accessToken.Name)
			// TODO: добавить обработчик, чтобы исключить ввод пробелов

			// получаем набор для значения токена доступа
			kitWndATValue := MakeSettingEntryKit("Значение", accessToken.Value)

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
				// TODO: добавить проверку уникальности названия
				updatedAccessToken.Value = kitWndATValue.Entry.Text()

				err := updatedAccessToken.updateInDB()
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
	var dbKit DataBaseKit
	subjects, err := dbKit.selectTableSubject()
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
	var dbKit DataBaseKit
	accessTokens, err := dbKit.selectTableAccessToken()
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
	windowTitle := fmt.Sprintf("Добавление нового субъекта")
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
	kitSAdditionName := MakeSettingEntryKit("Название", "")

	// получаем набор для ввода идентификатора субъекта в базе данных ВК
	kitSAdditionSubjectID := MakeSettingEntryKit("Идентификатор в ВК", "")
	// и привязываем его к обработке ввода, чтобы кроме чисел ничего не вводилось
	kitSAdditionSubjectID.Entry.OnChanged(func(*ui.Entry) {
		NumericEntriesHandler(kitSAdditionSubjectID.Entry)
	})

	// описываем группу для общих установок субъекта
	groupSAdditionGeneral := ui.NewGroup("Общие")
	groupSAdditionGeneral.SetMargined(true)
	boxSAdditionGeneral := ui.NewVerticalBox()
	boxSAdditionGeneral.Append(kitSAdditionName.Box, false)
	boxSAdditionGeneral.Append(kitSAdditionSubjectID.Box, false)
	groupSAdditionGeneral.SetChild(boxSAdditionGeneral)

	// получаем набор для ввода идентификатора получателя сообщений в модуле wall_post_monitor
	kitSAdditionSendToinWPM := MakeSettingEntryKit("Получатель", "")
	// и привязываем его к обработке ввода, чтобы кроме чисел ничего не вводилось
	kitSAdditionSendToinWPM.Entry.OnChanged(func(*ui.Entry) {
		NumericEntriesHandler(kitSAdditionSendToinWPM.Entry)
	})

	// получаем набор для выбора токена доступа для метода wall.get в модуле wall_post_monitor
	kitSAdditionWGinWPM := MakeSettingComboboxKit("Ключ доступа для \"wall.get\"", accessTokensNames, "")

	// получаем набор для выбора токена доступа для метода users.get в модуле wall_post_monitor
	kitSAdditionUGinWPM := MakeSettingComboboxKit("Ключ доступа для \"users.get\"", accessTokensNames, "")

	// получаем набор для выбора токена доступа для метода groups.getById в модуле wall_post_monitor
	kitSAdditionGGBIinWPM := MakeSettingComboboxKit("Ключ доступа для \"groups.getById\"", accessTokensNames, "")

	// получаем набор для выбора токена доступа для метода messages.send в модуле wall_post_monitor
	kitSAdditionMSinWPM := MakeSettingComboboxKit("Ключ доступа для \"messages.send\"", accessTokensNames, "")

	// описываем группу для установок модуля wall_post_monitor субъекта
	groupSAdditionWPM := ui.NewGroup("Посты на стене")
	groupSAdditionWPM.SetMargined(true)
	boxSAdditionWPM := ui.NewVerticalBox()
	boxSAdditionWPM.Append(kitSAdditionSendToinWPM.Box, false)
	boxSAdditionWPM.Append(kitSAdditionWGinWPM.Box, false)
	boxSAdditionWPM.Append(kitSAdditionUGinWPM.Box, false)
	boxSAdditionWPM.Append(kitSAdditionGGBIinWPM.Box, false)
	boxSAdditionWPM.Append(kitSAdditionMSinWPM.Box, false)
	groupSAdditionWPM.SetChild(boxSAdditionWPM)

	// получаем набор для ввода идентификатора получателя сообщений в модуле album_photo_monitor
	kitSAdditionSendToinAPM := MakeSettingEntryKit("Получатель", "")
	// и привязываем его к обработке ввода, чтобы кроме чисел ничего не вводилось
	kitSAdditionSendToinAPM.Entry.OnChanged(func(*ui.Entry) {
		NumericEntriesHandler(kitSAdditionSendToinAPM.Entry)
	})

	// получаем набор для выбора токена доступа для метода photos.get в модуле album_photo_monitor
	kitSAdditionPGinAPM := MakeSettingComboboxKit("Ключ доступа для \"photos.get\"", accessTokensNames, "")

	// получаем набор для выбора токена доступа для метода photos.getAlbums в модуле album_photo_monitor
	kitSAdditionPGAinAPM := MakeSettingComboboxKit("Ключ доступа для \"photos.getAlbums\"", accessTokensNames, "")

	// получаем набор для выбора токена доступа для метода users.get в модуле album_photo_monitor
	kitSAdditionUGinAPM := MakeSettingComboboxKit("Ключ доступа для \"users.get\"", accessTokensNames, "")

	// получаем набор для выбора токена доступа для метода groups.getById в модуле album_photo_monitor
	kitSAdditionGGBIinAPM := MakeSettingComboboxKit("Ключ доступа для \"groups.getById\"", accessTokensNames, "")

	// получаем набор для выбора токена доступа для метода messages.send в модуле album_photo_monitor
	kitSAdditionMSinAPM := MakeSettingComboboxKit("Ключ доступа для \"messages.send\"", accessTokensNames, "")

	// описываем группу для установок модуля album_photo_monitor субъекта
	groupSAdditionAPM := ui.NewGroup("Фото в альбомах")
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
	kitSAdditionSendToinVM := MakeSettingEntryKit("Получатель", "")
	// и привязываем его к обработке ввода, чтобы кроме чисел ничего не вводилось
	kitSAdditionSendToinVM.Entry.OnChanged(func(*ui.Entry) {
		NumericEntriesHandler(kitSAdditionSendToinVM.Entry)
	})

	// получаем набор для выбора токена доступа для метода video.get в модуле video_monitor
	kitSAdditionVGinVM := MakeSettingComboboxKit("Ключ доступа для \"video.get\"", accessTokensNames, "")

	// получаем набор для выбора токена доступа для метода users.get в модуле video_monitor
	kitSAdditionUGinVM := MakeSettingComboboxKit("Ключ доступа для \"users.get\"", accessTokensNames, "")

	// получаем набор для выбора токена доступа для метода groups.getById в модуле video_monitor
	kitSAdditionGGBIinVM := MakeSettingComboboxKit("Ключ доступа для \"groups.getById\"", accessTokensNames, "")

	// получаем набор для выбора токена доступа для метода messages.send в модуле video_monitor
	kitSAdditionMSinVM := MakeSettingComboboxKit("Ключ доступа для \"messages.send\"", accessTokensNames, "")

	// описываем группу для установок модуля video_monitor субъекта
	groupSAdditionVM := ui.NewGroup("Видео в альбомах")
	groupSAdditionVM.SetMargined(true)
	boxSAdditionVM := ui.NewVerticalBox()
	boxSAdditionVM.Append(kitSAdditionSendToinVM.Box, false)
	boxSAdditionVM.Append(kitSAdditionVGinVM.Box, false)
	boxSAdditionVM.Append(kitSAdditionUGinVM.Box, false)
	boxSAdditionVM.Append(kitSAdditionGGBIinVM.Box, false)
	boxSAdditionVM.Append(kitSAdditionMSinVM.Box, false)
	groupSAdditionVM.SetChild(boxSAdditionVM)

	// получаем набор для ввода идентификатора получателя сообщений в модуле photo_comment_monitor
	kitSAdditionSendToinPCM := MakeSettingEntryKit("Получатель", "")
	// и привязываем его к обработке ввода, чтобы кроме чисел ничего не вводилось
	kitSAdditionSendToinPCM.Entry.OnChanged(func(*ui.Entry) {
		NumericEntriesHandler(kitSAdditionSendToinPCM.Entry)
	})

	// получаем набор для выбора токена доступа для метода photos.getAllComments в модуле photo_comment_monitor
	kitSAdditionPGACinPCM := MakeSettingComboboxKit("Ключ доступа для \"photos.getAllComments\"", accessTokensNames, "")

	// получаем набор для выбора токена доступа для метода users.get в модуле photo_comment_monitor
	kitSAdditionUGinPCM := MakeSettingComboboxKit("Ключ доступа для \"users.get\"", accessTokensNames, "")

	// получаем набор для выбора токена доступа для метода groups.getById в модуле photo_comment_monitor
	kitSAdditionGGBIinPCM := MakeSettingComboboxKit("Ключ доступа для \"groups.getById\"", accessTokensNames, "")

	// получаем набор для выбора токена доступа для метода messages.send в модуле photo_comment_monitor
	kitSAdditionMSinPCM := MakeSettingComboboxKit("Ключ доступа для \"messages.send\"", accessTokensNames, "")

	// описываем группу для установок модуля photo_comment_monitor субъекта
	groupSAdditionPCM := ui.NewGroup("Комментарии под фото")
	groupSAdditionPCM.SetMargined(true)
	boxSAdditionPCM := ui.NewVerticalBox()
	boxSAdditionPCM.Append(kitSAdditionSendToinPCM.Box, false)
	boxSAdditionPCM.Append(kitSAdditionPGACinPCM.Box, false)
	boxSAdditionPCM.Append(kitSAdditionUGinPCM.Box, false)
	boxSAdditionPCM.Append(kitSAdditionGGBIinPCM.Box, false)
	boxSAdditionPCM.Append(kitSAdditionMSinPCM.Box, false)
	groupSAdditionPCM.SetChild(boxSAdditionPCM)

	// получаем набор для ввода идентификатора получателя сообщений в модуле video_comment_monitor
	kitSAdditionSendToinVCM := MakeSettingEntryKit("Получатель", "")
	// и привязываем его к обработке ввода, чтобы кроме чисел ничего не вводилось
	kitSAdditionSendToinVCM.Entry.OnChanged(func(*ui.Entry) {
		NumericEntriesHandler(kitSAdditionSendToinVCM.Entry)
	})

	// получаем набор для выбора токена доступа для метода video.getComments в модуле video_comment_monitor
	kitSAdditionVGCinVCM := MakeSettingComboboxKit("Ключ доступа для \"video.getComments\"", accessTokensNames, "")

	// получаем набор для выбора токена доступа для метода users.get в модуле video_comment_monitor
	kitSAdditionUGinVCM := MakeSettingComboboxKit("Ключ доступа для \"users.get\"", accessTokensNames, "")

	// получаем набор для выбора токена доступа для метода groups.getById в модуле video_comment_monitor
	kitSAdditionGGBIinVCM := MakeSettingComboboxKit("Ключ доступа для \"groups.getById\"", accessTokensNames, "")

	// получаем набор для выбора токена доступа для метода video.get в модуле video_comment_monitor
	kitSAdditionVGinVCM := MakeSettingComboboxKit("Ключ доступа для \"video.get\"", accessTokensNames, "")

	// получаем набор для выбора токена доступа для метода messages.send в модуле video_comment_monitor
	kitSAdditionMSinVCM := MakeSettingComboboxKit("Ключ доступа для \"messages.send\"", accessTokensNames, "")

	// описываем группу для установок модуля video_comment_monitor субъекта
	groupSAdditionVCM := ui.NewGroup("Комментарии под видео")
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
	kitSAdditionSendToinTM := MakeSettingEntryKit("Получатель", "")
	// и привязываем его к обработке ввода, чтобы кроме чисел ничего не вводилось
	kitSAdditionSendToinTM.Entry.OnChanged(func(*ui.Entry) {
		NumericEntriesHandler(kitSAdditionSendToinTM.Entry)
	})

	// получаем набор для выбора токена доступа для метода board.getComments в модуле topic_monitor
	kitSAdditionBGCinTM := MakeSettingComboboxKit("Ключ доступа для \"board.getComments\"", accessTokensNames, "")

	// получаем набор для выбора токена доступа для метода board.getTopics в модуле topic_monitor
	kitSAdditionBGTinTM := MakeSettingComboboxKit("Ключ доступа для \"board.getTopics\"", accessTokensNames, "")

	// получаем набор для выбора токена доступа для метода users.get в модуле topic_monitor
	kitSAdditionUGinTM := MakeSettingComboboxKit("Ключ доступа для \"users.get\"", accessTokensNames, "")

	// получаем набор для выбора токена доступа для метода groups.getById в модуле topic_monitor
	kitSAdditionGGBIinTM := MakeSettingComboboxKit("Ключ доступа для \"groups.getById\"", accessTokensNames, "")

	// получаем набор для выбора токена доступа для метода messages.send в модуле topic_monitor
	kitSAdditionMSinTM := MakeSettingComboboxKit("Ключ доступа для \"messages.send\"", accessTokensNames, "")

	// описываем группу для установок модуля topic_monitor субъекта
	groupSAdditionTM := ui.NewGroup("Комментарии в обсуждениях")
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
	kitSAdditionSendToinWPCM := MakeSettingEntryKit("Получатель", "")
	// и привязываем его к обработке ввода, чтобы кроме чисел ничего не вводилось
	kitSAdditionSendToinWPCM.Entry.OnChanged(func(*ui.Entry) {
		NumericEntriesHandler(kitSAdditionSendToinWPCM.Entry)
	})

	// получаем набор для выбора токена доступа для метода wall.getComments в модуле wall_post_comment_monitor
	kitSAdditionWGCsinWPCM := MakeSettingComboboxKit("Ключ доступа для \"wall.getComments\"", accessTokensNames, "")

	// получаем набор для выбора токена доступа для метода users.get в модуле wall_post_comment_monitor
	kitSAdditionUGinWPCM := MakeSettingComboboxKit("Ключ доступа для \"users.get\"", accessTokensNames, "")

	// получаем набор для выбора токена доступа для метода groups.getById в модуле wall_post_comment_monitor
	kitSAdditionGGBIinWPCM := MakeSettingComboboxKit("Ключ доступа для \"groups.getById\"", accessTokensNames, "")

	// получаем набор для выбора токена доступа для метода wall.get в модуле wall_post_comment_monitor
	kitSAdditionWGinWPCM := MakeSettingComboboxKit("Ключ доступа для \"wall.get\"", accessTokensNames, "")

	// получаем набор для выбора токена доступа для метода wall.getComment в модуле wall_post_comment_monitor
	kitSAdditionWGCinWPCM := MakeSettingComboboxKit("Ключ доступа для \"wall.getComment\"", accessTokensNames, "")

	// получаем набор для выбора токена доступа для метода messages.send в модуле wall_post_comment_monitor
	kitSAdditionMSinWPCM := MakeSettingComboboxKit("Ключ доступа для \"messages.send\"", accessTokensNames, "")

	// описываем группу для установок модуля wall_post_comment_monitor субъекта
	groupSAdditionWPCM := ui.NewGroup("Комментарии под постами")
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
			warningTitle := "Поле \"Идентификатор в ВК\" не должно быть пустым"
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
			warningTitle := "Поле \"Название\" не должно быть пустым"
			ShowWarningWindow(warningTitle)
			return
		}
		newSubjectData.Name = kitSAdditionName.Entry.Text()
		// TODO: добавить проверку уникальности названия

		var listNewMethodData ListNewMethodData
		var listNewMonitorModuleData ListNewMonitorModuleData

		monitorsNames := []string{"wall_post_monitor", "album_photo_monitor", "video_monitor",
			"photo_comment_monitor", "video_comment_monitor", "topic_monitor", "wall_post_comment_monitor"}

		for _, monitorName := range monitorsNames {
			switch monitorName {
			case "wall_post_monitor":

				if kitSAdditionWGinWPM.Combobox.Selected() == -1 {
					warningTitle := "Нужно выбрать элемент в списке " +
						"\"Ключ доступа для \"wall.get\"\""
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
					warningTitle := "Нужно выбрать элемент в списке " +
						"\"Ключ доступа для \"Users.get\"\""
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
					warningTitle := "Нужно выбрать элемент в списке " +
						"\"Ключ доступа для \"Groups.getById\"\""
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
					warningTitle := "Нужно выбрать элемент в списке " +
						"\"Ключ доступа для \"messages.send\"\""
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
					warningTitle := "Поле \"Получатель\" не должно быть пустым."
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
					warningTitle := "Нужно выбрать элемент в списке " +
						"\"Ключ доступа для \"photos.get\"\""
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
					warningTitle := "Нужно выбрать элемент в списке " +
						"\"Ключ доступа для \"photos.getAlbums\"\""
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
					warningTitle := "Нужно выбрать элемент в списке " +
						"\"Ключ доступа для \"users.get\"\""
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
					warningTitle := "Нужно выбрать элемент в списке " +
						"\"Ключ доступа для \"groups.getById\"\""
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
					warningTitle := "Нужно выбрать элемент в списке " +
						"\"Ключ доступа для \"messages.send\"\""
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
					warningTitle := "Поле \"Получатель\" не должно быть пустым."
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
					warningTitle := "Нужно выбрать элемент в списке " +
						"\"Ключ доступа для \"video.get\"\""
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
					warningTitle := "Нужно выбрать элемент в списке " +
						"\"Ключ доступа для \"users.get\"\""
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
					warningTitle := "Нужно выбрать элемент в списке " +
						"\"Ключ доступа для \"groups.getById\"\""
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
					warningTitle := "Нужно выбрать элемент в списке " +
						"\"Ключ доступа для \"messages.send\"\""
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
					warningTitle := "Поле \"Получатель\" не должно быть пустым."
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
					warningTitle := "Нужно выбрать элемент в списке " +
						"\"Ключ доступа для \"photos.getAllComments\"\""
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
					warningTitle := "Нужно выбрать элемент в списке " +
						"\"Ключ доступа для \"users.get\"\""
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
					warningTitle := "Нужно выбрать элемент в списке " +
						"\"Ключ доступа для \"groups.getById\"\""
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
					warningTitle := "Нужно выбрать элемент в списке " +
						"\"Ключ доступа для \"messages.send\"\""
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
					warningTitle := "Поле \"Получатель\" не должно быть пустым."
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
					warningTitle := "Нужно выбрать элемент в списке " +
						"\"Ключ доступа для \"video.getComments\"\""
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
					warningTitle := "Нужно выбрать элемент в списке " +
						"\"Ключ доступа для \"users.get\"\""
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
					warningTitle := "Нужно выбрать элемент в списке " +
						"\"Ключ доступа для \"groups.getById\"\""
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
					warningTitle := "Нужно выбрать элемент в списке " +
						"\"Ключ доступа для \"video.get\"\""
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
					warningTitle := "Нужно выбрать элемент в списке " +
						"\"Ключ доступа для \"messages.send\"\""
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
					warningTitle := "Поле \"Получатель\" не должно быть пустым."
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
					warningTitle := "Нужно выбрать элемент в списке " +
						"\"Ключ доступа для \"board.getComments\"\""
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
					warningTitle := "Нужно выбрать элемент в списке " +
						"\"Ключ доступа для \"board.getTopicsgetComments\"\""
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
					warningTitle := "Нужно выбрать элемент в списке " +
						"\"Ключ доступа для \"users.get\"\""
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
					warningTitle := "Нужно выбрать элемент в списке " +
						"\"Ключ доступа для \"groups.getById\"\""
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
					warningTitle := "Нужно выбрать элемент в списке " +
						"\"Ключ доступа для \"messages.send\"\""
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
					warningTitle := "Поле \"Получатель\" не должно быть пустым."
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
					warningTitle := "Нужно выбрать элемент в списке " +
						"\"Ключ доступа для \"wall.getComments\"\""
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
					warningTitle := "Нужно выбрать элемент в списке " +
						"\"Ключ доступа для \"users.get\"\""
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
					warningTitle := "Нужно выбрать элемент в списке " +
						"\"Ключ доступа для \"groups.getById\"\""
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
					warningTitle := "Нужно выбрать элемент в списке " +
						"\"Ключ доступа для \"wall.get\"\""
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
					warningTitle := "Нужно выбрать элемент в списке " +
						"\"Ключ доступа для \"wall.getComment\"\""
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
					warningTitle := "Нужно выбрать элемент в списке " +
						"\"Ключ доступа для \"messages.send\"\""
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
					warningTitle := "Поле \"Получатель\" не должно быть пустым."
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

		var dbKit DataBaseKit
		subjects, err := dbKit.selectTableSubject()
		if err != nil {
			ToLogFile(err.Error(), string(debug.Stack()))
			panic(err.Error())
		}
		id := subjects[len(subjects)-1].ID

		for _, monitorName := range monitorsNames {
			switch monitorName {
			case "wall_post_monitor":
				additionNewMonitor(monitorName, id)

				var monitor Monitor
				err := monitor.selectFromDBByNameAndBySubjectID(monitorName, id)
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

				var monitor Monitor
				err := monitor.selectFromDBByNameAndBySubjectID(monitorName, id)
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

				var monitor Monitor
				err := monitor.selectFromDBByNameAndBySubjectID(monitorName, id)
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

				var monitor Monitor
				err := monitor.selectFromDBByNameAndBySubjectID(monitorName, id)
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

				var monitor Monitor
				err := monitor.selectFromDBByNameAndBySubjectID(monitorName, id)
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

				var monitor Monitor
				err := monitor.selectFromDBByNameAndBySubjectID(monitorName, id)
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

				var monitor Monitor
				err := monitor.selectFromDBByNameAndBySubjectID(monitorName, id)
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

	err := subject.insertIntoDB()
	if err != nil {
		ToLogFile(err.Error(), string(debug.Stack()))
		panic(err.Error())
	}
}

func additionNewMonitor(monitorName string, subjectID int) {
	var monitor Monitor

	monitor.Name = monitorName
	monitor.SubjectID = subjectID

	err := monitor.insertIntoDB()
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

	err := method.insertIntoDB()
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

	err := wallPostMonitorParam.insertIntoDB()
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

	err := albumPhotoMonitorParam.insertIntoDB()
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

	err := videoMonitorParam.insertIntoDB()
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

	err := photoCommentMonitorParam.insertIntoDB()
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

	err := videoCommentMonitorParam.insertIntoDB()
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

	err := topicMonitorParam.insertIntoDB()
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

	err := wallPostCommentMonitorParam.insertIntoDB()
	if err != nil {
		ToLogFile(err.Error(), string(debug.Stack()))
		panic(err.Error())
	}
}

func makeSubjectAdditionalSettingsBox(subjectData Subject) *ui.Box {
	boxSubjectAdditionalSettingsBox := ui.NewVerticalBox()

	// создаем список с названиями кнопок для вызова окна доп. с установками
	btnsNames := []string{"Общие", "Наблюдатель постов на стене", "Наблюдатель фото в альбомах",
		"Наблюдатель видео в альбомах",
		"Наблюдатель комментариев под фото", "Наблюдатель комментариев под видео",
		"Наблюдатель комментариев в обсуждениях", "Наблюдатель комментариев под постами"}

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
		btnSettingsSection := ui.NewButton("Настройки...")
		boxBtnSettingsSection.Append(btnSettingsSection, false)
		boxSettingsSection.Append(boxBtnSettingsSection, true)

		// привязываем к кнопке отображения окна с доп. установками соответствующую процедуру
		switch btnName {
		case "Общие":
			btnSettingsSection.OnClicked(func(*ui.Button) {
				showSubjectGeneralSettingWindow(subjectData.ID, btnName)
			})
		case "Наблюдатель постов на стене":
			btnSettingsSection.OnClicked(func(*ui.Button) {
				showSubjectWallPostSettingWindow(subjectData.ID, subjectData.Name, btnName)
			})
		case "Наблюдатель фото в альбомах":
			btnSettingsSection.OnClicked(func(*ui.Button) {
				showSubjectAlbumPhotoSettingWindow(subjectData.ID, subjectData.Name, btnName)
			})
		case "Наблюдатель видео в альбомах":
			btnSettingsSection.OnClicked(func(*ui.Button) {
				showSubjectVideoSettingWindow(subjectData.ID, subjectData.Name, btnName)
			})
		case "Наблюдатель комментариев под фото":
			btnSettingsSection.OnClicked(func(*ui.Button) {
				showSubjectPhotoCommentSettingWindow(subjectData.ID, subjectData.Name, btnName)
			})
		case "Наблюдатель комментариев под видео":
			btnSettingsSection.OnClicked(func(*ui.Button) {
				showSubjectVideoCommentSettingWindow(subjectData.ID, subjectData.Name, btnName)
			})
		case "Наблюдатель комментариев в обсуждениях":
			btnSettingsSection.OnClicked(func(*ui.Button) {
				showSubjectTopicSettingWindow(subjectData.ID, subjectData.Name, btnName)
			})
		case "Наблюдатель комментариев под постами":
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
	var dbKit DataBaseKit
	subjects, err := dbKit.selectTableSubject()
	if err != nil {
		ToLogFile(err.Error(), string(debug.Stack()))
		panic(err.Error())
	}

	// запрашиваем список токенов доступа из базы данных
	accessTokens, err := dbKit.selectTableAccessToken()
	if err != nil {
		ToLogFile(err.Error(), string(debug.Stack()))
		panic(err.Error())
	}
	// формируем список с названиями токенов доступа
	var accessTokensNames []string
	for _, accessToken := range accessTokens {
		accessTokensNames = append(accessTokensNames, accessToken.Name)
	}

	// получаем набор для отображения общих установок субъекта мониторинга
	kitWindowGeneralSettings := MakeSettingWindowKit("", 300, 100)

	// перечисляем субъекты
	for _, subject := range subjects {
		// ищем субъект с подходящим идентификатором
		if subject.ID == IDSubject {

			// устанавливаем заголовок окна в соответствии с названием субъекта и назначением установок
			windowTitle := fmt.Sprintf("%v: настройки для %v", btnName, subject.Name)
			kitWindowGeneralSettings.Window.SetTitle(windowTitle)

			boxWndSMain := ui.NewHorizontalBox()
			boxWndSMain.SetPadded(true)

			boxWndSLeft := ui.NewVerticalBox()
			boxWndSLeft.SetPadded(true)
			boxWndSCenter := ui.NewVerticalBox()
			boxWndSCenter.SetPadded(true)
			boxWndSRight := ui.NewVerticalBox()
			boxWndSRight.SetPadded(true)

			// получаем набор для названия субъекта мониторинга
			kitWndSName := MakeSettingEntryKit("Название", subject.Name)
			// TODO: добавить обработчик, чтобы исключить ввод пробелов

			// получаем набор для идентификатора субъекта мониторинга в базе ВК
			kitWndSSubjectID := MakeSettingEntryKit("Идентификатор в ВК", strconv.Itoa(subject.SubjectID))
			// и привязываем его к обработке ввода, чтобы кроме чисел ничего не вводилось
			kitWndSSubjectID.Entry.OnChanged(func(*ui.Entry) {
				NumericEntriesHandler(kitWndSSubjectID.Entry)
			})

			// описываем коробку и группу для общих установок субъекта
			boxWndSGeneral := ui.NewVerticalBox()
			groupWndSGeneral := ui.NewGroup("Общие")
			groupWndSGeneral.SetMargined(true)
			boxWndSGeneral.Append(kitWndSName.Box, false)
			boxWndSGeneral.Append(kitWndSSubjectID.Box, false)
			groupWndSGeneral.SetChild(boxWndSGeneral)

			// запрашиваем из базы данных данные по модулю мониторинга wall_post_monitor соответствующего субъекта
			var monitorWPM Monitor
			monitorWPM.selectFromDBByNameAndBySubjectID("wall_post_monitor", subject.ID)

			// запрашиваем из базы данных данные по методу wall.get соответствующих субъекта и модуля мониторинга
			var methodWGforWPM Method
			methodWGforWPM.selectFromDBByNameAndBySubjectIDAndByMonitorID("wall.get", subject.ID, monitorWPM.ID)
			var currentATInWPMWallGet string
			for _, accessToken := range accessTokens {
				if methodWGforWPM.AccessTokenID == accessToken.ID {
					currentATInWPMWallGet = accessToken.Name
				}
			}
			// получаем набор для выбора токена доступа для метода wall.get в модуле wall_post_monitor
			kitSWGinWPM := MakeSettingComboboxKit("Ключ доступа для \"wall.get\"", accessTokensNames, currentATInWPMWallGet)

			// запрашиваем из базы данных данные по методу users.get соответствующих субъекта и модуля мониторинга
			var methodUGforWPM Method
			methodUGforWPM.selectFromDBByNameAndBySubjectIDAndByMonitorID("users.get", subject.ID, monitorWPM.ID)
			var currentATInWPMUsersGet string
			for _, accessToken := range accessTokens {
				if methodUGforWPM.AccessTokenID == accessToken.ID {
					currentATInWPMUsersGet = accessToken.Name
				}
			}
			// получаем набор для выбора токена доступа для метода users.get в модуле wall_post_monitor
			kitSUGinWPM := MakeSettingComboboxKit("Ключ доступа для \"users.get\"", accessTokensNames, currentATInWPMUsersGet)

			// запрашиваем из базы данных данные по методу groups.getById соответствующих субъекта и модуля мониторинга
			var methodGGBIforWPM Method
			methodGGBIforWPM.selectFromDBByNameAndBySubjectIDAndByMonitorID("groups.getById", subject.ID, monitorWPM.ID)
			var currentATInWPMGroupsGetByID string
			for _, accessToken := range accessTokens {
				if methodGGBIforWPM.AccessTokenID == accessToken.ID {
					currentATInWPMGroupsGetByID = accessToken.Name
				}
			}
			// получаем набор для выбора токена доступа для метода groups.getById в модуле wall_post_monitor
			kitSGGBIinWPM := MakeSettingComboboxKit("Ключ доступа для \"groups.getById\"", accessTokensNames, currentATInWPMGroupsGetByID)

			// запрашиваем из базы данных данные по методу messages.send соответствующих субъекта и модуля мониторинга
			var methodMSforWPM Method
			methodMSforWPM.selectFromDBByNameAndBySubjectIDAndByMonitorID("messages.send", subject.ID, monitorWPM.ID)
			var currentATInWPMMessagesSend string
			for _, accessToken := range accessTokens {
				if methodMSforWPM.AccessTokenID == accessToken.ID {
					currentATInWPMMessagesSend = accessToken.Name
				}
			}
			// получаем набор для выбора токена доступа для метода messages.send в модуле wall_post_monitor
			kitSMSinWPM := MakeSettingComboboxKit("Ключ доступа для \"messages.send\"", accessTokensNames, currentATInWPMMessagesSend)

			// описываем коробку и группу для установок модуля wall_post_monitor субъекта
			boxWndSWPM := ui.NewVerticalBox()
			groupWndSWPM := ui.NewGroup("Посты на стене")
			groupWndSWPM.SetMargined(true)
			boxWndSWPM.Append(kitSWGinWPM.Box, false)
			boxWndSWPM.Append(kitSUGinWPM.Box, false)
			boxWndSWPM.Append(kitSGGBIinWPM.Box, false)
			boxWndSWPM.Append(kitSMSinWPM.Box, false)
			groupWndSWPM.SetChild(boxWndSWPM)

			// запрашиваем из базы данных данные по модулю мониторинга album_photo_monitor соответствующего субъекта
			var monitorAPM Monitor
			monitorAPM.selectFromDBByNameAndBySubjectID("album_photo_monitor", subject.ID)

			// запрашиваем из базы данных данные по методу photos.get соответствующих субъекта и модуля мониторинга
			var methodPGforAPM Method
			methodPGforAPM.selectFromDBByNameAndBySubjectIDAndByMonitorID("photos.get", subject.ID, monitorAPM.ID)
			var currentATInAPMPhotosGet string
			for _, accessToken := range accessTokens {
				if methodPGforAPM.AccessTokenID == accessToken.ID {
					currentATInAPMPhotosGet = accessToken.Name
				}
			}
			// получаем набор для выбора токена доступа для метода photos.get в модуле album_photo_monitor
			kitSPGinAPM := MakeSettingComboboxKit("Ключ доступа для \"photos.get\"", accessTokensNames, currentATInAPMPhotosGet)

			// запрашиваем из базы данных данные по методу photos.getAlbums соответствующих субъекта и модуля мониторинга
			var methodPGAforAPM Method
			methodPGAforAPM.selectFromDBByNameAndBySubjectIDAndByMonitorID("photos.getAlbums", subject.ID, monitorAPM.ID)
			var currentATInAPMPhotosGetAlbums string
			for _, accessToken := range accessTokens {
				if methodPGAforAPM.AccessTokenID == accessToken.ID {
					currentATInAPMPhotosGetAlbums = accessToken.Name
				}
			}
			// получаем набор для выбора токена доступа для метода photos.getAlbums в модуле album_photo_monitor
			kitSPGAinAPM := MakeSettingComboboxKit("Ключ доступа для \"photos.getAlbums\"", accessTokensNames, currentATInAPMPhotosGetAlbums)

			// запрашиваем из базы данных данные по методу users.get соответствующих субъекта и модуля мониторинга
			var methodUGforAPM Method
			methodUGforAPM.selectFromDBByNameAndBySubjectIDAndByMonitorID("users.get", subject.ID, monitorAPM.ID)
			var currentATInAPMUsersGet string
			for _, accessToken := range accessTokens {
				if methodUGforAPM.AccessTokenID == accessToken.ID {
					currentATInAPMUsersGet = accessToken.Name
				}
			}
			// получаем набор для выбора токена доступа для метода users.get в модуле album_photo_monitor
			kitSUGinAPM := MakeSettingComboboxKit("Ключ доступа для \"users.get\"", accessTokensNames, currentATInAPMUsersGet)

			// запрашиваем из базы данных данные по методу groups.getById соответствующих субъекта и модуля мониторинга
			var methodGGBIforAPM Method
			methodGGBIforAPM.selectFromDBByNameAndBySubjectIDAndByMonitorID("groups.getById", subject.ID, monitorAPM.ID)
			var currentATInAPMGroupsGetByID string
			for _, accessToken := range accessTokens {
				if methodGGBIforAPM.AccessTokenID == accessToken.ID {
					currentATInAPMGroupsGetByID = accessToken.Name
				}
			}
			// получаем набор для выбора токена доступа для метода groups.getById в модуле album_photo_monitor
			kitSGGBIinAPM := MakeSettingComboboxKit("Ключ доступа для \"groups.getById\"", accessTokensNames, currentATInAPMGroupsGetByID)

			// запрашиваем из базы данных данные по методу messages.send соответствующих субъекта и модуля мониторинга
			var methodMSforAPM Method
			methodMSforAPM.selectFromDBByNameAndBySubjectIDAndByMonitorID("messages.send", subject.ID, monitorAPM.ID)
			var currentATInAPMMessagesSend string
			for _, accessToken := range accessTokens {
				if methodMSforAPM.AccessTokenID == accessToken.ID {
					currentATInAPMMessagesSend = accessToken.Name
				}
			}
			// получаем набор для выбора токена доступа для метода messages.send в модуле album_photo_monitor
			kitSMSinAPM := MakeSettingComboboxKit("Ключ доступа для \"messages.send\"", accessTokensNames, currentATInAPMMessagesSend)

			// описываем группу для токенов доступа модуля album_photo_monitor субъекта
			boxWndSAPM := ui.NewVerticalBox()
			groupWndSAPM := ui.NewGroup("Фото в альбомах")
			groupWndSAPM.SetMargined(true)
			boxWndSAPM.Append(kitSPGinAPM.Box, false)
			boxWndSAPM.Append(kitSPGAinAPM.Box, false)
			boxWndSAPM.Append(kitSUGinAPM.Box, false)
			boxWndSAPM.Append(kitSGGBIinAPM.Box, false)
			boxWndSAPM.Append(kitSMSinAPM.Box, false)
			groupWndSAPM.SetChild(boxWndSAPM)

			// запрашиваем из базы данных данные по модулю мониторинга video_monitor соответствующего субъекта
			var monitorVM Monitor
			monitorVM.selectFromDBByNameAndBySubjectID("video_monitor", subject.ID)

			// запрашиваем из базы данных данные по методу video.get соответствующих субъекта и модуля мониторинга
			var methodVGforVM Method
			methodVGforVM.selectFromDBByNameAndBySubjectIDAndByMonitorID("video.get", subject.ID, monitorVM.ID)
			var currentATInVMVideoGet string
			for _, accessToken := range accessTokens {
				if methodVGforVM.AccessTokenID == accessToken.ID {
					currentATInVMVideoGet = accessToken.Name
				}
			}
			// получаем набор для выбора токена доступа для метода video.get в модуле video_monitor
			kitSVGinVM := MakeSettingComboboxKit("Ключ доступа для \"video.get\"", accessTokensNames, currentATInVMVideoGet)

			// запрашиваем из базы данных данные по методу users.get соответствующих субъекта и модуля мониторинга
			var methodUGforVM Method
			methodUGforVM.selectFromDBByNameAndBySubjectIDAndByMonitorID("users.get", subject.ID, monitorVM.ID)
			var currentATInVMUsersGet string
			for _, accessToken := range accessTokens {
				if methodUGforVM.AccessTokenID == accessToken.ID {
					currentATInVMUsersGet = accessToken.Name
				}
			}
			// получаем набор для выбора токена доступа для метода users.get в модуле video_monitor
			kitSUGinVM := MakeSettingComboboxKit("Ключ доступа для \"users.get\"", accessTokensNames, currentATInVMUsersGet)

			// запрашиваем из базы данных данные по методу groups.getById соответствующих субъекта и модуля мониторинга
			var methodGGBIforVM Method
			methodGGBIforVM.selectFromDBByNameAndBySubjectIDAndByMonitorID("groups.getById", subject.ID, monitorVM.ID)
			var currentATInVMGroupsGetByID string
			for _, accessToken := range accessTokens {
				if methodGGBIforVM.AccessTokenID == accessToken.ID {
					currentATInVMGroupsGetByID = accessToken.Name
				}
			}
			// получаем набор для выбора токена доступа для метода groups.getById в модуле video_monitor
			kitSGGBIinVM := MakeSettingComboboxKit("Ключ доступа для \"groups.getById\"", accessTokensNames, currentATInVMGroupsGetByID)

			// запрашиваем из базы данных данные по методу messages.send соответствующих субъекта и модуля мониторинга
			var methodMSforVM Method
			methodMSforVM.selectFromDBByNameAndBySubjectIDAndByMonitorID("messages.send", subject.ID, monitorVM.ID)
			var currentATInVMMessagesSend string
			for _, accessToken := range accessTokens {
				if methodMSforVM.AccessTokenID == accessToken.ID {
					currentATInVMMessagesSend = accessToken.Name
				}
			}
			// получаем набор для выбора токена доступа для метода messages.send в модуле video_monitor
			kitSMSinVM := MakeSettingComboboxKit("Ключ доступа для \"messages.send\"", accessTokensNames, currentATInVMMessagesSend)

			// описываем группу для токенов доступа модуля video_monitor субъекта
			boxWndSVM := ui.NewVerticalBox()
			groupWndSVM := ui.NewGroup("Видео в альбомах")
			groupWndSVM.SetMargined(true)
			boxWndSVM.Append(kitSVGinVM.Box, false)
			boxWndSVM.Append(kitSUGinVM.Box, false)
			boxWndSVM.Append(kitSGGBIinVM.Box, false)
			boxWndSVM.Append(kitSMSinVM.Box, false)
			groupWndSVM.SetChild(boxWndSVM)

			// запрашиваем из базы данных данные по модулю мониторинга photo_comment_monitor соответствующего субъекта
			var monitorPCM Monitor
			monitorPCM.selectFromDBByNameAndBySubjectID("photo_comment_monitor", subject.ID)

			// запрашиваем из базы данных данные по методу photos.getAllComments соответствующих субъекта и модуля мониторинга
			var methodPGACforPCM Method
			methodPGACforPCM.selectFromDBByNameAndBySubjectIDAndByMonitorID("photos.getAllComments", subject.ID, monitorPCM.ID)
			var currentATInPCMPhotosGetAllComments string
			for _, accessToken := range accessTokens {
				if methodPGACforPCM.AccessTokenID == accessToken.ID {
					currentATInPCMPhotosGetAllComments = accessToken.Name
				}
			}
			// получаем набор для выбора токена доступа для метода photos.getAllComments в модуле photo_comment_monitor
			kitSPGACinPCM := MakeSettingComboboxKit("Ключ доступа для \"photos.getAllComments\"", accessTokensNames, currentATInPCMPhotosGetAllComments)

			// запрашиваем из базы данных данные по методу users.get соответствующих субъекта и модуля мониторинга
			var methodUGforPCM Method
			methodUGforPCM.selectFromDBByNameAndBySubjectIDAndByMonitorID("users.get", subject.ID, monitorPCM.ID)
			var currentATInPCMUsersGet string
			for _, accessToken := range accessTokens {
				if methodUGforPCM.AccessTokenID == accessToken.ID {
					currentATInPCMUsersGet = accessToken.Name
				}
			}
			// получаем набор для выбора токена доступа для метода users.get в модуле photo_comment_monitor
			kitSUGinPCM := MakeSettingComboboxKit("Ключ доступа для \"users.get\"", accessTokensNames, currentATInPCMUsersGet)

			// запрашиваем из базы данных данные по методу groups.getById соответствующих субъекта и модуля мониторинга
			var methodGGBIforPCM Method
			methodGGBIforPCM.selectFromDBByNameAndBySubjectIDAndByMonitorID("groups.getById", subject.ID, monitorPCM.ID)
			var currentATInPCMGroupsGetByID string
			for _, accessToken := range accessTokens {
				if methodGGBIforPCM.AccessTokenID == accessToken.ID {
					currentATInPCMGroupsGetByID = accessToken.Name
				}
			}
			// получаем набор для выбора токена доступа для метода groups.getById в модуле photo_comment_monitor
			kitSGGBIinPCM := MakeSettingComboboxKit("Ключ доступа для \"groups.getById\"", accessTokensNames, currentATInPCMGroupsGetByID)

			// запрашиваем из базы данных данные по методу messages.send соответствующих субъекта и модуля мониторинга
			var methodMSforPCM Method
			methodMSforPCM.selectFromDBByNameAndBySubjectIDAndByMonitorID("messages.send", subject.ID, monitorPCM.ID)
			var currentATInPCMMessagesSend string
			for _, accessToken := range accessTokens {
				if methodMSforPCM.AccessTokenID == accessToken.ID {
					currentATInPCMMessagesSend = accessToken.Name
				}
			}
			// получаем набор для выбора токена доступа для метода messages.send в модуле photo_comment_monitor
			kitSMSinPCM := MakeSettingComboboxKit("Ключ доступа для \"messages.send\"", accessTokensNames, currentATInPCMMessagesSend)

			// описываем группу для токенов доступа модуля photo_comment_monitor субъекта
			boxWndSPCM := ui.NewVerticalBox()
			groupWndSPCM := ui.NewGroup("Комментарии под фото")
			groupWndSPCM.SetMargined(true)
			boxWndSPCM.Append(kitSPGACinPCM.Box, false)
			boxWndSPCM.Append(kitSUGinPCM.Box, false)
			boxWndSPCM.Append(kitSGGBIinPCM.Box, false)
			boxWndSPCM.Append(kitSMSinPCM.Box, false)
			groupWndSPCM.SetChild(boxWndSPCM)

			// запрашиваем из базы данных данные по модулю мониторинга video_comment_monitor соответствующего субъекта
			var monitorVCM Monitor
			monitorVCM.selectFromDBByNameAndBySubjectID("video_comment_monitor", subject.ID)

			// запрашиваем из базы данных данные по методу video.getComments соответствующих субъекта и модуля мониторинга
			var methodVGCforVCM Method
			methodVGCforVCM.selectFromDBByNameAndBySubjectIDAndByMonitorID("video.getComments", subject.ID, monitorVCM.ID)
			var currentATInVCMVideoGetComments string
			for _, accessToken := range accessTokens {
				if methodVGCforVCM.AccessTokenID == accessToken.ID {
					currentATInVCMVideoGetComments = accessToken.Name
				}
			}
			// получаем набор для выбора токена доступа для метода video.getComments в модуле video_comment_monitor
			kitSVGCinVCM := MakeSettingComboboxKit("Ключ доступа для \"video.getComments\"", accessTokensNames, currentATInVCMVideoGetComments)

			// запрашиваем из базы данных данные по методу users.get соответствующих субъекта и модуля мониторинга
			var methodUGforVCM Method
			methodUGforVCM.selectFromDBByNameAndBySubjectIDAndByMonitorID("users.get", subject.ID, monitorVCM.ID)
			var currentATInVCMUsersGet string
			for _, accessToken := range accessTokens {
				if methodUGforVCM.AccessTokenID == accessToken.ID {
					currentATInVCMUsersGet = accessToken.Name
				}
			}
			// получаем набор для выбора токена доступа для метода users.get в модуле video_comment_monitor
			kitSUGinVCM := MakeSettingComboboxKit("Ключ доступа для \"users.get\"", accessTokensNames, currentATInVCMUsersGet)

			// запрашиваем из базы данных данные по методу groups.getById соответствующих субъекта и модуля мониторинга
			var methodGGBIforVCM Method
			methodGGBIforVCM.selectFromDBByNameAndBySubjectIDAndByMonitorID("groups.getById", subject.ID, monitorVCM.ID)
			var currentATInVCMGroupsGetByID string
			for _, accessToken := range accessTokens {
				if methodGGBIforVCM.AccessTokenID == accessToken.ID {
					currentATInVCMGroupsGetByID = accessToken.Name
				}
			}
			// получаем набор для выбора токена доступа для метода groups.getById в модуле video_comment_monitor
			kitSGGBIinVCM := MakeSettingComboboxKit("Ключ доступа для \"groups.getById\"", accessTokensNames, currentATInVCMGroupsGetByID)

			// запрашиваем из базы данных данные по методу video.get соответствующих субъекта и модуля мониторинга
			var methodVGforVCM Method
			methodVGforVCM.selectFromDBByNameAndBySubjectIDAndByMonitorID("video.get", subject.ID, monitorVCM.ID)
			var currentATInVCMVideoGet string
			for _, accessToken := range accessTokens {
				if methodVGforVCM.AccessTokenID == accessToken.ID {
					currentATInVCMVideoGet = accessToken.Name
				}
			}
			// получаем набор для выбора токена доступа для метода video.get в модуле video_comment_monitor
			kitSVGinVCM := MakeSettingComboboxKit("Ключ доступа для \"video.get\"", accessTokensNames, currentATInVCMVideoGet)

			// запрашиваем из базы данных данные по методу messages.send соответствующих субъекта и модуля мониторинга
			var methodMSforVCM Method
			methodMSforVCM.selectFromDBByNameAndBySubjectIDAndByMonitorID("messages.send", subject.ID, monitorVCM.ID)
			var currentATInVCMMessagesSend string
			for _, accessToken := range accessTokens {
				if methodMSforVCM.AccessTokenID == accessToken.ID {
					currentATInVCMMessagesSend = accessToken.Name
				}
			}
			// получаем набор для выбора токена доступа для метода messages.send в модуле video_comment_monitor
			kitSMSinVCM := MakeSettingComboboxKit("Ключ доступа для \"messages.send\"", accessTokensNames, currentATInVCMMessagesSend)

			// описываем группу для токенов доступа модуля video_comment_monitor субъекта
			boxWndSVCM := ui.NewVerticalBox()
			groupWndSVCM := ui.NewGroup("Комментарии под видео")
			groupWndSVCM.SetMargined(true)
			boxWndSVCM.Append(kitSVGCinVCM.Box, false)
			boxWndSVCM.Append(kitSUGinVCM.Box, false)
			boxWndSVCM.Append(kitSGGBIinVCM.Box, false)
			boxWndSVCM.Append(kitSVGinVCM.Box, false)
			boxWndSVCM.Append(kitSMSinVCM.Box, false)
			groupWndSVCM.SetChild(boxWndSVCM)

			// запрашиваем из базы данных данные по модулю мониторинга topic_monitor соответствующего субъекта
			var monitorTM Monitor
			monitorTM.selectFromDBByNameAndBySubjectID("topic_monitor", subject.ID)

			// запрашиваем из базы данных данные по методу board.getComments соответствующих субъекта и модуля мониторинга
			var methodBGCforTM Method
			methodBGCforTM.selectFromDBByNameAndBySubjectIDAndByMonitorID("board.getComments", subject.ID, monitorTM.ID)
			var currentATInTMBoardGetComments string
			for _, accessToken := range accessTokens {
				if methodBGCforTM.AccessTokenID == accessToken.ID {
					currentATInTMBoardGetComments = accessToken.Name
				}
			}
			// получаем набор для выбора токена доступа для метода board.getComments в модуле topic_monitor
			kitSBGCinTM := MakeSettingComboboxKit("Ключ доступа для \"board.getComments\"", accessTokensNames, currentATInTMBoardGetComments)

			// запрашиваем из базы данных данные по методу board.getTopics соответствующих субъекта и модуля мониторинга
			var methodBGTforTM Method
			methodBGTforTM.selectFromDBByNameAndBySubjectIDAndByMonitorID("board.getTopics", subject.ID, monitorTM.ID)
			var currentATInTMBoardGetTopics string
			for _, accessToken := range accessTokens {
				if methodBGTforTM.AccessTokenID == accessToken.ID {
					currentATInTMBoardGetTopics = accessToken.Name
				}
			}
			// получаем набор для выбора токена доступа для метода board.getTopics в модуле topic_monitor
			kitSBGTinTM := MakeSettingComboboxKit("Ключ доступа для \"board.getTopics\"", accessTokensNames, currentATInTMBoardGetTopics)

			// запрашиваем из базы данных данные по методу users.get соответствующих субъекта и модуля мониторинга
			var methodUGforTM Method
			methodUGforTM.selectFromDBByNameAndBySubjectIDAndByMonitorID("users.get", subject.ID, monitorTM.ID)
			var currentATInTMUsersGet string
			for _, accessToken := range accessTokens {
				if methodUGforTM.AccessTokenID == accessToken.ID {
					currentATInTMUsersGet = accessToken.Name
				}
			}
			// получаем набор для выбора токена доступа для метода users.get в модуле topic_monitor
			kitSUGinTM := MakeSettingComboboxKit("Ключ доступа для \"users.get\"", accessTokensNames, currentATInTMUsersGet)

			// запрашиваем из базы данных данные по методу groups.getById соответствующих субъекта и модуля мониторинга
			var methodGGBIforTM Method
			methodGGBIforTM.selectFromDBByNameAndBySubjectIDAndByMonitorID("groups.getById", subject.ID, monitorTM.ID)
			var currentATInTMGroupsGetByID string
			for _, accessToken := range accessTokens {
				if methodGGBIforTM.AccessTokenID == accessToken.ID {
					currentATInTMGroupsGetByID = accessToken.Name
				}
			}
			// получаем набор для выбора токена доступа для метода groups.getById в модуле topic_monitor
			kitSGGBIinTM := MakeSettingComboboxKit("Ключ доступа для \"groups.getById\"", accessTokensNames, currentATInTMGroupsGetByID)

			// запрашиваем из базы данных данные по методу messages.send соответствующих субъекта и модуля мониторинга
			var methodMSforTM Method
			methodMSforTM.selectFromDBByNameAndBySubjectIDAndByMonitorID("messages.send", subject.ID, monitorTM.ID)
			var currentATInTMMessagesSend string
			for _, accessToken := range accessTokens {
				if methodMSforTM.AccessTokenID == accessToken.ID {
					currentATInTMMessagesSend = accessToken.Name
				}
			}
			// получаем набор для выбора токена доступа для метода messages.send в модуле topic_monitor
			kitSMSinTM := MakeSettingComboboxKit("Ключ доступа для \"messages.send\"", accessTokensNames, currentATInTMMessagesSend)

			// описываем группу для токенов доступа модуля topic_monitor субъекта
			boxWndSTM := ui.NewVerticalBox()
			groupWndSTM := ui.NewGroup("Комментарии в обсуждениях")
			groupWndSTM.SetMargined(true)
			boxWndSTM.Append(kitSBGCinTM.Box, false)
			boxWndSTM.Append(kitSBGTinTM.Box, false)
			boxWndSTM.Append(kitSUGinTM.Box, false)
			boxWndSTM.Append(kitSGGBIinTM.Box, false)
			boxWndSTM.Append(kitSMSinTM.Box, false)
			groupWndSTM.SetChild(boxWndSTM)

			// запрашиваем из базы данных данные по модулю мониторинга wall_post_comment_monitor соответствующего субъекта
			var monitorWPCM Monitor
			monitorWPCM.selectFromDBByNameAndBySubjectID("wall_post_comment_monitor", subject.ID)

			// запрашиваем из базы данных данные по методу wall.getComments соответствующих субъекта и модуля мониторинга
			var methodWGCsforWPCM Method
			methodWGCsforWPCM.selectFromDBByNameAndBySubjectIDAndByMonitorID("wall.getComments", subject.ID, monitorWPCM.ID)
			var currentATInWPCMWallGetComments string
			for _, accessToken := range accessTokens {
				if methodWGCsforWPCM.AccessTokenID == accessToken.ID {
					currentATInWPCMWallGetComments = accessToken.Name
				}
			}
			// получаем набор для выбора токена доступа для метода wall.getComments в модуле wall_post_comment_monitor
			kitSWGCsinWPCM := MakeSettingComboboxKit("Ключ доступа для \"wall.getComments\"", accessTokensNames, currentATInWPCMWallGetComments)

			// запрашиваем из базы данных данные по методу users.get соответствующих субъекта и модуля мониторинга
			var methodUGforWPCM Method
			methodUGforWPCM.selectFromDBByNameAndBySubjectIDAndByMonitorID("users.get", subject.ID, monitorWPCM.ID)
			var currentATInWPCMUsersGet string
			for _, accessToken := range accessTokens {
				if methodUGforWPCM.AccessTokenID == accessToken.ID {
					currentATInWPCMUsersGet = accessToken.Name
				}
			}
			// получаем набор для выбора токена доступа для метода users.get в модуле wall_post_comment_monitor
			kitSUGinWPCM := MakeSettingComboboxKit("Ключ доступа для \"users.get\"", accessTokensNames, currentATInWPCMUsersGet)

			// запрашиваем из базы данных данные по методу groups.getById соответствующих субъекта и модуля мониторинга
			var methodGGBIforWPCM Method
			methodGGBIforWPCM.selectFromDBByNameAndBySubjectIDAndByMonitorID("groups.getById", subject.ID, monitorWPCM.ID)
			var currentATInWPCMGroupsGetByID string
			for _, accessToken := range accessTokens {
				if methodGGBIforWPCM.AccessTokenID == accessToken.ID {
					currentATInWPCMGroupsGetByID = accessToken.Name
				}
			}
			// получаем набор для выбора токена доступа для метода groups.getById в модуле wall_post_comment_monitor
			kitSGGBIinWPCM := MakeSettingComboboxKit("Ключ доступа для \"groups.getById\"", accessTokensNames, currentATInWPCMGroupsGetByID)

			// запрашиваем из базы данных данные по методу wall.get соответствующих субъекта и модуля мониторинга
			var methodWGforWPCM Method
			methodWGforWPCM.selectFromDBByNameAndBySubjectIDAndByMonitorID("wall.get", subject.ID, monitorWPCM.ID)
			var currentATInWPCMWallGet string
			for _, accessToken := range accessTokens {
				if methodWGforWPCM.AccessTokenID == accessToken.ID {
					currentATInWPCMWallGet = accessToken.Name
				}
			}
			// получаем набор для выбора токена доступа для метода wall.get в модуле wall_post_comment_monitor
			kitSWGinWPCM := MakeSettingComboboxKit("Ключ доступа для \"wall.get\"", accessTokensNames, currentATInWPCMWallGet)

			// запрашиваем из базы данных данные по методу wall.getComment соответствующих субъекта и модуля мониторинга
			var methodWGCforWPCM Method
			methodWGCforWPCM.selectFromDBByNameAndBySubjectIDAndByMonitorID("wall.getComment", subject.ID, monitorWPCM.ID)
			var currentATInWPCMWallGetComment string
			for _, accessToken := range accessTokens {
				if methodWGCforWPCM.AccessTokenID == accessToken.ID {
					currentATInWPCMWallGetComment = accessToken.Name
				}
			}
			// получаем набор для выбора токена доступа для метода wall.getComment в модуле wall_post_comment_monitor
			kitSWGCinWPCM := MakeSettingComboboxKit("Ключ доступа для \"wall.getComment\"", accessTokensNames, currentATInWPCMWallGetComment)

			// запрашиваем из базы данных данные по методу messages.send соответствующих субъекта и модуля мониторинга
			var methodMSforWPCM Method
			methodMSforWPCM.selectFromDBByNameAndBySubjectIDAndByMonitorID("messages.send", subject.ID, monitorWPCM.ID)
			var currentATInWPCMMessagesSend string
			for _, accessToken := range accessTokens {
				if methodMSforWPCM.AccessTokenID == accessToken.ID {
					currentATInWPCMMessagesSend = accessToken.Name
				}
			}
			// получаем набор для выбора токена доступа для метода messages.send в модуле wall_post_comment_monitor
			kitSMSinWPCM := MakeSettingComboboxKit("Ключ доступа для \"messages.send\"", accessTokensNames, currentATInWPCMMessagesSend)

			// описываем группу для токенов доступа модуля wall_post_comment_monitor субъекта
			boxWndSWGCM := ui.NewVerticalBox()
			groupWndWGCM := ui.NewGroup("Комментарии под постами")
			groupWndWGCM.SetMargined(true)
			boxWndSWGCM.Append(kitSWGCsinWPCM.Box, false)
			boxWndSWGCM.Append(kitSUGinWPCM.Box, false)
			boxWndSWGCM.Append(kitSGGBIinWPCM.Box, false)
			boxWndSWGCM.Append(kitSWGinWPCM.Box, false)
			boxWndSWGCM.Append(kitSWGCinWPCM.Box, false)
			boxWndSWGCM.Append(kitSMSinWPCM.Box, false)
			groupWndWGCM.SetChild(boxWndSWGCM)

			// добавляем все заполненные группы на
			// левую коробку
			boxWndSLeft.Append(groupWndSGeneral, false)
			boxWndSLeft.Append(groupWndSWPM, false)
			boxWndSLeft.Append(groupWndSAPM, false)

			// коробку посередине
			boxWndSCenter.Append(groupWndSVM, false)
			boxWndSCenter.Append(groupWndSPCM, false)
			boxWndSCenter.Append(groupWndSVCM, false)

			// и правую коробку (для экономии места из-за отсутствия в данной библиотеке для gui скроллинга)
			boxWndSRight.Append(groupWndSTM, false)
			boxWndSRight.Append(groupWndWGCM, false)

			// затем добавляем левую, центральную и правую коробки на одну общую
			boxWndSMain.Append(boxWndSLeft, false)
			boxWndSMain.Append(boxWndSCenter, false)
			boxWndSMain.Append(boxWndSRight, false)

			// а ее - на основную
			kitWindowGeneralSettings.Box.Append(boxWndSMain, false)

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
					warningTitle := "Поле \"Идентификатор в ВК\" не должно быть пустым"
					ShowWarningWindow(warningTitle)
					return
				}
				updatedSubject.SubjectID, err = strconv.Atoi(kitWndSSubjectID.Entry.Text())
				if err != nil {
					ToLogFile(err.Error(), string(debug.Stack()))
					panic(err.Error())
				}
				if len(kitWndSName.Entry.Text()) == 0 {
					warningTitle := "Поле \"Название\" не должно быть пустым"
					ShowWarningWindow(warningTitle)
					return
				}
				updatedSubject.Name = kitWndSName.Entry.Text()
				updatedSubject.BackupWikipage = subject.BackupWikipage
				updatedSubject.LastBackup = subject.LastBackup
				err := updatedSubject.updateInDB()
				if err != nil {
					ToLogFile(err.Error(), string(debug.Stack()))
					panic(err.Error())
				}

				var updatedMethodWGforWPM = methodWGforWPM
				updatedMethodWGforWPM.AccessTokenID = accessTokens[kitSWGinWPM.Combobox.Selected()].ID
				err = updatedMethodWGforWPM.updateInDB()
				if err != nil {
					ToLogFile(err.Error(), string(debug.Stack()))
					panic(err.Error())
				}

				var updatedMethodUGforWPM = methodUGforWPM
				updatedMethodUGforWPM.AccessTokenID = accessTokens[kitSUGinWPM.Combobox.Selected()].ID
				err = updatedMethodUGforWPM.updateInDB()
				if err != nil {
					ToLogFile(err.Error(), string(debug.Stack()))
					panic(err.Error())
				}

				var updatedMethodGGBIforWPM = methodGGBIforWPM
				updatedMethodGGBIforWPM.AccessTokenID = accessTokens[kitSGGBIinWPM.Combobox.Selected()].ID
				err = updatedMethodGGBIforWPM.updateInDB()
				if err != nil {
					ToLogFile(err.Error(), string(debug.Stack()))
					panic(err.Error())
				}

				var updatedMethodMSforWPM = methodMSforWPM
				updatedMethodMSforWPM.AccessTokenID = accessTokens[kitSMSinWPM.Combobox.Selected()].ID
				err = updatedMethodMSforWPM.updateInDB()
				if err != nil {
					ToLogFile(err.Error(), string(debug.Stack()))
					panic(err.Error())
				}

				var updatedMethodPGforAPM = methodPGforAPM
				updatedMethodPGforAPM.AccessTokenID = accessTokens[kitSPGinAPM.Combobox.Selected()].ID
				err = updatedMethodPGforAPM.updateInDB()
				if err != nil {
					ToLogFile(err.Error(), string(debug.Stack()))
					panic(err.Error())
				}

				var updatedMethodPGAforAPM = methodPGAforAPM
				updatedMethodPGAforAPM.AccessTokenID = accessTokens[kitSPGAinAPM.Combobox.Selected()].ID
				err = updatedMethodPGAforAPM.updateInDB()
				if err != nil {
					ToLogFile(err.Error(), string(debug.Stack()))
					panic(err.Error())
				}

				var updatedMethodUGforAPM = methodUGforAPM
				updatedMethodUGforAPM.AccessTokenID = accessTokens[kitSUGinAPM.Combobox.Selected()].ID
				err = updatedMethodUGforAPM.updateInDB()
				if err != nil {
					ToLogFile(err.Error(), string(debug.Stack()))
					panic(err.Error())
				}

				var updatedMethodGGBIforAPM = methodGGBIforAPM
				updatedMethodGGBIforAPM.AccessTokenID = accessTokens[kitSGGBIinAPM.Combobox.Selected()].ID
				err = updatedMethodGGBIforAPM.updateInDB()
				if err != nil {
					ToLogFile(err.Error(), string(debug.Stack()))
					panic(err.Error())
				}

				var updatedMethodMSforAPM = methodMSforAPM
				updatedMethodMSforAPM.AccessTokenID = accessTokens[kitSMSinAPM.Combobox.Selected()].ID
				err = updatedMethodMSforAPM.updateInDB()
				if err != nil {
					ToLogFile(err.Error(), string(debug.Stack()))
					panic(err.Error())
				}

				var updatedMethodVGforVM = methodVGforVM
				updatedMethodVGforVM.AccessTokenID = accessTokens[kitSVGinVM.Combobox.Selected()].ID
				err = updatedMethodVGforVM.updateInDB()
				if err != nil {
					ToLogFile(err.Error(), string(debug.Stack()))
					panic(err.Error())
				}

				var updatedMethodUGforVM = methodUGforVM
				updatedMethodUGforVM.AccessTokenID = accessTokens[kitSUGinVM.Combobox.Selected()].ID
				err = updatedMethodUGforVM.updateInDB()
				if err != nil {
					ToLogFile(err.Error(), string(debug.Stack()))
					panic(err.Error())
				}

				var updatedMethodGGBIforVM = methodGGBIforVM
				updatedMethodGGBIforVM.AccessTokenID = accessTokens[kitSGGBIinVM.Combobox.Selected()].ID
				err = updatedMethodGGBIforVM.updateInDB()
				if err != nil {
					ToLogFile(err.Error(), string(debug.Stack()))
					panic(err.Error())
				}

				var updatedMethodMSforVM = methodMSforVM
				updatedMethodMSforVM.AccessTokenID = accessTokens[kitSMSinVM.Combobox.Selected()].ID
				err = updatedMethodMSforVM.updateInDB()
				if err != nil {
					ToLogFile(err.Error(), string(debug.Stack()))
					panic(err.Error())
				}

				var updatedMethodPGACforPCM = methodPGACforPCM
				updatedMethodPGACforPCM.AccessTokenID = accessTokens[kitSPGACinPCM.Combobox.Selected()].ID
				err = updatedMethodPGACforPCM.updateInDB()
				if err != nil {
					ToLogFile(err.Error(), string(debug.Stack()))
					panic(err.Error())
				}

				var updatedMethodUGforPCM = methodUGforPCM
				updatedMethodUGforPCM.AccessTokenID = accessTokens[kitSUGinPCM.Combobox.Selected()].ID
				err = updatedMethodUGforPCM.updateInDB()
				if err != nil {
					ToLogFile(err.Error(), string(debug.Stack()))
					panic(err.Error())
				}

				var updatedMethodGGBIforPCM = methodGGBIforPCM
				updatedMethodGGBIforPCM.AccessTokenID = accessTokens[kitSGGBIinPCM.Combobox.Selected()].ID
				err = updatedMethodGGBIforPCM.updateInDB()
				if err != nil {
					ToLogFile(err.Error(), string(debug.Stack()))
					panic(err.Error())
				}

				var updatedMethodMSforPCM = methodMSforPCM
				updatedMethodMSforPCM.AccessTokenID = accessTokens[kitSMSinPCM.Combobox.Selected()].ID
				err = updatedMethodMSforPCM.updateInDB()
				if err != nil {
					ToLogFile(err.Error(), string(debug.Stack()))
					panic(err.Error())
				}

				var updatedMethodVGCforVCM = methodVGCforVCM
				updatedMethodVGCforVCM.AccessTokenID = accessTokens[kitSVGCinVCM.Combobox.Selected()].ID
				err = updatedMethodVGCforVCM.updateInDB()
				if err != nil {
					ToLogFile(err.Error(), string(debug.Stack()))
					panic(err.Error())
				}

				var updatedMethodUGforVCM = methodUGforVCM
				updatedMethodUGforVCM.AccessTokenID = accessTokens[kitSUGinVCM.Combobox.Selected()].ID
				err = updatedMethodUGforVCM.updateInDB()
				if err != nil {
					ToLogFile(err.Error(), string(debug.Stack()))
					panic(err.Error())
				}

				var updatedMethodGGBIforVCM = methodGGBIforVCM
				updatedMethodGGBIforVCM.AccessTokenID = accessTokens[kitSGGBIinVCM.Combobox.Selected()].ID
				err = updatedMethodGGBIforVCM.updateInDB()
				if err != nil {
					ToLogFile(err.Error(), string(debug.Stack()))
					panic(err.Error())
				}

				var updatedMethodVGforVCM = methodVGforVCM
				updatedMethodVGforVCM.AccessTokenID = accessTokens[kitSVGinVCM.Combobox.Selected()].ID
				err = updatedMethodVGforVCM.updateInDB()
				if err != nil {
					ToLogFile(err.Error(), string(debug.Stack()))
					panic(err.Error())
				}

				var updatedMethodMSforVCM = methodMSforVCM
				updatedMethodMSforVCM.AccessTokenID = accessTokens[kitSMSinVCM.Combobox.Selected()].ID
				err = updatedMethodMSforVCM.updateInDB()
				if err != nil {
					ToLogFile(err.Error(), string(debug.Stack()))
					panic(err.Error())
				}

				var updatedMethodBGCforTM = methodBGCforTM
				updatedMethodBGCforTM.AccessTokenID = accessTokens[kitSBGCinTM.Combobox.Selected()].ID
				err = updatedMethodBGCforTM.updateInDB()
				if err != nil {
					ToLogFile(err.Error(), string(debug.Stack()))
					panic(err.Error())
				}

				var updatedMethodBGTforTM = methodBGTforTM
				updatedMethodBGTforTM.AccessTokenID = accessTokens[kitSBGTinTM.Combobox.Selected()].ID
				err = updatedMethodBGTforTM.updateInDB()
				if err != nil {
					ToLogFile(err.Error(), string(debug.Stack()))
					panic(err.Error())
				}

				var updatedMethodUGforTM = methodUGforTM
				updatedMethodUGforTM.AccessTokenID = accessTokens[kitSUGinTM.Combobox.Selected()].ID
				err = updatedMethodUGforTM.updateInDB()
				if err != nil {
					ToLogFile(err.Error(), string(debug.Stack()))
					panic(err.Error())
				}

				var updatedMethodGGBIforTM = methodGGBIforTM
				updatedMethodGGBIforTM.AccessTokenID = accessTokens[kitSGGBIinTM.Combobox.Selected()].ID
				err = updatedMethodGGBIforTM.updateInDB()
				if err != nil {
					ToLogFile(err.Error(), string(debug.Stack()))
					panic(err.Error())
				}

				var updatedMethodMSforTM = methodMSforTM
				updatedMethodMSforTM.AccessTokenID = accessTokens[kitSMSinTM.Combobox.Selected()].ID
				err = updatedMethodMSforTM.updateInDB()
				if err != nil {
					ToLogFile(err.Error(), string(debug.Stack()))
					panic(err.Error())
				}

				var updatedMethodWGCsforWPCM = methodWGCsforWPCM
				updatedMethodWGCsforWPCM.AccessTokenID = accessTokens[kitSWGCsinWPCM.Combobox.Selected()].ID
				err = updatedMethodWGCsforWPCM.updateInDB()
				if err != nil {
					ToLogFile(err.Error(), string(debug.Stack()))
					panic(err.Error())
				}

				var updatedMethodUGforWPCM = methodUGforWPCM
				updatedMethodUGforWPCM.AccessTokenID = accessTokens[kitSUGinWPCM.Combobox.Selected()].ID
				err = updatedMethodUGforWPCM.updateInDB()
				if err != nil {
					ToLogFile(err.Error(), string(debug.Stack()))
					panic(err.Error())
				}

				var updatedMethodGGBIforWPCM = methodGGBIforWPCM
				updatedMethodGGBIforWPCM.AccessTokenID = accessTokens[kitSGGBIinWPCM.Combobox.Selected()].ID
				err = updatedMethodGGBIforWPCM.updateInDB()
				if err != nil {
					ToLogFile(err.Error(), string(debug.Stack()))
					panic(err.Error())
				}

				var updatedMethodWGforWPCM = methodWGforWPCM
				updatedMethodWGforWPCM.AccessTokenID = accessTokens[kitSWGinWPCM.Combobox.Selected()].ID
				err = updatedMethodWGforWPCM.updateInDB()
				if err != nil {
					ToLogFile(err.Error(), string(debug.Stack()))
					panic(err.Error())
				}

				var updatedMethodWGCforWPCM = methodWGCforWPCM
				updatedMethodWGCforWPCM.AccessTokenID = accessTokens[kitSWGCinWPCM.Combobox.Selected()].ID
				err = updatedMethodWGCforWPCM.updateInDB()
				if err != nil {
					ToLogFile(err.Error(), string(debug.Stack()))
					panic(err.Error())
				}

				var updatedMethodMSforWPCM = methodMSforWPCM
				updatedMethodMSforWPCM.AccessTokenID = accessTokens[kitSMSinWPCM.Combobox.Selected()].ID
				err = updatedMethodMSforWPCM.updateInDB()
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
	var wallPostMonitorParam WallPostMonitorParam
	err := wallPostMonitorParam.selectFromDBBySubjectID(IDSubject)
	if err != nil {
		ToLogFile(err.Error(), string(debug.Stack()))
		panic(err.Error())
	}

	// получаем набор для отображения установок модуля мониторинга постов на стене
	windowTitle := fmt.Sprintf("%v: настройки для %v", btnName, nameSubject)
	kitWindowWallPostSettings := MakeSettingWindowKit(windowTitle, 300, 100)

	boxWndWP := ui.NewVerticalBox()

	// получаем набор для флага необходимости активировать модуль мониторинга
	kitWndWPMonitoring := MakeSettingCheckboxKit("Разрешить наблюдение", wallPostMonitorParam.NeedMonitoring)

	// получаем набор для интервала между запусками функции мониторинга
	kitWndWPInterval := MakeSettingSpinboxKit("Интервал (сек.)", 5, 21600, wallPostMonitorParam.Interval)

	// получаем набор для количества проверяемых постов
	kitWndWPSendTo := MakeSettingEntryKit("Получатель", strconv.Itoa(wallPostMonitorParam.SendTo))
	// и привязываем его к обработке ввода, чтобы кроме чисел ничего не вводилось
	kitWndWPSendTo.Entry.OnChanged(func(*ui.Entry) {
		NumericEntriesHandler(kitWndWPSendTo.Entry)
	})

	// получаем набор для фильтра получаемых для проверки постов
	listPostsFilters := []string{"Все", "От пользователей", "От сообщества", "Предложенные"}
	listPostsFiltersEn := []string{"all", "others", "owner", "suggests"}
	var currentFilter string
	for i, item := range listPostsFiltersEn {
		if wallPostMonitorParam.Filter == item {
			currentFilter = listPostsFilters[i]
			break
		}
	}
	kitWndWPFilter := MakeSettingComboboxKit("Фильтр", listPostsFilters, currentFilter)

	// получаем набор для количества проверяемых постов
	kitWndWPPostsCount := MakeSettingSpinboxKit("Количество постов", 1, 100, wallPostMonitorParam.PostsCount)

	// получаем набор для списка ключевых слов для отбора постов
	kitWndWPKeywordsForMonitoring := MakeSettingEntryListKit("Ключевые слова", wallPostMonitorParam.KeywordsForMonitoring)

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
			warningTitle := "Поле \"Получатель\" не должно быть пустым."
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
			warningTitle := "Нужно выбрать элемент в списке \"Фильтр\""
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

		err = updatedWallPostMonitorParam.updateInDB()
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
	var albumPhotoMonitorParam AlbumPhotoMonitorParam
	err := albumPhotoMonitorParam.selectFromDBBySubjectID(IDSubject)
	if err != nil {
		ToLogFile(err.Error(), string(debug.Stack()))
		panic(err.Error())
	}

	// получаем набор для отображения установок модуля мониторинга фотографий в альбомах
	windowTitle := fmt.Sprintf("%v: настройки для %v", btnName, nameSubject)
	kitWindowAlbumPhotoSettings := MakeSettingWindowKit(windowTitle, 300, 100)

	boxWndAP := ui.NewVerticalBox()

	// получаем набор для флага необходимости активировать модуль мониторинга
	kitWndAPMonitoring := MakeSettingCheckboxKit("Разрешить наблюдение", albumPhotoMonitorParam.NeedMonitoring)

	// получаем набор для интервала между запусками функции мониторинга
	kitWndAPInterval := MakeSettingSpinboxKit("Интервал (сек.)", 5, 21600, albumPhotoMonitorParam.Interval)

	// получаем набор для идентификатора получателя сообщений
	kitWndAPSendTo := MakeSettingEntryKit("Получатель", strconv.Itoa(albumPhotoMonitorParam.SendTo))
	// и привязываем его к обработке ввода, чтобы кроме чисел ничего не вводилось
	kitWndAPSendTo.Entry.OnChanged(func(*ui.Entry) {
		NumericEntriesHandler(kitWndAPSendTo.Entry)
	})

	// получаем набор для количества проверяемых фото
	kitWndApPhotosCount := MakeSettingSpinboxKit("Количество фотографий", 1, 1000, albumPhotoMonitorParam.PhotosCount)

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
			warningTitle := "Поле \"Получатель\" не должно быть пустым."
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

		err = updatedAlbumPhotoMonitorParam.updateInDB()
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
	var videoMonitorParam VideoMonitorParam
	err := videoMonitorParam.selectFromDBBySubjectID(IDSubject)
	if err != nil {
		ToLogFile(err.Error(), string(debug.Stack()))
		panic(err.Error())
	}

	// получаем набор для отображения установок модуля мониторинга видео
	windowTitle := fmt.Sprintf("%v: настройки для %v", btnName, nameSubject)
	kitWindowVideoSettings := MakeSettingWindowKit(windowTitle, 300, 100)

	boxWndV := ui.NewVerticalBox()

	// получаем набор для флага необходимости активировать модуль мониторинга
	kitWndVMonitoring := MakeSettingCheckboxKit("Разрешить наблюдение", videoMonitorParam.NeedMonitoring)

	// получаем набор для интервала между запусками функции мониторинга
	kitWndVInterval := MakeSettingSpinboxKit("Интервал (сек.)", 5, 21600, videoMonitorParam.Interval)

	// получаем набор для идентификатора получателя сообщений
	kitWndVSendTo := MakeSettingEntryKit("Получатель", strconv.Itoa(videoMonitorParam.SendTo))
	// и привязываем его к обработке ввода, чтобы кроме чисел ничего не вводилось
	kitWndVSendTo.Entry.OnChanged(func(*ui.Entry) {
		NumericEntriesHandler(kitWndVSendTo.Entry)
	})

	// получаем набор для количества проверяемых видео
	kitWndVVideoCount := MakeSettingSpinboxKit("Количество видеозаписей", 1, 1000, videoMonitorParam.VideoCount)

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
			warningTitle := "Поле \"Получатель\" не должно быть пустым."
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

		err = updatedVideoMonitorParam.updateInDB()
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
	var photoCommentMonitorParam PhotoCommentMonitorParam
	err := photoCommentMonitorParam.selectFromDBBySubjectID(IDSubject)
	if err != nil {
		ToLogFile(err.Error(), string(debug.Stack()))
		panic(err.Error())
	}

	// получаем набор для отображения установок модуля мониторинга комментариев под фотками
	windowTitle := fmt.Sprintf("%v: настройки для %v", btnName, nameSubject)
	kitWindowPhotoCommentSettings := MakeSettingWindowKit(windowTitle, 300, 100)

	boxWndPC := ui.NewVerticalBox()

	// получаем набор для флага необходимости активировать модуль мониторинга
	kitWndPCMonitoring := MakeSettingCheckboxKit("Разрешить наблюдение", photoCommentMonitorParam.NeedMonitoring)

	// получаем набор для интервала между запусками функции мониторинга
	kitWndPCInterval := MakeSettingSpinboxKit("Интервал (сек.)", 5, 21600, photoCommentMonitorParam.Interval)

	// получаем набор для идентификатора получателя сообщений
	kitWndPCSendTo := MakeSettingEntryKit("Получатель", strconv.Itoa(photoCommentMonitorParam.SendTo))
	// и привязываем его к обработке ввода, чтобы кроме чисел ничего не вводилось
	kitWndPCSendTo.Entry.OnChanged(func(*ui.Entry) {
		NumericEntriesHandler(kitWndPCSendTo.Entry)
	})

	// получаем набор для количества проверяемых комментариев
	kitWndPCCommentsCount := MakeSettingSpinboxKit("Количество комментариев", 1, 1000, photoCommentMonitorParam.CommentsCount)

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
			warningTitle := "Поле \"Получатель\" не должно быть пустым."
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

		err = updatedPhotoCommentMonitorParam.updateInDB()
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
	var videoCommentMonitorParam VideoCommentMonitorParam
	err := videoCommentMonitorParam.selectFromDBBySubjectID(IDSubject)
	if err != nil {
		ToLogFile(err.Error(), string(debug.Stack()))
		panic(err.Error())
	}

	// получаем набор для отображения установок модуля мониторинга комментариев в обсуждениях
	windowTitle := fmt.Sprintf("%v: настройки для %v", btnName, nameSubject)
	kitWindowVideoCommentSettings := MakeSettingWindowKit(windowTitle, 300, 100)

	boxWndVC := ui.NewVerticalBox()

	// получаем набор для флага необходимости активировать модуль мониторинга
	kitWndVCMonitoring := MakeSettingCheckboxKit("Разрешить наблюдение", videoCommentMonitorParam.NeedMonitoring)

	// получаем набор для интервала между запусками функции мониторинга
	kitWndVCInterval := MakeSettingSpinboxKit("Интервал (сек.)", 5, 21600, videoCommentMonitorParam.Interval)

	// получаем набор для идентификатора получателя сообщений
	kitWndVCSendTo := MakeSettingEntryKit("Получатель", strconv.Itoa(videoCommentMonitorParam.SendTo))
	// и привязываем его к обработке ввода, чтобы кроме чисел ничего не вводилось
	kitWndVCSendTo.Entry.OnChanged(func(*ui.Entry) {
		NumericEntriesHandler(kitWndVCSendTo.Entry)
	})

	// получаем набор для количества проверяемых видео
	kitWndVCVideosCount := MakeSettingSpinboxKit("Количество видеозаписей", 1, 200, videoCommentMonitorParam.VideosCount)

	// получаем набор для количества проверяемых комментариев
	kitWndVCCommentsCount := MakeSettingSpinboxKit("Количество комментариев", 1, 100, videoCommentMonitorParam.CommentsCount)

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
			warningTitle := "Поле \"Получатель\" не должно быть пустым."
			ShowWarningWindow(warningTitle)
			return
		}
		updatedVideoCommentMonitorParam.SendTo, err = strconv.Atoi(kitWndVCSendTo.Entry.Text())
		if err != nil {
			// FIXME: тут старый способ обработки ошибок
			date := UnixTimeStampToDate(int(time.Now().Unix()))
			log.Fatal(fmt.Errorf("> [%v] WARNING! Error: %v", date, err))
		}
		updatedVideoCommentMonitorParam.Interval = kitWndVCInterval.Spinbox.Value()
		updatedVideoCommentMonitorParam.LastDate = videoCommentMonitorParam.LastDate
		updatedVideoCommentMonitorParam.CommentsCount = kitWndVCCommentsCount.Spinbox.Value()
		updatedVideoCommentMonitorParam.VideosCount = kitWndVCVideosCount.Spinbox.Value()

		err = updatedVideoCommentMonitorParam.updateInDB()
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
	var topicMonitorParam TopicMonitorParam
	err := topicMonitorParam.selectFromDBBySubjectID(IDSubject)
	if err != nil {
		ToLogFile(err.Error(), string(debug.Stack()))
		panic(err.Error())
	}

	// получаем набор для отображения установок модуля мониторинга комментариев в обсуждениях
	windowTitle := fmt.Sprintf("%v: настройки для %v", btnName, nameSubject)
	kitWindowTopicSettings := MakeSettingWindowKit(windowTitle, 300, 100)

	boxWndT := ui.NewVerticalBox()

	// получаем набор для флага необходимости активировать модуль мониторинга
	kitWndTMonitoring := MakeSettingCheckboxKit("Разрешить наблюдение", topicMonitorParam.NeedMonitoring)

	// получаем набор для интервала между запусками функции мониторинга
	kitWndTInterval := MakeSettingSpinboxKit("Интервал (сек.)", 5, 21600, topicMonitorParam.Interval)

	// получаем набор для идентификатора получателя сообщений
	kitWndTSendTo := MakeSettingEntryKit("Получатель", strconv.Itoa(topicMonitorParam.SendTo))
	// и привязываем его к обработке ввода, чтобы кроме чисел ничего не вводилось
	kitWndTSendTo.Entry.OnChanged(func(*ui.Entry) {
		NumericEntriesHandler(kitWndTSendTo.Entry)
	})

	// получаем набор для количества проверяемых топиков обсуждений
	kitWndTTopicsCount := MakeSettingSpinboxKit("Количество обсуждений", 1, 100, topicMonitorParam.TopicsCount)

	// получаем набор для количества проверяемых комментариев
	kitWndTCommentsCount := MakeSettingSpinboxKit("Количество комментариев", 1, 100, topicMonitorParam.TopicsCount)

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
			warningTitle := "Поле \"Получатель\" не должно быть пустым."
			ShowWarningWindow(warningTitle)
			return
		}
		updatedTopicMonitorParam.SendTo, err = strconv.Atoi(kitWndTSendTo.Entry.Text())
		if err != nil {
			// FIXME: тут старый способ обработки ошибок
			date := UnixTimeStampToDate(int(time.Now().Unix()))
			log.Fatal(fmt.Errorf("> [%v] WARNING! Error: %v", date, err))
		}
		updatedTopicMonitorParam.Interval = kitWndTInterval.Spinbox.Value()
		updatedTopicMonitorParam.LastDate = topicMonitorParam.LastDate
		updatedTopicMonitorParam.CommentsCount = kitWndTCommentsCount.Spinbox.Value()
		updatedTopicMonitorParam.TopicsCount = kitWndTTopicsCount.Spinbox.Value()

		err = updatedTopicMonitorParam.updateInDB()
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
	var wallPostCommentMonitorParam WallPostCommentMonitorParam
	err := wallPostCommentMonitorParam.selectFromDBBySubjectID(IDSubject)
	if err != nil {
		ToLogFile(err.Error(), string(debug.Stack()))
		panic(err.Error())
	}

	// получаем набор для отображения установок модуля мониторинга комментариев под постами
	windowTitle := fmt.Sprintf("%v: настройки для %v", btnName, nameSubject)
	kitWindowWallPostCommentSettings := MakeSettingWindowKit(windowTitle, 300, 100)

	boxWndWPC := ui.NewVerticalBox()

	// получаем набор для флага необходимости активировать модуль мониторинга
	kitWndWPCMonitoring := MakeSettingCheckboxKit("Разрешить наблюдение", wallPostCommentMonitorParam.NeedMonitoring)

	// получаем набор для интервала между запусками функции мониторинга
	kitWndWPCInterval := MakeSettingSpinboxKit("Интервал (сек.)", 5, 21600, wallPostCommentMonitorParam.Interval)

	// получаем набор для идентификатора получателя сообщений
	kitWndWPCSendTo := MakeSettingEntryKit("Получатель", strconv.Itoa(wallPostCommentMonitorParam.SendTo))
	// и привязываем его к обработке ввода, чтобы кроме чисел ничего не вводилось
	kitWndWPCSendTo.Entry.OnChanged(func(*ui.Entry) {
		NumericEntriesHandler(kitWndWPCSendTo.Entry)
	})

	// получаем набор для количества проверяемых постов
	kitWndWPCPostsCount := MakeSettingSpinboxKit("Количество постов", 1, 100, wallPostCommentMonitorParam.PostsCount)

	// получаем набор для количества проверяемых комментариев
	kitWndWPCCommentsCount := MakeSettingSpinboxKit("Количество комментариев", 1, 100, wallPostCommentMonitorParam.CommentsCount)

	// получаем набор для фильтров постов для проверки комментариев
	listPostsFilters := []string{"Все", "От пользователей", "От сообщества"}
	listPostsFiltersEn := []string{"all", "others", "owner"}
	var currentFilter string
	for i, item := range listPostsFiltersEn {
		if wallPostCommentMonitorParam.Filter == item {
			currentFilter = listPostsFilters[i]
			break
		}
	}
	kitWndWPCFilter := MakeSettingComboboxKit("Фильтр", listPostsFilters, currentFilter)

	// получаем набор для флага необходимости проверять все комментарии без исключения
	kitWndWPCMonitoringAll := MakeSettingCheckboxKit("Собирать все комментарии", wallPostCommentMonitorParam.MonitoringAll)

	// получаем набор для флага необходимости проверять комментарии от сообществ
	kitWndWPCMonitorByCommunity := MakeSettingCheckboxKit("Собирать комментарии от сообществ", wallPostCommentMonitorParam.MonitorByCommunity)

	// получаем набор для списка ключевых слов для поиска комментариев
	kitWndWPCKeywordsForMonitoring := MakeSettingEntryListKit("Ключевые слова", wallPostCommentMonitorParam.KeywordsForMonitoring)

	// получаем набор для списка комментариев для поиска
	kitWndWPCSmallCommentsForMonitoring := MakeSettingEntryListKit("Комментарии", wallPostCommentMonitorParam.SmallCommentsForMonitoring)

	// получаем набор для списка имен и фамилий авторов комментариев для поиска комментариев
	kitWndWPCUsersNamesForMonitoring := MakeSettingEntryListKit("Наблюдать комментаторов по имени", wallPostCommentMonitorParam.UsersNamesForMonitoring)

	// получаем набор для списка идентификаторов авторов комментариев для поиска комментариев
	kitWndWPCUsersIdsForMonitoring := MakeSettingEntryListKit("Наблюдать комментаторов по ид-у в ВК", wallPostCommentMonitorParam.UsersIDsForMonitoring)

	// получаем набор для списка идентификаторов авторов комментариев для их игнорирования при проверке комментариев
	kitWndWPCUsersIdsForIgnore := MakeSettingEntryListKit("Игнорировать комментаторов по ид-у в ВК", wallPostCommentMonitorParam.UsersIDsForIgnore)

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
			warningTitle := "Поле \"Получатель\" не должно быть пустым."
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
			warningTitle := "Нужно выбрать элемент в списке \"Фильтр\""
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

		err = updatedWallPostCommentMonitorParam.updateInDB()
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
