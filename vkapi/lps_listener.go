package vkapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/VitJRBOG/GroupsObserver/tools"
)

type ResponseLongPollServer struct {
	TS      string                     `json:"ts"`
	Updates []UpdateFromLongPollServer `json:"updates"`
}

type UpdateFromLongPollServer struct {
	Type   string                 `json:"type"`
	Object map[string]interface{} `json:"object"`
}

func ListenLongPollServer(accessToken string, wardVkID, lastTS int) ResponseLongPollServer {
	lpsConnectionData := getLongPollServerConnectionData(accessToken, wardVkID)
	queryForLongPollServer := makeQueryToLongPollServer(lpsConnectionData, lastTS)
	response := sendRequestToLongPollServer(queryForLongPollServer)
	respLPS, err := parseResponseLongPollServer(response)
	if err != nil {
		switch true {
		case strings.Contains(strings.ToLower(err.Error()),
			"the event history is out of date or has been partially lost"):
			return ListenLongPollServer(accessToken, wardVkID, lastTS+1)
		case strings.Contains(strings.ToLower(err.Error()), "key expired"):
			return ListenLongPollServer(accessToken, wardVkID, lastTS)
		case strings.Contains(strings.ToLower(err.Error()), "information lost"):
			return ListenLongPollServer(accessToken, wardVkID, lastTS)
		default:
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
	}
	return respLPS
}

func makeQueryToLongPollServer(lpsConnectionData longPollServerConnectionData, lastTS int) string {
	url := fmt.Sprintf("%s?act=a_check&key=%s&wait=%d&ts=%s",
		lpsConnectionData.Server, lpsConnectionData.Key, 5, strconv.Itoa(lastTS))
	return url
}

func sendRequestToLongPollServer(url string) []byte {
	response, err := http.Get(url)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "connection reset by peer") {
			time.Sleep(5 * time.Second)
			return sendRequestToLongPollServer(url)
		}
		tools.WriteToLog(err, debug.Stack())
		panic(err.Error())
	}
	defer func() {
		err := response.Body.Close()
		if err != nil {
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
	}()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		tools.WriteToLog(err, debug.Stack())
		panic(err.Error())
	}
	return body
}

func parseResponseLongPollServer(response []byte) (ResponseLongPollServer, error) {
	var respValues map[string]interface{}

	err := json.Unmarshal(response, &respValues)
	if err != nil {
		tools.WriteToLog(err, debug.Stack())
		panic(err.Error())
	}

	var respLPS ResponseLongPollServer
	if updates, exist := respValues["updates"]; exist {
		respLPS.TS = respValues["ts"].(string)

		for _, o := range updates.([]interface{}) {
			item := o.(map[string]interface{})

			var u UpdateFromLongPollServer
			u.Type = item["type"].(string)
			u.Object = item["object"].(map[string]interface{})

			respLPS.Updates = append(respLPS.Updates, u)
		}
	} else {
		if errorCode, exist := respValues["failed"]; exist {
			var errorMsg string
			switch errorCode.(float64) {
			case 1:
				errorMsg = "the event history is out of date or has been partially lost"
			case 2:
				errorMsg = "key expired"
			case 3:
				errorMsg = "information lost"
			default:
				errorMsg = "unknown error"
			}

			err = fmt.Errorf("code %v: %s", errorCode.(float64), errorMsg)
			return respLPS, err
		}
		tools.WriteToLog(err, debug.Stack())
		panic(err.Error())
	}

	return respLPS, nil
}
