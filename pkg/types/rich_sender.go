package types

// RichSender is the interface needed to implement to send rich notifications
type RichSender interface {
	SendItems(items []MessageItem, params Params) error
}
