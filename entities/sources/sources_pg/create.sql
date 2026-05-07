CREATE TABLE sources (
    id                BIGSERIAL                PRIMARY KEY,
    key               TEXT                     NOT NULL,
    options           TEXT                     NOT NULL,
    type_key          TEXT                     NOT NULL,
    content           TEXT                     NOT NULL,
    history           TEXT                     NOT NULL,
    created_at        TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_sources_identity_key ON sources(identity_key);

CREATE INDEX idx_sources_type_key     ON sources(type_key);

CREATE INDEX idx_sources_created_at   ON sources(created_at);



