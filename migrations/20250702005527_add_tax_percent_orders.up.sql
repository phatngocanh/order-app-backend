ALTER TABLE orders 
ADD COLUMN tax_percent INT DEFAULT 0 COMMENT 'Phần trăm thuế của đơn hàng (%)';
