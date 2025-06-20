CREATE TABLE order_images (
    id INT AUTO_INCREMENT PRIMARY KEY,
    order_id INT NOT NULL,
    image_url TEXT NOT NULL,
    image_type VARCHAR(50) DEFAULT 'delivery', -- delivery, receipt, other
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE
); 