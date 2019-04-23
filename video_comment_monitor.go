package main

import (
	"fmt"
	"strings"
)

// VideoCommentMonitor проверяет комментарии под видео
func VideoCommentMonitor(subject Subject) error {
	sender := fmt.Sprintf("%v's video comment monitoring", subject.Name)

	// запрашиваем структуру с параметрами модуля мониторинга фотографий
	videoCommentMonitorParam, err := SelectDBVideoCommentMonitorParam(subject.ID)
	if err != nil {
		return err
	}

	// запрашиваем структуру с видеозаписями
	videos, err := getVideosForComments(sender, subject, videoCommentMonitorParam)
	if err != nil {
		return err
	}

	// запрашиваем комментарии из этих видео
	videoComments, err := getVideoComments(sender, subject, videoCommentMonitorParam, videos)
	if err != nil {
		return err
	}

	var targetVideoComments []VideoComment

	// отфильтровываем старые
	for _, videoComment := range videoComments {
		if videoComment.Date > videoCommentMonitorParam.LastDate {
			targetVideoComments = append(targetVideoComments, videoComment)
		}
	}

	// если после фильтрации что-то осталось, продолжаем
	if len(targetVideoComments) > 0 {

		// сортируем в порядке от раннего к позднему
		for j := 0; j < len(targetVideoComments); j++ {
			f := 0
			for i := 0; i < len(targetVideoComments)-1-j; i++ {
				if targetVideoComments[i].Date > targetVideoComments[i+1].Date {
					x := targetVideoComments[i]
					y := targetVideoComments[i+1]
					targetVideoComments[i+1] = x
					targetVideoComments[i] = y
					f = 1
				}
			}
			if f == 0 {
				break
			}
		}

		// перебираем отсортированный список
		for _, videoComment := range targetVideoComments {

			// формируем строку с данными для карты для отправки сообщения
			messageParameters, err := makeMessageVideoComment(sender, subject,
				videoCommentMonitorParam, videoComment)
			if err != nil {
				return err
			}

			// отправляем сообщение с полученными данными
			if err := SendMessage(sender, messageParameters, subject); err != nil {
				return err
			}

			// выводим в консоль сообщение о новой фотографии
			outputReportAboutNewVideoComment(sender, videoComment)

			// обновляем дату последнего проверенного поста в БД
			if err := UpdateDBVideoCommentMonitorLastDate(subject.ID, videoComment.Date); err != nil {
				return err
			}
		}
	}

	return nil
}

// outputReportAboutNewVideoComment выводит сообщение о новом комментарии
func outputReportAboutNewVideoComment(sender string, videoComment VideoComment) {
	creationDate := UnixTimeStampToDate(videoComment.Date + 18000) // цифру изменить под свой часовой пояс (тут 5 часов)
	message := fmt.Sprintf("New comment under video at %v.", creationDate)
	OutputMessage(sender, message)
}

// Video
// структура уже описана в модуле мониторинга новых видео,
// поэтому тут этот момент опустим

