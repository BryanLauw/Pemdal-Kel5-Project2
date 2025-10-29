-- ========================================
-- Cassandra CQL Queries - Copy & Paste
-- ========================================
-- Buka terminal dan jalankan: docker exec -it cassandra cqlsh
-- Atau gunakan DataStax DevCenter / DBeaver
-- Copy-paste query di bawah satu per satu


-- ========================================
-- UPDATE 1: Update Status Pesanan Expired
-- ========================================
-- Note: Cassandra TIDAK SUPPORT complex WHERE dengan time comparison
-- Harus dilakukan dalam 2 langkah:

-- Step 1: SELECT data yang akan diupdate (dengan ALLOW FILTERING)
SELECT id_pesanan, waktu_pemesanan, status_pemesanan
FROM rumahsakit.pemesanan_obat
WHERE status_pemesanan = 'belum dibayar'
ALLOW FILTERING;

-- Step 2: UPDATE satu per satu (ganti PESANAN_ID dengan ID dari hasil SELECT)
-- UPDATE rumahsakit.pemesanan_obat
-- SET status_pemesanan = 'dibatalkan'
-- WHERE id_pesanan = 'PESANAN_ID';

-- Contoh jika ada pesanan dengan id 'PO001':
-- UPDATE rumahsakit.pemesanan_obat
-- SET status_pemesanan = 'dibatalkan'
-- WHERE id_pesanan = 'PO001';


-- ========================================
-- LIMITATION CASSANDRA
-- ========================================
-- Query SQL yang TIDAK BISA dilakukan di Cassandra:
-- UPDATE pemesanan_obat
-- SET status_pemesanan = 'dibatalkan'
-- WHERE status_pemesanan = 'belum dibayar'
--   AND waktu_pemesanan < NOW() - INTERVAL 2 DAY;
--
-- Alasan:
-- 1. Cassandra tidak support time-based filtering di WHERE
-- 2. Cassandra tidak support multiple conditions tanpa partition key
-- 3. Harus fetch data dulu, filter di aplikasi, lalu update satu per satu
-- 4. Ini adalah trade-off untuk mendapatkan high write performance


-- ========================================
-- QUERY LAINNYA (Bonus)
-- ========================================

-- Lihat semua obat
SELECT * FROM rumahsakit.obat LIMIT 10;

-- Lihat obat berdasarkan ID
SELECT * FROM rumahsakit.obat WHERE id_obat = 'O0001';

-- Lihat log aktivitas dari perangkat tertentu
SELECT * FROM rumahsakit.log_aktivitas 
WHERE id_perangkat = 'BAYMIN-0001'
LIMIT 10;

-- Lihat layanan di rumah sakit tertentu
SELECT * FROM rumahsakit.lokasi_layanan
WHERE id_rs = 'RS001';

-- Count total obat (SLOW - full table scan)
SELECT COUNT(*) FROM rumahsakit.obat;


-- ========================================
-- INSERT EXAMPLES (Cassandra Style)
-- ========================================

-- Insert obat baru
INSERT INTO rumahsakit.obat (id_obat, nama, label, harga, stok)
VALUES ('O9999', 'Paracetamol 500mg', 'Pereda Nyeri', 15000, 200);

-- Insert log aktivitas
INSERT INTO rumahsakit.log_aktivitas (id_perangkat, waktu_aktivitas, detail_aktivitas)
VALUES ('BAYMIN-0001', toTimestamp(now()), 'Detak jantung: 72 bpm');

-- Insert pemesanan layanan
INSERT INTO rumahsakit.pemesanan_layanan (id_pesanan, email_pemesan, waktu_pemesanan, jadwal_pelaksanaan, status_pemesanan)
VALUES ('PL001', 'andi@example.com', toTimestamp(now()), '2025-11-01 10:00:00', 'terkonfirmasi');


-- ========================================
-- UPDATE EXAMPLES
-- ========================================

-- Update stok obat
UPDATE rumahsakit.obat
SET stok = 150
WHERE id_obat = 'O9999';

-- Update status pesanan (harus pake partition key)
UPDATE rumahsakit.pemesanan_obat
SET status_pemesanan = 'dikirim'
WHERE id_pesanan = 'PO001';


-- ========================================
-- DELETE EXAMPLES
-- ========================================

-- Delete obat
DELETE FROM rumahsakit.obat
WHERE id_obat = 'O9999';

-- Delete log aktivitas tertentu
DELETE FROM rumahsakit.log_aktivitas
WHERE id_perangkat = 'BAYMIN-0001'
AND waktu_aktivitas = '2025-10-29 10:30:00';


-- ========================================
-- SCHEMA QUERIES
-- ========================================

-- Lihat semua keyspaces
DESCRIBE KEYSPACES;

-- Lihat struktur keyspace rumahsakit
DESCRIBE KEYSPACE rumahsakit;

-- Lihat semua tabel di keyspace
DESCRIBE TABLES;

-- Lihat struktur tabel obat
DESCRIBE TABLE rumahsakit.obat;

-- Lihat struktur tabel pemesanan_obat
DESCRIBE TABLE rumahsakit.pemesanan_obat;


-- ========================================
-- ANALISIS DATA
-- ========================================

-- Count obat per label (ALLOW FILTERING - SLOW!)
SELECT label, COUNT(*) 
FROM rumahsakit.obat 
GROUP BY label
ALLOW FILTERING;

-- Lihat pemesanan terbaru (time-series query - FAST!)
SELECT * FROM rumahsakit.log_aktivitas
WHERE id_perangkat = 'BAYMIN-0001'
ORDER BY waktu_aktivitas DESC
LIMIT 20;


-- ========================================
-- CLEANUP (Hati-hati!)
-- ========================================

-- Truncate semua data di tabel (DANGER!)
-- TRUNCATE rumahsakit.obat;
-- TRUNCATE rumahsakit.pemesanan_obat;

-- Drop tabel (DANGER!)
-- DROP TABLE rumahsakit.obat;


-- ========================================
-- TIPS
-- ========================================
-- 1. Cassandra TIDAK menampilkan execution time di cqlsh
-- 2. Untuk benchmark, gunakan file Go (update1.go) yang sudah ada timing
-- 3. ALLOW FILTERING sangat lambat untuk dataset besar
-- 4. Selalu query berdasarkan partition key untuk performa optimal
-- 5. Time-series queries (ORDER BY clustering key) sangat cepat
-- 6. Avoid COUNT(*) pada tabel besar
-- 7. Cassandra dioptimalkan untuk WRITE, bukan complex READ
