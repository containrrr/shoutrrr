package smtp

// Config is the configuration needed to send discord notifications
type Config struct {
	Host string
	Username string
	Password string
	Port uint16
	FromAddress string
	FromName string
}
