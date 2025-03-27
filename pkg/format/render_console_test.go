package format

import (
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/onsi/gomega/format"
)

var _ = ginkgo.Describe("RenderConsole", func() {
	format.CharactersAroundMismatchToInclude = 30
	renderer := ConsoleTreeRenderer{WithValues: false}

	ginkgo.It("should render the expected output based on config reflection/tags", func() {
		actual := testRenderTree(renderer, &struct {
			Name string `default:"notempty"`
			Host string `url:"host"`
		}{})

		expected := `
Host string                                                                       <URL: Host> <Required>
Name string                                                                       <Default: notempty>
`[1:]

		gomega.Expect(actual).To(gomega.Equal(expected))
	})

	ginkgo.It(`should render enum types as "option"`, func() {
		actual := testRenderTree(renderer, &testEnummer{})

		expected := `
Choice option                                                                       <Default: Maybe> [Yes, No, Maybe]
`[1:]

		gomega.Expect(actual).To(gomega.Equal(expected))
	})

	ginkgo.It("should render url paths in sorted order", func() {
		actual := testRenderTree(renderer, &struct {
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

		gomega.Expect(actual).To(gomega.Equal(expected))
	})
})
