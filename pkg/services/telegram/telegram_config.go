//go:generate go run ../../../cmd/shoutrrr-gen --lang go ../../../spec/telegram.yml

package telegram

// Scheme is the identifying part of this service's configuration URL
const (
	Scheme = "telegram"
)
