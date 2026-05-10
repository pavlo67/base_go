DROP TABLE IF EXISTS files;

CREATE TABLE files (
    path        TEXT PRIMARY KEY,
    is_dir      INTEGER NOT NULL,
    size        INTEGER NOT NULL,
    ctime       TEXT NOT NULL,
    mtime       TEXT NOT NULL,
    crc         INTEGER,
    mime_type   TEXT NOT NULL,
    created_at  TEXT NOT NULL,
    updated_at  TEXT NOT NULL
)
