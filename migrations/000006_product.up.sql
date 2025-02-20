CREATE TABLE IF NOT EXISTS category (
  id UUID PRIMARY KEY,
  name VARCHAR NOT NULL,
  created_at timestamp NOT NULL DEFAULT 'now()',
  updated_at timestamp NOT NULL DEFAULT 'now()'
);

CREATE TABLE IF NOT EXISTS product (
  id UUID PRIMARY KEY,
  category_id UUID NOT NULL REFERENCES category(id) ON DELETE CASCADE,
  name VARCHAR NOT NULL,
  description TEXT,
  price DECIMAL NOT NULL,
  image_url VARCHAR,
  created_at timestamp NOT NULL DEFAULT 'now()',
  updated_at timestamp NOT NULL DEFAULT 'now()'
);

CREATE TABLE IF NOT EXISTS banner (
  id UUID PRIMARY KEY,  
  title VARCHAR NOT NULL,
  image_url VARCHAR NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT 'now()',
  updated_at TIMESTAMP NOT NULL DEFAULT 'now()'
);