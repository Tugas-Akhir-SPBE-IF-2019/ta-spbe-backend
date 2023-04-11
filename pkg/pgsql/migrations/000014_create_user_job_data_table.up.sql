CREATE TABLE IF NOT EXISTS user_job_data(
  id uuid PRIMARY KEY,
  user_id uuid NOT NULL REFERENCES users(id),
  role TEXT NOT NULL,
  company TEXT NOT NULL,
  joined_date INT NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE
);

INSERT INTO user_job_data (id, user_id, role, company, joined_date, created_at)
VALUES
  ('ccd52961-fa4e-43ba-a6df-a4c97849d898', 'ccd52961-fa4e-43ba-a6df-a4c97849d899', 'Pengajar Matematika', 'Institut Teknologi Bandung', 2017, '2019-01-23T12:54:18.610Z'),
  ('ccd52961-fa4e-43ba-a6df-a4c97849d897', 'ccd52961-fa4e-43ba-a6df-a4c97849d899', 'Pengajar Fisika', 'Universitas Terbuka', 2022, '2019-01-23T12:54:18.610Z');
