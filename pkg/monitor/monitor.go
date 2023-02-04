package monitor

import (
	"log"
	"os"
	"time"

	"github.com/Setom29/CloudronMonitoring/pkg/dbConn"
	"github.com/Setom29/CloudronMonitoring/pkg/getServerStats"
)

type Monitor struct {
	db       *dbConn.DB
	duration time.Duration
	dbFile   string
	limit    float64
}

func NewMonitor(f parseFlags.Flags) *Monitor {
	db, err := dbConn.NewDB(f.dbFile)
	if err != nil {
		log.Fatal(err)
	}
	m := Monitor{
		db:       db,
		duration: duration,
		dbFile:   dbFile,
		limit:    limit,
	}
	return &m
}

func (monitor *Monitor) StartMonitoring(ch chan os.Signal) {
	db, err := dbConn.NewDB(monitor.dbFile)
	if err != nil {
		log.Fatal(err)
	}
	for {
		select {
		case <-ch:
			db.Close()
			return
		default:
			stat := dbConn.ServerStatus{CPUStatus: getServerStats.GetCpu(monitor.duration), RAMStatus: getServerStats.GetMem(), DiskStatus: getServerStats.GetDisk()}
			stat.Time = time.Now()
			if err != nil {
				log.Println(err)
			}
			monitor.db.Add(stat)
		}
	}
}
