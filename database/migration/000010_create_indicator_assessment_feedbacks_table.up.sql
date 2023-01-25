CREATE TABLE IF NOT EXISTS indicator_assessment_feedbacks(
  id uuid PRIMARY KEY,
  indicator_assessment_id uuid NOT NULL REFERENCES indicator_assessments(id),
  support_data_document_proof_id uuid REFERENCES support_data_document_proofs(id),
  level INT,
  feedback TEXT,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE
);

INSERT INTO indicator_assessment_feedbacks (id, indicator_assessment_id, support_data_document_proof_id, level, feedback, created_at)
VALUES
  ('fcd52961-fa4e-43ba-a6df-a4c97849d829', 'ecd52961-fa4e-43ba-a6df-a4c97849d899', 'fcd52961-fa4e-43ba-a6df-a4c97849d819', 5, 'Seharusnya level 5 berdasarkan data yang ada', '2019-01-23T12:54:18.610Z'),
  ('fcd52961-fa4e-43ba-a6df-a4c97849d828', 'ecd52961-fa4e-43ba-a6df-a4c97849d898', 'fcd52961-fa4e-43ba-a6df-a4c97849d818', 4, 'Level 4 adalah level yang paling sesuai', '2019-01-23T12:54:18.610Z'),
  ('fcd52961-fa4e-43ba-a6df-a4c97849d827', 'ecd52961-fa4e-43ba-a6df-a4c97849d897', 'fcd52961-fa4e-43ba-a6df-a4c97849d817', 5, 'Sudah pasti harusnya level 5', '2019-01-23T12:54:18.610Z');
