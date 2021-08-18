package observer

import (
	"fmt"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/VitJRBOG/Watcher/db"
	"github.com/VitJRBOG/Watcher/tools"
	"github.com/VitJRBOG/Watcher/vkapi"
)

type ModuleParams struct {
	Name      string
	Status    string
	Message   chan string
	BrakeFlag bool
	Ward      db.Ward
}

func MakeObservers() []*ModuleParams {
	var params []*ModuleParams
	wards := db.SelectWards()
	for _, ward := range wards {
		var accessToken db.AccessToken
		accessToken.SelectByID(ward.GetAccessTokenID)

		var p ModuleParams
		p.Name = fmt.Sprintf("%s observer", ward.Name)
		p.Status = "inactive"
		p.Ward = ward
		p.Message = make(chan string)
		params = append(params, &p)
	}
	return params
}

func StartObserver(params *ModuleParams) {
	// TODO: проверка, не был ли подопечный удален из БД
	afterTurningOff := true
	for true {
		if params.Ward.UnderObservation == 1 {
			params.Status = "active"
			if afterTurningOff {
				params.Message <- "It begins..."
				afterTurningOff = false
			}
			var accessToken db.AccessToken
			accessToken.SelectByID(params.Ward.GetAccessTokenID)

			longPollApiIsEnabled := checkLongPollAPIEnabling(accessToken.Value, params.Ward.VkID)
			if !longPollApiIsEnabled {
				vkapi.TurnOnLongPollAPI(accessToken.Value, params.Ward.VkID)
			}

			respLPS := vkapi.ListenLongPollServer(accessToken.Value, -(params.Ward.VkID), params.Ward.LastTS)
			params.Status = "processing"
			err := parseLongPollServerResponse(respLPS, &params.Ward)
			if err != nil {
				if strings.Contains(strings.ToLower(err.Error()), "too much messages sent to user") {
					params.Status = "waiting for 5 minutes"
					params.Message <- fmt.Sprintf("ERROR: «%s»", err.Error())
					time.Sleep(5 * time.Minute)
					params.Message <- "Let's get back to work..."
				} else {
					tools.WriteToLog(err, debug.Stack())
					panic(err.Error())
				}
			}
			if params.BrakeFlag {
				params.Status = "stopped"
				params.Message <- "Was stopped by user"
				break
			}
		} else {
			params.Status = "inactive"
			if !(afterTurningOff) {
				params.Message <- "My journey has reached a close..."
			}
			time.Sleep(5 * time.Second)
			afterTurningOff = true
		}
		params.Ward.SelectByID(params.Ward.ID)
	}
}

func parseLongPollServerResponse(respLPS vkapi.ResponseLongPollServer, ward *db.Ward) error {
	if len(respLPS.Updates) > 0 {
		for _, update := range respLPS.Updates {
			var err error

			switch update.Type {
			case "wall_post_new":
				err = handleWallPostUpdate(ward, update)
			case "photo_new":
				err = handlePhotoUpdate(ward, update)
			case "photo_comment_new":
				err = handlePhotoCommentUpdate(ward, update)
			case "video_new":
				err = handleVideoUpdate(ward, update)
			case "video_comment_new":
				err = handleVideoCommentUpdate(ward, update)
			case "board_post_new":
				err = handleBoardPostUpdate(ward, update)
			case "wall_reply_new":
				err = handleWallReplyUpdate(ward, update)
			}

			if err != nil {
				if strings.Contains(strings.ToLower(err.Error()),
					"too much messages sent to user") {
					return err
				} else {
					tools.WriteToLog(err, debug.Stack())
					panic(err.Error())
				}
			}
		}
		updateWard(ward, respLPS.TS)
	}
	return nil
}

func handleWallPostUpdate(ward *db.Ward, update vkapi.UpdateFromLongPollServer) error {
	var wallPost vkapi.WallPost
	wallPost.ParseData(update)
	targetWallPostType := checkWallPostType(ward.ID, wallPost)
	if targetWallPostType {
		getAT, sendAT, operatorVkID := getDataByCurrentObserver("wall_post", ward.ID)
		err := wallPost.SendWithMessage(getAT, sendAT, operatorVkID)
		if err != nil {
			if strings.Contains(strings.ToLower(err.Error()),
				"too much messages sent to user") {
				return err
			} else {
				tools.WriteToLog(err, debug.Stack())
				panic(err.Error())
			}
		}
	}
	return nil
}

func handlePhotoUpdate(ward *db.Ward, update vkapi.UpdateFromLongPollServer) error {
	var photo vkapi.Photo
	photo.ParseData(update)
	getAT, sendAT, operatorVkID := getDataByCurrentObserver("photo", ward.ID)
	err := photo.SendWithMessage(getAT, sendAT, operatorVkID)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()),
			"too much messages sent to user") {
			return err
		} else {
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
	}
	return nil
}

func handlePhotoCommentUpdate(ward *db.Ward, update vkapi.UpdateFromLongPollServer) error {
	var photoComment vkapi.PhotoComment
	photoComment.ParseData(update)
	getAT, sendAT, operatorVkID := getDataByCurrentObserver("photo_comment", ward.ID)
	err := photoComment.SendWithMessage(getAT, sendAT, operatorVkID)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()),
			"too much messages sent to user") {
			return err
		} else {
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
	}
	return nil
}

func handleVideoUpdate(ward *db.Ward, update vkapi.UpdateFromLongPollServer) error {
	var video vkapi.Video
	video.ParseData(update)
	getAT, sendAT, operatorVkID := getDataByCurrentObserver("video", ward.ID)
	err := video.SendWithMessage(getAT, sendAT, operatorVkID)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()),
			"too much messages sent to user") {
			return err
		} else {
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
	}
	return nil
}

func handleVideoCommentUpdate(ward *db.Ward, update vkapi.UpdateFromLongPollServer) error {
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
	return nil
}

func handleBoardPostUpdate(ward *db.Ward, update vkapi.UpdateFromLongPollServer) error {
	var boardPost vkapi.BoardPost
	boardPost.ParseData(update)
	getAT, sendAT, operatorVkID := getDataByCurrentObserver("board_post", ward.ID)
	err := boardPost.SendWithMessage(getAT, sendAT, operatorVkID)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()),
			"too much messages sent to user") {
			return err
		} else {
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
	}
	return nil
}

func handleWallReplyUpdate(ward *db.Ward, update vkapi.UpdateFromLongPollServer) error {
	var wallReply vkapi.WallReply
	wallReply.ParseData(update)
	getAT, sendAT, operatorVkID := getDataByCurrentObserver("wall_reply", ward.ID)
	err := wallReply.SendWithMessage(getAT, sendAT, operatorVkID)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()),
			"too much messages sent to user") {
			return err
		} else {
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
	}
	return nil
}

func checkLongPollAPIEnabling(accessToken string, wardVkId int) bool {
	lpApiSettings := vkapi.GetLongPollSettings(accessToken, wardVkId)
	return lpApiSettings.IsEnabled
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
	ward.UpdateInDB()
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
