CREATE TABLE IF NOT EXISTS indicator_assessment_proofs(
  id uuid PRIMARY KEY,
  indicator_assessment_id uuid NOT NULL REFERENCES indicator_assessments(id),
  image_url TEXT NOT NULL,
  document_url TEXT NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE
);
