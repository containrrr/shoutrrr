package plugins

type Plugin interface {
    Execute(config, message string) error
}