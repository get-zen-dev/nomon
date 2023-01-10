package cloudronStatus

import (
	"log"
	"os"
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
	wg          *sync.WaitGroup
}

func NewCloudronServerConnection(cfg parseData.Cfg, wg *sync.WaitGroup) *CloudronServer {

	return &CloudronServer{Domain: cfg.Domain, CpuStatus: 0, RamStatus: 0, DiskStatus: 0, LastUpdated: time.Now(), cfg: cfg, wg: wg}
}

func StartMonitoring(filepath string) {
	var wg sync.WaitGroup

	cfg, err := parseData.ParseConfig(filepath)
	if err != nil {
		log.Println("Error parsing config file")
		os.Exit(1)
	}

	cloudronServ := NewCloudronServerConnection(cfg, &wg)
	wg.Add(3)
	go cloudronServ.getCpuStatus()
	go cloudronServ.getDiskStatus()
	go cloudronServ.getRamStatus()

	wg.Wait()

}

func (cloudronServ *CloudronServer) getCpuStatus() []byte {
	defer cloudronServ.wg.Done()
	for {
		res, err := requests.MakeRequest(cloudronServ.cfg.Domain + cloudronServ.cfg.CpuUrl + "&access_token=" + cloudronServ.cfg.Token)
		if err != nil {
			log.Println("Handling error from getCpuStatus():", err)
		} else {
			log.Println("CPU")
			return res
		}
	}
}

func (cloudronServ *CloudronServer) getDiskStatus() []byte {
	defer cloudronServ.wg.Done()
	for {
		res, err := requests.MakeRequest(cloudronServ.cfg.Domain + cloudronServ.cfg.DiskUrl + "?access_token=" + cloudronServ.cfg.Token)
		if err != nil {
			log.Println("Handling error from getDiskStatus():", err)
		} else {
			log.Println("Disk")
			return res
		}
	}
}

func (cloudronServ *CloudronServer) getRamStatus() []byte {
	defer cloudronServ.wg.Done()
	for {
		res, err := requests.MakeRequest(cloudronServ.cfg.Domain + cloudronServ.cfg.RamUrl + "?access_token=" + cloudronServ.cfg.Token)
		if err != nil {
			log.Println("Handling error from getRamStatus():", err)
		} else {
			log.Println("RAM")
			var parseData.DiskResponceStruct
		}
	}
}
