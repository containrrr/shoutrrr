package urlpart_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/containrrr/shoutrrr/pkg/urlpart"
)

var _ = Describe("URLPart", func() {
	It("should return the expected URL part for each lookup key", func() {
		Expect(ParseOne("user")).To(Equal(User))
		Expect(ParseOne("pass")).To(Equal(Password))
		Expect(ParseOne("password")).To(Equal(Password))
		Expect(ParseOne("host")).To(Equal(Host))
		Expect(ParseOne("port")).To(Equal(Port))

		Expect(ParseOne("path")).To(Equal(Path1))
		Expect(ParseOne("path1")).To(Equal(Path1))
		Expect(ParseOne("path2")).To(Equal(Path2))
		Expect(ParseOne("path3")).To(Equal(Path3))
		Expect(ParseOne("path4")).To(Equal(Path4))

		Expect(ParseOne("query")).To(Equal(Query))
		Expect(ParseOne("")).To(Equal(Query))
	})
	It("should return the expected suffix for each URL part", func() {
		Expect(User.Suffix()).To(Equal(':'))
		Expect(Password.Suffix()).To(Equal('@'))
		Expect(Host.Suffix()).To(Equal(':'))
		Expect(Port.Suffix()).To(Equal('/'))
		Expect(Path1.Suffix()).To(Equal('/'))
	})
})
