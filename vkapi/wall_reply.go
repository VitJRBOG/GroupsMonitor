package vkapi

import (
	"fmt"
	"runtime/debug"
	"strings"

	"github.com/VitJRBOG/GroupsObserver/tools"
)

type WallReply struct {
	ID           int          `json:"id"`
	PostID       int          `json:"post_id"`
	OwnerID      int          `json:"owner_id"`
	FromID       int          `json:"from_id"`
	ParentsStack int          `json:"parents_stack"`
	Date         int          `json:"date"`
	Text         string       `json:"text"`
	Attachments  []attachment `json:"attachments"`
}

func (w *WallReply) ParseData(update UpdateFromLongPollServer) {
	item := update.Object

	w.ID = int(item["id"].(float64))
	w.PostID = int(item["post_id"].(float64))
	w.OwnerID = int(item["owner_id"].(float64))
	w.FromID = int(item["from_id"].(float64))
	w.parseParentsStack(item["parents_stack"].([]interface{}))
	w.Date = int(item["date"].(float64))
	w.Text = item["text"].(string)
	if attachments, exist := item["attachments"]; exist == true {
		w.Attachments = parseAttachmentsData(attachments.([]interface{}))
	}
}

func (w *WallReply) parseParentsStack(parentsStack []interface{}) {
	if len(parentsStack) > 0 {
		w.ParentsStack = int(parentsStack[0].(float64))
	}
}

func (w *WallReply) SendWithMessage(getAccessToken, sendAccessToken string, operatorVkID int) error {
	var vkMsg vkMessage
	vkMsg.PeerID = operatorVkID
	vkMsg.RandomID = w.Date + w.ID // чтобы исключить пропуск комментариев, которые вышли одновременно,
	// можно суммировать дату публикации с уникальным идентификаторам каждого комментария
	// и использовать в качестве random_id
	vkMsg.Header, vkMsg.Text, vkMsg.Footer = w.makeTextForMessage(getAccessToken)
	vkMsg.Attachments, vkMsg.Link = w.parseAttachmentsForMessage()
	vkMsg.ContentSource = w.parseContentSource()

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

func (w *WallReply) makeTextForMessage(getAccessToken string) (string, string, string) {
	hyperlinkToGroup := w.makeHyperlinkToGroup(getAccessToken, w.OwnerID)
	var hyperlinkToAuthor string
	if w.FromID > 0 {
		hyperlinkToAuthor = w.makeHyperlinkToUser(getAccessToken, w.FromID)
	} else {
		hyperlinkToAuthor = w.makeHyperlinkToGroup(getAccessToken, w.FromID)
	}
	date := tools.ConvertUnixTimeToDate(w.Date)
	urlToComment := w.makeURLToComment()

	msgHeader := fmt.Sprintf("Новый комментарий на стене\n"+
		"Расположение: %s\n"+
		"Автор: %s\n"+
		"Дата: %s",
		hyperlinkToGroup, hyperlinkToAuthor, date)
	msgText := w.Text
	msgFooter := urlToComment

	return msgHeader, msgText, msgFooter
}

func (w *WallReply) makeHyperlinkToGroup(getAccessToken string, groupID int) string {
	groupInfo := getGroupInfo(getAccessToken, groupID)

	hyperlink := fmt.Sprintf("@club%d (%s)", groupInfo.ID, groupInfo.Name)
	return hyperlink
}

func (w *WallReply) makeHyperlinkToUser(getAccessToken string, authorID int) string {
	userInfo := getUserInfo(getAccessToken, authorID)
	hyperlink := fmt.Sprintf("@id%d (%s %s)", userInfo.ID, userInfo.FirstName, userInfo.LastName)
	return hyperlink
}

func (w *WallReply) makeURLToComment() string {
	text := fmt.Sprintf("\nhttps://vk.com/wall%d_%d?reply=%d", w.OwnerID, w.PostID, w.ID)
	if w.ParentsStack > 0 {
		text = fmt.Sprintf("%s&thread=%d", text, w.ParentsStack)
	}
	return text
}

func (w *WallReply) parseAttachmentsForMessage() (string, string) {
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

func (w *WallReply) parseContentSource() string {
	contentSource := fmt.Sprintf(`{"type": "url", "url": "https://vk.com/wall%d_%d?reply=%d"}`,
		w.OwnerID, w.PostID, w.ID)
	return contentSource
}
