package format

import (
	"github.com/nicholas-fedor/shoutrrr/pkg/types"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Prop Key Resolver", func() {
	var (
		ts  *testStruct
		pkr PropKeyResolver
	)
	ginkgo.BeforeEach(func() {
		ts = &testStruct{}
		pkr = NewPropKeyResolver(ts)
		_ = pkr.SetDefaultProps(ts)
	})
	ginkgo.Describe("Updating config props from params", func() {
		ginkgo.When("a param matches a prop key", func() {
			ginkgo.It("should be updated in the config", func() {
				err := pkr.UpdateConfigFromParams(nil, &types.Params{"str": "newValue"})
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(ts.Str).To(gomega.Equal("newValue"))
			})
		})
		ginkgo.When("a param does not match a prop key", func() {
			ginkgo.It("should report the first error", func() {
				err := pkr.UpdateConfigFromParams(nil, &types.Params{"a": "z"})
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
			ginkgo.It("should process the other keys", func() {
				_ = pkr.UpdateConfigFromParams(nil, &types.Params{"signed": "1", "b": "c", "str": "val"})
				gomega.Expect(ts.Signed).To(gomega.Equal(1))
				gomega.Expect(ts.Str).To(gomega.Equal("val"))
			})
		})
	})
	ginkgo.Describe("Setting default props", func() {
		ginkgo.When("a default tag are set for a field", func() {
			ginkgo.It("should have that value as default", func() {
				gomega.Expect(ts.Str).To(gomega.Equal("notempty"))
			})
		})
		ginkgo.When("a default tag have an invalid value", func() {
			ginkgo.It("should have that value as default", func() {
				tsb := &testStructBadDefault{}
				pkr = NewPropKeyResolver(tsb)
				err := pkr.SetDefaultProps(tsb)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})
	})
})
