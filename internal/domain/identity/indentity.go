package identity_domain

type Identity struct {
	ID     string
	Traits map[string]any
}

type ListIdentitiesParams struct {
	Page                  int64
	PerPage               int64
	PageSize              int64
	PageToken             string
	CredentialsIdentifier string
}

type ListSessionsParams struct {
	PageSize  int64
	PageToken string
	Active    *bool
}

type IdentitySessionDevice struct {
	IPAddress string `json:"ip_address,omitempty"`
	UserAgent string `json:"user_agent,omitempty"`
}

type IdentitySession struct {
	IdentityID string
	Devices    []IdentitySessionDevice
}
