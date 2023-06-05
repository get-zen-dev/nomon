package parseConfig

import (
	"errors"
	"io/ioutil"

	"log"

	"github.com/Setom29/CloudronMonitoring/pkg/monitor"
	"gopkg.in/yaml.v3"
)

func Parse(cfgFile string) (monitor.Config, error) {
	file, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		log.Println("Error reading config file: ", err)
		return monitor.Config{}, err
	}

	var cfg monitor.Config

	err = yaml.Unmarshal(file, &cfg)
	if err != nil {
		return monitor.Config{}, err
	}

	if err = validateArgs(cfg); err != nil {
		return monitor.Config{}, err
	}
	cfg.CPUCycles = 1
	cfg.RAMCycles = 1
	cfg.DiskCycles = 1
	cfg.CheckTime = 5
	cfg.DBFile = "./data/sqlite.db"
	return cfg, nil
}

func validateArgs(cfg monitor.Config) error {
	if cfg.CpuAlertThreshold > 100 || cfg.CpuAlertThreshold < 0 {
		return errors.New("wrong value for CPU limit")
	}
	if cfg.RamAlertThreshold > 100 || cfg.RamAlertThreshold < 0 {
		return errors.New("wrong value for RAM limit")
	}
	if cfg.DiskAlertThreshold > 100 || cfg.DiskAlertThreshold < 0 {
		return errors.New("wrong value for Disk limit")
	}
	if cfg.CheckEvery < 0 {
		return errors.New("wrong value for duration")
	}
	// Parse time value
	if cfg.OldDataCleanup > 23 || cfg.OldDataCleanup < 0 {
		return errors.New("wrong value for db_clear_time")
	}
	if cfg.LogLevel > 6 || cfg.LogLevel < 1 {
		return errors.New("wrong value for monitor_log_level")
	}

	return nil
}
