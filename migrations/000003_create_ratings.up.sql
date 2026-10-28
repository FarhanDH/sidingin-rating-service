CREATE TABLE ratings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id TEXT NOT NULL REFERENCES rating_tickets (order_id),
    direction TEXT NOT NULL CHECK (direction IN ('customer_to_technician', 'customer_to_service', 'customer_to_app', 'technician_to_customer', 'technician_to_app')),
    stars SMALLINT NOT NULL CHECK(stars BETWEEN 1 AND 5),
    comment VARCHAR(255),
    is_comment_hidden BOOLEAN NOT NULL DEFAULT FALSE,
    is_requires_review BOOLEAN NOT NULL DEFAULT FALSE,
    app_platform TEXT, -- 'pwa', 'android', 'ios', only for ...to_app direction
    app_version TEXT, -- ...to_app direction
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(order_id, direction),  -- Ensure one rating per order direction.
    UNIQUE(id, direction) -- needed Composite FK for rating_selections
);
CREATE INDEX on ratings (created_at DESC);
CREATE INDEX on ratings (created_at) WHERE is_requires_review; -- partial index for admin queue
