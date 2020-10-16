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

type subjectAdditionalSettingsBox struct {
	box *ui.Box
}

func (sasb *subjectAdditionalSettingsBox) init() {
	sasb.box = ui.NewVerticalBox()
}

func (sasb *subjectAdditionalSettingsBox) appendBtnsBox(sasbb subjectAdditionalSettingsBtnBox) {
	sasb.box.Append(sasbb.box, false)
}

type subjectAdditionalSettingsBtnBox struct {
	box    *ui.Box
	label  *ui.Label
	button *ui.Button
}

func (sasbb *subjectAdditionalSettingsBtnBox) init() {
	sasbb.box = ui.NewHorizontalBox()
}

func (sasbb *subjectAdditionalSettingsBtnBox) initLabel(title string) {
	sasbb.label = ui.NewLabel(title)
	sasbb.box.Append(sasbb.label, true)
}

func (sasbb *subjectAdditionalSettingsBtnBox) initButton() {
	sasbb.button = ui.NewButton("Настройки...")
	sasbb.box.Append(sasbb.button, true)
}

func makeSubjectAdditionalSettingsBox(subjectData Subject) *ui.Box {
	var sasb subjectAdditionalSettingsBox
	sasb.init()

	var generalSettingsBtnBox subjectAdditionalSettingsBtnBox
	generalSettingsBtnBox.init()
	generalSettingsBtnBox.initLabel("Общие")
	generalSettingsBtnBox.initButton()
	generalSettingsBtnBox.button.OnClicked(func(*ui.Button) {
		showSubjectGeneralSettingWindow(subjectData.ID, "Общие")
	})
	sasb.appendBtnsBox(generalSettingsBtnBox)

	var wallPostSettingsBtnBox subjectAdditionalSettingsBtnBox
	wallPostSettingsBtnBox.init()
	wallPostSettingsBtnBox.initLabel("Наблюдатель постов на стене")
	wallPostSettingsBtnBox.initButton()
	wallPostSettingsBtnBox.button.OnClicked(func(*ui.Button) {
		showSubjectWallPostSettingWindow(subjectData.ID, subjectData.Name, "Наблюдатель постов на стене")
	})
	sasb.appendBtnsBox(wallPostSettingsBtnBox)

	var albumPhotoSettingsBtnBox subjectAdditionalSettingsBtnBox
	albumPhotoSettingsBtnBox.init()
	albumPhotoSettingsBtnBox.initLabel("Наблюдатель фото в альбомах")
	albumPhotoSettingsBtnBox.initButton()
	albumPhotoSettingsBtnBox.button.OnClicked(func(*ui.Button) {
		showSubjectAlbumPhotoSettingWindow(subjectData.ID, subjectData.Name, "Наблюдатель фото в альбомах")
	})
	sasb.appendBtnsBox(albumPhotoSettingsBtnBox)

	var videoSettingsBtnBox subjectAdditionalSettingsBtnBox
	videoSettingsBtnBox.init()
	videoSettingsBtnBox.initLabel("Наблюдатель видео в альбомах")
	videoSettingsBtnBox.initButton()
	videoSettingsBtnBox.button.OnClicked(func(*ui.Button) {
		showSubjectVideoSettingWindow(subjectData.ID, subjectData.Name, "Наблюдатель видео в альбомах")
	})
	sasb.appendBtnsBox(videoSettingsBtnBox)

	var photoCommentSettingsBtnBox subjectAdditionalSettingsBtnBox
	photoCommentSettingsBtnBox.init()
	photoCommentSettingsBtnBox.initLabel("Наблюдатель комментариев под фото")
	photoCommentSettingsBtnBox.initButton()
	photoCommentSettingsBtnBox.button.OnClicked(func(*ui.Button) {
		showSubjectPhotoCommentSettingWindow(subjectData.ID, subjectData.Name, "Наблюдатель комментариев под фото")
	})
	sasb.appendBtnsBox(photoCommentSettingsBtnBox)

	var videoCommentSettingsBtnBox subjectAdditionalSettingsBtnBox
	videoCommentSettingsBtnBox.init()
	videoCommentSettingsBtnBox.initLabel("Наблюдатель комментариев под видео")
	videoCommentSettingsBtnBox.initButton()
	videoCommentSettingsBtnBox.button.OnClicked(func(*ui.Button) {
		showSubjectVideoCommentSettingWindow(subjectData.ID, subjectData.Name, "Наблюдатель комментариев под видео")
	})
	sasb.appendBtnsBox(videoCommentSettingsBtnBox)

	var topicSettingsBtnBox subjectAdditionalSettingsBtnBox
	topicSettingsBtnBox.init()
	topicSettingsBtnBox.initLabel("Наблюдатель комментариев в обсуждениях")
	topicSettingsBtnBox.initButton()
	topicSettingsBtnBox.button.OnClicked(func(*ui.Button) {
		showSubjectTopicSettingWindow(subjectData.ID, subjectData.Name, "Наблюдатель комментариев в обсуждениях")
	})
	sasb.appendBtnsBox(topicSettingsBtnBox)

	var wallPostCommentSettingsBtnBox subjectAdditionalSettingsBtnBox
	wallPostCommentSettingsBtnBox.init()
	wallPostCommentSettingsBtnBox.initLabel("Наблюдатель комментариев под постами")
	wallPostCommentSettingsBtnBox.initButton()
	wallPostCommentSettingsBtnBox.button.OnClicked(func(*ui.Button) {
		showSubjectWallPostCommentSettings(subjectData.ID, subjectData.Name, "Наблюдатель комментариев под постами")
	})
	sasb.appendBtnsBox(wallPostCommentSettingsBtnBox)

	return sasb.box
}

type settingsWindow struct {
	window *ui.Window
}

func (sw *settingsWindow) init(windowTitle string, width, height int) {
	sw.window = ui.NewWindow(windowTitle, width, height, true)
	sw.window.OnClosing(func(*ui.Window) bool {
		sw.window.Disable()
		return true
	})
	sw.window.SetMargined(true)
	sw.window.Show()
}

func (sw *settingsWindow) setBox(bsw boxSettingsWnd) {
	sw.window.SetChild(bsw.box)
}

type boxSettingsWnd struct {
	box         *ui.Box
	internalBox *ui.Box
	group       *ui.Group
}

func (bsw *boxSettingsWnd) init() {
	bsw.box = ui.NewVerticalBox()
	bsw.box.SetPadded(true)
}

func (bsw *boxSettingsWnd) initGroup(groupTitle string) {
	bsw.group = ui.NewGroup(groupTitle)
	bsw.group.SetMargined(true)
	bsw.internalBox = ui.NewVerticalBox()
	bsw.internalBox.SetPadded(true)
	bsw.group.SetChild(bsw.internalBox)
	bsw.box.Append(bsw.group, false)
}

func (bsw *boxSettingsWnd) initFlexibleSpaceBox() {
	box := ui.NewHorizontalBox()
	bsw.box.Append(box, false)
}

func (bsw *boxSettingsWnd) appendCheckbox(cbk CheckboxKit) {
	bsw.internalBox.Append(cbk.Box, false)
}

func (bsw *boxSettingsWnd) appendSpinbox(sbk SpinboxKit) {
	bsw.internalBox.Append(sbk.Box, false)
}

func (bsw *boxSettingsWnd) appendEntry(ek EntryKit) {
	bsw.internalBox.Append(ek.Box, false)
}

func (bsw *boxSettingsWnd) appendEntryListKit(elk EntryListKit) { // FIXME: лучше объединить с appendEntry()
	bsw.internalBox.Append(elk.Box, false)
}

func (bsw *boxSettingsWnd) appendCombobox(cbk ComboboxKit) {
	bsw.internalBox.Append(cbk.Box, false)
}

func (bsw *boxSettingsWnd) appendPartedBox(pbsw partedBoxSettingsWnd) {
	bsw.box.Append(pbsw.box, false)
}

func (bsw *boxSettingsWnd) appendButtons(bk ButtonsKit) {
	bsw.box.Append(bk.Box, true)
}

type partedBoxSettingsWnd struct {
	box    *ui.Box
	left   *ui.Box
	center *ui.Box
	right  *ui.Box
}

func (pbsw *partedBoxSettingsWnd) init() {
	pbsw.box = ui.NewHorizontalBox()
	pbsw.box.SetPadded(true)
	pbsw.left = ui.NewVerticalBox()
	pbsw.box.Append(pbsw.left, false)
	pbsw.center = ui.NewVerticalBox()
	pbsw.box.Append(pbsw.center, false)
	pbsw.right = ui.NewVerticalBox()
	pbsw.box.Append(pbsw.right, false)
}

func (pbsw *partedBoxSettingsWnd) appendBoxToLeftPart(bsw boxSettingsWnd) {
	pbsw.left.Append(bsw.box, false)
}

func (pbsw *partedBoxSettingsWnd) appendBoxToCenterPart(bsw boxSettingsWnd) {
	pbsw.center.Append(bsw.box, false)
}

func (pbsw *partedBoxSettingsWnd) appendBoxToRightPart(bsw boxSettingsWnd) {
	pbsw.right.Append(bsw.box, false)
}

type subjectGeneralSettingsParams struct {
	kitSubjectName EntryKit
	kitSubjectID   EntryKit
}

func (sgsp *subjectGeneralSettingsParams) init(subjectParams Subject) {
	sgsp.kitSubjectName = MakeSettingEntryKit("Название", subjectParams.Name)
	sgsp.kitSubjectID = MakeSettingEntryKit("Идентификатор в ВК", strconv.Itoa(subjectParams.SubjectID))
	sgsp.kitSubjectID.Entry.OnChanged(func(*ui.Entry) {
		NumericEntriesHandler(sgsp.kitSubjectID.Entry)
	})
}

func (sgsp *subjectGeneralSettingsParams) updateParamsInDB(subjectParams Subject) bool {
	var updatedParams Subject
	updatedParams.ID = subjectParams.ID
	if len(sgsp.kitSubjectID.Entry.Text()) == 0 {
		warningTitle := "Поле \"Идентификатор в ВК\" не должно быть пустым"
		ShowWarningWindow(warningTitle)
		return false
	}
	var err error
	updatedParams.SubjectID, err = strconv.Atoi(sgsp.kitSubjectID.Entry.Text())
	if err != nil {
		ToLogFile(err.Error(), string(debug.Stack()))
		panic(err.Error())
	}
	if len(sgsp.kitSubjectName.Entry.Text()) == 0 {
		warningTitle := "Поле \"Название\" не должно быть пустым"
		ShowWarningWindow(warningTitle)
		return false
	}
	updatedParams.Name = sgsp.kitSubjectName.Entry.Text()
	updatedParams.BackupWikipage = subjectParams.BackupWikipage
	updatedParams.LastBackup = subjectParams.LastBackup
	err = updatedParams.updateInDB()
	if err != nil {
		ToLogFile(err.Error(), string(debug.Stack()))
		panic(err.Error())
	}

	return true
}

