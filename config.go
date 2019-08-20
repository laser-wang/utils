package utils

import (
	"fmt"

	"github.com/go-ini/ini"
)

const (
	CFG_SEP = ","
)

var (
	config *Config
)

type Config struct {
	profile string
}

func NewConfig(profile string) *Config {
	o := &Config{
		profile,
	}
	return o
}

func (self *Config) Load() (*ini.File, error) {
	// 判断文件是否存在
	//	if err := syscall.Access(self.profile, syscall.F_OK); err != nil {
	//		return nil, fmt.Errorf("%s: %s", err.Error(), self.profile)
	//	}
	if FileIsExist(self.profile) == false {
		return nil, fmt.Errorf("%s: %s", "文件不存在", self.profile)
	}
	return ini.LooseLoad(self.profile)
}

func MakeDefaultConfig(cfg *Config) {
	if cfg != nil {
		config = cfg
	}
}

func DefaultConfig() *Config {
	return config
}
