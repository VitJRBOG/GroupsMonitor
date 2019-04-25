package main

import (
	"fmt"
	"strings"
)

// TopicMonitor проверяет комментарии в топиках обсуждений
func TopicMonitor(subject Subject) error {
	sender := fmt.Sprintf("%v's topic monitoring", subject.Name)

	// запрашиваем структуру с параметрами модуля мониторинга топиков обсуждений
	topicMonitorParam, err := SelectDBTopicMonitorParam(subject.ID)
	if err != nil {
		return err
	}

	// запрашиваем структуру с топиками обсуждений
	topics, err := getTopics(sender, subject, topicMonitorParam)
	if err != nil {
		return err
	}

	var targetTopics []Topic

	// перебираем полученные топики обсуждений
	for _, topic := range topics {
		// отфильтровываем те, которые не обновлялись с момента последней проверки
		if topic.UpdateDate > topicMonitorParam.LastDate {
			targetTopics = append(targetTopics, topic)
		}
	}

	// если после фильтрации что-то осталось, продолжаем
	if len(targetTopics) > 0 {

		// запрашиваем комментарии из этих топиков
		topicComments, err := getTopicComments(sender, subject, topicMonitorParam, targetTopics)
		if err != nil {
			return err
		}

		var targetTopicComments []TopicComment

		// отфильтровываем старые комментарии
		for _, topicComment := range topicComments {
			if topicComment.Date > topicMonitorParam.LastDate {
				targetTopicComments = append(targetTopicComments, topicComment)
			}
		}

		// если после фильтрации что-то осталось, продолжаем
		if len(targetTopicComments) > 0 {

			// сортируем в порядке от раннего к позднему
			for j := 0; j < len(targetTopicComments); j++ {
				f := 0
				for i := 0; i < len(targetTopicComments)-1-j; i++ {
					if targetTopicComments[i].Date > targetTopicComments[i+1].Date {
						x := targetTopicComments[i]
						y := targetTopicComments[i+1]
						targetTopicComments[i+1] = x
						targetTopicComments[i] = y
						f = 1
					}
				}
				if f == 0 {
					break
				}
			}

			// перебираем отсортированный список
			for _, topicComment := range targetTopicComments {

				// формируем строку с данными для карты для отправки сообщения
				messageParameters, err := makeMessageTopicComment(sender, subject,
					topicMonitorParam, topicComment)
				if err != nil {
					return err
				}

				// отправляем сообщение с полученными данными
				if err := SendMessage(sender, messageParameters, subject); err != nil {
					return err
				}

				// выводим в консоль сообщение о новом комментарии
				outputReportAboutNewTopicComment(sender, topicComment)

				// обновляем дату последнего проверенного комментария в БД
				if err := UpdateDBTopicMonitorLastDate(subject.ID, topicComment.Date); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// outputReportAboutNewTopicComment выводит сообщение о новом комментарии
func outputReportAboutNewTopicComment(sender string, topicComment TopicComment) {
	creationDate := UnixTimeStampToDate(topicComment.Date)
	message := fmt.Sprintf("New comment in %v at %v.", topicComment.TopicName, creationDate)
	OutputMessage(sender, message)
}

// Topic - структура для данных о топиках обсуждений
type Topic struct {
	ID         int    `json:"id"`
	OwnerID    int    `json:"owner_id"`
	Title      string `json:"title"`
	UpdateDate int    `json:"updated"`
}

// getTopics формирует запрос на получение топиков обсуждений и посылает его к vk api
func getTopics(sender string, subject Subject,
	topicMonitorParam TopicMonitorParam) ([]Topic, error) {
	var topics []Topic

	// формируем карту с параметрами запроса
	jsonDump := fmt.Sprintf(`{
			"group_id": "%d",
			"count": "%d",
			"v": "5.95"
		}`, subject.SubjectID*-1, // умножаем на -1, потому что id паблика должен быть без минуса
		topicMonitorParam.TopicsCount)
	values, err := MakeJSON(jsonDump)
	if err != nil {
		return topics, err
	}

	// отправляем запрос, получаем ответ
	response, err := SendVKAPIQuery(sender, "board.getTopics", values, subject)
	if err != nil {
		return topics, err
	}

	topics = ParseTopicVkAPIMap(response["response"].(map[string]interface{}), subject)

	return topics, nil
}

// ParseTopicVkAPIMap парсит карты с данными о топиках обсуждений
func ParseTopicVkAPIMap(resp map[string]interface{}, subject Subject) []Topic {
	var topics []Topic

	// перебираем элементы с данными о фотографиях
	itemsMap := resp["items"].([]interface{})
	for _, itemMap := range itemsMap {
		item := itemMap.(map[string]interface{})

		// парсим данные из элемента в структуру
		var topic Topic
		topic.ID = int(item["id"].(float64))
		topic.OwnerID = subject.SubjectID
		topic.UpdateDate = int(item["updated"].(float64))
		topic.Title = item["title"].(string)

		topics = append(topics, topic)
	}

	return topics
}

// TopicComment - структура для данных о комментарии в топиках обсуждений
type TopicComment struct {
	ID           int          `json:"id"`
	TopicID      int          `json:"tid"`
	TopicName    string       `json:"title"`
	TopicOwnerID int          `json:"owner_id"`
	FromID       int          `json:"from_id"`
	Date         int          `json:"date"`
	Text         string       `json:"text"`
	Attachments  []Attachment `json:"attachments"`
}

// getTopicComments формирует запрос на получение комментариев в топиках обсуждений и посылает его к vk api
func getTopicComments(sender string, subject Subject,
	topicMonitorParam TopicMonitorParam, topics []Topic) ([]TopicComment, error) {
	var topicsComments []TopicComment

	for _, topic := range topics {

		// формируем карту с параметрами запроса
		jsonDump := fmt.Sprintf(`{
			"group_id": "%d",
			"topic_id": "%d",
			"count": "%d",
			"sort": "desc",
			"v": "5.95"
		}`, topic.OwnerID*-1, // умножаем на -1, потому что id паблика должен быть без минуса
			topic.ID, topicMonitorParam.CommentsCount)
		values, err := MakeJSON(jsonDump)
		if err != nil {
			return topicsComments, err
		}

		// отправляем запрос, получаем ответ
		response, err := SendVKAPIQuery(sender, "board.getComments", values, subject)
		if err != nil {
			return topicsComments, err
		}

		topicComments := parseTopicCommentVkAPIMap(response["response"].(map[string]interface{}), topic)
		topicsComments = append(topicsComments, topicComments...)
	}

	return topicsComments, nil
}

// parseTopicCommentVkAPIMap извлекает данные о комментариях из полученной карты vk api
func parseTopicCommentVkAPIMap(resp map[string]interface{}, topic Topic) []TopicComment {
	var topicComments []TopicComment

	// перебираем элементы с данными о комментариях
	itemsMap := resp["items"].([]interface{})
	for _, itemMap := range itemsMap {
		item := itemMap.(map[string]interface{})

		// парсим данные из элемента в структуру
		var topicComment TopicComment
		topicComment.ID = int(item["id"].(float64))
		topicComment.TopicName = topic.Title
		topicComment.TopicID = topic.ID
		topicComment.TopicOwnerID = topic.OwnerID
		topicComment.FromID = int(item["from_id"].(float64))
		topicComment.Date = int(item["date"].(float64))
		topicComment.Text = item["text"].(string)

		// если есть прикрепления, то вызываем парсер прикреплений
		if mediaContent, exist := item["attachments"]; exist == true {
			topicComment.Attachments = ParseAttachments(mediaContent.([]interface{}))
		}
		topicComments = append(topicComments, topicComment)
	}
	return topicComments
}

// makeMessageTopicComment собирает сообщение с данными о комментарии
func makeMessageTopicComment(sender string, subject Subject,
	topicMonitorParam TopicMonitorParam, topicComment TopicComment) (string, error) {
	// собираем данные для сигнатуры сообщения:
	// название топика обсуждения
	topicName := topicComment.TopicName

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
	if topicComment.FromID < 0 {
		vkCommunity, err := GetCommunityInfo(sender, subject, topicComment.FromID)
		authorHyperlink = MakeCommunityHyperlink(vkCommunity)
		if err != nil {
			return "", err
		}
	} else {
		vkUser, err := GetUserInfo(sender, subject, topicComment.FromID)
		authorHyperlink = MakeUserHyperlink(vkUser)
		if err != nil {
			return "", err
		}
	}

	// и когда был создан
	creationDate := UnixTimeStampToDate(topicComment.Date)

	// формируем строку с прикреплениями
	var attachments string
	var link string
	if len(topicComment.Attachments) > 0 {
		for _, attachment := range topicComment.Attachments {
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

	// собираем ссылку на комментарий в топике обсуждения
	topicCommentULR := fmt.Sprintf("https://vk.com/topic%d_%d?post=%d",
		topicComment.TopicOwnerID, topicComment.TopicID, topicComment.ID)

	// добавляем подготовленные фрагменты сообщения в общий текст
	// сначала сигнатуру
	text := fmt.Sprintf("New topic comment\\nTopic: %v\\nLocation: %v\\nAuthor: %v\\nCreated: %v",
		topicName, locationHyperlink, authorHyperlink, creationDate)

	// затем основной текст комментария, если он есть
	if len(topicComment.Text) > 0 {
		// но сначала обрезаем его из-за ограничения на длину запроса
		if len(topicComment.Text) > 800 {
			topicComment.Text = string(topicComment.Text[0:800])
		}
		// и экранируем все символы пропуска строки, потому что у json.Unmarshal с ними проблемы
		topicComment.Text = strings.Replace(topicComment.Text, "\n", "\\n", -1)
		text += fmt.Sprintf("\\n\\n%v", topicComment.Text)
	}

	// потом прикрепленную ссылку, если она есть
	if len(link) > 0 {
		if len(topicComment.Text) == 0 {
			text += "\\n"
		}
		text += fmt.Sprintf("\\n%v", link)
	}

	// и ссылку на сам комментарий
	text += fmt.Sprintf("\\n\\n%v", topicCommentULR)

	// экранируем все апострофы, чтобы не сломали нам json.Unmarshal
	text = strings.Replace(text, `"`, `\"`, -1)

	// далее формируем строку с данными для карты
	jsonDump := fmt.Sprintf(`{
		"peer_id": "%d",
		"message": "%v",
		"v": "5.68"
	`, topicMonitorParam.SendTo, text)

	// если прикрепления были, то добавляем их в json-строку отдельно
	if len(attachments) > 0 {
		jsonDump += fmt.Sprintf(`, "attachment": "%v"`, attachments)
	}

	// закрываем карту
	jsonDump += "}"

	return jsonDump, nil
}
