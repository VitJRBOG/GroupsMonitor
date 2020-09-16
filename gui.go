package main

import (
	"fmt"
	"regexp"
	"runtime/debug"
	"strings"

	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"
)

// RunGui запускает собранный GUI
func RunGui() error {
	err := ui.Main(initGui)
	if err != nil {
		return err
	}
	return nil
}

func initGui() {
	// проверяем наличие ресурсных файлов программы, и если их нет, то создаем
	initFiles()

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
		ToLogFile(err.Error(), string(debug.Stack()))
		panic(err.Error())
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
	slctd := -1
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

func showWarningWindow(warningTitle string) {
	// описываем окно с информацией об ошибке
	windowWarning := ui.NewWindow("WARNING!", 100, 100, true)
	windowWarning.SetMargined(true)
	windowWarning.OnClosing(func(*ui.Window) bool {
		windowWarning.Disable()
		return true
	})

	// описываем основную коробку
	boxWndWarning := ui.NewVerticalBox()
	boxWndWarning.SetPadded(true)
	windowWarning.SetChild(boxWndWarning)

	// описываем коробку для информации
	boxInfo := ui.NewVerticalBox()

	// описываем заголовок ошибки
	labelTitleWarning := ui.NewLabel(warningTitle)
	boxInfo.Append(labelTitleWarning, true)

	// описываем коробку для кнопки
	boxButton := ui.NewHorizontalBox()

	// описываем кнопку
	buttonOK := ui.NewButton("OK")
	buttonOK.OnClicked(func(*ui.Button) {
		windowWarning.Hide()
	})

	// описываем коробку для выравнивания кнопки и коробку с кнопкой
	boxLeftPartButtonBox := ui.NewVerticalBox()
	boxRightPartButtonBox := ui.NewVerticalBox()
	boxRightPartButtonBox.Append(buttonOK, false)

	// добавляем их на коробку для кнопки
	boxButton.Append(boxLeftPartButtonBox, true)
	boxButton.Append(boxRightPartButtonBox, false)

	// добавляем все эти компоненты на главную коробку
	boxWndWarning.Append(boxInfo, false)
	boxWndWarning.Append(boxButton, false)

	windowWarning.Show()
}

func numericEntriesHandler(numericEntry *ui.Entry) {
	// проверим, есть ли знак минуса в начале строки
	negativeNumber := false
	if len(numericEntry.Text()) > 0 {
		listEntryChars := strings.Split(numericEntry.Text(), "")
		if listEntryChars[0] == "-" {
			negativeNumber = true
		}
	}

	re := regexp.MustCompile("[0-9]+")
	correctValue := strings.Join(re.FindAllString(numericEntry.Text(), -1), "")

	if negativeNumber {
		numericEntry.SetText("-" + correctValue)
	} else {
		numericEntry.SetText(correctValue)
	}
}

func initFiles() {
	dbHasBeenCreated, err := CheckFiles()
	if err != nil {
		ToLogFile(err.Error(), string(debug.Stack()))
		panic(err.Error())
	}

	if dbHasBeenCreated {
		message := "File of database has been created just now. Database is empty. " +
			"Need to create new access token and new subject for monitoring."
		showWarningWindow(message)
	}
}

func createThreads() []*Thread {
	// запускаем функцию создания потоков с модулями проверки
	threads, err := MakeThreads()
	if err != nil {
		ToLogFile(err.Error(), string(debug.Stack()))
		panic(err.Error())
	}

	return threads
}
