package testutils_test

import (
	"net/url"
	"testing"

	. "github.com/containrrr/shoutrrr/internal/testutils"
	"github.com/containrrr/shoutrrr/pkg/services/standard"
	"github.com/containrrr/shoutrrr/pkg/types"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestTestUtils(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "Shoutrrr TestUtils Suite")
}

var _ = Describe("the testutils package", func() {
	When("calling function TestLogger", func() {
		It("should not return nil", func() {
			Expect(TestLogger()).NotTo(Equal(nil))
		})
		It(`should have the prefix "[Test] "`, func() {
			Expect(TestLogger().Prefix()).To(Equal("[Test] "))
		})
	})

	Describe("Must helpers", func() {
		Describe("URLMust", func() {
			It("should panic when an invalid URL is passed", func() {
				failures := InterceptGomegaFailures(func() { URLMust(":") })
				Expect(failures).To(HaveLen(1))
			})
		})

		Describe("JSONRespondMust", func() {
			It("should panic when an invalid struct is passed", func() {
				notAValidJSONSource := func() {}
				failures := InterceptGomegaFailures(func() { JSONRespondMust(200, notAValidJSONSource) })
				Expect(failures).To(HaveLen(1))
			})
		})
	})

	Describe("Config test helpers", func() {
		var config dummyConfig
		BeforeEach(func() {
			config = dummyConfig{}
		})
		Describe("TestConfigSetInvalidQueryValue", func() {
			It("should fail when not correctly implemented", func() {
				failures := InterceptGomegaFailures(func() {
					TestConfigSetInvalidQueryValue(&config, "mock://host?invalid=value")
				})
				Expect(failures).To(HaveLen(1))
			})
		})

		Describe("TestConfigGetInvalidQueryValue", func() {
			It("should fail when not correctly implemented", func() {
				failures := InterceptGomegaFailures(func() {
					TestConfigGetInvalidQueryValue(&config)
				})
				Expect(failures).To(HaveLen(1))
			})
		})

		Describe("TestConfigSetDefaultValues", func() {
			It("should fail when not correctly implemented", func() {
				failures := InterceptGomegaFailures(func() {
					TestConfigSetDefaultValues(&config)
				})
				Expect(failures).NotTo(BeEmpty())
			})
		})

		Describe("TestConfigGetEnumsCount", func() {
			It("should fail when not correctly implemented", func() {
				failures := InterceptGomegaFailures(func() {
					TestConfigGetEnumsCount(&config, 99)
				})
				Expect(failures).NotTo(BeEmpty())
			})
		})

		Describe("TestConfigGetFieldsCount", func() {
			It("should fail when not correctly implemented", func() {
				failures := InterceptGomegaFailures(func() {
					TestConfigGetFieldsCount(&config, 99)
				})
				Expect(failures).NotTo(BeEmpty())
			})
		})
	})

	Describe("Service test helpers", func() {
		var service dummyService
		BeforeEach(func() {
			service = dummyService{}
		})
		Describe("TestConfigSetInvalidQueryValue", func() {
			It("should fail when not correctly implemented", func() {
				failures := InterceptGomegaFailures(func() {
					TestServiceSetInvalidParamValue(&service, "invalid", "value")
				})
				Expect(failures).To(HaveLen(1))
			})
		})
	})
})

type dummyConfig struct {
	standard.EnumlessConfig
	Foo uint64 `key:"foo" default:"-1"`
}

func (dc *dummyConfig) GetURL() *url.URL           { return &url.URL{} }
func (dc *dummyConfig) SetURL(u *url.URL) error    { return nil }
func (dc *dummyConfig) Get(string) (string, error) { return "", nil }
func (dc *dummyConfig) Set(string, string) error   { return nil }
func (dc *dummyConfig) QueryFields() []string      { return []string{} }

type dummyService struct {
	standard.Standard
	Config dummyConfig
}

func (s *dummyService) Initialize(_ *url.URL, _ types.StdLogger) error { return nil }
func (s *dummyService) Send(_ string, _ *types.Params) error           { return nil }
