BEGIN;

CREATE TABLE events (
    id UUID PRIMARY KEY,
    user_id UUID,
    title TEXT NOT NULL,
    type TEXT NOT NULL,
    custom_type TEXT,
    description TEXT,
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()

    CONSTRAINT events_time_range_check CHECK (start_time < end_time),
    CONSTRAINT events_custom_type_check CHECK (
        (type = 'other' AND custom_type IS NOT NULL AND custom_type <> '')
        OR (type <> 'other' AND custom_type IS NULL)
    )
);

CREATE INDEX idx_events_user_id_start_time ON events (user_id, start_time);

COMMIT;