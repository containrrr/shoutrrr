package types

type Sender interface {
	Send(message string, params *map[string]string) error
}