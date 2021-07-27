package types

// TLSClient is the interface that needs to be implemented for custom TLS certificate support
type TLSClient interface {
	AddTrustedRootCertificate([]byte) bool
}
