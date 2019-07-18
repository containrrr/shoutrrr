package types

// EnumFormatter translate enums between strings and numbers
type EnumFormatter interface {
	Print(e int) string
	Parse(s string) int
	Names() []string
}
