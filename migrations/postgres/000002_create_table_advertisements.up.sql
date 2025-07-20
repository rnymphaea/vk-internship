CREATE TABLE IF NOT EXISTS advertisements (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  author_id UUID NOT NULL,
  caption VARCHAR(128) NOT NULL CHECK(length(caption) BETWEEN 3 AND 128),
  description VARCHAR(1024) NOT NULL CHECK(length(description) < 1024),
  image_url VARCHAR(1024) CHECK(length(image_url) < 1024),
  price INTEGER NOT NULL CHECK(price > 0),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMPTZ,
  FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE CASCADE
);
