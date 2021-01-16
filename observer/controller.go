package observer

import (
	"fmt"
	"github.com/VitJRBOG/GroupsMonitor_new/db"
	"github.com/VitJRBOG/GroupsMonitor_new/tools"
	"github.com/VitJRBOG/GroupsMonitor_new/vkapi"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
)

type ControllerParams struct {
	Name      string
	Status    string
	Message   string
	BrakeFlag bool
}

func RunControllerObservers() []*ControllerParams {
	var listControllerParams []*ControllerParams
	wards := db.SelectWards()
	for _, ward := range wards {
		var accessToken db.AccessToken
		accessToken.SelectByID(ward.GetAccessTokenID)

		var c ControllerParams
		go startObserver(&c, ward)
		listControllerParams = append(listControllerParams, &c)
	}
	return listControllerParams
}

func startObserver(ctrlParams *ControllerParams, ward db.Ward) {
	ctrlParams.Name = fmt.Sprintf("%s observer", ward.Name)
	for true {
		ctrlParams.Status = "active"
		var accessToken db.AccessToken
		accessToken.SelectByID(ward.GetAccessTokenID)
		respLPS := vkapi.ListenLongPollServer(accessToken.Value, -(ward.VkID), ward.LastTS)
		err := parseLongPollServerResponse(respLPS, &ward)
		if err != nil {
			if strings.Contains(strings.ToLower(err.Error()), "too much messages sent to user") {
				ctrlParams.Message = fmt.Sprintf("ERROR: «%s»", err.Error())
				ctrlParams.Status = "waiting for 5 minutes"
				time.Sleep(5 * time.Minute)
				ctrlParams.Message = "Let's get back to work..."
			} else {
				tools.WriteToLog(err, debug.Stack())
				panic(err.Error())
			}
		}
		if ctrlParams.BrakeFlag {
			ctrlParams.Message = "Was stopped by user"
			ctrlParams.Status = "stopped"
			break
		}
	}
}

func parseLongPollServerResponse(respLPS vkapi.ResponseLongPollServer, ward *db.Ward) error {
	if len(respLPS.Updates) > 0 {
		for _, update := range respLPS.Updates {
			switch update.Type {
			case "wall_post_new":
				var wallPost vkapi.WallPost
				wallPost.ParseData(update)
				targetWallPostType := checkWallPostType(ward.ID, wallPost)
				if targetWallPostType {
					getAT, sendAT, operatorVkID := getDataByCurrentObserver("wall_post", ward.ID)
					err := wallPost.SendWithMessage(getAT, sendAT, operatorVkID)
					if err != nil {
						if strings.Contains(strings.ToLower(err.Error()), "too much messages sent to user") {
							return err
						} else {
							tools.WriteToLog(err, debug.Stack())
							panic(err.Error())
						}
					}
				}
			case "wall_reply_new":
				var wallReply vkapi.WallReply
				wallReply.ParseData(update)
				getAT, sendAT, operatorVkID := getDataByCurrentObserver("wall_reply", ward.ID)
				err := wallReply.SendWithMessage(getAT, sendAT, operatorVkID)
				if err != nil {
					if strings.Contains(strings.ToLower(err.Error()), "too much messages sent to user") {
						return err
					} else {
						tools.WriteToLog(err, debug.Stack())
						panic(err.Error())
					}
				}
			case "photo_new":
				var photo vkapi.Photo
				photo.ParseData(update)
				getAT, sendAT, operatorVkID := getDataByCurrentObserver("photo", ward.ID)
				err := photo.SendWithMessage(getAT, sendAT, operatorVkID)
				if err != nil {
					if strings.Contains(strings.ToLower(err.Error()), "too much messages sent to user") {
						return err
					} else {
						tools.WriteToLog(err, debug.Stack())
						panic(err.Error())
					}
				}
			case "photo_comment_new":
				var photoComment vkapi.PhotoComment
				photoComment.ParseData(update)
				getAT, sendAT, operatorVkID := getDataByCurrentObserver("photo_comment", ward.ID)
				err := photoComment.SendWithMessage(getAT, sendAT, operatorVkID)
				if err != nil {
					if strings.Contains(strings.ToLower(err.Error()), "too much messages sent to user") {
						return err
					} else {
						tools.WriteToLog(err, debug.Stack())
						panic(err.Error())
					}
				}
			case "video_new":
				var video vkapi.Video
				video.ParseData(update)
				getAT, sendAT, operatorVkID := getDataByCurrentObserver("video", ward.ID)
				err := video.SendWithMessage(getAT, sendAT, operatorVkID)
				if err != nil {
					if strings.Contains(strings.ToLower(err.Error()), "too much messages sent to user") {
						return err
					} else {
						tools.WriteToLog(err, debug.Stack())
						panic(err.Error())
					}
				}
			case "video_comment_new":
				var videoComment vkapi.VideoComment
				videoComment.ParseData(update)
				getAT, sendAT, operatorVkID := getDataByCurrentObserver("video_comment", ward.ID)
				err := videoComment.SendWithMessage(getAT, sendAT, operatorVkID)
				if err != nil {
					if strings.Contains(strings.ToLower(err.Error()), "too much messages sent to user") {
						return err
					} else {
						tools.WriteToLog(err, debug.Stack())
						panic(err.Error())
					}
				}
			case "board_post_new":
				var boardPost vkapi.BoardPost
				boardPost.ParseData(update)
				getAT, sendAT, operatorVkID := getDataByCurrentObserver("board_post", ward.ID)
				err := boardPost.SendWithMessage(getAT, sendAT, operatorVkID)
				if err != nil {
					if strings.Contains(strings.ToLower(err.Error()), "too much messages sent to user") {
						return err
					} else {
						tools.WriteToLog(err, debug.Stack())
						panic(err.Error())
					}
				}
			}
		}
		updateWard(ward, respLPS.TS)
	}
	return nil
}

func checkWallPostType(wardID int, w vkapi.WallPost) bool {
	var observer db.Observer
	observer.SelectByNameAndWardID("wall_post", wardID)

	if w.PostType == observer.AdditionalParams.WallPost.PostType {
		return true
	}
	return false
}

func updateWard(ward *db.Ward, lastTS string) {
	var err error
	ward.LastTS, err = strconv.Atoi(lastTS)
	if err != nil {
		tools.WriteToLog(err, debug.Stack())
		panic(err.Error())
	}
	ward.UpdateDB()
}

func getDataByCurrentObserver(observerName string, wardID int) (string, string, int) {
	var observer db.Observer
	observer.SelectByNameAndWardID(observerName, wardID)

	var ward db.Ward
	ward.SelectByID(wardID)

	var getAccessToken db.AccessToken
	getAccessToken.SelectByID(ward.GetAccessTokenID)

	var sendAccessToken db.AccessToken
	sendAccessToken.SelectByID(observer.SendAccessTokenID)

	var operator db.Operator
	operator.SelectByID(observer.OperatorID)

	return getAccessToken.Value, sendAccessToken.Value, operator.VkID
}
