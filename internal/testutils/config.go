package testutils

import (
	"github.com/containrrr/shoutrrr/pkg/format"
	"net/url"

	Ω "github.com/onsi/gomega"

	"github.com/containrrr/shoutrrr/pkg/types"
)

// TestConfigGetInvalidQueryValue tests whether the config returns an error when an invalid query value is requested
func TestConfigGetInvalidQueryValue(config types.ServiceConfig) {
	value, err := format.GetConfigQueryResolver(config).Get("invalid query var")
	Ω.ExpectWithOffset(1, value).To(Ω.BeEmpty())
	Ω.ExpectWithOffset(1, err).To(Ω.HaveOccurred())
}

// TestConfigSetInvalidQueryValue tests whether the config returns an error when a URL with an invalid query value is parsed
func TestConfigSetInvalidQueryValue(config types.ServiceConfig, rawInvalidURL string) {
	invalidURL, err := url.Parse(rawInvalidURL)
	Ω.ExpectWithOffset(1, err).ToNot(Ω.HaveOccurred(), "the test URL did not parse correctly")

	err = config.SetURL(invalidURL)
	Ω.ExpectWithOffset(1, err).To(Ω.HaveOccurred())
}

// TestConfigGetEnumsCount tests whether the config.Enums returns the expected amount of items
func TestConfigGetEnumsCount(config types.ServiceConfig, expectedCount int) {
	enums := config.Enums()
	Ω.ExpectWithOffset(1, enums).To(Ω.HaveLen(expectedCount))
}

// TestConfigGetFieldsCount tests whether the config.QueryFields return the expected amount of fields
func TestConfigGetFieldsCount(config types.ServiceConfig, expectedCount int) {
	fields := format.GetConfigQueryResolver(config).QueryFields()
	Ω.ExpectWithOffset(1, fields).To(Ω.HaveLen(expectedCount))
}
