ALTER TABLE orders
ADD COLUMN status_transitioned_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT 'Thời gian trạng thái đơn hàng được chuyển đổi lần cuối';