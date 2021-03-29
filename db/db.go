package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"runtime/debug"

	"github.com/VitJRBOG/GroupsObserver/tools"
	_ "github.com/mattn/go-sqlite3"
)

type AccessToken struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

func (a *AccessToken) InsertIntoDB() {
	dbase := openDB()
	defer func() {
		err := dbase.Close()
		if err != nil {
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
	}()

	query := fmt.Sprintf(`INSERT INTO access_token (name, value) VALUES ('%s', '%s')`,
		a.Name, a.Value)
	_, err := dbase.Exec(query)
	if err != nil {
		tools.WriteToLog(err, debug.Stack())
		panic(err.Error())
	}
}

func (a *AccessToken) SelectByID(id int) {
	dbase := openDB()
	defer func() {
		err := dbase.Close()
		if err != nil {
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
	}()

	query := fmt.Sprintf("SELECT * FROM access_token WHERE id=%d", id)
	rows := sendSelectQuery(dbase, query)
	defer func() {
		err := rows.Close()
		if err != nil {
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
	}()

	for rows.Next() {
		err := rows.Scan(&a.ID, &a.Name, &a.Value)
		if err != nil {
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
	}
}

func (a *AccessToken) SelectByName(name string) {
	dbase := openDB()
	defer func() {
		err := dbase.Close()
		if err != nil {
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
	}()

	query := fmt.Sprintf("SELECT * FROM access_token WHERE name='%s'", name)
	rows := sendSelectQuery(dbase, query)
	defer func() {
		err := rows.Close()
		if err != nil {
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
	}()

	for rows.Next() {
		err := rows.Scan(&a.ID, &a.Name, &a.Value)
		if err != nil {
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
	}
}

func (a *AccessToken) UpdateInDB() {
	dbase := openDB()
	defer func() {
		err := dbase.Close()
		if err != nil {
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
	}()

	query := fmt.Sprintf(`UPDATE access_token SET name='%s', value='%s' WHERE id=%d`,
		a.Name, a.Value, a.ID)
	_, err := dbase.Exec(query)
	if err != nil {
		tools.WriteToLog(err, debug.Stack())
		panic(err.Error())
	}
}

func (a *AccessToken) DeleteFromDB() {
	dbase := openDB()
	defer func() {
		err := dbase.Close()
		if err != nil {
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
	}()

	query := fmt.Sprintf(`DELETE FROM access_token WHERE id=%d`, a.ID)
	_, err := dbase.Exec(query)
	if err != nil {
		tools.WriteToLog(err, debug.Stack())
		panic(err.Error())
	}
}

func SelectAccessTokens() []AccessToken {
	var accessTokens []AccessToken

	dbase := openDB()
	defer func() {
		err := dbase.Close()
		if err != nil {
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
	}()

	query := "SELECT * FROM access_token"
	rows := sendSelectQuery(dbase, query)
	defer func() {
		err := rows.Close()
		if err != nil {
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
	}()

	for rows.Next() {
		var a AccessToken
		err := rows.Scan(&a.ID, &a.Name, &a.Value)
		if err != nil {
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
		accessTokens = append(accessTokens, a)
	}

	return accessTokens
}

type Operator struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	VkID int    `json:"vk_id"`
}

func (o *Operator) InsertIntoDB() {
	dbase := openDB()
	defer func() {
		err := dbase.Close()
		if err != nil {
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
	}()

	query := fmt.Sprintf(`INSERT INTO operator (name, vk_id) VALUES ('%s', %d)`,
		o.Name, o.VkID)
	_, err := dbase.Exec(query)
	if err != nil {
		tools.WriteToLog(err, debug.Stack())
		panic(err.Error())
	}
}

func (o *Operator) UpdateInDB() {
	dbase := openDB()
	defer func() {
		err := dbase.Close()
		if err != nil {
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
	}()

	query := fmt.Sprintf(`UPDATE operator SET name='%s', vk_id=%d WHERE id=%d`,
		o.Name, o.VkID, o.ID)
	_, err := dbase.Exec(query)
	if err != nil {
		tools.WriteToLog(err, debug.Stack())
		panic(err.Error())
	}
}

func (o *Operator) SelectByID(id int) {
	dbase := openDB()
	defer func() {
		err := dbase.Close()
		if err != nil {
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
	}()

	query := fmt.Sprintf("SELECT * FROM operator WHERE id=%d", id)
	rows := sendSelectQuery(dbase, query)
	defer func() {
		err := rows.Close()
		if err != nil {
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
	}()

	for rows.Next() {
		err := rows.Scan(&o.ID, &o.Name, &o.VkID)
		if err != nil {
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
	}
}

func (o *Operator) SelectByName(name string) {
	dbase := openDB()
	defer func() {
		err := dbase.Close()
		if err != nil {
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
	}()

	query := fmt.Sprintf("SELECT * FROM operator WHERE name='%s'", name)
	rows := sendSelectQuery(dbase, query)
	defer func() {
		err := rows.Close()
		if err != nil {
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
	}()

	for rows.Next() {
		err := rows.Scan(&o.ID, &o.Name, &o.VkID)
		if err != nil {
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
	}
}

func (o *Operator) DeleteFromDB() {
	dbase := openDB()
	defer func() {
		err := dbase.Close()
		if err != nil {
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
	}()

	query := fmt.Sprintf(`DELETE FROM operator WHERE id=%d`, o.ID)
	_, err := dbase.Exec(query)
	if err != nil {
		tools.WriteToLog(err, debug.Stack())
		panic(err.Error())
	}
}

func SelectOperators() []Operator {
	var operators []Operator

	dbase := openDB()
	defer func() {
		err := dbase.Close()
		if err != nil {
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
	}()

	query := "SELECT * FROM operator"
	rows := sendSelectQuery(dbase, query)
	defer func() {
		err := rows.Close()
		if err != nil {
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
	}()

	for rows.Next() {
		var o Operator
		err := rows.Scan(&o.ID, &o.Name, &o.VkID)
		if err != nil {
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
		operators = append(operators, o)
	}

	return operators
}

type Ward struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	VkID             int    `json:"vk_id"`
	IsOwned          int    `json:"is_owned"`
	UnderObservation int    `json:"under_observation"`
	LastTS           int    `json:"last_ts"`
	GetAccessTokenID int    `json:"get_access_token_id"`
}

func (w *Ward) InsertIntoDB() {
	dbase := openDB()
	defer func() {
		err := dbase.Close()
		if err != nil {
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
	}()

	query := fmt.Sprintf("INSERT INTO ward (name, vk_id, is_owned, under_observation, "+
		"last_ts, get_access_token_id) VALUES ('%s', %d, %d, %d, %d, %d)",
		w.Name, w.VkID, w.IsOwned, w.UnderObservation, w.LastTS, w.GetAccessTokenID)
	_, err := dbase.Exec(query)
	if err != nil {
		tools.WriteToLog(err, debug.Stack())
		panic(err.Error())
	}
}

func (w *Ward) SelectByID(id int) {
	dbase := openDB()
	defer func() {
		err := dbase.Close()
		if err != nil {
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
	}()

	query := fmt.Sprintf("SELECT * FROM ward WHERE id=%d", id)
	rows := sendSelectQuery(dbase, query)
	defer func() {
		err := rows.Close()
		if err != nil {
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
	}()

	for rows.Next() {
		err := rows.Scan(&w.ID, &w.Name, &w.VkID, &w.IsOwned,
			&w.UnderObservation, &w.LastTS, &w.GetAccessTokenID)
		if err != nil {
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
	}
}

func (w *Ward) SelectByName(name string) {
	dbase := openDB()
	defer func() {
		err := dbase.Close()
		if err != nil {
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
	}()

	query := fmt.Sprintf("SELECT * FROM ward WHERE name='%s'", name)
	rows := sendSelectQuery(dbase, query)
	defer func() {
		err := rows.Close()
		if err != nil {
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
	}()

	for rows.Next() {
		err := rows.Scan(&w.ID, &w.Name, &w.VkID, &w.IsOwned,
			&w.UnderObservation, &w.LastTS, &w.GetAccessTokenID)
		if err != nil {
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
	}
}

func (w *Ward) UpdateInDB() {
	dbase := openDB()
	defer func() {
		err := dbase.Close()
		if err != nil {
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
	}()

	query := fmt.Sprintf(`UPDATE ward SET name='%s', vk_id=%d, is_owned=%d,  
		under_observation=%d, last_ts=%d, get_access_token_id=%d 
		WHERE id=%d`,
		w.Name, w.VkID, w.IsOwned, w.UnderObservation, w.LastTS, w.GetAccessTokenID, w.ID)
	sendUpdateQuery(dbase, query)
}

func (w *Ward) DeleteFromDB() {
	dbase := openDB()
	defer func() {
		err := dbase.Close()
		if err != nil {
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
	}()

	query := fmt.Sprintf(`DELETE FROM ward WHERE id=%d`, w.ID)
	_, err := dbase.Exec(query)
	if err != nil {
		tools.WriteToLog(err, debug.Stack())
		panic(err.Error())
	}
}

func SelectWards() []Ward {
	var wards []Ward

	dbase := openDB()
	defer func() {
		err := dbase.Close()
		if err != nil {
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
	}()

	query := "SELECT * FROM ward"
	rows := sendSelectQuery(dbase, query)
	defer func() {
		err := rows.Close()
		if err != nil {
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
	}()

	for rows.Next() {
		var w Ward
		err := rows.Scan(&w.ID, &w.Name, &w.VkID, &w.IsOwned,
			&w.UnderObservation, &w.LastTS, &w.GetAccessTokenID)
		if err != nil {
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
		wards = append(wards, w)
	}

	return wards
}

type Observer struct {
	ID                int    `json:"id"`
	Name              string `json:"name"`
	WardID            int    `json:"ward_id"`
	OperatorID        int    `json:"operator_id"`
	SendAccessTokenID int    `json:"send_access_token_id"`
	AdditionalParams  observerAdditionalParams
}

type observerAdditionalParams struct {
	WallPost wallPostObserverAdditionalParams
}

type wallPostObserverAdditionalParams struct {
	PostType string `json:"post_type"`
}

func (o *Observer) InsertIntoDB() {
	dbase := openDB()
	defer func() {
		err := dbase.Close()
		if err != nil {
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
	}()

	additionalParams := o.additionalParamsToJSON()

	query := fmt.Sprintf(`INSERT INTO observer (name, ward_id, operator_id, 
		send_access_token_id, additional_params) VALUES ('%s', %d, %d, %d, '%s')`,
		o.Name, o.WardID, o.OperatorID, o.SendAccessTokenID, additionalParams)
	_, err := dbase.Exec(query)
	if err != nil {
		tools.WriteToLog(err, debug.Stack())
		panic(err.Error())
	}
}

func (o *Observer) SelectByID(id int) {
	dbase := openDB()
	defer func() {
		err := dbase.Close()
		if err != nil {
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
	}()

	query := fmt.Sprintf("SELECT * FROM observer WHERE id=%d", id)
	rows := sendSelectQuery(dbase, query)
	defer func() {
		err := rows.Close()
		if err != nil {
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
	}()

	var additionalParams string
	for rows.Next() {
		err := rows.Scan(&o.ID, &o.Name, &o.WardID, &o.OperatorID,
			&o.SendAccessTokenID, &additionalParams)
		if err != nil {
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
	}
	o.parseAdditionalParams(additionalParams)
}

func (o *Observer) SelectByNameAndWardID(name string, wardID int) {
	dbase := openDB()
	defer func() {
		err := dbase.Close()
		if err != nil {
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
	}()

	query := fmt.Sprintf("SELECT * FROM observer WHERE name='%s' AND ward_id=%d", name, wardID)
	rows := sendSelectQuery(dbase, query)
	defer func() {
		err := rows.Close()
		if err != nil {
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
	}()

	var additionalParams string
	for rows.Next() {
		err := rows.Scan(&o.ID, &o.Name, &o.WardID, &o.OperatorID,
			&o.SendAccessTokenID, &additionalParams)
		if err != nil {
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
	}
	o.parseAdditionalParams(additionalParams)
}

func (o *Observer) UpdateInDB() {
	dbase := openDB()
	defer func() {
		err := dbase.Close()
		if err != nil {
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
	}()

	additionalParams := o.additionalParamsToJSON()

	query := fmt.Sprintf(`UPDATE observer SET name='%s', ward_id=%d, operator_id=%d, 
		send_access_token_id=%d, additional_params='%s' WHERE id=%d`,
		o.Name, o.WardID, o.OperatorID,
		o.SendAccessTokenID, additionalParams, o.ID)
	sendUpdateQuery(dbase, query)
}

func (o *Observer) DeleteFromDB() {
	dbase := openDB()
	defer func() {
		err := dbase.Close()
		if err != nil {
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
	}()

	query := fmt.Sprintf(`DELETE FROM observer WHERE id=%d`, o.ID)
	_, err := dbase.Exec(query)
	if err != nil {
		tools.WriteToLog(err, debug.Stack())
		panic(err.Error())
	}
}

func (o *Observer) parseAdditionalParams(additionalParams string) {
	if o.Name == "wall_post" {
		values := []byte(additionalParams)
		err := json.Unmarshal(values, &o.AdditionalParams.WallPost)
		if err != nil {
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
	}
}

func (o *Observer) additionalParamsToJSON() string {
	if o.AdditionalParams.WallPost.PostType != "" {
		p := fmt.Sprintf(`{"post_type": "%s"}`, o.AdditionalParams.WallPost.PostType)
		return p
	}
	return "{}"
}

func Initialization() bool {
	isExist := checkDBDirectory()
	if isExist {
		isExist := checkDBExistence()
		if isExist {
			return false
		}
		initDB()
		return true
	}
	initDBDir()
	return Initialization()
}

func checkDBDirectory() bool {
	path := getPathToDBDir()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func checkDBExistence() bool {
	path := getPathToDB()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func initDBDir() {
	path := getPathToDBDir()
	err := os.Mkdir(path, 0700)
	if err != nil {
		tools.WriteToLog(err, debug.Stack())
		panic(err.Error())
	}
}

func initDB() {
	dbase := openDB()
	defer func() {
		err := dbase.Close()
		if err != nil {
			tools.WriteToLog(err, debug.Stack())
			panic(err.Error())
		}
	}()

	query := `BEGIN TRANSACTION;
	CREATE TABLE IF NOT EXISTS "access_token" (
		"id"	INTEGER NOT NULL UNIQUE,
		"name"	TEXT NOT NULL UNIQUE,
		"value"	TEXT NOT NULL,
		PRIMARY KEY("id" AUTOINCREMENT)
	);
	CREATE TABLE IF NOT EXISTS "observer" (
		"id"	INTEGER NOT NULL UNIQUE,
		"name"	TEXT NOT NULL,
		"ward_id"	INTEGER NOT NULL,
		"operator_id"	INTEGER NOT NULL,
		"send_access_token_id"	INTEGER NOT NULL,
		"additional_params"	TEXT NOT NULL DEFAULT '{}',
		FOREIGN KEY("send_access_token_id") REFERENCES "access_token"("id"),
		FOREIGN KEY("ward_id") REFERENCES "ward"("id"),
		FOREIGN KEY("operator_id") REFERENCES "observer"("id"),
		PRIMARY KEY("id" AUTOINCREMENT)
	);
	CREATE TABLE IF NOT EXISTS "operator" (
		"id"	INTEGER NOT NULL UNIQUE,
		"name"	TEXT NOT NULL UNIQUE,
		"vk_id"	INTEGER NOT NULL,
		PRIMARY KEY("id" AUTOINCREMENT)
	);
	CREATE TABLE IF NOT EXISTS "ward" (
		"id"	INTEGER NOT NULL UNIQUE,
		"name"	TEXT NOT NULL UNIQUE,
		"vk_id"	INTEGER NOT NULL UNIQUE,
		"is_owned"	INTEGER NOT NULL,
		"under_observation"	INTEGER NOT NULL,
		"last_ts"	INTEGER NOT NULL,
		"get_access_token_id"	INTEGER NOT NULL,
		FOREIGN KEY("get_access_token_id") REFERENCES "access_token"("id"),
		PRIMARY KEY("id" AUTOINCREMENT)
	);
	COMMIT;`
	sendUpdateQuery(dbase, query)
}

func sendSelectQuery(dbase *sql.DB, query string) *sql.Rows {
	rows, err := dbase.Query(query)
	if err != nil {
		tools.WriteToLog(err, debug.Stack())
		panic(err.Error())
	}
	return rows
}

func sendUpdateQuery(dbase *sql.DB, query string) sql.Result {
	resp, err := dbase.Exec(query)
	if err != nil {
		tools.WriteToLog(err, debug.Stack())
		panic(err.Error())
	}
	return resp
}

func openDB() *sql.DB {
	path := getPathToDB()
	dbase, err := sql.Open("sqlite3", path)
	if err != nil {
		tools.WriteToLog(err, debug.Stack())
		panic(err.Error())
	}
	return dbase
}

func getPathToDBDir() string {
	absPathToDBDir := tools.GetPath("data")
	return absPathToDBDir
}

func getPathToDB() string {
	absPathToDB := tools.GetPath("data/groups_observer.db")
	return absPathToDB
}
