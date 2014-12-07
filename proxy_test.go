package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
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

func setupUpstreamServer(fileName string) *httptest.Server {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, fileName)
	})
	return httptest.NewServer(handler)
}

func TestHandler(t *testing.T) {
	assert, p, w := setup(t)
	server := setupUpstreamServer("test/images/gopher.png")
	defer server.Close()

	url := "https://front.com?url=" + server.URL
	req, _ := http.NewRequest("GET", url, nil)

	p.handler(w, req)

	assert.Equal(w.Code, 200, "status should be ok")
	assert.Equal(w.HeaderMap.Get("Content-Type"), "image/png")
	assert.NotEmpty(w.Body)
}

func TestHandlerWithoutUrl(t *testing.T) {
	assert, p, w := setup(t)

	url := "https://front.com"
	req, _ := http.NewRequest("GET", url, nil)

	p.handler(w, req)

	assert.Equal(w.Code, http.StatusBadRequest, "status should be bad request")
	assert.Equal(w.Body.String(), "No request url to proxy\n")
}

func TestHandlerWithContentTypeRegexMatching(t *testing.T) {
	config := Config{allowedContentTypes: "image/jpeg"}
	assert, p, w := setupWithConfig(t, config)
	server := setupUpstreamServer("test/images/money.jpg") // is image/jpeg
	defer server.Close()

	url := "https://front.com?url=" + server.URL
	req, _ := http.NewRequest("GET", url, nil)

	p.handler(w, req)

	assert.Equal(w.Code, 200, "status should be ok")
	assert.Equal(w.HeaderMap.Get("Content-Type"), "image/jpeg")
}

func TestHandlerWithContentTypeRegexNotMatching(t *testing.T) {
	config := Config{allowedContentTypes: "image/jpeg"}
	assert, p, w := setupWithConfig(t, config)
	server := setupUpstreamServer("test/images/gopher.png") // is image/png
	defer server.Close()

	url := "https://front.com?url=" + server.URL
	req, _ := http.NewRequest("GET", url, nil)

	p.handler(w, req)

	assert.Equal(w.Code, http.StatusBadRequest, "status should be bad request")
}

func TestProxyRequestWithNon200Body(t *testing.T) {
	assert, p, w := setup(t)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	}))
	defer server.Close()

	params := url.Values{}
	params.Set("url", server.URL)
	p.proxyRequest(w, params)

	assert.Equal(w.Code, http.StatusNotFound, "status should be not found")
}

func TestWriteResponseCopiesBody(t *testing.T) {
	assert, p, w := setup(t)
	server := setupUpstreamServer("test/hello.txt")
	defer server.Close()

	resp, _ := http.Get(server.URL)

	p.writeResponse(w, resp)

	assert.Equal(w.Code, 200, "status should be ok")
	assert.Equal(w.Body.String(), "Nothing to see here\n")
}
