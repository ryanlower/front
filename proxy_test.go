package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setup(t *testing.T) (*assert.Assertions, *Proxy, *httptest.ResponseRecorder) {
	config := Config{}
	return setupWithConfig(t, config)
}

func setupWithConfig(t *testing.T, config Config) (*assert.Assertions, *Proxy, *httptest.ResponseRecorder) {
	assert := assert.New(t)
	proxy := newProxy(config)
	recorder := httptest.NewRecorder()

	return assert, proxy, recorder
}

func extractBodyValue(w *httptest.ResponseRecorder, k string) string {
	values := map[string]string{}
	json.NewDecoder(w.Body).Decode(&values)
	return values[k]
}

func TestHandler(t *testing.T) {
	assert, p, w := setup(t)

	url := "https://glaze.com?url=http://httpbin.org/ip"
	req, _ := http.NewRequest("GET", url, nil)

	p.handler(w, req)

	assert.Equal(w.Code, 200, "status should be ok")
	assert.NotEmpty(w.Body)
}

func TestHandlerWithoutUrl(t *testing.T) {
	assert, p, w := setup(t)

	url := "https://glaze.com"
	req, _ := http.NewRequest("GET", url, nil)

	p.handler(w, req)

	assert.Equal(w.Code, http.StatusBadRequest, "status should be bad request")
	assert.Equal(w.Body.String(), "No request url to proxy\n")
}

func TestHandlerWithContentTypeRegexMatching(t *testing.T) {
	config := Config{allowedContentTypes: "^image/"}
	assert, p, w := setupWithConfig(t, config)

	// With matching Content-Type
	content_type := "image/png"
	req, _ := http.NewRequest("GET", "https://glaze.com?url=http://httpbin.org/response-headers?Content-Type=" + content_type, nil)

	p.handler(w, req)

	assert.Equal(w.Code, 200, "status should be ok")
}

func TestHandlerWithContentTypeRegexNotMatching(t *testing.T) {
	config := Config{allowedContentTypes: "^image/"}
	assert, p, w := setupWithConfig(t, config)

	// With matching Content-Type
	content_type := "text/plain"
	req, _ := http.NewRequest("GET", "https://glaze.com?url=http://httpbin.org/response-headers?Content-Type=" + content_type, nil)

	p.handler(w, req)

	assert.Equal(w.Code, http.StatusBadRequest, "status should be bad request")
}

func TestWriteResponseCopiesBody(t *testing.T) {
	assert, p, w := setup(t)

	resp, _ := http.Get("http://httpbin.org/ip") // { "origin": "xxx.xx.xx.xxx" }

	p.writeResponse(w, resp)

	assert.Equal(w.Code, 200, "status should be ok")
	assert.NotEmpty(extractBodyValue(w, "origin"), "body should contain origin")
}
