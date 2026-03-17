package permission

import (
	"context"
	"errors"
	"fmt"
	permission_domain "sso-bff/internal/domain/permission"
	permission_sdk "sso-bff/modules/permission/sdk"

	keto "github.com/ory/keto-client-go/v25"
)

type PermissionGatewayType struct {
	admin  *keto.APIClient
	public *keto.APIClient
}

func NewPermissionGateway(
	s *permission_sdk.PermissionSDK,
) (*PermissionGatewayType, error) {

	if s == nil || s.Admin == nil || s.Public == nil {
		return nil, errors.New("keto sdk not configured")
	}

	return &PermissionGatewayType{
		admin:  s.Admin,
		public: s.Public,
	}, nil
}

func (g *PermissionGatewayType) WriteTuple(
	ctx context.Context,
	t permission_domain.RelationTuple) error {
	if g == nil || g.admin == nil {
		return errors.New("keto admin client not configured")
	}

	body := keto.CreateRelationshipBody{
		Namespace: &t.Namespace,
		Object:    &t.Object,
		Relation:  &t.Relation,
		SubjectId: &t.SubjectID,
	}

	_, _, err := g.admin.RelationshipAPI.
		CreateRelationship(ctx).
		CreateRelationshipBody(body).
		Execute()
	if err != nil {
		return fmt.Errorf("keto write tuple: %w", err)
	}
	return nil
}

func (g *PermissionGatewayType) CheckTuple(
	ctx context.Context,
	t permission_domain.RelationTuple,
) (bool, error) {
	if g == nil || g.public == nil {
		return false, errors.New("keto public client not configured")
	}

	resp, _, err := g.public.PermissionAPI.
		CheckPermission(ctx).
		Namespace(t.Namespace).
		Object(t.Object).
		Relation(t.Relation).
		SubjectId(t.SubjectID).
		Execute()
	if err != nil {
		return false, fmt.Errorf("keto check tuple: %w", err)
	}

	return resp.GetAllowed(), nil
}

func (g *PermissionGatewayType) ListTuples(
	ctx context.Context,
	params permission_domain.ListTuplesParams,
) (permission_domain.ListTuplesResult, error) {

	if g == nil || g.public == nil {
		return permission_domain.ListTuplesResult{},
			errors.New("permission not configured")
	}

	req := g.public.RelationshipAPI.GetRelationships(ctx)

	if params.PageSize > 0 {
		req = req.PageSize(params.PageSize)
	}

	if params.PageToken != "" {
		req = req.PageToken(params.PageToken)
	}

	if params.Namespace != "" {
		req = req.Namespace(params.Namespace)
	}

	if params.Object != "" {
		req = req.Object(params.Object)
	}

	if params.Relation != "" {
		req = req.Relation(params.Relation)
	}

	if params.SubjectID != "" {
		req = req.SubjectId(params.SubjectID)
	}

	resp, _, err := req.Execute()

	if err != nil {
		return permission_domain.ListTuplesResult{},
			fmt.Errorf("permission list tuples: %w", err)
	}

	out := permission_domain.ListTuplesResult{
		Tuples:        make([]permission_domain.RelationTuple, 0),
		NextPageToken: resp.GetNextPageToken(),
	}

	for _, rel := range resp.GetRelationTuples() {
		tuple := permission_domain.RelationTuple{
			Namespace: rel.GetNamespace(),
			Object:    rel.GetObject(),
			Relation:  rel.GetRelation(),
		}

		if rel.SubjectId != nil {
			tuple.SubjectID = *rel.SubjectId
		}

		out.Tuples = append(out.Tuples, tuple)
	}

	return out, nil
}
