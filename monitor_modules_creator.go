package main

// createMonitorModules запускает алгоритмы создания модулей мониторинга
func createMonitorModules(monitorModule MonitorModule) error {
	switch monitorModule.Name {
	case "wall_post_monitor":
		err := createWallPostMonitor(monitorModule)
		if err != nil {
			return err
		}
	case "album_photo_monitor":
		err := createAlbumPhotoMonitor(monitorModule)
		if err != nil {
			return err
		}
	case "video_monitor":
		err := createVideoMonitor(monitorModule)
		if err != nil {
			return err
		}
	case "photo_comment_monitor":
		err := createPhotoCommentMonitor(monitorModule)
		if err != nil {
			return err
		}
	case "video_comment_monitor":
		err := createVideoCommentMonitor(monitorModule)
		if err != nil {
			return err
		}
	case "topic_monitor":
		err := createTopicMonitor(monitorModule)
		if err != nil {
			return err
		}
	case "wall_post_comment_monitor":
		err := createWallPostCommentMonitor(monitorModule)
		if err != nil {
			return err
		}
	}

	return nil
}

// createWallPostMonitor добавляет данные о новом модуле мониторинга wall_post_monitor
func createWallPostMonitor(monitorModule MonitorModule) error {
	var wallPostMonitorParam WallPostMonitorParam

	wallPostMonitorParam.SubjectID = monitorModule.SubjectID
	wallPostMonitorParam.NeedMonitoring = monitorModule.NeedMonitoring
	wallPostMonitorParam.Interval = monitorModule.Interval
	wallPostMonitorParam.SendTo = monitorModule.SendTo
	wallPostMonitorParam.Filter = monitorModule.Filter
	wallPostMonitorParam.LastDate = 0 // TODO лучше ставить текущую дату
	wallPostMonitorParam.PostsCount = monitorModule.MainCount
	wallPostMonitorParam.KeywordsForMonitoring = `{"list":[]}`
	wallPostMonitorParam.UsersIDsForIgnore = `{"list":[]}`

	err := InsertDBWallPostMonitor(wallPostMonitorParam)
	if err != nil {
		return err
	}

	return nil
}

// createAlbumPhotoMonitor добавляет данные о новом модуле мониторинга album_photo_monitor
func createAlbumPhotoMonitor(monitorModule MonitorModule) error {
	var albumPhotoMonitorParam AlbumPhotoMonitorParam

	albumPhotoMonitorParam.SubjectID = monitorModule.SubjectID
	albumPhotoMonitorParam.NeedMonitoring = monitorModule.NeedMonitoring
	albumPhotoMonitorParam.SendTo = monitorModule.SendTo
	albumPhotoMonitorParam.Interval = monitorModule.Interval
	albumPhotoMonitorParam.LastDate = 0 // TODO лучше ставить текущую дату
	albumPhotoMonitorParam.PhotosCount = monitorModule.MainCount

	err := InsertDBAlbumPhotoMonitor(albumPhotoMonitorParam)
	if err != nil {
		return err
	}

	return nil
}

// createVideoMonitor добавляет данные о новом модуле мониторинга video_monitor
func createVideoMonitor(monitorModule MonitorModule) error {
	var videoMonitorParam VideoMonitorParam

	videoMonitorParam.SubjectID = monitorModule.SubjectID
	videoMonitorParam.NeedMonitoring = monitorModule.NeedMonitoring
	videoMonitorParam.SendTo = monitorModule.SendTo
	videoMonitorParam.Interval = monitorModule.Interval
	videoMonitorParam.LastDate = 0 // TODO лучше ставить текущую дату
	videoMonitorParam.VideoCount = monitorModule.MainCount

	err := InsertDBVideoMonitor(videoMonitorParam)
	if err != nil {
		return err
	}

	return nil
}

