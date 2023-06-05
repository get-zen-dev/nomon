package monitor

import (
	"fmt"
	"math"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/Setom29/CloudronMonitoring/pkg/dbConn"
	"github.com/Setom29/CloudronMonitoring/pkg/getServerStats"
	"github.com/containrrr/shoutrrr"
	log "github.com/sirupsen/logrus"
)

type Monitor struct {
	DB            *dbConn.DB
	Cfg           Config
	Counters      Counters
	Message       string
	WG            *sync.WaitGroup
	RAMTotal      uint64
	DiskTotal     uint64
	LastCheck     time.Time
	LastClearTime time.Time
}
type Counters struct {
	CPUCounter  int
	RAMCounter  int
	DiskCounter int
}

type Config struct {
	CpuAlertThreshold  float64  `yaml:"cpu_alert_threshold"`
	RamAlertThreshold  float64  `yaml:"ram_alert_threshold"`
	DiskAlertThreshold float64  `yaml:"disk_alert_threshold"`
	CheckEvery         int      `yaml:"check_every"`
	Port               int      `yaml:"port"`
	OldDataCleanup     int      `yaml:"old_data_cleanup"`
	LogLevel           int      `yaml:"log_level"`
	URLS               []string `yaml:"urls"`
	CPUCycles          int
	RAMCycles          int
	DiskCycles         int
	CheckTime          int
	DBFile             string
}

func NewMonitor(config Config) *Monitor {
	log.Trace("monitor:NewMonitor")
	db, err := dbConn.NewDB(config.DBFile)
	if err != nil {
		log.Fatal("Error creating database: ", err)
	}
	ramTotal, _, diskTotal := getServerStats.GetTotalMetrics()
	t := time.Now()
	m := Monitor{
		DB:       db,
		Cfg:      config,
		Counters: Counters{CPUCounter: 0, RAMCounter: 0, DiskCounter: 0},
		WG:       &sync.WaitGroup{},

		RAMTotal:      ramTotal,
		DiskTotal:     diskTotal,
		LastCheck:     t,
		LastClearTime: t.AddDate(0, 0, -1),
	}
	log.Debug("New monitor created")
	return &m
}

// StartMonitoring creates new db connection and pushes statistics to the database
func (monitor *Monitor) StartMonitoring(ch chan os.Signal) {
	log.Trace("monitor:StartMonitoring")
	log.Info("Start monitoring")
	for {
		select {
		case s := <-ch:
			log.Debug("Got signal: ", s)
			monitor.DB.Close()
			os.Exit(0)
		default:
			stat := dbConn.ServerStatus{
				CPUUsed:  getServerStats.GetCpu(monitor.Cfg.CheckTime),
				RAMUsed:  getServerStats.GetMem(),
				DiskUsed: getServerStats.GetDisk(),
				Time:     time.Now(),
			}
			log.Debug("Getting server status: ", stat)
			err := monitor.DB.Add(stat)
			if err != nil {
				log.Error("Error adding stats: ", err)
			}
			monitor.Analyse()
			monitor.ClearDatabase()
		}
	}
}

// CheckStatus reads data from database and alerts if metric's usage limit exceeded
func (monitor *Monitor) Analyse() {
	log.Trace("monitor:Analyse")
	if time.Now().Before(monitor.LastCheck.Add(time.Duration(monitor.Cfg.CheckEvery) * time.Second)) {
		return
	} else {
		monitor.LastCheck = time.Now()
	}
	rows, err := monitor.DB.Sql.Query(fmt.Sprintf("SELECT * FROM serverStatus WHERE time >= Datetime('now', '-%d seconds', 'localtime');", monitor.Cfg.CheckEvery))
	if err != nil {
		log.Fatal("Error getting rows from db: ", err)
	}
	defer rows.Close()
	log.Trace("Calculating the cumulative sum for metrics")
	// cumulative sum for metrics
	counter := 0
	cpuUsedCumSum, ramUsedCumSum, diskUsedCumSum := 0.0, 0.0, 0.0
	for rows.Next() {
		counter++
		stat := dbConn.ServerStatus{}
		err := rows.Scan(&stat.Time, &stat.CPUUsed, &stat.RAMUsed, &stat.DiskUsed)
		if err != nil {
			log.Error("Error scanning rows: ", err)
			continue
		}
		cpuUsedCumSum += stat.CPUUsed
		ramUsedCumSum += float64(stat.RAMUsed)
		diskUsedCumSum += float64(stat.DiskUsed)
	}

	// alert check
	msg := monitor.AnilizeStatistics(cpuUsedCumSum/float64(counter), ramUsedCumSum/float64(counter), diskUsedCumSum/float64(counter))
	log.Debug("Alert check message: ", msg)

	if msg != "" {
		monitor.Message = msg
		log.Warn(msg)
		for ind := range monitor.Cfg.URLS {
			if err := shoutrrr.Send(monitor.Cfg.URLS[ind], msg); err != nil {
				log.Error("Error sending report, ", err)
			}
		}
		log.Trace("Message sent")
	}
}

