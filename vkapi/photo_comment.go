package vkapi

import (
	"fmt"
	"github.com/VitJRBOG/GroupsMonitor_new/tools"
	"runtime/debug"
	"strings"
)

type PhotoComment struct {
	ID           int                      `json:"id"`
	PhotoID      int                      `json:"photo_id"`
	PhotoOwnerID int                      `json:"photo_owner_id"`
	FromID       int                      `json:"from_id"`
	Date         int                      `json:"date"`
	Text         string                   `json:"text"`
	Attachments  []PhotoCommentAttachment `json:"attachments"`
}

func (c *PhotoComment) SendWithMessage(getAccessToken, sendAccessToken string, operatorVkID int) error {
	var vkMsg vkMessage
	vkMsg.PeerID = operatorVkID
	vkMsg.RandomID = c.Date + c.ID // чтобы исключить пропуск комментариев, которые вышли одновременно,
	// можно суммировать дату публикации с уникальным идентификаторам каждого комментария
	// и использовать в качестве random_id
	vkMsg.Header, vkMsg.Text, vkMsg.Footer = c.makeTextForMessage(getAccessToken)
	vkMsg.Attachments, vkMsg.Link = c.parseAttachmentsForMessage()

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

func (c *PhotoComment) makeTextForMessage(getAccessToken string) (string, string, string) {
	hyperlinkToGroup := c.makeHyperlinkToGroup(getAccessToken, c.PhotoOwnerID)
	var hyperlinkToAuthor string
	if c.FromID > 0 {
		hyperlinkToAuthor = c.makeHyperlinkToUser(getAccessToken, c.FromID)
	} else {
		hyperlinkToAuthor = c.makeHyperlinkToGroup(getAccessToken, c.FromID)
	}
	date := tools.ConvertUnixTimeToDate(c.Date)
	urlToComment := c.makeURLToComment()

	msgHeader := fmt.Sprintf("Новый комментарий под фото\n"+
		"Расположение: %s\n"+
		"Автор: %s\n"+
		"Дата: %s",
		hyperlinkToGroup, hyperlinkToAuthor, date)
	msgText := c.Text
	msgFooter := urlToComment

	return msgHeader, msgText, msgFooter
}

func (c *PhotoComment) makeHyperlinkToGroup(getAccessToken string, groupID int) string {
	groupInfo := getGroupInfo(getAccessToken, groupID)

	hyperlink := fmt.Sprintf("@club%d (%s)", groupInfo.ID, groupInfo.Name)
	return hyperlink
}

func (c *PhotoComment) makeHyperlinkToUser(getAccessToken string, authorID int) string {
	userInfo := getUserInfo(getAccessToken, authorID)
	hyperlink := fmt.Sprintf("@id%d (%s %s)", userInfo.ID, userInfo.FirstName, userInfo.LastName)
	return hyperlink
}

func (c *PhotoComment) makeURLToComment() string {
	text := fmt.Sprintf("\nhttps://vk.com/photo%d_%d",
		c.PhotoOwnerID, c.PhotoID)
	return text
}

func (c *PhotoComment) parseAttachmentsForMessage() (string, string) {
	var attachments string
	var link string
	for _, attachment := range c.Attachments {
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

type PhotoCommentAttachment struct {
	Type      string `json:"text"`
	ID        int    `json:"id"`
	OwnerID   int    `json:"owner_id"`
	AccessKey string `json:"access_key"`
	URl       string `json:"url"`
}

func ParsePhotoCommentData(update UpdateFromLongPollServer) PhotoComment {
	var c PhotoComment

	item := update.Object

	c.ID = int(item["id"].(float64))
	c.PhotoID = int(item["photo_id"].(float64))
	c.PhotoOwnerID = int(item["photo_owner_id"].(float64))
	c.FromID = int(item["from_id"].(float64))
	c.Date = int(item["date"].(float64))
	c.Text = item["text"].(string)
	if attachments, exist := item["attachments"]; exist == true {
		c.Attachments = parsePhotoCommentAttachmentsData(attachments.([]interface{}))
	}

	return c
}

func parsePhotoCommentAttachmentsData(attachments []interface{}) []PhotoCommentAttachment {
	var pcAttachments []PhotoCommentAttachment

	for _, m := range attachments {
		var a PhotoCommentAttachment

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
			pcAttachments = append(pcAttachments, a)
		}
	}

	return pcAttachments
}
