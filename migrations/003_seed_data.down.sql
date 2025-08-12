-- Remove seeded quotes
DELETE FROM quotes WHERE author_id IN (
    SELECT id FROM authors WHERE name IN ('Albert Einstein', 'Maya Angelou', 'Mark Twain')
);

-- Remove seeded authors
DELETE FROM authors WHERE name IN ('Albert Einstein', 'Maya Angelou', 'Mark Twain');
