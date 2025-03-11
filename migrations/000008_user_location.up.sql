CREATE TABLE IF NOT EXISTS user_location (
  id UUID PRIMARY KEY,
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  address VARCHAR,
  floor INT,
  door_number INT,
  entrance INT,
  latitude DECIMAL(10,8),
  longitude DECIMAL(11,8),
  created_at TIMESTAMP DEFAULT 'now()',
  updated_at TIMESTAMP DEFAULT 'now()'
);