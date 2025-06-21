CREATE TABLE products (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    spec INT, -- Quy cách
    original_price INT NOT NULL DEFAULT 0 COMMENT 'Giá gốc của sản phẩm (VND)'
);
