# SSO Portal And BFF Live Test Report

## Status

The SSO platform is now live and ready for testing.

- SSO Portal: `https://sso-staging.doj.gov.ph`
- BFF: `https://sso-bff-staging.gov.ph`

The intended integration model is:

- users interact with `sso-portal` for login
- client applications integrate with the BFF
- Hydra, Kratos, and Keto are used behind the BFF

Client applications should not integrate directly with Hydra, Kratos, or Keto when the BFF is already the standard gateway.

## Login Flow Based On The BFF

1. The client application calls the BFF login endpoint.
2. The BFF resolves the app by `dsn` or `app` id.
3. The BFF redirects the browser into the centralized SSO flow.
4. `sso-portal` handles the login UX.
5. Hydra completes the OAuth flow.
6. The BFF receives the callback at `/api/callback`.
7. The BFF creates the BFF session and sets the session cookie.
8. The client application reads the authenticated user state from `/api/session`.

## Minimum App Requirements

To test a new application against the live SSO flow, the app should have:

- a registered app entry in the BFF
- a valid `dsn`
- an allowed `redirect_path`
- access to the BFF base URL
- usage of the BFF discovery endpoint

The app should use:

- `GET /discoveries`
- `GET /api/login`
- `GET /api/session`
- `POST /api/session/refresh`
- `POST /api/logout`

## Discovery Endpoint

The BFF exposes an API discovery document at:

```text
GET /discoveries
```

Example:

```sh
curl -s https://sso-bff-staging.gov.ph/discoveries
```

This returns the routes and supported HTTP methods published by the BFF.

## Test Login Usage

Example login request by `dsn`:

```text
GET https://sso-bff-staging.gov.ph/api/login?dsn=https://your-app.example.com&redirect=/dashboard
```

Example login request by app id:

```text
GET https://sso-bff-staging.gov.ph/api/login?app=<app-id>&redirect=/dashboard
```

After login, the client application should call:

```text
GET https://sso-bff-staging.gov.ph/api/session
```

Example session response:

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

## Ory Stack Usage Through The BFF

### Hydra

Hydra is used for the OAuth2 and consent flow behind the BFF.

Typical client configuration for the BFF flow:

```json
{
  "client_id": "bff-client",
  "client_secret": "your-client-secret",
  "grant_types": [
    "authorization_code",
    "refresh_token",
    "client_credentials"
  ],
  "response_types": [
    "code"
  ],
  "scope": "openid offline api:read api:write",
  "redirect_uris": [
    "https://sso-bff-staging.gov.ph/api/callback"
  ]
}
```

Important notes:

- the redirect URI should point to the BFF callback, not directly to the app
- the app uses the BFF session after login, not the raw Hydra flow

### Kratos

Kratos is used for identity and authentication behind the BFF.

The BFF currently reads and exposes traits such as:

- `name.firstName`
- `name.lastName`
- `username`

The BFF also reads:

- `metadata_public.organization_id`

Example identity shape used by the BFF:

```json
{
  "schema_id": "user",
  "traits": {
    "name": {
      "firstName": "Juan",
      "lastName": "Dela Cruz"
    },
    "email": "juan.delacruz@example.com",
    "username": "doj123456"
  },
  "metadata_public": {
    "organization_id": "115f5f72-4afc-4f4f-8c7a-980322a84567"
  }
}
```

Minimum schema expectation for BFF compatibility:

- `traits.name.firstName`
- `traits.name.lastName`
- `traits.username`

Optional but currently used in session payload:

- `metadata_public.organization_id`

### Keto

Keto is used for authorization behind the BFF.

The BFF works with relationship tuples in this shape:

```json
{
  "namespace": "app",
  "object": "sso-portal",
  "relation": "admin",
  "subject_id": "identity-id"
}
```

The current BFF code explicitly checks this tuple for admin flows:

- `namespace`: `app`
- `object`: `sso-portal`
- `relation`: `admin`

That means a user must have the `app:sso-portal#admin` relationship to use admin endpoints such as:

- `GET /api/admin/identities`
- `POST /api/admin/identities`
- app registry management endpoints

## API Documentation

### Health And Discovery

- `GET /health`
- `GET /discoveries`

### OAuth And Session

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

- `POST /api/permissions/tuple`
- `GET /api/permissions/check`
- `GET /api/permissions`

### Apps

- `GET /api/apps`
- `POST /api/apps`
- `GET /api/apps/{id}`
- `PUT /api/apps/{id}`
- `DELETE /api/apps/{id}`

### Audit

- `GET /api/audit/events`

## Recommended Test Sequence

1. Call `/discoveries` to confirm the BFF routes are available.
2. Ensure the target application is registered in the BFF app registry.
3. Start login through `/api/login`.
4. Complete login in `sso-portal`.
5. Confirm the BFF session through `/api/session`.
6. If needed, verify authorization through `/api/permissions/check`.
7. Verify logout through `/api/logout`.

## Notes

- The BFF is the main integration point for application teams.
- `sso-portal` is the shared login experience.
- Hydra, Kratos, and Keto remain infrastructure components behind the BFF.
- API route discovery should start from `/discoveries`.
