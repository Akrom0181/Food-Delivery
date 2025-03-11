CREATE TABLE IF NOT EXISTS orders (
  id UUID PRIMARY KEY,
  user_id UUID NOT NULL,
  total_price DECIMAL NOT NULL,
  status VARCHAR NOT NULL CHECK (status IN ('pending', 'confirmed', 'cancelled', 'preparing', 'picked_up', 'delivered')) DEFAULT 'pending',
  delivery_status VARCHAR NOT NULL CHECK (delivery_status IN ('olib ketish', 'yetkazib berish')),
  address VARCHAR,
  floor INT,
  door_number INT,
  entrance INT,
  latitude DECIMAL(10,8),
  longitude DECIMAL(11,8),
  created_at TIMESTAMP DEFAULT now(),
  updated_at TIMESTAMP DEFAULT now()
);


CREATE TABLE IF NOT EXISTS orderitems (
  id UUID PRIMARY KEY,
  order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
  product_id UUID NOT NULL REFERENCES product(id) ON DELETE CASCADE,
  total_price DECIMAL(10,2) NOT NULL,
  quantity INT NOT NULL,
  price DECIMAL NOT NULL,
  created_at TIMESTAMP DEFAULT 'now()',
  updated_at TIMESTAMP DEFAULT 'now()'
);


ALTER TABLE orders ADD COLUMN branch_id UUID REFERENCES branch(id) ON DELETE CASCADE;
ALTER TABLE orders ADD COLUMN courier_id UUID REFERENCES couriers(id) ON DELETE CASCADE;

CREATE EXTENSION cube;  -- for earthdistance
CREATE EXTENSION earthdistance; 
