package types

// Generator is the interface for tools that generate service configurations from a user dialog
type Generator interface {
	Generate(service Service, props map[string]string, args []string) (ServiceConfig, error)
}
