CREATE TABLE IF NOT EXISTS indicators(
  id uuid PRIMARY KEY,
  indicator_number INT NOT NULL,
  aspect TEXT NOT NULL,
  domain TEXT NOT NULL,
  detail TEXT NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE
);

INSERT INTO indicators (id, indicator_number, aspect, domain, detail, created_at)
VALUES
  ('dcd52961-fa4e-43ba-a6df-a4c97849d899', 1, 'Kebijakan Internal Tata Kelola SPBE', 'Kebijakan Internal SPBE', 'Tingkat Kematangan Kebijakan Internal Arsitektur SPBE Instansi Pusat/Pemerintah Daerah', '2019-01-23T12:54:18.610Z'),
  ('dcd52961-fa4e-43ba-a6df-a4c97849d898', 2, 'Kebijakan Internal Tata Kelola SPBE', 'Kebijakan Internal SPBE', 'Tingkat Kematangan Kebijakan Internal Peta Rencana SPBE Instansi Pusat/Pemerintah Daerah', '2019-01-23T12:54:18.610Z'),
  ('dcd52961-fa4e-43ba-a6df-a4c97849d897', 3, 'Kebijakan Internal Tata Kelola SPBE', 'Kebijakan Internal SPBE', 'Tingkat Kematangan Kebijakan Internal Manajemen Data di Instansi Pusat/Pemerintah Daerah', '2019-01-23T12:54:18.610Z'),
  ('dcd52961-fa4e-43ba-a6df-a4c97849d896', 4, 'Kebijakan Internal Tata Kelola SPBE', 'Kebijakan Internal SPBE', 'Tingkat Kematangan Kebijakan Internal Pembangunan Aplikasi SPBE', '2019-01-23T12:54:18.610Z'),
  ('dcd52961-fa4e-43ba-a6df-a4c97849d895', 5, 'Kebijakan Internal Tata Kelola SPBE', 'Kebijakan Internal SPBE', 'Tingkat Kematangan Kebijakan Internal Layanan Pusat Data', '2019-01-23T12:54:18.610Z'),
  ('dcd52961-fa4e-43ba-a6df-a4c97849d894', 6, 'Kebijakan Internal Tata Kelola SPBE', 'Kebijakan Internal SPBE', 'Tingkat Kematangan Kebijakan Internal Layanan Jaringan Intra Instansi Pusat/Pemerintah Daerah', '2019-01-23T12:54:18.610Z'),
  ('dcd52961-fa4e-43ba-a6df-a4c97849d893', 7, 'Kebijakan Internal Tata Kelola SPBE', 'Kebijakan Internal SPBE', 'Tingkat Kematangan Kebijakan Internal Penggunaan Sistem Penghubung Layanan Instansi Pusat/Pemerintah Daerah', '2019-01-23T12:54:18.610Z'),
  ('dcd52961-fa4e-43ba-a6df-a4c97849d892', 8, 'Kebijakan Internal Tata Kelola SPBE', 'Kebijakan Internal SPBE', 'Tingkat Kematangan Kebijakan Internal Manajemen Keamanan Informasi. Evidence kebijakan internal terkait Manajemen Keamanan Informasi', '2019-01-23T12:54:18.610Z'),
  ('dcd52961-fa4e-43ba-a6df-a4c97849d891', 9, 'Kebijakan Internal Tata Kelola SPBE', 'Kebijakan Internal SPBE', 'Kebijakan internal Audit TIK', '2019-01-23T12:54:18.610Z'),
  ('dcd52961-fa4e-43ba-a6df-a4c97849d890', 10, 'Kebijakan Internal Tata Kelola SPBE', 'Kebijakan Internal SPBE', 'Tingkat Kematangan Kebijakan Internal Tim Koordinasi SPBE Instansi Pusat/Pemerintah Daerah', '2019-01-23T12:54:18.610Z');
