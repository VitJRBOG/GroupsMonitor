package main

import (
	"fmt"
	"strings"
)

// AlbumPhotoMonitor проверяет фотографии в альбомах
func AlbumPhotoMonitor(subject Subject) error {
	sender := fmt.Sprintf("%v's album photo monitoring", subject.Name)

	// запрашиваем структуру с параметрами модуля мониторинга фотографий
	albumPhotoMonitorParam, err := SelectDBAlbumPhotoMonitorParam(subject.ID)
	if err != nil {
		return err
	}

	// запрашиваем структуру с альбомами
	albums, err := getAlbums(sender, subject)
	if err != nil {
		return err
	}

	// запрашиваем фотографии из этих альбомов
	albumsPhotos, err := getAlbumsPhotos(sender, subject, albumPhotoMonitorParam, albums)
	if err != nil {
		return err
	}

	var targetAlbumsPhotos []AlbumPhoto

	// отфильтровываем старые
	for _, albumPhoto := range albumsPhotos {
		if albumPhoto.Date > albumPhotoMonitorParam.LastDate {
			targetAlbumsPhotos = append(targetAlbumsPhotos, albumPhoto)
		}
	}

	// если после фильтрации что-то осталось, продолжаем
	if len(targetAlbumsPhotos) > 0 {

		// сортируем в порядке от раннего к позднему
		for j := 0; j < len(targetAlbumsPhotos); j++ {
			f := 0
			for i := 0; i < len(targetAlbumsPhotos)-1-j; i++ {
				if targetAlbumsPhotos[i].Date > targetAlbumsPhotos[i+1].Date {
					x := targetAlbumsPhotos[i]
					y := targetAlbumsPhotos[i+1]
					targetAlbumsPhotos[i+1] = x
					targetAlbumsPhotos[i] = y
					f = 1
				}
			}
			if f == 0 {
				break
			}
		}

		// перебираем отсортированный список
		for _, albumPhoto := range targetAlbumsPhotos {

			// формируем строку с данными для карты для отправки сообщения
			messageParameters, err := makeMessageAlbumPhoto(sender, subject, albumPhotoMonitorParam, albumPhoto)
			if err != nil {
				return err
			}

			// отправляем сообщение с полученными данными
			if err := SendMessage(sender, messageParameters, subject); err != nil {
				return err
			}

			// выводим в консоль сообщение о новой фотографии
			outputReportAboutNewAlbumPhoto(sender, albumPhoto)

			// обновляем дату последнего проверенного поста в БД
			if err := UpdateDBAlbumPhotoMonitorLastDate(subject.ID, albumPhoto.Date,
				albumPhotoMonitorParam); err != nil {
				return err
			}
		}
	}

	return nil
}

// outputReportAboutNewAlbumPhoto выводит сообщение о новой фотографии
func outputReportAboutNewAlbumPhoto(sender string, albumPhoto AlbumPhoto) {
	creationDate := UnixTimeStampToDate(albumPhoto.Date)
	message := fmt.Sprintf("New photo in %v at %v.", albumPhoto.AlbumName, creationDate)
	OutputMessage(sender, message)
}

// Album - структура для данных об альбоме
type Album struct {
	ID      int    `json:"id"`
	OwnerID int    `json:"owner_id"`
	Name    string `json:"title"`
}

// getAlbums формирует запрос на получение альбомов и посылает его к vk api
func getAlbums(sender string, subject Subject) ([]Album, error) {
	var albums []Album

	// формируем карту с параметрами запроса
	jsonDump := fmt.Sprintf(`{
		"owner_id": "%d",
		"v": "5.95"
	}`, subject.SubjectID)
	values, err := MakeJSON(jsonDump)
	if err != nil {
		return albums, err
	}

	// отправляем запрос, получаем ответ
	response, err := SendVKAPIQuery(sender, "photos.getAlbums", values, subject)
	if err != nil {
		return albums, err
	}

	// парсим полученные данные об альбомах
	albums = parseAlbumVkAPIMap(response["response"].(map[string]interface{}))

	return albums, nil
}

// parseAlbumsVkAPIMap извлекает данные об альбомах из полученной карты vk api
func parseAlbumVkAPIMap(resp map[string]interface{}) []Album {
	var albums []Album

	// перебираем элементы с данными об альбомах
	itemsMap := resp["items"].([]interface{})
	for _, itemMap := range itemsMap {
		item := itemMap.(map[string]interface{})

		// парсим данные из элемента в структуру
		var album Album
		album.ID = int(item["id"].(float64))
		album.OwnerID = int(item["owner_id"].(float64))
		album.Name = item["title"].(string)

		albums = append(albums, album)
	}

	return albums
}

// AlbumPhoto - структура для данных о фотографии
type AlbumPhoto struct {
	ID        int    `json:"id"`
	OwnerID   int    `json:"owner_id"`
	UserID    int    `json:"user_id"`
	AlbumID   int    `json:"album_id"`
	AlbumName string `json:"title"`
	Text      string `json:"text"`
	Date      int    `json:"date"`
}

