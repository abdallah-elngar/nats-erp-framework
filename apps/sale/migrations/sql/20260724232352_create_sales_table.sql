-- name: create_sales_table
-- id: 20260724232352
-- created: 2026-07-24 23:23:52

-- up:
CREATE TABLE IF NOT EXISTS sales (
    id BIGSERIAL PRIMARY KEY,
    sa_name VARCHAR(255) NOT NULL UNIQUE,
    price DECIMAL(10,2),
    quantity INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- down:
DROP TABLE IF EXISTS sales;
