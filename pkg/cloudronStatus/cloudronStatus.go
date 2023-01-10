package cloudronStatus

import (
	"log"
	"sync"
	"time"

	"github.com/Setom29/CloudronMonitoring/pkg/parseData"
	"github.com/Setom29/CloudronMonitoring/pkg/requests"
)

type CloudronServer struct {
	Domain      string
	CpuStatus   int
	RamStatus   int
	DiskStatus  int
	LastUpdated time.Time
	cfg         parseData.Cfg
}

func NewCloudronServerConnection(cfg parseData.Cfg) *CloudronServer {

	return &CloudronServer{Domain: cfg.Domain, CpuStatus: 0, RamStatus: 0, DiskStatus: 0, LastUpdated: time.Now(), cfg: cfg}
}

func StartMonitoring(cfg parseData.Cfg, wg *sync.WaitGroup) {
	defer wg.Done()
	cloudronServ := NewCloudronServerConnection(cfg)
	log.Println(cloudronServ.getDiskStatus())
}

func (cloudronServ *CloudronServer) getCpuStatus() string {
	for {
		res, err := requests.MakeRequest(cloudronServ.cfg.Domain + cloudronServ.cfg.CpuUrl + "&access_token=" + cloudronServ.cfg.Token)
		if err != nil {
			log.Println("Handling error from getCpuStatus():", err)
		} else {
			return res
		}
	}
}

func (cloudronServ *CloudronServer) getDiskStatus() string {
	for {
		res, err := requests.MakeRequest(cloudronServ.cfg.Domain + cloudronServ.cfg.DiskUrl + "?access_token=" + cloudronServ.cfg.Token)
		if err != nil {
			log.Println("Handling error from getDiskStatus():", err)
		} else {
			return res
		}
	}
}

func (cloudronServ *CloudronServer) getRamStatus() string {
	for {
		res, err := requests.MakeRequest(cloudronServ.cfg.Domain + cloudronServ.cfg.RamUrl + "?access_token=" + cloudronServ.cfg.Token)
		if err != nil {
			log.Println("Handling error from getRamStatus():", err)
		} else {
			return res
		}
	}
}
