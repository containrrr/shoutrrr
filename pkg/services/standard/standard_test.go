package standard

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestStandard(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Shoutrrr Standard Suite")
}

var (
	logger       *Logger
	builder      *strings.Builder
	stringLogger *log.Logger
)

var _ = Describe("the standard logging implementation", func() {
	When("setlogger is called with nil", func() {
		It("should provide the logging API without any errors", func() {
			logger = &Logger{}
			logger.SetLogger(nil)
			logger.Log("discarded log message")

			Expect(logger.logger).ToNot(BeNil())
		})

	})
	When("setlogger is called with a proper logger", func() {
		BeforeEach(func() {
			logger = &Logger{}
			builder = &strings.Builder{}
			stringLogger = log.New(builder, "", 0)
		})
		When("when  logger.Log is called", func() {
			It("should log messages", func() {

				logger.SetLogger(stringLogger)
				logger.Log("foo")
				logger.Log("bar")

				Expect(builder.String()).To(Equal("foo\nbar\n"))
			})
		})
		When("when  logger.Logf is called", func() {
			It("should log messages", func() {

				logger.SetLogger(stringLogger)
				logger.Logf("foo %d", 7)

				Expect(builder.String()).To(Equal("foo 7\n"))
			})
		})
	})
})

var _ = Describe("the standard template implementation", func() {
	When("a template is being set from a file", func() {
		It("should load the template without any errors", func() {
			file, err := ioutil.TempFile("", "")
			if err != nil {
				Skip(fmt.Sprintf("Could not create temp file: %s", err))
				return
			}
			fileName := file.Name()
			defer os.Remove(fileName)

			_, err = io.WriteString(file, "template content")
			if err != nil {
				Skip(fmt.Sprintf("Could not write to temp file: %s", err))
				return
			}

			templater := &Templater{}
			err = templater.SetTemplateFile("foo", fileName)
			Expect(err).ToNot(HaveOccurred())
		})

	})
	When("a template is being set from a file that does not exist", func() {
		It("should return an error", func() {
			templater := &Templater{}
			err := templater.SetTemplateFile("foo", "filename_that_should_not_exist")
			Expect(err).To(HaveOccurred())
		})
	})
	When("a template is being set with a badly formatted string", func() {
		It("should return an error", func() {
			templater := &Templater{}
			err := templater.SetTemplateString("foo", "template {{ missing end tag")
			Expect(err).To(HaveOccurred())
		})
	})
	When("a template is being retrieved with a present ID", func() {
		It("should return the corresponding template", func() {
			templater := &Templater{}
			err := templater.SetTemplateString("bar", "template body")
			Expect(err).NotTo(HaveOccurred())

			tpl, found := templater.GetTemplate("bar")
			Expect(tpl).ToNot(BeNil())
			Expect(found).To(BeTrue())
		})
	})
	When("a template is being retrieved with an invalid ID", func() {
		It("should return an error", func() {
			templater := &Templater{}
			err := templater.SetTemplateString("bar", "template body")
			Expect(err).NotTo(HaveOccurred())

			tpl, found := templater.GetTemplate("bad ID")
			Expect(tpl).To(BeNil())
			Expect(found).ToNot(BeTrue())
		})
	})
})

var _ = Describe("the standard enumless config implementation", func() {
	When("it's enum method is called", func() {
		It("should return an empty map", func() {
			Expect((&EnumlessConfig{}).Enums()).To(BeEmpty())
		})
	})
})