// AnilizeStatistics gets statistics values, checks for alerts and returns message for reporting
func (monitor *Monitor) AnilizeStatistics(cpuStat, ramStat, diskStat float64) string {
	log.Trace("monitor:AnilizeStatistics")
	log.Infof("CPU: %f%%, RAM: %f%%, Disk: %f%%", cpuStat, ramStat/float64(monitor.RAMTotal)*100, diskStat/float64(monitor.DiskTotal)*100)
	msgArr := make([]string, 0)

	// CPU
	if cpuStat > monitor.Cfg.CpuAlertThreshold {
		monitor.Counters.CPUCounter++
		if monitor.Counters.CPUCounter >= monitor.Cfg.CPUCycles {
			msgArr = append(msgArr, fmt.Sprintf("CPU usage limit exceeded: %f%%", cpuStat))
		}
	} else {
		if monitor.Counters.CPUCounter >= monitor.Cfg.CPUCycles {
			msgArr = append(msgArr, "CPU usage within normal limits.")
		}
		monitor.Counters.CPUCounter = 0
	}

	// RAM
	if ramStat/float64(monitor.RAMTotal)*100 > monitor.Cfg.RamAlertThreshold {
		monitor.Counters.RAMCounter++
		if monitor.Counters.RAMCounter >= monitor.Cfg.RAMCycles {
			msgArr = append(msgArr, fmt.Sprintf("RAM usage limit exceeded: %f%% (%f GB /%f GB)",
				ramStat/float64(monitor.RAMTotal)*100,
				ramStat/math.Pow(1024, 3),
				float64(monitor.RAMTotal)/math.Pow(1024, 3)))
		}
	} else {
		if monitor.Counters.RAMCounter >= monitor.Cfg.RAMCycles {
			msgArr = append(msgArr, "RAM usage within normal limits.")
		}
		monitor.Counters.RAMCounter = 0
	}

	// Disk
	if diskStat/float64(monitor.DiskTotal)*100 > monitor.Cfg.DiskAlertThreshold {
		monitor.Counters.DiskCounter++
		if monitor.Counters.DiskCounter >= monitor.Cfg.DiskCycles {
			msgArr = append(msgArr, fmt.Sprintf("Disk usage limit exceeded: %f%% (%f GB /%f GB)",
				diskStat/float64(monitor.DiskTotal)*100,
				diskStat/math.Pow(1024, 3),
				float64(monitor.DiskTotal)/math.Pow(1024, 3)))
		}
	} else {
		if monitor.Counters.DiskCounter >= monitor.Cfg.DiskCycles {
			msgArr = append(msgArr, "Disk usage within normal limits.")
		}
		monitor.Counters.DiskCounter = 0
	}

	return strings.Join(msgArr, "\n")
}

func (monitor *Monitor) ClearDatabase() {
	log.Trace("monitor:ClearDatabase")
	if currTime := time.Now(); currTime.Hour() >= monitor.Cfg.OldDataCleanup && currTime.Hour() < monitor.Cfg.OldDataCleanup+1 &&
		currTime.After(monitor.LastClearTime.AddDate(0, 0, 1)) {
		log.Info("Start clearing outdated values")
		_, err := monitor.DB.Sql.Exec(fmt.Sprintf("DELETE FROM serverStatus WHERE time < Datetime('now', '-%d seconds', 'localtime');", monitor.Cfg.CheckEvery))
		if err != nil {
			log.Error("Error clearing db: ", err)
		}
		monitor.LastClearTime = monitor.LastClearTime.AddDate(0, 0, 1)
		log.Debug("Database clear time: ", monitor.LastClearTime)
	}
}
