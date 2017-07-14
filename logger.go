package main

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

const tsFmt = "2006-01-02 15:04:05"

func init() {
	if gin.Mode() == gin.DebugMode {
		log.SetFormatter(&log.TextFormatter{TimestampFormat: tsFmt, FullTimestamp: true})
	} else {
		log.SetFormatter(&log.JSONFormatter{TimestampFormat: tsFmt})
	}
}
