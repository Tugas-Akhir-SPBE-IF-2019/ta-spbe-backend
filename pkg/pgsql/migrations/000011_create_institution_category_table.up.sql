CREATE TABLE IF NOT EXISTS institution_category(
  id SERIAL PRIMARY KEY,
  category TEXT NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE
);

INSERT INTO institution_category (category)
VALUES
  ('NAD_ACEH'),
  ('SUMATERA_UTARA'),
  ('SUMATERA_BARAT'),
  ('SUMATERA_SELATAN'),
  ('RIAU'),
  ('KEPULAUAN_RIAU'),
  ('JAMBI'),
  ('BENGKULU'),
  ('BANGKA_BELITUNG'),
  ('LAMPUNG'),
  ('BANTEN'),
  ('JAWA_BARAT'),
  ('JAWA_TENGAH'),
  ('JAWA_TIMUR'),
  ('DKI_JAKARTA'),
  ('DAERAH_ISTIMEWA_YOGYAKARTA'),
  ('BALI'),
  ('NUSA_TENGGARA_BARAT'),
  ('NUSA_TENGGARA_TIMUR'),
  ('KALIMANTAN_BARAT'),
  ('KALIMANTAN_SELATAN'),
  ('KALIMANTAN_TENGAH'),
  ('KALIMANTAN_TIMUR'),
  ('KALIMANTAN_UTARA'),
  ('GORONTALO'),
  ('SULAWESI_SELATAN'),
  ('SULAWESI_TENGGARA'),
  ('SULAWESI_TENGAH'),
  ('SULAWESI_UTARA'),
  ('SULAWESI_BARAT'),
  ('MALUKU'),
  ('MALUKU_UTARA'),
  ('PAPUA'),
  ('PAPUA_BARAT'),
  ('INDUSTRI_TELEKOMUNIKASI_DAN_MEDIA'),
  ('INDUSTRI_ENERGI_MINYAK_DAN_GAS_BUMI'),
  ('INDUSTRI_MINERAL_DAN_BATUBARA'),
  ('INDUSTRI_PERKEBUNAN_DAN_KEHUTANAN'),
  ('INDUSTRI_PANGAN'),
  ('INDUSTRI_KESEHATAN'),
  ('INDUSTRI_MANUFAKTUR_DAN_SURVEI'),
  ('JASA_KEUANGAN'),
  ('JASA_ASURANSI_DAN_DANA_PENSIUN'),
  ('JASA_INFRASTRUKTUR'),
  ('JASA_LOGISTIK'),
  ('JASA_PARIWISATA_DAN_PENDUKUNG'),
  ('SUBKLASTER_DANAREKSA');
