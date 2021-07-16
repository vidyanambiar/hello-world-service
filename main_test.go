package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHelloWorldHandler(t *testing.T) {
	// Create a request to pass to the handler
	req, err := http.NewRequest("GET", "/api/hello-world-service/v0/ping", nil);

	if err != nil {
		t.Error("Received error in response", err);
	}

	// Create a ResponseRecorder to record the response 
	responseRecorder := httptest.NewRecorder()

	handler := http.HandlerFunc(helloWorld)

    handler.ServeHTTP(responseRecorder, req)

    // Verify the status code
    if status := responseRecorder.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: %v instead of %v",
            status, http.StatusOK)
    }

    // Verify the response body
    expectedResponse := "Hello world"
    if responseRecorder.Body.String() != expectedResponse {
        t.Errorf("handler returned unexpected body: %v instead of %v",
		responseRecorder.Body.String(), expectedResponse)
    }
}