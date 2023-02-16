package parseConfig

import (
	"errors"
	"io/ioutil"

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
	if args.CheckTime < 0 || args.CheckTime > args.Duration {
		return monitor.Args{}, errors.New("wrong value for check time")
	}
	args.DBFile = "sqlite.db"
	return args, nil
}
