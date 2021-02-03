package vkapi

import (
	"encoding/json"
	"fmt"
	"github.com/VitJRBOG/GroupsObserver/tools"
	"io/ioutil"
	"net/http"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
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
	respLPS := parseResponseLongPollServer(response)
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

func parseResponseLongPollServer(response []byte) ResponseLongPollServer {
	var respLPS ResponseLongPollServer
	err := json.Unmarshal(response, &respLPS)
	if err != nil {
		tools.WriteToLog(err, debug.Stack())
		panic(err.Error())
	}
	return respLPS
}
