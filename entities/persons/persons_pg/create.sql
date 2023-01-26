DROP TABLE IF EXISTS persons;

CREATE TABLE persons (
  id            BIGSERIAL                PRIMARY KEY,

  firstnames    VARCHAR(63)[]            NOT NULL,
  middlename    VARCHAR(63)              NOT NULL,
  lastname      TEXT                     NOT NULL,
  nicknames     VARCHAR(63)[]            NOT NULL,
  contacts      TEXT                     ,
  info          TEXT                     ,

  urn           TEXT                     NOT NULL,
  tags          TEXT[]                   ,
  relations_map TEXT                     ,
  owner_nss     TEXT                     NOT NULL,
  viewer_nss    TEXT                     NOT NULL,
  history       TEXT                     ,
  created_at    TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at    TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_lastname   ON persons(lastname);

