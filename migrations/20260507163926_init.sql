-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id UUID PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    google_id TEXT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ
);

CREATE TABLE sessions (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users (id),
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE organizations (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    slug TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ,
    CONSTRAINT organizations_slug_check CHECK (slug = LOWER(slug))
);

CREATE TABLE memberships (
    id UUID PRIMARY KEY,
    organization_id UUID NOT NULL REFERENCES organizations (id),
    user_id UUID NOT NULL REFERENCES users (id),
    role TEXT NOT NULL,
    permissions TEXT[],
    status TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ,
    UNIQUE (organization_id, user_id),
    CONSTRAINT memberships_role_check CHECK (
        role IN ('owner', 'admin', 'member')
    ),
    CONSTRAINT memberships_status_check CHECK (
        status IN ('active', 'inactive')
    )
);

CREATE TABLE invitations (
    id UUID PRIMARY KEY,
    organization_id UUID NOT NULL REFERENCES organizations (id),
    role TEXT NOT NULL,
    permissions TEXT[],
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT invitations_role_check CHECK (
        role IN ('owner', 'admin', 'member')
    )
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE organization_memberships;
DROP TABLE organizations;
DROP TABLE sessions;
DROP TABLE users;
-- +goose StatementEnd
