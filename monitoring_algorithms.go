package main

import (
	"fmt"
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
	wallPosts, err := GetWallPosts(sender, subject, wallPostMonitorParam)
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

			// формируем строку с данными для карты для отправки сообщения
			messageParameters, err := MakeMessageWallPost(sender, subject, wallPostMonitorParam, wallPost)
			if err != nil {
				return err
			}

			// отправляем сообщение с полученными данными
			if err := SendMessage(sender, messageParameters, subject); err != nil {
				return err
			}

			// обновляем дату последнего проверенного поста в БД
			if err := UpdateDBWallPostMonitorLastDate(subject.ID, wallPost.Date, wallPostMonitorParam); err != nil {
				return err
			}
		}
	}

	return nil
}

// func GetAlbumPhotos() {}
// func GetVideos() {}
// func GetPhotosComments() {}
// func GetVideosComments() {}
// func GetTopicPosts() {}
// func WallPostsComments() {}
//    чтобы достать комменты из веток,
//    можно запрашивать все комменты из поста старой версией метода wall.getComments,
//    а потом каждый пробивать новой версией wall.getComment по полученному id
//
//    из-за ограничений vk api нет иного способа достать 10-го коммента из ветки
