package permission_factory

import (
	"context"
	"testing"

	audit_domain "sso-bff/internal/domain/audit"
	"sso-bff/internal/httpx"
	permission_sdk "sso-bff/modules/permission/sdk"
)

type stubSessionResolver struct{}

func (stubSessionResolver) SubjectBySessionID(ctx context.Context, sessionID string) (string, error) {
	return "subject", nil
}

type noopAuditWriter struct{}

func (noopAuditWriter) Insert(ctx context.Context, e audit_domain.AuditEvent) error {
	return nil
}

func TestNewPermissionModule_Success(t *testing.T) {
	sdk := permission_sdk.NewPermissionSDK("http://admin", "http://public")
	module, err := NewPermissionModule(sdk, stubSessionResolver{}, noopAuditWriter{}, httpx.CookieConfig{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if module == nil || module.Service == nil || module.Handler == nil {
		t.Fatalf("expected module created")
	}
}

func TestNewPermissionModule_MissingSDK(t *testing.T) {
	_, err := NewPermissionModule(nil, stubSessionResolver{}, noopAuditWriter{}, httpx.CookieConfig{})
	if err == nil {
		t.Fatalf("expected error when sdk nil")
	}
}
