CREATE TABLE rating_tags (
    code TEXT PRIMARY KEY,
    direction TEXT NOT NULL CHECK (direction IN ('customer_to_technician', 'customer_to_service', 'customer_to_app', 'technician_to_customer', 'technician_to_app')),
    label TEXT NOT NULL,
    min_stars SMALLINT NOT NULL DEFAULT 1,
    max_stars SMALLINT NOT NULL DEFAULT 5,
    sort_order SMALLINT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    CHECK (min_stars BETWEEN 1 AND 5 AND max_stars BETWEEN min_stars AND 5),
    UNIQUE(code, direction) -- Needed for Composite Foreign Key in rating_selections table.
);
