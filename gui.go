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
	err := ui.Main(initGUI)
	if err != nil {
		return err
	}
	return nil
}

// mainWindow хранит данные о главном окне программы
type mainWindow struct {
	window *ui.Window
}

func (mw *mainWindow) init() {
	mw.window = ui.NewWindow("GroupsMonitor", 255, 160, true)
	mw.window.SetMargined(true)
	mw.window.OnClosing(func(*ui.Window) bool {
		ui.Quit()
		return true
	})
	ui.OnShouldQuit(func() bool {
		mw.window.Destroy()
		return true
	})
	mw.window.Show()
}

func (mw *mainWindow) setBox(bmw boxMainWnd) {
	mw.window.SetChild(bmw.box)
}

// boxMainWnd хранит данные о главном боксе главного окна программы
type boxMainWnd struct {
	box *ui.Box
}

func (bmw *boxMainWnd) init() {
	bmw.box = ui.NewVerticalBox()
	bmw.box.SetPadded(true)
}

func (bmw *boxMainWnd) appendUpperPart(upperPart upperPartBoxMainWnd) {
	bmw.box.Append(upperPart.box, false)
}

func (bmw *boxMainWnd) appendBottomPart(bottomPart bottomPartBoxMainWnd) {
	bmw.box.Append(bottomPart.group, false)
}

// upperPartBoxMainWnd хранит данные о верхней части главного бокса
type upperPartBoxMainWnd struct {
	box *ui.Box
}

func (upperPart *upperPartBoxMainWnd) init() {
	upperPart.box = ui.NewHorizontalBox()
}

func (upperPart *upperPartBoxMainWnd) initFlexibleSpaceBox() {
	box := ui.NewHorizontalBox()
	upperPart.box.Append(box, true)
}

func (upperPart *upperPartBoxMainWnd) appendButtonsBox(sbb selectingButtonsBox) {
	upperPart.box.Append(sbb.box, false)
}

// selectingButtonsBox хранит данные о кнопках переключения между разделами программы
type selectingButtonsBox struct {
	box               *ui.Box
	btnGeneral        *ui.Button
	btnThreadsControl *ui.Button
	btnSettings       *ui.Button
}

func (bb *selectingButtonsBox) init() {
	bb.box = ui.NewHorizontalBox()
}

func (bb *selectingButtonsBox) initBtnGeneral(bottomPart bottomPartBoxMainWnd, boxGeneral *ui.Box) {
	bb.btnGeneral = ui.NewButton("General")
	bb.btnGeneral.Disable()

	bb.btnGeneral.OnClicked(func(*ui.Button) {
		bottomPart.group.SetChild(boxGeneral)
		bottomPart.group.SetTitle("General")
		bb.btnGeneral.Disable()
		if !(bb.btnThreadsControl.Enabled()) {
			bb.btnThreadsControl.Enable()
		}
		if !(bb.btnSettings.Enabled()) {
			bb.btnSettings.Enable()
		}
	})

	bb.box.Append(bb.btnGeneral, false)
}

func (bb *selectingButtonsBox) initBtnThreadsControl(bottomPart bottomPartBoxMainWnd, boxThreadsControl *ui.Box) {
	bb.btnThreadsControl = ui.NewButton("Threads")

	bb.btnThreadsControl.OnClicked(func(*ui.Button) {
		bottomPart.group.SetChild(boxThreadsControl)
		bb.btnThreadsControl.Disable()
		bottomPart.group.SetTitle("Thread control")
		if !(bb.btnGeneral.Enabled()) {
			bb.btnGeneral.Enable()
		}
		if !(bb.btnSettings.Enabled()) {
			bb.btnSettings.Enable()
		}
	})

	bb.box.Append(bb.btnThreadsControl, false)
}

func (bb *selectingButtonsBox) initBtnSettings(bottomPart bottomPartBoxMainWnd, boxSettings *ui.Box) {
	bb.btnSettings = ui.NewButton("Settings")

	bb.btnSettings.OnClicked(func(*ui.Button) {
		bottomPart.group.SetChild(boxSettings)
		bottomPart.group.SetTitle("Settings")
		bb.btnSettings.Disable()
		if !(bb.btnGeneral.Enabled()) {
			bb.btnGeneral.Enable()
		}
		if !(bb.btnThreadsControl.Enabled()) {
			bb.btnThreadsControl.Enable()
		}
	})

	bb.box.Append(bb.btnSettings, false)
}

