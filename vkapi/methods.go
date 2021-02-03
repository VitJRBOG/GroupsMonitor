package vkapi

import (
	"encoding/json"
	govkapi "github.com/VitJRBOG/GoVkApi"
	"github.com/VitJRBOG/GroupsObserver/tools"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
)

type longPollServerConnectionData struct {
	Key    string `json:"key"`
	Server string `json:"server"`
	TS     string `json:"ts"`
}

func getLongPollServerConnectionData(accessToken string, wardVkId int) longPollServerConnectionData {
	values := map[string]string{
		"group_id": strconv.Itoa(wardVkId),
		"v":        "5.126",
	}

	response, err := callMethod("groups.getLongPollServer", accessToken, values)
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

type vkMessage struct {
	PeerID, RandomID                        int
	Header, Text, Link, Footer, Attachments string
}

func (m *vkMessage) sendMessage(accessToken string) error {
	text := m.makeTextForMessage()
	values := m.makeMessageValues(m.PeerID, m.RandomID, text, m.Attachments)
	// FIXME: не прикрепляет картинки
	_, err := callMethod("messages.send", accessToken, values)
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

func (m *vkMessage) makeMessageValues(peerID, randomID int, text, attachments string) map[string]string {
	values := map[string]string{
		"peer_id":     strconv.Itoa(peerID),
		"random_id":   strconv.Itoa(randomID),
		"message":     text,
		"attachments": attachments,
		"v":           "5.126",
	}
	return values
}

type groupInfo struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func getGroupInfo(accessToken string, groupVkID int) groupInfo {
	values := map[string]string{
		"group_ids": strconv.Itoa(-groupVkID),
		"v":         "5.126",
	}

	response, err := callMethod("groups.getById", accessToken, values)
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
	values := map[string]string{
		"user_ids": strconv.Itoa(userVkID),
		"v":        "5.126",
	}
	response, err := callMethod("users.get", accessToken, values)
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

func callMethod(methodName, accessToken string, values map[string]string) ([]byte, error) {
	response, err := govkapi.SendRequestVkApi(accessToken, methodName, values)
	if err != nil {
		switch true {
		case strings.Contains(strings.ToLower(err.Error()), "too many requests per second"):
			time.Sleep(340 * time.Millisecond)
			return callMethod(methodName, accessToken, values)
		case strings.Contains(strings.ToLower(err.Error()), "too much messages sent to user"):
			return nil, err
		default:
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
	}
	return response, nil
}
