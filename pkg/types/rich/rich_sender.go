package rich

import "github.com/containrrr/shoutrrr/pkg/types"

// Sender is the interface needed to implement to send rich notifications
type Sender interface {
	SendItems(items MessageItem, params *types.Params) error
}