type accessTokensWallPostSettingsParams struct {
	kitWallGetAccessToken       ComboboxKit
	kitUsersGetAccessToken      ComboboxKit
	kitGroupsGetByIDAccessToken ComboboxKit
	kitMessagesSendAccessToken  ComboboxKit
}

func (atwpsp *accessTokensWallPostSettingsParams) init(spfdb settingsParamsFromDB) {
	atNames := spfdb.accessTokensNames()
	wgSelectedAccessToken := spfdb.selectedAccessToken(spfdb.wallGetMethodParams)
	atwpsp.kitWallGetAccessToken = MakeSettingComboboxKit("Ключ доступа для \"wall.get\"", atNames, wgSelectedAccessToken)
	ugSelectedAccessToken := spfdb.selectedAccessToken(spfdb.usersGetMethodParams)
	atwpsp.kitUsersGetAccessToken = MakeSettingComboboxKit("Ключ доступа для \"users.get\"", atNames, ugSelectedAccessToken)
	ggbiSelectedAccessToken := spfdb.selectedAccessToken(spfdb.groupsGetByIDMethodParams)
	atwpsp.kitGroupsGetByIDAccessToken = MakeSettingComboboxKit("Ключ доступа для \"groups.getById\"", atNames, ggbiSelectedAccessToken)
	msSelectedAccessToken := spfdb.selectedAccessToken(spfdb.messagesSendMethodParams)
	atwpsp.kitMessagesSendAccessToken = MakeSettingComboboxKit("Ключ доступа для \"messages.send\"", atNames, msSelectedAccessToken)
}

func (atwpsp *accessTokensWallPostSettingsParams) updateParamsInDB(spfdb settingsParamsFromDB) {
	spfdb.updateMethodInDB(spfdb.wallGetMethodParams, atwpsp.kitWallGetAccessToken)
	spfdb.updateMethodInDB(spfdb.usersGetMethodParams, atwpsp.kitUsersGetAccessToken)
	spfdb.updateMethodInDB(spfdb.groupsGetByIDMethodParams, atwpsp.kitGroupsGetByIDAccessToken)
	spfdb.updateMethodInDB(spfdb.messagesSendMethodParams, atwpsp.kitMessagesSendAccessToken)
}

type accessTokensAlbumPhotoSettingsParams struct {
	kitPhotosGetAccessToken       ComboboxKit
	kitPhotosGetAlbumsAccessToken ComboboxKit
	kitUsersGetAccessToken        ComboboxKit
	kitGroupsGetByIDAccessToken   ComboboxKit
	kitMessagesSendAccessToken    ComboboxKit
}

func (atapsp *accessTokensAlbumPhotoSettingsParams) init(spfdb settingsParamsFromDB) {
	atNames := spfdb.accessTokensNames()
	pgSelectedAccessToken := spfdb.selectedAccessToken(spfdb.photosGetMethodParams)
	atapsp.kitPhotosGetAccessToken = MakeSettingComboboxKit("Ключ доступа для \"photos.get\"", atNames, pgSelectedAccessToken)
	pgaSelectedAccessToken := spfdb.selectedAccessToken(spfdb.photosGetAlbumsMethodParams)
	atapsp.kitPhotosGetAlbumsAccessToken = MakeSettingComboboxKit("Ключ доступа для \"photos.getAlbums\"", atNames, pgaSelectedAccessToken)
	ugSelectedAccessToken := spfdb.selectedAccessToken(spfdb.usersGetMethodParams)
	atapsp.kitUsersGetAccessToken = MakeSettingComboboxKit("Ключ доступа для \"users.get\"", atNames, ugSelectedAccessToken)
	ggbiSelectedAccessToken := spfdb.selectedAccessToken(spfdb.groupsGetByIDMethodParams)
	atapsp.kitGroupsGetByIDAccessToken = MakeSettingComboboxKit("Ключ доступа для \"groups.getById\"", atNames, ggbiSelectedAccessToken)
	msSelectedAccessToken := spfdb.selectedAccessToken(spfdb.messagesSendMethodParams)
	atapsp.kitMessagesSendAccessToken = MakeSettingComboboxKit("Ключ доступа для \"messages.send\"", atNames, msSelectedAccessToken)
}

func (atapsp *accessTokensAlbumPhotoSettingsParams) updateParamsInDB(spfdb settingsParamsFromDB) {
	spfdb.updateMethodInDB(spfdb.photosGetMethodParams, atapsp.kitPhotosGetAccessToken)
	spfdb.updateMethodInDB(spfdb.photosGetAlbumsMethodParams, atapsp.kitPhotosGetAlbumsAccessToken)
	spfdb.updateMethodInDB(spfdb.usersGetMethodParams, atapsp.kitUsersGetAccessToken)
	spfdb.updateMethodInDB(spfdb.groupsGetByIDMethodParams, atapsp.kitGroupsGetByIDAccessToken)
	spfdb.updateMethodInDB(spfdb.messagesSendMethodParams, atapsp.kitMessagesSendAccessToken)
}

type accessTokensVideoSettingsParams struct {
	kitVideoGetAccessToken      ComboboxKit
	kitUsersGetAccessToken      ComboboxKit
	kitGroupsGetByIDAccessToken ComboboxKit
	kitMessagesSendAccessToken  ComboboxKit
}

func (atvsp *accessTokensVideoSettingsParams) init(spfdb settingsParamsFromDB) {
	atNames := spfdb.accessTokensNames()
	vgSelectedAccessToken := spfdb.selectedAccessToken(spfdb.videoGetMethodParams)
	atvsp.kitVideoGetAccessToken = MakeSettingComboboxKit("Ключ доступа для \"video.get\"", atNames, vgSelectedAccessToken)
	ugSelectedAccessToken := spfdb.selectedAccessToken(spfdb.usersGetMethodParams)
	atvsp.kitUsersGetAccessToken = MakeSettingComboboxKit("Ключ доступа для \"users.get\"", atNames, ugSelectedAccessToken)
	ggbiSelectedAccessToken := spfdb.selectedAccessToken(spfdb.groupsGetByIDMethodParams)
	atvsp.kitGroupsGetByIDAccessToken = MakeSettingComboboxKit("Ключ доступа для \"groups.getById\"", atNames, ggbiSelectedAccessToken)
	msSelectedAccessToken := spfdb.selectedAccessToken(spfdb.messagesSendMethodParams)
	atvsp.kitMessagesSendAccessToken = MakeSettingComboboxKit("Ключ доступа для \"messages.send\"", atNames, msSelectedAccessToken)
}

func (atvsp *accessTokensVideoSettingsParams) updateParamsInDB(spfdb settingsParamsFromDB) {
	spfdb.updateMethodInDB(spfdb.videoGetMethodParams, atvsp.kitVideoGetAccessToken)
	spfdb.updateMethodInDB(spfdb.usersGetMethodParams, atvsp.kitUsersGetAccessToken)
	spfdb.updateMethodInDB(spfdb.groupsGetByIDMethodParams, atvsp.kitGroupsGetByIDAccessToken)
	spfdb.updateMethodInDB(spfdb.messagesSendMethodParams, atvsp.kitMessagesSendAccessToken)
}

type accessTokensPhotoCommentSettingsParams struct {
	kitPhotosGetAllCommentsAccessToken ComboboxKit
	kitUsersGetAccessToken             ComboboxKit
	kitGroupsGetByIDAccessToken        ComboboxKit
	kitMessagesSendAccessToken         ComboboxKit
}

func (atpcsp *accessTokensPhotoCommentSettingsParams) init(spfdb settingsParamsFromDB) {
	atNames := spfdb.accessTokensNames()
	pgacSelectedAccessToken := spfdb.selectedAccessToken(spfdb.photosGetAllCommentsMethodParams)
	atpcsp.kitPhotosGetAllCommentsAccessToken = MakeSettingComboboxKit("Ключ доступа для \"photos.getAllComments\"", atNames, pgacSelectedAccessToken)
	ugSelectedAccessToken := spfdb.selectedAccessToken(spfdb.usersGetMethodParams)
	atpcsp.kitUsersGetAccessToken = MakeSettingComboboxKit("Ключ доступа для \"users.get\"", atNames, ugSelectedAccessToken)
	ggbiSelectedAccessToken := spfdb.selectedAccessToken(spfdb.groupsGetByIDMethodParams)
	atpcsp.kitGroupsGetByIDAccessToken = MakeSettingComboboxKit("Ключ доступа для \"groups.getById\"", atNames, ggbiSelectedAccessToken)
	msSelectedAccessToken := spfdb.selectedAccessToken(spfdb.messagesSendMethodParams)
	atpcsp.kitMessagesSendAccessToken = MakeSettingComboboxKit("Ключ доступа для \"messages.send\"", atNames, msSelectedAccessToken)
}

func (atpcsp *accessTokensPhotoCommentSettingsParams) updateParamsInDB(spfdb settingsParamsFromDB) {
	spfdb.updateMethodInDB(spfdb.photosGetAllCommentsMethodParams, atpcsp.kitPhotosGetAllCommentsAccessToken)
	spfdb.updateMethodInDB(spfdb.usersGetMethodParams, atpcsp.kitUsersGetAccessToken)
	spfdb.updateMethodInDB(spfdb.groupsGetByIDMethodParams, atpcsp.kitGroupsGetByIDAccessToken)
	spfdb.updateMethodInDB(spfdb.messagesSendMethodParams, atpcsp.kitMessagesSendAccessToken)
}

type accessTokensVideoCommentSettingsParams struct {
	kitVideoGetCommentsAccessToken ComboboxKit
	kitUsersGetAccessToken         ComboboxKit
	kitGroupsGetByIDAccessToken    ComboboxKit
	kitVideoGetAccessToken         ComboboxKit
	kitMessagesSendAccessToken     ComboboxKit
}

