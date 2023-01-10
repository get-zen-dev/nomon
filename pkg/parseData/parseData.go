package parseData

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v3"
)

type Cfg struct {
	Token   string `yaml:"Token"`
	Domain  string `yaml:"Domain"`
	DiskUrl string `yaml:"DiskUrl"`
	RamUrl  string `yaml:"RamUrl"`
	CpuUrl  string `yaml:"CpuUrl"`
	AppsUrl string `yaml:"AppsUrl"`
}

func ParseConfig(cfgFile string) (Cfg, error) {
	file, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		return Cfg{}, err
	}

	var cfg Cfg

	err = yaml.Unmarshal(file, &cfg)
	if err != nil {
		return Cfg{}, err
	}
	return cfg, nil
}

type contentsElemStruct struct {
	Type  string `json:"type"`
	Id    string `json:"id"`
	Path  string `json:"path"`
	Usage int    `json:"usage"`
}

func ParseStruct(arr []byte) {
	var st contentsElemStruct
	json.Unmarshal(arr, &st)
	log.Println(st)
}

// type DiskResponceStruct struct {
// }

// type Contents struct {
// 	contents []contentsElemStruct
// }
