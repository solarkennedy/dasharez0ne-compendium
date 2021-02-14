package main

import (
	"github.com/gin-gonic/gin"
)

func index(c *gin.Context) {
	data := gin.H{
		"Page index (here)":               "/",
		"Specific macro (id: 1463183460)": "/macro/1463183460",
	}
	switch c.NegotiateFormat(gin.MIMEHTML, gin.MIMEJSON) {
	case gin.MIMEHTML:
		c.HTML(200, "index.tmpl", data)
	case gin.MIMEJSON:
		c.JSON(200, data)
	}
}

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/assets", "./assets")

	r.GET("/", index)
	// r.GET("/search", search)
	// r.GET("/tag/:tag_name", tag)
	// r.GET("/tags", tags)
	// r.GET("/macro/:id", macro)
	return r
}

func main() {
	r := SetupRouter()
	_ = r.Run()
}
