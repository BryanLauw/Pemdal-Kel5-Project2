// ========================================
// Neo4j Browser Queries - Copy & Paste
// ========================================
// Buka Neo4j Browser di http://localhost:7474
// Login: neo4j / password123
// Copy-paste query di bawah satu per satu


// ========================================
// INSERT 1: Menambahkan Pengguna Baru
// ========================================
CREATE (p:Pasien {
    email: 'andi@example.com',
    kata_sandi: 'hashed_password',
    nama_lengkap: 'Andi Setiawan',
    tanggal_lahir: '1995-04-21',
    nomor_telepon: '08123456789',
    provinsi: 'Jawa Barat',
    kota: 'Bandung',
    jalan: 'Jl. Merdeka 123'
})
RETURN p.email AS email, p.nama_lengkap AS nama;


// ========================================
// INSERT 2: Cari Pasien Berdasarkan Nama
// ========================================
MATCH (p:Pasien {nama_lengkap: 'Andi Setiawan'})
WITH p LIMIT 1
RETURN p.email AS email, p.nama_lengkap AS nama, p.kota AS kota;


// ========================================
// INSERT 3: Menambah Rumah Sakit Baru
// ========================================
CREATE (r:RumahSakit {
    id_rs: 'RS999',
    email: 'rs@example.com',
    nama_rumah_sakit: 'RS Sehat Selalu',
    no_telepon: '0221234567',
    provinsi: 'Jawa Barat',
    kota: 'Bandung',
    jalan: 'Jl. Kesehatan 10'
})
RETURN r.id_rs AS id_rs, r.nama_rumah_sakit AS nama;


// ========================================
// INSERT 4: Menambah Departemen di RS
// ========================================
MATCH (rs:RumahSakit {nama_rumah_sakit: 'RS Sehat Selalu'})
WITH rs LIMIT 1
CREATE (d:Departemen {
    nama_departemen: 'Kardiologi',
    gedung: 'Gedung A'
})
CREATE (rs)-[:memiliki_departemen]->(d)
RETURN d.nama_departemen AS departemen, rs.nama_rumah_sakit AS rumah_sakit, rs.id_rs AS id_rs;


// ========================================
// SPECIAL GRAPH: Cari Dokter Spesialis
// Keunggulan Graph: Relationship Traversal
// ========================================
MATCH (tm:TenagaMedis {profesi: 'Dokter Spesialis Anak'})-[:bekerja_di]->(d:Departemen)
      <-[:memiliki_departemen]-(rs:RumahSakit {kota: 'Bandung'})
RETURN 
    tm.nama_lengkap AS nama_dokter,
    tm.nomor_telepon AS telepon,
    d.nama_departemen AS departemen,
    rs.nama_rumah_sakit AS rumah_sakit,
    rs.jalan AS alamat
ORDER BY rs.nama_rumah_sakit, tm.nama_lengkap
LIMIT 50;


// ========================================
// BONUS: Visualisasi Graph
// ========================================

// Lihat semua Pasien (limit 25)
MATCH (p:Pasien)
RETURN p
LIMIT 25;

// Lihat Rumah Sakit dan Departemennya
MATCH (rs:RumahSakit)-[:memiliki_departemen]->(d:Departemen)
RETURN rs, d
LIMIT 50;

// Lihat Tenaga Medis dan tempat kerja mereka
MATCH (tm:TenagaMedis)-[:bekerja_di]->(d:Departemen)<-[:memiliki_departemen]-(rs:RumahSakit)
RETURN tm, d, rs
LIMIT 30;

// Lihat Pasien dan perangkat Baymin mereka
MATCH (p:Pasien)-[:memiliki_perangkat]->(b:Baymin)
RETURN p, b
LIMIT 25;


// ========================================
// ANALISIS QUERIES
// ========================================

// Count semua nodes per label
MATCH (n)
RETURN labels(n) AS label, count(*) AS total
ORDER BY total DESC;

// Count semua relationships per type
MATCH ()-[r]->()
RETURN type(r) AS relationship, count(*) AS total
ORDER BY total DESC;

// Cari Rumah Sakit dengan departemen terbanyak
MATCH (rs:RumahSakit)-[:memiliki_departemen]->(d:Departemen)
RETURN rs.nama_rumah_sakit AS rumah_sakit, rs.kota AS kota, count(d) AS jumlah_departemen
ORDER BY jumlah_departemen DESC
LIMIT 10;

// Cari profesi dengan jumlah tenaga medis terbanyak
MATCH (tm:TenagaMedis)
RETURN tm.profesi AS profesi, count(*) AS jumlah
ORDER BY jumlah DESC;


// ========================================
// CLEANUP (Hati-hati!)
// ========================================

// Hapus data test yang baru dibuat
MATCH (p:Pasien {email: 'andi@example.com'})
DETACH DELETE p;

MATCH (rs:RumahSakit {id_rs: 'RS999'})
DETACH DELETE rs;

// Hapus semua departemen yang tidak terhubung ke RS
MATCH (d:Departemen)
WHERE NOT (d)<-[:memiliki_departemen]-()
DELETE d;


// ========================================
// TIPS
// ========================================
// 1. Timing otomatis muncul di bawah hasil query
// 2. Klik node di visualization untuk lihat properties
// 3. Gunakan PROFILE untuk detailed performance metrics
// 4. Gunakan EXPLAIN untuk lihat execution plan tanpa run query
// 5. Double-click node untuk expand relationships
