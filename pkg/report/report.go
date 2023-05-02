package report

import (
	"fmt"

	"github.com/containrrr/shoutrrr"
	log "github.com/sirupsen/logrus"
)

type Report struct {
	Service           string `yaml:"service"`
	MatrixAccessToken string `yaml:"matrix_access_token"`
	MatrixRoomID      string `yaml:"matrix_room_id"`
	MatrixHostServer  string `yaml:"matrix_host_server"`
	URL               string
	Message           string
}

func (r *Report) SendMessage(msg string) {
	log.Trace("report:SendMessage")
	r.Message = msg
	r.MakeURL()
	r.Report()

}

func (r *Report) MakeURL() {
	log.Trace("report:MakeURL")
	if r.Service == "matrix" {
		r.URL = fmt.Sprintf("matrix://:%s@%s/?rooms=%s", r.MatrixAccessToken, r.MatrixHostServer, r.MatrixRoomID)
	}
}

func (r *Report) Report() {
	log.Trace("report:Report")
	if err := shoutrrr.Send(r.URL, r.Message); err != nil {
		log.Error("Error sending report, ", err)
	}
	log.Trace("Message sent")
}
