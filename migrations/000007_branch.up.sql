CREATE TABLE IF NOT EXISTS branch (
  id UUID PRIMARY KEY,
  name VARCHAR NOT NULL,
  address VARCHAR NOT NULL,
  latitude DECIMAL(10,8),
  longitude DECIMAL(11,8),
  phone_number VARCHAR,
  created_at TIMESTAMP DEFAULT 'now()',
  updated_at TIMESTAMP DEFAULT 'now()'
);