package pkr

import (
	t "github.com/containrrr/shoutrrr/pkg/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Prop Key Resolver", func() {
	var (
		ts  *testStruct
		pkr PropKeyResolver
	)
	BeforeEach(func() {
		ts = &testStruct{}
		pkr = NewPropKeyResolver(ts)
		_ = pkr.SetDefaultProps(ts)
	})
	Describe("Updating config props from params", func() {
		When("a param matches a prop key", func() {
			It("should be updated in the config", func() {
				err := pkr.UpdateConfigFromParams(nil, &t.Params{"str": "newValue"})
				Expect(err).NotTo(HaveOccurred())
				Expect(ts.Str).To(Equal("newValue"))
			})
		})
		When("a param does not match a prop key", func() {
			It("should report the first error", func() {
				err := pkr.UpdateConfigFromParams(nil, &t.Params{"a": "z"})
				Expect(err).To(HaveOccurred())
			})
			It("should process the other keys", func() {
				_ = pkr.UpdateConfigFromParams(nil, &t.Params{"signed": "1", "b": "c", "str": "val"})
				Expect(ts.Signed).To(Equal(1))
				Expect(ts.Str).To(Equal("val"))
			})
		})
	})
	Describe("Setting default props", func() {
		When("a default tag are set for a field", func() {
			It("should have that value as default", func() {
				Expect(ts.Str).To(Equal("notempty"))
			})
		})
		When("a default tag have an invalid value", func() {
			It("should have that value as default", func() {
				tsb := &testStructBadDefault{}
				pkr = NewPropKeyResolver(tsb)
				err := pkr.SetDefaultProps(tsb)
				Expect(err).To(HaveOccurred())
			})
		})
	})

})
