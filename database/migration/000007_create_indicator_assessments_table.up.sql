CREATE TABLE IF NOT EXISTS indicator_assessments(
  id uuid PRIMARY KEY,
  indicator_id uuid NOT NULL REFERENCES indicators(id),
  assessment_id uuid NOT NULL REFERENCES assessments(id),
  status INT NOT NULL REFERENCES indicator_assessment_statuses(id),
  level INT NOT NULL DEFAULT 0,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE
);

INSERT INTO indicator_assessments (id, indicator_id,assessment_id, status, level, created_at)
VALUES
  ('ecd52961-fa4e-43ba-a6df-a4c97849d899', 'dcd52961-fa4e-43ba-a6df-a4c97849d899','acd52961-fa4e-43ba-a6df-a4c97849d899', 2, 4,'2019-01-23T12:54:18.610Z'),
  ('ecd52961-fa4e-43ba-a6df-a4c97849d898', 'dcd52961-fa4e-43ba-a6df-a4c97849d898','acd52961-fa4e-43ba-a6df-a4c97849d899', 1, 0,'2019-01-23T12:54:18.610Z'),
  ('ecd52961-fa4e-43ba-a6df-a4c97849d897', 'dcd52961-fa4e-43ba-a6df-a4c97849d898','acd52961-fa4e-43ba-a6df-a4c97849d898', 3, 2,'2019-01-23T12:54:18.610Z');
