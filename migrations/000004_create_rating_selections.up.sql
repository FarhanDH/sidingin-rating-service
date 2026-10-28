CREATE TABLE rating_selections (
    rating_id UUID NOT NULL,
    direction TEXT NOT NULL,
    rating_tag TEXT NOT NULL,
    PRIMARY KEY (rating_id, rating_tag), -- Every rating one A, can only have one A Tag.
    FOREIGN KEY (rating_id, direction) REFERENCES ratings (id, direction) ON DELETE CASCADE, -- Ensure one rating per order direction. And delete the related rows in this table, when a row in ratings table deleted.
    FOREIGN KEY (code, direction) REFERENCES rating_tags (code, direction) -- Ensure the code is belong to direction in rating_tags table.
)
