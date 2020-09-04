package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3" // иначе драйвер для работы с SQLite не работает
)

// openDB читает и возвращает базу данных
func openDB() (*sql.DB, error) {
	// определяем путь к файлу базы данных
	path, err := ReadPathFile()
	if err != nil {
		return nil, err
	}
	pathToDB := path + "groupsmonitor_db.db"
	// читаем db sqlite
	db, err := sql.Open("sqlite3", pathToDB)
	if err != nil {
		// если словил ошибку, то преобразуем ее в текст
		errorMessage := fmt.Sprintf("%v", err)

		// и отправляем в обработчик ошибок
		typeError, causeError := DBIOError(errorMessage)

		// потом проверяем тип полученной ошибки
		switch typeError {

		// задержка, если получена ошибка, которую можно решить таким образом
		case "timeout error":
			interval := 1 * time.Second
			sender := "Database"
			message := fmt.Sprintf("Error: %v. Timeout for %v...", causeError, interval)
			OutputMessage(sender, message)
			time.Sleep(interval)
			return openDB()
		}

		// если пойманная ошибка не обрабатывается, то возвращаем ее стандартным путем
		return nil, err
	}

	return db, nil
}

// GetTableLen определяет количество записей в указанной таблице
func GetTableLen(tableName string) (int, error) {
	db, err := openDB()
	defer db.Close()
	if err != nil {
		return 0, err
	}

	// читаем данные из БД
	query := fmt.Sprintf("SELECT COUNT(*) FROM '%v'", tableName)
	rows, err := db.Query(query)
	defer rows.Close()
	if err != nil {
		return 0, err
	}

	var tableLenth int

	for rows.Next() {
		err = rows.Scan(&tableLenth)
		if err != nil {
			return 0, err
		}
	}
	return tableLenth, nil
}

// AccessToken - структура для полей из таблицы access_token
type AccessToken struct {
	ID    int
	Name  string
	Value string
}

