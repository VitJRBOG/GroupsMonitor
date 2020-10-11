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
	btc.groupThreadControlBtns.SetTitle(ssbd.title)
}

func (btc *boxThreadControl) setLabelNone() {
	label := ui.NewLabel("Субъекты не найдены")
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
	mcb.buttonOn = ui.NewButton("Вкл.")
	mcb.buttonOn.Disable()
	mcb.btnsBox.Append(mcb.buttonOn, true)
	mcb.buttonOff = ui.NewButton("Откл.")
	mcb.buttonOff.Disable()
	mcb.btnsBox.Append(mcb.buttonOff, true)
	mcb.statusLabel = ui.NewLabel("неактивен")
	mcb.box = ui.NewHorizontalBox()
	mcb.box.Append(mcb.monitorNameLabel, true)
	mcb.box.Append(mcb.btnsBox, true)
	mcb.box.Append(mcb.statusLabel, true)
}

func (mcb *monitorControlBox) setFuncToBtnOn(monitorName string, threadData *Thread) {
	mcb.buttonOn.OnClicked(func(*ui.Button) {
		switch monitorName {
		case "посты на стене":
			threadData.runWallPostMonitoring()
		case "фото в альбомах":
			threadData.runAlbumPhotoMonitoring()
		case "видео в альбомах":
			threadData.runVideoMonitoring()
		case "комментарии под фото":
			threadData.runPhotoCommentMonitoring()
		case "комментарии под видео":
			threadData.runVideoCommentMonitoring()
		case "комментарии в обсуждениях":
			threadData.runTopicMonitoring()
		case "комментарии под постами":
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
	monitorsNames := []string{"Посты на стене", "Фото в альбомах", "Видео в альбомах",
		"Комментарии под фото", "Комментарии под видео", "Комментарии в обсуждениях",
		"Комментарии под постами"}

	var listMCB []monitorControlBox
	for i := 0; i < len(monitorsNames); i++ {
		var mcb monitorControlBox
		mcb.init(monitorsNames[i])

		listMCB = append(listMCB, mcb)
	}

	for i := 0; i < len(listMCB); i++ {
		for n := 0; n < len(*threads); n++ {
			if (*threads)[n].Status != "неактивен" {
				if strings.ToLower((*threads)[n].Name) == strings.ToLower(subject.Name+": "+listMCB[i].title) {
					monitorName := strings.ReplaceAll((*threads)[i].Name, (*threads)[i].Subject.Name+": ", "")
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
		case "неактивен":
			ui.QueueMain(func() {
				mcb.buttonOn.Disable()
				mcb.buttonOff.Disable()
			})
		case "остановлен":
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
