package types

type QueuedSender interface {
	Enqueuef(format string, v ...interface{})
	Enqueue(message string)
	Flush(params *map[string]string)
	Service() Service
}