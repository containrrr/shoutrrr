package types

// TLSClient is the interface that needs to be implemented for custom TLS certificate support
type TLSClient interface {
	AddTrustedCertificate([]byte) bool
	// UseCustomRootCAs sets whether the HTTP client uses custom loaded certificates instead of the system ones
	// Note that on windows, enabling this will disable the system root CAs (for this service).
	// Because of this, custom root CAs are disabled on windows by default.
	UseCustomRootCAs(enabled bool)
}
