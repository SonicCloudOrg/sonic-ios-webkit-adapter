package WebKitProtocol

type Connection struct {
	Protocol *string `json:"protocol,omitempty"`
	Cipher   *string `json:"cipher,omitempty"`
}

type Certificate struct {
	Subject     *string   `json:"subject,omitempty"`
	ValidFrom   *Walltime `json:"validFrom,omitempty"`
	ValidUntil  *Walltime `json:"validUntil,omitempty"`
	DnsNames    *[]string `json:"dnsNames,omitempty"`
	IpAddresses *[]string `json:"ipAddresses,omitempty"`
}

type Security struct {
	Connection  *Connection  `json:"connection,omitempty"`
	Certificate *Certificate `json:"certificate,omitempty"`
}
