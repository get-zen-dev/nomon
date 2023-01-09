package main

import (
	"log"
	"os"

	"github.com/Setom29/CloudronMonitoring/pkg/parseYML"
)

func main() {
	cfg, err := parseYML.ParseYAML("./config/config.yml")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	log.Println(cfg.Domain)

	// requests.MakeRequest(cfg.Domain + cfg.Links[0] + )

}
