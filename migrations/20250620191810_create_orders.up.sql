CREATE TABLE orders (
    id INT AUTO_INCREMENT PRIMARY KEY,
    customer_id INT NOT NULL,
    order_date DATE NOT NULL,
    delivery_status VARCHAR(20) CHECK (delivery_status IN ('CHUA_GIAO', 'DA_GIAO')),
    debt_status VARCHAR(100),
    FOREIGN KEY (customer_id) REFERENCES customers(id)
);
