package util

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"runtime"
)

// ConfigureFallbackCertVerification will set the VerifyPeerCertificate callback to a custom function that tries to
// validate the peer certificate using the system certificate store if verifying using the root CAs in the config fails
//
// Workaround for https://github.com/golang/go/issues/16736
// Based on example https://pkg.go.dev/crypto/tls@go1.14#example-Config-VerifyPeerCertificate
func ConfigureFallbackCertVerification(conf *tls.Config) {

	if runtime.GOOS != `windows` {
		// Only needed on windows
		return
	}

	if conf.RootCAs == nil {
		// No custom certs have been added
		return
	}

	// Unless we enable this, the regular verification will still abort the handshake
	conf.InsecureSkipVerify = true

	conf.VerifyPeerCertificate = func(rawCerts [][]byte, _ [][]*x509.Certificate) error {
		certs := make([]*x509.Certificate, len(rawCerts))
		for i, asn1Data := range rawCerts {
			cert, err := x509.ParseCertificate(asn1Data)
			if err != nil {
				return errors.New("tls: failed to parse certificate from server: " + err.Error())
			}
			certs[i] = cert
		}

		opts := x509.VerifyOptions{
			Roots:         conf.RootCAs,
			DNSName:       conf.ServerName,
			Intermediates: x509.NewCertPool(),
		}

		targetCert := certs[0]

		for _, cert := range certs[1:] {
			opts.Intermediates.AddCert(cert)
		}

		_, err := targetCert.Verify(opts)
		if err != nil {
			// Try again using no root store as CryptoAPI will be used to verify instead
			opts.Roots = nil
			_, err = targetCert.Verify(opts)
		}

		return err
	}
}