func (atvcsp *accessTokensVideoCommentSettingsParams) init(spfdb settingsParamsFromDB) {
	atNames := spfdb.accessTokensNames()
	vgcSelectedAccessToken := spfdb.selectedAccessToken(spfdb.videoGetCommentsMethodParams)
	atvcsp.kitVideoGetCommentsAccessToken = MakeSettingComboboxKit("Ключ доступа для \"video.getComments\"", atNames, vgcSelectedAccessToken)
	ugSelectedAccessToken := spfdb.selectedAccessToken(spfdb.usersGetMethodParams)
	atvcsp.kitUsersGetAccessToken = MakeSettingComboboxKit("Ключ доступа для \"users.get\"", atNames, ugSelectedAccessToken)
	ggbiSelectedAccessToken := spfdb.selectedAccessToken(spfdb.groupsGetByIDMethodParams)
	atvcsp.kitGroupsGetByIDAccessToken = MakeSettingComboboxKit("Ключ доступа для \"groups.getById\"", atNames, ggbiSelectedAccessToken)
	vgSelectedAccessToken := spfdb.selectedAccessToken(spfdb.videoGetMethodParams)
	atvcsp.kitVideoGetAccessToken = MakeSettingComboboxKit("Ключ доступа для \"video.get\"", atNames, vgSelectedAccessToken)
	msSelectedAccessToken := spfdb.selectedAccessToken(spfdb.messagesSendMethodParams)
	atvcsp.kitMessagesSendAccessToken = MakeSettingComboboxKit("Ключ доступа для \"messages.send\"", atNames, msSelectedAccessToken)
}

func (atvcsp *accessTokensVideoCommentSettingsParams) updateParamsInDB(spfdb settingsParamsFromDB) {
	spfdb.updateMethodInDB(spfdb.videoGetCommentsMethodParams, atvcsp.kitVideoGetCommentsAccessToken)
	spfdb.updateMethodInDB(spfdb.usersGetMethodParams, atvcsp.kitUsersGetAccessToken)
	spfdb.updateMethodInDB(spfdb.groupsGetByIDMethodParams, atvcsp.kitGroupsGetByIDAccessToken)
	spfdb.updateMethodInDB(spfdb.videoGetMethodParams, atvcsp.kitVideoGetAccessToken)
	spfdb.updateMethodInDB(spfdb.messagesSendMethodParams, atvcsp.kitMessagesSendAccessToken)
}

type accessTokensTopicSettingsParams struct {
	kitBoardGetCommentsAccessToken ComboboxKit
	kitBoardGetTopicsAccessToken   ComboboxKit
	kitUsersGetAccessToken         ComboboxKit
	kitGroupsGetByIDAccessToken    ComboboxKit
	kitMessagesSendAccessToken     ComboboxKit
}

func (attsp *accessTokensTopicSettingsParams) init(spfdb settingsParamsFromDB) {
	atNames := spfdb.accessTokensNames()
	bgcSelectedAccessToken := spfdb.selectedAccessToken(spfdb.boardGetCommentsMethodParams)
	attsp.kitBoardGetCommentsAccessToken = MakeSettingComboboxKit("Ключ доступа для \"board.getComments\"", atNames, bgcSelectedAccessToken)
	bgtSelectedAccessToken := spfdb.selectedAccessToken(spfdb.boardGetTopicsMethodParams)
	attsp.kitBoardGetTopicsAccessToken = MakeSettingComboboxKit("Ключ доступа для \"board.getTopics\"", atNames, bgtSelectedAccessToken)
	ugSelectedAccessToken := spfdb.selectedAccessToken(spfdb.usersGetMethodParams)
	attsp.kitUsersGetAccessToken = MakeSettingComboboxKit("Ключ доступа для \"users.get\"", atNames, ugSelectedAccessToken)
	ggbiSelectedAccessToken := spfdb.selectedAccessToken(spfdb.groupsGetByIDMethodParams)
	attsp.kitGroupsGetByIDAccessToken = MakeSettingComboboxKit("Ключ доступа для \"groups.getById\"", atNames, ggbiSelectedAccessToken)
	msSelectedAccessToken := spfdb.selectedAccessToken(spfdb.messagesSendMethodParams)
	attsp.kitMessagesSendAccessToken = MakeSettingComboboxKit("Ключ доступа для \"messages.send\"", atNames, msSelectedAccessToken)
}

func (attsp *accessTokensTopicSettingsParams) updateParamsInDB(spfdb settingsParamsFromDB) {
	spfdb.updateMethodInDB(spfdb.boardGetCommentsMethodParams, attsp.kitBoardGetCommentsAccessToken)
	spfdb.updateMethodInDB(spfdb.boardGetTopicsMethodParams, attsp.kitBoardGetTopicsAccessToken)
	spfdb.updateMethodInDB(spfdb.usersGetMethodParams, attsp.kitUsersGetAccessToken)
	spfdb.updateMethodInDB(spfdb.groupsGetByIDMethodParams, attsp.kitGroupsGetByIDAccessToken)
	spfdb.updateMethodInDB(spfdb.messagesSendMethodParams, attsp.kitMessagesSendAccessToken)
}

type accessTokensWallPostCommentSettingsParams struct {
	kitWallGetCommentsAccessToken ComboboxKit
	kitUsersGetAccessToken        ComboboxKit
	kitGroupsGetByIDAccessToken   ComboboxKit
	kitWallGetAccessToken         ComboboxKit
	kitWallGetCommentAccessToken  ComboboxKit
	kitMessagesSendAccessToken    ComboboxKit
}

func (atwpcsp *accessTokensWallPostCommentSettingsParams) init(spfdb settingsParamsFromDB) {
	atNames := spfdb.accessTokensNames()
	wgcsSelectedAccessToken := spfdb.selectedAccessToken(spfdb.wallGetCommentsMethodParams)
	atwpcsp.kitWallGetCommentsAccessToken = MakeSettingComboboxKit("Ключ доступа для \"wall.getComments\"", atNames, wgcsSelectedAccessToken)
	ugSelectedAccessToken := spfdb.selectedAccessToken(spfdb.usersGetMethodParams)
	atwpcsp.kitUsersGetAccessToken = MakeSettingComboboxKit("Ключ доступа для \"users.get\"", atNames, ugSelectedAccessToken)
	ggbiSelectedAccessToken := spfdb.selectedAccessToken(spfdb.groupsGetByIDMethodParams)
	atwpcsp.kitGroupsGetByIDAccessToken = MakeSettingComboboxKit("Ключ доступа для \"groups.getById\"", atNames, ggbiSelectedAccessToken)
	wgSelectedAccessToken := spfdb.selectedAccessToken(spfdb.wallGetMethodParams)
	atwpcsp.kitWallGetAccessToken = MakeSettingComboboxKit("Ключ доступа для \"wall.get\"", atNames, wgSelectedAccessToken)
	wgcSelectedAccessToken := spfdb.selectedAccessToken(spfdb.wallGetCommentMethodParams)
	atwpcsp.kitWallGetCommentAccessToken = MakeSettingComboboxKit("Ключ доступа для \"wall.getComment\"", atNames, wgcSelectedAccessToken)
	msSelectedAccessToken := spfdb.selectedAccessToken(spfdb.messagesSendMethodParams)
	atwpcsp.kitMessagesSendAccessToken = MakeSettingComboboxKit("Ключ доступа для \"messages.send\"", atNames, msSelectedAccessToken)
}

func (atwpcsp *accessTokensWallPostCommentSettingsParams) updateParamsInDB(spfdb settingsParamsFromDB) {
	spfdb.updateMethodInDB(spfdb.wallGetCommentsMethodParams, atwpcsp.kitWallGetCommentsAccessToken)
	spfdb.updateMethodInDB(spfdb.usersGetMethodParams, atwpcsp.kitUsersGetAccessToken)
	spfdb.updateMethodInDB(spfdb.groupsGetByIDMethodParams, atwpcsp.kitGroupsGetByIDAccessToken)
	spfdb.updateMethodInDB(spfdb.wallGetMethodParams, atwpcsp.kitWallGetAccessToken)
	spfdb.updateMethodInDB(spfdb.wallGetCommentMethodParams, atwpcsp.kitWallGetCommentAccessToken)
	spfdb.updateMethodInDB(spfdb.messagesSendMethodParams, atwpcsp.kitMessagesSendAccessToken)
}

type settingsParamsFromDB struct {
	subjectParams                    Subject
	accessTokens                     []AccessToken
	monitorParams                    Monitor
	usersGetMethodParams             Method
	groupsGetByIDMethodParams        Method
	wallGetMethodParams              Method
	wallGetCommentMethodParams       Method
	wallGetCommentsMethodParams      Method
	photosGetMethodParams            Method
	photosGetAlbumsMethodParams      Method
	photosGetAllCommentsMethodParams Method
	videoGetMethodParams             Method
	videoGetCommentsMethodParams     Method
	boardGetTopicsMethodParams       Method
	boardGetCommentsMethodParams     Method
	messagesSendMethodParams         Method
}

func (spfdb *settingsParamsFromDB) setSubjectParams(subjectParams Subject) {
	spfdb.subjectParams = subjectParams
}

func (spfdb *settingsParamsFromDB) setAccessTokens(accessTokens []AccessToken) {
	spfdb.accessTokens = accessTokens
}

func (spfdb *settingsParamsFromDB) selectMonitorParams(monitorName string) {
	err := spfdb.monitorParams.selectFromDBByNameAndBySubjectID(monitorName, spfdb.subjectParams.ID)
	if err != nil {
		ToLogFile(err.Error(), string(debug.Stack()))
		panic(err.Error())
	}
}

func (spfdb *settingsParamsFromDB) selectUsersGetMethodParams() {
	err := spfdb.usersGetMethodParams.selectFromDBByNameAndBySubjectIDAndByMonitorID("users.get", spfdb.subjectParams.ID, spfdb.monitorParams.ID)
	if err != nil {
		ToLogFile(err.Error(), string(debug.Stack()))
		panic(err.Error())
	}
}

func (spfdb *settingsParamsFromDB) selectGroupsGetByIDMethodParams() {
	err := spfdb.groupsGetByIDMethodParams.selectFromDBByNameAndBySubjectIDAndByMonitorID("groups.getById", spfdb.subjectParams.ID, spfdb.monitorParams.ID)
	if err != nil {
		ToLogFile(err.Error(), string(debug.Stack()))
		panic(err.Error())
	}
}

