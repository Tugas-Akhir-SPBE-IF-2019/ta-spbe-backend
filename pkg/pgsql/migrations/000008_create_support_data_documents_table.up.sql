CREATE TABLE IF NOT EXISTS support_data_documents(
  id uuid PRIMARY KEY,
  assessment_id uuid NOT NULL REFERENCES assessments(id),
  indicator_assessment_id uuid NOT NULL REFERENCES indicator_assessments(id),
  document_name TEXT NOT NULL,
  document_url TEXT NOT NULL,
  document_original_name TEXT NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE
);

INSERT INTO support_data_documents (id, assessment_id, indicator_assessment_id, document_name, document_url, document_original_name, created_at)
VALUES
  ('fcd52961-fa4e-43ba-a6df-a4c97849d899', 'acd52961-fa4e-43ba-a6df-a4c97849d899', 'ecd52961-fa4e-43ba-a6df-a4c97849d899', 'Dokumen Peraturan Walikota 2021.pdf','https://www.image.com', 'Dokumen Peraturan Walikota 2021.pdf', '2019-01-23T12:54:18.610Z'),
  ('fcd52961-fa4e-43ba-a6df-a4c97849d898', 'acd52961-fa4e-43ba-a6df-a4c97849d898', 'ecd52961-fa4e-43ba-a6df-a4c97849d898', 'Draft Revisi Bupati 2020.pdf','https://www.image.com', 'Dokumen Peraturan Walikota 2021.pdf', '2019-01-23T12:54:18.610Z'),
  ('fcd52961-fa4e-43ba-a6df-a4c97849d897', 'acd52961-fa4e-43ba-a6df-a4c97849d897', 'ecd52961-fa4e-43ba-a6df-a4c97849d897', 'Surat Keputusan Walikota.pdf','https://www.image.com', 'Dokumen Peraturan Walikota 2021.pdf', '2019-01-23T12:54:18.610Z');
