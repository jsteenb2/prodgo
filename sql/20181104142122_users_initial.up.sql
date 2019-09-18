CREATE TABLE IF NOT EXISTS users (
  id VARCHAR(255) NOT NULL,
  email VARCHAR(255) UNIQUE NOT NULL,
  first_name VARCHAR(255) NOT NULL,
  last_name VARCHAR(255) NOT NULL,
  password VARCHAR(255) NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  PRIMARY KEY(id)
);

CREATE INDEX IF NOT EXISTS first_name_index on users(first_name);
CREATE INDEX IF NOT EXISTS last_name_index on users(first_name);
CREATE INDEX IF NOT EXISTS email_index on users(email);

