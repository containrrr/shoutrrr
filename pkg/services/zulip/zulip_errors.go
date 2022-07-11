package zulip

// ErrorMessage for error events within the zulip service
type ErrorMessage string

const (
	// MissingAPIKey from the service URL
	MissingAPIKey ErrorMessage = "botKey missing from config URL"
	// MissingHost from the service URL
	MissingHost ErrorMessage = "host missing from config URL"
	// MissingBotMail from the service URL
	MissingBotMail ErrorMessage = "botMail missing from config URL"
	// TopicTooLong if topic is more than 60 characters
	TopicTooLong ErrorMessage = "topic exceeds max length (%d characters): was %d characters"
)
