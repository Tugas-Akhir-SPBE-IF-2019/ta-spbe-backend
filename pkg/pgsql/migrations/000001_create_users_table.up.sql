CREATE TABLE IF NOT EXISTS users(
  id uuid PRIMARY KEY,
  email TEXT NOT NULL UNIQUE,
  name TEXT NOT NULL,
  role TEXT NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE
);

INSERT INTO users (id, email, name, role, created_at)
VALUES
  ('ccd52961-fa4e-43ba-a6df-a4c97849d899', '13519142@std.stei.itb.ac.id', 'test user 1', 'user', '2019-01-23T12:54:18.610Z'),
  ('ccd52961-fa4e-43ba-a6df-a4c97849d898', 'testuser2@mail.com', 'test user 2', 'guest', '2019-01-23T12:54:18.610Z'),
  ('ccd52961-fa4e-43ba-a6df-a4c97849d897', 'testuser3@mail.com', 'test user 3', 'admin', '2019-01-23T12:54:18.610Z');
