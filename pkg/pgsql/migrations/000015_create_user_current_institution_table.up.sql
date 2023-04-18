CREATE TABLE IF NOT EXISTS user_current_institutions(
  id uuid PRIMARY KEY,
  user_id uuid NOT NULL REFERENCES users(id),
  institution_id INT REFERENCES institution(id),
  institution_name TEXT NOT NULL,
  role TEXT NOT NULL,
  status TEXT NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE
);

INSERT INTO user_current_institutions(id, user_id, institution_id, institution_name, role, status)
VALUES
  ('e536b8f7-e287-4ec4-80be-73b984b0f861', 'ccd52961-fa4e-43ba-a6df-a4c97849d899', 1, 'Kabupaten Aceh Barat', 'Asesor Eksternal', 'VALID'),
  ('6b04a2b0-f5fa-4a20-abb3-29b2f4e2e3a7', 'ccd52961-fa4e-43ba-a6df-a4c97849d899', 2, 'Kabupaten Aceh Barat Daya', 'Asesor Eksternal', 'VALID')
