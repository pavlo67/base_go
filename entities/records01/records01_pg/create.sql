DROP TABLE IF EXISTS records01;

CREATE TABLE records01 (
  id            BIGSERIAL                PRIMARY KEY,

  title         TEXT                     NOT NULL,
  summary       TEXT                     NOT NULL,
  record_type   VARCHAR(63)              NOT NULL,
  data          TEXT                     NOT NULL,
  embedded      TEXT                     ,

  urn           TEXT                     NOT NULL,
  tags          TEXT[]                   ,
  relations_map TEXT                     ,
  owner_nss     TEXT                     NOT NULL,
  viewer_nss    TEXT                     NOT NULL,
  history       TEXT                     NOT NULL, -- !!!
  created_at    TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at    TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_records01_title   ON records01(record_type,title);

