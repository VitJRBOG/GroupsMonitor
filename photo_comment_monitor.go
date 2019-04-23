package main

import (
	"fmt"
	"strings"
)

// PhotoCommentMonitor проверяет комментарии под фотографиями
func PhotoCommentMonitor(subject Subject) error {

	sender := fmt.Sprintf("%v's photo comment monitoring", subject.Name)

	// запрашиваем структуру с параметрами модуля мониторинга постов
	photoCommentMonitorParam, err := SelectDBPhotoCommentMonitorParam(subject.ID)
	if err != nil {
		return err
	}

	// запрашиваем структуру с постами со стены субъекта
	photoComments, err := getPhotoComments(sender, subject, photoCommentMonitorParam)
	if err != nil {
		return err
	}

	var targetPhotoComments []PhotoComment

	// отфильтровываем старые
	for _, photoComment := range photoComments {
		if photoComment.Date > photoCommentMonitorParam.LastDate {
			targetPhotoComments = append(targetPhotoComments, photoComment)
		}
	}

	// если после фильтрации что-то осталось, продолжаем
	if len(targetPhotoComments) > 0 {

		// сортируем в порядке от раннего к позднему
		for j := 0; j < len(targetPhotoComments); j++ {
			f := 0
			for i := 0; i < len(targetPhotoComments)-1-j; i++ {
				if targetPhotoComments[i].Date > targetPhotoComments[i+1].Date {
					x := targetPhotoComments[i]
					y := targetPhotoComments[i+1]
					targetPhotoComments[i+1] = x
					targetPhotoComments[i] = y
					f = 1
				}
			}
			if f == 0 {
				break
			}
		}

		// перебираем отсортированный список
		for _, photoComment := range targetPhotoComments {

			// формируем строку с данными для карты для отправки сообщения
			messageParameters, err := makeMessagePhotoComment(sender, subject,
				photoCommentMonitorParam, photoComment)
			if err != nil {
				return err
			}

			// отправляем сообщение с полученными данными
			if err := SendMessage(sender, messageParameters, subject); err != nil {
				return err
			}

			// выводим в консоль сообщение о новом посте
			outputReportAboutNewPhotoComment(sender, photoComment)

			// обновляем дату последнего проверенного поста в БД
			if err := UpdateDBPhotoCommentMonitorLastDate(subject.ID, photoComment.Date); err != nil {
				return err
			}
		}
	}

	return nil
}

// outputReportAboutNewPhotoComment выводит сообщение о новом комментарии под фотографией
func outputReportAboutNewPhotoComment(sender string, photoComment PhotoComment) {
	creationDate := UnixTimeStampToDate(photoComment.Date + 18000) // цифру изменить под свой часовой пояс (тут 5 часов)
	message := fmt.Sprintf("New comment under photo at %v.", creationDate)
	OutputMessage(sender, message)
}

// PhotoComment - структура для данных о комментариях под фотографиями
type PhotoComment struct {
	ID           int          `json:"id"`
	PhotoOwnerID int          `json:"owner_id"`
	PhotoID      int          `json:"pid"`
	FromID       int          `json:"from_id"`
	Date         int          `json:"date"`
	Text         string       `json:"text"`
	Attachments  []Attachment `json:"attachments"`
}

// getPhotoComments формирует запрос на получение комментариев под фотографиями и посылает его к vk api
func getPhotoComments(sender string, subject Subject,
	photoCommentMonitorParam PhotoCommentMonitorParam) ([]PhotoComment, error) {
	var photoComments []PhotoComment

	// формируем карту с параметрами запроса
	jsonDump := fmt.Sprintf(`{
			"owner_id": "%d",
			"count": "%d",
			"v": "5.95"
		}`, subject.SubjectID, photoCommentMonitorParam.CommentsCount)
	values, err := MakeJSON(jsonDump)
	if err != nil {
		return photoComments, err
	}

	// отправляем запрос, получаем ответ
	response, err := SendVKAPIQuery(sender, "photos.getAllComments", values, subject)
	if err != nil {
		return photoComments, err
	}

	// парсим полученные данные о комментариях
	photoComments = parsePhotoCommentsVkAPIMap(response["response"].(map[string]interface{}), subject)

	return photoComments, nil
}

