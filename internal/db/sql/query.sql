
-- name: InsertAuditEvent :exec
INSERT INTO audit_events (
  identity_id, client_id, event_type, ip_address, user_agent
) VALUES (
  $1, $2, $3, $4, $5
);

-- name: ListRecentAuditEvents :many
SELECT id, identity_id, client_id, event_type, ip_address, user_agent, created_at
FROM audit_events
WHERE identity_id = $1
ORDER BY created_at DESC
LIMIT $2;

-- name: CreateAppRegistry :one
INSERT INTO app_registry (
  dsn, redirect_path
) VALUES (
  $1, $2
)
RETURNING id, dsn, redirect_path, created_at, updated_at;

-- name: UpdateAppRegistry :one
UPDATE app_registry
SET dsn = $2,
    redirect_path = $3,
    updated_at = now()
WHERE id = $1
RETURNING id, dsn, redirect_path, created_at, updated_at;

-- name: GetAppRegistry :one
SELECT id, dsn, redirect_path, created_at, updated_at
FROM app_registry
WHERE id = $1;

-- name: ListAppRegistries :many
SELECT id, dsn, redirect_path, created_at, updated_at
FROM app_registry
ORDER BY created_at DESC;

-- name: DeleteAppRegistry :exec
DELETE FROM app_registry
WHERE id = $1;