func (spfdb *settingsParamsFromDB) selectWallGetMethodParams() {
	err := spfdb.wallGetMethodParams.selectFromDBByNameAndBySubjectIDAndByMonitorID("wall.get", spfdb.subjectParams.ID, spfdb.monitorParams.ID)
	if err != nil {
		ToLogFile(err.Error(), string(debug.Stack()))
		panic(err.Error())
	}
}

func (spfdb *settingsParamsFromDB) selectWallGetCommentMethodParams() {
	err := spfdb.wallGetCommentMethodParams.selectFromDBByNameAndBySubjectIDAndByMonitorID("wall.getComment", spfdb.subjectParams.ID, spfdb.monitorParams.ID)
	if err != nil {
		ToLogFile(err.Error(), string(debug.Stack()))
		panic(err.Error())
	}
}

func (spfdb *settingsParamsFromDB) selectWallGetCommentsMethodParams() {
	err := spfdb.wallGetCommentsMethodParams.selectFromDBByNameAndBySubjectIDAndByMonitorID("wall.getComments", spfdb.subjectParams.ID, spfdb.monitorParams.ID)
	if err != nil {
		ToLogFile(err.Error(), string(debug.Stack()))
		panic(err.Error())
	}
}

func (spfdb *settingsParamsFromDB) selectPhotosGetMethodParams() {
	err := spfdb.photosGetMethodParams.selectFromDBByNameAndBySubjectIDAndByMonitorID("photos.get", spfdb.subjectParams.ID, spfdb.monitorParams.ID)
	if err != nil {
		ToLogFile(err.Error(), string(debug.Stack()))
		panic(err.Error())
	}
}

func (spfdb *settingsParamsFromDB) selectPhotosGetAlbumsMethodParams() {
	err := spfdb.photosGetAlbumsMethodParams.selectFromDBByNameAndBySubjectIDAndByMonitorID("photos.getAlbums", spfdb.subjectParams.ID, spfdb.monitorParams.ID)
	if err != nil {
		ToLogFile(err.Error(), string(debug.Stack()))
		panic(err.Error())
	}
}

func (spfdb *settingsParamsFromDB) selectPhotosGetAllCommentsMethodParams() {
	err := spfdb.photosGetAllCommentsMethodParams.selectFromDBByNameAndBySubjectIDAndByMonitorID("photos.getAllComments", spfdb.subjectParams.ID, spfdb.monitorParams.ID)
	if err != nil {
		ToLogFile(err.Error(), string(debug.Stack()))
		panic(err.Error())
	}
}

func (spfdb *settingsParamsFromDB) selectVideoGetMethodParams() {
	err := spfdb.videoGetMethodParams.selectFromDBByNameAndBySubjectIDAndByMonitorID("video.get", spfdb.subjectParams.ID, spfdb.monitorParams.ID)
	if err != nil {
		ToLogFile(err.Error(), string(debug.Stack()))
		panic(err.Error())
	}
}

func (spfdb *settingsParamsFromDB) selectVideoGetCommentsMethodParams() {
	err := spfdb.videoGetCommentsMethodParams.selectFromDBByNameAndBySubjectIDAndByMonitorID("video.getComments", spfdb.subjectParams.ID, spfdb.monitorParams.ID)
	if err != nil {
		ToLogFile(err.Error(), string(debug.Stack()))
		panic(err.Error())
	}
}

func (spfdb *settingsParamsFromDB) selectBoardGetTopicsMethodParams() {
	err := spfdb.boardGetTopicsMethodParams.selectFromDBByNameAndBySubjectIDAndByMonitorID("board.getTopics", spfdb.subjectParams.ID, spfdb.monitorParams.ID)
	if err != nil {
		ToLogFile(err.Error(), string(debug.Stack()))
		panic(err.Error())
	}
}

func (spfdb *settingsParamsFromDB) selectBoardGetCommentsMethodParams() {
	err := spfdb.boardGetCommentsMethodParams.selectFromDBByNameAndBySubjectIDAndByMonitorID("board.getComments", spfdb.subjectParams.ID, spfdb.monitorParams.ID)
	if err != nil {
		ToLogFile(err.Error(), string(debug.Stack()))
		panic(err.Error())
	}
}

func (spfdb *settingsParamsFromDB) selectMessagesSendMethodParams() {
	err := spfdb.messagesSendMethodParams.selectFromDBByNameAndBySubjectIDAndByMonitorID("messages.send", spfdb.subjectParams.ID, spfdb.monitorParams.ID)
	if err != nil {
		ToLogFile(err.Error(), string(debug.Stack()))
		panic(err.Error())
	}
}

func (spfdb *settingsParamsFromDB) accessTokensNames() []string {
	var accessTokensNames []string
	for _, accessToken := range spfdb.accessTokens {
		accessTokensNames = append(accessTokensNames, accessToken.Name)
	}
	return accessTokensNames
}

func (spfdb *settingsParamsFromDB) selectedAccessToken(methodParams Method) string {
	var currentAccessToken string
	for _, accessToken := range spfdb.accessTokens {
		if methodParams.AccessTokenID == accessToken.ID {
			currentAccessToken = accessToken.Name
		}
	}
	return currentAccessToken
}

func (spfdb *settingsParamsFromDB) updateMethodInDB(methodParams Method, kitMethodAccessToken ComboboxKit) {
	var updatedMethodParams Method
	updatedMethodParams.ID = methodParams.ID
	updatedMethodParams.Name = methodParams.Name
	updatedMethodParams.SubjectID = methodParams.SubjectID
	updatedMethodParams.MonitorID = methodParams.MonitorID
	updatedMethodParams.AccessTokenID = spfdb.accessTokens[kitMethodAccessToken.Combobox.Selected()].ID
	err := updatedMethodParams.updateInDB()
	if err != nil {
		ToLogFile(err.Error(), string(debug.Stack()))
		panic(err.Error())
	}
}

