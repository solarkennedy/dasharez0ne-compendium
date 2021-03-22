package main

import (
	"sort"
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
		"Random macro":                    "/random",
		"Search Query (bearclaws)":        "/search/bearclaws",
		"All Macros":                      "/all",
	}
	switch c.NegotiateFormat(gin.MIMEHTML, gin.MIMEJSON) {
	case gin.MIMEHTML:
		md, _ := c.MustGet("MacroData").(MacroData)
		m := md.getRandomMacro()
		content := gin.H{
			"urls":        data,
			"macro":       m,
			"total_count": len(md.AllMacros),
			"full_path":   FullURL(c),
		}
		c.HTML(200, "index.tmpl", content)
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
// @Router /macro/{id} [get]
func macro(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.HTML(404, "", err)
	}

	md, _ := c.MustGet("MacroData").(MacroData)
	m, ok := md.getMacro(idInt)
	if !ok {
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Macro not found"})
		return
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

// random godoc
// @Summary Shows a random macro
// @Description Get a random macro with all the associated data
// @Accept  json
// @Produce  json
// @Success 200 {object} Macro
// @Router /random [get]
func random(c *gin.Context) {
	md, _ := c.MustGet("MacroData").(MacroData)
	m := md.getRandomMacro()
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

type tagRow struct {
	Tag     string
	Count   int
	Example Macro
}

// tags godoc
// @Summary Shows tags
// @Description Show all tags
// @Accept  json
// @Produce  json
// @Success 200 {object} object
// @Router /tags [get]
func tags(c *gin.Context) {
	md, _ := c.MustGet("MacroData").(MacroData)
	tags := md.getTags()
	switch c.NegotiateFormat(gin.MIMEHTML, gin.MIMEJSON) {
	case gin.MIMEHTML:
		tagRows := []tagRow{}
		keys := make([]string, 0, len(tags))
		for k := range tags {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, tag := range keys {
			count := tags[tag]
			row := tagRow{
				Tag:     tag,
				Count:   count,
				Example: md.GetRandomExample(tag),
			}
			tagRows = append(tagRows, row)
		}
		data := gin.H{
			"tag_rows":  tagRows,
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
// @Router /tag/{tag} [get]
func tag(c *gin.Context) {
	md, _ := c.MustGet("MacroData").(MacroData)
	tagName := c.Param("tagName")
	tagged := md.getTagged(tagName)
	switch c.NegotiateFormat(gin.MIMEHTML, gin.MIMEJSON) {
	case gin.MIMEHTML:
		data := gin.H{
			"number":    len(tagged),
			"tag":       tagName,
			"tagged":    tagged,
			"full_path": FullURL(c),
		}
		c.HTML(200, "tag.tmpl", data)
	case gin.MIMEJSON:
		c.JSON(200, tagged)
	}
}

// all godoc
// @Summary Shows all macros
// @Description Shows all macros
// @Accept  json
// @Produce  json
// @Success 200 {object} object
// @Router /all [get]
func all(c *gin.Context) {
	md, _ := c.MustGet("MacroData").(MacroData)
	switch c.NegotiateFormat(gin.MIMEHTML, gin.MIMEJSON) {
	case gin.MIMEHTML:
		data := gin.H{
			"number":    len(md.AllMacros),
			"macros":    md.AllMacros,
			"full_path": FullURL(c),
		}
		c.HTML(200, "all.tmpl", data)
	case gin.MIMEJSON:
		c.JSON(200, md.AllMacros)
	}
}

// search godoc
// @Summary Search macros
// @Description Shows all macros by keyword. Searches captions, original_text, and tags
// @Accept  json
// @Produce  json
// @Param keyword path string true "Keyword"
// @Success 200 {object} object
// @Router /search/{keyword} [get]
func search(c *gin.Context) {
	md, _ := c.MustGet("MacroData").(MacroData)
	keyword := c.Param("keyword")
	results, err := md.search(keyword)
	if err != nil {
		c.JSON(500, err.Error())
	}
	switch c.NegotiateFormat(gin.MIMEHTML, gin.MIMEJSON) {
	case gin.MIMEHTML:
		data := gin.H{
			"number":    len(results),
			"keyword":   keyword,
			"results":   results,
			"full_path": FullURL(c),
		}
		c.HTML(200, "search.tmpl", data)
	case gin.MIMEJSON:
		c.JSON(200, results)
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
	r.GET("/search/:keyword", search)
	r.GET("/random", random)
	r.GET("/all", all)

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
