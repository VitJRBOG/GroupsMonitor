package vkapi

import (
	"fmt"
	"github.com/VitJRBOG/GroupsObserver/tools"
	"runtime/debug"
	"strings"
)

type WallPost struct {
	ID          int          `json:"id"`
	OwnerID     int          `json:"owner_id"`
	FromID      int          `json:"from_id"`
	SignerID    int          `json:"signer_id"`
	Date        int          `json:"date"`
	PostType    string       `json:"post_type"`
	Text        string       `json:"text"`
	Attachments []attachment `json:"attachments"`
}

func (w *WallPost) ParseData(update UpdateFromLongPollServer) {
	item := update.Object

	w.ID = int(item["id"].(float64))
	w.OwnerID = int(item["owner_id"].(float64))
	w.FromID = int(item["from_id"].(float64))
	if signerID, exist := item["signer_id"]; exist == true {
		w.SignerID = int(signerID.(float64))
	} else {
		w.SignerID = 0
	}
	w.Date = int(item["date"].(float64))
	w.PostType = item["post_type"].(string)
	w.Text = item["text"].(string)
	if attachments, exist := item["attachments"]; exist == true {
		w.Attachments = parseAttachmentsData(attachments.([]interface{}))
	}
}

func (w *WallPost) SendWithMessage(getAccessToken, sendAccessToken string, operatorVkID int) error {
	var vkMsg vkMessage
	vkMsg.PeerID = operatorVkID
	vkMsg.RandomID = w.Date + w.ID // чтобы исключить пропуск постов, которые вышли одновременно,
	// можно суммировать дату публикации с уникальным идентификаторам каждого поста
	// и использовать в качестве random_id
	vkMsg.Header, vkMsg.Text, vkMsg.Footer = w.makeTextForMessage(getAccessToken)
	vkMsg.Attachments, vkMsg.Link = w.parseAttachmentsForMessage()

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

func (w *WallPost) makeTextForMessage(getAccessToken string) (string, string, string) {
	postType := w.selectPostType()
	hyperlinkToGroup := w.makeHyperlinkToGroup(getAccessToken)
	var hyperlinkToUser string
	if w.SignerID != 0 {
		hyperlinkToUser = w.makeHyperlinkToUser(getAccessToken, w.SignerID)
	} else {
		if w.OwnerID != w.FromID {
			hyperlinkToUser = w.makeHyperlinkToUser(getAccessToken, w.FromID)
		} else {
			hyperlinkToUser = "[нет данных]"
		}
	}
	date := tools.ConvertUnixTimeToDate(w.Date)
	urlToPost := w.makeURLToPost()

	msgHeader := fmt.Sprintf("Новый %s на стене\n"+
		"Расположение: %s\n"+
		"Автор: %s\n"+
		"Дата: %s",
		postType, hyperlinkToGroup, hyperlinkToUser, date)
	msgText := w.Text
	msgFooter := urlToPost

	return msgHeader, msgText, msgFooter
}

func (w *WallPost) selectPostType() string {
	var postType string
	switch w.PostType {
	case "post":
		postType = "опубликованный пост"
	case "suggest":
		postType = "предложенный пост"
	case "postpone":
		postType = "отложенный пост"
	default:
		postType = "пост"
	}
	return postType
}

func (w *WallPost) makeHyperlinkToGroup(getAccessToken string) string {
	groupInfo := getGroupInfo(getAccessToken, w.OwnerID)

	hyperlink := fmt.Sprintf("@club%d (%s)", groupInfo.ID, groupInfo.Name)
	return hyperlink
}

func (w *WallPost) makeHyperlinkToUser(getAccessToken string, authorID int) string {
	userInfo := getUserInfo(getAccessToken, authorID)
	hyperlink := fmt.Sprintf("@id%d (%s %s)", userInfo.ID, userInfo.FirstName, userInfo.LastName)
	return hyperlink
}

func (w *WallPost) makeURLToPost() string {
	text := fmt.Sprintf("\nhttps://vk.com/wall%d_%d", w.OwnerID, w.ID)
	return text
}

func (w *WallPost) parseAttachmentsForMessage() (string, string) {
	var attachments string
	var link string
	for _, attachment := range w.Attachments {
		if attachment.Type != "link" {
			attachments += fmt.Sprintf("%s%d_%d",
				attachment.Type, attachment.OwnerID, attachment.ID)
			if len(attachment.AccessKey) > 0 {
				attachments += fmt.Sprintf("_%s", attachment.AccessKey)
			}
			attachments += ","
		} else {
			link = attachment.URL
		}
	}
	if len(attachments) > 0 {
		attachments = attachments[:len(attachments)-1]
	}

	return attachments, link
}
