package handler

import "net/http"

func observersPageHandler(w http.ResponseWriter, _ *http.Request) {
	t := getHtmlTemplates()
	err := t.ExecuteTemplate(w, "observers", nil)
	if err != nil {
		panic(err.Error())
	}
}
