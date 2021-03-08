package vkapi

import (
	"encoding/json"
	"net/url"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	govkapi "github.com/VitJRBOG/GoVkApi/v3"
	"github.com/VitJRBOG/GroupsObserver/tools"
)

type longPollServerConnectionData struct {
	Key    string `json:"key"`
	Server string `json:"server"`
	TS     string `json:"ts"`
}

func getLongPollServerConnectionData(accessToken string, wardVkId int) longPollServerConnectionData {
	values := url.Values{
		"access_token": {accessToken},
		"group_id":     {strconv.Itoa(wardVkId)},
		"v":            {"5.126"},
	}

	response, err := callMethod("groups.getLongPollServer", values)
	if err != nil {
		tools.WriteToLog(err, debug.Stack())
		panic(err.Error())
	}

	lpsConnectionData := parseConnectionDataForLongPollServer(response)
	return lpsConnectionData
}

func parseConnectionDataForLongPollServer(response []byte) longPollServerConnectionData {
	var lpsConnectionData longPollServerConnectionData
	err := json.Unmarshal(response, &lpsConnectionData)
	if err != nil {
		tools.WriteToLog(err, debug.Stack())
		panic(err.Error())
	}
	return lpsConnectionData
}

type LongPollApiSettings struct {
	IsEnabled       bool   `json:"is_enabled"`
	ApiVersion      string `json:"api_version"`
	WallPostNew     int    `json:"wall_post_new"`
	WallReplyNew    int    `json:"wall_reply_new"`
	PhotoNew        int    `json:"photo_new"`
	PhotoCommentNew int    `json:"photo_comment_new"`
	VideoNew        int    `json:"video_new"`
	VideoCommentNew int    `json:"video_comment_new"`
	BoardPostNew    int    `json:"board_post_new"`
}

func GetLongPollSettings(accessToken string, wardVkId int) LongPollApiSettings {

	var groupID int
	if wardVkId < 0 {
		groupID = wardVkId * -1
	} else {
		groupID = wardVkId
	}

	values := url.Values{
		"access_token": {accessToken},
		"group_id":     {strconv.Itoa(groupID)},
		"v":            {"5.126"},
	}
	response, err := callMethod("groups.getLongPollSettings", values)
	if err != nil {
		tools.WriteToLog(err, debug.Stack())
		panic(err.Error())
	}

	lpApiSettings := parseSettingsForLongPollAPI(response)
	return lpApiSettings
}

func parseSettingsForLongPollAPI(response []byte) LongPollApiSettings {
	var lpApiSettings LongPollApiSettings

	var f interface{}
	err := json.Unmarshal(response, &f)
	if err != nil {
		tools.WriteToLog(err, debug.Stack())
		panic(err.Error())
	}

	s := f.(map[string]interface{})

	lpApiSettings.IsEnabled = s["is_enabled"].(bool)
	lpApiSettings.ApiVersion = s["api_version"].(string)
	lpApiSettings.WallPostNew = int(s["events"].(map[string]interface{})["wall_post_new"].(float64))
	lpApiSettings.WallReplyNew = int(s["events"].(map[string]interface{})["wall_reply_new"].(float64))
	lpApiSettings.PhotoNew = int(s["events"].(map[string]interface{})["photo_new"].(float64))
	lpApiSettings.PhotoCommentNew = int(s["events"].(map[string]interface{})["photo_comment_new"].(float64))
	lpApiSettings.VideoNew = int(s["events"].(map[string]interface{})["video_new"].(float64))
	lpApiSettings.VideoCommentNew = int(s["events"].(map[string]interface{})["video_comment_new"].(float64))
	lpApiSettings.BoardPostNew = int(s["events"].(map[string]interface{})["board_post_new"].(float64))

	return lpApiSettings
}

