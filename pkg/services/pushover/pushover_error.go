package pushover

// ErrorMessage for error events within the pushover service
type ErrorMessage string

const (
	// UserMissing should be used when a config URL is missing a user
	UserMissing ErrorMessage = "user missing from config URL"
	// TokenMissing should be used when a config URL is missing a token
	TokenMissing ErrorMessage = "token missing from config URL"
)
