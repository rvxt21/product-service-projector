CREATE DATABASE products;
\c products;
CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100),
    description TEXT,
    price INT,
    quantity INT,
    category VARCHAR(100),
    is_available BOOLEAN
);
CREATE TABLE IF NOT EXISTS category (
    idCategory SERIAL PRIMARY KEY,
    nameCategory VARCHAR(100) NOT NULL,
    descriptionCategory TEXT NOT NULL
);