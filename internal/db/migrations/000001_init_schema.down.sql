-- enable extension (required for uuid generation)
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE public.app_registry (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    dsn TEXT NOT NULL,
    redirect_path TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE public.audit_events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    identity_id TEXT NOT NULL,
    client_id TEXT,
    event_type TEXT NOT NULL,
    ip_address TEXT,
    user_agent TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);