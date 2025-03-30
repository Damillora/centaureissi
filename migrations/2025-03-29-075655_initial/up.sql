-- Your SQL goes here
CREATE TABLE users (
    id INTEGER PRIMARY KEY NOT NULL,
    username varchar(100) UNIQUE NOT NULL,
    password varchar(100) NOT NULL,

    created_at datetime DEFAULT CURRENT_TIMESTAMP,
    updated_at datetime
);
SELECT diesel_manage_updated_at('users');

CREATE TABLE user_tokens (
    id INTEGER PRIMARY KEY NOT NULL,
    user_id INTEGER NOT NULL,
    token varchar(100) NOT NULL,
    
    revoked_at datetime NULL,
    created_at datetime DEFAULT CURRENT_TIMESTAMP,
    updated_at datetime,

    FOREIGN KEY (user_id) REFERENCES users(user_id)
);
SELECT diesel_manage_updated_at('user_tokens');


CREATE TABLE messages (
    id INTEGER PRIMARY KEY NOT NULL,
    user_id INTEGER NOT NULL,
    content_hash varchar(128) NOT NULL,
    
    created_at datetime DEFAULT CURRENT_TIMESTAMP,
    updated_at datetime,

    FOREIGN KEY (user_id) REFERENCES users(user_id)
);

CREATE INDEX messages_content_hash ON messages (content_hash);

SELECT diesel_manage_updated_at('messages');
