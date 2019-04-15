package slack

type ErrorMessage string
const (
	TokenAMissing   ErrorMessage = "first part of the API token is missing"
	TokenBMissing   ErrorMessage = "second part of the API token is missing"
	TokenCMissing   ErrorMessage = "third part of the API token is missing."
	TokenAMalformed ErrorMessage = "first part of the API token is malformed"
	TokenBMalformed ErrorMessage = "second part of the API token is malformed"
	TokenCMalformed ErrorMessage = "third part of the API token is malformed"
	NotEnoughArguments ErrorMessage = "the url does not include enough arguments"
	CouldNotCreateConfig ErrorMessage = "could not generate a config for slack"
)
