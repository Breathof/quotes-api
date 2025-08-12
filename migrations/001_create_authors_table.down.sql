-- Drop trigger
DROP TRIGGER IF EXISTS update_authors_updated_at ON authors;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop index
DROP INDEX IF EXISTS idx_authors_name;

-- Drop table
DROP TABLE IF EXISTS authors;
