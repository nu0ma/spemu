-- Test schema for spemu integration tests
-- This file defines the database schema used in tests

CREATE TABLE users (
  id INT64 NOT NULL,
  name STRING(100) NOT NULL,
  email STRING(255) NOT NULL,
  created_at TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp=true)
) PRIMARY KEY (id);

CREATE TABLE posts (
  id INT64 NOT NULL,
  user_id INT64 NOT NULL,
  title STRING(200) NOT NULL,
  content STRING(MAX),
  created_at TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp=true),
  FOREIGN KEY (user_id) REFERENCES users (id)
) PRIMARY KEY (id);

CREATE TABLE comments (
  id INT64 NOT NULL,
  post_id INT64 NOT NULL,
  user_id INT64 NOT NULL,
  content STRING(1000) NOT NULL,
  created_at TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp=true),
  FOREIGN KEY (post_id) REFERENCES posts (id),
  FOREIGN KEY (user_id) REFERENCES users (id)
) PRIMARY KEY (id);

-- Additional test tables for integration tests
CREATE TABLE test_table (
  id INT64 NOT NULL,
  name STRING(100) NOT NULL,
  value STRING(255),
  created_at TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp=true)
) PRIMARY KEY (id);

CREATE TABLE benchmark_table (
  id INT64 NOT NULL,
  value STRING(100),
  created_at TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp=true)
) PRIMARY KEY (id);