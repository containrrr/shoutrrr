package types

// ConfigProp interface is used to de-/serialize structs from/to a string representation
type ConfigProp interface {
	SetFromProp(propValue string) error
	GetPropValue() (string, error)
}
