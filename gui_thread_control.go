package main

import (
	"runtime/debug"
	"strings"
	"time"

	"github.com/andlabs/ui"
)

// boxThreadControl хранит данные о боксе для управления потоками мониторинга субъектов
type boxThreadControl struct {
	box                    *ui.Box
	groupSubjSelectionBtns *ui.Group
	groupThreadControlBtns *ui.Group
}

func (btc *boxThreadControl) init() {
	btc.box = ui.NewHorizontalBox()
	btc.box.SetPadded(true)
	btc.groupSubjSelectionBtns = ui.NewGroup("")
	btc.groupThreadControlBtns = ui.NewGroup("")
	btc.box.Append(btc.groupSubjSelectionBtns, false)
	btc.box.Append(btc.groupThreadControlBtns, true)
}

func (btc *boxThreadControl) setSubjSelectionBtnsBox(ssbb subjSelectingBtnsBox) {
	btc.groupSubjSelectionBtns.SetChild(ssbb.box)
}

func (btc *boxThreadControl) setThreadControlBtnsBox(ssbd subjSelectingBtnData) {
	btc.groupThreadControlBtns.SetChild(ssbd.box)
}

func (btc *boxThreadControl) setLabelNone() {
	label := ui.NewLabel("Subjects not found")
	btc.box.Append(label, true)
}

// subjSelectingBtnsBox хранит данные о боксе с кнопками для переключения между субъектами
type subjSelectingBtnsBox struct {
	box *ui.Box
}

func (ssbb *subjSelectingBtnsBox) init() {
	ssbb.box = ui.NewVerticalBox()
}

func (ssbb *subjSelectingBtnsBox) appendButtons(listSSBD []subjSelectingBtnData) {
	for i := 0; i < len(listSSBD); i++ {
		ssbb.box.Append(listSSBD[i].button, false)
	}
}

// subjSelectingBtnData хранит данные о кнопках для переключения между субъектами
type subjSelectingBtnData struct {
	title  string
	button *ui.Button
	box    *ui.Box
}

func (ssbd *subjSelectingBtnData) init(threads []*Thread, subjectName string) {
	ssbd.title = subjectName
	ssbd.button = ui.NewButton(subjectName)
	ssbd.box = makeSubjectBox(subjectName, threads)
}

func (ssbd *subjSelectingBtnData) setFuncToBtn(listSSBD []subjSelectingBtnData, btc boxThreadControl) {
	ssbd.button.OnClicked(func(*ui.Button) {
		btc.groupThreadControlBtns.SetChild(ssbd.box)
		btc.groupThreadControlBtns.SetTitle(ssbd.title)

		for i := 0; i < len(listSSBD); i++ {
			if !(listSSBD[i].button.Enabled()) {
				listSSBD[i].button.Enable()
			}
		}
		ssbd.button.Disable()
	})
}

// makeThreadControlBox собирает бокс для контроля потоков мониторинга субъектов
func makeThreadControlBox(threads []*Thread) *ui.Box {

	subjectNames := getSubjectsNames()

	var btc boxThreadControl
	btc.init()

	var listSSBD []subjSelectingBtnData
	if len(subjectNames) > 0 {
		for i := 0; i < len(subjectNames); i++ {
			var ssbd subjSelectingBtnData
			ssbd.init(threads, subjectNames[i])
			listSSBD = append(listSSBD, ssbd)
		}
	} else {
		btc.setLabelNone()
	}

	var ssbb subjSelectingBtnsBox
	ssbb.init()

	if len(listSSBD) > 0 {
		for i := 0; i < len(listSSBD); i++ {
			listSSBD[i].setFuncToBtn(listSSBD, btc)
		}

		ssbb.appendButtons(listSSBD)
		btc.setSubjSelectionBtnsBox(ssbb)

		btc.setThreadControlBtnsBox(listSSBD[0])
		listSSBD[0].button.Disable()
	}

	return btc.box
}

func getSubjectsNames() []string {
	var subjectsNames []string

	// запрашиваем список субъектов из базы данных
	var dbKit DataBaseKit
	subjects, err := dbKit.selectTableSubject()
	if err != nil {
		ToLogFile(err.Error(), string(debug.Stack()))
		panic(err.Error())
	}
	// и добавляем их названия в список
	for _, subject := range subjects {
		subjectsNames = append(subjectsNames, subject.Name)
	}

	return subjectsNames
}

// subjectBoxData хранит данные о боксе с кнопками для управления потоками мониторинга одного из субъектов
type subjectBoxData struct {
	box *ui.Box
}

func (sbd *subjectBoxData) init() {
	sbd.box = ui.NewVerticalBox()
}

func (sbd *subjectBoxData) appendMonitorControlBoxes(listMCB []monitorControlBox) {
	for i := 0; i < len(listMCB); i++ {
		sbd.box.Append(listMCB[i].box, false)
	}
}

// monitorControlBox хранит данные о кнопке для управления одним из потоков мониторинга субъекта
type monitorControlBox struct {
	box    *ui.Box
	title  string
	button *ui.Button
	label  *ui.Label
}

func (mcb *monitorControlBox) init(buttonName string) {
	mcb.box = ui.NewHorizontalBox()
	mcb.title = buttonName
	mcb.button = ui.NewButton(buttonName)
	mcb.button.Disable()
	mcb.label = ui.NewLabel("inactive")
	mcb.box.Append(mcb.button, true)
	mcb.box.Append(mcb.label, true)
}

func (mcb *monitorControlBox) setFuncToBtn(threadData *Thread) {
	mcb.button.OnClicked(func(*ui.Button) {
		if threadData.ActionFlag == 3 {
			threadData.ActionFlag = 2
		} else {
			if threadData.ActionFlag != 1 {
				threadData.ActionFlag = 3
			}
		}
	})

	mcb.button.Enable()
	go threadStatusChecking(mcb.label, threadData)
}

// makeSubjectBox собирает бокс для управления потоками мониторинга одного из субъектов
func makeSubjectBox(subjectName string, threads []*Thread) *ui.Box {
	monitorsNames := []string{"Wall post monitoring", "Album photo monitoring", "Video monitoring",
		"Photo comment monitoring", "Video comment monitoring", "Topic monitoring",
		"Wall post comment monitoring"}

	var listMCB []monitorControlBox
	for i := 0; i < len(monitorsNames); i++ {
		var mcb monitorControlBox
		mcb.init(monitorsNames[i])

		listMCB = append(listMCB, mcb)
	}

	for i := 0; i < len(listMCB); i++ {
		for n := 0; n < len(threads); n++ {
			if strings.Contains(strings.ToLower(threads[n].Name), strings.ToLower(subjectName)) {
				if strings.ToLower(threads[n].Name) == strings.ToLower(subjectName+"'s "+listMCB[i].title) {
					listMCB[i].setFuncToBtn(threads[n])
				}
			}
		}
	}

	var sbd subjectBoxData
	sbd.init()
	sbd.appendMonitorControlBoxes(listMCB)

	return sbd.box
}

func threadStatusChecking(statusLabel *ui.Label, threadData *Thread) {
	for true {
		if statusLabel.Text() != threadData.Status {
			statusLabel.SetText(threadData.Status)
		}
		time.Sleep(1 * time.Second)
	}
}
