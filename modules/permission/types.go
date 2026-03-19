package permission

import (
	"context"
	permission_domain "sso-bff/internal/domain/permission"
)

type PermissionClient interface {
	WriteTuple(
		ctx context.Context,
		tuple permission_domain.RelationTuple,
	) error

	CheckTuple(
		ctx context.Context,
		tuple permission_domain.RelationTuple,
	) (bool, error)

	ListTuples(
		ctx context.Context,
		params permission_domain.ListTuplesParams,
	) (permission_domain.ListTuplesResult, error)
}

type RolePayload struct {
	Object string `json:"object"`
	Role   string `json:"role"`
}
