package main

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestHelloWorldService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "HelloWorldService Suite")
}
