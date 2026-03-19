package identity_factory

import (
	"context"
	"testing"

	"sso-bff/internal/httpx"
	identity_sdk "sso-bff/modules/identity/sdk"
	"sso-bff/modules/oauth"
	oauth_gateway "sso-bff/modules/oauth/gateway"

	audit_domain "sso-bff/internal/domain/audit"
	permission_domain "sso-bff/internal/domain/permission"
)

type stubSessionResolver struct{}

func (stubSessionResolver) SubjectBySessionID(ctx context.Context, sessionID string) (string, error) {
	return "sub", nil
}

func (stubSessionResolver) GetSession(ctx context.Context, sessionID string) (oauth.Session, error) {
	return oauth.Session{}, nil
}

type stubPermissionChecker struct{}

func (stubPermissionChecker) CheckTuple(ctx context.Context, tuple permission_domain.RelationTuple) (bool, error) {
	return true, nil
}

func (stubPermissionChecker) ListTuples(ctx context.Context, params permission_domain.ListTuplesParams) (permission_domain.ListTuplesResult, error) {
	return permission_domain.ListTuplesResult{}, nil
}

type noopAuditWriter struct{}

func (noopAuditWriter) Insert(ctx context.Context, e audit_domain.AuditEvent) error {
	return nil
}

func TestNewModule_Success(t *testing.T) {
	sdk := identity_sdk.NewIdentitySDK("http://public", "http://admin")
	module, err := NewModule(
		sdk,
		oauth_gateway.NewOauthAdminGateway(nil),
		stubSessionResolver{},
		stubPermissionChecker{},
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
		stubSessionResolver{},
		stubPermissionChecker{},
		noopAuditWriter{},
		httpx.CookieConfig{},
	)
	if err == nil {
		t.Fatalf("expected error when sdk is nil")
	}
}
