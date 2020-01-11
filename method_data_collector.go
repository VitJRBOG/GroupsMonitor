package main

import (
	"fmt"
	"strconv"
)

// CollectionMethodData
//собирает данные обо всех используемых в данном мониторе методах vk api
func CollectionMethodData(monitor *Monitor) (*[]Method, error) {
	var methods []Method

	methodsNames := getMethodsNames(monitor.Name)

	for _, methodName := range methodsNames {
		var method Method
		accessTokenID, err := getAccessTokenID(monitor.Name, methodName)
		if err != nil {
			return nil, err
		}
		method.AccessTokenID = accessTokenID
		method.Name = methodName
		method.SubjectID = monitor.SubjectID
		method.MonitorID = monitor.ID

		methods = append(methods, method)
	}

	return &methods, nil
}

// getMethodsNames получает список названий необходимых методов vk api
func getMethodsNames(monitorName string) []string {
	var methodsNames []string

	switch monitorName {
	case "wall_post_monitor":
		methodsNames = append(methodsNames, "wall.get", "users.get",
			"groups.getById", "messages.send")
	case "album_photo_monitor":
		methodsNames = append(methodsNames, "photos.get",
			"photos.getAlbums", "users.get", "groups.getById", "messages.send")
	case "video_monitor":
		methodsNames = append(methodsNames, "video.get", "users.get",
			"groups.getById", "messages.send")
	case "photo_comment_monitor":
		methodsNames = append(methodsNames, "photos.getAllComments",
			"users.get", "groups.getById", "messages.send")
	case "video_comment_monitor":
		methodsNames = append(methodsNames, "video.get",
			"video.getComments", "users.get", "groups.getById", "messages.send")
	case "topic_monitor":
		methodsNames = append(methodsNames, "board.getTopics",
			"board.getComments", "users.get", "groups.getById", "messages.send")
	case "wall_post_comment_monitor":
		methodsNames = append(methodsNames, "wall.get", "wall.getComments",
			"wall.getComment", "users.get", "groups.getById", "messages.send")
	}

	return methodsNames
}

// getAccessTokenID
//запрашивает и пользователя идентификатор в БД необходимого токена доступа для конкретного метода vk api
func getAccessTokenID(monitorName string, methodName string) (int, error) {
	// получаем из БД данные обо всех токенах доступа
	accessTokens, err := SelectDBAccessTokens()
	if err != nil {
		return 0, err
	}

	// перебираем токены доступа и отображаем их данные для пользователя
	for _, accessToken := range accessTokens {
		fmt.Printf("> [Access token]: name = %v, ID = %d\n", accessToken.Name,
			accessToken.ID)
	}

	// запрашиваем у пользователя идентификатор необходимого токена доступа для данного метода vk api
	fmt.Printf("> [Get access token's ID for \"%v\" of \"%v\"]: ", methodName,
		monitorName)
	var userAnswer string
	_, err = fmt.Scan(&userAnswer)
	if err != nil {
		return 0, err
	}
	accessTokenID, err := strconv.Atoi(userAnswer)
	if err != nil {
		return 0, err
	}
	return accessTokenID, nil
}
