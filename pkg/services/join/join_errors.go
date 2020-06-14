package join

// ErrorMessage for error events within the pushover service
type ErrorMessage string

const (
	// APIKeyMissing should be used when a config URL is missing a token
	APIKeyMissing ErrorMessage = "API key missing from config URL"

	// DevicesMissing should be used when a config URL is missing devices
	DevicesMissing ErrorMessage = "devices missing from config URL"
)
