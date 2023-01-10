package cloudronStatus

import (
	"encoding/json"
	"log"
	"os"
	"sync"
	"time"

	"github.com/Setom29/CloudronMonitoring/pkg/parseData"
	"github.com/Setom29/CloudronMonitoring/pkg/requests"
	"github.com/tidwall/gjson"
)

type CloudronServer struct {
	Domain      string
	CpuStatus   float32
	RamStatus   float32
	DiskStatus  float32
	LastUpdated time.Time
	cfg         parseData.Cfg
	wg          *sync.WaitGroup
}

// Create a new CloudronServer instance
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
	wg.Add(2)
	go cloudronServ.getCpuAndRamStatus()
	go cloudronServ.getDiskStatus()

	wg.Wait()

}

func (cloudronServ *CloudronServer) getCpuAndRamStatus() {
	defer cloudronServ.wg.Done()
	for {
		res, err := requests.MakeRequest(cloudronServ.cfg.Domain + cloudronServ.cfg.CpuUrl + "&access_token=" + cloudronServ.cfg.Token)
		if err != nil {
			log.Println("Handling error from getCpuStatus():", err)
		} else {
			cpu := gjson.ParseBytes(res).Get("cpu").Array()
			ram := gjson.ParseBytes(res).Get("ram").Array()
			log.Println(cpu)
			log.Println(ram)
			break
		}
	}
}

func (cloudronServ *CloudronServer) getDiskStatus() {
	defer cloudronServ.wg.Done()
	for {
		res, err := requests.MakeRequest(cloudronServ.cfg.Domain + cloudronServ.cfg.DiskUrl + "?access_token=" + cloudronServ.cfg.Token)
		if err != nil {
			log.Println("Handling error from getDiskStatus():", err)
		} else {
			var diskRespStruct parseData.DiskResponceStruct
			err = json.Unmarshal(res, &diskRespStruct)

			cloudronServ.DiskStatus = diskRespStruct.Usage.Disks.Sda.Capacity
			log.Println(cloudronServ.DiskStatus)
			break
		}
	}
}

// func (cloudronServ *CloudronServer) getRamStatus() []byte {
// 	defer cloudronServ.wg.Done()
// 	for {
// 		res, err := requests.MakeRequest(cloudronServ.cfg.Domain + cloudronServ.cfg.RamUrl + "?access_token=" + cloudronServ.cfg.Token)
// 		if err != nil {
// 			log.Println("Handling error from getRamStatus():", err)
// 		} else {
// 			log.Println("RAM")
// 			return res
// 		}
// 	}
// }
