
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
    FOREIGN KEY (category) REFERENCES categories(idCategory) ON DELETE SET NULL
);

INSERT INTO categories (nameCategory, descriptionCategory) VALUES
('Electronics', 'Devices and gadgets'),
('Clothing', 'Apparel and accessories'),
('Books', 'Printed and digital books'),
('Home Appliances', 'Appliances for household use'),
('Toys', 'Toys and games for children'),
('Furniture', 'Home and office furniture'),
('Groceries', 'Food and beverages'),
('Beauty', 'Beauty and personal care products'),
('Sports', 'Sports equipment and accessories'),
('Automotive', 'Automotive parts and accessories');

INSERT INTO products (name, description, price, quantity, category, is_available) VALUES
('Smartphone', 'Latest model smartphone', 699, 50, 1, TRUE),
('Jeans', 'Denim jeans', 49, 100, 2, TRUE),
('Cookbook', 'Recipe book for various cuisines', 25, 200, 3, TRUE),
('Microwave', 'Compact microwave oven', 120, 30, 4, TRUE),
('Action Figure', 'Superhero action figure', 20, 150, 5, TRUE),
('Office Chair', 'Ergonomic office chair', 150, 40, 6, TRUE),
('Organic Pasta', 'Whole grain organic pasta', 5, 300, 7, TRUE),
('Shampoo', 'Herbal shampoo', 10, 250, 8, TRUE),
('Tennis Racket', 'Professional tennis racket', 80, 60, 9, TRUE),
('Car Battery', '12V car battery', 100, 20, 10, TRUE);
