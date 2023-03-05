CREATE TABLE IF NOT EXISTS assessment_status_histories(
  id uuid PRIMARY KEY,
  assessment_id uuid NOT NULL REFERENCES assessments(id),
  status INT NOT NULL REFERENCES assessment_statuses(id),
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE
);

INSERT INTO assessment_status_histories (id, assessment_id, status, created_at)
VALUES
  ('bcd52961-fa4e-43ba-a6df-a4c97849d899', 'acd52961-fa4e-43ba-a6df-a4c97849d899', 2, '2019-01-23T12:54:18.610Z'),
  ('bcd52961-fa4e-43ba-a6df-a4c97849d898', 'acd52961-fa4e-43ba-a6df-a4c97849d898', 3, '2019-01-23T12:54:18.610Z');