// getAlbumPhotos формирует запрос на получение фотографий из альбомов и посылает его к vk api
func getAlbumsPhotos(sender string, subject Subject,
	albumPhotoMonitorParam AlbumPhotoMonitorParam, albums []Album) ([]AlbumPhoto, error) {
	var albumsPhotos []AlbumPhoto

	for _, album := range albums {

		// формируем карту с параметрами запроса
		jsonDump := fmt.Sprintf(`{
			"owner_id": "%d",
			"album_id": "%d",
			"rev": "1",
			"count": "%d",
			"v": "5.95"
		}`, album.OwnerID, album.ID, albumPhotoMonitorParam.PhotosCount)
		values, err := MakeJSON(jsonDump)
		if err != nil {
			return albumsPhotos, err
		}

		// отправляем запрос, получаем ответ
		response, err := SendVKAPIQuery(sender, "photos.get", values, subject)
		if err != nil {
			return albumsPhotos, err
		}

		albumPhotos := parseAlbumPhotoVkAPIMap(response["response"].(map[string]interface{}), album)
		albumsPhotos = append(albumsPhotos, albumPhotos...)
	}

	return albumsPhotos, nil
}

// parseAlbumPhotoVkAPIMap извлекает данные о фотографиях из полученной карты vk api
func parseAlbumPhotoVkAPIMap(resp map[string]interface{}, album Album) []AlbumPhoto {
	var albumPhotos []AlbumPhoto

	// перебираем элементы с данными о фотографиях
	itemsMap := resp["items"].([]interface{})
	for _, itemMap := range itemsMap {
		item := itemMap.(map[string]interface{})

		// парсим данные из элемента в структуру
		var albumPhoto AlbumPhoto
		albumPhoto.ID = int(item["id"].(float64))
		albumPhoto.AlbumID = int(item["album_id"].(float64))
		albumPhoto.OwnerID = int(item["owner_id"].(float64))
		albumPhoto.UserID = int(item["user_id"].(float64))
		albumPhoto.Date = int(item["date"].(float64))
		albumPhoto.Text = item["text"].(string)
		albumPhoto.AlbumName = album.Name

		albumPhotos = append(albumPhotos, albumPhoto)
	}

	return albumPhotos
}

// makeMessageAlbumPhoto собирает сообщение с данными о фотографии
func makeMessageAlbumPhoto(sender string, subject Subject,
	albumPhotoMonitorParam AlbumPhotoMonitorParam, albumPhoto AlbumPhoto) (string, error) {

	// собираем данные для сигнатуры сообщения:
	// название альбома
	albumName := albumPhoto.AlbumName

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

	// кто загрузил, если это не пользователь id100
	var authorHyperlink string
	if albumPhoto.UserID == 100 {
		authorHyperlink = locationHyperlink
	} else {
		vkUser, err := GetUserInfo(sender, subject, albumPhoto.UserID)
		authorHyperlink = MakeUserHyperlink(vkUser)
		if err != nil {
			return "", err
		}
	}

	// и когда загрузил
	creationDate := UnixTimeStampToDate(albumPhoto.Date)

	// формируем строку с прикреплением
	attachment := fmt.Sprintf("photo%d_%d", albumPhoto.OwnerID, albumPhoto.ID)

	// собираем ссылку на фото
	photoURL := fmt.Sprintf("https://vk.com/photo%d_%d", albumPhoto.OwnerID, albumPhoto.ID)

	// добавляем подготовленные фрагменты сообщения в общий текст
	// сначала сигнатуру
	text := fmt.Sprintf("New album photo\\nAlbum: %v\\nLocation: %v\\nAuthor: %v\\nCreated: %v",
		albumName, locationHyperlink, authorHyperlink, creationDate)

	// затем описание к фотографии, если оно есть
	if len(albumPhoto.Text) > 0 {
		// но сначала обрезаем его из-за ограничения на длину запроса
		if len(albumPhoto.Text) > 800 {
			albumPhoto.Text = string(albumPhoto.Text[0:800]) + "\\n[long_text]"
		}
		// и экранируем все символы пропуска строки, потому что у json.Unmarshal с ними проблемы
		albumPhoto.Text = strings.Replace(albumPhoto.Text, "\n", "\\n", -1)
		text += fmt.Sprintf("\\n\\n%v", albumPhoto.Text)
	}

	// и добавляем ссылку на саму фотку
	text += fmt.Sprintf("\\n\\n%v", photoURL)

	// экранируем все апострофы, чтобы не сломали нам json.Unmarshal
	text = strings.Replace(text, `"`, `\"`, -1)
	// экранируем все обратные слэши, чтобы также не сломали json.Unmarshal
	text = strings.Replace(text, `\`, `\\`, -1)
	// и возвращаем символы пропуска строки после экранировки обратных слэшей
	text = strings.Replace(text, `\\n`, `\n`, -1)

	// далее формируем строку с данными для карты
	jsonDump := fmt.Sprintf(`{
		"peer_id": "%d",
		"message": "%v",
		"attachment": "%v",
		"v": "5.68"
	}`, albumPhotoMonitorParam.SendTo, text, attachment)

	return jsonDump, nil
}