func showSubjectGeneralSettingWindow(IDSubject int, btnName string) {
	var subjectParams Subject
	err := subjectParams.selectFromDBByID(IDSubject)
	if err != nil {
		ToLogFile(err.Error(), string(debug.Stack()))
		panic(err.Error())
	}

	var dbKit DataBaseKit
	accessTokens, err := dbKit.selectTableAccessToken()
	if err != nil {
		ToLogFile(err.Error(), string(debug.Stack()))
		panic(err.Error())
	}

	windowTitle := fmt.Sprintf("%v: настройки для %v", btnName, subjectParams.Name)
	var sw settingsWindow
	sw.init(windowTitle, 300, 100)

	var partedBsw partedBoxSettingsWnd
	partedBsw.init()

	var sgsp subjectGeneralSettingsParams
	sgsp.init(subjectParams)
	var generalMonitorBsw boxSettingsWnd
	generalMonitorBsw.init()
	generalMonitorBsw.initGroup("Обшие")
	generalMonitorBsw.appendEntry(sgsp.kitSubjectName)
	generalMonitorBsw.appendEntry(sgsp.kitSubjectID)
	partedBsw.appendBoxToLeftPart(generalMonitorBsw)

	var wallPostMonitorSpfdb settingsParamsFromDB
	wallPostMonitorSpfdb.setSubjectParams(subjectParams)
	wallPostMonitorSpfdb.setAccessTokens(accessTokens)
	wallPostMonitorSpfdb.selectMonitorParams("wall_post_monitor")
	wallPostMonitorSpfdb.selectWallGetMethodParams()
	wallPostMonitorSpfdb.selectUsersGetMethodParams()
	wallPostMonitorSpfdb.selectGroupsGetByIDMethodParams()
	wallPostMonitorSpfdb.selectMessagesSendMethodParams()
	var atwpsp accessTokensWallPostSettingsParams
	atwpsp.init(wallPostMonitorSpfdb)
	var wallPostMonitorBsw boxSettingsWnd
	wallPostMonitorBsw.init()
	wallPostMonitorBsw.initGroup("Посты на стене")
	wallPostMonitorBsw.appendCombobox(atwpsp.kitWallGetAccessToken)
	wallPostMonitorBsw.appendCombobox(atwpsp.kitUsersGetAccessToken)
	wallPostMonitorBsw.appendCombobox(atwpsp.kitGroupsGetByIDAccessToken)
	wallPostMonitorBsw.appendCombobox(atwpsp.kitMessagesSendAccessToken)
	partedBsw.appendBoxToLeftPart(wallPostMonitorBsw)

	var albumPhotoMonitorSpfdb settingsParamsFromDB
	albumPhotoMonitorSpfdb.setSubjectParams(subjectParams)
	albumPhotoMonitorSpfdb.setAccessTokens(accessTokens)
	albumPhotoMonitorSpfdb.selectMonitorParams("album_photo_monitor")
	albumPhotoMonitorSpfdb.selectPhotosGetMethodParams()
	albumPhotoMonitorSpfdb.selectPhotosGetAlbumsMethodParams()
	albumPhotoMonitorSpfdb.selectUsersGetMethodParams()
	albumPhotoMonitorSpfdb.selectGroupsGetByIDMethodParams()
	albumPhotoMonitorSpfdb.selectMessagesSendMethodParams()
	var atapsp accessTokensAlbumPhotoSettingsParams
	atapsp.init(albumPhotoMonitorSpfdb)
	var albumPhotoMonitorBsw boxSettingsWnd
	albumPhotoMonitorBsw.init()
	albumPhotoMonitorBsw.initGroup("Фото в альбомах")
	albumPhotoMonitorBsw.appendCombobox(atapsp.kitPhotosGetAccessToken)
	albumPhotoMonitorBsw.appendCombobox(atapsp.kitPhotosGetAlbumsAccessToken)
	albumPhotoMonitorBsw.appendCombobox(atapsp.kitUsersGetAccessToken)
	albumPhotoMonitorBsw.appendCombobox(atapsp.kitGroupsGetByIDAccessToken)
	albumPhotoMonitorBsw.appendCombobox(atapsp.kitMessagesSendAccessToken)
	partedBsw.appendBoxToLeftPart(albumPhotoMonitorBsw)

	var videoMonitorSpfdb settingsParamsFromDB
	videoMonitorSpfdb.setSubjectParams(subjectParams)
	videoMonitorSpfdb.setAccessTokens(accessTokens)
	videoMonitorSpfdb.selectMonitorParams("video_monitor")
	videoMonitorSpfdb.selectVideoGetMethodParams()
	videoMonitorSpfdb.selectUsersGetMethodParams()
	videoMonitorSpfdb.selectGroupsGetByIDMethodParams()
	videoMonitorSpfdb.selectMessagesSendMethodParams()
	var atvsp accessTokensVideoSettingsParams
	atvsp.init(videoMonitorSpfdb)
	var videoMonitorBsw boxSettingsWnd
	videoMonitorBsw.init()
	videoMonitorBsw.initGroup("Видео в альбомах")
	videoMonitorBsw.appendCombobox(atvsp.kitVideoGetAccessToken)
	videoMonitorBsw.appendCombobox(atvsp.kitUsersGetAccessToken)
	videoMonitorBsw.appendCombobox(atvsp.kitGroupsGetByIDAccessToken)
	videoMonitorBsw.appendCombobox(atvsp.kitMessagesSendAccessToken)
	partedBsw.appendBoxToCenterPart(videoMonitorBsw)

	var photoCommentMonitorSpfdb settingsParamsFromDB
	photoCommentMonitorSpfdb.setSubjectParams(subjectParams)
	photoCommentMonitorSpfdb.setAccessTokens(accessTokens)
	photoCommentMonitorSpfdb.selectMonitorParams("photo_comment_monitor")
	photoCommentMonitorSpfdb.selectPhotosGetAllCommentsMethodParams()
	photoCommentMonitorSpfdb.selectUsersGetMethodParams()
	photoCommentMonitorSpfdb.selectGroupsGetByIDMethodParams()
	photoCommentMonitorSpfdb.selectMessagesSendMethodParams()
	var atpcsp accessTokensPhotoCommentSettingsParams
	atpcsp.init(photoCommentMonitorSpfdb)
	var photoCommentMonitorBsw boxSettingsWnd
	photoCommentMonitorBsw.init()
	photoCommentMonitorBsw.initGroup("Комментарии под фото")
	photoCommentMonitorBsw.appendCombobox(atpcsp.kitPhotosGetAllCommentsAccessToken)
	photoCommentMonitorBsw.appendCombobox(atpcsp.kitUsersGetAccessToken)
	photoCommentMonitorBsw.appendCombobox(atpcsp.kitGroupsGetByIDAccessToken)
	photoCommentMonitorBsw.appendCombobox(atpcsp.kitMessagesSendAccessToken)
	partedBsw.appendBoxToCenterPart(photoCommentMonitorBsw)

	var videoCommentMonitorSpfdb settingsParamsFromDB
	videoCommentMonitorSpfdb.setSubjectParams(subjectParams)
	videoCommentMonitorSpfdb.setAccessTokens(accessTokens)
	videoCommentMonitorSpfdb.selectMonitorParams("video_comment_monitor")
	videoCommentMonitorSpfdb.selectVideoGetCommentsMethodParams()
	videoCommentMonitorSpfdb.selectUsersGetMethodParams()
	videoCommentMonitorSpfdb.selectGroupsGetByIDMethodParams()
	videoCommentMonitorSpfdb.selectVideoGetMethodParams()
	videoCommentMonitorSpfdb.selectMessagesSendMethodParams()
	var atvcsp accessTokensVideoCommentSettingsParams
	atvcsp.init(videoCommentMonitorSpfdb)
	var videoCommentMonitorBsw boxSettingsWnd
	videoCommentMonitorBsw.init()
	videoCommentMonitorBsw.initGroup("Комментарии под видео")
	videoCommentMonitorBsw.appendCombobox(atvcsp.kitVideoGetCommentsAccessToken)
	videoCommentMonitorBsw.appendCombobox(atvcsp.kitUsersGetAccessToken)
	videoCommentMonitorBsw.appendCombobox(atvcsp.kitGroupsGetByIDAccessToken)
	videoCommentMonitorBsw.appendCombobox(atvcsp.kitVideoGetAccessToken)
	videoCommentMonitorBsw.appendCombobox(atvcsp.kitMessagesSendAccessToken)
	partedBsw.appendBoxToCenterPart(videoCommentMonitorBsw)

	var topicMonitorSpfdb settingsParamsFromDB
	topicMonitorSpfdb.setSubjectParams(subjectParams)
	topicMonitorSpfdb.setAccessTokens(accessTokens)
	topicMonitorSpfdb.selectMonitorParams("topic_monitor")
	topicMonitorSpfdb.selectBoardGetCommentsMethodParams()
	topicMonitorSpfdb.selectBoardGetTopicsMethodParams()
	topicMonitorSpfdb.selectUsersGetMethodParams()
	topicMonitorSpfdb.selectGroupsGetByIDMethodParams()
	topicMonitorSpfdb.selectMessagesSendMethodParams()
	var attsp accessTokensTopicSettingsParams
	attsp.init(topicMonitorSpfdb)
	var topicMonitorBsw boxSettingsWnd
	topicMonitorBsw.init()
	topicMonitorBsw.initGroup("Комментарии в обсуждениях")
	topicMonitorBsw.appendCombobox(attsp.kitBoardGetCommentsAccessToken)
	topicMonitorBsw.appendCombobox(attsp.kitBoardGetTopicsAccessToken)
	topicMonitorBsw.appendCombobox(attsp.kitUsersGetAccessToken)
	topicMonitorBsw.appendCombobox(attsp.kitGroupsGetByIDAccessToken)
	topicMonitorBsw.appendCombobox(attsp.kitMessagesSendAccessToken)
	partedBsw.appendBoxToRightPart(topicMonitorBsw)

	var wallPostCommentMonitorSpfdb settingsParamsFromDB
	wallPostCommentMonitorSpfdb.setSubjectParams(subjectParams)
	wallPostCommentMonitorSpfdb.setAccessTokens(accessTokens)
	wallPostCommentMonitorSpfdb.selectMonitorParams("wall_post_comment_monitor")
	wallPostCommentMonitorSpfdb.selectWallGetCommentsMethodParams()
	wallPostCommentMonitorSpfdb.selectUsersGetMethodParams()
	wallPostCommentMonitorSpfdb.selectGroupsGetByIDMethodParams()
	wallPostCommentMonitorSpfdb.selectWallGetMethodParams()
	wallPostCommentMonitorSpfdb.selectWallGetCommentMethodParams()
	wallPostCommentMonitorSpfdb.selectMessagesSendMethodParams()
	var atwpcsp accessTokensWallPostCommentSettingsParams
	atwpcsp.init(wallPostCommentMonitorSpfdb)
	var wallPostCommentBsw boxSettingsWnd
	wallPostCommentBsw.init()
	wallPostCommentBsw.initGroup("Комментарии под постами")
	wallPostCommentBsw.appendCombobox(atwpcsp.kitWallGetCommentsAccessToken)
	wallPostCommentBsw.appendCombobox(atwpcsp.kitUsersGetAccessToken)
	wallPostCommentBsw.appendCombobox(atwpcsp.kitGroupsGetByIDAccessToken)
	wallPostCommentBsw.appendCombobox(atwpcsp.kitWallGetAccessToken)
	wallPostCommentBsw.appendCombobox(atwpcsp.kitWallGetCommentAccessToken)
	wallPostCommentBsw.appendCombobox(atwpcsp.kitMessagesSendAccessToken)
	partedBsw.appendBoxToRightPart(wallPostCommentBsw)

	var bsw boxSettingsWnd
	bsw.init()
	bsw.appendPartedBox(partedBsw)
	kitButtons := MakeSettingButtonsKit()
	kitButtons.ButtonApply.OnClicked(func(*ui.Button) {
		updated := sgsp.updateParamsInDB(subjectParams)
		if updated {
			atwpsp.updateParamsInDB(wallPostMonitorSpfdb)
			atapsp.updateParamsInDB(albumPhotoMonitorSpfdb)
			atvsp.updateParamsInDB(videoMonitorSpfdb)
			atpcsp.updateParamsInDB(photoCommentMonitorSpfdb)
			atvcsp.updateParamsInDB(videoCommentMonitorSpfdb)
			attsp.updateParamsInDB(topicMonitorSpfdb)
			atwpcsp.updateParamsInDB(wallPostCommentMonitorSpfdb)
			sw.window.Disable()
			sw.window.Hide()
		}
	})
	kitButtons.ButtonCancel.OnClicked(func(*ui.Button) {
		sw.window.Disable()
		sw.window.Hide()
	})
	bsw.appendButtons(kitButtons)
	bsw.initFlexibleSpaceBox()

	sw.setBox(bsw)
}

type wallPostSettingsParams struct {
	kitNeedMonitoring        CheckboxKit
	kitInterval              SpinboxKit
	kitSendTo                EntryKit
	kitFilter                ComboboxKit
	kitPostsCount            SpinboxKit
	kitKeywordsForMonitoring EntryListKit
}

func (wpsp *wallPostSettingsParams) init(params WallPostMonitorParam) {
	wpsp.kitNeedMonitoring = MakeSettingCheckboxKit("Разрешить наблюдение", params.NeedMonitoring)
	wpsp.kitInterval = MakeSettingSpinboxKit("Интервал (сек.)", 5, 21600, params.Interval)
	wpsp.kitSendTo = MakeSettingEntryKit("Получатель", strconv.Itoa(params.SendTo))
	wpsp.kitSendTo.Entry.OnChanged(func(*ui.Entry) {
		NumericEntriesHandler(wpsp.kitSendTo.Entry)
	})
	listPostsFilters := []string{"Все", "От пользователей", "От сообщества", "Предложенные"}
	listPostsFiltersEn := []string{"all", "others", "owner", "suggests"}
	var currentFilter string
	for i, item := range listPostsFiltersEn {
		if params.Filter == item {
			currentFilter = listPostsFilters[i]
			break
		}
	}
	wpsp.kitFilter = MakeSettingComboboxKit("Фильтр", listPostsFilters, currentFilter)
	wpsp.kitPostsCount = MakeSettingSpinboxKit("Количество постов", 1, 100, params.PostsCount)
	wpsp.kitKeywordsForMonitoring = MakeSettingEntryListKit("Ключевые слова", params.KeywordsForMonitoring)
}

