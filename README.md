# SSO BFF

This BFF is the integration layer in front of Hydra, Kratos, and Keto for `sso-portal`.

If an application wants to use the centralized SSO flow, it should integrate with:

- `sso-portal` for the login experience
- this BFF for session, identity, and authorization APIs

It should not integrate with Hydra, Kratos, or Keto directly if the BFF is already the standard gateway.

## Integration Model

The intended model is:

1. The client application calls the BFF.
2. The BFF redirects the browser into the configured SSO flow.
3. `sso-portal` handles the login UX.
4. Hydra, Kratos, and Keto stay behind the BFF and portal.
5. The BFF issues the session used by the client application.

For a consuming app, the BFF is the stable entrypoint.

## What An Integrating App Needs

If a team wants to onboard an app into the SSO flow, the app should only need:

- a registered app entry in the BFF
- a valid `dsn`
- an allowed redirect path
- access to the BFF base URL
- access to the BFF discovery endpoint

The app does not need to:

- talk to Hydra directly
- talk to Kratos directly
- talk to Keto directly
- host its own login flow if `sso-portal` is the shared login UI

## Service Discovery

The BFF exposes:

- `GET /discoveries`

This returns the available routes and methods so consuming applications can discover the BFF contract dynamically.

Example:

```sh
curl -s https://sso-staging.doj.gov.ph/discoveries
```

## Standard Login Flow

### 1. Register The App In The BFF

An admin registers the application in the BFF app registry.

Relevant routes:

- `POST /api/apps`
- `GET /api/apps`
- `GET /api/apps/{id}`
- `PUT /api/apps/{id}`
- `DELETE /api/apps/{id}`

App registration includes:

- `dsn`
- `redirect_path`

Example payload:

```json
{
  "dsn": "https://your-app.example.com",
  "redirect_path": "/dashboard"
}
```

### 2. Start Login From The Client App

The client app redirects the browser to the BFF login route:

```text
GET /api/login?dsn=https://your-app.example.com&redirect=/dashboard
```

You can also call it by explicit app id:

```text
GET /api/login?app=<app-id>&redirect=/dashboard
```

The BFF uses the configured app registry and redirects into the central SSO flow.

### 3. Let `sso-portal` Handle Login

The user authenticates through the existing portal flow.

At this stage:

- the app should not present a separate login screen
- the app should rely on the shared portal login managed by the BFF flow

### 4. BFF Receives The Callback

The BFF callback endpoint is:

- `GET /api/callback`

After the OAuth callback completes:

- the BFF creates its own session
- the BFF sets the session cookie
- the BFF resolves the final redirect target back to the client app

### 5. The App Reads Session State From The BFF

The client app checks:

- `GET /api/session`

Typical response fields include:

- `authenticated`
- `sub`
- `exp`
- `profile.name.first_name`
- `profile.name.last_name`
- `profile.username`
- `organization_id`
- `roles` when tuple lookup is requested

Example:

```json
{
  "authenticated": true,
  "sub": "user-123",
  "exp": "2026-03-31T10:00:00Z",
  "profile": {
    "name": {
      "first_name": "Paul",
      "last_name": "Test"
    },
    "username": "paul.test"
  },
  "organization_id": "org-456"
}
```

### 6. Refresh Or Logout Through The BFF

Routes:

- `POST /api/session/refresh`
- `POST /api/logout`

The app should use these BFF routes instead of implementing its own OAuth token/session flow.

## Authorization Model

Permissions are also mediated by the BFF.

Relevant routes:

- `GET /api/permissions`
- `GET /api/permissions/check`
- `POST /api/permissions/tuple`

The app can ask the BFF for authorization-related data instead of calling Keto directly.

## Identity Model

Identity-related operations are also exposed through the BFF.

Routes:

- `POST /api/identity/login`
- `GET /api/identity/settings`
- `POST /api/identity/settings`
- `GET /api/admin/identities`
- `POST /api/admin/identities`

For normal app integration, the main identity route used after login is usually:

- `GET /api/session`

## Minimum Requirements For A New App

If a team wants to connect a new app to the SSO flow, the practical requirements are:

