-- +goose Up
-- +goose StatementBegin
CREATE TABLE channels (
    id UUID PRIMARY KEY,
    organization_id UUID NOT NULL REFERENCES organizations (id),
    provider TEXT NOT NULL,
    name TEXT NOT NULL,
    url TEXT NOT NULL,
    credentials TEXT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ,
    CONSTRAINT channels_provider_check CHECK (
        provider IN ('woocommerce', 'shopify')
    )
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE channels;
-- +goose StatementEnd
