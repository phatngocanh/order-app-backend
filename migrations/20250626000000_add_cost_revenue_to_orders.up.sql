ALTER TABLE orders 
ADD COLUMN total_original_cost INT DEFAULT 0 COMMENT 'Tổng chi phí gốc của đơn hàng (VND)',
ADD COLUMN total_sales_revenue INT DEFAULT 0 COMMENT 'Tổng doanh thu bán hàng của đơn hàng (VND)'; 