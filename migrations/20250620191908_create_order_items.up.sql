CREATE TABLE order_items (
    id INT AUTO_INCREMENT PRIMARY KEY,
    order_id INT NOT NULL,
    product_id INT NOT NULL,
    number_of_boxes INT DEFAULT NULL COMMENT 'Số thùng',
    spec INT DEFAULT NULL COMMENT 'Quy cách mỗi thùng',
    quantity INT NOT NULL COMMENT 'Số lượng cuối cùng (có thể tính từ số thùng * quy cách hoặc nhập trực tiếp)',
    selling_price INT NOT NULL COMMENT 'Giá bán của sản phẩm (VND)',
    discount INT DEFAULT 0 COMMENT 'Chiết khấu (%)',
    final_amount INT COMMENT 'Số tiền cuối cùng sau khi trừ chiết khấu (VND)',
    FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES products(id)
);