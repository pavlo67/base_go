CREATE TABLE storage (
  id           BIGSERIAL                PRIMARY KEY,
  data_key     TEXT                     NOT NULL,
  url          TEXT                     NOT NULL,
  title        TEXT                     NOT NULL,
  summary      TEXT                     NOT NULL,
  embedded     TEXT                     ,
  tags         TEXT                     ,
  type_key     TEXT                     NOT NULL,
  content      TEXT                     ,
  history      TEXT                     ,
  created_at   TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at   TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_storage_key   ON storage(data_key);

CREATE INDEX idx_storage_title ON storage(type_key, title);

--

CREATE TABLE flow (
  id           BIGSERIAL                PRIMARY KEY,
  data_key     TEXT                     NOT NULL,
  url          TEXT                     NOT NULL,
  title        TEXT                     NOT NULL,
  summary      TEXT                     NOT NULL,
  embedded     TEXT                     ,
  tags         TEXT                     ,
  type_key     TEXT                     NOT NULL,
  content      TEXT                     ,
  history      TEXT                     ,
  created_at   TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at   TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_flow_key   ON flow(data_key);

CREATE INDEX idx_flow_title ON flow(type_key, title);
