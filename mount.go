package main

import (
	"github.com/gin-gonic/gin"
	"github.com/philipbo/gin_project_tmpl/html_render"
)

func mount(router *gin.Engine) {
	htmlRender := HTMLRender.New(HTMLRender.Options{
		Layout:     "layouts/layout",
		Extensions: []string{".html"},
	})

	router.HTMLRender = htmlRender

	router.GET("/", func(c *gin.Context) {
		c.String(200, "hello world\n")
	})

	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong\n")
	})
}
