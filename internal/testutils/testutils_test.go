package testutils_test

import (
	"net/url"
	"testing"

	"github.com/nicholas-fedor/shoutrrr/internal/testutils"
	"github.com/nicholas-fedor/shoutrrr/pkg/services/standard"
	"github.com/nicholas-fedor/shoutrrr/pkg/types"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

func TestTestUtils(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)

	ginkgo.RunSpecs(t, "Shoutrrr TestUtils Suite")
}

var _ = ginkgo.Describe("the testutils package", func() {
	ginkgo.When("calling function TestLogger", func() {
		ginkgo.It("should not return nil", func() {
			gomega.Expect(testutils.TestLogger()).NotTo(gomega.Equal(nil))
		})
		ginkgo.It(`should have the prefix "[Test] "`, func() {
			gomega.Expect(testutils.TestLogger().Prefix()).To(gomega.Equal("[Test] "))
		})
	})

	ginkgo.Describe("Must helpers", func() {
		ginkgo.Describe("URLMust", func() {
			ginkgo.It("should panic when an invalid URL is passed", func() {
				failures := gomega.InterceptGomegaFailures(func() { testutils.URLMust(":") })
				gomega.Expect(failures).To(gomega.HaveLen(1))
			})
		})

		ginkgo.Describe("JSONRespondMust", func() {
			ginkgo.It("should panic when an invalid struct is passed", func() {
				notAValidJSONSource := func() {}
				failures := gomega.InterceptGomegaFailures(func() { testutils.JSONRespondMust(200, notAValidJSONSource) })
				gomega.Expect(failures).To(gomega.HaveLen(1))
			})
		})
	})

	ginkgo.Describe("Config test helpers", func() {
		var config dummyConfig
		ginkgo.BeforeEach(func() {
			config = dummyConfig{}
		})
		ginkgo.Describe("TestConfigSetInvalidQueryValue", func() {
			ginkgo.It("should fail when not correctly implemented", func() {
				failures := gomega.InterceptGomegaFailures(func() {
					testutils.TestConfigSetInvalidQueryValue(&config, "mock://host?invalid=value")
				})
				gomega.Expect(failures).To(gomega.HaveLen(1))
			})
		})

		ginkgo.Describe("TestConfigGetInvalidQueryValue", func() {
			ginkgo.It("should fail when not correctly implemented", func() {
				failures := gomega.InterceptGomegaFailures(func() {
					testutils.TestConfigGetInvalidQueryValue(&config)
				})
				gomega.Expect(failures).To(gomega.HaveLen(1))
			})
		})

		ginkgo.Describe("TestConfigSetDefaultValues", func() {
			ginkgo.It("should fail when not correctly implemented", func() {
				failures := gomega.InterceptGomegaFailures(func() {
					testutils.TestConfigSetDefaultValues(&config)
				})
				gomega.Expect(failures).NotTo(gomega.BeEmpty())
			})
		})

		ginkgo.Describe("TestConfigGetEnumsCount", func() {
			ginkgo.It("should fail when not correctly implemented", func() {
				failures := gomega.InterceptGomegaFailures(func() {
					testutils.TestConfigGetEnumsCount(&config, 99)
				})
				gomega.Expect(failures).NotTo(gomega.BeEmpty())
			})
		})

		ginkgo.Describe("TestConfigGetFieldsCount", func() {
			ginkgo.It("should fail when not correctly implemented", func() {
				failures := gomega.InterceptGomegaFailures(func() {
					testutils.TestConfigGetFieldsCount(&config, 99)
				})
				gomega.Expect(failures).NotTo(gomega.BeEmpty())
			})
		})
	})

	ginkgo.Describe("Service test helpers", func() {
		var service dummyService
		ginkgo.BeforeEach(func() {
			service = dummyService{}
		})
		ginkgo.Describe("TestConfigSetInvalidQueryValue", func() {
			ginkgo.It("should fail when not correctly implemented", func() {
				failures := gomega.InterceptGomegaFailures(func() {
					testutils.TestServiceSetInvalidParamValue(&service, "invalid", "value")
				})
				gomega.Expect(failures).To(gomega.HaveLen(1))
			})
		})
	})
})

type dummyConfig struct {
	standard.EnumlessConfig
	Foo uint64 `default:"-1" key:"foo"`
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
func (s *dummyService) GetID() string                                  { return "dummy" }
