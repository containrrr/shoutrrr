package types

type ConfigProp interface {
	SetFromProp(propValue string) error
	GetPropValue() (string, error)
}
