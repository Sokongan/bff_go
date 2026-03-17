package permission

import (
	"context"
	permission_domain "sso-bff/internal/domain/permission"
)

type TupleLister interface {
	ListTuples(
		ctx context.Context,
		params permission_domain.ListTuplesParams,
	) (permission_domain.ListTuplesResult, error)
}

type PermissionChecker interface {
	CheckTuple(
		ctx context.Context,
		tuple permission_domain.RelationTuple) (bool, error)

	ListTuples(
		ctx context.Context,
		params permission_domain.ListTuplesParams,
	) (permission_domain.ListTuplesResult, error)
}
