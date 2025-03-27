package format

import (
	"net/url"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Query Formatter", func() {
	var pkr PropKeyResolver
	ginkgo.BeforeEach(func() {
		ts = &testStruct{}
		pkr = NewPropKeyResolver(ts)
		_ = pkr.SetDefaultProps(ts)
	})
	ginkgo.Describe("Creating a service URL query from a config", func() {
		ginkgo.When("a config property has been changed from default", func() {
			ginkgo.It("should be included in the query string", func() {
				ts.Str = "test"
				query := BuildQuery(&pkr)
				// (pkr, )
				gomega.Expect(query).To(gomega.Equal("str=test"))
			})
		})
		ginkgo.When("a custom query key conflicts with a config property key", func() {
			ginkgo.It("should include both values, with the custom escaped", func() {
				ts.Str = "service"
				customQuery := url.Values{"str": {"custom"}}
				query := BuildQueryWithCustomFields(&pkr, customQuery)
				gomega.Expect(query.Encode()).To(gomega.Equal("__str=custom&str=service"))
			})
		})
	})
	ginkgo.Describe("Setting prop values from query", func() {
		ginkgo.When("a custom query key conflicts with a config property key", func() {
			ginkgo.It("should set the config prop from the regular and return the custom one unescaped", func() {
				ts.Str = "service"
				serviceQuery := url.Values{"__str": {"custom"}, "str": {"service"}}
				query, err := SetConfigPropsFromQuery(&pkr, serviceQuery)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(ts.Str).To(gomega.Equal("service"))
				gomega.Expect(query.Get("str")).To(gomega.Equal("custom"))
			})
		})
	})
})
