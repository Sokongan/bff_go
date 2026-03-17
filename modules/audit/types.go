package audit

import (
	"context"
	audit_domain "sso-bff/internal/domain/audit"
)

type AuditWriter interface {
	Insert(ctx context.Context, e audit_domain.AuditEvent) error
}
