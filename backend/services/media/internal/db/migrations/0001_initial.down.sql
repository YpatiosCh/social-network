
-- Drop trigger
DROP TRIGGER IF EXISTS update_images_updated_at ON images;

-- Drop trigger function
DROP FUNCTION IF EXISTS update_timestamp();

-- Drop table
DROP TABLE IF EXISTS images;
