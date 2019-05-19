package slack

// ErrorMessage for error events within the slack plugin
type ErrorMessage string

const (
	// TokenAMissing from the service URL
	TokenAMissing   ErrorMessage = "first part of the API token is missing"
	// TokenBMissing from the service URL
	TokenBMissing   ErrorMessage = "second part of the API token is missing"
	// TokenCMissing from the service URL
	TokenCMissing   ErrorMessage = "third part of the API token is missing."
	// TokenAMalformed inthe service URL
	TokenAMalformed ErrorMessage = "first part of the API token is malformed"
	// TokenBMalformed inthe service URL
	TokenBMalformed ErrorMessage = "second part of the API token is malformed"
	// TokenCMalformed inthe service URL
	TokenCMalformed ErrorMessage = "third part of the API token is malformed"
	// NotEnoughArguments provided in the service URL
	NotEnoughArguments ErrorMessage = "the url does not include enough arguments"
)
