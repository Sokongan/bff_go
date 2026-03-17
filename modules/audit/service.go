package audit

import (
	"context"
	"errors"
	audit_domain "sso-bff/internal/domain/audit"
	"sso-bff/modules/oauth"
)

var ErrAuditMisconfigured = errors.New("audit service misconfigured")

type AuditService struct {
	repo AuditRepository
}

func NewAuditService(repo AuditRepository) *AuditService {
	return &AuditService{repo: repo}
}

func (s *AuditService) Insert(
	ctx context.Context,
	e audit_domain.AuditEvent) error {
	if s == nil || s.repo == nil {
		return oauth.ErrServiceMisconfigured
	}
	return s.repo.Insert(ctx, e)
}

func (s *AuditService) ListRecent(
	ctx context.Context,
	identityID string,
	limit int32,
) ([]audit_domain.AuditEvent, error) {
	if s == nil || s.repo == nil {
		return nil, ErrAuditMisconfigured
	}
	return s.repo.ListRecent(ctx, identityID, limit)
}
