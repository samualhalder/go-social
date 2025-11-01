CREATE TABLE IF NOT EXISTS roles (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    level INT DEFAULT 0,
    description TEXT
);

INSERT INTO roles (name, level, description)
VALUES
('user', 1, 'User can create and delete their own posts and comments'),
('moderator', 2, 'Moderator can update any post and comment'),
('admin', 3, 'Admin can update and delete any post and comment');
