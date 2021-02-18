package handler

import (
	"html/template"
	"net/http"

	"github.com/VitJRBOG/GroupsObserver/tools"
	"github.com/gorilla/mux"
)

func InitHandler() {
	rtr := mux.NewRouter()

	// General part //

	rtr.HandleFunc("/", homePageHandler).Methods("GET")

	// Observers part //

	rtr.HandleFunc("/observers", observersPageHandler).Methods("GET")
	rtr.HandleFunc("/observers/{id:[0-9]+}", observerControlPageHandler).Methods("GET", "POST")
	rtr.HandleFunc("/observers/{id:[0-9]+}/turn_on", observerTurnOnPageHandler).Methods("POST")
	rtr.HandleFunc("/observers/{id:[0-9]+}/turn_off", observerTurnOffPageHandler).Methods("POST")

	// Settings part //

	rtr.HandleFunc("/settings", settingsPageHandler).Methods("GET")

	rtr.HandleFunc("/settings/access_tokens", accessTokensPageHandler).Methods("GET", "POST")
	rtr.HandleFunc("/settings/access_tokens/new", accessTokenNewPageHandler).Methods("GET", "POST")
	rtr.HandleFunc("/settings/access_tokens/new/create",
		accessTokenCreateNewPageHandler).Methods("POST")
	rtr.HandleFunc("/settings/access_tokens/{id:[0-9]+}",
		accessTokenSettingsPageHandler).Methods("GET", "POST")
	rtr.HandleFunc("/settings/access_tokens/{id:[0-9]+}/delete",
		accessTokenDeletePageHandler).Methods("POST")
	rtr.HandleFunc("/settings/access_tokens/{id:[0-9]+}/update",
		accessTokenUpdatePageHandler).Methods("POST")

	rtr.HandleFunc("/settings/operators", operatorsPageHandler).Methods("GET", "POST")
	rtr.HandleFunc("/settings/operators/new", operatorNewPageHandler).Methods("GET", "POST")
	rtr.HandleFunc("/settings/operators/new/create", operatorCreateNewPageHandler).Methods("POST")
	rtr.HandleFunc("/settings/operators/{id:[0-9]+}",
		operatorSettingsPageHandler).Methods("GET", "POST")
	rtr.HandleFunc("/settings/operators/{id:[0-9]+}/delete",
		operatorDeletePageHandler).Methods("POST")
	rtr.HandleFunc("/settings/operators/{id:[0-9]+}/update",
		operatorUpdatePageHandler).Methods("POST")

	rtr.HandleFunc("/settings/wards", wardsPageHandler).Methods("GET", "POST")
	rtr.HandleFunc("/settings/wards/new", wardNewPageHandler).Methods("GET", "POST")
	rtr.HandleFunc("/settings/wards/new/create",
		wardCreateNewPageHandler).Methods("POST")
	rtr.HandleFunc("/settings/wards/{id:[0-9]+}",
		wardSettingsPageHandler).Methods("GET", "POST")
	rtr.HandleFunc("/settings/wards/{id:[0-9]+}/delete",
		wardDeletePageHandler).Methods("POST")
	rtr.HandleFunc("/settings/wards/{id:[0-9]+}/update",
		wardUpdatePageHandler).Methods("POST")

	//

	pathToResourcesWebview := tools.GetPath("ui/webview/")

	http.Handle("/static/", http.StripPrefix("/static/",
		http.FileServer(http.Dir(pathToResourcesWebview+"./static/"))))
	http.Handle("/", rtr)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err.Error())
	}
}

func getHtmlTemplates() *template.Template {
	pathToResourcesWebview := tools.GetPath("ui/webview/")

	t, err := template.ParseFiles(
		pathToResourcesWebview+"html/header.html",
		pathToResourcesWebview+"html/general/index.html",
		pathToResourcesWebview+"html/observers/observers.html",
		pathToResourcesWebview+"html/observers/observer_control.html",
		pathToResourcesWebview+"html/settings/settings.html",
		pathToResourcesWebview+"html/settings/access_tokens.html",
		pathToResourcesWebview+"html/settings/access_token_new.html",
		pathToResourcesWebview+"html/settings/access_token_settings.html",
		pathToResourcesWebview+"html/settings/operators.html",
		pathToResourcesWebview+"html/settings/operator_new.html",
		pathToResourcesWebview+"html/settings/operator_settings.html",
		pathToResourcesWebview+"html/settings/wards.html",
		pathToResourcesWebview+"html/settings/ward_new.html",
		pathToResourcesWebview+"html/settings/ward_settings.html",
		pathToResourcesWebview+"html/footer.html")
	if err != nil {
		panic(err.Error())
	}
	return t
}
