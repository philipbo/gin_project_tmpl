package main

import "github.com/gin-gonic/gin"

func mount(router *gin.Engine) {
	router.GET("/", func(c *gin.Context) {
		c.String(200, "hello world\n")
	})

	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong\n")
	})
}
