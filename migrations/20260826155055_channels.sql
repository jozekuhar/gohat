-- +goose Up
-- +goose StatementBegin
CREATE TABLE channels (
    id UUID PRIMARY KEY,
    organization_id UUID NOT NULL REFERENCES organizations (id),
    provider TEXT NOT NULL,
    name TEXT NOT NULL,
    credentials JSONB NOT NULL,
    status TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ,
    UNIQUE (organization_id, name),
    CONSTRAINT channels_provider_check CHECK (
        provider IN ('woocommerce', 'shopify')
    ),
    CONSTRAINT channels_status_check CHECK (
        provider IN ('active', 'inactive')
    )
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE channels;
-- +goose StatementEnd
