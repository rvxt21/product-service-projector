
-- CREATE DATABASE products;

\c products;

CREATE TABLE IF NOT EXISTS categories (
    idCategory SERIAL PRIMARY KEY,
    nameCategory VARCHAR(100) NOT NULL,
    descriptionCategory TEXT NOT NULL
);


CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100),
    description TEXT,
    price INT,
    quantity INT,
    category INT,
    is_available BOOLEAN,
    FOREIGN KEY (category) REFERENCES categories(idCategory) 
);
