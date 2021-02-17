package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"

	"io/ioutil"
	"os"

	"github.com/blevesearch/bleve/v2"
	"github.com/gin-gonic/gin"
)

func dataMiddleware(md MacroData) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("MacroData", md)
		c.Next()
	}
}

func annotateData(macros []Macro) {
	for i := range macros {
		macros[i].EditURL = fmt.Sprintf("https://github.com/solarkennedy/dasharez0ne-compendium/wiki/%d/_edit", macros[i].Id)
	}
}

func loadData() map[int]Macro {
	jsonFile, err := os.Open("resources/data.json")
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()
	macros := []Macro{}
	b, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(b, &macros)
	if err != nil {
		panic(err)
	}
	annotateData(macros)
	fmt.Printf("Loaded %d macros\n", len(macros))
	m := map[int]Macro{}
	for _, i := range macros {
		m[i.Id] = i
	}
	return m
}

type MacroData struct {
	AllMacros   map[int]Macro
	SearchIndex bleve.Index
}

func NewMacroData() MacroData {
	macros := loadData()
	mapping := bleve.NewIndexMapping()
	index, err := bleve.New("", mapping)
	if err != nil {
		panic(err)
	}
	for _, m := range macros {
		err = index.Index(fmt.Sprintf("%d", m.Id), m.Caption)
		if err != nil {
			panic(err)
		}
	}
	l, _ := index.DocCount()
	fmt.Printf("Indexed %d documents\n", l)
	return MacroData{
		AllMacros:   macros,
		SearchIndex: index,
	}
}

func (md MacroData) getMacro(id int) (Macro, error) {
	m, ok := md.AllMacros[id]
	if !ok {
		return Macro{}, fmt.Errorf("Not found")
	}
	return m, nil
}

func (md MacroData) getTags() map[string]int {
	allTags := map[string]int{}
	for _, m := range md.AllMacros {
		for _, tag := range m.Tags {
			allTags[tag]++
		}
	}
	return allTags
}

func (md MacroData) getTagged(tagName string) []Macro {
	t := []Macro{}
	for _, m := range md.AllMacros {
		if contains(m.Tags, tagName) {
			t = append(t, m)
		}
	}
	return t
}

func (md MacroData) GetRandomExample(tagName string) Macro {
	options := md.getTagged(tagName)
	i := rand.Intn(len(options))
	return options[i]
}

func (md MacroData) getRandomMacro() Macro {
	for _, m := range md.AllMacros {
		return m
	}
	return md.AllMacros[1]
}

func (md MacroData) search(keyword string) ([]Macro, error) {
	r := []Macro{}
	query := bleve.NewMatchQuery(keyword)
	search := bleve.NewSearchRequest(query)
	search.Highlight = bleve.NewHighlight()
	searchResults, err := md.SearchIndex.Search(search)
	if err != nil {
		return nil, err
	}
	for _, hit := range searchResults.Hits {
		idInt, _ := strconv.Atoi(hit.ID)
		r = append(r, md.AllMacros[idInt])
	}
	return r, nil
}

func contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}
