package types

// Sender is the interface needed to implement to send notifications
type Sender interface {
	Send(message string, params *Params) error

	// Rich sender API:
	// SendItems(items []MessageItem, params *Params) error
}
