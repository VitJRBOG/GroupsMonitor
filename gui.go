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

func (wsk *WindowSettingsKit) init(windowTitle string, width, height int) {
	wsk.Window = ui.NewWindow(windowTitle, width, height, true)
	wsk.Window.OnClosing(func(*ui.Window) bool {
		wsk.Window.Disable()
		return true
	})
	wsk.Window.SetMargined(true)
}

func (wsk *WindowSettingsKit) initBox() {
	wsk.Box = ui.NewVerticalBox()
	wsk.Box.SetPadded(true)
	wsk.Window.SetChild(wsk.Box)
}

// MakeSettingWindowKit создает набор для окна с установками
func MakeSettingWindowKit(windowTitle string, width, height int) WindowSettingsKit {
	var wsk WindowSettingsKit
	wsk.init(windowTitle, width, height)
	wsk.initBox()

	return wsk
}

// CheckboxKit хранит ссылки на объекты для параметров с переключателями
type CheckboxKit struct {
	Box      *ui.Box
	CheckBox *ui.Checkbox
}

func (cbk *CheckboxKit) init() {
	cbk.Box = ui.NewHorizontalBox()
	cbk.Box.SetPadded(true)
}

func (cbk *CheckboxKit) initLabel(labelTitle string) {
	label := ui.NewLabel(labelTitle)
	cbk.Box.Append(label, true)
}

func (cbk *CheckboxKit) initCheckbox(needMonitoringFlag int) {
	cbk.CheckBox = ui.NewCheckbox("")
	if needMonitoringFlag == 1 {
		cbk.CheckBox.SetChecked(true)
	} else {
		cbk.CheckBox.SetChecked(false)
	}
	cbk.Box.Append(cbk.CheckBox, true)
}

// MakeSettingCheckboxKit создает набор для поля с переключателем
func MakeSettingCheckboxKit(labelTitle string, needMonitoringFlag int) CheckboxKit {
	var cbk CheckboxKit
	cbk.init()
	cbk.initLabel(labelTitle)
	cbk.initCheckbox(needMonitoringFlag)

	return cbk
}

// SpinboxKit хранит ссылки на объекты для параметров с спинбоксом
type SpinboxKit struct {
	Box     *ui.Box
	Spinbox *ui.Spinbox
}

func (sbk *SpinboxKit) init() {
	sbk.Box = ui.NewHorizontalBox()
	sbk.Box.SetPadded(true)
}

func (sbk *SpinboxKit) initLabel(labelTitle string) {
	label := ui.NewLabel(labelTitle)
	sbk.Box.Append(label, true)
}

func (sbk *SpinboxKit) initSpinbox(minValue, maxValue, currentValue int) {
	sbk.Spinbox = ui.NewSpinbox(minValue, maxValue)
	sbk.Spinbox.SetValue(currentValue)
	sbk.Box.Append(sbk.Spinbox, true)
}

// MakeSettingSpinboxKit создает набор для спинбокса
func MakeSettingSpinboxKit(labelTitle string, minValue, maxValue, currentValue int) SpinboxKit {
	var sbk SpinboxKit
	sbk.init()
	sbk.initLabel(labelTitle)
	sbk.initSpinbox(minValue, maxValue, currentValue)

	return sbk
}

// EntryKit хранит ссылки на объекты для параметров с полями для ввода текста
type EntryKit struct {
	Box   *ui.Box
	Entry *ui.Entry
}

func (ek *EntryKit) init() {
	ek.Box = ui.NewHorizontalBox()
	ek.Box.SetPadded(true)
}

func (ek *EntryKit) initLabel(labelTitle string) {
	label := ui.NewLabel(labelTitle)
	ek.Box.Append(label, true)
}

func (ek *EntryKit) initEntry(entryValue string) {
	ek.Entry = ui.NewEntry()
	ek.Entry.SetText(entryValue)
	ek.Box.Append(ek.Entry, true)
}

// MakeSettingEntryKit создает набор для текстового поля
func MakeSettingEntryKit(labelTitle string, entryValue string) EntryKit {
	var ek EntryKit
	ek.init()
	ek.initLabel(labelTitle)
	ek.initEntry(entryValue)

	return ek
}

// EntryListKit хранит ссылки на объекты для параметров со списком в поле для ввода текста
type EntryListKit struct {
	Box   *ui.Box
	Entry *ui.Entry
}

func (elk *EntryListKit) init() {
	elk.Box = ui.NewHorizontalBox()
	elk.Box.SetPadded(true)
}

func (elk *EntryListKit) initLabel(labelTitle string) {
	label := ui.NewLabel(labelTitle)
	elk.Box.Append(label, true)
}

func (elk *EntryListKit) initEntry(jsonDump string) {
	elk.Entry = ui.NewEntry()
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
		elk.Entry.SetText(list)
	}
	elk.Box.Append(elk.Entry, true)
}

// MakeSettingEntryListKit создает набор для текстового поля с перечислением
func MakeSettingEntryListKit(labelTitle, jsonDump string) EntryListKit {
	var elk EntryListKit
	elk.init()
	elk.initLabel(labelTitle)
	elk.initEntry(jsonDump)

	return elk
}

