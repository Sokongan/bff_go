package audit

import (
	"context"
	audit_domain "sso-bff/internal/domain/audit"
)

type AuditRepository interface {
	Insert(ctx context.Context, e audit_domain.AuditEvent) error
	ListRecent(
		ctx context.Context,
		identityID string,
		limit int32) ([]audit_domain.AuditEvent, error)
}
