package parseData

import (
	"io/ioutil"

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

type DiskResponceStruct struct {
	Usage struct {
		TS    int `json:"ts"`
		Disks struct {
			Sda diskSda `json:"/dev/sda"`
		} `json:"disks"`
	} `json:"usage"`
}
type diskSda struct {
	Filesystem string               `json:"filesystem"`
	Type       string               `json:"type"`
	Size       int                  `json:"size"`
	Used       int                  `json:"used"`
	Available  int                  `json:"available"`
	Capacity   float32              `json:"capacity"`
	Mountpoint string               `json:"mountpoint"`
	Contents   []contentsElemStruct `json:"contents"`
}

type contentsElemStruct struct {
	Type  string `json:"type"`
	Id    string `json:"id"`
	Path  string `json:"path"`
	Usage int    `json:"usage"`
}
