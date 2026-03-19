package permission_factory_modules

import (
	"context"

	permission_domain "sso-bff/internal/domain/permission"
	"sso-bff/modules/permission"
)

// PermissionService contains business logic around permission tuples.
// It delegates to a PermissionClient (e.g., Keto gateway) provided at construction time.
type PermissionService struct {
	perm permission.PermissionClient
}

func NewPermissionService(perm permission.PermissionClient) *PermissionService {
	return &PermissionService{perm: perm}
}

func (s *PermissionService) WriteTuple(
	ctx context.Context,
	tuple permission_domain.RelationTuple,
) error {
	if s == nil || s.perm == nil {
		return permission.ErrMisconfig
	}
	return s.perm.WriteTuple(ctx, tuple)
}

func (s *PermissionService) CheckTuple(
	ctx context.Context,
	tuple permission_domain.RelationTuple,
) (bool, error) {
	if s == nil || s.perm == nil {
		return false, permission.ErrMisconfig
	}
	return s.perm.CheckTuple(ctx, tuple)
}

func (s *PermissionService) ListTuples(
	ctx context.Context,
	params permission_domain.ListTuplesParams,
) (permission_domain.ListTuplesResult, error) {
	if s == nil || s.perm == nil {
		return permission_domain.ListTuplesResult{}, permission.ErrMisconfig
	}
	return s.perm.ListTuples(ctx, params)
}
