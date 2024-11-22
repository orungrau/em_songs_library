CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE songs (
                       id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                       title VARCHAR(255) NOT NULL,
                       text TEXT,
                       link VARCHAR(255),
                       "group" VARCHAR(255) NOT NULL,
                       release_date TIMESTAMP NOT NULL,
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                       deleted_at TIMESTAMP
);

CREATE INDEX idx_songs_title ON songs (title);
CREATE INDEX idx_songs_group ON songs ("group");
CREATE INDEX idx_songs_release_date ON songs (release_date);
