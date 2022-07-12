package smtp

// type encMethodVals struct {
// 	// None means no encryption
// 	None encMethod
// 	// ExplicitTLS means that TLS needs to be initiated by using StartTLS
// 	ExplicitTLS encMethod
// 	// ImplicitTLS means that TLS is used for the whole session
// 	ImplicitTLS encMethod
// 	// Auto means that TLS will be implicitly used for port 465, otherwise explicit TLS will be used if its supported
// 	Auto encMethod

// 	// Enum is the EnumFormatter instance for EncMethods
// 	Enum types.EnumFormatter
// }

func useImplicitTLS(encryption encryptionOption, port uint16) bool {
	switch encryption {
	case EncryptionOptions.ImplicitTLS:
		return true
	case EncryptionOptions.Auto:
		return port == ImplicitTLSPort
	default:
		return false
	}
}

// ImplicitTLSPort is de facto standard SMTPS port
const ImplicitTLSPort = 465