// ComboboxKit хранит ссылки на объекты для параметров с выпадающим списком
type ComboboxKit struct {
	Box      *ui.Box
	Combobox *ui.Combobox
}

func (cbk *ComboboxKit) init() {
	cbk.Box = ui.NewHorizontalBox()
	cbk.Box.SetPadded(true)
}

func (cbk *ComboboxKit) initLabel(labelTitle string) {
	label := ui.NewLabel(labelTitle)
	cbk.Box.Append(label, true)
}

func (cbk *ComboboxKit) initCombobox(comboboxValues []string, currentValue string) {
	cbk.Combobox = ui.NewCombobox()
	slctd := -1
	for i, item := range comboboxValues {
		cbk.Combobox.Append(item)
		if currentValue == item {
			slctd = i
		}
	}
	cbk.Combobox.SetSelected(slctd)
	cbk.Box.Append(cbk.Combobox, true)
}

// MakeSettingComboboxKit создает набор для поля с выпадающим списком
func MakeSettingComboboxKit(labelTitle string, comboboxValues []string, currentValue string) ComboboxKit {
	var cbk ComboboxKit
	cbk.init()
	cbk.initLabel(labelTitle)
	cbk.initCombobox(comboboxValues, currentValue)

	return cbk
}

// ButtonsKit хранит ссылки на объекты для кнопок принятия и отмены изменений в установках
type ButtonsKit struct {
	Box          *ui.Box
	ButtonApply  *ui.Button
	ButtonCancel *ui.Button
}

func (bk *ButtonsKit) init() {
	bk.Box = ui.NewHorizontalBox()
	bk.Box.SetPadded(true)
}

func (bk *ButtonsKit) initFlexibleSpaceBox() {
	box := ui.NewHorizontalBox()
	bk.Box.Append(box, false)
}

func (bk *ButtonsKit) initButtons() {
	buttonsBox := ui.NewHorizontalBox()
	bk.ButtonApply = ui.NewButton("Apply")
	bk.ButtonCancel = ui.NewButton("Cancel")
	buttonsBox.Append(bk.ButtonApply, false)
	buttonsBox.Append(bk.ButtonCancel, false)
	bk.Box.Append(buttonsBox, false)
}

// MakeSettingButtonsKit создает набор для кнопок отмены и принятия изменений
func MakeSettingButtonsKit() ButtonsKit {
	var bk ButtonsKit
	bk.init()
	bk.initFlexibleSpaceBox()
	bk.initFlexibleSpaceBox()
	bk.initButtons()

	return bk
}

// warningWindow хранит данные об окне с сообщением для пользователя
type warningWindow struct {
	window *ui.Window
}

func (ww *warningWindow) init() {
	ww.window = ui.NewWindow("WARNING!", 100, 100, true)
	ww.window.SetMargined(true)
	ww.window.OnClosing(func(*ui.Window) bool {
		ww.window.Disable()
		return true
	})
	ww.window.Show()
}

func (ww *warningWindow) setBox(box boxWarningWnd) {
	ww.window.SetChild(box.box)
}

// boxWarningWnd хранит данные о главном боксе для окна с сообщением
type boxWarningWnd struct {
	box *ui.Box
}

func (bww *boxWarningWnd) init() {
	bww.box = ui.NewVerticalBox()
	bww.box.SetPadded(true)
}

func (bww *boxWarningWnd) appendInfoBox(bwi boxWarningInfo) {
	bww.box.Append(bwi.box, false)
}

func (bww *boxWarningWnd) appendBtnsBox(bwb boxWarningBtn) {
	bww.box.Append(bwb.box, false)
}

// boxWarningInfo хранит данные о информационной части главного бокса
type boxWarningInfo struct {
	box *ui.Box
}

func (bwi *boxWarningInfo) init(warningTitle string) {
	bwi.box = ui.NewVerticalBox()
	titleLabel := ui.NewLabel(warningTitle)
	bwi.box.Append(titleLabel, true)
}

// boxWarningBtn хранит данные о части главного бокса, на которой размещена кнопка OK
type boxWarningBtn struct {
	box   *ui.Box
	btnOk *ui.Button
}

func (bwb *boxWarningBtn) init() {
	bwb.box = ui.NewHorizontalBox()
}

func (bwb *boxWarningBtn) initFlexibleSpaceBox() {
	box := ui.NewVerticalBox()
	bwb.box.Append(box, true)
}

func (bwb *boxWarningBtn) initBtnOk(ww warningWindow) {
	bwb.btnOk = ui.NewButton("OK")
	bwb.btnOk.OnClicked(func(*ui.Button) {
		ww.window.Hide()
	})
	bwb.box.Append(bwb.btnOk, false)
}

// ShowWarningWindow отображает окно с сообщением для пользователя
func ShowWarningWindow(warningTitle string) {
	var bww boxWarningWnd
	bww.init()

	var ww warningWindow
	ww.init()
	ww.setBox(bww)

	var bwi boxWarningInfo
	bwi.init(warningTitle)

	var bwb boxWarningBtn
	bwb.init()
	bwb.initFlexibleSpaceBox()
	bwb.initBtnOk(ww)

	bww.appendInfoBox(bwi)
	bww.appendBtnsBox(bwb)
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
