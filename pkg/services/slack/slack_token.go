package slack

import (
	"errors"
	"regexp"
	"strings"
)

// Token is a three part string split into A, B and C
type Token []string

func ValidateToken(token Token) error {
	if err := tokenPartsAreNotEmpty(token); err != nil {
		return err
	} else if err := tokenPartsAreValidFormat(token); err != nil {
		return err
	}
	return nil
}

func tokenPartsAreNotEmpty(token Token) error {
	if token[0] == "" {
		return errors.New(string(TokenAMissing))
	} else if token[1] == "" {
		return errors.New(string(TokenBMissing))
	} else if token[2] == "" {
		return errors.New(string(TokenCMissing))
	}
	return nil
}

func tokenPartsAreValidFormat(token Token) error {
	if !matchesPattern("[A-Z0-9]{9}", token[0]) {
		return errors.New(string(TokenAMalformed))
	} else if !matchesPattern("[A-Z0-9]{9}", token[1]) {
		return errors.New(string(TokenBMalformed))
	} else if !matchesPattern("[A-Za-z0-9]{24}", token[2]) {
		return errors.New(string(TokenCMalformed))
	}
	return nil
}

func matchesPattern(pattern string, part string) bool {
	matched, err := regexp.Match(pattern, []byte(part))
	if matched != true || err != nil {
		return false
	}
	return true
}

func (t Token) String() string {
	return strings.Join(t, "-")
}

// ParseToken creates a Token from a sting representation
func ParseToken(s string) Token {
	token := strings.Split(s, "-")
	return token
}
