package permission

import (
	"context"
	permission_domain "sso-bff/internal/domain/permission"
)

type PermissionService struct {
	perm PermissionClient
}

func NewPermissionService(perm PermissionClient) *PermissionService {
	return &PermissionService{perm: perm}
}

func (s *PermissionService) WriteTuple(
	ctx context.Context,
	tuple permission_domain.RelationTuple,
) error {
	if s == nil || s.perm == nil {
		return ErrMisconfig
	}
	return s.perm.WriteTuple(ctx, tuple)
}

func (s *PermissionService) CheckTuple(
	ctx context.Context,
	tuple permission_domain.RelationTuple,
) (bool, error) {

	if s == nil || s.perm == nil {
		return false, ErrMisconfig
	}

	return s.perm.CheckTuple(ctx, tuple)
}

func (s *PermissionService) ListTuples(
	ctx context.Context,
	params permission_domain.ListTuplesParams,
) (permission_domain.ListTuplesResult, error) {

	if s == nil || s.perm == nil {
		return permission_domain.ListTuplesResult{},
			ErrMisconfig
	}
	return s.perm.ListTuples(ctx, params)
}
