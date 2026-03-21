-- +goose Up
-- -- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE account (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    description TEXT,

    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 1. Content Types: Defines the "Blueprint" for projects
-- Example: slug='account-uuid_blog-post', 'account-uuid_portfolio-item', or 'account-uuid_bulletin'
CREATE TABLE content_types (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    account_id UUID REFERENCES account(id) ON DELETE CASCADE,

    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) UNIQUE NOT NULL,
    description TEXT,

    -- JSON schema to validate the 'data' in entries
    schema_definition JSONB NOT NULL,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 2. Entries: The actual content
CREATE TABLE entries (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    content_type_id UUID REFERENCES content_types(id) ON DELETE CASCADE,

    slug VARCHAR(255) NOT NULL,
    title VARCHAR(255) NOT NULL,

    -- This is where the magic happens: store nested JSON for different needs
    content_data JSONB NOT NULL DEFAULT '{}',
    status VARCHAR(50) DEFAULT 'draft', -- draft, published, archived

    version INT DEFAULT 1,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

    -- Ensure slug is unique within a specific content type
    UNIQUE(content_type_id, slug)
);

-- Index for high-performance JSONB filtering
CREATE INDEX idx_entries_content_data ON entries USING GIN (content_data);
CREATE INDEX idx_entries_status ON entries(status);

-- 3. Webhook Subscriptions
CREATE TABLE webhook_subscriptions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    account_id UUID REFERENCES account(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,

    target_url TEXT NOT NULL,
    secret_token TEXT NOT NULL, -- For HMAC signing in Go

    -- Array of events to listen to: e.g., {'bulletin.created', 'entry.published'}
    subscribed_events TEXT[] NOT NULL,

    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 4. Event & Audit Logs
-- Tracks what happened and the result of webhook deliveries
CREATE TABLE event_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    account_id UUID REFERENCES account(id) ON DELETE CASCADE,

    event_type VARCHAR(100) NOT NULL, -- e.g., 'ENTRY_PUBLISHED'
    entity_id UUID NOT NULL,          -- ID of the Entry or ContentType
    payload JSONB NOT NULL,           -- The state of the data at that time
    webhook_responses JSONB,          -- Stores HTTP status codes from subscribers
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- +goose Down
SELECT 'down SQL query';
