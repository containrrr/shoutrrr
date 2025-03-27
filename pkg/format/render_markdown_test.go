package format

import (
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/onsi/gomega/format"
)

var _ = ginkgo.Describe("RenderMarkdown", func() {
	format.CharactersAroundMismatchToInclude = 10

	ginkgo.It("should render the expected output based on config reflection/tags", func() {
		actual := testRenderTree(MarkdownTreeRenderer{HeaderPrefix: `### `}, &struct {
			Name string `default:"notempty"`
			Host string `url:"host"`
		}{})

		expected := `
### URL Fields

*  __Host__ (**Required**)  
  URL part: <code class="service-url">mock://<strong>host</strong>/</code>  
### Query/Param Props


*  __Name__  
  Default: `[1:] + "`notempty`" + `  

`

		gomega.Expect(actual).To(gomega.Equal(expected))
	})

	ginkgo.It("should render url paths in sorted order", func() {
		actual := testRenderTree(MarkdownTreeRenderer{HeaderPrefix: `### `}, &struct {
			Host  string `url:"host"`
			Path1 string `url:"path1"`
			Path3 string `url:"path3"`
			Path2 string `url:"path2"`
		}{})

		expected := `
### URL Fields

*  __Host__ (**Required**)  
  URL part: <code class="service-url">mock://<strong>host</strong>/path1/path2/path3</code>  
*  __Path1__ (**Required**)  
  URL part: <code class="service-url">mock://host/<strong>path1</strong>/path2/path3</code>  
*  __Path2__ (**Required**)  
  URL part: <code class="service-url">mock://host/path1/<strong>path2</strong>/path3</code>  
*  __Path3__ (**Required**)  
  URL part: <code class="service-url">mock://host/path1/path2/<strong>path3</strong></code>  
### Query/Param Props


`[1:] // Remove initial newline

		gomega.Expect(actual).To(gomega.Equal(expected))
	})

	ginkgo.It("should render prop aliases", func() {
		actual := testRenderTree(MarkdownTreeRenderer{HeaderPrefix: `### `}, &struct {
			Name string `key:"name,handle,title,target"`
		}{})

		expected := `
### URL Fields

### Query/Param Props


*  __Name__ (**Required**)  
  Aliases: `[1:] + "`handle`, `title`, `target`" + `  

`

		gomega.Expect(actual).To(gomega.Equal(expected))
	})

	ginkgo.It("should render possible enum values", func() {
		actual := testRenderTree(MarkdownTreeRenderer{HeaderPrefix: `### `}, &testEnummer{})

		expected := `
### URL Fields

### Query/Param Props


*  __Choice__  
  Default: `[1:] + "`Maybe`" + `  
  Possible values: ` + "`Yes`, `No`, `Maybe`" + `  

`

		gomega.Expect(actual).To(gomega.Equal(expected))
	})

	ginkgo.When("there are no query props", func() {
		ginkgo.It("should prepend an empty-message instead of props description", func() {
			actual := testRenderTree(MarkdownTreeRenderer{
				HeaderPrefix:      `### `,
				PropsDescription:  "Feel free to set these:",
				PropsEmptyMessage: "There is nothing to set!",
			}, &struct {
				Host string `url:"host"`
			}{})

			expected := `
### URL Fields

*  __Host__ (**Required**)  
  URL part: <code class="service-url">mock://<strong>host</strong>/</code>  
### Query/Param Props

There is nothing to set!
`[1:] // Remove initial newline

			//_ = actual == expected
			//println()
			//println(actual)
			gomega.Expect(actual).To(gomega.Equal(expected))
		})
	})

	ginkgo.When("there are query props", func() {
		ginkgo.It("should prepend the props description", func() {
			actual := testRenderTree(MarkdownTreeRenderer{
				HeaderPrefix:      `### `,
				PropsDescription:  "Feel free to set these:",
				PropsEmptyMessage: "There is nothing to set!",
			}, &struct {
				Host     string `url:"host"`
				CoolMode bool   `key:"coolmode" optional:""`
			}{})

			expected := `
### URL Fields

*  __Host__ (**Required**)  
  URL part: <code class="service-url">mock://<strong>host</strong>/</code>  
### Query/Param Props

Feel free to set these:
*  __CoolMode__  
  Default: *empty*  

`[1:] // Remove initial newline

			gomega.Expect(actual).To(gomega.Equal(expected))
		})
	})
})
