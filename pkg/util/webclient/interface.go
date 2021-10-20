package webclient

type Client interface {
	Get(url string, response interface{}) error
	Post(url string, request interface{}, response interface{}) error
	Headers() map[string]string
	ErrorResponse(err error, response interface{}) bool
}

type StringClient interface {
	Get(url string, response *string) error
	Post(url string, request string, response *string) error
	Headers() map[string]string
	ErrorResponse(err error, response string) bool
}