func SetLongPollSettings(accessToken string, wardVkId int, newParam map[string]string) {
	var groupID int
	if wardVkId < 0 {
		groupID = wardVkId * -1
	} else {
		groupID = wardVkId
	}

	values := url.Values{
		"access_token":         {accessToken},
		"group_id":             {strconv.Itoa(groupID)},
		"api_version":          {"5.126"},
		newParam["event_name"]: {newParam["mode"]},
		"v":                    {"5.126"},
	}

	_, err := callMethod("groups.setLongPollSettings", values)
	if err != nil {
		tools.WriteToLog(err, debug.Stack())
		panic(err.Error())
	}
}

type vkMessage struct {
	PeerID, RandomID                                       int
	Header, Text, Link, Footer, Attachments, ContentSource string
}

func (m *vkMessage) sendMessage(accessToken string) error {
	text := m.makeTextForMessage()
	values := m.makeMessageValues(m.PeerID, m.RandomID,
		accessToken, text, m.Attachments)
	_, err := callMethod("messages.send", values)
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

func (m *vkMessage) makeTextForMessage() string {
	var text string
	text += m.Header + "\n\n"
	if len(m.Text) > 0 {
		if len(m.Text) > 800 {
			text += m.Text[:800] + "...\n[много текста]\n\n"
		} else {
			text += m.Text + "\n\n"
		}
	} else {
		text += "\n[текста нет]\n\n"
	}
	if len(m.Link) > 0 {
		text += m.Link + "\n"
	}
	text += m.Footer

	return text
}

func (m *vkMessage) makeMessageValues(peerID, randomID int,
	accessToken, text, attachments string) url.Values {
	values := url.Values{
		"access_token":   {accessToken},
		"peer_id":        {strconv.Itoa(peerID)},
		"random_id":      {strconv.Itoa(randomID)},
		"message":        {text},
		"attachment":     {attachments},
		"content_source": {m.ContentSource},
		"v":              {"5.126"},
	}
	return values
}

type groupInfo struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func getGroupInfo(accessToken string, groupVkID int) groupInfo {
	values := url.Values{
		"access_token": {accessToken},
		"group_ids":    {strconv.Itoa(-groupVkID)},
		"v":            {"5.126"},
	}

	response, err := callMethod("groups.getById", values)
	if err != nil {
		tools.WriteToLog(err, debug.Stack())
		panic(err.Error())
	}
	groupInfo := parseGroupInfo(response)
	return groupInfo
}

func parseGroupInfo(response []byte) groupInfo {
	var g groupInfo

	var f interface{}
	err := json.Unmarshal(response, &f)

	if err != nil {
		tools.WriteToLog(err, debug.Stack())
		panic(err.Error())
	}
	item := f.([]interface{})

	g.ID = int(item[0].(map[string]interface{})["id"].(float64))
	g.Name = item[0].(map[string]interface{})["name"].(string)

	return g
}

type userInfo struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func getUserInfo(accessToken string, userVkID int) userInfo {
	values := url.Values{
		"access_token": {accessToken},
		"user_ids":     {strconv.Itoa(userVkID)},
		"v":            {"5.126"},
	}
	response, err := callMethod("users.get", values)
	if err != nil {
		tools.WriteToLog(err, debug.Stack())
		panic(err.Error())
	}
	userInfo := parseUserInfo(response)
	return userInfo
}

func parseUserInfo(response []byte) userInfo {
	var u userInfo

	var f interface{}
	err := json.Unmarshal(response, &f)
	if err != nil {
		tools.WriteToLog(err, debug.Stack())
		panic(err.Error())
	}
	item := f.([]interface{})

	u.ID = int(item[0].(map[string]interface{})["id"].(float64))
	u.FirstName = item[0].(map[string]interface{})["first_name"].(string)
	u.LastName = item[0].(map[string]interface{})["last_name"].(string)

	return u
}

func callMethod(methodName string, values url.Values) ([]byte, error) {
	response, err := govkapi.Method(methodName, values)
	if err != nil {
		switch true {
		case strings.Contains(strings.ToLower(err.Error()), "too many requests per second"):
			time.Sleep(340 * time.Millisecond)
			return callMethod(methodName, values)
		case strings.Contains(strings.ToLower(err.Error()), "too much messages sent to user"):
			return nil, err
		default:
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
	}
	return response, nil
}
