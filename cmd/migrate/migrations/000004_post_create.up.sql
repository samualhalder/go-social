CREATE TABLE IF NOT EXISTS posts (
    id bigserial PRIMARY KEY,
    title text NOT NULL,
    context text NOT NULL,
    user_id bigint,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);
