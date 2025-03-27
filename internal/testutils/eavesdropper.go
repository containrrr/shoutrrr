package testutils

// Eavesdropper is an interface that provides a way to get a summarized output of a connection RX and TX.
type Eavesdropper interface {
	GetConversation(includeGreeting bool) string
	GetClientSentences() []string
}
