-- +goose Up

CREATE TYPE maintenance_status AS ENUM (
    'Pending',
    'Assigned',
    'In Progress',
    'Resolved',
    'Rejected'
);

CREATE TYPE escalation_level AS ENUM (
    'None',
    'Level1',
    'Level2',
    'Level3'
);

CREATE TABLE maintenance_categories (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL UNIQUE,
    description TEXT NOT NULL DEFAULT '',
    manager_email TEXT NOT NULL,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TABLE maintenance_subcategories (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    category_id UUID NOT NULL REFERENCES maintenance_categories(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    supervisor_email TEXT NOT NULL,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TABLE maintenance_details (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    subcategory_id UUID NOT NULL REFERENCES maintenance_subcategories(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TABLE maintenance_workers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL,
    user_email TEXT NOT NULL DEFAULT '',
    phone TEXT NOT NULL,
    subcategory_id UUID NOT NULL REFERENCES maintenance_subcategories(id) ON DELETE RESTRICT,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TABLE maintenance_requests (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    requester_email TEXT NOT NULL,
    requester_name TEXT NOT NULL,
    location TEXT NOT NULL,
    category_id UUID NOT NULL REFERENCES maintenance_categories(id) ON DELETE RESTRICT,
    subcategory_id UUID NOT NULL REFERENCES maintenance_subcategories(id) ON DELETE RESTRICT,
    detail_id UUID REFERENCES maintenance_details(id) ON DELETE SET NULL,
    description TEXT NOT NULL,
    status maintenance_status NOT NULL DEFAULT 'Pending',
    escalation_level escalation_level NOT NULL DEFAULT 'None',
    last_escalated_at TIMESTAMP WITH TIME ZONE,
    assigned_worker_id UUID REFERENCES maintenance_workers(id) ON DELETE SET NULL,
    assigned_at TIMESTAMP WITH TIME ZONE,
    resolved_at TIMESTAMP WITH TIME ZONE,
    resolution_notes TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_maintenance_requests_requester ON maintenance_requests(requester_email);
CREATE INDEX idx_maintenance_requests_category ON maintenance_requests(category_id);
CREATE INDEX idx_maintenance_requests_subcategory ON maintenance_requests(subcategory_id);
CREATE INDEX idx_maintenance_requests_status ON maintenance_requests(status);
CREATE INDEX idx_maintenance_requests_escalation ON maintenance_requests(escalation_level);
CREATE INDEX idx_maintenance_requests_status_created ON maintenance_requests(status, created_at);

CREATE TABLE maintenance_status_log (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    request_id UUID NOT NULL REFERENCES maintenance_requests(id) ON DELETE CASCADE,
    user_email TEXT NOT NULL,
    action TEXT NOT NULL,
    previous_status maintenance_status,
    new_status maintenance_status,
    comments TEXT NOT NULL DEFAULT '',
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_maintenance_status_log_request ON maintenance_status_log(request_id);

CREATE TABLE maintenance_escalation_log (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    request_id UUID NOT NULL REFERENCES maintenance_requests(id) ON DELETE CASCADE,
    escalation_level escalation_level NOT NULL,
    notified_emails TEXT[] NOT NULL,
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_maintenance_escalation_log_request ON maintenance_escalation_log(request_id);

CREATE TABLE maintenance_config (
    id INTEGER PRIMARY KEY DEFAULT 1 CHECK (id = 1),
    dean_email TEXT NOT NULL DEFAULT '',
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

INSERT INTO maintenance_config (id, dean_email) VALUES (1, '');

-- +goose Down
DROP TABLE IF EXISTS maintenance_escalation_log;
DROP TABLE IF EXISTS maintenance_status_log;
DROP TABLE IF EXISTS maintenance_requests;
DROP TABLE IF EXISTS maintenance_workers;
DROP TABLE IF EXISTS maintenance_details;
DROP TABLE IF EXISTS maintenance_subcategories;
DROP TABLE IF EXISTS maintenance_categories;
DROP TABLE IF EXISTS maintenance_config;
DROP TYPE IF EXISTS maintenance_status;
DROP TYPE IF EXISTS escalation_level;