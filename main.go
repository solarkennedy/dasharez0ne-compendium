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
		"API Docs (swagger)":              "/api/index.html",
		"Tag list (all tags)":             "/tags",
		"Tag Example (acrostic)":          "/tag/acrostic",
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

// tags godoc
// @Summary Shows tags
// @Description Show all tags
// @Accept  json
// @Produce  json
// @Success 200 {object} object
// @Router / [get]
func tags(c *gin.Context) {
	md, _ := c.MustGet("MacroData").(MacroData)
	tags := md.getTags()
	switch c.NegotiateFormat(gin.MIMEHTML, gin.MIMEJSON) {
	case gin.MIMEHTML:
		data := gin.H{
			"tags":      tags,
			"full_path": FullURL(c),
		}
		c.HTML(200, "tags.tmpl", data)
	case gin.MIMEJSON:
		c.JSON(200, tags)
	}
}

// tag godoc
// @Summary Shows all macros with a particular tag
// @Description Shows all macros with a particular tag
// @Accept  json
// @Produce  json
// @Param tag path string true "Tag name"
// @Success 200 {object} object
// @Router / [get]
func tag(c *gin.Context) {
	md, _ := c.MustGet("MacroData").(MacroData)
	tagName := c.Param("tagName")
	tagged := md.getTagged(tagName)
	switch c.NegotiateFormat(gin.MIMEHTML, gin.MIMEJSON) {
	case gin.MIMEHTML:
		tagHash := md.getExamplesOf(tagged)
		data := gin.H{
			"number":    len(tagHash),
			"tag":       tagName,
			"tagHash":   tagHash,
			"full_path": FullURL(c),
		}
		c.HTML(200, "tag.tmpl", data)
	case gin.MIMEJSON:
		c.JSON(200, tagged)
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
	r.GET("/tags", tags)
	r.GET("/tag/:tagName", tag)

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
