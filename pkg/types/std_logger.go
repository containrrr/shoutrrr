package types

// StdLogger is an interface for outputting log information from services that are non-fatal.
type StdLogger interface {
	Print(args ...any)
	Printf(format string, args ...any)
	Println(args ...any)
}
