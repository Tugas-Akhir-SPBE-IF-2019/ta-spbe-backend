CREATE TABLE IF NOT EXISTS user_evaluation_data(
  id uuid PRIMARY KEY,
  user_id uuid NOT NULL REFERENCES users(id),
  role TEXT NOT NULL,
  institution_id INT NOT NULL REFERENCES institution(id),
  evaluation_year INT NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE
);

INSERT INTO user_evaluation_data (id, user_id, role, institution_id, evaluation_year, created_at)
VALUES
  ('ccd52961-fa4e-43ba-a6df-a4c97849d898', 'ccd52961-fa4e-43ba-a6df-a4c97849d899', 'Evaluator Eksternal', 1, 2022, '2019-01-23T12:54:18.610Z'),
  ('ccd52961-fa4e-43ba-a6df-a4c97849d897', 'ccd52961-fa4e-43ba-a6df-a4c97849d899', 'Evaluator Eksternal', 2, 2022, '2019-01-23T12:54:18.610Z');
