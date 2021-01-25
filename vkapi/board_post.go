package vkapi

import (
	"fmt"
	"github.com/VitJRBOG/GroupsMonitor/tools"
	"runtime/debug"
	"strings"
)

type BoardPost struct {
	ID           int          `json:"id"`
	TopicID      int          `json:"topic_id"`
	TopicOwnerID int          `json:"topic_owner_id"`
	FromID       int          `json:"from_id"`
	Date         int          `json:"date"`
	Text         string       `json:"text"`
	Attachments  []attachment `json:"attachments"`
}

func (b *BoardPost) ParseData(update UpdateFromLongPollServer) {
	item := update.Object

	b.ID = int(item["id"].(float64))
	b.TopicID = int(item["topic_id"].(float64))
	b.TopicOwnerID = int(item["topic_owner_id"].(float64))
	b.FromID = int(item["from_id"].(float64))
	b.Date = int(item["date"].(float64))
	b.Text = item["text"].(string)
	if attachments, exist := item["attachments"]; exist == true {
		b.Attachments = parseAttachmentsData(attachments.([]interface{}))
	}
}

func (b *BoardPost) SendWithMessage(getAccessToken, sendAccessToken string, operatorVkID int) error {
	var vkMsg vkMessage
	vkMsg.PeerID = operatorVkID
	vkMsg.RandomID = b.Date + b.ID // чтобы исключить пропуск постов, которые вышли одновременно,
	// можно суммировать дату публикации с уникальным идентификаторам каждого поста
	// и использовать в качестве random_id
	vkMsg.Header, vkMsg.Text, vkMsg.Footer = b.makeTextForMessage(getAccessToken)
	vkMsg.Attachments, vkMsg.Link = b.parseAttachmentsForMessage()

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

func (b *BoardPost) makeTextForMessage(getAccessToken string) (string, string, string) {
	hyperlinkToGroup := b.makeHyperlinkToGroup(getAccessToken, b.TopicOwnerID)
	var hyperlinkToAuthor string
	if b.FromID > 0 {
		hyperlinkToAuthor = b.makeHyperlinkToUser(getAccessToken, b.FromID)
	} else {
		hyperlinkToAuthor = b.makeHyperlinkToGroup(getAccessToken, b.FromID)
	}
	date := tools.ConvertUnixTimeToDate(b.Date)
	urlToComment := b.makeURLToComment()

	msgHeader := fmt.Sprintf("Новый пост в обсуждениях\n"+
		"Расположение: %s\n"+
		"Автор: %s\n"+
		"Дата: %s",
		hyperlinkToGroup, hyperlinkToAuthor, date)
	msgText := b.Text
	msgFooter := urlToComment

	return msgHeader, msgText, msgFooter
}

func (b *BoardPost) makeHyperlinkToGroup(getAccessToken string, groupID int) string {
	groupInfo := getGroupInfo(getAccessToken, groupID)

	hyperlink := fmt.Sprintf("@club%d (%s)", groupInfo.ID, groupInfo.Name)
	return hyperlink
}

func (b *BoardPost) makeHyperlinkToUser(getAccessToken string, authorID int) string {
	userInfo := getUserInfo(getAccessToken, authorID)
	hyperlink := fmt.Sprintf("@id%d (%s %s)", userInfo.ID, userInfo.FirstName, userInfo.LastName)
	return hyperlink
}

func (b *BoardPost) makeURLToComment() string {
	text := fmt.Sprintf("\nhttps://vk.com/topic%d_%d?post=%d",
		b.TopicOwnerID, b.TopicID, b.ID)
	return text
}

func (b *BoardPost) parseAttachmentsForMessage() (string, string) {
	var attachments string
	var link string
	for _, attachment := range b.Attachments {
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
