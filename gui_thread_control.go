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

func (ssbd *subjSelectingBtnData) init(threads *[]*Thread, subject Subject) {
	ssbd.title = subject.Name
	ssbd.button = ui.NewButton(subject.Name)
	ssbd.box = makeSubjectBox(subject, threads)
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
func makeThreadControlBox(threads *[]*Thread) *ui.Box {

	var btc boxThreadControl
	btc.init()

	var dbKit DataBaseKit
	subjects, err := dbKit.selectTableSubject()
	if err != nil {
		ToLogFile(err.Error(), string(debug.Stack()))
		panic(err.Error())
	}

	var listSSBD []subjSelectingBtnData
	if len(subjects) > 0 {
		for i := 0; i < len(subjects); i++ {
			var ssbd subjSelectingBtnData
			ssbd.init(threads, subjects[i])
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

// monitorControlBox хранит данные о кнопках включения и отключения одного из потоков мониторинга субъекта
type monitorControlBox struct {
	title            string
	box              *ui.Box
	monitorNameLabel *ui.Label
	btnsBox          *ui.Box
	buttonOn         *ui.Button
	buttonOff        *ui.Button
	statusLabel      *ui.Label
}

func (mcb *monitorControlBox) init(monitorName string) {
	mcb.title = monitorName
	mcb.monitorNameLabel = ui.NewLabel(monitorName)
	mcb.btnsBox = ui.NewHorizontalBox()
	mcb.buttonOn = ui.NewButton("On")
	mcb.buttonOn.Disable()
	mcb.btnsBox.Append(mcb.buttonOn, true)
	mcb.buttonOff = ui.NewButton("Off")
	mcb.buttonOff.Disable()
	mcb.btnsBox.Append(mcb.buttonOff, true)
	mcb.statusLabel = ui.NewLabel("inactive")
	mcb.box = ui.NewHorizontalBox()
	mcb.box.Append(mcb.monitorNameLabel, true)
	mcb.box.Append(mcb.btnsBox, true)
	mcb.box.Append(mcb.statusLabel, true)
}

func (mcb *monitorControlBox) setFuncToBtnOn(monitorName string, threadData *Thread) {
	mcb.buttonOn.OnClicked(func(*ui.Button) {
		switch monitorName {
		case "wall post monitoring":
			threadData.runWallPostMonitoring()
		case "album photo monitoring":
			threadData.runAlbumPhotoMonitoring()
		case "video monitoring":
			threadData.runVideoMonitoring()
		case "photo comment monitoring":
			threadData.runPhotoCommentMonitoring()
		case "video comment monitoring":
			threadData.runVideoCommentMonitoring()
		case "topic monitoring":
			threadData.runTopicMonitoring()
		case "wall post comment monitoring":
			threadData.runWallPostCommentMonitoring()
		}

		mcb.buttonOn.Disable()
		mcb.buttonOff.Enable()
	})
}

func (mcb *monitorControlBox) setFuncToBtnOff(threadData *Thread) {
	mcb.buttonOff.OnClicked(func(*ui.Button) {
		threadData.ActionFlag = 1

		mcb.buttonOn.Enable()
		mcb.buttonOff.Disable()
	})
}

func (mcb *monitorControlBox) setThreadStatusChecking(threadData *Thread) {
	go threadStatusChecking(mcb, threadData)
}

// makeSubjectBox собирает бокс для управления потоками мониторинга одного из субъектов
func makeSubjectBox(subject Subject, threads *[]*Thread) *ui.Box {
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
		for n := 0; n < len(*threads); n++ {
			if (*threads)[n].Status != "inactive" {
				if strings.ToLower((*threads)[n].Name) == strings.ToLower(subject.Name+"'s "+listMCB[i].title) {
					monitorName := strings.ReplaceAll((*threads)[i].Name, (*threads)[i].Subject.Name+"'s ", "")
					listMCB[i].setFuncToBtnOn(monitorName, (*threads)[n])
					listMCB[i].setFuncToBtnOff((*threads)[n])
					listMCB[i].setThreadStatusChecking((*threads)[n])
				}
			}
		}
	}

	var sbd subjectBoxData
	sbd.init()
	sbd.appendMonitorControlBoxes(listMCB)

	return sbd.box
}

func threadStatusChecking(mcb *monitorControlBox, threadData *Thread) {
	for true {
		if mcb.statusLabel.Text() != threadData.Status {
			ui.QueueMain(func() {
				mcb.statusLabel.SetText(threadData.Status)
			})
		}
		switch threadData.Status {
		case "inactive":
			ui.QueueMain(func() {
				mcb.buttonOn.Disable()
				mcb.buttonOff.Disable()
			})
		case "stopped":
			ui.QueueMain(func() {
				mcb.buttonOn.Enable()
				mcb.buttonOff.Disable()
			})
		default:
			ui.QueueMain(func() {
				mcb.buttonOn.Disable()
				mcb.buttonOff.Enable()
			})
		}
		time.Sleep(1 * time.Second)
	}
}
