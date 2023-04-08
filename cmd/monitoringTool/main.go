package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Setom29/CloudronMonitoring/pkg/monitor"
	"github.com/Setom29/CloudronMonitoring/pkg/parseConfig"
)

func main() {
	f, r, err := parseConfig.Parse("./data/config.yml")
	if err != nil {
		log.Println(err)
		return
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT)

	m := monitor.NewMonitor(f, r)
	go m.StartMonitoring(sigChan)

	http.HandleFunc("/", makeIndexHandler(m))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", f.Port), nil))
}

type IndexData struct {
	Message   string
	LastCheck string
}

func makeIndexHandler(m *monitor.Monitor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := IndexData{
			Message:   "Last check:",
			LastCheck: m.LastCheck.Format("2006-01-02 15:04:05 UTC"),
		}
		// parse the HTML template
		tmpl, err := template.ParseFiles("index.html")
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
