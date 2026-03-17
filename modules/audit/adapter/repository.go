package adapter

import (
	"context"
	"database/sql"
	dbgen "sso-bff/internal/db/gen"
	audit_domain "sso-bff/internal/domain/audit"
)

type Repo struct {
	q *dbgen.Queries
}

func New(db *sql.DB) *Repo {
	return &Repo{q: dbgen.New(db)}
}

func (r *Repo) Insert(
	ctx context.Context,
	e audit_domain.AuditEvent,
) error {

	return r.q.InsertAuditEvent(ctx, dbgen.InsertAuditEventParams{
		IdentityID: e.IdentityID,
		ClientID:   toNull(e.ClientID),
		EventType:  e.EventType,
		IpAddress:  toNull(e.IPAddress),
		UserAgent:  toNull(e.UserAgent),
	})
}

func (r *Repo) ListRecent(
	ctx context.Context,
	identityID string,
	limit int32,
) ([]audit_domain.AuditEvent, error) {

	if limit <= 0 {
		limit = 50
	}

	rows, err := r.q.ListRecentAuditEvents(
		ctx,
		dbgen.ListRecentAuditEventsParams{
			IdentityID: identityID,
			Limit:      limit,
		})

	if err != nil {
		return nil, err
	}

	out := make([]audit_domain.AuditEvent, 0, len(rows))
	for _, r := range rows {
		out = append(out, audit_domain.AuditEvent{
			IdentityID: r.IdentityID,
			ClientID:   r.ClientID.String,
			EventType:  r.EventType,
			IPAddress:  r.IpAddress.String,
			UserAgent:  r.UserAgent.String,
		})
	}
	return out, nil
}

func toNull(v string) sql.NullString {
	if v == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: v, Valid: true}
}