// getVideos формирует запрос на получение видео и посылает его к vk api
func getVideosForComments(sender string, subject Subject,
	videoCommentMonitorParam VideoCommentMonitorParam) ([]Video, error) {
	var videos []Video

	// формируем карту с параметрами запроса
	jsonDump := fmt.Sprintf(`{
			"owner_id": "%d",
			"count": "%d",
			"v": "5.95"
		}`, subject.SubjectID, videoCommentMonitorParam.VideosCount)
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

// ParseVideoVkAPIMap
// функция парсинга карты с данными о видео уже описана в модуле мониторинга новых видео,
// поэтому тут этот момент опустим

// VideoComment - структура для данных о комментарии под видео
type VideoComment struct {
	ID           int          `json:"id"`
	VideoOwnerID int          `json:"owner_id"`
	VideoID      int          `json:"vid"`
	FromID       int          `json:"from_id"`
	Date         int          `json:"date"`
	Text         string       `json:"text"`
	Attachments  []Attachment `json:"attachments"`
}

// getVideoComments формирует запрос на получение комментариев под видео и посылает его к vk api
func getVideoComments(sender string, subject Subject,
	videoCommentMonitorParam VideoCommentMonitorParam, videos []Video) ([]VideoComment, error) {
	var videosComments []VideoComment

	for _, video := range videos {

		// формируем карту с параметрами запроса
		jsonDump := fmt.Sprintf(`{
			"owner_id": "%d",
			"video_id": "%d",
			"count": "%d",
			"v": "5.95"
		}`, video.OwnerID, video.ID, videoCommentMonitorParam.CommentsCount)
		values, err := MakeJSON(jsonDump)
		if err != nil {
			return videosComments, err
		}

		// отправляем запрос, получаем ответ
		response, err := SendVKAPIQuery(sender, "video.getComments", values, subject)
		if err != nil {
			return videosComments, err
		}

		videoComments := parseVideoCommentVkAPIMap(response["response"].(map[string]interface{}), video)
		videosComments = append(videosComments, videoComments...)
	}

	return videosComments, nil
}

// parseVideoCommentVkAPIMap извлекает данные о комментариях из полученной карты vk api
func parseVideoCommentVkAPIMap(resp map[string]interface{}, video Video) []VideoComment {
	var videoComments []VideoComment

	// перебираем элементы с данными о комментариях
	itemsMap := resp["items"].([]interface{})
	for _, itemMap := range itemsMap {
		item := itemMap.(map[string]interface{})

		// парсим данные из элемента в структуру
		var videoComment VideoComment
		videoComment.ID = int(item["id"].(float64))
		videoComment.VideoID = video.ID
		videoComment.VideoOwnerID = video.OwnerID
		videoComment.FromID = int(item["from_id"].(float64))
		videoComment.Date = int(item["date"].(float64))
		videoComment.Text = item["text"].(string)

		// если есть прикрепления, то вызываем парсер прикреплений
		if mediaContent, exist := item["attachments"]; exist == true {
			videoComment.Attachments = ParseAttachments(mediaContent.([]interface{}))
		}
		videoComments = append(videoComments, videoComment)
	}
	return videoComments
}

// makeMessageAlbumPhoto собирает сообщение с данными о комментарии
func makeMessageVideoComment(sender string, subject Subject,
	videoCommentMonitorParam VideoCommentMonitorParam, videoComment VideoComment) (string, error) {
	// собираем данные для сигнатуры сообщения:
	// где комментарий был обнаружен
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

	// кто автор
	var authorHyperlink string
	if videoComment.FromID < 0 {
		vkCommunity, err := GetCommunityInfo(sender, subject, videoComment.FromID)
		authorHyperlink = MakeCommunityHyperlink(vkCommunity)
		if err != nil {
			return "", err
		}
	} else {
		vkUser, err := GetUserInfo(sender, subject, videoComment.FromID)
		authorHyperlink = MakeUserHyperlink(vkUser)
		if err != nil {
			return "", err
		}
	}

	// и когда был создан
	creationDate := UnixTimeStampToDate(videoComment.Date + 18000) // цифру изменить под свой часовой пояс (тут 5 часов)

	// формируем строку с прикреплениями
	var attachments string
	var link string
	if len(videoComment.Attachments) > 0 {
		for _, attachment := range videoComment.Attachments {
			// отдельно проверим прикрепленную к посту ссылку
			if attachment.Type == "link" {
				link = attachment.URL
			} else { // если не ссылка, значит - всё остальное
				attachments = fmt.Sprintf("%v%d_%d", attachment.Type, attachment.OwnerID, attachment.ID)
				if len(attachment.AccessKey) > 0 {
					attachments += fmt.Sprintf("_%v", attachment.AccessKey)
				}
				// заранее добавляем запятую, чтобы следующее прикрепление не слилось к предыдущим
				attachments += ","
			}
		}

		// если кроме ссылки в прикреплениях было что-то еще, то нужно удалить лишнюю запятую в конце строки
		if len(attachments) > 0 {
			if string(attachments[len(attachments)-1]) == "," {
				attachments = string(attachments[0 : len(attachments)-1])
			}
		}
	}

	// собираем ссылку на видео, потому что ссылку на комментарий под видео сделать нельзя
	videoURL := fmt.Sprintf("https://vk.com/video%d_%d", videoComment.VideoOwnerID, videoComment.VideoID)

	// добавляем подготовленные фрагменты сообщения в общий текст
	// сначала сигнатуру
	text := fmt.Sprintf("New video comment\\nLocation: %v\\nAuthor: %v\\nCreated: %v",
		locationHyperlink, authorHyperlink, creationDate)

	// затем основной текст комментария, если он есть
	if len(videoComment.Text) > 0 {
		// но сначала экранируем все символы пропуска строки, потому что у json.Unmarshal с ними проблемы
		videoComment.Text = strings.Replace(videoComment.Text, "\n", "\\n", -1)
		text += fmt.Sprintf("\\n\\n%v", videoComment.Text)
	}

	// потом прикрепленную ссылку, если она есть
	if len(link) > 0 {
		if len(videoComment.Text) == 0 {
			text += "\\n"
		}
		text += fmt.Sprintf("\\n%v", link)
	}

	// и ссылку на само видео
	text += fmt.Sprintf("\\n\\n%v", videoURL)

	// экранируем все апострофы, чтобы не сломали нам json.Unmarshal
	text = strings.Replace(text, `"`, `\"`, -1)

	// далее формируем строку с данными для карты
	jsonDump := fmt.Sprintf(`{
		"peer_id": "%d",
		"message": "%v",
		"v": "5.68"
	`, videoCommentMonitorParam.SendTo, text)

	// если прикрепления были, то добавляем их в json-строку отдельно
	if len(attachments) > 0 {
		jsonDump += fmt.Sprintf(`, "attachment": "%v"`, attachments)
	}

	// закрываем карту
	jsonDump += "}"

	return jsonDump, nil
}
