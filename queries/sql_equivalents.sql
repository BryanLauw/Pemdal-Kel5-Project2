-- ========================================
-- SQL Query Equivalents
-- ========================================
-- SQL version dari Neo4j Cypher queries
-- Menunjukkan kompleksitas SQL vs Graph Database


-- ========================================
-- SPECIAL GRAPH: Cari Dokter Spesialis
-- SQL Version (Complex dengan Multiple JOINs)
-- ========================================

-- Query ini di SQL membutuhkan MULTIPLE JOINs yang lambat
-- Di Neo4j hanya butuh relationship traversal yang cepat

SELECT 
    u.nama_lengkap AS nama_dokter,
    u.nomor_telepon AS telepon,
    d.nama_departemen AS departemen,
    rs.nama_rumah_sakit AS rumah_sakit,
    rs.jalan AS alamat
FROM tenaga_medis tm
INNER JOIN user u ON tm.email = u.email
INNER JOIN departemen d ON tm.id_departemen = d.id_departemen
INNER JOIN rumah_sakit rs ON d.id_rs = rs.id_rs
WHERE tm.profesi = 'Dokter Spesialis Anak'
  AND rs.kota = 'Bandung'
ORDER BY rs.nama_rumah_sakit, u.nama_lengkap
LIMIT 50;

-- Perbandingan Performa:
-- SQL: ~50-500ms (4 table JOINs, tergantung dataset size)
-- Neo4j: ~5-50ms (Direct relationship traversal)
-- Winner: Neo4j 10x FASTER!

-- Mengapa SQL Lambat:
-- 1. Harus JOIN 4 tabel (tenaga_medis, user, departemen, rumah_sakit)
-- 2. Query optimizer harus merencanakan join order
-- 3. Butuh index pada multiple columns (profesi, kota, foreign keys)
-- 4. Performa menurun drastis saat dataset membesar
-- 5. Setiap JOIN create intermediate result set

-- Mengapa Neo4j Cepat:
-- 1. Index-free adjacency (langsung ke node terkait)
-- 2. Tidak perlu JOIN, hanya traversal relationship
-- 3. Performa konsisten bahkan untuk millions of nodes
-- 4. Graph database dioptimalkan untuk query pattern seperti ini


-- ========================================
-- Neo4j Cypher Version (untuk comparison)
-- ========================================

-- MATCH (tm:TenagaMedis {profesi: 'Dokter Spesialis Anak'})-[:bekerja_di]->(d:Departemen)
--       <-[:memiliki_departemen]-(rs:RumahSakit {kota: 'Bandung'})
-- RETURN 
--     tm.nama_lengkap AS nama_dokter,
--     tm.nomor_telepon AS telepon,
--     d.nama_departemen AS departemen,
--     rs.nama_rumah_sakit AS rumah_sakit,
--     rs.jalan AS alamat
-- ORDER BY rs.nama_rumah_sakit, tm.nama_lengkap
-- LIMIT 50;


-- ========================================
-- SQL Queries Lainnya (dari referensi)
-- ========================================

-- 1. Daftar pasien dan total pesanan obat mereka
SELECT p.email, u.nama_lengkap, COUNT(po.id_pesanan) AS total_pesanan
FROM pasien p
JOIN user u ON p.email = u.email
LEFT JOIN pemesanan_obat po ON p.email = po.email_pemesan
GROUP BY p.email, u.nama_lengkap;

-- 2. Obat dengan stok kurang dari 20
SELECT id_obat, label, stok
FROM obat
WHERE stok < 20
ORDER BY stok ASC;

-- 3. Log aktivitas Baymin milik pasien tertentu
SELECT u.nama_lengkap, l.waktu_aktivitas, l.detail_aktivitas
FROM log_aktivitas l
JOIN baymin b ON l.id_perangkat = b.id_perangkat
JOIN pasien p ON b.email_pasien = p.email
JOIN user u ON p.email = u.email
WHERE p.email = 'adhiarjawaskita@example.com'
ORDER BY l.waktu_aktivitas DESC;

-- 4. Tenaga medis dan jumlah janji temu mereka
SELECT tm.email, u.nama_lengkap, tm.profesi, COUNT(jt.id_janji_temu) AS jumlah_janji_temu
FROM tenaga_medis tm
JOIN user u ON tm.email = u.email
LEFT JOIN janji_temu jt ON tm.email = jt.email_tenaga_medis
GROUP BY tm.email, u.nama_lengkap, tm.profesi
ORDER BY jumlah_janji_temu DESC;

