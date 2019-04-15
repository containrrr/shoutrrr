package slack

import (
    "errors"
    "regexp"
)

// Token is a three part string split into A, B and C
type Token struct {
    A string
    B string
    C string
}

func validateToken(token Token) error {
    if err := tokenPartsAreNotEmpty(token); err != nil {
        return err
    } else if err := tokenPartsAreValidFormat(token); err != nil {
        return err
    }
    return nil
}

func tokenPartsAreNotEmpty(token Token) error {
    if token.A == "" {
        return errors.New(string(TokenAMissing))
    } else if token.B == "" {
        return errors.New(string(TokenBMissing))
    } else if token.C == "" {
        return errors.New(string(TokenCMissing))
    }
    return nil
}

func tokenPartsAreValidFormat(token Token) error {
    if !matchesPattern("[A-Z0-9]{9}", token.A) {
        return errors.New(string(TokenAMalformed))
    } else if !matchesPattern("[A-Z0-9]{9}", token.B) {
        return errors.New(string(TokenBMalformed))
    } else if !matchesPattern("[A-Za-z0-9]{24}", token.C) {
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
