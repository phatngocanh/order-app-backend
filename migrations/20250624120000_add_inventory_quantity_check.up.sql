-- Add CHECK constraint to ensure inventory quantity is never negative
ALTER TABLE inventory ADD CONSTRAINT check_quantity_non_negative CHECK (quantity >= 0); 