CREATE TABLE orders (
order_id BIGINT PRIMARY KEY,
user_id INTEGER NOT NULL,
status TEXT NOT NULL,
accrual NUMERIC(10, 2) DEFAULT 0 CHECK (accrual >= 0)
);