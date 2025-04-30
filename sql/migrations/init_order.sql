CREATE TABLE IF NOT EXISTS accruals (
accrual_order_id BIGINT PRIMARY KEY,
user_id INTEGER,
status TEXT NOT NULL,
accrual NUMERIC(10, 2) DEFAULT 0 CHECK (accrual >= 0),
uploaded_at TIMESTAMP DEFAULT NOW()
);
