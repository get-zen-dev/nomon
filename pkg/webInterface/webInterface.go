package webInterface

import (
	"net/http"
	"text/template"
)

type IndexData struct {
	Message   string
	LastCheck string
}

func IndexHandler(message, lastCheck string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// define a struct that contains the values to be passed to the HTML template
		data := struct {
			Message   string
			LastCheck string
		}{
			Message:   message,
			LastCheck: lastCheck,
		}

		// parse the HTML template
		tmpl, err := template.ParseFiles("./index.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// execute the HTML template with the struct instance
		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
