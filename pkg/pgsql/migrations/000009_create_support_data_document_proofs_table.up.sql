CREATE TABLE IF NOT EXISTS support_data_document_proofs(
  id uuid PRIMARY KEY,
  indicator_assessment_id uuid NOT NULL REFERENCES indicator_assessments(id),
  support_data_document_id uuid NOT NULL REFERENCES support_data_documents(id),
  proof TEXT NOT NULL,
  image_url TEXT,
  specific_page_document_url TEXT,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE
);

INSERT INTO support_data_document_proofs (id, indicator_assessment_id, support_data_document_id, proof, image_url, specific_page_document_url, created_at)
VALUES
  ('fcd52961-fa4e-43ba-a6df-a4c97849d819', 'ecd52961-fa4e-43ba-a6df-a4c97849d899', 'fcd52961-fa4e-43ba-a6df-a4c97849d899', '<p>Berdasarkan <b>rancangan arsitektur SPBE</b></p>', 'https://www.image.com', 'https://www.document.com', '2019-01-23T12:54:18.610Z'),
  ('fcd52961-fa4e-43ba-a6df-a4c97849d818', 'ecd52961-fa4e-43ba-a6df-a4c97849d898', 'fcd52961-fa4e-43ba-a6df-a4c97849d899', '<p>Mengingat <b>peraturan</b> yang telah diberikan</p>', 'https://www.image.com', 'https://www.document.com', '2019-01-23T12:54:18.610Z'),
  ('fcd52961-fa4e-43ba-a6df-a4c97849d817', 'ecd52961-fa4e-43ba-a6df-a4c97849d897', 'fcd52961-fa4e-43ba-a6df-a4c97849d899', '<p>Menimbang <b>keputusan</b> yang diberlakukan<p>', 'https://www.image.com', 'https://www.document.com', '2019-01-23T12:54:18.610Z');
