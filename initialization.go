package main

import (
	"io/ioutil"
	"os"
)

// CheckFiles проверяет наличие файлов, необходимых для нормальной работы программы
func CheckFiles() error {
	// проверяем наличие файла, в котором хранится путь к остальным файлам
	err := checkPathFileExistence()
	if err != nil {
		return err
	}

	// проверяем наличие файла БД
	err = checkDBFileExistence()
	if err != nil {
		return err
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
		sender := "Initialization"
		message := "File \"path.txt\" has been created."
		OutputMessage(sender, message)
	}

	return nil
}

// checkDBFileExistence проверяет наличие файла БД
func checkDBFileExistence() error {
	path, err := ReadPathFile()
	if err != nil {
		return err
	}
	// проверяем
	if _, err := os.Stat(path + "groupsmonitor_db.db"); os.IsNotExist(err) {
		// если файл БД отсутствует, создаем
		InitDB()
	}

	return nil
}
