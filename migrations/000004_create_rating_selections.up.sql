CREATE TABLE rating_selections (
    rating_id UUID NOT NULL,
    direction TEXT NOT NULL,
    tag_code TEXT NOT NULL,
    PRIMARY KEY (rating_id, tag_code), -- Every one rating A, can only have one A Tag.
    FOREIGN KEY (rating_id, direction) REFERENCES ratings (id, direction) ON DELETE CASCADE, -- Ensure one rating per order direction. And delete the related rows in this table, when a row in ratings table deleted.
    FOREIGN KEY (tag_code, direction) REFERENCES rating_tags (code, direction) -- Ensure the code is belong to direction in rating_tags table.
)
