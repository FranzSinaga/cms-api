CREATE TABLE posts (
  id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  title         VARCHAR(225) NOT NULL,
  slug          VARCHAR(255) NOT NULL UNIQUE,
  description   TEXT,
  thumbnail_url VARCHAR(500),
  content       TEXT NOT NULL,
  status        VARCHAR(50) NOT NULL DEFAULT 'draft',
  created_by    UUID NOT NULL REFERENCES users(id),
  updated_by    UUID NOT NULL REFERENCES users(id),
  published_at  TIMESTAMP,
  created_at    TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at    TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_posts_slug ON posts(slug);
CREATE INDEX idx_posts_status ON posts(status);