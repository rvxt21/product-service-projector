CREATE DATABASE products;
\c products;
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100),
    description TEXT,
    price INT,
    quantity INT,
    category VARCHAR(100),
    is_available BOOLEAN
);