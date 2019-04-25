package main

import (
	"fmt"
	"strings"
)

// WallPostMonitor проверяет посты со стены
func WallPostMonitor(subject Subject) error {
	sender := fmt.Sprintf("%v's wall post monitoring", subject.Name)

	// запрашиваем структуру с параметрами модуля мониторинга постов
	wallPostMonitorParam, err := SelectDBWallPostMonitorParam(subject.ID)
	if err != nil {
		return err
	}

	// запрашиваем структуру с постами со стены субъекта
	wallPosts, err := getWallPosts(sender, subject, wallPostMonitorParam)
	if err != nil {
		return err
	}

	var targetWallPosts []WallPost

	// отфильтровываем старые
	for _, wallPost := range wallPosts {
		if wallPost.Date > wallPostMonitorParam.LastDate {
			targetWallPosts = append(targetWallPosts, wallPost)
		}
	}

	// если после фильтрации что-то осталось, продолжаем
	if len(targetWallPosts) > 0 {

		// сортируем в порядке от раннего к позднему
		for j := 0; j < len(targetWallPosts); j++ {
			f := 0
			for i := 0; i < len(targetWallPosts)-1-j; i++ {
				if targetWallPosts[i].Date > targetWallPosts[i+1].Date {
					x := targetWallPosts[i]
					y := targetWallPosts[i+1]
					targetWallPosts[i+1] = x
					targetWallPosts[i] = y
					f = 1
				}
			}
			if f == 0 {
				break
			}
		}

		// перебираем отсортированный список
		for _, wallPost := range targetWallPosts {

			// проверяем пост на соответствие критериям
			match, err := checkTargetWallPost(wallPostMonitorParam, wallPost)
			if err != nil {
				return err
			}

			// если соответствует, то отправляем пост
			if match {

				// формируем строку с данными для карты для отправки сообщения
				messageParameters, err := makeMessageWallPost(sender, subject, wallPostMonitorParam, wallPost)
				if err != nil {
					return err
				}

				// отправляем сообщение с полученными данными
				if err := SendMessage(sender, messageParameters, subject); err != nil {
					return err
				}

				// выводим в консоль сообщение о новом посте
				outputReportAboutNewWallPost(sender, wallPost)
			}

			// обновляем дату последнего проверенного поста в БД
			if err := UpdateDBWallPostMonitorLastDate(subject.ID, wallPost.Date, wallPostMonitorParam); err != nil {
				return err
			}
		}
	}

	return nil
}

// checkTargetWallPost проверяет пост на соответствие критериям
func checkTargetWallPost(wallPostMonitorParam WallPostMonitorParam, wallPost WallPost) (bool, error) {
	var match bool
	keywords, err := MakeParamList(wallPostMonitorParam.KeywordsForMonitoring)
	if err != nil {
		return match, err
	}
	if len(keywords.List) > 0 {
		for _, keyword := range keywords.List {
			match = strings.Contains(wallPost.Text, keyword)
			if match {
				return match, nil
			}
		}
	} else {
		match = true
	}

	return match, nil
}

// outputReportAboutNewWallPost выводит сообщение о новом посте
func outputReportAboutNewWallPost(sender string, wallPost WallPost) {
	creationDate := UnixTimeStampToDate(wallPost.Date)
	message := fmt.Sprintf("New %v at %v.", wallPost.PostType, creationDate)
	OutputMessage(sender, message)
}

// WallPost - структура для данных о посте со стены
type WallPost struct {
	ID          int          `json:"id"`
	OwnerID     int          `json:"owner_id"`
	FromID      int          `json:"from_id"`
	Date        int          `json:"date"`
	PostType    string       `json:"post_type"`
	Text        string       `json:"text"`
	Attachments []Attachment `json:"attachments"`
}

// getWallPosts формирует запрос на получение постов со стены и посылает его к vk api
func getWallPosts(sender string, subject Subject, wallPostMonitorParam WallPostMonitorParam) ([]WallPost, error) {
	var wallPosts []WallPost

	// формируем карту с параметрами запроса
	jsonDump := fmt.Sprintf(`{
			"owner_id": "%d",
			"count": "%d",
			"filter": "%v",
			"v": "5.95"
		}`, subject.SubjectID, wallPostMonitorParam.PostsCount, wallPostMonitorParam.Filter)
	values, err := MakeJSON(jsonDump)
	if err != nil {
		return wallPosts, err
	}

	// отправляем запрос, получаем ответ
	response, err := SendVKAPIQuery(sender, "wall.get", values, subject)
	if err != nil {
		return wallPosts, err
	}

	// парсим полученные данные о постах
	wallPosts = ParseWallPostsVkAPIMap(response["response"].(map[string]interface{}))

	return wallPosts, nil
}

