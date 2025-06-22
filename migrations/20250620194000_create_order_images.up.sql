CREATE TABLE order_images (
    id INT AUTO_INCREMENT PRIMARY KEY,
    order_id INT NOT NULL,
    image_url TEXT NOT NULL,
    FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE
); 