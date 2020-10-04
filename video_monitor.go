package main

import (
	"fmt"
	"strings"
)

// VideoMonitor проверяет видео
func VideoMonitor(subject Subject) (*VideoMonitorParam, error) {
	sender := fmt.Sprintf("%v's video monitoring", subject.Name)

	// запрашиваем структуру с параметрами модуля мониторинга видео
	var videoMonitorParam VideoMonitorParam
	err := videoMonitorParam.selectFromDBBySubjectID(subject.ID)
	if err != nil {
		return nil, err
	}

	// запрашиваем структуру с видео из альбомов субъекта
	videos, err := getVideos(sender, subject, videoMonitorParam)
	if err != nil {
		return nil, err
	}

	var targetVideos []Video

	// отфильтровываем старые
	for _, video := range videos {
		if video.Date > videoMonitorParam.LastDate {
			targetVideos = append(targetVideos, video)
		}
	}

	// если после фильтрации что-то осталось, продолжаем
	if len(targetVideos) > 0 {
		// сортируем в порядке от раннего к позднему
		for j := 0; j < len(targetVideos); j++ {
			f := 0
			for i := 0; i < len(targetVideos)-1-j; i++ {
				if targetVideos[i].Date > targetVideos[i+1].Date {
					x := targetVideos[i]
					y := targetVideos[i+1]
					targetVideos[i+1] = x
					targetVideos[i] = y
					f = 1
				}
			}
			if f == 0 {
				break
			}
		}

		// перебираем отсортированный список
		for _, video := range targetVideos {
			// формируем строку с данными для карты для отправки сообщения
			messageParameters, err := makeMessageVideo(sender, subject, videoMonitorParam, video)
			if err != nil {
				return nil, err
			}

			// отправляем сообщение с полученными данными
			if err := SendMessage(sender, messageParameters, subject); err != nil {
				return nil, err
			}

			// выводим в консоль сообщение о новом видео
			outputReportAboutNewVideo(sender, video)

			// обновляем дату последнего проверенного поста в БД
			if err := videoMonitorParam.updateInDBFieldLastDate(subject.ID, video.Date); err != nil {
				return nil, err
			}
		}
	}

	return &videoMonitorParam, nil
}

// outputReportAboutNewVideo выводит сообщение о новом видео
func outputReportAboutNewVideo(sender string, video Video) {
	creationDate := UnixTimeStampToDate(video.Date)
	message := fmt.Sprintf("New video at %v.", creationDate)
	OutputMessage(sender, message)
}

// Video - структура для данных о видео
type Video struct {
	ID          int    `json:"id"`
	OwnerID     int    `json:"owner_id"`
	FromID      int    `json:"user_id"`
	Description string `json:"description"`
	Date        int    `json:"date"`
}

// getVideos формирует запрос на получение видео и посылает его к vk api
func getVideos(sender string, subject Subject, videoMonitorParam VideoMonitorParam) ([]Video, error) {
	var videos []Video

	// формируем карту с параметрами запроса
	jsonDump := fmt.Sprintf(`{
			"owner_id": "%d",
			"count": "%d",
			"v": "5.95"
		}`, subject.SubjectID, videoMonitorParam.VideoCount)
	values, err := MakeJSON(jsonDump)
	if err != nil {
		return videos, err
	}

	// отправляем запрос, получаем ответ
	response, err := SendVKAPIQuery(sender, "video.get", values, subject)
	if err != nil {
		return videos, err
	}

	videos = ParseVideoVkAPIMap(response["response"].(map[string]interface{}))

	return videos, nil
}

// ParseVideoVkAPIMap извлекает данные о видео из полученной карты vk api
func ParseVideoVkAPIMap(resp map[string]interface{}) []Video {
	var videos []Video

	// перебираем элементы с данными о фотографиях
	itemsMap := resp["items"].([]interface{})
	for _, itemMap := range itemsMap {
		item := itemMap.(map[string]interface{})

		// парсим данные из элемента в структуру
		var video Video
		video.ID = int(item["id"].(float64))
		video.OwnerID = int(item["owner_id"].(float64))
		if fromID, exist := item["user_id"]; exist == true {
			video.FromID = int(fromID.(float64))
		}
		video.Date = int(item["date"].(float64))
		video.Description = item["description"].(string)

		videos = append(videos, video)
	}

	return videos
}

// makeMessageVideo собирает сообщение с данными о видео
func makeMessageVideo(sender string, subject Subject,
	videoMonitorParam VideoMonitorParam, video Video) (string, error) {
	// собираем данные для сигнатуры сообщения:

	// где находится
	var locationHyperlink string
	if subject.SubjectID < 0 {
		vkCommunity, err := GetCommunityInfo(sender, subject, subject.SubjectID)
		locationHyperlink = MakeCommunityHyperlink(vkCommunity)
		if err != nil {
			return "", err
		}
	} else {
		vkUser, err := GetUserInfo(sender, subject, subject.SubjectID)
		locationHyperlink = MakeUserHyperlink(vkUser)
		if err != nil {
			return "", err
		}
	}

	// кто загрузил
	var authorHyperlink string
	if video.FromID == 0 {
		authorHyperlink = locationHyperlink
	} else {
		vkUser, err := GetUserInfo(sender, subject, video.FromID)
		authorHyperlink = MakeUserHyperlink(vkUser)
		if err != nil {
			return "", err
		}
	}

	// и когда загрузил
	creationDate := UnixTimeStampToDate(video.Date)

	// формируем строку с прикреплением
	attachment := fmt.Sprintf("video%d_%d", video.OwnerID, video.ID)

	// собираем ссылку на видео
	videoURL := fmt.Sprintf("https://vk.com/video%d_%d", video.OwnerID, video.ID)

	// добавляем подготовленные фрагменты сообщения в общий текст
	// сначала сигнатуру
	text := fmt.Sprintf("New video\\nLocation: %v\\nAuthor: %v\\nCreated: %v",
		locationHyperlink, authorHyperlink, creationDate)

	// затем описание к видео, если оно есть
	if len(video.Description) > 0 {
		// но сначала обрезаем его из-за ограничения на длину запроса
		if len(video.Description) > 800 {
			video.Description = string(video.Description[0:800]) + "\\n[long_text]"
		}
		// и экранируем все символы пропуска строки, потому что у json.Unmarshal с ними проблемы
		video.Description = strings.Replace(video.Description, "\n", "\\n", -1)
		text += fmt.Sprintf("\\n\\n%v", video.Description)
	}

	// и добавляем ссылку на само видео
	text += fmt.Sprintf("\\n\\n%v", videoURL)

	// экранируем все обратные слэши, не сломали json.Unmarshal
	text = strings.Replace(text, `\`, `\\`, -1)
	// экранируем все апострофы, чтобы не сломали нам json.Unmarshal
	text = strings.Replace(text, `"`, `\"`, -1)
	// и возвращаем символы пропуска строки после экранировки обратных слэшей
	text = strings.Replace(text, `\\n`, `\n`, -1)

	// далее формируем строку с данными для карты
	jsonDump := fmt.Sprintf(`{
		"peer_id": "%d",
		"message": "%v",
		"attachment": "%v",
		"v": "5.68"
	}`, videoMonitorParam.SendTo, text, attachment)

	return jsonDump, nil
}
