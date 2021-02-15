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

func loadData() []Macro {
	jsonFile, err := os.Open("data.json")
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
