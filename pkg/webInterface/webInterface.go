package webInterface

import (
	"net/http"
	"text/template"
	"time"

	"github.com/Setom29/CloudronMonitoring/pkg/monitor"
	log "github.com/sirupsen/logrus"
)

type IndexData struct {
	Message   string
	LastCheck time.Time
}

func MakeIndexHandler(m *monitor.Monitor) http.HandlerFunc {
	log.Trace("webInterface:MakeIndexHandler")
	return func(w http.ResponseWriter, r *http.Request) {
		// define a struct that contains the values to be passed to the HTML template
		data := IndexData{
			Message:   m.Message,
			LastCheck: m.LastCheck,
		}

		// parse the HTML template
		tmpl, err := template.ParseFiles("./index.html")
		if err != nil {
			log.Error("Error creating template:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// execute the HTML template with the struct instance
		err = tmpl.Execute(w, data)
		if err != nil {
			log.Error("Error applying template:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Debug("http.HandlerFunc created.")
	}

}
