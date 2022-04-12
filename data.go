package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"sort"

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

func annotateData(macros []*Macro) {
	for i := range macros {
		macros[i].EditURL = fmt.Sprintf("https://github.com/solarkennedy/dasharez0ne-compendium/wiki/%d/_edit", macros[i].Id)
	}
}

func loadData() map[int]*Macro {
	jsonFile, err := os.Open("resources/data.json")
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()
	macros := []*Macro{}
	b, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(b, &macros)
	if err != nil {
		panic(err)
	}
	annotateData(macros)
	m := map[int]*Macro{}
	for _, i := range macros {
		tags := i.Tags
		if !(contains(tags, "merch") || contains(tags, "misc") || contains(tags, "donation") || contains(tags, "retweet")) {
			m[i.Id] = i
		}
	}
	for k, v := range m {
		if v.DupeOf != 0 {
			orig := v.DupeOf
			if _, ok := m[orig]; ok {
				m[orig].Dupes = append(m[orig].Dupes, k)
			}
		}
	}
	fmt.Printf("Loaded %d macros\n", len(m))
	return m
}

type MacroData struct {
	AllMacros   map[int]*Macro
	SearchIndex bleve.Index
}

func NewMacroData() MacroData {
	macros := loadData()
	// mapping := bleve.NewIndexMapping()
	// index, err := bleve.New("", mapping)
	// if err != nil {
	// 	panic(err)
	// }
	// for _, m := range macros {
	// 	err = index.Index(fmt.Sprintf("%d", m.Id), m.Caption)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }
	// l, _ := index.DocCount()
	// fmt.Printf("Indexed %d documents\n", l)
	return MacroData{
		AllMacros:   macros,
		SearchIndex: nil,
	}
}

func (md MacroData) getMacro(id int) (*Macro, bool) {
	m, ok := md.AllMacros[id]
	return m, ok
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

func (md MacroData) getTagged(tagName string) []*Macro {
	t := []*Macro{}
	for _, m := range md.AllMacros {
		if contains(m.Tags, tagName) {
			t = append(t, m)
		}
	}
	sort.Slice(t, func(i, j int) bool {
		return t[i].Id < t[j].Id
	})
	return t
}

func (md MacroData) GetRandomExample(tagName string) *Macro {
	options := md.getTagged(tagName)
	i := rand.Intn(len(options))
	return options[i]
}

func (md MacroData) getRandomMacro() *Macro {
	for _, m := range md.AllMacros {
		return m
	}
	return md.AllMacros[1]
}

func (md MacroData) search(keyword string) ([]*Macro, error) {
	r := []*Macro{}
	return r, nil
	// query := bleve.NewFuzzyQuery(keyword)
	// searchRequest := bleve.NewSearchRequest(query)
	// searchResults, err := md.SearchIndex.Search(searchRequest)
	// if err != nil {
	// 	return nil, err
	// }
	// for _, hit := range searchResults.Hits {
	// 	idInt, _ := strconv.Atoi(hit.ID)
	// 	r = append(r, md.AllMacros[idInt])
	// }
	// return r, nil
}

func contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}
