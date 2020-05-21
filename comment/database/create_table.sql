CREATE TABLE comments (
  id SERIAL PRIMARY KEY,
  org VARCHAR(64) NOT NULL,
  author VARCHAR(64) NOT NULL,
  comment VARCHAR(512) NOT NULL,
  is_deleted BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMP,
  updated_at TIMESTAMP
);