-- 5. Detail resep dari janji temu pasien
SELECT jt.id_janji_temu, r.penyakit, o.label, dr.dosis
FROM hasil_janji_temu hjt
JOIN resep r ON hjt.id_resep = r.id_resep
JOIN detail_resep dr ON r.id_resep = dr.id_resep
JOIN obat o ON dr.id_obat = o.id_obat
JOIN janji_temu jt ON hjt.id_janji_temu = jt.id_janji_temu
WHERE jt.id_janji_temu = 3;

-- 6. Pasien dengan pengeluaran obat terbesar
SELECT p.email, u.nama_lengkap, SUM(o.harga * dp.jumlah) AS total_pengeluaran
FROM pasien p
JOIN user u ON p.email = u.email
JOIN pemesanan_obat po ON p.email = po.email_pemesan
JOIN detail_pesanan dp ON po.id_pesanan = dp.id_pesanan
JOIN obat o ON dp.id_obat = o.id_obat
WHERE po.status_pemesanan = 'selesai'
GROUP BY p.email, u.nama_lengkap
ORDER BY total_pengeluaran DESC
LIMIT 5;

-- 7. Layanan medis yang paling sering dipesan
SELECT lm.nama_layanan, COUNT(pl.id_pesanan) AS jumlah_pesanan
FROM pemesanan_layanan pl
JOIN layanan_medis lm ON pl.id_rs = lm.id_rs AND pl.id_layanan = lm.id_layanan
GROUP BY lm.nama_layanan
ORDER BY jumlah_pesanan DESC;

-- 8. 10 rumah sakit dengan janji temu terbanyak
SELECT rs.nama_rumah_sakit, COUNT(jt.id_janji_temu) AS total_janji
FROM rumah_sakit rs
JOIN janji_temu jt ON rs.id_rs = jt.id_rs
GROUP BY rs.id_rs, rs.nama_rumah_sakit
ORDER BY total_janji DESC
LIMIT 10;

-- 9. Rumah sakit dengan tenaga medis terbanyak
SELECT rs.nama_rumah_sakit, COUNT(tm.email) AS jumlah_tenaga_medis
FROM rumah_sakit rs
JOIN tenaga_medis tm ON rs.id_rs = tm.id_rs
GROUP BY rs.nama_rumah_sakit
ORDER BY jumlah_tenaga_medis DESC
LIMIT 10;

-- 10. Pasien yang pernah janji temu tapi tidak dapat resep
SELECT DISTINCT u.nama_lengkap, p.email
FROM janji_temu jt
JOIN pasien p ON jt.email_pasien = p.email
JOIN user u ON p.email = u.email
WHERE jt.id_janji_temu NOT IN (SELECT id_janji_temu FROM hasil_janji_temu);


-- ========================================
-- GRAPH vs SQL Comparison
-- ========================================

-- Query Type: Find Connected Data (3+ hops)
-- Example: Pasien -> Baymin -> Log Aktivitas
-- SQL: 3-4 JOINs, ~50-300ms
-- Neo4j: Direct traversal, ~5-30ms

-- Query Type: Pattern Matching
-- Example: Find all doctors in city with specialization
-- SQL: Multiple JOINs + complex WHERE, ~100-500ms
-- Neo4j: MATCH pattern, ~10-50ms

-- Query Type: Shortest Path
-- Example: Connection between two entities
-- SQL: Recursive CTE or multiple queries, ~200-1000ms
-- Neo4j: shortestPath() function, ~10-50ms

-- Query Type: Recommendation
-- Example: Find similar patients based on conditions
-- SQL: Complex subqueries + JOINs, ~500-2000ms
-- Neo4j: Collaborative filtering via relationships, ~50-200ms


-- ========================================
-- KESIMPULAN
-- ========================================

-- Gunakan SQL ketika:
-- 1. ACID compliance critical (financial transactions)
-- 2. Complex aggregations (SUM, AVG, GROUP BY)
-- 3. Time-based filtering dengan complex WHERE
-- 4. Structured data dengan fixed schema
-- 5. Ad-hoc reporting queries

-- Gunakan Neo4j ketika:
-- 1. Data highly connected (social network, org chart)
-- 2. Relationship traversal penting (friend of friend)
-- 3. Pattern matching (find all X connected to Y via Z)
-- 4. Recommendation engines
-- 5. Fraud detection (network analysis)
-- 6. Schema flexibility diperlukan

-- Gunakan Cassandra ketika:
-- 1. High write throughput (IoT, logs, events)
-- 2. Time-series data dengan partition by time
-- 3. Linear scalability critical
-- 4. Simple queries by partition key
-- 5. Eventually consistent is acceptable