// parsePhotoCommentsVkAPIMap извлекает данные о комментариях из полученной карты vk api
func parsePhotoCommentsVkAPIMap(resp map[string]interface{}, subject Subject) []PhotoComment {
	var photoComments []PhotoComment

	// перебираем элементы с данными о комментариях
	itemsMap := resp["items"].([]interface{})
	for _, itemMap := range itemsMap {
		item := itemMap.(map[string]interface{})

		// парсим данные из элемента в структуру
		var photoComment PhotoComment
		photoComment.ID = int(item["id"].(float64))
		photoComment.PhotoID = int(item["pid"].(float64))
		photoComment.PhotoOwnerID = subject.SubjectID
		photoComment.FromID = int(item["from_id"].(float64))
		photoComment.Date = int(item["date"].(float64))
		photoComment.Text = item["text"].(string)

		// если есть прикрепления, то вызываем парсер прикреплений
		if mediaContent, exist := item["attachments"]; exist == true {
			photoComment.Attachments = ParseAttachments(mediaContent.([]interface{}))
		}
		photoComments = append(photoComments, photoComment)
	}
	return photoComments
}

// makeMessagePhotoComment собирает сообщение с данными о комментарии
func makeMessagePhotoComment(sender string, subject Subject,
	photoCommentMonitorParam PhotoCommentMonitorParam, photoComment PhotoComment) (string, error) {
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
	if photoComment.FromID < 0 {
		vkCommunity, err := GetCommunityInfo(sender, subject, photoComment.FromID)
		authorHyperlink = MakeCommunityHyperlink(vkCommunity)
		if err != nil {
			return "", err
		}
	} else {
		vkUser, err := GetUserInfo(sender, subject, photoComment.FromID)
		authorHyperlink = MakeUserHyperlink(vkUser)
		if err != nil {
			return "", err
		}
	}

	// и когда был создан
	creationDate := UnixTimeStampToDate(photoComment.Date + 18000) // цифру изменить под свой часовой пояс (тут 5 часов)

	// формируем строку с прикреплениями
	var attachments string
	var link string
	if len(photoComment.Attachments) > 0 {
		for _, attachment := range photoComment.Attachments {
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

	// собираем ссылку на фото, потому что ссылку на комментарий под фото сделать нельзя
	photoURL := fmt.Sprintf("https://vk.com/photo%d_%d", photoComment.PhotoOwnerID, photoComment.PhotoID)

	// добавляем подготовленные фрагменты сообщения в общий текст
	// сначала сигнатуру
	text := fmt.Sprintf("New photo comment\\nLocation: %v\\nAuthor: %v\\nCreated: %v",
		locationHyperlink, authorHyperlink, creationDate)

	// затем основной текст комментария, если он есть
	if len(photoComment.Text) > 0 {
		// но сначала экранируем все символы пропуска строки, потому что у json.Unmarshal с ними проблемы
		photoComment.Text = strings.Replace(photoComment.Text, "\n", "\\n", -1)
		text += fmt.Sprintf("\\n\\n%v", photoComment.Text)
	}

	// потом прикрепленную ссылку, если она есть
	if len(link) > 0 {
		if len(photoComment.Text) == 0 {
			text += "\\n"
		}
		text += fmt.Sprintf("\\n%v", link)
	}

	// и ссылку на саму фотку
	text += fmt.Sprintf("\\n\\n%v", photoURL)

	// далее формируем строку с данными для карты
	jsonDump := fmt.Sprintf(`{
		"peer_id": "%d",
		"message": "%v",
		"v": "5.68"
	`, photoCommentMonitorParam.SendTo, text)

	// если прикрепления были, то добавляем их в json-строку отдельно
	if len(attachments) > 0 {
		jsonDump += fmt.Sprintf(`, "attachment": "%v"`, attachments)
	}

	// закрываем карту
	jsonDump += "}"

	return jsonDump, nil
}
