-- Seed some initial data for development
INSERT INTO authors (name, bio) VALUES
    ('Albert Einstein', 'German-born theoretical physicist who developed the theory of relativity.'),
    ('Maya Angelou', 'American poet, memoirist, and civil rights activist.'),
    ('Mark Twain', 'American writer, humorist, entrepreneur, publisher, and lecturer.')
ON CONFLICT DO NOTHING;

-- Insert quotes only if authors exist
INSERT INTO quotes (content, author_id, source, tags)
SELECT 
    'Imagination is more important than knowledge.',
    a.id,
    'What Life Means to Einstein (1929)',
    ARRAY['imagination', 'knowledge', 'wisdom']
FROM authors a WHERE a.name = 'Albert Einstein'
ON CONFLICT DO NOTHING;

INSERT INTO quotes (content, author_id, source, tags)
SELECT 
    'I''ve learned that people will forget what you said, people will forget what you did, but people will never forget how you made them feel.',
    a.id,
    NULL,
    ARRAY['wisdom', 'feelings', 'impact']
FROM authors a WHERE a.name = 'Maya Angelou'
ON CONFLICT DO NOTHING;

INSERT INTO quotes (content, author_id, source, tags)
SELECT 
    'The two most important days in your life are the day you are born and the day you find out why.',
    a.id,
    NULL,
    ARRAY['life', 'purpose', 'wisdom']
FROM authors a WHERE a.name = 'Mark Twain'
ON CONFLICT DO NOTHING;
