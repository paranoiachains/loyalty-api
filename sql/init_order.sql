CREATE TABLE IF NOT EXISTS accruals (
accrual_order_id BIGINT PRIMARY KEY,
user_id INTEGER,
status TEXT NOT NULL,
accrual NUMERIC(10, 2) DEFAULT 0 CHECK (accrual >= 0),
uploaded_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS withdrawals (
order_id BIGINT PRIMARY KEY,
user_id INTEGER,
sum NUMERIC(10, 2) NOT NULL CHECK (sum > 0),
processed_at TIMESTAMP DEFAULT NOW()
);