package pushover

// ErrorMessage for error events within the pushover service
type ErrorMessage string

const (
	UserMissing ErrorMessage = "user missing from config URL"
	TokenMissing ErrorMessage = "token missing from config URL"
)