func (wpsp *wallPostSettingsParams) updateParamsInDB(params WallPostMonitorParam) bool {
	var updatedParams WallPostMonitorParam
	updatedParams.ID = params.ID
	updatedParams.SubjectID = params.SubjectID
	if wpsp.kitNeedMonitoring.CheckBox.Checked() {
		updatedParams.NeedMonitoring = 1
	} else {
		updatedParams.NeedMonitoring = 0
	}
	updatedParams.Interval = wpsp.kitInterval.Spinbox.Value()
	if len(wpsp.kitSendTo.Entry.Text()) == 0 {
		warningTitle := "Поле \"Получатель\" не должно быть пустым."
		ShowWarningWindow(warningTitle)
		return false
	}
	var err error
	updatedParams.SendTo, err = strconv.Atoi(wpsp.kitSendTo.Entry.Text())
	if err != nil {
		ToLogFile(err.Error(), string(debug.Stack()))
		panic(err.Error())
	}
	listPostsFilters := []string{"all", "others", "owner", "suggests"}
	if wpsp.kitFilter.Combobox.Selected() == -1 {
		warningTitle := "Нужно выбрать элемент в списке \"Фильтр\""
		ShowWarningWindow(warningTitle)
		return false
	}
	updatedParams.Filter = listPostsFilters[wpsp.kitFilter.Combobox.Selected()]
	updatedParams.LastDate = params.LastDate
	updatedParams.PostsCount = wpsp.kitPostsCount.Spinbox.Value()
	// TODO: проверка соответствия оформления требованиям json
	jsonDump := fmt.Sprintf("{\"list\":[%v]}", wpsp.kitKeywordsForMonitoring.Entry.Text())
	updatedParams.KeywordsForMonitoring = jsonDump
	updatedParams.UsersIDsForIgnore = params.UsersIDsForIgnore

	err = updatedParams.updateInDB()
	if err != nil {
		ToLogFile(err.Error(), string(debug.Stack()))
		panic(err.Error())
	}

	return true
}

func showSubjectWallPostSettingWindow(IDSubject int, nameSubject, btnName string) {
	var params WallPostMonitorParam
	err := params.selectFromDBBySubjectID(IDSubject)
	if err != nil {
		ToLogFile(err.Error(), string(debug.Stack()))
		panic(err.Error())
	}

	var wpsp wallPostSettingsParams
	wpsp.init(params)

	windowTitle := fmt.Sprintf("%v: настройки для %v", btnName, nameSubject)
	var sw settingsWindow
	sw.init(windowTitle, 300, 100)

	var bsw boxSettingsWnd
	bsw.init()
	bsw.initGroup("")
	bsw.appendCheckbox(wpsp.kitNeedMonitoring)
	bsw.appendSpinbox(wpsp.kitInterval)
	bsw.appendEntry(wpsp.kitSendTo)
	bsw.appendCombobox(wpsp.kitFilter)
	bsw.appendSpinbox(wpsp.kitPostsCount)
	bsw.appendEntryListKit(wpsp.kitKeywordsForMonitoring)
	kitButtons := MakeSettingButtonsKit()
	kitButtons.ButtonApply.OnClicked(func(*ui.Button) {
		updated := wpsp.updateParamsInDB(params)
		if updated {
			sw.window.Disable()
			sw.window.Hide()
		}
	})
	kitButtons.ButtonCancel.OnClicked(func(*ui.Button) {
		sw.window.Disable()
		sw.window.Hide()
	})
	bsw.appendButtons(kitButtons)
	bsw.initFlexibleSpaceBox()

	sw.setBox(bsw)
}

type albumPhotoSettingsParams struct {
	kitNeedMonitoring CheckboxKit
	kitInterval       SpinboxKit
	kitSendTo         EntryKit
	kitPhotoCount     SpinboxKit
}

func (apsp *albumPhotoSettingsParams) init(params AlbumPhotoMonitorParam) {
	apsp.kitNeedMonitoring = MakeSettingCheckboxKit("Разрешить наблюдение", params.NeedMonitoring)
	apsp.kitInterval = MakeSettingSpinboxKit("Интервал (сек.)", 5, 21600, params.Interval)
	apsp.kitSendTo = MakeSettingEntryKit("Получатель", strconv.Itoa(params.SendTo))
	apsp.kitSendTo.Entry.OnChanged(func(*ui.Entry) {
		NumericEntriesHandler(apsp.kitSendTo.Entry)
	})
	apsp.kitPhotoCount = MakeSettingSpinboxKit("Количество фотографий", 1, 1000, params.PhotosCount)
}

func (apsp *albumPhotoSettingsParams) updateParamsInDB(params AlbumPhotoMonitorParam) bool {
	var updatedParams AlbumPhotoMonitorParam
	updatedParams.ID = params.ID
	updatedParams.SubjectID = params.SubjectID
	if apsp.kitNeedMonitoring.CheckBox.Checked() {
		updatedParams.NeedMonitoring = 1
	} else {
		updatedParams.NeedMonitoring = 0
	}
	if len(apsp.kitSendTo.Entry.Text()) == 0 {
		warningTitle := "Поле \"Получатель\" не должно быть пустым."
		ShowWarningWindow(warningTitle)
		return false
	}
	var err error
	updatedParams.SendTo, err = strconv.Atoi(apsp.kitSendTo.Entry.Text())
	if err != nil {
		ToLogFile(err.Error(), string(debug.Stack()))
		panic(err.Error())
	}
	updatedParams.Interval = apsp.kitInterval.Spinbox.Value()
	updatedParams.LastDate = params.LastDate
	updatedParams.PhotosCount = apsp.kitPhotoCount.Spinbox.Value()

	err = updatedParams.updateInDB()
	if err != nil {
		ToLogFile(err.Error(), string(debug.Stack()))
		panic(err.Error())
	}

	return true
}

func showSubjectAlbumPhotoSettingWindow(IDSubject int, nameSubject, btnName string) {
	var params AlbumPhotoMonitorParam
	err := params.selectFromDBBySubjectID(IDSubject)
	if err != nil {
		ToLogFile(err.Error(), string(debug.Stack()))
		panic(err.Error())
	}

	var apsp albumPhotoSettingsParams
	apsp.init(params)

	windowTitle := fmt.Sprintf("%v: настройки для %v", btnName, nameSubject)
	var sw settingsWindow
	sw.init(windowTitle, 300, 100)

	var bsw boxSettingsWnd
	bsw.init()
	bsw.initGroup("")
	bsw.appendCheckbox(apsp.kitNeedMonitoring)
	bsw.appendSpinbox(apsp.kitInterval)
	bsw.appendEntry(apsp.kitSendTo)
	bsw.appendSpinbox(apsp.kitPhotoCount)
	kitButtons := MakeSettingButtonsKit()
	kitButtons.ButtonApply.OnClicked(func(*ui.Button) {
		updated := apsp.updateParamsInDB(params)
		if updated {
			sw.window.Disable()
			sw.window.Hide()
		}
	})
	kitButtons.ButtonCancel.OnClicked(func(*ui.Button) {
		sw.window.Disable()
		sw.window.Hide()
	})
	bsw.appendButtons(kitButtons)
	bsw.initFlexibleSpaceBox()

	sw.setBox(bsw)
}

type videoSettingsParams struct {
	kitNeedMonitoring CheckboxKit
	kitInterval       SpinboxKit
	kitSendTo         EntryKit
	kitVideoCount     SpinboxKit
}

func (vsp *videoSettingsParams) init(params VideoMonitorParam) {
	vsp.kitNeedMonitoring = MakeSettingCheckboxKit("Разрешить наблюдение", params.NeedMonitoring)
	vsp.kitInterval = MakeSettingSpinboxKit("Интервал (сек.)", 5, 21600, params.Interval)
	vsp.kitSendTo = MakeSettingEntryKit("Получатель", strconv.Itoa(params.SendTo))
	vsp.kitSendTo.Entry.OnChanged(func(*ui.Entry) {
		NumericEntriesHandler(vsp.kitSendTo.Entry)
	})
	vsp.kitVideoCount = MakeSettingSpinboxKit("Количество видеозаписей", 1, 1000, params.VideoCount)
}

func (vsp *videoSettingsParams) updateParamsInDB(params VideoMonitorParam) bool {
	var updatedParams VideoMonitorParam
	updatedParams.ID = params.ID
	updatedParams.SubjectID = params.SubjectID
	if vsp.kitNeedMonitoring.CheckBox.Checked() {
		updatedParams.NeedMonitoring = 1
	} else {
		updatedParams.NeedMonitoring = 0
	}
	if len(vsp.kitSendTo.Entry.Text()) == 0 {
		warningTitle := "Поле \"Получатель\" не должно быть пустым."
		ShowWarningWindow(warningTitle)
		return false
	}
	var err error
	updatedParams.SendTo, err = strconv.Atoi(vsp.kitSendTo.Entry.Text())
	if err != nil {
		ToLogFile(err.Error(), string(debug.Stack()))
		panic(err.Error())
	}
	updatedParams.Interval = vsp.kitInterval.Spinbox.Value()
	updatedParams.LastDate = params.LastDate
	updatedParams.VideoCount = vsp.kitInterval.Spinbox.Value()

	err = updatedParams.updateInDB()
	if err != nil {
		ToLogFile(err.Error(), string(debug.Stack()))
		panic(err.Error())
	}

	return true
}

func showSubjectVideoSettingWindow(IDSubject int, nameSubject, btnName string) {
	var params VideoMonitorParam
	err := params.selectFromDBBySubjectID(IDSubject)
	if err != nil {
		ToLogFile(err.Error(), string(debug.Stack()))
		panic(err.Error())
	}

	var vsp videoSettingsParams
	vsp.init(params)

	windowTitle := fmt.Sprintf("%v: настройки для %v", btnName, nameSubject)
	var sw settingsWindow
	sw.init(windowTitle, 300, 100)

	var bsw boxSettingsWnd
	bsw.init()
	bsw.initGroup("")
	bsw.appendCheckbox(vsp.kitNeedMonitoring)
	bsw.appendSpinbox(vsp.kitInterval)
	bsw.appendEntry(vsp.kitSendTo)
	bsw.appendSpinbox(vsp.kitVideoCount)
	kitButtons := MakeSettingButtonsKit()
	kitButtons.ButtonApply.OnClicked(func(*ui.Button) {
		updated := vsp.updateParamsInDB(params)
		if updated {
			sw.window.Disable()
			sw.window.Hide()
		}
	})
	kitButtons.ButtonCancel.OnClicked(func(*ui.Button) {
		sw.window.Disable()
		sw.window.Hide()
	})
	bsw.appendButtons(kitButtons)
	bsw.initFlexibleSpaceBox()

	sw.setBox(bsw)
}

