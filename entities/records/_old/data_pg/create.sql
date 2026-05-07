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
  owner_key    TEXT                     NOT NULL,
  viewer_key   TEXT                     NOT NULL,
  history      TEXT                     ,
  created_at   TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at   TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_storage_key    ON storage(data_key);

CREATE INDEX idx_storage_owner  ON storage(owner_key, type_key, created_at);

CREATE INDEX idx_storage_viewer ON storage(viewer_key, type_key, created_at);
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
  owner_key    TEXT                     NOT NULL,
  viewer_key   TEXT                     NOT NULL,
  history      TEXT                     ,
  created_at   TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at   TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_flow_key   ON flow(data_key);

CREATE INDEX idx_flow_owner  ON flow(owner_key, type_key, created_at);

CREATE INDEX idx_flow_viewer ON flow(viewer_key, type_key, created_at);