// createPhotoCommentMonitor добавляет данные о новом модуле мониторинга photo_comment_monitor
func createPhotoCommentMonitor(monitorModule MonitorModule) error {
	var photoCommentMonitorParam PhotoCommentMonitorParam

	photoCommentMonitorParam.SubjectID = monitorModule.SubjectID
	photoCommentMonitorParam.NeedMonitoring = monitorModule.NeedMonitoring
	photoCommentMonitorParam.CommentsCount = monitorModule.MainCount
	photoCommentMonitorParam.LastDate = 0 // TODO лучше ставить текущую дату
	photoCommentMonitorParam.Interval = monitorModule.Interval
	photoCommentMonitorParam.SendTo = monitorModule.SendTo

	err := InsertDBPhotoCommentMonitor(photoCommentMonitorParam)
	if err != nil {
		return err
	}

	return nil
}

// createPhotoCommentMonitor добавляет данные о новом модуле мониторинга video_comment_monitor
func createVideoCommentMonitor(monitorModule MonitorModule) error {
	var videoCommentMonitorParam VideoCommentMonitorParam

	videoCommentMonitorParam.SubjectID = monitorModule.SubjectID
	videoCommentMonitorParam.NeedMonitoring = monitorModule.NeedMonitoring
	videoCommentMonitorParam.VideosCount = monitorModule.MainCount
	videoCommentMonitorParam.Interval = monitorModule.Interval
	videoCommentMonitorParam.CommentsCount = monitorModule.SecondCount
	videoCommentMonitorParam.SendTo = monitorModule.SendTo
	videoCommentMonitorParam.LastDate = 0 // TODO лучше ставить текущую дату

	err := InsertDBVideoCommentMonitor(videoCommentMonitorParam)
	if err != nil {
		return err
	}

	return nil
}

// createTopicMonitor добавляет данные о новом модуле мониторинга topic_monitor
func createTopicMonitor(monitorModule MonitorModule) error {
	var topicMonitorParam TopicMonitorParam

	topicMonitorParam.SubjectID = monitorModule.SubjectID
	topicMonitorParam.NeedMonitoring = monitorModule.NeedMonitoring
	topicMonitorParam.TopicsCount = monitorModule.MainCount
	topicMonitorParam.CommentsCount = monitorModule.SecondCount
	topicMonitorParam.Interval = monitorModule.Interval
	topicMonitorParam.SendTo = monitorModule.SendTo
	topicMonitorParam.LastDate = 0 // TODO лучше ставить текущую дату

	err := InsertDBTopicMonitor(topicMonitorParam)
	if err != nil {
		return err
	}

	return nil
}

func createWallPostCommentMonitor(monitorModule MonitorModule) error {
	var wallPostCommentMonitorParam WallPostCommentMonitorParam

	wallPostCommentMonitorParam.SubjectID = monitorModule.SubjectID
	wallPostCommentMonitorParam.NeedMonitoring = monitorModule.NeedMonitoring
	wallPostCommentMonitorParam.PostsCount = monitorModule.MainCount
	wallPostCommentMonitorParam.CommentsCount = monitorModule.SecondCount
	wallPostCommentMonitorParam.MonitoringAll = 1
	wallPostCommentMonitorParam.UsersIDsForMonitoring = `{"list":[]}`
	wallPostCommentMonitorParam.UsersNamesForMonitoring = `{"list":[]}`
	wallPostCommentMonitorParam.AttachmentsTypesForMonitoring = `{"list":["photo", "video", "audio", "doc", "poll", "link"]}`
	wallPostCommentMonitorParam.UsersIDsForIgnore = `{"list":[]}`
	wallPostCommentMonitorParam.Interval = monitorModule.Interval
	wallPostCommentMonitorParam.SendTo = monitorModule.SendTo
	wallPostCommentMonitorParam.Filter = monitorModule.Filter
	wallPostCommentMonitorParam.LastDate = 0 // TODO лучше ставить текущую дату
	wallPostCommentMonitorParam.KeywordsForMonitoring = `{"list":[]}`
	wallPostCommentMonitorParam.SmallCommentsForMonitoring = `{"list":[]}`
	wallPostCommentMonitorParam.DigitsCountForCardNumberMonitoring = `{"list":["16"]}`
	wallPostCommentMonitorParam.DigitsCountForPhoneNumberMonitoring = `{"list":["6","11"]}`
	wallPostCommentMonitorParam.MonitorByCommunity = 1

	err := InsertDBWallPostCommentMonitor(wallPostCommentMonitorParam)
	if err != nil {
		return err
	}

	return nil
}
