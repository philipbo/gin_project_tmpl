package config

import (
	"github.com/BurntSushi/toml"
	"github.com/sirupsen/logrus"
)

var (
	ConfigFile string
	GConfig    Config
)

type Common struct {
	Base string
	Addr string
	War  string
}

type Config struct {
	Common Common
}

func DecodeConfig(fpath string) Config {
	var conf Config
	if _, err := toml.DecodeFile(fpath, &conf); err != nil {
		logrus.Fatalln("decoe config file error %v", err)
	}

	return conf
}
