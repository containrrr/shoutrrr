package slack

import "errors"

var (
	// ErrorInvalidToken is returned whenever the specified token does not match any known formats
	ErrorInvalidToken = errors.New("invalid slack token format")

	// ErrorMismatchedTokenSeparators is returned if the token uses different separators between parts (of the recognized `/-,`)
	ErrorMismatchedTokenSeparators = errors.New("invalid webhook token format")
)
