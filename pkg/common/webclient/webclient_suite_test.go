package webclient_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestWebClient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "WebClient Suite")
}
