package audit_domain

type AuditEvent struct {
	IdentityID string
	ClientID   string
	EventType  string
	IPAddress  string
	UserAgent  string
}
