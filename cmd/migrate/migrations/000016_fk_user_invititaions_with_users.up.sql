ALTER TABLE user_invitations
ADD CONSTRAINT fk_user_invitations_user
FOREIGN KEY (user_id) REFERENCES users(id)
ON DELETE CASCADE
ON UPDATE CASCADE;
