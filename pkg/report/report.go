package report

import (
	"fmt"

	"github.com/containrrr/shoutrrr"
)

type Report struct {
	Service           string `yaml:"service"`
	MatrixAccessToken string `yaml:"matrix_access_token"`
	MatrixRoomID      string `yaml:"matrix_room_id"`
	MatrixHostServer  string `yaml:"matrix_host_server"`
	URL               string
	Message           string
}

func (r *Report) SendMessage() {
	r.MakeURL()
	r.Report()

}

func (r *Report) MakeURL() {
	if r.Service == "matrix" {
		r.URL = fmt.Sprintf("matrix://%s@%s/?rooms=%s", r.MatrixAccessToken, r.MatrixHostServer, r.MatrixRoomID)
	}
}

func (r *Report) Report() {
	shoutrrr.Send(r.URL, r.Message)
}
