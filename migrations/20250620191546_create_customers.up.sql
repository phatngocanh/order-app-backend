CREATE TABLE customers (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    phone VARCHAR(20),
    address varchar(255),
    location_type VARCHAR(20) CHECK (location_type IN ('TINH', 'THANH_PHO'))
);
