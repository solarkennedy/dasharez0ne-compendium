package main

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/solarkennedy/dasharez0ne-compendium/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const (
	canonicalURL = "dasharez0ne-compendium.fly.dev"
)

type Macro struct {
	Id           int      `json:"id"`
	OriginalText string   `json:"original_text"`
	Url          string   `json:"url"`
	Tags         []string `json:"tags"`
	Image        string   `json:"image"`
	Caption      string   `json:"caption"`
	EditURL      string   `json:"edit_url"`
}

func FullURL(c *gin.Context) string {
	return "https://" + canonicalURL + c.FullPath()
}

// Index godoc
// @Summary Show an index
// @Description Get an index of available urls
// @Accept  json
// @Produce  json
// @Success 200 {object} object
// @Router / [get]
func index(c *gin.Context) {
	data := gin.H{
		"Page index (here)":               "/",
		"Specific macro (id: 1463183460)": "/macro/1463183460",
		"API Docs":                        "/api/index.html",
	}
	switch c.NegotiateFormat(gin.MIMEHTML, gin.MIMEJSON) {
	case gin.MIMEHTML:
		data["full_path"] = FullURL(c)
		c.HTML(200, "index.tmpl", data)
	case gin.MIMEJSON:
		c.JSON(200, data)
	}
}

// macro godoc
// @Summary Shows a macro
// @Description Get a macro with all the associated data
// @Accept  json
// @Produce  json
// @Param id path int true "Macro ID"
// @Success 200 {object} Macro
// @Router / [get]
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
		"id":        m.Id,
		"macro":     m,
		"full_path": FullURL(c),
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
	r.LoadHTMLGlob("resources/templates/*")
	r.Static("/assets", "resources/assets")

	r.GET("/", index)
	r.GET("/macro/:id", macro)

	url := ginSwagger.URL("https://" + canonicalURL + "/api/doc.json") // The url pointing to API definition
	r.GET("/api/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	docs.SwaggerInfo.Title = "dasharez0ne Compendium API"
	docs.SwaggerInfo.Description = "This is da motherfucken share z0ne api"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = ""
	docs.SwaggerInfo.Host = canonicalURL
	docs.SwaggerInfo.Schemes = []string{"https"}

	return r
}

func main() {
	r := SetupRouter()
	_ = r.Run()
}
