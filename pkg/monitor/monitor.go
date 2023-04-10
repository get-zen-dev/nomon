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
	"github.com/Setom29/CloudronMonitoring/pkg/report"
)

type Monitor struct {
	DB            *dbConn.DB
	F             Args
	R             report.Report
	WG            *sync.WaitGroup
	RAMTotal      uint64
	DiskTotal     uint64
	LastCheck     time.Time
	LastClearTime time.Time
}

type Args struct {
	CPULimit    float64 `yaml:"cpu_limit"`
	RAMLimit    float64 `yaml:"ram_limit"`
	DiskLimit   float64 `yaml:"disk_limit"`
	Duration    int     `yaml:"duration"`
	Port        int     `yaml:"port"`
	DBClearTime int     `yaml:"db_clear_time"`
	CheckTime   int
	DBFile      string
}

func NewMonitor(f Args, r report.Report) *Monitor {
	db, err := dbConn.NewDB(f.DBFile)
	if err != nil {
		log.Fatal(err)
	}
	ramTotal, _, diskTotal := getServerStats.GetTotalMetrics()
	t := time.Now().UTC()
	m := Monitor{
		DB: db,
		F:  f,
		R:  r,
		WG: &sync.WaitGroup{},

		RAMTotal:      ramTotal,
		DiskTotal:     diskTotal,
		LastCheck:     t,
		LastClearTime: t,
	}
	return &m
}

// StartMonitoring creates new db connection and pushes statistics to the database
func (monitor *Monitor) StartMonitoring(ch chan os.Signal) {
	log.Println("Start monitoring")

	for {
		select {
		case <-ch:
			monitor.DB.Close()
			os.Exit(0)
		default:

			stat := dbConn.ServerStatus{
				CPUUsed:  getServerStats.GetCpu(monitor.F.CheckTime),
				RAMUsed:  getServerStats.GetMem(),
				DiskUsed: getServerStats.GetDisk(),
				Time:     time.Now()}
			err := monitor.DB.Add(stat)
			if err != nil {
				log.Println("Error adding stats:", err)
			}
			monitor.Analyse()
			monitor.ClearDatabase()
		}
	}
}

// CheckStatus reads data from database and alerts if metric's usage limit exceeded
func (monitor *Monitor) Analyse() {
	if time.Now().UTC().Before(monitor.LastCheck.Add(time.Duration(monitor.F.Duration) * time.Second)) {
		return
	} else {
		monitor.LastCheck = time.Now().UTC()
	}
	rows, err := monitor.DB.Sql.Query(fmt.Sprintf("SELECT * FROM serverStatus WHERE time >= Datetime('now', '-%d seconds', 'localtime');", monitor.F.Duration))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// cumulative sum for metrics
	counter := 0
	cpuUsedCumSum, ramUsedCumSum, diskUsedCumSum := 0.0, 0.0, 0.0
	for rows.Next() {
		counter++
		stat := dbConn.ServerStatus{}
		err := rows.Scan(&stat.Time, &stat.CPUUsed, &stat.RAMUsed, &stat.DiskUsed)
		if err != nil {
			log.Println("Error scanning server status:", err)
			continue
		}
		cpuUsedCumSum += stat.CPUUsed
		ramUsedCumSum += float64(stat.RAMUsed)
		diskUsedCumSum += float64(stat.DiskUsed)
	}
	// alert check
	msg := ""
	if cpuStatus := cpuUsedCumSum / float64(counter); cpuStatus > monitor.F.CPULimit {
		msg += fmt.Sprintf("CPU usage limit exceeded: %f%%\n", cpuStatus)
	}
	if ramStatus := ramUsedCumSum / float64(counter); ramStatus/float64(monitor.RAMTotal)*100 > monitor.F.RAMLimit {
		msg += fmt.Sprintf("RAM usage limit exceeded: %f%% (%f GB /%f GB)\n",
			ramStatus/float64(monitor.RAMTotal)*100,
			ramStatus/math.Pow(1024, 3),
			float64(monitor.RAMTotal)/math.Pow(1024, 3))
	}
	if diskStatus := diskUsedCumSum / float64(counter); diskStatus/float64(monitor.DiskTotal)*100 > monitor.F.DiskLimit {
		msg += fmt.Sprintf("Disk usage limit exceeded: %f%% (%f GB /%f GB)",
			diskStatus/float64(monitor.DiskTotal)*100,
			diskStatus/math.Pow(1024, 3),
			float64(monitor.DiskTotal)/math.Pow(1024, 3))
	}
	if msg != "" {
		log.Println(msg)
		monitor.R.SendMessage(msg)
	}
}

func (monitor *Monitor) ClearDatabase() {
	if currTime := time.Now().UTC(); currTime.Hour() <= monitor.F.DBClearTime && monitor.F.DBClearTime <= currTime.Hour()+1 &&
		currTime.After(monitor.LastClearTime.AddDate(0, 0, 1)) {
		log.Println("Start clearing outdated values")
		_, err := monitor.DB.Sql.Exec(fmt.Sprintf("DELETE FROM serverStatus WHERE time < Datetime('now', '-%d seconds', 'localtime');", monitor.F.Duration))
		if err != nil {
			log.Println(err)
		}

	}
}
