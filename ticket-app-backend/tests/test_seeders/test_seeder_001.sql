INSERT INTO events (id, name, room, price)
VALUES (1, 'Spider-man', 'Hall 1', 150)
    ON CONFLICT (id) DO NOTHING;

INSERT INTO seats (id, event_id, row_number, line_number, status)
VALUES (2, 1, 1, 10, 'available')
    ON CONFLICT (id) DO NOTHING;

INSERT INTO users (id, username, email, password_hash)
SELECT
    i,
    'user_' || i,
    'user_' || i || '@test.com',
    'test_hash_123456'
FROM generate_series(1001, 1050) AS i
    ON CONFLICT (id) DO NOTHING;