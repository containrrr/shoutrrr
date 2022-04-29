package testutils

import (
	"github.com/containrrr/shoutrrr/pkg/types"

	Ω "github.com/onsi/gomega"
)

// TestServiceSetInvalidParamValue tests whether the service returns an error when an invalid param key/value is passed through Send
func TestServiceSetInvalidParamValue(service types.Service, key string, value string) {
	err := service.Send("TestMessage", &types.Params{key: value})
	Ω.ExpectWithOffset(1, err).To(Ω.HaveOccurred())
}
