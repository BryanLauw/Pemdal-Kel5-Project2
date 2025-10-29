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
// CLEANUP (Hati-hati!)
// ========================================

// ========================================
// HAPUS DATA INSERT 1-4 (untuk re-run queries)
// ========================================

// Hapus Pasien dengan email 'andi@example.com' (INSERT 1)
MATCH (p:Pasien {email: 'andi@example.com'})
DETACH DELETE p;

// Hapus Rumah Sakit 'RS Sehat Selalu' beserta departemennya (INSERT 3 & 4)
MATCH (rs:RumahSakit {id_rs: 'RS999'})
OPTIONAL MATCH (rs)-[:memiliki_departemen]->(d:Departemen)
DETACH DELETE rs, d;

// Atau hapus hanya departemen 'Kardiologi' di RS tertentu (INSERT 4)
MATCH (rs:RumahSakit {nama_rumah_sakit: 'RS Sehat Selalu'})-[:memiliki_departemen]->(d:Departemen {nama_departemen: 'Kardiologi'})
DETACH DELETE d;

// Hapus semua departemen yang tidak terhubung ke RS manapun
MATCH (d:Departemen)
WHERE NOT (d)<-[:memiliki_departemen]-()
DELETE d;

// ========================================
// QUICK CLEANUP - Hapus semua data INSERT 1-4 sekaligus
// ========================================
MATCH (p:Pasien {email: 'andi@example.com'})
DETACH DELETE p;

MATCH (rs:RumahSakit {id_rs: 'RS999'})
OPTIONAL MATCH (rs)-[r]->(d:Departemen)
DETACH DELETE rs, d;


// ========================================
// TIPS
// ========================================
// 1. Timing otomatis muncul di bawah hasil query
// 2. Klik node di visualization untuk lihat properties
// 3. Gunakan PROFILE untuk detailed performance metrics
// 4. Gunakan EXPLAIN untuk lihat execution plan tanpa run query
// 5. Double-click node untuk expand relationships
