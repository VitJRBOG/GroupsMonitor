package main

import (
	"io"
	"os"
)

// ReadPathFile читает содержимое файла path.txt
func ReadPathFile() (string, error) {
	path, err := ReadTextFile("path.txt")
	if err != nil {
		return "", err
	}

	if string(path[len(path)-1]) != "/" {
		path += "/"
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
