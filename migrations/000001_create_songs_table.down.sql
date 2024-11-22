DROP INDEX IF EXISTS idx_songs_title;
DROP INDEX IF EXISTS idx_songs_group;
DROP INDEX IF EXISTS idx_songs_release_date;

DROP TABLE IF EXISTS songs;

DROP EXTENSION IF EXISTS "pgcrypto";