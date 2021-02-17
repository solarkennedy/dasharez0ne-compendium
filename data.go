package main

import (
	"encoding/json"
	"fmt"

	"io/ioutil"
	"os"

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

func loadData() []Macro {
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
	return macros
}

type MacroData struct {
	AllMacros []Macro
}

func NewMacroData() MacroData {
	macros := loadData()
	return MacroData{
		AllMacros: macros,
	}
}

func (md MacroData) getMacro(id int) (Macro, error) {
	fmt.Printf("Searching for id %d out of %d macros...\n", id, len(md.AllMacros))
	// TODO: don't iterate
	for _, m := range md.AllMacros {
		if m.Id == id {
			return m, nil
		}
	}
	return Macro{}, fmt.Errorf("Not found")
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

func (md MacroData) getTagged(tagName string) []int {
	t := []int{}
	for _, m := range md.AllMacros {
		if contains(m.Tags, tagName) {
			t = append(t, m.Id)
		}
	}
	return t
}

func (md MacroData) getExamplesOf(input []int) map[int]string {
	t := map[int]string{}
	for _, e := range input {
		m, _ := md.getMacro(e)
		t[e] = m.Image
	}
	return t
}

func contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}
