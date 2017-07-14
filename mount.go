package main

import (
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/philipbo/gin_project_tmpl/ctx"
	"github.com/philipbo/gin_project_tmpl/db"
	"github.com/philipbo/gin_project_tmpl/html_render"
)

func mount(router *gin.Engine, war string, ds *db.Ds) {
	htmlRender := HTMLRender.New(HTMLRender.Options{
		Layout:     "layouts/layout",
		Extensions: []string{".html"},
	})

	router.HTMLRender = htmlRender

	router.Use(func(c *gin.Context) {
		c.Set(ctx.CTX_DS, ds)
		c.Next()
	})

	router.Use(static.Serve("/", static.LocalFile(war, true)))

	router.GET("/", func(c *gin.Context) {
		c.String(200, "hello world\n")
	})

	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong\n")
	})
}
