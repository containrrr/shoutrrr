package format

import (
	"net/url"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Query Formatter", func() {
	var pkr PropKeyResolver
	BeforeEach(func() {
		ts = &testStruct{}
		pkr = NewPropKeyResolver(ts)
		_ = pkr.SetDefaultProps(ts)
	})
	Describe("Creating a service URL query from a config", func() {
		When("a config property has been changed from default", func() {
			It("should be included in the query string", func() {
				ts.Str = "test"
				query := BuildQuery(&pkr)
				// (pkr, )
				Expect(query).To(Equal("str=test"))
			})
		})
		When("a custom query key conflicts with a config property key", func() {
			It("should include both values, with the custom escaped", func() {
				ts.Str = "service"
				customQuery := url.Values{"str": {"custom"}}
				query := BuildQueryWithCustomFields(&pkr, customQuery)
				Expect(query.Encode()).To(Equal("__str=custom&str=service"))
			})
		})
	})
	Describe("Setting prop values from query", func() {
		When("a custom query key conflicts with a config property key", func() {
			It("should set the config prop from the regular and return the custom one unescaped", func() {
				ts.Str = "service"
				serviceQuery := url.Values{"__str": {"custom"}, "str": {"service"}}
				query, err := SetConfigPropsFromQuery(&pkr, serviceQuery)
				Expect(err).NotTo(HaveOccurred())
				Expect(ts.Str).To(Equal("service"))
				Expect(query.Get("str")).To(Equal("custom"))
			})
		})
	})

})
