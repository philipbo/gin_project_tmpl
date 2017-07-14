package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/philipbo/gin_project_tmpl/config"
	"github.com/philipbo/gin_project_tmpl/gin_mid"
	"github.com/sirupsen/logrus"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
	flag.StringVar(&config.ConfigFile, "config", "./config/config_dev.toml", "config file path")
}

func main() {
	logrus.Infoln("start", gin.Mode())

	flag.Parse()
	conf := config.DecodeConfig(config.ConfigFile)
	config.GConfig = conf

	router := gin.New()
	router.Use(gin_mid.Logger(), gin_mid.Recovery())

	mount(router)

	logrus.Infof("start web server on %s", conf.Common.Addr)
	logrus.Fatal(router.Run(conf.Common.Addr))

}
