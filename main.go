package main

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type Macro struct {
	Id           int      `json:"id"`
	OriginalText string   `json:"original_text"`
	Url          string   `json:"url"`
	Tags         []string `json:"tags"`
	Image        string   `json:"image"`
	Caption      string   `json:"caption"`
}

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

func macro(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.HTML(404, "", err)
	}

	md, _ := c.MustGet("MacroData").(MacroData)
	m, err := md.getMacro(idInt)
	if err != nil {
		panic(err)
	}
	data := gin.H{
		"id":    m.Id,
		"macro": m,
	}
	switch c.NegotiateFormat(gin.MIMEHTML, gin.MIMEJSON) {
	case gin.MIMEHTML:
		c.HTML(200, "macro.tmpl", data)
	case gin.MIMEJSON:
		c.JSON(200, m)
	}
}

func SetupRouter() *gin.Engine {
	r := gin.Default()
	md := NewMacroData()
	r.Use(dataMiddleware(md))
	r.LoadHTMLGlob("templates/*")
	r.Static("/assets", "./assets")

	r.GET("/", index)
	// r.GET("/search", search)
	// r.GET("/tag/:tag_name", tag)
	// r.GET("/tags", tags)
	r.GET("/macro/:id", macro)
	return r
}

func main() {
	r := SetupRouter()
	_ = r.Run()
}
