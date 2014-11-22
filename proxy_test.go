package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func extractBodyValue(w *httptest.ResponseRecorder, k string) string {
	values := map[string]string{}
	json.NewDecoder(w.Body).Decode(&values)
	return values[k]
}

func TestProxy(t *testing.T) {
	assert := assert.New(t)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "https://glaze.com?url=http://httpbin.org/ip", nil)

	proxy(w, req)

	assert.Equal(w.Code, 200, "status should be ok")
	assert.NotEmpty(w.Body)
}

func TestProxyWithoutUrl(t *testing.T) {
	assert := assert.New(t)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "https://glaze.com", nil)

	proxy(w, req)

	assert.Equal(w.Code, 400, "status should be bad request")
	assert.Equal(w.Body.String(), "No request url to proxy\n")
}

func TestWriteResponseCopiesBody(t *testing.T) {
	assert := assert.New(t)
	w := httptest.NewRecorder()
	resp, _ := http.Get("http://httpbin.org/ip") // { "origin": "xxx.xx.xx.xxx" }

	writeResponse(w, resp)

	assert.Equal(w.Code, 200, "status should be ok")
	assert.NotEmpty(extractBodyValue(w, "origin"), "body should contain origin")
}
