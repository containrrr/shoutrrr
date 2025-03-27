package testutils

import "net/http"

// MockClientService is used to allow mocking the HTTP client when testing.
type MockClientService interface {
	GetHTTPClient() *http.Client
}
