package testutils

import (
	"github.com/nicholas-fedor/shoutrrr/pkg/types"

	"github.com/onsi/gomega"
)

// TestServiceSetInvalidParamValue tests whether the service returns an error
// when an invalid param key/value is passed through Send.
func TestServiceSetInvalidParamValue(service types.Service, key string, value string) {
	err := service.Send("TestMessage", &types.Params{key: value})
	gomega.ExpectWithOffset(1, err).To(gomega.HaveOccurred())
}
