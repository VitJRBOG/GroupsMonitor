package main

import (
	"io/ioutil"
	"os"
)

// CheckFiles проверяет наличие файлов, необходимых для нормальной работы программы
func CheckFiles() (bool, error) {
	// проверяем наличие файла, в котором хранится лог-вывод программы
	err := checkLogFileExistence()
	if err != nil {
		return false, err
	}

	// проверяем наличие файла, в котором хранится путь к остальным файлам
	err = checkPathFileExistence()
	if err != nil {
		return false, err
	}

	// проверяем наличие файла БД
	dbHasBeenCreated, err := checkDBFileExistence()
	if err != nil {
		return false, err
	}

	return dbHasBeenCreated, nil
}

// checkLogFileExistence проверяет наличие файла, где хранится лог-вывод программы
func checkLogFileExistence() error {
	// проверяем
	if _, err := os.Stat("log.txt"); os.IsNotExist(err) {

		// если отсутствует, то создаем новый
		err = WriteTextFile("log.txt", "")
		if err != nil {
			return err
		}
	}

	return nil
}

// checkPathFile проверяет наличие файла, где хранится путь
func checkPathFileExistence() error {
	// проверяем
	if _, err := os.Stat("path.txt"); os.IsNotExist(err) {

		// если отсутствует, то создаем новый
		valuesBytes := []byte("")
		err = ioutil.WriteFile("path.txt", valuesBytes, 0644)
		if err != nil {
			return err
		}
	}

	return nil
}

// checkDBFileExistence проверяет наличие файла БД
func checkDBFileExistence() (bool, error) {
	hasBeenCreated := false

	path, err := ReadPathFile()
	if err != nil {
		return hasBeenCreated, err
	}
	// проверяем
	if _, err := os.Stat(path + "groupsmonitor_db.db"); os.IsNotExist(err) {
		// если файл БД отсутствует, создаем
		var dbKit DataBaseKit
		err = dbKit.initDB()
		if err != nil {
			return hasBeenCreated, err
		}
		hasBeenCreated = true
	}

	return hasBeenCreated, nil
}
