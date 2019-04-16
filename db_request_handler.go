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

// AccessToken - структура для полей из таблицы access_token
type AccessToken struct {
	ID    int
	Name  string
	Value string
}

// SelectDBAccessTokens читает и возвращает поля из таблицы access_token
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
