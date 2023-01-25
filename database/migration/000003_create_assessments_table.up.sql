CREATE TABLE IF NOT EXISTS assessments(
  id uuid PRIMARY KEY,
  user_id uuid NOT NULL REFERENCES users(id),
  status INT NOT NULL REFERENCES assessment_statuses(id),
  institution_name TEXT NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE
);

INSERT INTO assessments (id, user_id, status, institution_name, created_at)
VALUES
  ('acd52961-fa4e-43ba-a6df-a4c97849d899', 'ccd52961-fa4e-43ba-a6df-a4c97849d899', 1, 'Kabupaten Lamongan', '2019-01-23T12:54:18.610Z'),
  ('acd52961-fa4e-43ba-a6df-a4c97849d898', 'ccd52961-fa4e-43ba-a6df-a4c97849d898', 2, 'Kota Bandung', '2019-01-23T12:54:18.610Z'),
  ('acd52961-fa4e-43ba-a6df-a4c97849d897', 'ccd52961-fa4e-43ba-a6df-a4c97849d897', 3, 'RRI', '2019-01-23T12:54:18.610Z');
