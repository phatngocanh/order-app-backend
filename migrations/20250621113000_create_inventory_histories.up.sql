CREATE TABLE inventory_histories (
    id INT AUTO_INCREMENT PRIMARY KEY,
    product_id INT NOT NULL,
    quantity INT NOT NULL COMMENT 'Số lượng nhập',
    importer_name VARCHAR(255) NOT NULL COMMENT 'Tên người nhập',
    imported_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT 'Thời gian nhập',
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
); 