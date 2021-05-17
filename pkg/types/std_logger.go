package types

// StdLogger is an interface of a subset of the stdlib log.Logger used for
// outputting log information from services that are non-fatal
type StdLogger interface {
	Print(...interface{})
	Printf(string, ...interface{})
	Println(...interface{})
}
