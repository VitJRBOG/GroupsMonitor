package main

import (
	"io"
	"io/ioutil"
	"os"
)

// ReadPathFile читает содержимое файла path.txt
func ReadPathFile() (string, error) {
	path, err := ReadTextFile("path.txt")
	if err != nil {
		return "", err
	}

	if len(path) > 0 {
		if string(path[len(path)-1]) != "/" {
			path += "/"
		}
	}

	return path, nil
}

// ReadTextFile читает текстовые файлы
func ReadTextFile(path string) (string, error) {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		return "", err
	}

	data := make([]byte, 64)
	text := ""

	for {
		n, err := file.Read(data)
		if err == io.EOF {
			break
		}
		text += string(data[:n])
	}

	return text, nil
}

// WriteTextFile сохраняет текст в текстовые файлы
func WriteTextFile(path string, text string) error {
	valuesBytes := []byte(text)
	err := ioutil.WriteFile(path, valuesBytes, 0644)
	if err != nil {
		return err
	}

	return nil
}
