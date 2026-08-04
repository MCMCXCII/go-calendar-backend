BEGIN;

CREATE TABLE events (
    id UUID PRIMARY KEY,
    user_id UUID,
    title TEXT NOT NULL,
    category TEXT NOT NULL,
    custom_category TEXT,
    descriptions TEXT,
    start_time TIMESTAMPTZ,
    end_time TIMESTAMPTZ,
);

COMMIT;