// InsertDBAccessToken добавляет новое поле в таблицу access_token
func InsertDBAccessToken(accessToken AccessToken) error {
	// получаем ссылку на db
	db, err := openDB()
	defer db.Close()
	if err != nil {
		return err
	}

	// добавляем новое поле в таблицу
	query := fmt.Sprintf(`INSERT INTO access_token (name, value) VALUES ('%v', '%v')`,
		accessToken.Name, accessToken.Value)
	_, err = db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

// SelectDBAccessTokens извлекает поля из таблицы access_token
func SelectDBAccessTokens() ([]AccessToken, error) {
	// получаем ссылку на db
	db, err := openDB()
	defer db.Close()
	if err != nil {
		return nil, err
	}

	// читаем данные из БД
	rows, err := db.Query("SELECT * FROM access_token")
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	// считываем данные из rows
	var accessTokens []AccessToken
	for rows.Next() {
		var accessToken AccessToken
		err = rows.Scan(&accessToken.ID, &accessToken.Name, &accessToken.Value)
		if err != nil {
			return nil, err
		}
		accessTokens = append(accessTokens, accessToken)
	}
	return accessTokens, nil
}

// SelectDBAccessTokenByID извлекает поле из таблицы access_token по id
func SelectDBAccessTokenByID(accessTokenID int) (AccessToken, error) {
	var accessToken AccessToken
	// получаем ссылку на db
	db, err := openDB()
	defer db.Close()
	if err != nil {
		return accessToken, err
	}

	// читаем данные из БД
	query := fmt.Sprintf("SELECT * FROM access_token WHERE id=%d", accessTokenID)
	rows, err := db.Query(query)
	defer rows.Close()
	if err != nil {
		return accessToken, err
	}

	// считываем данные из rows
	for rows.Next() {
		err = rows.Scan(&accessToken.ID, &accessToken.Name, &accessToken.Value)
		if err != nil {
			return accessToken, err
		}
	}
	return accessToken, nil
}

// SelectDBAccessTokenByName извлекает поле из таблицы access_token по name
func SelectDBAccessTokenByName(nameAccessToken string) (AccessToken, error) {
	var accessToken AccessToken
	// получаем ссылку на db
	db, err := openDB()
	defer db.Close()
	if err != nil {
		return accessToken, err
	}

	// читаем данные из БД
	query := fmt.Sprintf("SELECT * FROM access_token WHERE name='%v'", nameAccessToken)
	rows, err := db.Query(query)
	defer rows.Close()
	if err != nil {
		return accessToken, err
	}

	// считываем данные из rows
	for rows.Next() {
		err = rows.Scan(&accessToken.ID, &accessToken.Name, &accessToken.Value)
		if err != nil {
			return accessToken, err
		}
	}
	return accessToken, nil
}

// UpdateDBAccessToken обновляет значения в поле таблицы access_token
func UpdateDBAccessToken(accessToken AccessToken) error {
	// получаем ссылку на db
	db, err := openDB()
	defer db.Close()
	if err != nil {
		return err
	}

	// обновляем значения в конкретном поле
	query := fmt.Sprintf(`UPDATE access_token SET name='%v', value='%v' WHERE id=%d`,
		accessToken.Name, accessToken.Value, accessToken.ID)
	_, err = db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

// Subject - структура для полей из таблицы subject
type Subject struct {
	ID             int    `json:"id"`
	SubjectID      int    `json:"subject_id"`
	Name           string `json:"name"`
	BackupWikipage string `json:"backup_wikipage"`
	LastBackup     int    `json:"last_backup"`
}

// InsertDBSubject добавляет новое поле в таблицу subject
func InsertDBSubject(subject Subject) error {
	// получаем ссылку на db
	db, err := openDB()
	defer db.Close()
	if err != nil {
		return err
	}

	// добавляем новое поле в таблицу
	query := fmt.Sprintf(`INSERT INTO subject (subject_id, name, 
		backup_wikipage, last_backup) VALUES ('%d', '%v', '%v', '%d')`,
		subject.SubjectID, subject.Name, subject.BackupWikipage, subject.LastBackup)
	_, err = db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

// SelectDBSubjects извлекает поля из таблицы subject
func SelectDBSubjects() ([]Subject, error) {
	// получаем ссылку на db
	db, err := openDB()
	defer db.Close()
	if err != nil {
		return nil, err
	}

	// читаем данные из БД
	rows, err := db.Query("SELECT * FROM subject")
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	// считываем данные из rows
	var subjects []Subject
	for rows.Next() {
		var subject Subject
		err = rows.Scan(&subject.ID, &subject.SubjectID, &subject.Name,
			&subject.BackupWikipage, &subject.LastBackup)
		if err != nil {
			return nil, err
		}
		subjects = append(subjects, subject)
	}
	return subjects, nil
}

// Monitor - структура для полей из таблицы monitor
type Monitor struct {
	ID        int
	Name      string
	SubjectID int
}

// InsertDBMonitor добавляет новое поле в таблицу monitor
func InsertDBMonitor(monitor Monitor) error {
	// получаем ссылку на db
	db, err := openDB()
	defer db.Close()
	if err != nil {
		return err
	}

	// добавляем новое поле в таблицу
	query := fmt.Sprintf(`INSERT INTO monitor (name, subject_id) 
		VALUES ('%v', '%d')`,
		monitor.Name, monitor.SubjectID)
	_, err = db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

// SelectDBMonitor извлекает из таблицы monitor поле с указанными name и subject_id
func SelectDBMonitor(monitorName string, subjectID int) (*Monitor, error) {
	// получаем ссылку на db
	db, err := openDB()
	defer db.Close()
	if err != nil {
		return nil, err
	}

	// читаем данные из БД
	query := fmt.Sprintf("SELECT * FROM monitor WHERE name='%v' and subject_id=%d", monitorName, subjectID)
	rows, err := db.Query(query)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	// считываем данные из rows
	var monitor Monitor
	for rows.Next() {
		err = rows.Scan(&monitor.ID, &monitor.Name, &monitor.SubjectID)
		if err != nil {
			return nil, err
		}
	}
	return &monitor, nil
}

// Method - структура для полей из таблицы method
type Method struct {
	ID            int
	Name          string
	SubjectID     int
	AccessTokenID int
	MonitorID     int
}

// InsertDBMethod добавляет новое поле в таблицу method
func InsertDBMethod(method Method) error {
	// получаем ссылку на db
	db, err := openDB()
	defer db.Close()
	if err != nil {
		return err
	}

	// добавляем новое поле в таблицу
	query := fmt.Sprintf(`INSERT INTO method (name, subject_id, 
		access_token_id, monitor_id) 
		VALUES ('%v', '%d', '%d', '%d')`,
		method.Name, method.SubjectID, method.AccessTokenID, method.MonitorID)
	_, err = db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

// SelectDBMethod извлекает из таблицы method поле с указанным name, subject_id и monitor_id
func SelectDBMethod(methodName string, subjectID, monitorID int) (*Method, error) {
	// получаем ссылку на db
	db, err := openDB()
	defer db.Close()
	if err != nil {
		return nil, err
	}

	// читаем данные из БД
	query := fmt.Sprintf("SELECT * FROM method WHERE name='%v' and subject_id=%d and monitor_id=%d",
		methodName, subjectID, monitorID)
	rows, err := db.Query(query)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	// считываем данные из rows
	var method Method
	for rows.Next() {
		err = rows.Scan(&method.ID, &method.Name,
			&method.SubjectID, &method.AccessTokenID,
			&method.MonitorID)
		if err != nil {
			return nil, err
		}
	}
	return &method, nil
}

// WallPostMonitorParam - структура для полей из таблицы wall_post_monitor
type WallPostMonitorParam struct {
	ID                    int
	SubjectID             int
	NeedMonitoring        int
	Interval              int
	SendTo                int
	Filter                string
	LastDate              int
	PostsCount            int
	KeywordsForMonitoring string
	UsersIDsForIgnore     string
}

// InsertDBWallPostMonitor добавляет новое поле в таблицу wall_post_monitor
func InsertDBWallPostMonitor(wPMP WallPostMonitorParam) error {
	// получаем ссылку на db
	db, err := openDB()
	defer db.Close()
	if err != nil {
		return err
	}

	// добавляем новое поле в таблицу
	query := fmt.Sprintf(`INSERT INTO wall_post_monitor (subject_id, need_monitoring, interval, 
		send_to, filter, last_date, posts_count, keywords_for_monitoring, users_ids_for_ignore) 
		VALUES ('%d', '%d', '%d', '%d', '%v', '%d', '%d', '%v', '%v')`,
		wPMP.SubjectID, wPMP.NeedMonitoring, wPMP.Interval, wPMP.SendTo,
		wPMP.Filter, wPMP.LastDate, wPMP.PostsCount, wPMP.KeywordsForMonitoring, wPMP.UsersIDsForIgnore)
	_, err = db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

// SelectDBWallPostMonitorParam извлекает поля из таблицы wall_post_monitor
func SelectDBWallPostMonitorParam(subjectID int) (WallPostMonitorParam, error) {
	var wallPostMonitorParam WallPostMonitorParam
	// получаем ссылку на db
	db, err := openDB()
	defer db.Close()
	if err != nil {
		return wallPostMonitorParam, err
	}

	// читаем данные из БД
	query := fmt.Sprintf("SELECT * FROM wall_post_monitor WHERE subject_id=%d", subjectID)
	rows, err := db.Query(query)
	defer rows.Close()
	if err != nil {
		return wallPostMonitorParam, err
	}

	// считываем данные из rows
	for rows.Next() {
		err = rows.Scan(&wallPostMonitorParam.ID, &wallPostMonitorParam.SubjectID,
			&wallPostMonitorParam.NeedMonitoring, &wallPostMonitorParam.Interval,
			&wallPostMonitorParam.SendTo, &wallPostMonitorParam.Filter,
			&wallPostMonitorParam.LastDate, &wallPostMonitorParam.PostsCount,
			&wallPostMonitorParam.KeywordsForMonitoring, &wallPostMonitorParam.UsersIDsForIgnore)
		if err != nil {
			return wallPostMonitorParam, err
		}
	}

	return wallPostMonitorParam, nil
}

// UpdateDBWallPostMonitorLastDate обновляет значение в поле таблицы wall_post_monitor
func UpdateDBWallPostMonitorLastDate(subjectID int, newLastDate int,
	wallPostMonitorParam WallPostMonitorParam) error {
	// получаем ссылку на db
	db, err := openDB()
	defer db.Close()
	if err != nil {
		return err
	}

	// обновляем значения в конкретном поле
	query := fmt.Sprintf(`UPDATE wall_post_monitor SET last_date=%d WHERE subject_id=%d`, newLastDate, subjectID)
	_, err = db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

// AlbumPhotoMonitorParam - структура для полей из таблицы album_photo_monitor
type AlbumPhotoMonitorParam struct {
	ID             int `json:"id"`
	SubjectID      int `json:"subject_id"`
	NeedMonitoring int `json:"need_monitoring"`
	SendTo         int `json:"send_to"`
	Interval       int `json:"interval"`
	LastDate       int `json:"last_date"`
	PhotosCount    int `json:"photos_count"`
}

// InsertDBAlbumPhotoMonitor добавляет новое поле в таблицу album_photo_monitor
func InsertDBAlbumPhotoMonitor(aPMP AlbumPhotoMonitorParam) error {
	// получаем ссылку на db
	db, err := openDB()
	defer db.Close()
	if err != nil {
		return err
	}

	// добавляем новое поле в таблицу
	query := fmt.Sprintf(`INSERT INTO album_photo_monitor (subject_id, 
		need_monitoring, send_to, interval, 
		last_date, photos_count) 
		VALUES ('%d', '%d', '%d', '%d', '%d', '%d')`,
		aPMP.SubjectID, aPMP.NeedMonitoring, aPMP.SendTo, aPMP.Interval,
		aPMP.LastDate, aPMP.PhotosCount)
	_, err = db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

// SelectDBAlbumPhotoMonitorParam извлекает поля из таблицы album_photo_monitor
func SelectDBAlbumPhotoMonitorParam(subjectID int) (AlbumPhotoMonitorParam, error) {
	var albumPhotoMonitorParam AlbumPhotoMonitorParam
	// получаем ссылку на db
	db, err := openDB()
	defer db.Close()
	if err != nil {
		return albumPhotoMonitorParam, err
	}

	// читаем данные из БД
	query := fmt.Sprintf("SELECT * FROM album_photo_monitor WHERE subject_id=%d", subjectID)
	rows, err := db.Query(query)
	defer rows.Close()
	if err != nil {
		return albumPhotoMonitorParam, err
	}

	// считываем данные из rows
	for rows.Next() {
		err = rows.Scan(&albumPhotoMonitorParam.ID, &albumPhotoMonitorParam.SubjectID,
			&albumPhotoMonitorParam.NeedMonitoring, &albumPhotoMonitorParam.SendTo,
			&albumPhotoMonitorParam.Interval, &albumPhotoMonitorParam.LastDate,
			&albumPhotoMonitorParam.PhotosCount)
		if err != nil {
			return albumPhotoMonitorParam, err
		}
	}

	return albumPhotoMonitorParam, nil
}

// UpdateDBAlbumPhotoMonitorLastDate обновляет значение в поле таблицы album_photo_monitor
func UpdateDBAlbumPhotoMonitorLastDate(subjectID int, newLastDate int,
	albumPhotoMonitorParam AlbumPhotoMonitorParam) error {
	// получаем ссылку на db
	db, err := openDB()
	defer db.Close()
	if err != nil {
		return err
	}

	// обновляем значения в конкретном поле
	query := fmt.Sprintf(`UPDATE album_photo_monitor SET last_date=%d WHERE subject_id=%d`, newLastDate, subjectID)
	_, err = db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

// VideoMonitorParam - структура для полей из таблицы video_monitor
type VideoMonitorParam struct {
	ID             int `json:"id"`
	SubjectID      int `json:"subject_id"`
	NeedMonitoring int `json:"need_monitoring"`
	SendTo         int `json:"send_to"`
	Interval       int `json:"interval"`
	LastDate       int `json:"last_date"`
	VideoCount     int `json:"video_count"`
}

// InsertDBVideoMonitor добавляет новое поле в таблицу video_monitor
func InsertDBVideoMonitor(vMP VideoMonitorParam) error {
	// получаем ссылку на db
	db, err := openDB()
	defer db.Close()
	if err != nil {
		return err
	}

	// добавляем новое поле в таблицу
	query := fmt.Sprintf(`INSERT INTO video_monitor (subject_id, 
		need_monitoring, send_to, interval, last_date, video_count) 
		VALUES ('%d', '%d', '%d', '%d', '%d', '%d')`,
		vMP.SubjectID, vMP.NeedMonitoring, vMP.SendTo, vMP.Interval,
		vMP.LastDate, vMP.VideoCount)
	_, err = db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

// SelectDBVideoMonitorParam извлекает поля из таблицы video_monitor
func SelectDBVideoMonitorParam(subjectID int) (VideoMonitorParam, error) {
	var videoMonitorParam VideoMonitorParam
	// получаем ссылку на db
	db, err := openDB()
	defer db.Close()
	if err != nil {
		return videoMonitorParam, err
	}

	// читаем данные из БД
	query := fmt.Sprintf("SELECT * FROM video_monitor WHERE subject_id=%d", subjectID)
	rows, err := db.Query(query)
	defer rows.Close()
	if err != nil {
		return videoMonitorParam, err
	}

	// считываем данные из rows
	for rows.Next() {
		err = rows.Scan(&videoMonitorParam.ID, &videoMonitorParam.SubjectID,
			&videoMonitorParam.NeedMonitoring, &videoMonitorParam.SendTo,
			&videoMonitorParam.VideoCount, &videoMonitorParam.LastDate,
			&videoMonitorParam.Interval)
		if err != nil {
			return videoMonitorParam, err
		}
	}

	return videoMonitorParam, nil
}

// UpdateDBVideoMonitorLastDate обновляет значение в поле таблицы video_monitor
func UpdateDBVideoMonitorLastDate(subjectID int, newLastDate int) error {
	// получаем ссылку на db
	db, err := openDB()
	defer db.Close()
	if err != nil {
		return err
	}

	// обновляем значения в конкретном поле
	query := fmt.Sprintf(`UPDATE video_monitor SET last_date=%d WHERE subject_id=%d`, newLastDate, subjectID)
	_, err = db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

// PhotoCommentMonitorParam - структура для полей из таблицы photo_comment_monitor
type PhotoCommentMonitorParam struct {
	ID             int `json:"id"`
	SubjectID      int `json:"subject_id"`
	NeedMonitoring int `json:"need_monitoring"`
	CommentsCount  int `json:"comments_count"`
	LastDate       int `json:"last_date"`
	Interval       int `json:"interval"`
	SendTo         int `json:"send_to"`
}

// InsertDBPhotoCommentMonitor добавляет новое поле в таблицу photo_comment_monitor
func InsertDBPhotoCommentMonitor(pCMP PhotoCommentMonitorParam) error {
	// получаем ссылку на db
	db, err := openDB()
	defer db.Close()
	if err != nil {
		return err
	}

	// добавляем новое поле в таблицу
	query := fmt.Sprintf(`INSERT INTO photo_comment_monitor (subject_id, 
		need_monitoring, comments_count, last_date, interval, send_to) 
		VALUES ('%d', '%d', '%d', '%d', '%d', '%d')`,
		pCMP.SubjectID, pCMP.NeedMonitoring, pCMP.CommentsCount, pCMP.LastDate,
		pCMP.Interval, pCMP.SendTo)
	_, err = db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

// SelectDBPhotoCommentMonitorParam извлекает поля из таблицы photo_comment_monitor
func SelectDBPhotoCommentMonitorParam(subjectID int) (PhotoCommentMonitorParam, error) {
	var photoCommentMonitorParam PhotoCommentMonitorParam
	// получаем ссылку на db
	db, err := openDB()
	defer db.Close()
	if err != nil {
		return photoCommentMonitorParam, err
	}

	// читаем данные из БД
	query := fmt.Sprintf("SELECT * FROM photo_comment_monitor WHERE subject_id=%d", subjectID)
	rows, err := db.Query(query)
	defer rows.Close()
	if err != nil {
		return photoCommentMonitorParam, err
	}

	// считываем данные из rows
	for rows.Next() {
		err = rows.Scan(&photoCommentMonitorParam.ID, &photoCommentMonitorParam.SubjectID,
			&photoCommentMonitorParam.NeedMonitoring, &photoCommentMonitorParam.CommentsCount,
			&photoCommentMonitorParam.LastDate, &photoCommentMonitorParam.Interval,
			&photoCommentMonitorParam.SendTo)
		if err != nil {
			return photoCommentMonitorParam, err
		}
	}

	return photoCommentMonitorParam, nil
}

// UpdateDBPhotoCommentMonitorLastDate обновляет значение в поле таблицы photo_comment_monitor
func UpdateDBPhotoCommentMonitorLastDate(subjectID int, newLastDate int) error {
	// получаем ссылку на db
	db, err := openDB()
	defer db.Close()
	if err != nil {
		return err
	}

	// обновляем значения в конкретном поле
	query := fmt.Sprintf(`UPDATE photo_comment_monitor SET last_date=%d WHERE subject_id=%d`, newLastDate, subjectID)
	_, err = db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

// VideoCommentMonitorParam - структура для полей из таблицы video_comment_monitor
type VideoCommentMonitorParam struct {
	ID             int `json:"id"`
	SubjectID      int `json:"subject_id"`
	NeedMonitoring int `json:"need_monitoring"`
	VideosCount    int `json:"videos_count"`
	Interval       int `json:"interval"`
	CommentsCount  int `json:"comments_count"`
	SendTo         int `json:"send_to"`
	LastDate       int `json:"last_date"`
}

// InsertDBVideoCommentMonitor добавляет новое поле в таблицу video_comment_monitor
func InsertDBVideoCommentMonitor(vCMP VideoCommentMonitorParam) error {
	// получаем ссылку на db
	db, err := openDB()
	defer db.Close()
	if err != nil {
		return err
	}

	// добавляем новое поле в таблицу
	query := fmt.Sprintf(`INSERT INTO video_comment_monitor (subject_id, need_monitoring, videos_count, 
		interval, comments_count, send_to, last_date) 
		VALUES ('%d', '%d', '%d', '%d', '%d', '%d', '%d')`,
		vCMP.SubjectID, vCMP.NeedMonitoring, vCMP.VideosCount, vCMP.Interval, vCMP.CommentsCount, vCMP.SendTo,
		vCMP.LastDate)
	_, err = db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

// SelectDBVideoCommentMonitorParam извлекает поля из таблицы video_comment_monitor
func SelectDBVideoCommentMonitorParam(subjectID int) (VideoCommentMonitorParam, error) {
	var videoCommentMonitorParam VideoCommentMonitorParam
	// получаем ссылку на db
	db, err := openDB()
	defer db.Close()
	if err != nil {
		return videoCommentMonitorParam, err
	}

	// читаем данные из БД
	query := fmt.Sprintf("SELECT * FROM video_comment_monitor WHERE subject_id=%d", subjectID)
	rows, err := db.Query(query)
	defer rows.Close()
	if err != nil {
		return videoCommentMonitorParam, err
	}

	// считываем данные из rows
	for rows.Next() {
		err = rows.Scan(&videoCommentMonitorParam.ID, &videoCommentMonitorParam.SubjectID,
			&videoCommentMonitorParam.NeedMonitoring, &videoCommentMonitorParam.VideosCount,
			&videoCommentMonitorParam.Interval, &videoCommentMonitorParam.CommentsCount,
			&videoCommentMonitorParam.SendTo, &videoCommentMonitorParam.LastDate)
		if err != nil {
			return videoCommentMonitorParam, err
		}
	}

	return videoCommentMonitorParam, nil
}

// UpdateDBVideoCommentMonitorLastDate обновляет значение в поле таблицы video_comment_monitor
func UpdateDBVideoCommentMonitorLastDate(subjectID int, newLastDate int) error {
	// получаем ссылку на db
	db, err := openDB()
	defer db.Close()
	if err != nil {
		return err
	}

	// обновляем значения в конкретном поле
	query := fmt.Sprintf(`UPDATE video_comment_monitor SET last_date=%d WHERE subject_id=%d`,
		newLastDate, subjectID)
	_, err = db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

// TopicMonitorParam - структура для полей из таблицы topic_monitor
type TopicMonitorParam struct {
	ID             int `json:"id"`
	SubjectID      int `json:"subject_id"`
	NeedMonitoring int `json:"need_monitoring"`
	TopicsCount    int `json:"topics_count"`
	CommentsCount  int `json:"comments_count"`
	Interval       int `json:"interval"`
	SendTo         int `json:"send_to"`
	LastDate       int `json:"last_date"`
}

// InsertDBTopicMonitor добавляет новое поле в таблицу topic_monitor
func InsertDBTopicMonitor(tMP TopicMonitorParam) error {
	// получаем ссылку на db
	db, err := openDB()
	defer db.Close()
	if err != nil {
		return err
	}

	// добавляем новое поле в таблицу
	query := fmt.Sprintf(`INSERT INTO topic_monitor (subject_id, 
		need_monitoring, topics_count, comments_count, 
		interval, send_to, last_date) 
		VALUES ('%d', '%d', '%d', '%d', '%d', '%d', '%d')`,
		tMP.SubjectID, tMP.NeedMonitoring, tMP.TopicsCount, tMP.CommentsCount,
		tMP.Interval, tMP.SendTo, tMP.LastDate)
	_, err = db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

// SelectDBTopicMonitorParam извлекает поля из таблицы topic_monitor
func SelectDBTopicMonitorParam(subjectID int) (TopicMonitorParam, error) {
	var topicMonitorParam TopicMonitorParam
	// получаем ссылку на db
	db, err := openDB()
	defer db.Close()
	if err != nil {
		return topicMonitorParam, err
	}

	// читаем данные из БД
	query := fmt.Sprintf("SELECT * FROM topic_monitor WHERE subject_id=%d", subjectID)
	rows, err := db.Query(query)
	defer rows.Close()
	if err != nil {
		return topicMonitorParam, err
	}

	// считываем данные из rows
	for rows.Next() {
		err = rows.Scan(&topicMonitorParam.ID, &topicMonitorParam.SubjectID,
			&topicMonitorParam.NeedMonitoring, &topicMonitorParam.TopicsCount,
			&topicMonitorParam.CommentsCount, &topicMonitorParam.Interval,
			&topicMonitorParam.SendTo, &topicMonitorParam.LastDate)
		if err != nil {
			return topicMonitorParam, err
		}
	}

	return topicMonitorParam, nil
}

// UpdateDBTopicMonitorLastDate обновляет значение в поле таблицы topic_monitor
func UpdateDBTopicMonitorLastDate(subjectID int, newLastDate int) error {
	// получаем ссылку на db
	db, err := openDB()
	defer db.Close()
	if err != nil {
		return err
	}

	// обновляем значения в конкретном поле
	query := fmt.Sprintf(`UPDATE topic_monitor SET last_date=%d WHERE subject_id=%d`,
		newLastDate, subjectID)
	_, err = db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

// WallPostCommentMonitorParam - структура для полей из таблицы wall_post_comment_monitor
type WallPostCommentMonitorParam struct {
	ID                                  int    `json:"id"`
	SubjectID                           int    `json:"subject_id"`
	NeedMonitoring                      int    `json:"need_monitoring"`
	PostsCount                          int    `json:"posts_count"`
	CommentsCount                       int    `json:"comments_count"`
	MonitoringAll                       int    `json:"monitoring_all"`
	UsersIDsForMonitoring               string `json:"users_ids_for_monitoring"`
	UsersNamesForMonitoring             string `json:"users_names_for_monitoring"`
	AttachmentsTypesForMonitoring       string `json:"attachments_types_for_monitoring"`
	UsersIDsForIgnore                   string `json:"users_ids_for_ignore"`
	Interval                            int    `json:"interval"`
	SendTo                              int    `json:"send_to"`
	Filter                              string `json:"filter"`
	LastDate                            int    `json:"last_date"`
	KeywordsForMonitoring               string `json:"keywords_for_monitoring"`
	SmallCommentsForMonitoring          string `json:"small_comments_for_monitoring"`
	DigitsCountForCardNumberMonitoring  string `json:"digits_count_for_card_number_monitoring"`
	DigitsCountForPhoneNumberMonitoring string `json:"digits_count_for_phone_number_monitoring"`
	MonitorByCommunity                  int    `json:"monitor_by_community"`
}

// InsertDBWallPostCommentMonitor добавляет новое поле в таблицу wall_post_comment_monitor
func InsertDBWallPostCommentMonitor(wPCMP WallPostCommentMonitorParam) error {
	// получаем ссылку на db
	db, err := openDB()
	defer db.Close()
	if err != nil {
		return err
	}

	// добавляем новое поле в таблицу
	query := fmt.Sprintf(`INSERT INTO wall_post_comment_monitor (subject_id, 
		need_monitoring, posts_count, comments_count, monitoring_all, 
   		users_ids_for_monitoring, users_names_for_monitoring, attachments_types_for_monitoring, 
		users_ids_for_ignore, interval, send_to, filter, last_date, 
		keywords_for_monitoring, small_comments_for_monitoring, 
		digits_count_for_card_number_monitoring, 
		digits_count_for_phone_number_monitoring, monitor_by_community)
		VALUES ('%d', '%d', '%d', '%d', '%d', '%v', '%v', '%v', '%v', '%d', 
		'%d', '%v', '%d', '%v', '%v', '%v', '%v', '%d')`,
		wPCMP.SubjectID, wPCMP.NeedMonitoring, wPCMP.PostsCount, wPCMP.CommentsCount,
		wPCMP.MonitoringAll, wPCMP.UsersIDsForMonitoring, wPCMP.UsersNamesForMonitoring,
		wPCMP.AttachmentsTypesForMonitoring, wPCMP.UsersIDsForIgnore, wPCMP.Interval,
		wPCMP.SendTo, wPCMP.Filter, wPCMP.LastDate, wPCMP.KeywordsForMonitoring,
		wPCMP.SmallCommentsForMonitoring, wPCMP.DigitsCountForCardNumberMonitoring,
		wPCMP.DigitsCountForPhoneNumberMonitoring, wPCMP.MonitorByCommunity)
	_, err = db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

// SelectDBWallPostCommentMonitorParam извлекает поля из таблицы wall_post_comment_monitor
func SelectDBWallPostCommentMonitorParam(subjectID int) (WallPostCommentMonitorParam, error) {
	var wallPostCommentMonitorParam WallPostCommentMonitorParam
	// получаем ссылку на db
	db, err := openDB()
	defer db.Close()
	if err != nil {
		return wallPostCommentMonitorParam, err
	}

	// читаем данные из БД
	query := fmt.Sprintf("SELECT * FROM wall_post_comment_monitor WHERE subject_id=%d", subjectID)
	rows, err := db.Query(query)
	defer rows.Close()
	if err != nil {
		return wallPostCommentMonitorParam, err
	}

	// считываем данные из rows
	for rows.Next() {
		err = rows.Scan(&wallPostCommentMonitorParam.ID, &wallPostCommentMonitorParam.SubjectID,
			&wallPostCommentMonitorParam.NeedMonitoring, &wallPostCommentMonitorParam.PostsCount,
			&wallPostCommentMonitorParam.CommentsCount, &wallPostCommentMonitorParam.MonitoringAll,
			&wallPostCommentMonitorParam.UsersIDsForMonitoring,
			&wallPostCommentMonitorParam.UsersNamesForMonitoring,
			&wallPostCommentMonitorParam.AttachmentsTypesForMonitoring,
			&wallPostCommentMonitorParam.UsersIDsForIgnore,
			&wallPostCommentMonitorParam.Interval, &wallPostCommentMonitorParam.SendTo,
			&wallPostCommentMonitorParam.Filter, &wallPostCommentMonitorParam.LastDate,
			&wallPostCommentMonitorParam.KeywordsForMonitoring,
			&wallPostCommentMonitorParam.SmallCommentsForMonitoring,
			&wallPostCommentMonitorParam.DigitsCountForCardNumberMonitoring,
			&wallPostCommentMonitorParam.DigitsCountForPhoneNumberMonitoring,
			&wallPostCommentMonitorParam.MonitorByCommunity)
		if err != nil {
			return wallPostCommentMonitorParam, err
		}
	}

	return wallPostCommentMonitorParam, nil
}

// UpdateDBWallPostCommentMonitorLastDate обновляет значение в поле таблицы wall_post_comment_monitor
func UpdateDBWallPostCommentMonitorLastDate(subjectID int, newLastDate int) error {
	// получаем ссылку на db
	db, err := openDB()
	defer db.Close()
	if err != nil {
		return err
	}

	// обновляем значения в конкретном поле
	query := fmt.Sprintf(`UPDATE wall_post_comment_monitor SET last_date=%d WHERE subject_id=%d`,
		newLastDate, subjectID)
	_, err = db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

// InitDB создает таблицы и связи в базе данных (желательно, пустой)
func InitDB() error {
	// получаем ссылку на db
	db, err := openDB()
	defer db.Close()
	if err != nil {
		return err
	}

	// формируем запрос на создание таблиц и связей в БД
	query := fmt.Sprintf(`BEGIN TRANSACTION;
		CREATE TABLE IF NOT EXISTS "access_token" (
			"id"	INTEGER NOT NULL UNIQUE,
			"name"	TEXT,
			"value"	TEXT,
			PRIMARY KEY("id" AUTOINCREMENT)
		);
		CREATE TABLE IF NOT EXISTS "wall_post_monitor" (
			"id"	INTEGER NOT NULL UNIQUE,
			"subject_id"	INTEGER NOT NULL,
			"need_monitoring"	INTEGER NOT NULL DEFAULT 0,
			"interval"	INTEGER,
			"send_to"	INTEGER,
			"filter"	TEXT,
			"last_date"	INTEGER,
			"posts_count"	INTEGER,
			"keywords_for_monitoring"	TEXT,
			"users_ids_for_ignore"	TEXT,
			FOREIGN KEY("subject_id") REFERENCES "subject"("id"),
			PRIMARY KEY("id" AUTOINCREMENT)
		);
		CREATE TABLE IF NOT EXISTS "video_comment_monitor" (
			"id"	INTEGER NOT NULL UNIQUE,
			"subject_id"	INTEGER NOT NULL,
			"need_monitoring"	INTEGER NOT NULL DEFAULT 0,
			"videos_count"	INTEGER,
			"interval"	INTEGER,
			"comments_count"	INTEGER,
			"send_to"	INTEGER,
			"last_date"	INTEGER,
			FOREIGN KEY("subject_id") REFERENCES "subject"("id"),
			PRIMARY KEY("id" AUTOINCREMENT)
		);
		CREATE TABLE IF NOT EXISTS "photo_comment_monitor" (
			"id"	INTEGER NOT NULL UNIQUE,
			"subject_id"	INTEGER NOT NULL,
			"need_monitoring"	INTEGER NOT NULL DEFAULT 0,
			"comments_count"	INTEGER,
			"last_date"	INTEGER,
			"interval"	INTEGER,
			"send_to"	INTEGER,
			FOREIGN KEY("subject_id") REFERENCES "subject"("id"),
			PRIMARY KEY("id" AUTOINCREMENT)
		);
		CREATE TABLE IF NOT EXISTS "album_photo_monitor" (
			"id"	INTEGER NOT NULL UNIQUE,
			"subject_id"	INTEGER NOT NULL,
			"need_monitoring"	INTEGER NOT NULL DEFAULT 0,
			"send_to"	INTEGER,
			"interval"	INTEGER,
			"last_date"	INTEGER,
			"photos_count"	INTEGER,
			FOREIGN KEY("subject_id") REFERENCES "subject"("id"),
			PRIMARY KEY("id" AUTOINCREMENT)
		);
		CREATE TABLE IF NOT EXISTS "subject" (
			"id"	INTEGER NOT NULL UNIQUE,
			"subject_id"	INTEGER,
			"name"	TEXT,
			"backup_wikipage"	TEXT,
			"last_backup"	INTEGER,
			PRIMARY KEY("id" AUTOINCREMENT)
		);
		CREATE TABLE IF NOT EXISTS "monitor" (
			"id"	INTEGER NOT NULL UNIQUE,
			"name"	TEXT,
			"subject_id"	INTEGER NOT NULL,
			FOREIGN KEY("subject_id") REFERENCES "subject"("id"),
			PRIMARY KEY("id" AUTOINCREMENT)
		);
		CREATE TABLE IF NOT EXISTS "method" (
			"id"	INTEGER NOT NULL UNIQUE,
			"name"	TEXT,
			"subject_id"	INTEGER NOT NULL,
			"access_token_id"	INTEGER NOT NULL,
			"monitor_id"	INTEGER,
			FOREIGN KEY("monitor_id") REFERENCES "monitor"("id"),
			FOREIGN KEY("access_token_id") REFERENCES "access_token"("id"),
			FOREIGN KEY("subject_id") REFERENCES "subject"("id"),
			PRIMARY KEY("id" AUTOINCREMENT)
		);
		CREATE TABLE IF NOT EXISTS "video_monitor" (
			"id"	INTEGER NOT NULL UNIQUE,
			"subject_id"	INTEGER NOT NULL,
			"need_monitoring"	INTEGER NOT NULL DEFAULT 0,
			"send_to"	INTEGER,
			"video_count"	INTEGER,
			"last_date"	INTEGER,
			"interval"	INTEGER,
			FOREIGN KEY("subject_id") REFERENCES "subject"("id"),
			PRIMARY KEY("id" AUTOINCREMENT)
		);
		CREATE TABLE IF NOT EXISTS "topic_monitor" (
			"id"	INTEGER NOT NULL UNIQUE,
			"subject_id"	INTEGER NOT NULL,
			"need_monitoring"	INTEGER NOT NULL DEFAULT 0,
			"topics_count"	INTEGER,
			"comments_count"	INTEGER,
			"interval"	INTEGER,
			"send_to"	INTEGER,
			"last_date"	INTEGER,
			FOREIGN KEY("subject_id") REFERENCES "subject"("id"),
			PRIMARY KEY("id" AUTOINCREMENT)
		);
		CREATE TABLE IF NOT EXISTS "wall_post_comment_monitor" (
			"id"	INTEGER NOT NULL UNIQUE,
			"subject_id"	INTEGER NOT NULL,
			"need_monitoring"	INTEGER NOT NULL DEFAULT 0,
			"posts_count"	INTEGER,
			"comments_count"	INTEGER,
			"monitoring_all"	INTEGER,
			"users_ids_for_monitoring"	TEXT,
			"users_names_for_monitoring"	TEXT,
			"attachments_types_for_monitoring"	TEXT,
			"users_ids_for_ignore"	TEXT,
			"interval"	INTEGER,
			"send_to"	INTEGER,
			"filter"	TEXT,
			"last_date"	INTEGER,
			"keywords_for_monitoring"	TEXT,
			"small_comments_for_monitoring"	TEXT,
			"digits_count_for_card_number_monitoring"	INTEGER,
			"digits_count_for_phone_number_monitoring"	INTEGER,
			"monitor_by_community"	INTEGER,
			FOREIGN KEY("subject_id") REFERENCES "subject"("id"),
			PRIMARY KEY("id" AUTOINCREMENT)
		);
		COMMIT;`)
	_, err = db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}
