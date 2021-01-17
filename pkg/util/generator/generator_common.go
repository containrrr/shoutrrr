package generator

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	re "regexp"
	"strconv"
)

var invalidFormatError = errors.New("invalid format")

func ValidateFormat(validator func(string) bool) func(string) error {
	return func(answer string) error {
		if validator(answer) {
			return nil
		}
		return invalidFormatError
	}
}

type UserDialog struct {
	reader  io.Reader
	writer  io.Writer
	scanner *bufio.Scanner
	props   map[string]string
}

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

func (ud *UserDialog) Write(message string, v ...interface{}) {
	if _, err := fmt.Fprintf(ud.writer, message, v...); err != nil {
		fmt.Printf("failed to write to output: %v", err)
	}
}

func (ud *UserDialog) Writeln(format string, v ...interface{}) {
	ud.Write(format+"\n", v...)
}

func (ud *UserDialog) Query(prompt string, validator *re.Regexp, key string) (groups []string) {
	ud.QueryString(prompt, ValidateFormat(func(answer string) bool {
		groups = validator.FindStringSubmatch(answer)
		return groups != nil
	}), key)

	return groups
}

func (ud *UserDialog) QueryAll(prompt string, validator *re.Regexp, key string, maxMatches int) (matches [][]string) {
	ud.QueryString(prompt, ValidateFormat(func(answer string) bool {
		matches = validator.FindAllStringSubmatch(answer, maxMatches)
		return matches != nil
	}), key)

	return matches
}

func (ud *UserDialog) QueryString(prompt string, validator func(string) error, key string) string {

	if validator == nil {
		validator = func(string) error {
			return nil
		}
	}

	answer, foundProp := ud.props[key]
	if foundProp {
		err := validator(answer)
		if err == nil {
			ud.Writeln("Using prop value '%v' for '%v", answer, key)
			return answer
		}
		ud.Writeln("Supplied prop value '%v' is not valid for '%v': %v", answer, key, err)
	}

	for {
		ud.Write(prompt)
		if !ud.scanner.Scan() {
			if err := ud.scanner.Err(); err != nil {
				ud.Writeln(err.Error())
				continue
			}

			// Input closed, so let's just return an empty string
			return ""
		}
		answer = ud.scanner.Text()

		if err := validator(answer); err != nil {
			ud.Writeln("%v", err)
			continue
		}
		return answer
	}
}

func (ud *UserDialog) QueryStringPattern(prompt string, validator *re.Regexp, key string) (answer string) {

	if validator == nil {
		panic("validator cannot be nil")
	}

	return ud.QueryString(prompt, func(s string) error {
		if validator.MatchString(s) {
			return nil
		}
		return invalidFormatError
	}, key)
}

func (ud *UserDialog) QueryInt(prompt string, key string, bitSize int) (value int64) {
	validator := re.MustCompile(`^((0x|#)([0-9a-fA-F]+))|(-?[0-9]+)$`)
	ud.QueryString(prompt, func(answer string) error {
		groups := validator.FindStringSubmatch(answer)
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
