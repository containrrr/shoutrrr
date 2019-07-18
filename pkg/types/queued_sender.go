package types

// QueuedSender is the interface for a proxied sender that queues messages before sending
type QueuedSender interface {
	Enqueuef(format string, v ...interface{})
	Enqueue(message string)
	Flush(params *map[string]string)
	Service() Service
}
