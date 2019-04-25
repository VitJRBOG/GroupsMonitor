package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// WallPostCommentMonitor проверяет комментарии под постами на стене
func WallPostCommentMonitor(subject Subject) error {

	sender := fmt.Sprintf("%v's wall post comment monitoring", subject.Name)

	// запрашиваем структуру с параметрами модуля мониторинга комментариев под постами
	wallPostCommentMonitorParam, err := SelectDBWallPostCommentMonitorParam(subject.ID)
	if err != nil {
		return err
	}

	// запрашиваем структуру с постами со стены субъекта
	wallPosts, err := getWallPostsForComments(sender, subject, wallPostCommentMonitorParam)
	if err != nil {
		return err
	}

	// запрашиваем комментарии из под каждого поста
	wallPostComments, err := getWallPostComments(sender, subject, wallPostCommentMonitorParam, wallPosts)
	if err != nil {
		return err
	}

	var targetWallPostComments []WallPostComment

	// отфильтровываем старые
	for _, wallPostComment := range wallPostComments {
		if wallPostComment.Date > wallPostCommentMonitorParam.LastDate {
			targetWallPostComments = append(targetWallPostComments, wallPostComment)
		}
	}

	// если после фильтрации что-то осталось, продолжаем
	if len(targetWallPostComments) > 0 {

		// запрашиваем данные по веткам комментариев
		targetWallPostComments, err := getWallPostComment(sender, subject,
			wallPostCommentMonitorParam, targetWallPostComments)
		if err != nil {
			return err
		}

		// сортируем в порядке от раннего к позднему
		for j := 0; j < len(targetWallPostComments); j++ {
			f := 0
			for i := 0; i < len(targetWallPostComments)-1-j; i++ {
				if targetWallPostComments[i].Date > targetWallPostComments[i+1].Date {
					x := targetWallPostComments[i]
					y := targetWallPostComments[i+1]
					targetWallPostComments[i+1] = x
					targetWallPostComments[i] = y
					f = 1
				}
			}
			if f == 0 {
				break
			}
		}

		// перебираем отсортированный список
		for _, wallPostComment := range targetWallPostComments {

			// проверяем комментарий на соответствие критериям
			match, err := checkTargetWallPostComment(wallPostCommentMonitorParam, wallPostComment, subject)
			if err != nil {
				return err
			}

			// если соответствует, то отправляем комментарий
			if match {

				// формируем строку с данными для карты для отправки сообщения
				messageParameters, err := makeMessageWallPostComment(sender, subject,
					wallPostCommentMonitorParam, wallPostComment)
				if err != nil {
					return err
				}

				// отправляем сообщение с полученными данными
				if err := SendMessage(sender, messageParameters, subject); err != nil {
					return err
				}

				// выводим в консоль сообщение о новом комментарии под постом
				outputReportAboutNewWallPostComment(sender, wallPostComment)
			}

			// обновляем дату последнего проверенного комментария в БД
			if err := UpdateDBWallPostCommentMonitorLastDate(subject.ID, wallPostComment.Date); err != nil {
				return err
			}
		}
	}

	return nil
}

// checkTargetWallPostComment проверяет комментарий на соответствие критериям
func checkTargetWallPostComment(wallPostCommentMonitorParam WallPostCommentMonitorParam,
	wallPostComment WallPostComment, subject Subject) (bool, error) {
	var match bool

	if wallPostCommentMonitorParam.MonitoringAll == 1 {
		match = true
		return match, nil
	}

	ignore, err := checkByIgnoreUsersIDs(wallPostCommentMonitorParam, wallPostComment)
	if err != nil {
		return false, err
	}
	if ignore == true {
		return false, nil
	}

	if wallPostCommentMonitorParam.MonitorByCommunity == 1 {
		match = checkByCommentsFromCommunity(wallPostComment)
	}
	if match == true {
		return match, nil
	}

	match, err = checkByAttachments(wallPostCommentMonitorParam, wallPostComment)
	if err != nil {
		return false, err
	}
	if match == true {
		return match, nil
	}

	match, err = checkByUsersIDs(wallPostCommentMonitorParam, wallPostComment)
	if err != nil {
		return false, err
	}
	if match == true {
		return match, nil
	}

	match, err = checkByUsersNames(wallPostCommentMonitorParam, wallPostComment, subject)
	if err != nil {
		return false, err
	}
	if match == true {
		return match, nil
	}

	match, err = checkBySmallComments(wallPostCommentMonitorParam, wallPostComment, 0)
	if err != nil {
		return false, err
	}
	if match == true {
		return match, nil
	}

	match, err = checkByKeywords(wallPostCommentMonitorParam, wallPostComment, 0)
	if err != nil {
		return false, err
	}
	if match == true {
		return match, nil
	}

	match, err = checkByPhoneNumber(wallPostCommentMonitorParam, wallPostComment)
	if err != nil {
		return false, err
	}
	if match == true {
		return match, nil
	}

	match, err = checkByCardNumber(wallPostCommentMonitorParam, wallPostComment)
	if err != nil {
		return false, err
	}
	if match == true {
		return match, nil
	}

	return match, nil
}

