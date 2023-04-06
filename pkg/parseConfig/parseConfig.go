package parseConfig

import (
	"errors"
	"io/ioutil"
	"time"

	"github.com/Setom29/CloudronMonitoring/pkg/monitor"
	"gopkg.in/yaml.v3"
)

func Parse(cfgFile string) (monitor.Args, error) {
	file, err := ioutil.ReadFile(cfgFile)

	if err != nil {
		return monitor.Args{}, err
	}

	var args monitor.Args

	err = yaml.Unmarshal(file, &args)
	if err != nil {
		return monitor.Args{}, err
	}
	if args.CPULimit > 100 || args.CPULimit < 0 {
		return monitor.Args{}, errors.New("wrong value for CPU limit")
	}
	if args.RAMLimit > 100 || args.RAMLimit < 0 {
		return monitor.Args{}, errors.New("wrong value for RAM limit")
	}
	if args.CPULimit > 100 || args.CPULimit < 0 {
		return monitor.Args{}, errors.New("wrong value for Disk limit")
	}
	if args.Duration < 0 {
		return monitor.Args{}, errors.New("wrong value for duration")
	}
	// Parse time value
	_, err = time.Parse("15:04:05", args.DBClearTime)
	if err != nil {
		return monitor.Args{}, errors.New("wrong value for db_clear_time")
	}
	args.CheckTime = 5
	args.DBFile = "./cfg/sqlite.db"
	return args, nil
}
