CREATE TABLE inventory (
    id INT AUTO_INCREMENT PRIMARY KEY,
    product_id INT NOT NULL,
    quantity INT NOT NULL DEFAULT 0,
    version VARCHAR(36) NOT NULL COMMENT 'UUID version của inventory, thay đổi mỗi khi có thay đổi quantity để đảm bảo consistency khi FE tạo đơn hàng',
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE,
    UNIQUE KEY unique_product_inventory (product_id)
);