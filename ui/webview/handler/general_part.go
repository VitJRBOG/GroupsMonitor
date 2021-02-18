package handler

import "net/http"

func homePageHandler(w http.ResponseWriter, _ *http.Request) {
	t := getHtmlTemplates()
	err := t.ExecuteTemplate(w, "index", nil)
	if err != nil {
		panic(err.Error())
	}
}
