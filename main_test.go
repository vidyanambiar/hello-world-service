package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/identitatem/idp-configs-api/config"
	"github.com/onsi/gomega"
)

func TestHelloWorld(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	req, err := http.NewRequest("GET", "/api/idp-configs-api/v0/ping", nil);

	g.Expect(err).ShouldNot(gomega.HaveOccurred())

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(helloWorld)
    handler.ServeHTTP(responseRecorder, req)

    // Verify the status code (200)
	g.Expect(responseRecorder.Code).To(gomega.Equal(http.StatusOK))

    // Verify the response body (Hello World)
	g.Expect(responseRecorder.Body.String()).To(gomega.Equal("Hello world"))
}

func TestHealthCheck(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	req, err := http.NewRequest("GET", "/", nil);

	g.Expect(err).ShouldNot(gomega.HaveOccurred())

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(statusOK)
    handler.ServeHTTP(responseRecorder, req)

    // Verify the status code to be 200
	g.Expect(responseRecorder.Code).To(gomega.Equal(http.StatusOK))
	// Verify the response header
	g.Expect(responseRecorder.Header().Get("Content-Type")).To(gomega.Equal("text/plain; charset=utf-8"))
}

func TestOpenApiSpec(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	
	// Initialize config for test
	config.Init()

	req, err := http.NewRequest("GET", "/api/idp-configs-api/v0/openapi.json", nil);

	g.Expect(err).ShouldNot(gomega.HaveOccurred())

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(serveOpenAPISpec)
    handler.ServeHTTP(responseRecorder, req)

    // Verify the status code to be 200
	g.Expect(responseRecorder.Code).To(gomega.Equal(http.StatusOK))

	// Verify content of response (property "openapi" should have value 3.0.0)
	type OpenAPI struct {
		Openapi string	// OpenAPI version specified in the spec (3.0.0)
	}
	var openAPIspec OpenAPI	
	json.Unmarshal(responseRecorder.Body.Bytes(), &openAPIspec)
	
	g.Expect(openAPIspec.Openapi).To(gomega.Equal("3.0.0"))
}

