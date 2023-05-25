CREATE TABLE IF NOT EXISTS orders (
    uid VARCHAR(255) NOT NULL UNIQUE,
    details JSON NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_orders_uid ON orders(uid);
