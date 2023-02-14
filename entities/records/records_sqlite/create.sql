DROP TABLE IF EXISTS records;

CREATE TABLE records (
  id           INTEGER    PRIMARY KEY AUTOINCREMENT,
  title        TEXT       NOT NULL,
  summary      TEXT       NOT NULL,
  type_key     TEXT       NOT NULL,
  data         TEXT       NOT NULL,
  embedded     TEXT       NOT NULL,

  urn          TEXT               ,
  tags         TEXT       NOT NULL,
  owner_nss    TEXT       NOT NULL,
  viewer_nss   TEXT       NOT NULL,
  history      TEXT       NOT NULL,
  created_at   TIMESTAMP  NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at   TIMESTAMP
);

CREATE UNIQUE INDEX idx_records_urn   ON records(urn)                          WHERE urn IS NOT NULL;
CREATE INDEX idx_records_viewer_title ON records(viewer_nss, type_key, title);
CREATE INDEX idx_records_owner_title  ON records(owner_nss,  type_key, title);
CREATE INDEX idx_records_type_title   ON records(            type_key, title);

