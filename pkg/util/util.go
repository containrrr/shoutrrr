package util

import (
	"log"

	"github.com/onsi/ginkgo"
)

func Min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

func Max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

func TestLogger() *log.Logger {
	return log.New(ginkgo.GinkgoWriter, "Test", log.LstdFlags)
}