// bottomPartBoxMainWnd хранит данные о нижней части главного бокса
type bottomPartBoxMainWnd struct {
	group *ui.Group
}

func (bottomPart *bottomPartBoxMainWnd) init() {
	bottomPart.group = ui.NewGroup("General")
	bottomPart.group.SetMargined(true)
}

func (bottomPart *bottomPartBoxMainWnd) setBox(box *ui.Box) {
	bottomPart.group.SetChild(box)
}

// initGUI собирает GUI
func initGUI() {
	initFiles()

	threads := createThreads()

	boxGeneral := makeGeneralBox(threads)
	boxThreadsControl := makeThreadControlBox(threads)
	boxSettings := makeSettingsBox()

	var bmw boxMainWnd
	bmw.init()

	var bottomPart bottomPartBoxMainWnd
	bottomPart.init()
	bottomPart.setBox(boxGeneral)

	var btnsBox selectingButtonsBox
	btnsBox.init()
	btnsBox.initBtnGeneral(bottomPart, boxGeneral)
	btnsBox.initBtnThreadsControl(bottomPart, boxThreadsControl)
	btnsBox.initBtnSettings(bottomPart, boxSettings)

	var upperPart upperPartBoxMainWnd
	upperPart.init()
	upperPart.initFlexibleSpaceBox()
	upperPart.appendButtonsBox(btnsBox)
	upperPart.initFlexibleSpaceBox()

	bmw.appendUpperPart(upperPart)
	bmw.appendBottomPart(bottomPart)

	var mw mainWindow
	mw.init()
	mw.setBox(bmw)
}

// WindowSettingsKit хранит ссылки на объекты окна с установками модулей мониторинга
type WindowSettingsKit struct {
	Window *ui.Window
	Box    *ui.Box
}

// MakeSettingWindowKit создает набор для окна с установками
func MakeSettingWindowKit(windowTitle string, width, height int) WindowSettingsKit {
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

// MakeSettingCheckboxKit создает набор для поля с переключателем
func MakeSettingCheckboxKit(labelTitle string, needMonitoringFlag int) CheckboxKit {
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

// MakeSettingSpinboxKit создает набор для спинбокса
func MakeSettingSpinboxKit(labelTitle string, minValue, maxValue, currentValue int) SpinboxKit {
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

// MakeSettingEntryKit создает набор для текстового поля
func MakeSettingEntryKit(labelTitle string, entryValue string) EntryKit {
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

// MakeSettingEntryListKit создает набор для текстового поля с перечислением
func MakeSettingEntryListKit(labelTitle, jsonDump string) EntryListKit {
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

// MakeSettingComboboxKit создает набор для поля с выпадающим списком
func MakeSettingComboboxKit(labelTitle string, comboboxValues []string, currentValue string) ComboboxKit {
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

// MakeSettingButtonsKit создает набор для кнопок отмены и принятия изменений
func MakeSettingButtonsKit() ButtonsKit {
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

// ShowWarningWindow отображает окно с сообщением для пользователя
func ShowWarningWindow(warningTitle string) {
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

// NumericEntriesHandler обработчик текстовых полей, предназначенных для ввода числа
func NumericEntriesHandler(numericEntry *ui.Entry) {
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

// initFiles запускает процесс инициализации файлов
func initFiles() {
	dbHasBeenCreated, err := CheckFiles()
	if err != nil {
		ToLogFile(err.Error(), string(debug.Stack()))
		panic(err.Error())
	}

	if dbHasBeenCreated {
		message := "File of database has been created just now. Database is empty. " +
			"Need to create new access token and new subject for monitoring."
		ShowWarningWindow(message)
	}
}

// createThreads запускает процесс создания потоков с модулями проверки
func createThreads() []*Thread {
	threads, err := MakeThreads()
	if err != nil {
		ToLogFile(err.Error(), string(debug.Stack()))
		panic(err.Error())
	}

	return threads
}