func checkByIgnoreUsersIDs(wallPostCommentMonitorParam WallPostCommentMonitorParam,
	wallPostComment WallPostComment) (bool, error) {
	usersIDs, err := MakeParamList(wallPostCommentMonitorParam.UsersIDsForIgnore)
	if err != nil {
		return false, err
	}
	if len(usersIDs.List) > 0 {
		for _, userID := range usersIDs.List {
			numUser, err := strconv.Atoi(userID)
			if err != nil {
				return false, err
			}
			if wallPostComment.FromID == numUser {
				return true, nil
			}
		}
	}
	return false, nil
}

func checkByCommentsFromCommunity(wallPostComment WallPostComment) bool {
	if wallPostComment.FromID < 0 {
		return true
	}
	return false
}

func checkByAttachments(wallPostCommentMonitorParam WallPostCommentMonitorParam,
	wallPostComment WallPostComment) (bool, error) {
	if len(wallPostComment.Attachments) > 0 {
		attachmentsFromParam, err := MakeParamList(wallPostCommentMonitorParam.AttachmentsTypesForMonitoring)
		if err != nil {
			return false, err
		}
		if len(attachmentsFromParam.List) > 0 {
			for _, attachmentFromParam := range attachmentsFromParam.List {
				for _, attachment := range wallPostComment.Attachments {
					if attachment.Type == attachmentFromParam {
						return true, nil
					}
				}
			}
		}
	}
	return false, nil
}

func checkByUsersIDs(wallPostCommentMonitorParam WallPostCommentMonitorParam,
	wallPostComment WallPostComment) (bool, error) {
	usersIDs, err := MakeParamList(wallPostCommentMonitorParam.UsersIDsForMonitoring)
	if err != nil {
		return false, err
	}
	if len(usersIDs.List) > 0 {
		for _, userID := range usersIDs.List {
			numUser, err := strconv.Atoi(userID)
			if err != nil {
				return false, err
			}
			if wallPostComment.FromID == numUser {
				return true, nil
			}
		}
	}
	return false, nil
}

func checkByUsersNames(wallPostCommentMonitorParam WallPostCommentMonitorParam,
	wallPostComment WallPostComment, subject Subject) (bool, error) {
	sender := fmt.Sprintf("%v's wall post comment monitoring", subject.Name)
	usersNames, err := MakeParamList(wallPostCommentMonitorParam.UsersNamesForMonitoring)
	if err != nil {
		return false, err
	}
	if len(usersNames.List) > 0 {
		for _, userName := range usersNames.List {
			var authorHyperlink string
			if wallPostComment.FromID < 0 {
				vkCommunity, err := GetCommunityInfo(sender, subject, wallPostComment.FromID)
				authorHyperlink = MakeCommunityHyperlink(vkCommunity)
				if err != nil {
					return false, err
				}
			} else {
				vkUser, err := GetUserInfo(sender, subject, wallPostComment.FromID)
				authorHyperlink = MakeUserHyperlink(vkUser)
				if err != nil {
					return false, err
				}
			}
			match := strings.Contains(authorHyperlink, userName)
			if match {
				return true, nil
			}
		}
	}
	return false, nil
}

func checkBySmallComments(wallPostCommentMonitorParam WallPostCommentMonitorParam,
	wallPostComment WallPostComment, step int) (bool, error) {
	smallComments, err := MakeParamList(wallPostCommentMonitorParam.SmallCommentsForMonitoring)
	if err != nil {
		return false, err
	}
	if len(smallComments.List) > 0 {
		for _, smallComment := range smallComments.List {
			if wallPostComment.Text == smallComment {
				return true, nil
			}
		}
	}
	switch step {
	case 0:
		step++
		wallPostComment.Text = CharChange(wallPostComment.Text, "lat_to_cyr")
		return checkBySmallComments(wallPostCommentMonitorParam, wallPostComment, step)
	case 1:
		step++
		wallPostComment.Text = CharChange(wallPostComment.Text, "cyr_to_lat")
		return checkBySmallComments(wallPostCommentMonitorParam, wallPostComment, step)
	}

	return false, nil
}

