package vkapi

import (
	"fmt"
	"runtime/debug"
	"strings"

	"github.com/VitJRBOG/Watcher/tools"
)

type Video struct {
	ID          int    `json:"id"`
	OwnerID     int    `json:"owner_id"`
	UserID      int    `json:"user_id"`
	Date        int    `json:"date"`
	Description string `json:"description"`
}

func (v *Video) ParseData(update UpdateFromLongPollServer) {
	item := update.Object

	v.ID = int(item["id"].(float64))
	v.OwnerID = int(item["owner_id"].(float64))
	if userID, exist := item["user_id"]; exist == true {
		v.UserID = int(userID.(float64))
	} else {
		v.UserID = 0
	}
	v.Date = int(item["date"].(float64))
	v.Description = item["description"].(string)
}

func (v *Video) SendWithMessage(getAccessToken, sendAccessToken string, operatorVkID int) error {
	var vkMsg vkMessage
	vkMsg.PeerID = operatorVkID
	vkMsg.RandomID = v.Date + v.ID // чтобы исключить пропуск видео, которые загрузились одновременно,
	// можно суммировать дату публикации с уникальным идентификаторам каждго видео
	// и использовать в качестве random_id
	vkMsg.Header, vkMsg.Text, vkMsg.Footer = v.makeTextForMessage(getAccessToken)
	vkMsg.Attachments = v.parseAttachmentsForMessage()
	vkMsg.ContentSource = v.parseContentSource()

	err := vkMsg.sendMessage(sendAccessToken)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "too much messages sent to user") {
			return err
		} else {
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
	}
	return nil
}

func (v *Video) makeTextForMessage(getAccessToken string) (string, string, string) {
	hyperlinkToGroup := v.makeHyperlinkToGroup(getAccessToken, v.OwnerID)
	var hyperlinkToAuthor string
	if v.UserID > 0 {
		hyperlinkToAuthor = v.makeHyperlinkToUser(getAccessToken, v.UserID)
	} else {
		hyperlinkToAuthor = v.makeHyperlinkToGroup(getAccessToken, v.OwnerID)
	}
	date := tools.ConvertUnixTimeToDate(v.Date)
	urlToComment := v.makeURLToVideo()

	msgHeader := fmt.Sprintf("Новое видео в альбомах\n"+
		"Расположение: %s\n"+
		"Автор: %s\n"+
		"Дата: %s",
		hyperlinkToGroup, hyperlinkToAuthor, date)
	msgText := v.Description
	msgFooter := urlToComment

	return msgHeader, msgText, msgFooter
}

func (v *Video) makeHyperlinkToGroup(getAccessToken string, groupID int) string {
	groupInfo := getGroupInfo(getAccessToken, groupID)

	hyperlink := fmt.Sprintf("@club%d (%s)", groupInfo.ID, groupInfo.Name)
	return hyperlink
}

func (v *Video) makeHyperlinkToUser(getAccessToken string, authorID int) string {
	userInfo := getUserInfo(getAccessToken, authorID)
	hyperlink := fmt.Sprintf("@id%d (%s %s)", userInfo.ID, userInfo.FirstName, userInfo.LastName)
	return hyperlink
}

func (v *Video) makeURLToVideo() string {
	text := fmt.Sprintf("\nhttps://vk.com/video%d_%d", v.OwnerID, v.ID)
	return text
}

func (v *Video) parseAttachmentsForMessage() string {
	attachments := fmt.Sprintf("photo%d_%d", v.OwnerID, v.ID)
	return attachments
}

func (v *Video) parseContentSource() string {
	contentSource := fmt.Sprintf(`{"type": "url", "url": "https://vk.com/video%d_%d"}`,
		v.OwnerID, v.ID)
	return contentSource
}
