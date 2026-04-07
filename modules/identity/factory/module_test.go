package identity_factory

import (
	"context"
	"testing"

	permission_domain "sso-bff/internal/domain/permission"
	"sso-bff/internal/httpx"
	identity_sdk "sso-bff/modules/identity/sdk"
	oauth_gateway "sso-bff/modules/oauth/gateway"

	audit_domain "sso-bff/internal/domain/audit"
)

type noopAuditWriter struct{}

func (noopAuditWriter) Insert(ctx context.Context, e audit_domain.AuditEvent) error {
	return nil
}

type stubTupleLister struct{}

func (stubTupleLister) ListTuples(ctx context.Context, params permission_domain.ListTuplesParams) (permission_domain.ListTuplesResult, error) {
	return permission_domain.ListTuplesResult{}, nil
}

func TestNewModule_Success(t *testing.T) {
	sdk := identity_sdk.NewIdentitySDK("http://public", "http://admin")
	module, err := NewModule(
		sdk,
		oauth_gateway.NewOauthAdminGateway(nil),
		stubTupleLister{},
		noopAuditWriter{},
		httpx.CookieConfig{},
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if module == nil {
		t.Fatalf("expected module not nil")
	}
	if module.Login == nil || module.Admin == nil || module.Settings == nil {
		t.Fatalf("module handlers missing")
	}
}

func TestNewModule_MissingSDK(t *testing.T) {
	_, err := NewModule(
		nil,
		oauth_gateway.NewOauthAdminGateway(nil),
		stubTupleLister{},
		noopAuditWriter{},
		httpx.CookieConfig{},
	)
	if err == nil {
		t.Fatalf("expected error when sdk is nil")
	}
}
