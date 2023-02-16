package monitor

import (
	"fmt"
	"log"
	"math"
	"os"
	"sync"
	"time"

	"github.com/Setom29/CloudronMonitoring/pkg/dbConn"
	"github.com/Setom29/CloudronMonitoring/pkg/getServerStats"
)

type Monitor struct {
	DB        *dbConn.DB
	F         Args
	WG        *sync.WaitGroup
	RAMTotal  uint64
	DiskTotal uint64
}

type Args struct {
	CPULimit  float64 `yaml:"cpulimit"`
	RAMLimit  float64 `yaml:"ramlimit"`
	DiskLimit float64 `yaml:"disklimit"`
	Duration  int     `yaml:"duration"`
	CheckTime int     `yaml:"checktime"`
	DBFile    string
}

func NewMonitor(f Args) *Monitor {
	db, err := dbConn.NewDB("sqlite.db")
	if err != nil {
		log.Fatal(err)
	}
	ramTotal, _, diskTotal := getServerStats.GetTotalMetrics()
	m := Monitor{
		DB: db,
		F:  f,
		WG: &sync.WaitGroup{},

		RAMTotal:  ramTotal,
		DiskTotal: diskTotal,
	}
	return &m
}

// StartMonitoring creates new db connection and pushes statistics to the database
func (monitor *Monitor) StartMonitoring(ch chan os.Signal) {
	defer monitor.WG.Done()
	log.Println("Starting monitoring")
	db, err := dbConn.NewDB(monitor.F.DBFile)

	go monitor.CheckStatus()
	go monitor.ClearDatabase()

	if err != nil {
		log.Fatal(err)
	}
	for {
		select {
		case <-ch:
			db.Close()
			return
		default:
			stat := dbConn.ServerStatus{
				CPUUsed:  getServerStats.GetCpu(monitor.F.CheckTime),
				RAMUsed:  getServerStats.GetMem(),
				DiskUsed: getServerStats.GetDisk(),
				Time:     time.Now()}
			if err != nil {
				log.Println(err)
			}
			err = monitor.DB.Add(stat)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

// CheckStatus reads data from database and alerts if metric's usage limit exceeded
func (monitor *Monitor) CheckStatus() {
	for {
		time.Sleep(time.Duration(monitor.F.Duration) * time.Second)
		rows, err := monitor.DB.Sql.Query(fmt.Sprintf("SELECT * FROM serverStatus WHERE time >= Datetime('now', '-%d seconds', 'localtime');", monitor.F.Duration))
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		counter := 0
		cpuUsedCumSum, ramUsedCumSum, diskUsedCumSum := 0.0, 0.0, 0.0
		for rows.Next() {
			counter++
			stat := dbConn.ServerStatus{}
			err := rows.Scan(&stat.Time, &stat.CPUUsed, &stat.RAMUsed, &stat.DiskUsed)
			if err != nil {
				fmt.Println(err)
				continue
			}
			cpuUsedCumSum += stat.CPUUsed
			ramUsedCumSum += float64(stat.RAMUsed)
			diskUsedCumSum += float64(stat.DiskUsed)
		}
		if cpuStatus := cpuUsedCumSum / float64(counter); cpuStatus > monitor.F.CPULimit {
			log.Printf("CPU usage limit exceeded: %f%%\n", cpuStatus)
		}
		if ramStatus := ramUsedCumSum / float64(counter); ramStatus/float64(monitor.RAMTotal)*100 > monitor.F.RAMLimit {
			log.Printf("RAM usage limit exceeded: %f%% (%f GB /%f GB)\n", ramStatus/float64(monitor.RAMTotal)*100, ramStatus/math.Pow(1024, 3), float64(monitor.RAMTotal)/math.Pow(1024, 3))
		}
		if diskStatus := diskUsedCumSum / float64(counter); diskStatus/float64(monitor.DiskTotal)*100 > monitor.F.DiskLimit {
			log.Printf("Disk usage limit exceeded: %f%% (%f GB /%f GB)\n", diskStatus/float64(monitor.DiskTotal)*100, diskStatus/math.Pow(1024, 3), float64(monitor.DiskTotal)/math.Pow(1024, 3))
		}
	}

}

func (monitor *Monitor) ClearDatabase() {
	time.Sleep(time.Hour)
	if currTime := time.Now(); currTime.Hour() < 3 && currTime.Hour() > 2 {
		monitor.DB.Mutex.Lock()
		defer monitor.DB.Mutex.Unlock()
		monitor.DB.Sql.Exec(fmt.Sprintf("DELETE FROM serverStatus WHERE time < Datetime('now', '-%d seconds', 'localtime');", monitor.F.Duration))

	}
}
