package format_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/containrrr/shoutrrr/pkg/format"
)

var _ = Describe("URLPart", func() {
	It("should return the expected URL part for each lookup key", func() {
		Expect(ParseURLPart("user")).To(Equal(URLUser))
		Expect(ParseURLPart("pass")).To(Equal(URLPassword))
		Expect(ParseURLPart("password")).To(Equal(URLPassword))
		Expect(ParseURLPart("host")).To(Equal(URLHost))
		Expect(ParseURLPart("port")).To(Equal(URLPort))

		Expect(ParseURLPart("path")).To(Equal(URLPath))
		Expect(ParseURLPart("path1")).To(Equal(URLPath))
		Expect(ParseURLPart("path2")).To(Equal(URLPath + 1))
		Expect(ParseURLPart("path3")).To(Equal(URLPath + 2))
		Expect(ParseURLPart("path4")).To(Equal(URLPath + 3))

		Expect(ParseURLPart("query")).To(Equal(URLQuery))
		Expect(ParseURLPart("")).To(Equal(URLQuery))
	})
	It("should return the expected suffix for each URL part", func() {
		Expect(URLUser.Suffix()).To(Equal(':'))
		Expect(URLPassword.Suffix()).To(Equal('@'))
		Expect(URLHost.Suffix()).To(Equal(':'))
		Expect(URLPort.Suffix()).To(Equal('/'))
		Expect(URLPath.Suffix()).To(Equal('/'))
	})
})
