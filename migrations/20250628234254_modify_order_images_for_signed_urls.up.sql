-- Modify order_images table to support signed URLs
-- Add s3_key column to store the S3 object key instead of public URL
ALTER TABLE order_images ADD COLUMN s3_key VARCHAR(500) NOT NULL DEFAULT '';