package generator_test

import (
	"fmt"
	"github.com/containrrr/shoutrrr/pkg/util/generator"
	"github.com/mattn/go-colorable"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	re "regexp"
	"strings"
	"testing"
)

func TestGenerator(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Generator Suite")
}

var (
	client  *generator.UserDialog
	userOut *gbytes.Buffer
	userIn  *gbytes.Buffer
)

func mockTyped(a ...interface{}) {
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

var _ = Describe("GeneratorCommon", func() {
	BeforeEach(func() {
		userOut = gbytes.NewBuffer()
		userIn = gbytes.NewBuffer()
		userInMono := colorable.NewNonColorable(userIn)
		client = generator.NewUserDialog(userOut, userInMono, map[string]string{"propKey": "propVal"})
	})

	It("reprompt upon invalid answers", func() {
		defer dumpBuffers()
		answer := make(chan string)
		go func() {
			answer <- client.QueryString("name:", generator.Required, "")
		}()

		mockTyped("")
		mockTyped("Normal Human Name")

		Eventually(userIn).Should(gbytes.Say(`name: `))

		Eventually(userIn).Should(gbytes.Say(`field is required`))
		Eventually(userIn).Should(gbytes.Say(`name: `))
		Eventually(answer).Should(Receive(Equal("Normal Human Name")))
	})

	It("should accept any input when validator is nil", func() {
		defer dumpBuffers()
		answer := make(chan string)
		go func() {
			answer <- client.QueryString("name:", nil, "")
		}()
		mockTyped("")
		Eventually(answer).Should(Receive(BeEmpty()))
	})

	It("should use predefined prop value if key is present", func() {
		defer dumpBuffers()
		answer := make(chan string)
		go func() {
			answer <- client.QueryString("name:", generator.Required, "propKey")
		}()
		Eventually(answer).Should(Receive(Equal("propVal")))
	})

	Describe("Query", func() {
		It("should prompt until a valid answer is provided", func() {
			defer dumpBuffers()
			answer := make(chan []string)
			query := "pick foo or bar:"
			go func() {
				answer <- client.Query(query, re.MustCompile("(foo|bar)"), "")
			}()

			mockTyped("")
			mockTyped("foo")

			Eventually(userIn).Should(gbytes.Say(query))
			Eventually(userIn).Should(gbytes.Say(`invalid format`))
			Eventually(userIn).Should(gbytes.Say(query))
			Eventually(answer).Should(Receive(ContainElement("foo")))
		})
	})

	Describe("QueryAll", func() {
		It("should prompt until a valid answer is provided", func() {
			defer dumpBuffers()
			answer := make(chan [][]string)
			query := "pick foo or bar:"
			go func() {
				answer <- client.QueryAll(query, re.MustCompile(`foo(ba[rz])`), "", -1)
			}()

			mockTyped("foobar foobaz")

			Eventually(userIn).Should(gbytes.Say(query))
			var matches [][]string
			Eventually(answer).Should(Receive(&matches))
			Expect(matches).To(ContainElement([]string{"foobar", "bar"}))
			Expect(matches).To(ContainElement([]string{"foobaz", "baz"}))
		})
	})

	Describe("QueryStringPattern", func() {
		It("should prompt until a valid answer is provided", func() {
			defer dumpBuffers()
			answer := make(chan string)
			query := "type of bar:"
			go func() {
				answer <- client.QueryStringPattern(query, re.MustCompile(".*bar"), "")
			}()

			mockTyped("foo")
			mockTyped("foobar")

			Eventually(userIn).Should(gbytes.Say(query))
			Eventually(userIn).Should(gbytes.Say(`invalid format`))
			Eventually(userIn).Should(gbytes.Say(query))
			Eventually(answer).Should(Receive(Equal("foobar")))
		})
	})

	Describe("QueryInt", func() {
		It("should prompt until a valid answer is provided", func() {
			defer dumpBuffers()
			answer := make(chan int64)
			query := "number:"
			go func() {
				answer <- client.QueryInt(query, "", 64)
			}()

			mockTyped("x")
			mockTyped("0x20")

			Eventually(userIn).Should(gbytes.Say(query))
			Eventually(userIn).Should(gbytes.Say(`not a number`))
			Eventually(userIn).Should(gbytes.Say(query))
			Eventually(answer).Should(Receive(Equal(int64(32))))
		})
	})

	Describe("QueryBool", func() {
		It("should prompt until a valid answer is provided", func() {
			defer dumpBuffers()
			answer := make(chan bool)
			query := "cool?"
			go func() {
				answer <- client.QueryBool(query, "")
			}()

			mockTyped("maybe")
			mockTyped("y")

			Eventually(userIn).Should(gbytes.Say(query))
			Eventually(userIn).Should(gbytes.Say(`answer using yes or no`))
			Eventually(userIn).Should(gbytes.Say(query))
			Eventually(answer).Should(Receive(BeTrue()))
		})
	})
})
