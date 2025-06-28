ALTER TABLE order_items ADD COLUMN original_price INTEGER NOT NULL DEFAULT 0 AFTER selling_price;
