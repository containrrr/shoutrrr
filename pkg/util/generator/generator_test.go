package generator_test

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/mattn/go-colorable"
	"github.com/nicholas-fedor/shoutrrr/pkg/util/generator"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

func TestGenerator(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Generator Suite")
}

var (
	client  *generator.UserDialog
	userOut *gbytes.Buffer
	userIn  *gbytes.Buffer
)

func mockTyped(a ...any) {
	_, _ = fmt.Fprint(userOut, a...)
	_, _ = fmt.Fprint(userOut, "\n")
}

func dumpBuffers() {
	for _, line := range strings.Split(string(userIn.Contents()), "\n") {
		println(">", line)
	}

	for _, line := range strings.Split(string(userOut.Contents()), "\n") {
		println("<", line)
	}
}

var _ = ginkgo.Describe("GeneratorCommon", func() {
	ginkgo.BeforeEach(func() {
		userOut = gbytes.NewBuffer()
		userIn = gbytes.NewBuffer()
		userInMono := colorable.NewNonColorable(userIn)
		client = generator.NewUserDialog(userOut, userInMono, map[string]string{"propKey": "propVal"})
	})

	ginkgo.It("reprompt upon invalid answers", func() {
		defer dumpBuffers()
		answer := make(chan string)
		go func() {
			answer <- client.QueryString("name:", generator.Required, "")
		}()

		mockTyped("")
		mockTyped("Normal Human Name")

		gomega.Eventually(userIn).Should(gbytes.Say(`name: `))

		gomega.Eventually(userIn).Should(gbytes.Say(`field is required`))
		gomega.Eventually(userIn).Should(gbytes.Say(`name: `))
		gomega.Eventually(answer).Should(gomega.Receive(gomega.Equal("Normal Human Name")))
	})

	ginkgo.It("should accept any input when validator is nil", func() {
		defer dumpBuffers()
		answer := make(chan string)
		go func() {
			answer <- client.QueryString("name:", nil, "")
		}()
		mockTyped("")
		gomega.Eventually(answer).Should(gomega.Receive(gomega.BeEmpty()))
	})

	ginkgo.It("should use predefined prop value if key is present", func() {
		defer dumpBuffers()
		answer := make(chan string)
		go func() {
			answer <- client.QueryString("name:", generator.Required, "propKey")
		}()
		gomega.Eventually(answer).Should(gomega.Receive(gomega.Equal("propVal")))
	})

	ginkgo.Describe("Query", func() {
		ginkgo.It("should prompt until a valid answer is provided", func() {
			defer dumpBuffers()
			answer := make(chan []string)
			query := "pick foo or bar:"
			go func() {
				answer <- client.Query(query, regexp.MustCompile("(foo|bar)"), "")
			}()

			mockTyped("")
			mockTyped("foo")

			gomega.Eventually(userIn).Should(gbytes.Say(query))
			gomega.Eventually(userIn).Should(gbytes.Say(`invalid format`))
			gomega.Eventually(userIn).Should(gbytes.Say(query))
			gomega.Eventually(answer).Should(gomega.Receive(gomega.ContainElement("foo")))
		})
	})

	ginkgo.Describe("QueryAll", func() {
		ginkgo.It("should prompt until a valid answer is provided", func() {
			defer dumpBuffers()
			answer := make(chan [][]string)
			query := "pick foo or bar:"
			go func() {
				answer <- client.QueryAll(query, regexp.MustCompile(`foo(ba[rz])`), "", -1)
			}()

			mockTyped("foobar foobaz")

			gomega.Eventually(userIn).Should(gbytes.Say(query))
			var matches [][]string
			gomega.Eventually(answer).Should(gomega.Receive(&matches))
			gomega.Expect(matches).To(gomega.ContainElement([]string{"foobar", "bar"}))
			gomega.Expect(matches).To(gomega.ContainElement([]string{"foobaz", "baz"}))
		})
	})

	ginkgo.Describe("QueryStringPattern", func() {
		ginkgo.It("should prompt until a valid answer is provided", func() {
			defer dumpBuffers()
			answer := make(chan string)
			query := "type of bar:"
			go func() {
				answer <- client.QueryStringPattern(query, regexp.MustCompile(".*bar"), "")
			}()

			mockTyped("foo")
			mockTyped("foobar")

			gomega.Eventually(userIn).Should(gbytes.Say(query))
			gomega.Eventually(userIn).Should(gbytes.Say(`invalid format`))
			gomega.Eventually(userIn).Should(gbytes.Say(query))
			gomega.Eventually(answer).Should(gomega.Receive(gomega.Equal("foobar")))
		})
	})

	ginkgo.Describe("QueryInt", func() {
		ginkgo.It("should prompt until a valid answer is provided", func() {
			defer dumpBuffers()
			answer := make(chan int64)
			query := "number:"
			go func() {
				answer <- client.QueryInt(query, "", 64)
			}()

			mockTyped("x")
			mockTyped("0x20")

			gomega.Eventually(userIn).Should(gbytes.Say(query))
			gomega.Eventually(userIn).Should(gbytes.Say(`not a number`))
			gomega.Eventually(userIn).Should(gbytes.Say(query))
			gomega.Eventually(answer).Should(gomega.Receive(gomega.Equal(int64(32))))
		})
	})

	ginkgo.Describe("QueryBool", func() {
		ginkgo.It("should prompt until a valid answer is provided", func() {
			defer dumpBuffers()
			answer := make(chan bool)
			query := "cool?"
			go func() {
				answer <- client.QueryBool(query, "")
			}()

			mockTyped("maybe")
			mockTyped("y")

			gomega.Eventually(userIn).Should(gbytes.Say(query))
			gomega.Eventually(userIn).Should(gbytes.Say(`answer using yes or no`))
			gomega.Eventually(userIn).Should(gbytes.Say(query))
			gomega.Eventually(answer).Should(gomega.Receive(gomega.BeTrue()))
		})
	})
})
