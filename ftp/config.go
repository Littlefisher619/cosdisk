package ftp

// PortRange defines a port-range
// ... used only for the passive transfer listening range at this stage.
type PortRange struct {
	Start int `json:"start"` // Start of the range
	End   int `json:"end"`   // End of the range
}

// Logging defines how we will log accesses
type Logging struct {
	FtpExchanges bool `json:"ftp_exchanges"` // Log all ftp exchanges
	FileAccesses bool `json:"file_accesses"` // Log all file accesses
}

// TLS define the TLS Config
type TLS struct {
	ServerCert *ServerCert `json:"server_cert"` // Server certificates
}

// ServerCert defines the TLS server certificate config
type ServerCert struct {
	Cert string `json:"cert"` // Public certificate(s)
	Key  string `json:"key"`  // Private key
}

// Config defines the content of the config file
type Config struct {
	Version                  int        `json:"version"`                     // File format version
	ListenAddress            string     `json:"listen_address"`              // Address to listen on
	PublicHost               string     `json:"public_host"`                 // Public host to listen on
	MaxClients               int        `json:"max_clients"`                 // Maximum clients who can connect
	PassiveTransferPortRange *PortRange `json:"passive_transfer_port_range"` // Listen port range
	Logging                  Logging    `json:"logging"`                     // Logging parameters
	TLS                      *TLS       `json:"tls"`                         // TLS Config
}