type photoCommentSettingsParams struct {
	kitNeedMonitoring CheckboxKit
	kitInterval       SpinboxKit
	kitSendTo         EntryKit
	kitCommentsCount  SpinboxKit
}

func (pcsp *photoCommentSettingsParams) init(params PhotoCommentMonitorParam) {
	pcsp.kitNeedMonitoring = MakeSettingCheckboxKit("Разрешить наблюдение", params.NeedMonitoring)
	pcsp.kitInterval = MakeSettingSpinboxKit("Интервал (сек.)", 5, 21600, params.Interval)
	pcsp.kitSendTo = MakeSettingEntryKit("Получатель", strconv.Itoa(params.SendTo))
	pcsp.kitSendTo.Entry.OnChanged(func(*ui.Entry) {
		NumericEntriesHandler(pcsp.kitSendTo.Entry)
	})
	pcsp.kitCommentsCount = MakeSettingSpinboxKit("Количество комментариев", 1, 1000, params.CommentsCount)
}

func (pcsp *photoCommentSettingsParams) updateParamsInDB(params PhotoCommentMonitorParam) bool {
	var updatedParams PhotoCommentMonitorParam
	updatedParams.ID = params.ID
	updatedParams.SubjectID = params.SubjectID
	if pcsp.kitNeedMonitoring.CheckBox.Checked() {
		updatedParams.NeedMonitoring = 1
	} else {
		updatedParams.NeedMonitoring = 0
	}
	if len(pcsp.kitSendTo.Entry.Text()) == 0 {
		warningTitle := "Поле \"Получатель\" не должно быть пустым."
		ShowWarningWindow(warningTitle)
		return false
	}
	var err error
	updatedParams.SendTo, err = strconv.Atoi(pcsp.kitSendTo.Entry.Text())
	if err != nil {
		ToLogFile(err.Error(), string(debug.Stack()))
		panic(err.Error())
	}
	updatedParams.Interval = pcsp.kitInterval.Spinbox.Value()
	updatedParams.LastDate = params.LastDate
	updatedParams.CommentsCount = pcsp.kitCommentsCount.Spinbox.Value()

	err = updatedParams.updateInDB()
	if err != nil {
		ToLogFile(err.Error(), string(debug.Stack()))
		panic(err.Error())
	}

	return true
}

func showSubjectPhotoCommentSettingWindow(IDSubject int, nameSubject, btnName string) {
	var params PhotoCommentMonitorParam
	err := params.selectFromDBBySubjectID(IDSubject)
	if err != nil {
		ToLogFile(err.Error(), string(debug.Stack()))
		panic(err.Error())
	}

	var pcsp photoCommentSettingsParams
	pcsp.init(params)

	windowTitle := fmt.Sprintf("%v: настройки для %v", btnName, nameSubject)
	var sw settingsWindow
	sw.init(windowTitle, 300, 100)

	var bsw boxSettingsWnd
	bsw.init()
	bsw.initGroup("")
	bsw.appendCheckbox(pcsp.kitNeedMonitoring)
	bsw.appendSpinbox(pcsp.kitInterval)
	bsw.appendEntry(pcsp.kitSendTo)
	bsw.appendSpinbox(pcsp.kitCommentsCount)
	kitButtons := MakeSettingButtonsKit()
	kitButtons.ButtonApply.OnClicked(func(*ui.Button) {
		updated := pcsp.updateParamsInDB(params)
		if updated {
			sw.window.Disable()
			sw.window.Hide()
		}
	})
	kitButtons.ButtonCancel.OnClicked(func(*ui.Button) {
		sw.window.Disable()
		sw.window.Hide()
	})
	bsw.appendButtons(kitButtons)
	bsw.initFlexibleSpaceBox()

	sw.setBox(bsw)
}

type videoCommentSettingsParams struct {
	kitNeedMonitoring CheckboxKit
	kitInterval       SpinboxKit
	kitSendTo         EntryKit
	kitVideoCount     SpinboxKit
	kitCommentsCount  SpinboxKit
}

func (vcsp *videoCommentSettingsParams) init(params VideoCommentMonitorParam) {
	vcsp.kitNeedMonitoring = MakeSettingCheckboxKit("Разрешить наблюдение", params.NeedMonitoring)
	vcsp.kitInterval = MakeSettingSpinboxKit("Интервал (сек.)", 5, 21600, params.Interval)
	vcsp.kitSendTo = MakeSettingEntryKit("Получатель", strconv.Itoa(params.SendTo))
	vcsp.kitSendTo.Entry.OnChanged(func(*ui.Entry) {
		NumericEntriesHandler(vcsp.kitSendTo.Entry)
	})
	vcsp.kitVideoCount = MakeSettingSpinboxKit("Количество видеозаписей", 1, 200, params.VideosCount)
	vcsp.kitCommentsCount = MakeSettingSpinboxKit("Количество комментариев", 1, 100, params.CommentsCount)
}

func (vcsp *videoCommentSettingsParams) updateParamsInDB(params VideoCommentMonitorParam) bool {
	var updatedParams VideoCommentMonitorParam
	updatedParams.ID = params.ID
	updatedParams.SubjectID = params.SubjectID
	if vcsp.kitNeedMonitoring.CheckBox.Checked() {
		updatedParams.NeedMonitoring = 1
	} else {
		updatedParams.NeedMonitoring = 0
	}
	if len(vcsp.kitSendTo.Entry.Text()) == 0 {
		warningTitle := "Поле \"Получатель\" не должно быть пустым."
		ShowWarningWindow(warningTitle)
		return false
	}
	var err error
	updatedParams.SendTo, err = strconv.Atoi(vcsp.kitSendTo.Entry.Text())
	if err != nil {
		// FIXME: тут старый способ обработки ошибок
		date := UnixTimeStampToDate(int(time.Now().Unix()))
		log.Fatal(fmt.Errorf("> [%v] WARNING! Error: %v", date, err))
	}
	updatedParams.Interval = vcsp.kitInterval.Spinbox.Value()
	updatedParams.LastDate = params.LastDate
	updatedParams.CommentsCount = vcsp.kitCommentsCount.Spinbox.Value()
	updatedParams.VideosCount = vcsp.kitVideoCount.Spinbox.Value()

	err = updatedParams.updateInDB()
	if err != nil {
		ToLogFile(err.Error(), string(debug.Stack()))
		panic(err.Error())
	}

	return true
}

func showSubjectVideoCommentSettingWindow(IDSubject int, nameSubject, btnName string) {
	var params VideoCommentMonitorParam
	err := params.selectFromDBBySubjectID(IDSubject)
	if err != nil {
		ToLogFile(err.Error(), string(debug.Stack()))
		panic(err.Error())
	}

	var vcsp videoCommentSettingsParams
	vcsp.init(params)

	windowTitle := fmt.Sprintf("%v: настройки для %v", btnName, nameSubject)
	var sw settingsWindow
	sw.init(windowTitle, 300, 100)

	var bsw boxSettingsWnd
	bsw.init()
	bsw.initGroup("")
	bsw.appendCheckbox(vcsp.kitNeedMonitoring)
	bsw.appendSpinbox(vcsp.kitInterval)
	bsw.appendEntry(vcsp.kitSendTo)
	bsw.appendSpinbox(vcsp.kitVideoCount)
	bsw.appendSpinbox(vcsp.kitCommentsCount)
	kitButtons := MakeSettingButtonsKit()
	kitButtons.ButtonApply.OnClicked(func(*ui.Button) {
		updated := vcsp.updateParamsInDB(params)
		if updated {
			sw.window.Disable()
			sw.window.Hide()
		}
	})
	kitButtons.ButtonCancel.OnClicked(func(*ui.Button) {
		sw.window.Disable()
		sw.window.Hide()
	})
	bsw.appendButtons(kitButtons)
	bsw.initFlexibleSpaceBox()

	sw.setBox(bsw)
}

type topicSettingsParams struct {
	kitNeedMonitoring CheckboxKit
	kitInterval       SpinboxKit
	kitSendTo         EntryKit
	kitTopicsCount    SpinboxKit
	kitCommentsCount  SpinboxKit
}

func (tsp *topicSettingsParams) init(params TopicMonitorParam) {
	tsp.kitNeedMonitoring = MakeSettingCheckboxKit("Разрешить наблюдение", params.NeedMonitoring)
	tsp.kitInterval = MakeSettingSpinboxKit("Интервал (сек.)", 5, 21600, params.Interval)
	tsp.kitSendTo = MakeSettingEntryKit("Получатель", strconv.Itoa(params.SendTo))
	tsp.kitSendTo.Entry.OnChanged(func(*ui.Entry) {
		NumericEntriesHandler(tsp.kitSendTo.Entry)
	})
	tsp.kitTopicsCount = MakeSettingSpinboxKit("Количество обсуждений", 1, 100, params.TopicsCount)
	tsp.kitCommentsCount = MakeSettingSpinboxKit("Количество комментариев", 1, 100, params.TopicsCount) // FIXME: поменять params.TopicsCount на params.CommentsCount
}

func (tsp *topicSettingsParams) updateParamsInDB(params TopicMonitorParam) bool {
	var updatedParams TopicMonitorParam
	updatedParams.ID = params.ID
	updatedParams.SubjectID = params.SubjectID
	if tsp.kitNeedMonitoring.CheckBox.Checked() {
		updatedParams.NeedMonitoring = 1
	} else {
		updatedParams.NeedMonitoring = 0
	}
	if len(tsp.kitSendTo.Entry.Text()) == 0 {
		warningTitle := "Поле \"Получатель\" не должно быть пустым."
		ShowWarningWindow(warningTitle)
		return false
	}
	var err error
	updatedParams.SendTo, err = strconv.Atoi(tsp.kitSendTo.Entry.Text())
	if err != nil {
		// FIXME: тут старый способ обработки ошибок
		date := UnixTimeStampToDate(int(time.Now().Unix()))
		log.Fatal(fmt.Errorf("> [%v] WARNING! Error: %v", date, err))
	}
	updatedParams.Interval = tsp.kitInterval.Spinbox.Value()
	updatedParams.LastDate = params.LastDate
	updatedParams.CommentsCount = tsp.kitCommentsCount.Spinbox.Value()
	updatedParams.TopicsCount = tsp.kitTopicsCount.Spinbox.Value()

	err = updatedParams.updateInDB()
	if err != nil {
		ToLogFile(err.Error(), string(debug.Stack()))
		panic(err.Error())
	}

	return true
}

