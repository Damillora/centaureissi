-- This file should undo anything in `up.sql`

DROP INDEX messages_content_hash;
DROP TABLE messages;
DROP TABLE user_tokens;
DROP TABLE users;