- the app must be registered in the BFF app registry
- the app must know the BFF base URL
- the app must use `/discoveries` to confirm available routes if dynamic discovery is needed
- the app must redirect users to `/api/login`
- the app must rely on the BFF session cookie
- the app must call `/api/session` to determine login state
- the app must use `/api/logout` to end the session

If those are already in place, the app can use the centralized `sso-portal` flow without integrating Hydra, Kratos, or Keto on its own.

## Core Routes

### Discovery And Health

- `GET /discoveries`
- `GET /health`

### Login And Session

- `GET /api/login`
- `GET /api/callback`
- `GET /api/consent`
- `GET /api/launch`
- `GET /api/session`
- `POST /api/session/refresh`
- `POST /api/logout`

### Identity

- `POST /api/identity/login`
- `GET /api/identity/settings`
- `POST /api/identity/settings`
- `GET /api/admin/identities`
- `POST /api/admin/identities`

### Permissions

- `GET /api/permissions`
- `GET /api/permissions/check`
- `POST /api/permissions/tuple`

### Apps

- `GET /api/apps`
- `POST /api/apps`
- `GET /api/apps/{id}`
- `PUT /api/apps/{id}`
- `DELETE /api/apps/{id}`

### Audit

- `GET /api/audit/events`

## Runtime Requirements

The BFF itself still depends on:

- PostgreSQL
- Redis
- OAuth/Hydra
- Kratos
- Keto

Required environment variables include:

```env
DSN=postgres://postgres:postgres@localhost:5432/bff_db?sslmode=disable
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=your-redis-password

OAUTH_ADMIN=http://oauth-admin.example.com
OAUTH_PUBLIC=https://oauth.example.com
OAUTH_PRIVATE=http://oauth.example.com:4444
OIDC_ISSUER=https://oauth.example.com

BFF_CLIENT_ID=bff-client
BFF_CLIENT_SECRET=supersecret123
BFF_REDIRECT_URL=https://bff.example.com/api/callback

M2M_CLIENT_ID=m2m-client
M2M_CLIENT_SECRET=supersecret123
ALLOWED_CLIENT_IDS=bff-client,m2m-client
ALLOWED_SCOPES=openid,offline,api:read,api:write

IDENTITY_ADMIN=https://identity-admin.example.com
IDENTITY_PUBLIC=https://identity.example.com

PERMISSION_ADMIN=https://permission.example.com
PERMISSION_PUBLIC=https://permission.example.com

HTTP_ADDR=:8080
```

The service loads `.env` automatically only when `ENV != prod`.

## Local Run

```sh
go run ./cmd/main.go
```

Default listen address:

```text
:8080
```

## Docker

Available image files:

- `Dockerfile`
- `bff.Dockerfile`

`bff.Dockerfile` matches Compose setups that build the BFF from a nested app directory such as `./bff_go`.

## Database Migration Note

Migration files are in:

- `internal/db/migrations`

Current caveat:

- `000001_init_schema.down.sql` contains the `CREATE TABLE` statements
- `000001_init_schema.up.sql` contains the `DROP TABLE` statements

So the files are currently reversed.

If you are using the `migrate` CLI, fix or swap them first before running:

```sh
migrate -path ./internal/db/migrations \
  -database "postgres://postgres:postgres@localhost:5433/bff_db?sslmode=disable" \
  up
```

## Troubleshooting

### Missing BFF Redirect URL

Error:

```text
BFF client configuration missing. Set BFF_CLIENT_ID, BFF_CLIENT_SECRET, and BFF_REDIRECT_URL
```

Cause:

- `BFF_REDIRECT_URL` is missing

### Wrong OAuth Variable Names

The code expects:

- `OAUTH_ADMIN`
- `OAUTH_PUBLIC`
- `OAUTH_PRIVATE`

Not:

- `OAUTH_ADMIN_URL`
- `OAUTH_PUBLIC_URL`
- `OAUTH_PRIVATE_URL`

### Kratos Admin Returns `undefined response type`

This usually means the BFF is receiving HTML or another non-JSON response from the Kratos admin endpoint, often because of:

- nginx allow/deny rules
- source IP restrictions
- using an external admin URL instead of an internal trusted admin URL

If the endpoint works from the host but returns `403` from inside the BFF container, the problem is usually server/proxy access control, not frontend code and not the browser Kratos session.