func parseCopyHistory(copyHistory []interface{}) Attachment {
	var attachment Attachment

	itemsMap := copyHistory[0].(map[string]interface{})
	// извлекаем данные об одном единственном репосте из карты vk api
	attachment.Type = itemsMap["post_type"].(string)
	attachment.OwnerID = int(itemsMap["owner_id"].(float64))
	attachment.ID = int(itemsMap["id"].(float64))
	if accessKey, exist := itemsMap["access_key"]; exist {
		attachment.AccessKey = accessKey.(string)
	}
	return attachment
}

// ParseWallPostsVkAPIMap извлекает данные о постах из полученной карты vk api
func ParseWallPostsVkAPIMap(resp map[string]interface{}) []WallPost {
	var wallPosts []WallPost

	// перебираем элементы с данными о постах
	itemsMap := resp["items"].([]interface{})
	for _, itemMap := range itemsMap {
		item := itemMap.(map[string]interface{})

		// парсим данные из элемента в структуру
		var wallPost WallPost
		wallPost.ID = int(item["id"].(float64))
		wallPost.OwnerID = int(item["owner_id"].(float64))
		if _, exist := item["signer_id"]; exist == true {
			wallPost.FromID = int(item["signer_id"].(float64))
		} else {
			wallPost.FromID = int(item["from_id"].(float64))
		}
		wallPost.Date = int(item["date"].(float64))
		wallPost.PostType = item["post_type"].(string)
		wallPost.Text = item["text"].(string)

		// если есть прикрепления, то вызываем парсер прикреплений
		if mediaContent, exist := item["attachments"]; exist == true {
			wallPost.Attachments = ParseAttachments(mediaContent.([]interface{}))
		}

		// если есть репосты, то вызываем парсер репостов
		if copyHistory, exist := item["copy_history"]; exist == true {
			wallPost.Attachments = append(wallPost.Attachments, parseCopyHistory(
				copyHistory.([]interface{})))
		}
		wallPosts = append(wallPosts, wallPost)
	}
	return wallPosts
}

// MakeMessageWallPost собирает сообщение с данными о посте
func makeMessageWallPost(sender string, subject Subject,
	wallPostMonitorParam WallPostMonitorParam, wallPost WallPost) (string, error) {

	// собираем данные для сигнатуры сообщения:
	// тип поста
	postType := fmt.Sprintf("%v", wallPost.PostType)

	// где он был обнаружен
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

	// кто автор, если это пользователь
	var authorHyperlink string
	if wallPost.FromID < 0 {
		authorHyperlink = "[no_data]"
	} else {
		vkUser, err := GetUserInfo(sender, subject, wallPost.FromID)
		authorHyperlink = MakeUserHyperlink(vkUser)
		if err != nil {
			return "", err
		}
	}

	// и когда был создан
	creationDate := UnixTimeStampToDate(wallPost.Date)

	// формируем строку с прикреплениями
	var attachments string
	var link string
	if len(wallPost.Attachments) > 0 {
		for _, attachment := range wallPost.Attachments {
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

	// собираем ссылку на пост
	postURL := fmt.Sprintf("https://vk.com/wall%d_%d", wallPost.OwnerID, wallPost.ID)

	// добавляем подготовленные фрагменты сообщения в общий текст
	// сначала сигнатуру
	text := fmt.Sprintf("New %v\\nLocation: %v\\nAuthor: %v\\nCreated: %v",
		postType, locationHyperlink, authorHyperlink, creationDate)

	// затем основной текст поста, если он есть
	if len(wallPost.Text) > 0 {
		// но сначала обрезаем его из-за ограничения на длину запроса
		if len(wallPost.Text) > 800 {
			wallPost.Text = string(wallPost.Text[0:800])
		}
		// и экранируем все символы пропуска строки, потому что у json.Unmarshal с ними проблемы
		wallPost.Text = strings.Replace(wallPost.Text, "\n", "\\n", -1)
		text += fmt.Sprintf("\\n\\n%v", wallPost.Text)
	}

	// потом прикрепленную ссылку, если она есть
	if len(link) > 0 {
		if len(wallPost.Text) == 0 {
			text += "\\n"
		}
		text += fmt.Sprintf("\\n%v", link)
	}

	// и ссылку на сам пост
	text += fmt.Sprintf("\\n\\n%v", postURL)

	// экранируем все апострофы, чтобы не сломали нам json.Unmarshal
	text = strings.Replace(text, `"`, `\"`, -1)

	// далее формируем строку с данными для карты
	jsonDump := fmt.Sprintf(`{
		"peer_id": "%d",
		"message": "%v",
		"v": "5.68"
	`, wallPostMonitorParam.SendTo, text)

	// если прикрепления были, то добавляем их в json-строку отдельно
	if len(attachments) > 0 {
		jsonDump += fmt.Sprintf(`, "attachment": "%v"`, attachments)
	}

	// закрываем карту
	jsonDump += "}"

	return jsonDump, nil
}
