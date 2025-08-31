ALTER TABLE posts RENAME COLUMN context TO content;

ALTER TABLE posts ADD COLUMN updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW();

ALTER TABLE posts ADD COLUMN tags varchar(200)[];
