package main

import (
	"database/sql"
	"fmt"

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

// UpdateDBAccessToken обновляет значение в поле таблицы access_token
func UpdateDBAccessToken(accessToken AccessToken) error {
	// получаем ссылку на db
	db, err := openDB()
	defer db.Close()
	if err != nil {
		return err
	}

	// обновляем значения в конкретном поле
	query := fmt.Sprintf(`UPDATE access_token SET value='%v' WHERE id=%d`, accessToken.Value, accessToken.ID)
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
