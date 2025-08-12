-- Drop trigger
DROP TRIGGER IF EXISTS update_quotes_updated_at ON quotes;

-- Drop indexes
DROP INDEX IF EXISTS idx_quotes_tags;
DROP INDEX IF EXISTS idx_quotes_author_id;

-- Drop table
DROP TABLE IF EXISTS quotes;
