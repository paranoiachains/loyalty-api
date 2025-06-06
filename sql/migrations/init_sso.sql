CREATE TABLE users (
user_id SERIAL PRIMARY KEY,
login TEXT UNIQUE NOT NULL,
password TEXT NOT NULL,
balance NUMERIC(10, 2) DEFAULT 0 CONSTRAINT balance_nonnegative CHECK (balance >= 0),
withdrawn NUMERIC(10, 2) DEFAULT 0 CHECK (withdrawn >= 0)
);

CREATE TABLE IF NOT EXISTS withdrawals (
order_id BIGINT PRIMARY KEY,
user_id INTEGER NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
sum NUMERIC(10, 2) NOT NULL CHECK (sum > 0),
processed_at TIMESTAMP DEFAULT NOW()
);