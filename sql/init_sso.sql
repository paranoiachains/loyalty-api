CREATE TABLE users (
user_id SERIAL PRIMARY KEY,
login TEXT UNIQUE NOT NULL,
password TEXT NOT NULL,
balance NUMERIC(10, 2) DEFAULT 0 CHECK (balance >= 0),
withdrawn NUMERIC(10, 2) DEFAULT 0 CHECK (withdrawn >= 0)
);