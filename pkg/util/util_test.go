package util_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"

	"github.com/nicholas-fedor/shoutrrr/internal/meta"
	"github.com/nicholas-fedor/shoutrrr/pkg/util"
)

func TestUtil(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Shoutrrr Util Suite")
}

const (
	a = 10
	b = 20
)

var _ = ginkgo.Describe("the util package", func() {
	ginkgo.When("calling function Min", func() {
		ginkgo.It("should return the smallest of two integers", func() {
			gomega.Expect(util.Min(a, b)).To(gomega.Equal(a))
			gomega.Expect(util.Min(b, a)).To(gomega.Equal(a))
		})
	})

	ginkgo.When("calling function Max", func() {
		ginkgo.It("should return the largest of two integers", func() {
			gomega.Expect(util.Max(a, b)).To(gomega.Equal(b))
			gomega.Expect(util.Max(b, a)).To(gomega.Equal(b))
		})
	})

	ginkgo.When("checking if a supplied kind is of the signed integer kind", func() {
		ginkgo.It("should be true if the kind is Int", func() {
			gomega.Expect(util.IsSignedInt(reflect.Int)).To(gomega.BeTrue())
		})
		ginkgo.It("should be false if the kind is String", func() {
			gomega.Expect(util.IsSignedInt(reflect.String)).To(gomega.BeFalse())
		})
	})

	ginkgo.When("checking if a supplied kind is of the unsigned integer kind", func() {
		ginkgo.It("should be true if the kind is Uint", func() {
			gomega.Expect(util.IsUnsignedInt(reflect.Uint)).To(gomega.BeTrue())
		})
		ginkgo.It("should be false if the kind is Int", func() {
			gomega.Expect(util.IsUnsignedInt(reflect.Int)).To(gomega.BeFalse())
		})
	})

	ginkgo.When("checking if a supplied kind is of the collection kind", func() {
		ginkgo.It("should be true if the kind is slice", func() {
			gomega.Expect(util.IsCollection(reflect.Slice)).To(gomega.BeTrue())
		})
		ginkgo.It("should be false if the kind is map", func() {
			gomega.Expect(util.IsCollection(reflect.Map)).To(gomega.BeFalse())
		})
	})

	ginkgo.When("calling function StripNumberPrefix", func() {
		ginkgo.It("should return the default base if none is found", func() {
			_, base := util.StripNumberPrefix("46")
			gomega.Expect(base).To(gomega.Equal(0))
		})
		ginkgo.It("should remove # prefix and return base 16 if found", func() {
			number, base := util.StripNumberPrefix("#ab")
			gomega.Expect(number).To(gomega.Equal("ab"))
			gomega.Expect(base).To(gomega.Equal(16))
		})
	})

	ginkgo.When("checking if a supplied kind is numeric", func() {
		ginkgo.It("should be true if supplied a constant integer", func() {
			gomega.Expect(util.IsNumeric(reflect.TypeOf(5).Kind())).To(gomega.BeTrue())
		})
		ginkgo.It("should be true if supplied a constant float", func() {
			gomega.Expect(util.IsNumeric(reflect.TypeOf(2.5).Kind())).To(gomega.BeTrue())
		})
		ginkgo.It("should be false if supplied a constant string", func() {
			gomega.Expect(util.IsNumeric(reflect.TypeOf("3").Kind())).To(gomega.BeFalse())
		})
	})

	ginkgo.When("calling function DocsURL", func() {
		ginkgo.It("should return the expected URL", func() {
			expectedBase := fmt.Sprintf(`https://nicholas-fedor.github.io/shoutrrr/%s/`, meta.DocsVersion)
			gomega.Expect(util.DocsURL(``)).To(gomega.Equal(expectedBase))
			gomega.Expect(util.DocsURL(`services/logger`)).To(gomega.Equal(expectedBase + `services/logger`))
		})
		ginkgo.It("should strip the leading slash from the path", func() {
			gomega.Expect(util.DocsURL(`/foo`)).To(gomega.Equal(util.DocsURL(`foo`)))
		})
	})
})
