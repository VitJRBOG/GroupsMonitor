package main

import (
	"fmt"
	"strconv"
	"strings"
)

// MonitorModule - хранит данные о модулях мониторов
type MonitorModule struct {
	Name           string
	SubjectID      int
	SendTo         int
	NeedMonitoring int
	Filter         string
	Interval       int
	MainCount      int
	SecondCount    int
}

// CollectionMonitorData собирает данные обо всех мониторах
func CollectionMonitorData(subjectID int) (*[]MonitorModule, error) {
	monitorsNames := []string{"wall_post_monitor", "album_photo_monitor",
		"video_monitor", "photo_comment_monitor", "video_comment_monitor",
		"topic_monitor", "wall_post_comment_monitor"}

	var monitors []MonitorModule

	// запускаем перебор названий мониторов и сбор данных для каждого
	for _, monitorName := range monitorsNames {
		var monitor MonitorModule
		monitor.Name = monitorName
		monitor.SubjectID = subjectID

		// получаем идентификатор получателя сообщений
		sendTo, err := getSendTo(monitorName)
		if err != nil {
			return nil, err
		}
		monitor.SendTo = sendTo

		// получаем значение флажка для монитора
		needMonitoring, err := getNeedMonitoring(monitorName)
		if err != nil {
			return nil, err
		}
		monitor.NeedMonitoring = needMonitoring

		// если запуск монитора разрешен, то получаем значение фильтра проверки
		if needMonitoring == 1 {
			filter, err := getFilter(monitorName)
			if err != nil {
				return nil, err
			}
			monitor.Filter = filter
		} else {
			// если нет, то присваиваем значение по умолчанию
			monitor.Filter = "all"
		}

		// если запуск монитора разрешен,
		//то получаем значение интервала проверки
		if needMonitoring == 1 {
			interval, err := getInterval(monitorName)
			if err != nil {
				return nil, err
			}
			monitor.Interval = interval
		} else {
			// если нет, то присваиваем значение по умолчанию
			monitor.Interval = 60
		}

		// если запуск монитора разрешен,
		//то получаем количество проверяемых объектов
		if needMonitoring == 1 {
			mainCount, err := getCount(monitorName, "main")
			if err != nil {
				return nil, err
			}
			monitor.MainCount = mainCount
		} else {
			// если нет, то присваиваем значение по умолчанию
			monitor.MainCount = 5
		}

		// проверяем текущее название монитора и, если требуется,
		//получаем дополнительное значение количества проверяемых объектов
		switch monitorName {
		case "video_comment_monitor":
			if needMonitoring == 1 {
				secondCount, err := getCount(monitorName, "second")
				if err != nil {
					return nil, err
				}
				monitor.SecondCount = secondCount
			} else {
				monitor.SecondCount = 5
			}
		case "topic_monitor":
			if needMonitoring == 1 {
				secondCount, err := getCount(monitorName, "second")
				if err != nil {
					return nil, err
				}
				monitor.SecondCount = secondCount
			} else {
				monitor.SecondCount = 5
			}
		case "wall_post_comment_monitor":
			if needMonitoring == 1 {
				secondCount, err := getCount(monitorName, "second")
				if err != nil {
					return nil, err
				}
				monitor.SecondCount = secondCount
			} else {
				monitor.SecondCount = 5
			}
		}

		monitors = append(monitors, monitor)
	}

	return &monitors, nil
}

// getSendTo запрашивает у пользователя идентификатор получателя сообщений в ВК
func getSendTo(monitorName string) (int, error) {
	monitorName = strings.ReplaceAll(monitorName, "_", " ")
	monitorName = strings.Replace(monitorName, string(monitorName[0]),
		strings.ToUpper(string(monitorName[0])), 1)
	sender := fmt.Sprintf("> [%v -> Get id of send to]: ", monitorName)
	fmt.Print(sender)

	var userAnswer string
	_, err := fmt.Scan(&userAnswer)
	if err != nil {
		return 0, err
	}
	sendTo, err := strconv.Atoi(userAnswer)
	if err != nil {
		return 0, err
	}
	return sendTo, nil
}

// getNeedMonitoring запрашивает у пользователя флаг для активации монитора
func getNeedMonitoring(monitorName string) (int, error) {
	monitorName = strings.ReplaceAll(monitorName, "_", " ")
	monitorName = strings.Replace(monitorName, string(monitorName[0]),
		strings.ToUpper(string(monitorName[0])), 1)
	sender := fmt.Sprintf("> [%v -> Get flag of need monitoring]: ", monitorName)
	fmt.Print(sender)

	var userAnswer string
	_, err := fmt.Scan(&userAnswer)
	if err != nil {
		return 0, err
	}
	needMonitoring, err := strconv.Atoi(userAnswer)
	if err != nil {
		return 0, err
	}
	return needMonitoring, nil
}

// getFilter запрашивает у пользователя значение для фильтра проверки
func getFilter(monitorName string) (string, error) {
	monitorName = strings.ReplaceAll(monitorName, "_", " ")
	monitorName = strings.Replace(monitorName, string(monitorName[0]),
		strings.ToUpper(string(monitorName[0])), 1)
	sender := fmt.Sprintf("> [%v -> Get filter]: ", monitorName)
	fmt.Print(sender)

	var userAnswer string
	_, err := fmt.Scan(&userAnswer)
	if err != nil {
		return "", err
	}
	return userAnswer, nil
}

// getInterval запрашивает у пользователя значение интервала проверки
func getInterval(monitorName string) (int, error) {
	monitorName = strings.ReplaceAll(monitorName, "_", " ")
	monitorName = strings.Replace(monitorName, string(monitorName[0]),
		strings.ToUpper(string(monitorName[0])), 1)
	sender := fmt.Sprintf("> [%v -> Get interval]: ", monitorName)
	fmt.Print(sender)

	var userAnswer string
	_, err := fmt.Scan(&userAnswer)
	if err != nil {
		return 0, err
	}
	interval, err := strconv.Atoi(userAnswer)
	if err != nil {
		return 0, err
	}
	return interval, nil
}

// getCount запрашивает у пользователя значение количества проверяемых объектов
func getCount(monitorName string, typeCount string) (int, error) {
	// определяем назначение цифры количества, которую будем запрашивать у пользователя
	switch monitorName {
	case "wall_post_monitor":
		typeCount = "posts count"
	case "album_photo_monitor":
		typeCount = "photos count"
	case "video_monitor":
		typeCount = "video count"
	case "photo_comment_monitor":
		typeCount = "comments count"
	case "video_comment_monitor":
		typeCount = "videos count"
	case "topic_monitor":
		switch typeCount {
		case "main":
			typeCount = "topics count"
		case "second":
			typeCount = "comments count"
		}
	case "wall_post_comment_monitor":
		switch typeCount {
		case "main":
			typeCount = "posts count"
		case "second":
			typeCount = "comments count"
		}
	}

	// делаем строчку с названием монитора более приятной глазу
	monitorName = strings.ReplaceAll(monitorName, "_", " ")
	monitorName = strings.Replace(monitorName, string(monitorName[0]),
		strings.ToUpper(string(monitorName[0])), 1)

	sender := fmt.Sprintf("> [%v -> Get %v]: ", monitorName, typeCount)
	fmt.Print(sender)

	var userAnswer string
	_, err := fmt.Scan(&userAnswer)
	if err != nil {
		return 0, err
	}
	count, err := strconv.Atoi(userAnswer)
	if err != nil {
		return 0, err
	}
	return count, nil
}
