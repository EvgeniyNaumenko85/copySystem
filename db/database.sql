CREATE DATABASE files_db;

CREATE TABLE IF NOT EXISTS users
(
    id                  BIGSERIAL PRIMARY KEY,
    user_name           VARCHAR(30) NOT NULL UNIQUE,
    email               VARCHAR(30) NOT NULL UNIQUE,
    password_hash       VARCHAR(255) NOT NULL,
    user_role           VARCHAR(30) NOT NULL,
    file_size_lim       INT DEFAULT 20,
    storage_size_lim    INT DEFAULT 100,
);

CREATE TABLE IF NOT EXISTS files
(
    id               BIGSERIAL PRIMARY KEY,
    user_id          INTEGER REFERENCES users (id) ON DELETE CASCADE,
    file_name        VARCHAR(30)  NOT NULL,
    extension        VARCHAR(10),
    file_path        VARCHAR(255) NOT NULL UNIQUE,
    file_size        DOUBLE PRECISION NOT NULL,
    deleted          BOOLEAN      NOT NULL DEFAULT false,
    access_to_all    BOOLEAN      NOT NULL DEFAULT false,
    added            TIMESTAMP without time zone NOT NULL DEFAULT to_timestamp('2023-06-25 12:53:00', 'YYYY-MM-DD HH24:MI:SS'),
);

CREATE TABLE IF NOT EXISTS access
(
    id               BIGSERIAL PRIMARY KEY,
    file_id          INTEGER REFERENCES files (id) ON DELETE CASCADE,
    user_id          INTEGER REFERENCES users (id) ON DELETE CASCADE,
    CONSTRAINT unique_user_file_pair UNIQUE (user_id, file_id)
    );