package main

import (
	"github.com/gin-gonic/gin"
	"github.com/philipbo/gin_project_tmpl/gin_mid"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.Infoln("start", gin.Mode())

	router := gin.New()
	router.Use(gin_mid.Logger(), gin_mid.Recovery())

	mount(router)

	logrus.Infof("start web server on %s", ":8000")
	logrus.Fatal(router.Run(":8000"))

}
