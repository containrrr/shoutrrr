package standard

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/nicholas-fedor/shoutrrr/internal/failures"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

func TestStandard(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Shoutrrr Standard Suite")
}

var (
	logger       *Logger
	builder      *strings.Builder
	stringLogger *log.Logger
)

var _ = ginkgo.Describe("the standard logging implementation", func() {
	ginkgo.When("setlogger is called with nil", func() {
		ginkgo.It("should provide the logging API without any errors", func() {
			logger = &Logger{}
			logger.SetLogger(nil)
			logger.Log("discarded log message")

			gomega.Expect(logger.logger).ToNot(gomega.BeNil())
		})
	})
	ginkgo.When("setlogger is called with a proper logger", func() {
		ginkgo.BeforeEach(func() {
			logger = &Logger{}
			builder = &strings.Builder{}
			stringLogger = log.New(builder, "", 0)
		})
		ginkgo.When("when  logger.Log is called", func() {
			ginkgo.It("should log messages", func() {
				logger.SetLogger(stringLogger)
				logger.Log("foo")
				logger.Log("bar")

				gomega.Expect(builder.String()).To(gomega.Equal("foo\nbar\n"))
			})
		})
		ginkgo.When("when  logger.Logf is called", func() {
			ginkgo.It("should log messages", func() {
				logger.SetLogger(stringLogger)
				logger.Logf("foo %d", 7)

				gomega.Expect(builder.String()).To(gomega.Equal("foo 7\n"))
			})
		})
	})
})

var _ = ginkgo.Describe("the standard template implementation", func() {
	ginkgo.When("a template is being set from a file", func() {
		ginkgo.It("should load the template without any errors", func() {
			file, err := os.CreateTemp("", "")
			if err != nil {
				ginkgo.Skip(fmt.Sprintf("Could not create temp file: %s", err))

				return
			}
			fileName := file.Name()
			defer os.Remove(fileName)

			_, err = io.WriteString(file, "template content")
			if err != nil {
				ginkgo.Skip(fmt.Sprintf("Could not write to temp file: %s", err))

				return
			}

			templater := &Templater{}
			err = templater.SetTemplateFile("foo", fileName)
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
		})
	})
	ginkgo.When("a template is being set from a file that does not exist", func() {
		ginkgo.It("should return an error", func() {
			templater := &Templater{}
			err := templater.SetTemplateFile("foo", "filename_that_should_not_exist")
			gomega.Expect(err).To(gomega.HaveOccurred())
		})
	})
	ginkgo.When("a template is being set with a badly formatted string", func() {
		ginkgo.It("should return an error", func() {
			templater := &Templater{}
			err := templater.SetTemplateString("foo", "template {{ missing end tag")
			gomega.Expect(err).To(gomega.HaveOccurred())
		})
	})
	ginkgo.When("a template is being retrieved with a present ID", func() {
		ginkgo.It("should return the corresponding template", func() {
			templater := &Templater{}
			err := templater.SetTemplateString("bar", "template body")
			gomega.Expect(err).NotTo(gomega.HaveOccurred())

			tpl, found := templater.GetTemplate("bar")
			gomega.Expect(tpl).ToNot(gomega.BeNil())
			gomega.Expect(found).To(gomega.BeTrue())
		})
	})
	ginkgo.When("a template is being retrieved with an invalid ID", func() {
		ginkgo.It("should return an error", func() {
			templater := &Templater{}
			err := templater.SetTemplateString("bar", "template body")
			gomega.Expect(err).NotTo(gomega.HaveOccurred())

			tpl, found := templater.GetTemplate("bad ID")
			gomega.Expect(tpl).To(gomega.BeNil())
			gomega.Expect(found).ToNot(gomega.BeTrue())
		})
	})
})

var _ = ginkgo.Describe("the standard enumless config implementation", func() {
	ginkgo.When("it's enum method is called", func() {
		ginkgo.It("should return an empty map", func() {
			gomega.Expect((&EnumlessConfig{}).Enums()).To(gomega.BeEmpty())
		})
	})
})

var _ = ginkgo.Describe("the standard failure implementation", func() {
	ginkgo.Describe("Failure function", func() {
		ginkgo.When("called with FailParseURL", func() {
			ginkgo.It("should return a failure with the correct message", func() {
				err := errors.New("invalid URL")
				failure := Failure(FailParseURL, err)
				gomega.Expect(failure.ID()).To(gomega.Equal(FailParseURL))
				gomega.Expect(failure.Error()).To(gomega.ContainSubstring("error parsing Service URL"))
				gomega.Expect(failure.Error()).To(gomega.ContainSubstring("invalid URL"))
			})
		})
		ginkgo.When("called with FailUnknown", func() {
			ginkgo.It("should return a failure with the unknown error message", func() {
				err := errors.New("something went wrong")
				failure := Failure(FailUnknown, err)
				gomega.Expect(failure.ID()).To(gomega.Equal(FailUnknown))
				gomega.Expect(failure.Error()).To(gomega.ContainSubstring("an unknown error occurred"))
				gomega.Expect(failure.Error()).To(gomega.ContainSubstring("something went wrong"))
			})
		})
		ginkgo.When("called with an unrecognized FailureID", func() {
			ginkgo.It("should fallback to the unknown error message", func() {
				err := errors.New("unrecognized error")
				failure := Failure(failures.FailureID(999), err) // Arbitrary unknown ID
				gomega.Expect(failure.ID()).To(gomega.Equal(failures.FailureID(999)))
				gomega.Expect(failure.Error()).To(gomega.ContainSubstring("an unknown error occurred"))
				gomega.Expect(failure.Error()).To(gomega.ContainSubstring("unrecognized error"))
			})
		})
		ginkgo.When("called with additional arguments", func() {
			ginkgo.It("should include formatted arguments in the error", func() {
				err := errors.New("base error")
				failure := Failure(FailParseURL, err, "extra info: %s", "details")
				gomega.Expect(failure.Error()).To(gomega.ContainSubstring("error parsing Service URL extra info: details"))
				gomega.Expect(failure.Error()).To(gomega.ContainSubstring("base error"))
			})
		})
	})

	ginkgo.Describe("IsTestSetupFailure function", func() {
		ginkgo.When("called with a FailTestSetup failure", func() {
			ginkgo.It("should return true and the correct message", func() {
				err := errors.New("setup issue")
				failure := Failure(FailTestSetup, err)
				msg, isSetupFailure := IsTestSetupFailure(failure)
				gomega.Expect(isSetupFailure).To(gomega.BeTrue())
				gomega.Expect(msg).To(gomega.ContainSubstring("test setup failed: setup issue"))
			})
		})
		ginkgo.When("called with a different failure", func() {
			ginkgo.It("should return false and an empty message", func() {
				err := errors.New("parse issue")
				failure := Failure(FailParseURL, err)
				msg, isSetupFailure := IsTestSetupFailure(failure)
				gomega.Expect(isSetupFailure).To(gomega.BeFalse())
				gomega.Expect(msg).To(gomega.BeEmpty())
			})
		})
		ginkgo.When("called with nil", func() {
			ginkgo.It("should return false and an empty message", func() {
				msg, isSetupFailure := IsTestSetupFailure(nil)
				gomega.Expect(isSetupFailure).To(gomega.BeFalse())
				gomega.Expect(msg).To(gomega.BeEmpty())
			})
		})
	})
})
