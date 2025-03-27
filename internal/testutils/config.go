package testutils

import (
	"net/url"

	"github.com/nicholas-fedor/shoutrrr/pkg/format"
	"github.com/nicholas-fedor/shoutrrr/pkg/types"

	"github.com/onsi/gomega"
)

// TestConfigGetInvalidQueryValue tests whether the config returns
// an error when an invalid query value is requested.
func TestConfigGetInvalidQueryValue(config types.ServiceConfig) {
	value, err := format.GetConfigQueryResolver(config).Get("invalid query var")
	gomega.ExpectWithOffset(1, value).To(gomega.BeEmpty())
	gomega.ExpectWithOffset(1, err).To(gomega.HaveOccurred())
}

// TestConfigSetInvalidQueryValue tests whether the config returns
// an error when a URL with an invalid query value is parsed.
func TestConfigSetInvalidQueryValue(config types.ServiceConfig, rawInvalidURL string) {
	invalidURL, err := url.Parse(rawInvalidURL)
	gomega.ExpectWithOffset(1, err).ToNot(gomega.HaveOccurred(), "the test URL did not parse correctly")

	err = config.SetURL(invalidURL)
	gomega.ExpectWithOffset(1, err).To(gomega.HaveOccurred())
}

// TestConfigSetDefaultValues tests whether setting the default values
// can be set for an empty config without any errors.
func TestConfigSetDefaultValues(config types.ServiceConfig) {
	pkr := format.NewPropKeyResolver(config)
	gomega.ExpectWithOffset(1, pkr.SetDefaultProps(config)).To(gomega.Succeed())
}

// TestConfigGetEnumsCount tests whether the config.Enums returns the expected amount of items.
func TestConfigGetEnumsCount(config types.ServiceConfig, expectedCount int) {
	enums := config.Enums()
	gomega.ExpectWithOffset(1, enums).To(gomega.HaveLen(expectedCount))
}

// TestConfigGetFieldsCount tests whether the config.QueryFields return the expected amount of fields.
func TestConfigGetFieldsCount(config types.ServiceConfig, expectedCount int) {
	fields := format.GetConfigQueryResolver(config).QueryFields()
	gomega.ExpectWithOffset(1, fields).To(gomega.HaveLen(expectedCount))
}