func checkByKeywords(wallPostCommentMonitorParam WallPostCommentMonitorParam,
	wallPostComment WallPostComment, step int) (bool, error) {
	keywords, err := MakeParamList(wallPostCommentMonitorParam.KeywordsForMonitoring)
	if err != nil {
		return false, err
	}
	if len(keywords.List) > 0 {
		for _, keyword := range keywords.List {
			match := strings.Contains(wallPostComment.Text, keyword)
			if match {
				return true, nil
			}
		}
	}
	switch step {
	case 0:
		step++
		wallPostComment.Text = CharChange(wallPostComment.Text, "lat_to_cyr")
		return checkBySmallComments(wallPostCommentMonitorParam, wallPostComment, step)
	case 1:
		step++
		wallPostComment.Text = CharChange(wallPostComment.Text, "cyr_to_lat")
		return checkBySmallComments(wallPostCommentMonitorParam, wallPostComment, step)
	}

	return false, nil
}

func checkByPhoneNumber(wallPostCommentMonitorParam WallPostCommentMonitorParam,
	wallPostComment WallPostComment) (bool, error) {
	digitsCounts, err := MakeParamList(wallPostCommentMonitorParam.DigitsCountForPhoneNumberMonitoring)
	if err != nil {
		return false, err
	}
	if len(digitsCounts.List) > 0 {
		reDigits := regexp.MustCompile("[0-9]+")
		reSymbols := regexp.MustCompile("[() -]+")
		var repeats int
		var interrupts int
		for i := 0; i < len(wallPostComment.Text); i++ {
			if len(reDigits.FindAllString(string(wallPostComment.Text[i]), -1)) > 0 {
				repeats++
				interrupts = 0
			} else {
				if len(reSymbols.FindAllString(string(wallPostComment.Text[i]), -1)) > 0 {
					if interrupts == 0 {
						interrupts++
					} else {
						interrupts = 0
						repeats = 0
					}
				} else {
					repeats = 0
					interrupts = 0
				}
			}
			for _, needDigits := range digitsCounts.List {
				numNeedDigits, err := strconv.Atoi(needDigits)
				if err != nil {
					return false, err
				}
				if repeats == numNeedDigits {
					if i < len(wallPostComment.Text)-1 {
						if len(reDigits.FindAllString(string(wallPostComment.Text[i+1]), -1)) == 0 {
							return true, nil
						}
					} else {
						return true, nil
					}
				}
			}
		}
	}
	return false, nil
}

func checkByCardNumber(wallPostCommentMonitorParam WallPostCommentMonitorParam,
	wallPostComment WallPostComment) (bool, error) {
	digitsCounts, err := MakeParamList(wallPostCommentMonitorParam.DigitsCountForCardNumberMonitoring)
	if err != nil {
		return false, err
	}
	if len(digitsCounts.List) > 0 {
		reDigits := regexp.MustCompile("[0-9]+")
		reSymbols := regexp.MustCompile("[ ]+")
		var repeats int
		var interrupts int
		for i := 0; i < len(wallPostComment.Text); i++ {
			if len(reDigits.FindAllString(string(wallPostComment.Text[i]), -1)) > 0 {
				repeats++
				interrupts = 0
			} else {
				if len(reSymbols.FindAllString(string(wallPostComment.Text[i]), -1)) > 0 {
					if interrupts == 0 {
						interrupts++
					} else {
						interrupts = 0
						repeats = 0
					}
				} else {
					repeats = 0
					interrupts = 0
				}
			}
			for _, needDigits := range digitsCounts.List {
				numNeedDigits, err := strconv.Atoi(needDigits)
				if err != nil {
					return false, err
				}
				if repeats == numNeedDigits {
					if i < len(wallPostComment.Text)-1 {
						if len(reDigits.FindAllString(string(wallPostComment.Text[i+1]), -1)) == 0 {
							return true, nil
						}
					} else {
						return true, nil
					}
				}
			}
		}
	}
	return false, nil
}

// outputReportAboutNewWallPostComment выводит сообщение о новом комментарии под постом
func outputReportAboutNewWallPostComment(sender string, wallPostComment WallPostComment) {
	creationDate := UnixTimeStampToDate(wallPostComment.Date)
	message := fmt.Sprintf("New comment under wall post at %v.", creationDate)
	OutputMessage(sender, message)
}

