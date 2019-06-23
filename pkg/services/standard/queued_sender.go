package standard

import (
	"fmt"
	"strings"

	"github.com/containrrr/shoutrrr/pkg/types"
)

// QueuedSender implements the standard queue sender interface
type QueuedSender struct {
	queue []string
}

// Enqueuef adds a formatted message to an internal queue and sends it when SendQueued is invoked
func (qs *QueuedSender) Enqueuef(format string, v ...interface{}) {
	qs.Enqueue(fmt.Sprintf(format, v...))
}

// Enqueue adds the message to an internal queue and sends it when SendQueued is invoked
func (qs *QueuedSender) Enqueue(message string) {
	qs.queue = append(qs.queue, message)
}

// Flush sends all messages that have been queued up as a combined message. This method should be deferred!
func (qs *QueuedSender) Flush(params *map[string]string) {
	var anonService interface{} = qs
	service := anonService.(types.Service)

	// Since this method is supposed to be deferred we just have to ignore errors
	_ = service.Send(strings.Join(qs.queue, "\n"), params)
}