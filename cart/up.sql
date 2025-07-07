CREATE TABLE IF NOT EXISTS carts (
  id CHAR(27) PRIMARY KEY,
  account_id CHAR(27) NOT NULL UNIQUE,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS cart_items (
  cart_id CHAR(27) REFERENCES carts(id) ON DELETE CASCADE,
  product_id CHAR(27) NOT NULL,
  quantity INT NOT NULL CHECK (quantity > 0),
  PRIMARY KEY (cart_id, product_id)
);
