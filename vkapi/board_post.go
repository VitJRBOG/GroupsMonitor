package vkapi

import (
	"fmt"
	"github.com/VitJRBOG/GroupsMonitor_new/tools"
	"runtime/debug"
	"strings"
)

type BoardPost struct {
	ID           int                   `json:"id"`
	TopicID      int                   `json:"topic_id"`
	TopicOwnerID int                   `json:"topic_owner_id"`
	FromID       int                   `json:"from_id"`
	Date         int                   `json:"date"`
	Text         string                `json:"text"`
	Attachments  []BoardPostAttachment `json:"attachments"`
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
	text := fmt.Sprintf("\n%s%d_%d?post=%d", "https://vk.com/topic",
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
			link = attachment.URl
		}
	}
	if len(attachments) > 0 {
		attachments = attachments[:len(attachments)-1]
	}

	return attachments, link
}

type BoardPostAttachment struct {
	Type      string `json:"text"`
	ID        int    `json:"id"`
	OwnerID   int    `json:"owner_id"`
	AccessKey string `json:"access_key"`
	URl       string `json:"url"`
}

func ParseBoardPostData(update UpdateFromLongPollServer) BoardPost {
	var b BoardPost

	item := update.Object

	b.ID = int(item["id"].(float64))
	b.TopicID = int(item["topic_id"].(float64))
	b.TopicOwnerID = int(item["topic_owner_id"].(float64))
	b.FromID = int(item["from_id"].(float64))
	b.Date = int(item["date"].(float64))
	b.Text = item["text"].(string)
	if attachments, exist := item["attachments"]; exist == true {
		b.Attachments = parseBoardPostAttachmentsData(attachments.([]interface{}))
	}

	return b
}

func parseBoardPostAttachmentsData(attachments []interface{}) []BoardPostAttachment {
	var bpAttachments []BoardPostAttachment

	for _, m := range attachments {
		var a BoardPostAttachment

		item := m.(map[string]interface{})

		a.Type = item["type"].(string)
		if a.Type == "photo" || a.Type == "video" || a.Type == "audio" ||
			a.Type == "doc" || a.Type == "poll" || a.Type == "link" {
			if a.Type == "link" {
				a.URl = item["link"].(map[string]interface{})["url"].(string)
			} else {
				a.OwnerID = int(item[a.Type].(map[string]interface{})["owner_id"].(float64))
				a.ID = int(item[a.Type].(map[string]interface{})["id"].(float64))
				if accessKey, exist := item[a.Type].(map[string]interface{})["access_key"]; exist {
					a.AccessKey = accessKey.(string)
				}
			}
			bpAttachments = append(bpAttachments, a)
		}
	}

	return bpAttachments
}
