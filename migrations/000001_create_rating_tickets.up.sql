CREATE TABLE rating_tickets (
    order_id TEXT PRIMARY KEY,
    customer_id TEXT NOT NULL,
    technician_id TEXT NOT NULL,
    completed_order_at TIMESTAMPTZ NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    reminder_count SMALLINT NOT NULL DEFAULT 0,
    last_reminder_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX ON rating_tickets (technician_id);
CREATE INDEX ON rating_tickets (customer_id, completed_order_at DESC);
