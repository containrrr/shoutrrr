package types

// Sender is the interface needed to implement to send notifications
type Sender interface {
	Send(message string, params *Params) error

	SendItems(items []MessageItem, params *Params) error
}
