package format

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	gformat "github.com/onsi/gomega/format"
)

var _ = Describe("RenderConsole", func() {
	gformat.CharactersAroundMismatchToInclude = 30
	renderer := ConsoleTreeRenderer{WithValues: false}

	It("should render the expected output based on config reflection/tags", func() {
		actual := testRenderTee(renderer, &struct {
			Name string `default:"notempty"`
			Host string `url:"host"`
		}{})

		expected := `
Host string                                                                       <URL: Host> <Required>
Name string                                                                       <Default: notempty>
`[1:]
		println()
		println(actual)

		Expect(actual).To(Equal(expected))
	})

	It("should render url paths in sorted order", func() {
		actual := testRenderTee(renderer, &struct {
			Host  string `url:"host"`
			Path1 string `url:"path1"`
			Path3 string `url:"path3"`
			Path2 string `url:"path2"`
		}{})

		expected := `
Host  string                                                                       <URL: Host> <Required>
Path1 string                                                                       <URL: Path> <Required>
Path2 string                                                                       <URL: Path> <Required>
Path3 string                                                                       <URL: Path> <Required>
`[1:]

		println()
		println(actual)

		Expect(actual).To(Equal(expected))
	})
})

/*

*  __TestEnum__
  Default: `+"`None`"+`
  Possible values: `+"`None`, `Foo`, `Bar`"+`
*/
