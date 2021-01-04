package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3" // иначе драйвер для работы с SQLite не работает
)

// DataBaseKit хранит ссылку на объект базы данных
type DataBaseKit struct {
	db *sql.DB
}

func (dbKit *DataBaseKit) initDB() error {
	err := dbKit.openDB()
	defer dbKit.db.Close()
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
	_, err = dbKit.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func (dbKit *DataBaseKit) openDB() error {
	path, err := ReadPathFile()
	if err != nil {
		return err
	}
	pathToDB := path + "groupsmonitor_db.db"
	dbKit.db, err = sql.Open("sqlite3", pathToDB)
	if err != nil {
		errorMessage := fmt.Sprintf("%v", err)

		typeError, causeError := DBIOError(errorMessage)

		switch typeError {

		case "timeout error":
			interval := 1 * time.Second
			sender := "Database"
			message := fmt.Sprintf("Error: %v. Timeout for %v...", causeError, interval)
			OutputMessage(sender, message)
			time.Sleep(interval)
			return dbKit.openDB()
		}

		return err
	}

	return nil
}

func (dbKit *DataBaseKit) selectTableAccessToken() ([]AccessToken, error) {
	err := dbKit.openDB()
	defer dbKit.db.Close()
	if err != nil {
		return nil, err
	}

	rows, err := dbKit.db.Query("SELECT * FROM access_token")
	defer rows.Close()
	if err != nil {
		return nil, err
	}

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

func (dbKit *DataBaseKit) selectTableSubject() ([]Subject, error) {
	err := dbKit.openDB()
	defer dbKit.db.Close()
	if err != nil {
		return nil, err
	}

	// читаем данные из БД
	rows, err := dbKit.db.Query("SELECT * FROM subject")
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

// AccessToken - структура для полей из таблицы access_token
type AccessToken struct {
	ID    int
	Name  string
	Value string
}

func (at *AccessToken) insertIntoDB() error {
	var dbKit DataBaseKit
	err := dbKit.openDB()
	defer dbKit.db.Close()
	if err != nil {
		return err
	}

	query := fmt.Sprintf(`INSERT INTO access_token (name, value) VALUES ('%v', '%v')`,
		at.Name, at.Value)
	_, err = dbKit.db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func (at *AccessToken) selectFromDBByID(accessTokenID int) error {
	var dbKit DataBaseKit
	err := dbKit.openDB()
	defer dbKit.db.Close()
	if err != nil {
		return err
	}

	query := fmt.Sprintf("SELECT * FROM access_token WHERE id=%d", accessTokenID)
	rows, err := dbKit.db.Query(query)
	defer rows.Close()
	if err != nil {
		return err
	}

	for rows.Next() {
		err = rows.Scan(&at.ID, &at.Name, &at.Value)
		if err != nil {
			return err
		}
	}
	return nil
}

func (at *AccessToken) selectFromDBByName(accessTokenName string) error {
	var dbKit DataBaseKit
	err := dbKit.openDB()
	defer dbKit.db.Close()
	if err != nil {
		return err
	}

	query := fmt.Sprintf("SELECT * FROM access_token WHERE name='%v'", accessTokenName)
	rows, err := dbKit.db.Query(query)
	defer rows.Close()
	if err != nil {
		return err
	}

	for rows.Next() {
		err = rows.Scan(&at.ID, &at.Name, &at.Value)
		if err != nil {
			return err
		}
	}
	return nil
}

func (at *AccessToken) updateInDB() error {
	var dbKit DataBaseKit
	err := dbKit.openDB()
	defer dbKit.db.Close()
	if err != nil {
		return err
	}

	query := fmt.Sprintf(`UPDATE access_token SET name='%v', value='%v' WHERE id=%d`,
		at.Name, at.Value, at.ID)
	_, err = dbKit.db.Exec(query)
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

func (s *Subject) insertIntoDB() error {
	var dbKit DataBaseKit
	err := dbKit.openDB()
	defer dbKit.db.Close()
	if err != nil {
		return err
	}

	query := fmt.Sprintf(`INSERT INTO subject (subject_id, name, 
		backup_wikipage, last_backup) VALUES ('%d', '%v', '%v', '%d')`,
		s.SubjectID, s.Name, s.BackupWikipage, s.LastBackup)
	_, err = dbKit.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func (s *Subject) selectFromDBByID(id int) error {
	var dbKit DataBaseKit
	err := dbKit.openDB()
	defer dbKit.db.Close()
	if err != nil {
		return err
	}

	query := fmt.Sprintf("SELECT * FROM subject WHERE id=%d", id)
	rows, err := dbKit.db.Query(query)
	defer rows.Close()
	if err != nil {
		return err
	}

	for rows.Next() {
		err = rows.Scan(&s.ID, &s.SubjectID, &s.Name, &s.BackupWikipage, &s.LastBackup)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Subject) updateInDB() error {
	var dbKit DataBaseKit
	err := dbKit.openDB()
	defer dbKit.db.Close()
	if err != nil {
		return err
	}

	query := fmt.Sprintf(`UPDATE subject 
		SET subject_id='%d', name='%v', backup_wikipage='%v', last_backup='%d' WHERE id=%d`,
		s.SubjectID, s.Name, s.BackupWikipage, s.LastBackup, s.ID)
	_, err = dbKit.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

// Monitor - структура для полей из таблицы monitor
type Monitor struct {
	ID        int
	Name      string
	SubjectID int
}

func (m *Monitor) insertIntoDB() error {
	var dbKit DataBaseKit
	err := dbKit.openDB()
	defer dbKit.db.Close()
	if err != nil {
		return err
	}
	query := fmt.Sprintf(`INSERT INTO monitor (name, subject_id) 
		VALUES ('%v', '%d')`,
		m.Name, m.SubjectID)
	_, err = dbKit.db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func (m *Monitor) selectFromDBByNameAndBySubjectID(monitorName string, subjectID int) error {
	var dbKit DataBaseKit
	err := dbKit.openDB()
	defer dbKit.db.Close()
	if err != nil {
		return err
	}

	query := fmt.Sprintf("SELECT * FROM monitor WHERE name='%v' and subject_id=%d", monitorName, subjectID)
	rows, err := dbKit.db.Query(query)
	defer rows.Close()
	if err != nil {
		return err
	}

	for rows.Next() {
		err = rows.Scan(&m.ID, &m.Name, &m.SubjectID)
		if err != nil {
			return err
		}
	}
	return nil
}

// Method - структура для полей из таблицы method
type Method struct {
	ID            int
	Name          string
	SubjectID     int
	AccessTokenID int
	MonitorID     int
}

func (m *Method) insertIntoDB() error {
	var dbKit DataBaseKit
	err := dbKit.openDB()
	defer dbKit.db.Close()
	if err != nil {
		return err
	}

	query := fmt.Sprintf(`INSERT INTO method (name, subject_id, 
		access_token_id, monitor_id) 
		VALUES ('%v', '%d', '%d', '%d')`,
		m.Name, m.SubjectID, m.AccessTokenID, m.MonitorID)
	_, err = dbKit.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func (m *Method) updateInDB() error {
	var dbKit DataBaseKit
	err := dbKit.openDB()
	defer dbKit.db.Close()
	if err != nil {
		return err
	}

	query := fmt.Sprintf(`UPDATE method 
		SET name='%v', subject_id='%d', access_token_id='%d', monitor_id='%d' WHERE id=%d`,
		m.Name, m.SubjectID, m.AccessTokenID, m.MonitorID, m.ID)
	_, err = dbKit.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func (m *Method) selectFromDBByNameAndBySubjectIDAndByMonitorID(methodName string, subjectID, monitorID int) error {
	var dbKit DataBaseKit
	err := dbKit.openDB()
	if err != nil {
		return err
	}

	query := fmt.Sprintf("SELECT * FROM method WHERE name='%v' and subject_id=%d and monitor_id=%d",
		methodName, subjectID, monitorID)
	rows, err := dbKit.db.Query(query)
	defer rows.Close()
	if err != nil {
		return err
	}

	for rows.Next() {
		err = rows.Scan(&m.ID, &m.Name,
			&m.SubjectID, &m.AccessTokenID,
			&m.MonitorID)
		if err != nil {
			return err
		}
	}
	return nil
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

func (wpmp *WallPostMonitorParam) insertIntoDB() error {
	var dbKit DataBaseKit
	err := dbKit.openDB()
	defer dbKit.db.Close()
	if err != nil {
		return err
	}

	query := fmt.Sprintf(`INSERT INTO wall_post_monitor (subject_id, need_monitoring, interval, 
		send_to, filter, last_date, posts_count, keywords_for_monitoring, users_ids_for_ignore) 
		VALUES ('%d', '%d', '%d', '%d', '%v', '%d', '%d', '%v', '%v')`,
		wpmp.SubjectID, wpmp.NeedMonitoring, wpmp.Interval, wpmp.SendTo,
		wpmp.Filter, wpmp.LastDate, wpmp.PostsCount, wpmp.KeywordsForMonitoring, wpmp.UsersIDsForIgnore)
	_, err = dbKit.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func (wpmp *WallPostMonitorParam) selectFromDBBySubjectID(subjectID int) error {
	var dbKit DataBaseKit
	err := dbKit.openDB()
	defer dbKit.db.Close()
	if err != nil {
		return err
	}

	query := fmt.Sprintf("SELECT * FROM wall_post_monitor WHERE subject_id=%d", subjectID)
	rows, err := dbKit.db.Query(query)
	defer rows.Close()
	if err != nil {
		return err
	}

	for rows.Next() {
		err = rows.Scan(&wpmp.ID, &wpmp.SubjectID,
			&wpmp.NeedMonitoring, &wpmp.Interval,
			&wpmp.SendTo, &wpmp.Filter,
			&wpmp.LastDate, &wpmp.PostsCount,
			&wpmp.KeywordsForMonitoring, &wpmp.UsersIDsForIgnore)
		if err != nil {
			return err
		}
	}

	return nil
}

func (wpmp *WallPostMonitorParam) updateInDB() error {
	var dbKit DataBaseKit
	err := dbKit.openDB()
	defer dbKit.db.Close()
	if err != nil {
		return err
	}

	query := fmt.Sprintf(`UPDATE wall_post_monitor 
		SET subject_id='%d', need_monitoring='%d', interval='%d', send_to='%d', filter='%v', 
		last_date='%d', posts_count='%d', keywords_for_monitoring='%v',
		users_ids_for_ignore='%v' WHERE id=%d`,
		wpmp.SubjectID, wpmp.NeedMonitoring,
		wpmp.Interval, wpmp.SendTo, wpmp.Filter,
		wpmp.LastDate, wpmp.PostsCount, wpmp.KeywordsForMonitoring,
		wpmp.UsersIDsForIgnore, wpmp.ID)
	_, err = dbKit.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func (wpmp *WallPostMonitorParam) updateInDBFieldLastDate(subjectID, newLastDate int) error {
	var dbKit DataBaseKit
	err := dbKit.openDB()
	defer dbKit.db.Close()
	if err != nil {
		return err
	}

	query := fmt.Sprintf(`UPDATE wall_post_monitor SET last_date=%d WHERE subject_id=%d`,
		newLastDate, subjectID)
	_, err = dbKit.db.Exec(query)
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

func (apmp *AlbumPhotoMonitorParam) insertIntoDB() error {
	var dbKit DataBaseKit
	err := dbKit.openDB()
	defer dbKit.db.Close()
	if err != nil {
		return err
	}

	query := fmt.Sprintf(`INSERT INTO album_photo_monitor (subject_id, 
		need_monitoring, send_to, interval, 
		last_date, photos_count) 
		VALUES ('%d', '%d', '%d', '%d', '%d', '%d')`,
		apmp.SubjectID, apmp.NeedMonitoring, apmp.SendTo, apmp.Interval,
		apmp.LastDate, apmp.PhotosCount)
	_, err = dbKit.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func (apmp *AlbumPhotoMonitorParam) selectFromDBBySubjectID(subjectID int) error {
	var dbKit DataBaseKit
	err := dbKit.openDB()
	defer dbKit.db.Close()
	if err != nil {
		return err
	}

	query := fmt.Sprintf("SELECT * FROM album_photo_monitor WHERE subject_id=%d", subjectID)
	rows, err := dbKit.db.Query(query)
	defer rows.Close()
	if err != nil {
		return err
	}

	for rows.Next() {
		err = rows.Scan(&apmp.ID, &apmp.SubjectID,
			&apmp.NeedMonitoring, &apmp.SendTo,
			&apmp.Interval, &apmp.LastDate,
			&apmp.PhotosCount)
		if err != nil {
			return err
		}
	}

	return nil
}

func (apmp *AlbumPhotoMonitorParam) updateInDB() error {
	var dbKit DataBaseKit
	err := dbKit.openDB()
	defer dbKit.db.Close()
	if err != nil {
		return err
	}

	// обновляем значения в конкретном поле
	query := fmt.Sprintf(`UPDATE album_photo_monitor 
		SET subject_id='%d', need_monitoring='%d', send_to='%d', 
		interval='%d', last_date='%d', photos_count='%d'
		WHERE id=%d`,
		apmp.SubjectID, apmp.NeedMonitoring, apmp.SendTo,
		apmp.Interval, apmp.LastDate, apmp.PhotosCount,
		apmp.ID)
	_, err = dbKit.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func (apmp *AlbumPhotoMonitorParam) updateInDBFieldLastDate(subjectID, newLastDate int) error {
	var dbKit DataBaseKit
	err := dbKit.openDB()
	defer dbKit.db.Close()
	if err != nil {
		return err
	}

	query := fmt.Sprintf(`UPDATE album_photo_monitor SET last_date=%d WHERE subject_id=%d`, newLastDate, subjectID)
	_, err = dbKit.db.Exec(query)
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

func (vmp *VideoMonitorParam) insertIntoDB() error {
	var dbKit DataBaseKit
	err := dbKit.openDB()
	defer dbKit.db.Close()
	if err != nil {
		return err
	}

	query := fmt.Sprintf(`INSERT INTO video_monitor (subject_id, 
		need_monitoring, send_to, interval, last_date, video_count) 
		VALUES ('%d', '%d', '%d', '%d', '%d', '%d')`,
		vmp.SubjectID, vmp.NeedMonitoring, vmp.SendTo, vmp.Interval,
		vmp.LastDate, vmp.VideoCount)
	_, err = dbKit.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func (vmp *VideoMonitorParam) selectFromDBBySubjectID(subjectID int) error {
	var dbKit DataBaseKit
	err := dbKit.openDB()
	defer dbKit.db.Close()
	if err != nil {
		return err
	}

	query := fmt.Sprintf("SELECT * FROM video_monitor WHERE subject_id=%d", subjectID)
	rows, err := dbKit.db.Query(query)
	defer rows.Close()
	if err != nil {
		return err
	}

	for rows.Next() {
		err = rows.Scan(&vmp.ID, &vmp.SubjectID,
			&vmp.NeedMonitoring, &vmp.SendTo,
			&vmp.VideoCount, &vmp.LastDate,
			&vmp.Interval)
		if err != nil {
			return err
		}
	}

	return nil
}

func (vmp *VideoMonitorParam) updateInDB() error {
	var dbKit DataBaseKit
	err := dbKit.openDB()
	defer dbKit.db.Close()
	if err != nil {
		return err
	}

	query := fmt.Sprintf(`UPDATE video_monitor 
		SET subject_id='%d', need_monitoring='%d', send_to='%d', 
		video_count='%d', last_date='%d', interval='%d'
		WHERE id=%d`,
		vmp.SubjectID, vmp.NeedMonitoring, vmp.SendTo,
		vmp.VideoCount, vmp.LastDate, vmp.Interval,
		vmp.ID)
	_, err = dbKit.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func (vmp *VideoMonitorParam) updateInDBFieldLastDate(subjectID, newLastDate int) error {
	var dbKit DataBaseKit
	err := dbKit.openDB()
	defer dbKit.db.Close()
	if err != nil {
		return err
	}

	query := fmt.Sprintf(`UPDATE video_monitor SET last_date=%d WHERE subject_id=%d`, newLastDate, subjectID)
	_, err = dbKit.db.Exec(query)
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

func (pcmp *PhotoCommentMonitorParam) insertIntoDB() error {
	var dbKit DataBaseKit
	err := dbKit.openDB()
	defer dbKit.db.Close()
	if err != nil {
		return err
	}

	query := fmt.Sprintf(`INSERT INTO photo_comment_monitor (subject_id, 
		need_monitoring, comments_count, last_date, interval, send_to) 
		VALUES ('%d', '%d', '%d', '%d', '%d', '%d')`,
		pcmp.SubjectID, pcmp.NeedMonitoring, pcmp.CommentsCount, pcmp.LastDate,
		pcmp.Interval, pcmp.SendTo)
	_, err = dbKit.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func (pcmp *PhotoCommentMonitorParam) selectFromDBBySubjectID(subjectID int) error {
	var dbKit DataBaseKit
	err := dbKit.openDB()
	defer dbKit.db.Close()
	if err != nil {
		return err
	}

	query := fmt.Sprintf("SELECT * FROM photo_comment_monitor WHERE subject_id=%d", subjectID)
	rows, err := dbKit.db.Query(query)
	defer rows.Close()
	if err != nil {
		return err
	}

	for rows.Next() {
		err = rows.Scan(&pcmp.ID, &pcmp.SubjectID,
			&pcmp.NeedMonitoring, &pcmp.CommentsCount,
			&pcmp.LastDate, &pcmp.Interval,
			&pcmp.SendTo)
		if err != nil {
			return err
		}
	}

	return nil
}

func (pcmp *PhotoCommentMonitorParam) updateInDB() error {
	var dbKit DataBaseKit
	err := dbKit.openDB()
	defer dbKit.db.Close()
	if err != nil {
		return err
	}

	query := fmt.Sprintf(`UPDATE photo_comment_monitor 
		SET subject_id='%d', need_monitoring='%d', send_to='%d', 
		comments_count='%d', last_date='%d', interval='%d'
		WHERE id=%d`,
		pcmp.SubjectID, pcmp.NeedMonitoring, pcmp.SendTo,
		pcmp.CommentsCount, pcmp.LastDate, pcmp.Interval,
		pcmp.ID)
	_, err = dbKit.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func (pcmp *PhotoCommentMonitorParam) updateInDBFieldLastDate(subjectID, newLastDate int) error {
	var dbKit DataBaseKit
	err := dbKit.openDB()
	defer dbKit.db.Close()
	if err != nil {
		return err
	}

	// обновляем значения в конкретном поле
	query := fmt.Sprintf(`UPDATE photo_comment_monitor SET last_date=%d WHERE subject_id=%d`, newLastDate, subjectID)
	_, err = dbKit.db.Exec(query)
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

func (vcmp *VideoCommentMonitorParam) insertIntoDB() error {
	var dbKit DataBaseKit
	err := dbKit.openDB()
	defer dbKit.db.Close()
	if err != nil {
		return err
	}

	query := fmt.Sprintf(`INSERT INTO video_comment_monitor (subject_id, need_monitoring, videos_count, 
		interval, comments_count, send_to, last_date) 
		VALUES ('%d', '%d', '%d', '%d', '%d', '%d', '%d')`,
		vcmp.SubjectID, vcmp.NeedMonitoring, vcmp.VideosCount, vcmp.Interval, vcmp.CommentsCount, vcmp.SendTo,
		vcmp.LastDate)
	_, err = dbKit.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func (vcmp *VideoCommentMonitorParam) selectFromDBBySubjectID(subjectID int) error {
	var dbKit DataBaseKit
	err := dbKit.openDB()
	defer dbKit.db.Close()
	if err != nil {
		return err
	}

	query := fmt.Sprintf("SELECT * FROM video_comment_monitor WHERE subject_id=%d", subjectID)
	rows, err := dbKit.db.Query(query)
	defer rows.Close()
	if err != nil {
		return err
	}

	for rows.Next() {
		err = rows.Scan(&vcmp.ID, &vcmp.SubjectID,
			&vcmp.NeedMonitoring, &vcmp.VideosCount,
			&vcmp.Interval, &vcmp.CommentsCount,
			&vcmp.SendTo, &vcmp.LastDate)
		if err != nil {
			return err
		}
	}

	return nil
}

func (vcmp *VideoCommentMonitorParam) updateInDB() error {
	var dbKit DataBaseKit
	err := dbKit.openDB()
	defer dbKit.db.Close()
	if err != nil {
		return err
	}

	query := fmt.Sprintf(`UPDATE video_comment_monitor 
		SET subject_id='%d', need_monitoring='%d', send_to='%d', 
		comments_count='%d', last_date='%d', interval='%d',
		videos_count='%d' WHERE id=%d`,
		vcmp.SubjectID, vcmp.NeedMonitoring, vcmp.SendTo,
		vcmp.CommentsCount, vcmp.LastDate, vcmp.Interval,
		vcmp.VideosCount, vcmp.ID)
	_, err = dbKit.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func (vcmp *VideoCommentMonitorParam) updateInDBFieldLastDate(subjectID, newLastDate int) error {
	var dbKit DataBaseKit
	err := dbKit.openDB()
	defer dbKit.db.Close()
	if err != nil {
		return err
	}

	query := fmt.Sprintf(`UPDATE video_comment_monitor SET last_date=%d WHERE subject_id=%d`,
		newLastDate, subjectID)
	_, err = dbKit.db.Exec(query)
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

func (tmp *TopicMonitorParam) insertIntoDB() error {
	var dbKit DataBaseKit
	err := dbKit.openDB()
	defer dbKit.db.Close()
	if err != nil {
		return err
	}

	query := fmt.Sprintf(`INSERT INTO topic_monitor (subject_id, 
		need_monitoring, topics_count, comments_count, 
		interval, send_to, last_date) 
		VALUES ('%d', '%d', '%d', '%d', '%d', '%d', '%d')`,
		tmp.SubjectID, tmp.NeedMonitoring, tmp.TopicsCount, tmp.CommentsCount,
		tmp.Interval, tmp.SendTo, tmp.LastDate)
	_, err = dbKit.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func (tmp *TopicMonitorParam) selectFromDBBySubjectID(subjectID int) error {
	var dbKit DataBaseKit
	err := dbKit.openDB()
	defer dbKit.db.Close()
	if err != nil {
		return err
	}

	query := fmt.Sprintf("SELECT * FROM topic_monitor WHERE subject_id=%d", subjectID)
	rows, err := dbKit.db.Query(query)
	defer rows.Close()
	if err != nil {
		return err
	}

	for rows.Next() {
		err = rows.Scan(&tmp.ID, &tmp.SubjectID,
			&tmp.NeedMonitoring, &tmp.TopicsCount,
			&tmp.CommentsCount, &tmp.Interval,
			&tmp.SendTo, &tmp.LastDate)
		if err != nil {
			return err
		}
	}

	return nil
}

func (tmp *TopicMonitorParam) updateInDB() error {
	var dbKit DataBaseKit
	err := dbKit.openDB()
	defer dbKit.db.Close()
	if err != nil {
		return err
	}

	query := fmt.Sprintf(`UPDATE topic_monitor 
		SET subject_id='%d', need_monitoring='%d', send_to='%d', 
		comments_count='%d', last_date='%d', interval='%d',
		topics_count='%d' WHERE id=%d`,
		tmp.SubjectID, tmp.NeedMonitoring, tmp.SendTo,
		tmp.CommentsCount, tmp.LastDate, tmp.Interval,
		tmp.TopicsCount, tmp.ID)
	_, err = dbKit.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func (tmp *TopicMonitorParam) updateInDBFieldLastDate(subjectID, newLastDate int) error {
	var dbKit DataBaseKit
	err := dbKit.openDB()
	defer dbKit.db.Close()
	if err != nil {
		return err
	}

	query := fmt.Sprintf(`UPDATE topic_monitor SET last_date=%d WHERE subject_id=%d`,
		newLastDate, subjectID)
	_, err = dbKit.db.Exec(query)
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

func (wpcmp *WallPostCommentMonitorParam) insertIntoDB() error {
	var dbKit DataBaseKit
	err := dbKit.openDB()
	defer dbKit.db.Close()
	if err != nil {
		return err
	}

	query := fmt.Sprintf(`INSERT INTO wall_post_comment_monitor (subject_id, 
		need_monitoring, posts_count, comments_count, monitoring_all, 
   		users_ids_for_monitoring, users_names_for_monitoring, attachments_types_for_monitoring, 
		users_ids_for_ignore, interval, send_to, filter, last_date, 
		keywords_for_monitoring, small_comments_for_monitoring, 
		digits_count_for_card_number_monitoring, 
		digits_count_for_phone_number_monitoring, monitor_by_community)
		VALUES ('%d', '%d', '%d', '%d', '%d', '%v', '%v', '%v', '%v', '%d', 
		'%d', '%v', '%d', '%v', '%v', '%v', '%v', '%d')`,
		wpcmp.SubjectID, wpcmp.NeedMonitoring, wpcmp.PostsCount, wpcmp.CommentsCount,
		wpcmp.MonitoringAll, wpcmp.UsersIDsForMonitoring, wpcmp.UsersNamesForMonitoring,
		wpcmp.AttachmentsTypesForMonitoring, wpcmp.UsersIDsForIgnore, wpcmp.Interval,
		wpcmp.SendTo, wpcmp.Filter, wpcmp.LastDate, wpcmp.KeywordsForMonitoring,
		wpcmp.SmallCommentsForMonitoring, wpcmp.DigitsCountForCardNumberMonitoring,
		wpcmp.DigitsCountForPhoneNumberMonitoring, wpcmp.MonitorByCommunity)
	_, err = dbKit.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func (wpcmp *WallPostCommentMonitorParam) selectFromDBBySubjectID(subjectID int) error {
	var dbKit DataBaseKit
	err := dbKit.openDB()
	defer dbKit.db.Close()
	if err != nil {
		return err
	}

	query := fmt.Sprintf("SELECT * FROM wall_post_comment_monitor WHERE subject_id=%d", subjectID)
	rows, err := dbKit.db.Query(query)
	defer rows.Close()
	if err != nil {
		return err
	}

	for rows.Next() {
		err = rows.Scan(&wpcmp.ID, &wpcmp.SubjectID,
			&wpcmp.NeedMonitoring, &wpcmp.PostsCount,
			&wpcmp.CommentsCount, &wpcmp.MonitoringAll,
			&wpcmp.UsersIDsForMonitoring,
			&wpcmp.UsersNamesForMonitoring,
			&wpcmp.AttachmentsTypesForMonitoring,
			&wpcmp.UsersIDsForIgnore,
			&wpcmp.Interval, &wpcmp.SendTo,
			&wpcmp.Filter, &wpcmp.LastDate,
			&wpcmp.KeywordsForMonitoring,
			&wpcmp.SmallCommentsForMonitoring,
			&wpcmp.DigitsCountForCardNumberMonitoring,
			&wpcmp.DigitsCountForPhoneNumberMonitoring,
			&wpcmp.MonitorByCommunity)
		if err != nil {
			return err
		}
	}

	return nil
}

func (wpcmp *WallPostCommentMonitorParam) updateInDB() error {
	var dbKit DataBaseKit
	err := dbKit.openDB()
	defer dbKit.db.Close()
	if err != nil {
		return err
	}

	query := fmt.Sprintf(`UPDATE wall_post_comment_monitor SET subject_id='%d', need_monitoring='%d',
		posts_count='%d', comments_count='%d', monitoring_all='%d', users_ids_for_monitoring='%v',
		users_names_for_monitoring='%v', attachments_types_for_monitoring='%v',
		users_ids_for_ignore='%v', interval='%d', send_to='%d', filter='%v',
		last_date='%d', keywords_for_monitoring='%v', small_comments_for_monitoring='%v',
		digits_count_for_card_number_monitoring='%v', digits_count_for_phone_number_monitoring='%v', 
		monitor_by_community='%d'
		WHERE id=%d`,
		wpcmp.SubjectID, wpcmp.NeedMonitoring,
		wpcmp.PostsCount, wpcmp.CommentsCount, wpcmp.MonitoringAll, wpcmp.UsersIDsForMonitoring,
		wpcmp.UsersNamesForMonitoring, wpcmp.AttachmentsTypesForMonitoring,
		wpcmp.UsersIDsForIgnore, wpcmp.Interval, wpcmp.SendTo, wpcmp.Filter,
		wpcmp.LastDate, wpcmp.KeywordsForMonitoring, wpcmp.SmallCommentsForMonitoring,
		wpcmp.DigitsCountForCardNumberMonitoring, wpcmp.DigitsCountForPhoneNumberMonitoring,
		wpcmp.MonitorByCommunity,
		wpcmp.ID)
	_, err = dbKit.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func (wpcmp *WallPostCommentMonitorParam) updateInDBFieldLastDate(subjectID, newLastDate int) error {
	var dbKit DataBaseKit
	err := dbKit.openDB()
	defer dbKit.db.Close()
	if err != nil {
		return err
	}

	// обновляем значения в конкретном поле
	query := fmt.Sprintf(`UPDATE wall_post_comment_monitor SET last_date=%d WHERE subject_id=%d`,
		newLastDate, subjectID)
	_, err = dbKit.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}