// WallPost
// структура уже описана в модуле мониторинга новых постов,
// поэтому тут этот момент опустим

// getWallPostsForComments формирует запрос на получение постов со стены и посылает его к vk api
func getWallPostsForComments(sender string, subject Subject,
	wallPostCommentMonitorParam WallPostCommentMonitorParam) ([]WallPost, error) {

	var wallPosts []WallPost

	// формируем карту с параметрами запроса
	jsonDump := fmt.Sprintf(`{
			"owner_id": "%d",
			"count": "%d",
			"filter": "%v",
			"v": "5.95"
		}`, subject.SubjectID, wallPostCommentMonitorParam.PostsCount,
		wallPostCommentMonitorParam.Filter)
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

// ParseWallPostsVkAPIMap
// функция парсинга карты с данными о постах уже описана в модуле мониторинга новых постов,
// поэтому тут этот момент опустим

// WallPostComment - структура для данных о комментариях под постами
type WallPostComment struct {
	ID          int          `json:"id"`
	ThreadID    int          `json:"parents_stack"`
	PostOwnerID int          `json:"owner_id"`
	PostID      int          `json:"post_id"`
	FromID      int          `json:"from_id"`
	Date        int          `json:"date"`
	Text        string       `json:"text"`
	Attachments []Attachment `json:"attachments"`
}

// getWallPostComments формирует запрос на получение комментариев под постами со стены и посылает его к vk api
func getWallPostComments(sender string, subject Subject,
	wallPostCommentMonitorParam WallPostCommentMonitorParam,
	wallPosts []WallPost) ([]WallPostComment, error) {
	var wallPostsComments []WallPostComment

	for _, wallPost := range wallPosts {

		// формируем карту с параметрами запроса
		jsonDump := fmt.Sprintf(`{
			"owner_id": "%d",
			"post_id": "%d",
			"count": "%d",
			"sort": "desc",
			"v": "5.84"
		}`, wallPost.OwnerID, wallPost.ID, wallPostCommentMonitorParam.CommentsCount)
		values, err := MakeJSON(jsonDump)
		if err != nil {
			return wallPostsComments, err
		}

		// отправляем запрос, получаем ответ
		response, err := SendVKAPIQuery(sender, "wall.getComments", values, subject)
		if err != nil {
			return wallPostsComments, err
		}

		wallPostComment := parseWallPostCommentVkAPIMap(response["response"].(map[string]interface{}), wallPost)
		wallPostsComments = append(wallPostsComments, wallPostComment...)
	}

	return wallPostsComments, nil
}

// parseWallPostCommentVkAPIMap извлекает данные о комментариях из полученной карты vk api
func parseWallPostCommentVkAPIMap(resp map[string]interface{}, wallPost WallPost) []WallPostComment {
	var wallPostComments []WallPostComment

	// перебираем элементы с данными о комментариях
	itemsMap := resp["items"].([]interface{})
	for _, itemMap := range itemsMap {
		item := itemMap.(map[string]interface{})

		// парсим данные из элемента в структуру
		var wallPostComment WallPostComment
		wallPostComment.ID = int(item["id"].(float64))
		wallPostComment.PostID = wallPost.ID
		wallPostComment.PostOwnerID = wallPost.OwnerID
		wallPostComment.FromID = int(item["from_id"].(float64))
		wallPostComment.Date = int(item["date"].(float64))
		wallPostComment.Text = item["text"].(string)

		// если есть прикрепления, то вызываем парсер прикреплений
		if mediaContent, exist := item["attachments"]; exist == true {
			wallPostComment.Attachments = ParseAttachments(mediaContent.([]interface{}))
		}
		wallPostComments = append(wallPostComments, wallPostComment)
	}
	return wallPostComments
}

// getWallPostComment формирует запрос на получение комментария под постом со стены и посылает его к vk api
func getWallPostComment(sender string, subject Subject,
	wallPostCommentMonitorParam WallPostCommentMonitorParam,
	wallPostComments []WallPostComment) ([]WallPostComment, error) {
	var updatedWallPostComments []WallPostComment

	for _, wallPostComment := range wallPostComments {

		// формируем карту с параметрами запроса
		jsonDump := fmt.Sprintf(`{
			"owner_id": "%d",
			"comment_id": "%d",
			"v": "5.95"
		}`, wallPostComment.PostOwnerID, wallPostComment.ID)
		values, err := MakeJSON(jsonDump)
		if err != nil {
			return wallPostComments, err
		}

		// отправляем запрос, получаем ответ
		response, err := SendVKAPIQuery(sender, "wall.getComment", values, subject)
		if err != nil {
			return wallPostComments, err
		}

		setParentsStack(response["response"].(map[string]interface{}), &wallPostComment)

		updatedWallPostComments = append(updatedWallPostComments, wallPostComment)
	}

	return updatedWallPostComments, nil
}

// setParentsStack извлекает из карты с данными о комментарии идентификатор ветки комментариев
func setParentsStack(resp map[string]interface{}, wallPostComment *WallPostComment) {

	// перебираем элементы с данными о комментариях
	itemsMap := resp["items"].([]interface{})
	for _, itemMap := range itemsMap {
		item := itemMap.(map[string]interface{})

		// проверяем длину списка с идентификаторами веток
		parentsStack := item["parents_stack"].([]interface{})

		// если элементы присутствуют, то добавляем оттуда данные в структуру
		if len(parentsStack) > 0 {
			wallPostComment.ThreadID = int(parentsStack[0].(float64))
		}
	}
}

// makeMessageWallPostComment собирает сообщение с данными о комментарии под постом
func makeMessageWallPostComment(sender string, subject Subject,
	wallPostCommentMonitorParam WallPostCommentMonitorParam, wallPostComment WallPostComment) (string, error) {
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
	if wallPostComment.FromID < 0 {
		vkCommunity, err := GetCommunityInfo(sender, subject, wallPostComment.FromID)
		authorHyperlink = MakeCommunityHyperlink(vkCommunity)
		if err != nil {
			return "", err
		}
	} else {
		vkUser, err := GetUserInfo(sender, subject, wallPostComment.FromID)
		authorHyperlink = MakeUserHyperlink(vkUser)
		if err != nil {
			return "", err
		}
	}

	// и когда был создан
	creationDate := UnixTimeStampToDate(wallPostComment.Date)

	// формируем строку с прикреплениями
	var attachments string
	var link string
	if len(wallPostComment.Attachments) > 0 {
		for _, attachment := range wallPostComment.Attachments {
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

	// собираем ссылку на комментарий
	commentURL := fmt.Sprintf("https://vk.com/wall%d_%d?reply=%d",
		wallPostComment.PostOwnerID, wallPostComment.PostID, wallPostComment.ID)

	// если коммент находится в ветке, то добавляем id ветки
	if wallPostComment.ThreadID > 0 {
		commentURL += fmt.Sprintf("&thread=%d", wallPostComment.ThreadID)
	}

	// добавляем подготовленные фрагменты сообщения в общий текст
	// сначала сигнатуру
	text := fmt.Sprintf("New post comment\\nLocation: %v\\nAuthor: %v\\nCreated: %v",
		locationHyperlink, authorHyperlink, creationDate)

	// затем основной текст комментария, если он есть
	if len(wallPostComment.Text) > 0 {
		// но сначала обрезаем его из-за ограничения на длину запроса
		if len(wallPostComment.Text) > 800 {
			wallPostComment.Text = string(wallPostComment.Text[0:800]) + "\\n[long_text]"
		}
		// и экранируем все символы пропуска строки, потому что у json.Unmarshal с ними проблемы
		wallPostComment.Text = strings.Replace(wallPostComment.Text, "\n", "\\n", -1)
		text += fmt.Sprintf("\\n\\n%v", wallPostComment.Text)
	}

	// потом прикрепленную ссылку, если она есть
	if len(link) > 0 {
		if len(wallPostComment.Text) == 0 {
			text += "\\n"
		}
		text += fmt.Sprintf("\\n%v", link)
	}

	// и ссылку на сам комментарий
	text += fmt.Sprintf("\\n\\n%v", commentURL)

	// экранируем все апострофы, чтобы не сломали нам json.Unmarshal
	text = strings.Replace(text, `"`, `\"`, -1)

	// далее формируем строку с данными для карты
	jsonDump := fmt.Sprintf(`{
		"peer_id": "%d",
		"message": "%v",
		"v": "5.68"
	`, wallPostCommentMonitorParam.SendTo, text)

	// если прикрепления были, то добавляем их в json-строку отдельно
	if len(attachments) > 0 {
		jsonDump += fmt.Sprintf(`, "attachment": "%v"`, attachments)
	}

	// закрываем карту
	jsonDump += "}"

	return jsonDump, nil
}
