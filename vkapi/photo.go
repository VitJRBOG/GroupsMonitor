package vkapi

import (
	"fmt"
	"github.com/VitJRBOG/GroupsMonitor/tools"
	"runtime/debug"
	"strings"
)

type Photo struct {
	ID      int    `json:"id"`
	AlbumID int    `json:"album_id"`
	OwnerID int    `json:"owner_id"`
	UserID  int    `json:"user_id"`
	Date    int    `json:"date"`
	Text    string `json:"text"`
}

func (p *Photo) ParseData(update UpdateFromLongPollServer) {
	item := update.Object

	p.ID = int(item["id"].(float64))
	p.AlbumID = int(item["album_id"].(float64))
	p.OwnerID = int(item["owner_id"].(float64))
	p.UserID = int(item["user_id"].(float64))
	p.Date = int(item["date"].(float64))
	p.Text = item["text"].(string)
}

func (p *Photo) SendWithMessage(getAccessToken, sendAccessToken string, operatorVkID int) error {
	var vkMsg vkMessage
	vkMsg.PeerID = operatorVkID
	vkMsg.RandomID = p.Date + p.ID // чтобы исключить пропуск фотографий, которые загрузились одновременно,
	// можно суммировать дату публикации с уникальным идентификаторам каждой фотографии
	// и использовать в качестве random_id
	vkMsg.Header, vkMsg.Text, vkMsg.Footer = p.makeTextForMessage(getAccessToken)
	vkMsg.Attachments = p.parseAttachmentsForMessage()

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

func (p *Photo) makeTextForMessage(getAccessToken string) (string, string, string) {
	hyperlinkToGroup := p.makeHyperlinkToGroup(getAccessToken, p.OwnerID)
	var hyperlinkToAuthor string
	if p.UserID != 100 {
		hyperlinkToAuthor = p.makeHyperlinkToUser(getAccessToken, p.UserID)
	} else {
		hyperlinkToAuthor = p.makeHyperlinkToGroup(getAccessToken, p.OwnerID)
	}
	date := tools.ConvertUnixTimeToDate(p.Date)
	urlToComment := p.makeURLToPhoto()

	msgHeader := fmt.Sprintf("Новое фото в альбомах\n"+
		"Расположение: %s\n"+
		"Автор: %s\n"+
		"Дата: %s",
		hyperlinkToGroup, hyperlinkToAuthor, date)
	msgText := p.Text
	msgFooter := urlToComment

	return msgHeader, msgText, msgFooter
}

func (p *Photo) makeHyperlinkToGroup(getAccessToken string, groupID int) string {
	groupInfo := getGroupInfo(getAccessToken, groupID)

	hyperlink := fmt.Sprintf("@club%d (%s)", groupInfo.ID, groupInfo.Name)
	return hyperlink
}

func (p *Photo) makeHyperlinkToUser(getAccessToken string, authorID int) string {
	userInfo := getUserInfo(getAccessToken, authorID)
	hyperlink := fmt.Sprintf("@id%d (%s %s)", userInfo.ID, userInfo.FirstName, userInfo.LastName)
	return hyperlink
}

func (p *Photo) makeURLToPhoto() string {
	text := fmt.Sprintf("\nhttps://vk.com/photo%d_%d", p.OwnerID, p.ID)
	return text
}

func (p *Photo) parseAttachmentsForMessage() string {
	attachments := fmt.Sprintf("photo%d_%d", p.OwnerID, p.ID)
	return attachments
}