func showSubjectTopicSettingWindow(IDSubject int, nameSubject, btnName string) {
	var params TopicMonitorParam
	err := params.selectFromDBBySubjectID(IDSubject)
	if err != nil {
		ToLogFile(err.Error(), string(debug.Stack()))
		panic(err.Error())
	}

	var tsp topicSettingsParams
	tsp.init(params)

	windowTitle := fmt.Sprintf("%v: настройки для %v", btnName, nameSubject)
	var sw settingsWindow
	sw.init(windowTitle, 300, 100)

	var bsw boxSettingsWnd
	bsw.init()
	bsw.initGroup("")
	bsw.appendCheckbox(tsp.kitNeedMonitoring)
	bsw.appendSpinbox(tsp.kitInterval)
	bsw.appendEntry(tsp.kitSendTo)
	bsw.appendSpinbox(tsp.kitTopicsCount)
	bsw.appendSpinbox(tsp.kitCommentsCount)
	kitButtons := MakeSettingButtonsKit()
	kitButtons.ButtonApply.OnClicked(func(*ui.Button) {
		updated := tsp.updateParamsInDB(params)
		if updated {
			sw.window.Disable()
			sw.window.Hide()
		}
	})
	kitButtons.ButtonCancel.OnClicked(func(*ui.Button) {
		sw.window.Disable()
		sw.window.Hide()
	})
	bsw.appendButtons(kitButtons)
	bsw.initFlexibleSpaceBox()

	sw.setBox(bsw)
}

type wallPostCommentSettingsParams struct {
	kitNeedMonitoring             CheckboxKit
	kitInterval                   SpinboxKit
	kitSendTo                     EntryKit
	kitPostsCount                 SpinboxKit
	kitCommentsCount              SpinboxKit
	kitFilter                     ComboboxKit
	kitMonitoringAll              CheckboxKit
	kitMonitorByCommunity         CheckboxKit
	kitKeywordsForMonitoring      EntryListKit
	kitSmallCommentsForMonitoring EntryListKit
	kitUsersNamesForMonitoring    EntryListKit
	kitUsersIdsForMonitoring      EntryListKit
	kitUsersIdsForIgnore          EntryListKit
}

func (wpcsp *wallPostCommentSettingsParams) init(params WallPostCommentMonitorParam) {
	wpcsp.kitNeedMonitoring = MakeSettingCheckboxKit("Разрешить наблюдение", params.NeedMonitoring)
	wpcsp.kitInterval = MakeSettingSpinboxKit("Интервал (сек.)", 5, 21600, params.Interval)
	wpcsp.kitSendTo = MakeSettingEntryKit("Получатель", strconv.Itoa(params.SendTo))
	wpcsp.kitSendTo.Entry.OnChanged(func(*ui.Entry) {
		NumericEntriesHandler(wpcsp.kitSendTo.Entry)
	})
	wpcsp.kitPostsCount = MakeSettingSpinboxKit("Количество постов", 1, 100, params.PostsCount)
	wpcsp.kitCommentsCount = MakeSettingSpinboxKit("Количество комментариев", 1, 100, params.CommentsCount)
	listPostsFilters := []string{"Все", "От пользователей", "От сообщества"}
	listPostsFiltersEn := []string{"all", "others", "owner"}
	var currentFilter string
	for i, item := range listPostsFiltersEn {
		if params.Filter == item {
			currentFilter = listPostsFilters[i]
			break
		}
	}
	wpcsp.kitFilter = MakeSettingComboboxKit("Фильтр", listPostsFilters, currentFilter)
	wpcsp.kitMonitoringAll = MakeSettingCheckboxKit("Собирать все комментарии", params.MonitoringAll)
	wpcsp.kitMonitorByCommunity = MakeSettingCheckboxKit("Собирать комментарии от сообществ", params.MonitorByCommunity)
	wpcsp.kitKeywordsForMonitoring = MakeSettingEntryListKit("Ключевые слова", params.KeywordsForMonitoring)
	wpcsp.kitSmallCommentsForMonitoring = MakeSettingEntryListKit("Комментарии", params.SmallCommentsForMonitoring)
	wpcsp.kitUsersNamesForMonitoring = MakeSettingEntryListKit("Наблюдать комментаторов по имени", params.UsersNamesForMonitoring)
	wpcsp.kitUsersIdsForMonitoring = MakeSettingEntryListKit("Наблюдать комментаторов по ид-у в ВК", params.UsersIDsForMonitoring)
	wpcsp.kitUsersIdsForIgnore = MakeSettingEntryListKit("Игнорировать комментаторов по ид-у в ВК", params.UsersIDsForIgnore)
}

func (wpcsp *wallPostCommentSettingsParams) updateParamsInDB(params WallPostCommentMonitorParam) bool {
	var updatedParams WallPostCommentMonitorParam
	updatedParams.ID = params.ID
	updatedParams.SubjectID = params.SubjectID
	if wpcsp.kitNeedMonitoring.CheckBox.Checked() {
		updatedParams.NeedMonitoring = 1
	} else {
		updatedParams.NeedMonitoring = 0
	}
	updatedParams.PostsCount = wpcsp.kitPostsCount.Spinbox.Value()
	updatedParams.CommentsCount = wpcsp.kitCommentsCount.Spinbox.Value()
	if wpcsp.kitMonitoringAll.CheckBox.Checked() {
		updatedParams.MonitoringAll = 1
	} else {
		updatedParams.MonitoringAll = 0
	}
	jsonDump := fmt.Sprintf("{\"list\":[%v]}", wpcsp.kitUsersIdsForMonitoring.Entry.Text())
	updatedParams.UsersIDsForMonitoring = jsonDump
	jsonDump = fmt.Sprintf("{\"list\":[%v]}", wpcsp.kitUsersNamesForMonitoring.Entry.Text())
	updatedParams.UsersNamesForMonitoring = jsonDump
	updatedParams.AttachmentsTypesForMonitoring = params.AttachmentsTypesForMonitoring
	jsonDump = fmt.Sprintf("{\"list\":[%v]}", wpcsp.kitUsersIdsForIgnore.Entry.Text())
	updatedParams.UsersIDsForIgnore = jsonDump
	updatedParams.Interval = wpcsp.kitInterval.Spinbox.Value()
	if len(wpcsp.kitSendTo.Entry.Text()) == 0 {
		warningTitle := "Поле \"Получатель\" не должно быть пустым."
		ShowWarningWindow(warningTitle)
		return false
	}
	var err error
	updatedParams.SendTo, err = strconv.Atoi(wpcsp.kitSendTo.Entry.Text())
	if err != nil {
		ToLogFile(err.Error(), string(debug.Stack()))
		panic(err.Error())
	}
	listPostsFilters := []string{"all", "others", "owner"}
	if wpcsp.kitFilter.Combobox.Selected() == -1 {
		warningTitle := "Нужно выбрать элемент в списке \"Фильтр\""
		ShowWarningWindow(warningTitle)
		return false
	}
	updatedParams.Filter = listPostsFilters[wpcsp.kitFilter.Combobox.Selected()]
	updatedParams.LastDate = params.LastDate
	// TODO: проверка соответствия оформления требованиям json
	jsonDump = fmt.Sprintf("{\"list\":[%v]}", wpcsp.kitKeywordsForMonitoring.Entry.Text())
	updatedParams.KeywordsForMonitoring = jsonDump
	// TODO: проверка соответствия оформления требованиям json
	jsonDump = fmt.Sprintf("{\"list\":[%v]}", wpcsp.kitSmallCommentsForMonitoring.Entry.Text())
	updatedParams.SmallCommentsForMonitoring = jsonDump
	updatedParams.DigitsCountForCardNumberMonitoring = params.DigitsCountForCardNumberMonitoring
	updatedParams.DigitsCountForPhoneNumberMonitoring = params.DigitsCountForPhoneNumberMonitoring
	if wpcsp.kitMonitorByCommunity.CheckBox.Checked() {
		updatedParams.MonitorByCommunity = 1
	} else {
		updatedParams.MonitorByCommunity = 0
	}

	err = updatedParams.updateInDB()
	if err != nil {
		ToLogFile(err.Error(), string(debug.Stack()))
		panic(err.Error())
	}

	return true
}

func showSubjectWallPostCommentSettings(IDSubject int, nameSubject, btnName string) {
	var params WallPostCommentMonitorParam
	err := params.selectFromDBBySubjectID(IDSubject)
	if err != nil {
		ToLogFile(err.Error(), string(debug.Stack()))
		panic(err.Error())
	}

	var wpcsp wallPostCommentSettingsParams
	wpcsp.init(params)

	windowTitle := fmt.Sprintf("%v: настройки для %v", btnName, nameSubject)
	var sw settingsWindow
	sw.init(windowTitle, 300, 100)

	var bsw boxSettingsWnd
	bsw.init()
	bsw.initGroup("")
	bsw.appendCheckbox(wpcsp.kitNeedMonitoring)
	bsw.appendSpinbox(wpcsp.kitInterval)
	bsw.appendEntry(wpcsp.kitSendTo)
	bsw.appendSpinbox(wpcsp.kitPostsCount)
	bsw.appendSpinbox(wpcsp.kitCommentsCount)
	bsw.appendCombobox(wpcsp.kitFilter)
	bsw.appendCheckbox(wpcsp.kitMonitoringAll)
	bsw.appendCheckbox(wpcsp.kitMonitorByCommunity)
	bsw.appendEntryListKit(wpcsp.kitKeywordsForMonitoring)
	bsw.appendEntryListKit(wpcsp.kitSmallCommentsForMonitoring)
	bsw.appendEntryListKit(wpcsp.kitUsersNamesForMonitoring)
	bsw.appendEntryListKit(wpcsp.kitUsersIdsForMonitoring)
	bsw.appendEntryListKit(wpcsp.kitUsersIdsForIgnore)
	kitButtons := MakeSettingButtonsKit()
	kitButtons.ButtonApply.OnClicked(func(*ui.Button) {
		updated := wpcsp.updateParamsInDB(params)
		if updated {
			sw.window.Disable()
			sw.window.Hide()
		}
	})
	kitButtons.ButtonCancel.OnClicked(func(*ui.Button) {
		sw.window.Disable()
		sw.window.Hide()
	})
	bsw.appendButtons(kitButtons)
	bsw.initFlexibleSpaceBox()

	sw.setBox(bsw)
}
