package zulip

// ErrorMessage for error events within the zulip service
type ErrorMessage string

const (
	// MissingAPIKey from the service URL
	MissingAPIKey ErrorMessage = "missing API key"
	// MissingHost from the service URL
	MissingHost ErrorMessage = "missing Zulip host"
	// MissingBotMail from the service URL
	MissingBotMail ErrorMessage = "missing Bot mail address"
	// TopicTooLong if topic is more than 60 characters
	TopicTooLong ErrorMessage = "topic exceeds max length (%d characters): was %d characters"
)
