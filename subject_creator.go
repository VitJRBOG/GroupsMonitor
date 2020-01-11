package main

import (
	"fmt"
	"strconv"
)

// CreateSubject создает субъект и наполняет его параметрами
func CreateSubject() (int, error) {
	var subject Subject

	// получаем идентификатор субъекта
	subjectID, err := getSubjectID()
	if err != nil {
		return 0, err
	}
	subject.SubjectID = subjectID

	// получаем имя субъекта
	name, err := getName()
	if err != nil {
		return 0, err
	}
	subject.Name = name

	// получаем идентификатор вики-странички для бэкапа
	backupWikipage, err := getBackupWikipage()
	if err != nil {
		return 0, err
	}
	subject.BackupWikipage = backupWikipage

	// указываем дату последнего бэкапа
	subject.LastBackup = 0

	// создаем новое поле в таблице
	err = InsertDBSubject(subject)
	if err != nil {
		return 0, err
	}

	// запрашиваем обновленную таблицу с субъектами
	subjects, err := SelectDBSubjects()
	if err != nil {
		return 0, err
	}

	// извлекаем идентификатор БД для нового субъекта
	id := subjects[len(subjects)-1].ID

	return id, nil
}

// getSubjectID запрашивает у пользователя идентификатор субъекта в базе ВК
func getSubjectID() (int, error) {
	fmt.Print("> [Create subject -> Get subject ID]: ")
	var userAnswer string
	_, err := fmt.Scan(&userAnswer)
	if err != nil {
		return 0, err
	}
	subjectID, err := strconv.Atoi(userAnswer)
	if err != nil {
		return 0, err
	}
	return subjectID, nil
}

// getName запрашивает у пользователя имя субъекта для своей базы данных
func getName() (string, error) {
	fmt.Print("> [Create subject -> Get name]: ")
	var userAnswer string
	_, err := fmt.Scan(&userAnswer)
	if err != nil {
		return "", err
	}
	return userAnswer, nil
}

// getBackupWikipage
//запрашивает у пользователя идентификатор вики-страницы с бэкапом
func getBackupWikipage() (string, error) {
	fmt.Print("> [Create subject -> Get backup wikipage id]: ")
	var userAnswer string
	_, err := fmt.Scan(&userAnswer)
	if err != nil {
		return "", err
	}
	return userAnswer, nil
}
