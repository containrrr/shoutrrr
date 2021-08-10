package webclient

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"github.com/containrrr/shoutrrr/pkg/types"
	"net/http"
)

// ParserFunc are functions that deserialize a struct from the passed bytes
type ParserFunc func(raw []byte, v interface{}) error

// WriterFunc are functions that serialize the passed struct into a byte stream
type WriterFunc func(v interface{}) ([]byte, error)

var _ types.TLSClient = &ClientService{}
var _ types.HTTPService = &ClientService{}

// ClientService is a Composable that adds a generic web request client to the service
type ClientService struct {
	client   *client
	certPool *x509.CertPool
}

// HTTPClient returns the underlying http.WebClient used in the Service
func (s *ClientService) HTTPClient() *http.Client {
	s.Initialize()
	return s.client.HTTPClient()
}

// WebClient returns the WebClient instance, initializing it if necessary
func (s *ClientService) WebClient() WebClient {
	s.Initialize()
	return s.client
}

// Initialize sets up the WebClient in the default state using JSON serialization and headers
func (s *ClientService) Initialize() {
	if s.client != nil {
		return
	}

	s.client = &client{
		httpClient: http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{},
			},
		},
		headers: http.Header{
			"Content-Type": []string{JSONContentType},
		},
		parse: json.Unmarshal,
		write: func(v interface{}) ([]byte, error) {
			return json.MarshalIndent(v, "", s.client.indent)
		},
	}
}

// AddTrustedRootCertificate adds the specified PEM certificate to the pool of trusted root CAs
func (s *ClientService) AddTrustedRootCertificate(caPEM []byte) bool {
	s.Initialize()
	if s.certPool == nil {
		certPool, err := x509.SystemCertPool()
		if err != nil {
			certPool = x509.NewCertPool()
		}
		s.certPool = certPool
		if tp, ok := s.client.httpClient.Transport.(*http.Transport); ok {
			tp.TLSClientConfig.RootCAs = s.certPool
		}
	}

	return s.certPool.AppendCertsFromPEM(caPEM)
}
