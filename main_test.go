package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func performRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func performJSONRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	req.Header.Add("Accept", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
func TestIndexJSON(t *testing.T) {
	router := SetupRouter()
	w := performJSONRequest(router, "GET", "/")
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)

	value, exists := response["Page index (here)"]
	assert.True(t, exists)
	assert.Equal(t, "/", value)
}

func TestIndexHTML(t *testing.T) {
	router := SetupRouter()
	w := performRequest(router, "GET", "/")
	assert.Equal(t, http.StatusOK, w.Code)

	body := w.Body.String()
	assert.Contains(t, body, "Welcome to dasharez0ne Compendium!")
	assert.Contains(t, body, "Page index")
}

func TestMacroJSON(t *testing.T) {
	router := SetupRouter()
	w := performJSONRequest(router, "GET", "/macro/1463183460")
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)

	value, exists := response["original_text"]
	assert.True(t, exists)
	assert.Equal(t, "FOLLOW ON FACebook https://t.co/68sSXQ7VBJ Da Motha F**CKin Share Z0ne...", value)
}

func TestMacroHTML(t *testing.T) {
	router := SetupRouter()
	w := performRequest(router, "GET", "/macro/1463183460")
	assert.Equal(t, http.StatusOK, w.Code)

	body := w.Body.String()
	assert.Contains(t, body, "I EAT ASS HOLES LIKE YOU")
}
