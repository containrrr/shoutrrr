package queue

import (
	"fmt"
	"strings"

	"github.com/containrrr/shoutrrr/pkg/types"
)

// QueuedSender implements the standard queue sender interface
type queuedSender struct {
	queue []string
	sender types.Service
}

// Enqueuef adds a formatted message to an internal queue and sends it when SendQueued is invoked
func (qs *queuedSender) Enqueuef(format string, v ...interface{}) {
	qs.Enqueue(fmt.Sprintf(format, v...))
}

// Enqueue adds the message to an internal queue and sends it when SendQueued is invoked
func (qs *queuedSender) Enqueue(message string) {
	qs.queue = append(qs.queue, message)
}

func GetQueued(sender types.Service) types.QueuedSender {
	qs := &queuedSender{
		sender: sender,
	}
	return qs
}

// Flush sends all messages that have been queued up as a combined message. This method should be deferred!
func (qs *queuedSender) Flush(params *map[string]string) {
	var anonService interface{} = qs
	service, ok := anonService.(types.Service)
	if ok {
		// Since this method is supposed to be deferred we just have to ignore errors
		_ = service.Send(strings.Join(qs.queue, "\n"), params)
	}
}

func (qs *queuedSender) Service() types.Service {
	return qs.sender
}