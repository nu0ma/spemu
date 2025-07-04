-- Sample seed data for spemu
-- Users table
INSERT INTO users (id, name, email, created_at) VALUES
  (1, 'John Doe', 'john.doe@example.com', '2024-01-01T00:00:00Z'),
  (2, 'Jane Smith', 'jane.smith@example.com', '2024-01-02T00:00:00Z'),
  (3, 'Bob Wilson', 'bob.wilson@example.com', '2024-01-03T00:00:00Z');

-- Posts table
INSERT INTO posts (id, user_id, title, content, created_at) VALUES
  (1, 1, 'First Post', 'Hello World! Testing with Spanner emulator.', '2024-01-01T01:00:00Z'),
  (2, 2, 'About Spanner', 'Cloud Spanner is an amazing database service.', '2024-01-02T02:00:00Z'),
  (3, 1, 'Second Post', 'Development with emulator is very convenient.', '2024-01-03T03:00:00Z'),
  (4, 3, 'Hello Everyone', 'Nice to meet you all!', '2024-01-03T04:00:00Z');

-- Comments table
INSERT INTO comments (id, post_id, user_id, content, created_at) VALUES
  (1, 1, 2, 'Great post!', '2024-01-01T02:00:00Z'),
  (2, 1, 3, 'Nice to meet you!', '2024-01-01T03:00:00Z'),
  (3, 2, 1, 'Very informative, thanks!', '2024-01-02T03:00:00Z'),
  (4, 3, 2, 'Emulator is really handy.', '2024-01-03T05:00:00Z');