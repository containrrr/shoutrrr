package types

type Generator interface {
	Generate(service Service, props map[string]string, args []string) (ServiceConfig, error)
}
