package main_test

import (
	"io"
	"net/http"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var hostName = "http://localhost:8080" // Make configurable to test in local/QE or prod environments
var resp *http.Response
var err error

var _ = Describe("HelloWorldService", func() {
	Context("The /api/hello-world-service/v0/ping endpoint responds successfully", func() {		
		BeforeEach(func() {			
			if (os.Getenv("APP_BASE_URL") != "") {
				hostName = os.Getenv("APP_BASE_URL")
			}
			resp, err = http.Get(hostName + "/api/hello-world-service/v0/ping")
		})

		// Verify status code
		It ("/api/hello-world-service/v0/ping responds with 200", func () {
			defer resp.Body.Close()
			Expect(resp.StatusCode).To(Equal(200))
			Expect(err).NotTo(HaveOccurred())
		})		

		// Verify response
		It ("/api/hello-world-service/v0/ping serves up Hello World", func () {
			expectedResponse := "Hello world"
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			Expect(string(body)).To(Equal(expectedResponse))
			Expect(err).NotTo(HaveOccurred())
		})				
	})	
})
