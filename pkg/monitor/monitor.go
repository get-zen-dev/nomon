package monitor

import (
	"log"
	"os"
	"time"

	"github.com/Setom29/CloudronMonitoring/pkg/dbConn"
	"github.com/Setom29/CloudronMonitoring/pkg/flags"
	"github.com/Setom29/CloudronMonitoring/pkg/getServerStats"
)

type Monitor struct {
	DB *dbConn.DB
	F  flags.Flags
}

func NewMonitor(f flags.Flags) *Monitor {
	db, err := dbConn.NewDB(f.DBFile)
	if err != nil {
		log.Fatal(err)
	}
	m := Monitor{
		DB: db,
		F:  f,
	}
	return &m
}

func (monitor *Monitor) StartMonitoring(ch chan os.Signal) {
	log.Println("Starting monitoring")
	db, err := dbConn.NewDB(monitor.F.DBFile)
	if err != nil {
		log.Fatal(err)
	}
	for {
		select {
		case <-ch:
			db.Close()
			return
		default:
			stat := dbConn.ServerStatus{CPUStatus: getServerStats.GetCpu(monitor.F.CheckTime), RAMStatus: getServerStats.GetMem(), DiskStatus: getServerStats.GetDisk()}
			stat.Time = time.Now()
			if err != nil {
				log.Println(err)
			}
			monitor.DB.Add(stat)
		}
	}
}
