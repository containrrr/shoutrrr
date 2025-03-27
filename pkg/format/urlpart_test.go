package format_test

import (
	"github.com/nicholas-fedor/shoutrrr/pkg/format"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("URLPart", func() {
	ginkgo.It("should return the expected URL part for each lookup key", func() {
		gomega.Expect(format.ParseURLPart("user")).To(gomega.Equal(format.URLUser))
		gomega.Expect(format.ParseURLPart("pass")).To(gomega.Equal(format.URLPassword))
		gomega.Expect(format.ParseURLPart("password")).To(gomega.Equal(format.URLPassword))
		gomega.Expect(format.ParseURLPart("host")).To(gomega.Equal(format.URLHost))
		gomega.Expect(format.ParseURLPart("port")).To(gomega.Equal(format.URLPort))
		gomega.Expect(format.ParseURLPart("path")).To(gomega.Equal(format.URLPath))
		gomega.Expect(format.ParseURLPart("path1")).To(gomega.Equal(format.URLPath))
		gomega.Expect(format.ParseURLPart("path2")).To(gomega.Equal(format.URLPath + 1))
		gomega.Expect(format.ParseURLPart("path3")).To(gomega.Equal(format.URLPath + 2))
		gomega.Expect(format.ParseURLPart("path4")).To(gomega.Equal(format.URLPath + 3))
		gomega.Expect(format.ParseURLPart("query")).To(gomega.Equal(format.URLQuery))
		gomega.Expect(format.ParseURLPart("")).To(gomega.Equal(format.URLQuery))
	})
	ginkgo.It("should return the expected suffix for each URL part", func() {
		gomega.Expect(format.URLUser.Suffix()).To(gomega.Equal(':'))
		gomega.Expect(format.URLPassword.Suffix()).To(gomega.Equal('@'))
		gomega.Expect(format.URLHost.Suffix()).To(gomega.Equal(':'))
		gomega.Expect(format.URLPort.Suffix()).To(gomega.Equal('/'))
		gomega.Expect(format.URLPath.Suffix()).To(gomega.Equal('/'))
	})
})
