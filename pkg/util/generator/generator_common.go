package generator

import (
	"bufio"
	"errors"
	"fmt"
	f "github.com/containrrr/shoutrrr/pkg/format"
	"github.com/fatih/color"
	"io"
	re "regexp"
	"strconv"
)

var errInvalidFormat = errors.New("invalid format")

// ValidateFormat is a validation wrapper turning false bool results into errors
func ValidateFormat(validator func(string) bool) func(string) error {
	return func(answer string) error {
		if validator(answer) {
			return nil
		}
		return errInvalidFormat
	}
}

var errRequired = errors.New("field is required")

// Required is a validator that checks whether the input contains any characters
func Required(answer string) error {
	if answer == "" {
		return errRequired
	}
	return nil
}

// UserDialog is an abstraction for question/answer based user interaction
type UserDialog struct {
	reader  io.Reader
	writer  io.Writer
	scanner *bufio.Scanner
	props   map[string]string
}

// NewUserDialog initializes a UserDialog with safe defaults
func NewUserDialog(reader io.Reader, writer io.Writer, props map[string]string) *UserDialog {
	if props == nil {
		props = map[string]string{}
	}
	return &UserDialog{
		reader:  reader,
		writer:  writer,
		scanner: bufio.NewScanner(reader),
		props:   props,
	}
}

// Write message to user
func (ud *UserDialog) Write(message string, v ...interface{}) {
	if _, err := fmt.Fprintf(ud.writer, message, v...); err != nil {
		fmt.Printf("failed to write to output: %v", err)
	}
}

// Writeln writes a message to the user that completes a line
func (ud *UserDialog) Writeln(format string, v ...interface{}) {
	ud.Write(format+"\n", v...)
}

// Query writes the prompt to the user and returns the regex groups if it matches the validator pattern
func (ud *UserDialog) Query(prompt string, validator *re.Regexp, key string) (groups []string) {
	ud.QueryString(prompt, ValidateFormat(func(answer string) bool {
		groups = validator.FindStringSubmatch(answer)
		return groups != nil
	}), key)

	return groups
}

// QueryAll is a version of Query that can return multiple matches
func (ud *UserDialog) QueryAll(prompt string, validator *re.Regexp, key string, maxMatches int) (matches [][]string) {
	ud.QueryString(prompt, ValidateFormat(func(answer string) bool {
		matches = validator.FindAllStringSubmatch(answer, maxMatches)
		return matches != nil
	}), key)

	return matches
}

// QueryString writes the prompt to the user and returns the answer if it passes the validator function
func (ud *UserDialog) QueryString(prompt string, validator func(string) error, key string) string {

	if validator == nil {
		validator = func(string) error {
			return nil
		}
	}

	answer, foundProp := ud.props[key]
	if foundProp {
		err := validator(answer)
		colAnswer := f.ColorizeValue(answer, false)
		colKey := f.ColorizeProp(key)
		if err == nil {
			ud.Writeln("Using prop value %v for %v", colAnswer, colKey)
			return answer
		}
		ud.Writeln("Supplied prop value %v is not valid for %v: %v", colAnswer, colKey, err)
	}

	for {
		ud.Write("%v ", prompt)
		color.Set(color.FgHiWhite)
		if !ud.scanner.Scan() {
			if err := ud.scanner.Err(); err != nil {
				ud.Writeln(err.Error())
				continue
			}

			// Input closed, so let's just return an empty string
			return ""
		}
		answer = ud.scanner.Text()
		color.Unset()

		if err := validator(answer); err != nil {
			ud.Writeln("%v", err)
			ud.Writeln("")
			continue
		}
		return answer
	}
}

// QueryStringPattern is a version of QueryString taking a regular expression pattern as the validator
func (ud *UserDialog) QueryStringPattern(prompt string, validator *re.Regexp, key string) (answer string) {

	if validator == nil {
		panic("validator cannot be nil")
	}

	return ud.QueryString(prompt, func(s string) error {
		if validator.MatchString(s) {
			return nil
		}
		return errInvalidFormat
	}, key)
}

// QueryInt writes the prompt to the user and returns the answer if it can be parsed as an integer
func (ud *UserDialog) QueryInt(prompt string, key string, bitSize int) (value int64) {
	validator := re.MustCompile(`^((0x|#)([0-9a-fA-F]+))|(-?[0-9]+)$`)
	ud.QueryString(prompt, func(answer string) error {
		groups := validator.FindStringSubmatch(answer)
		if len(groups) < 1 {
			return errors.New("not a number")
		}
		number := groups[0]
		base := 0
		if groups[2] == "#" {
			// Explicitly treat #ffa080 as hexadecimal
			base = 16
			number = groups[3]
		}

		var err error
		value, err = strconv.ParseInt(number, base, bitSize)

		return err
	}, key)
	return value
}

// QueryBool writes the prompt to the user and returns the answer if it can be parsed as a boolean
func (ud *UserDialog) QueryBool(prompt string, key string) (value bool) {
	ud.QueryString(prompt, func(answer string) error {
		parsed, ok := f.ParseBool(answer, false)
		if ok {
			value = parsed
			return nil
		}
		return fmt.Errorf("answer using %v or %v", f.ColorizeTrue("yes"), f.ColorizeFalse("no"))
	}, key)
	return value
}
