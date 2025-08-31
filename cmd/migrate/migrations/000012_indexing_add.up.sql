-- Create the extension and indexes for full-text search
-- Check article: https://niallburkley.com/blog/index-columns-for-like-in-postgres/
CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- Note: CREATE INDEX does NOT support IF NOT EXISTS directly in PostgreSQL.
-- So we proceed with normal CREATE INDEX statements.

CREATE INDEX idx_comments_content ON comments USING gin (content gin_trgm_ops);
CREATE INDEX idx_posts_title ON posts USING gin (title gin_trgm_ops);
CREATE INDEX idx_posts_tags ON posts USING gin (tags);

CREATE INDEX idx_users_username ON users (username);
CREATE INDEX idx_posts_user_id ON posts (user_id);
CREATE INDEX idx_comments_post_id ON comments (post_id);
