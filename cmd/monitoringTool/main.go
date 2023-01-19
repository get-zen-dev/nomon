package main

import (
	"github.com/Setom29/CloudronMonitoring/pkg/getServerStats"
)

func main() {
	getServerStats.StartMonitoring()
	// log.Println("Starting server on http://127.0.0.1:8080/")

	// r := newRouter()
	// log.Fatal(http.ListenAndServe(":8080", r))
}
