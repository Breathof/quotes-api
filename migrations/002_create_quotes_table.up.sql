-- Create quotes table
CREATE TABLE IF NOT EXISTS quotes (
    id BIGSERIAL PRIMARY KEY,
    content TEXT NOT NULL,
    author_id BIGINT NOT NULL,
    source VARCHAR(500),
    tags TEXT[],
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- Foreign key constraint
    CONSTRAINT fk_quotes_author
        FOREIGN KEY (author_id) 
        REFERENCES authors(id)
        ON DELETE RESTRICT
);

-- Create indexes
CREATE INDEX idx_quotes_author_id ON quotes(author_id);
CREATE INDEX idx_quotes_tags ON quotes USING GIN(tags);

-- Create updated_at trigger
CREATE TRIGGER update_quotes_updated_at BEFORE UPDATE
    ON quotes FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
