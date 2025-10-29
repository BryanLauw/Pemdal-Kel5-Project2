# Pemdal-Kel5-Project2 - Sistem Informasi Rumah Sakit

Proyek ini adalah sistem informasi rumah sakit yang menggunakan **hybrid database** dengan **Cassandra** (untuk data transaksional dan time-series) dan **Neo4j** (untuk data relasional).

---

## Table of Contents
- [Requirements](#-requirements)
- [Instalasi](#-instalasi)
- [Setup & Konfigurasi](#-setup--konfigurasi)
- [Cara Menjalankan](#-cara-menjalankan)
- [Cara Membuat Query Baru](#-cara-membuat-query-baru)
- [Benchmark Queries](#-benchmark-queries)
- [Struktur Database](#-struktur-database)
- [Troubleshooting](#-troubleshooting)

---

## Requirements

Pastikan sistem kamu sudah terinstall:
- **Go** 1.25.3 atau lebih baru ([Download Go](https://go.dev/dl/))
- **Docker** & **Docker Compose** ([Download Docker](https://www.docker.com/products/docker-desktop))
- **Git** (untuk clone repository)

---

## Instalasi

### 1. Clone Repository
```powershell
git clone https://github.com/BryanLauw/Pemdal-Kel5-Project2.git
cd Pemdal-Kel5-Project2
```

### 2. Install Go Dependencies
```powershell
go mod download
```

Ini akan menginstall dependencies:
- `github.com/gocql/gocql` - Driver Cassandra
- `github.com/neo4j/neo4j-go-driver/v5` - Driver Neo4j
- `github.com/go-faker/faker/v4` - Library untuk generate data dummy

---

## Setup & Konfigurasi

### 1. Jalankan Docker Containers

Jalankan Cassandra dan Neo4j menggunakan Docker Compose:

```powershell
docker-compose up -d
```

Tunggu beberapa saat hingga container ready. Cek status dengan:

```powershell
docker-compose ps
```

Pastikan kedua service sudah `running`.

### 2. Verifikasi Koneksi Database

**PENTING:** Tunggu 2-3 menit setelah `docker-compose up -d` agar Cassandra fully initialized!

**Cek Status Container:**
```powershell
docker-compose ps
```

Tunggu hingga Cassandra status berubah dari `(health: starting)` menjadi `(healthy)`.

**Cassandra:**
```powershell
# Tunggu hingga Cassandra ready (cek setiap 30 detik)
docker exec -it cassandra nodetool status

# Jika output menampilkan "UN" (Up Normal), coba connect:
docker exec -it cassandra cqlsh
```

Jika berhasil, kamu akan masuk ke CQL shell. Ketik `exit` untuk keluar.

💡 **Jika masih error "Connection refused":**
- Tunggu 1-2 menit lagi, Cassandra masih initializing
- Cek logs: `docker logs cassandra`

**Neo4j:**

1. Buka browser: [http://localhost:7474](http://localhost:7474)

2. Di halaman connect, masukkan:
   - **Connect URL:** `bolt://localhost:7687` (atau `neo4j://localhost:7687`)
   - **Authentication type:** Username / Password
   - **Username:** `neo4j`
   - **Password:** `password123`

3. Klik **Connect**

💡 **Jika connection error:**
   - Pastikan menggunakan URL `bolt://localhost:7687` (bukan http)
   - Cek Neo4j sudah running: `docker ps | findstr neo4j`
   - Cek logs: `docker logs neo4j --tail 20`

### 2.5 Test Connection (Optional tapi Recommended)

Sebelum lanjut, test koneksi ke kedua database:

```powershell
go run testConnection.go
```

Output yang diharapkan:
```
Testing Database Connections...

Testing Cassandra Connection...
   	SUCCESS: Connected to Cassandra
	WARNING: Keyspace 'rumahsakit' belum dibuat
   	Jalankan: go run initSchema.go

Testing Neo4j Connection...
	SUCCESS: Connected to Neo4j
	Response: Connection OK
	WARNING: No constraints found
	Jalankan: go run initSchema.go

✅ All tests completed!
```

Jika ada error di tahap ini, lihat section [Troubleshooting](#-troubleshooting).

### 3. Inisialisasi Schema Database

**Pastikan Cassandra sudah fully ready sebelum jalankan script ini!**

Jalankan script untuk membuat keyspace, tables, dan constraints:

```powershell
go run initSchema.go
```

Output yang diharapkan:
```
Creating Cassandra keyspace and tables (denormalized model)...
Keyspace 'rumahsakit' ready.
Cassandra denormalized schema created successfully.
Creating Neo4j constraints and relationships...
Neo4j constraints created successfully.
```

### 4. Seed Data (Optional)

Untuk mengisi database dengan data dummy (1000 pasien, 500 tenaga medis, dll):

```powershell
go run seed.go
```

**Warning:** Proses seeding bisa memakan waktu 5-15 menit tergantung spesifikasi komputer.

---

## Cara Menjalankan

Setelah setup selesai, kamu bisa:

1. **Membuat query custom** (lihat section berikutnya)
2. **Menggunakan Neo4j Browser** untuk visualisasi data graph
3. **Menggunakan CQL shell** untuk query Cassandra

---

## Cara Membuat Query Baru

### Contoh: Menambahkan User Baru ke Database

#### **Opsi 1: Menambahkan Pasien ke Neo4j**

Buat file baru `queries/addPasien.go`:

```go
package main

import (
	"fmt"
	"log"
	"src/neo4j"
)

func main() {
	// Inisialisasi koneksi Neo4j
	neo4j.InitNeo4j()
	defer neo4j.CloseNeo4j()

	// Data pasien baru
	pasienData := map[string]interface{}{
		"email":         "john.doe@example.com",
		"kata_sandi":    "securepass123",
		"nama_lengkap":  "John Doe",
		"tanggal_lahir": "1990-05-15",
		"nomor_telepon": "081234567890",
		"provinsi":      "DKI Jakarta",
		"kota":          "Jakarta",
		"jalan":         "Jl. Sudirman No. 100",
	}

	// Query Cypher untuk membuat node Pasien
	query := `
		CREATE (p:Pasien {
			email: $email,
			kata_sandi: $kata_sandi,
			nama_lengkap: $nama_lengkap,
			tanggal_lahir: $tanggal_lahir,
			nomor_telepon: $nomor_telepon,
			provinsi: $provinsi,
			kota: $kota,
			jalan: $jalan
		})
		RETURN p.email AS email, p.nama_lengkap AS nama
	`

	// Eksekusi query
	err := neo4j.CreateNeo4j(query, pasienData)
	if err != nil {
		log.Fatalf("Gagal menambahkan pasien: %v", err)
	}

	fmt.Println("Pasien berhasil ditambahkan!")
	fmt.Printf("   Email: %s\n", pasienData["email"])
	fmt.Printf("   Nama: %s\n", pasienData["nama_lengkap"])
}
```

**Jalankan:**
```powershell
go run queries/addPasien.go
```

---

#### **Opsi 2: Menambahkan Obat ke Cassandra**

Buat file baru `queries/addObat.go`:

```go
package main

import (
	"fmt"
	"log"
	"src/cassandra"
)

func main() {
	// Inisialisasi koneksi Cassandra
	cassandra.InitCassandra()
	defer cassandra.Close()

	// Data obat baru
	idObat := "O9999"
	nama := "Paracetamol 500mg"
	label := "Pereda Nyeri"
	harga := 15000.0
	stok := 200

	// Query CQL untuk insert obat
	query := `
		INSERT INTO rumahsakit.obat (id_obat, nama, label, harga, stok)
		VALUES (?, ?, ?, ?, ?)
	`

	// Eksekusi query
	err := cassandra.InsertCassandra(query, idObat, nama, label, harga, stok)
	if err != nil {
		log.Fatalf("Gagal menambahkan obat: %v", err)
	}

	fmt.Println("Obat berhasil ditambahkan!")
	fmt.Printf("   ID: %s\n", idObat)
	fmt.Printf("   Nama: %s\n", nama)
	fmt.Printf("   Harga: Rp %.0f\n", harga)
	fmt.Printf("   Stok: %d\n", stok)
}
```

**Jalankan:**
```powershell
go run queries/addObat.go
```

---

#### **Opsi 3: Query Read - Mendapatkan Data Pasien**

Buat file `queries/getPasien.go`:

```go
package main

import (
	"fmt"
	"log"
	"src/neo4j"
)

func main() {
	neo4j.InitNeo4j()
	defer neo4j.CloseNeo4j()

	// Query untuk mendapatkan 10 pasien pertama
	query := `
		MATCH (p:Pasien)
		RETURN p.email AS email, p.nama_lengkap AS nama, p.kota AS kota
		LIMIT 10
	`

	results, err := neo4j.ReadNeo4j(query, nil)
	if err != nil {
		log.Fatalf("Query gagal: %v", err)
	}

	fmt.Println("\nDaftar Pasien:")
	fmt.Println("=====================================")
	for i, record := range results {
		fmt.Printf("%d. %s - %s (%s)\n",
			i+1,
			record["nama"],
			record["email"],
			record["kota"],
		)
	}
}
```

**Jalankan:**
```powershell
go run queries/getPasien.go
```

---

### Template Query Lainnya

Untuk query yang lebih kompleks, kamu bisa:

1. **Update Data:**
```go
// Neo4j - Update profesi tenaga medis
query := `
	MATCH (t:TenagaMedis {email: $email})
	SET t.profesi = $profesi_baru
	RETURN t
`
params := map[string]interface{}{
	"email": "tm1@rs.com",
	"profesi_baru": "Dokter Spesialis Bedah",
}
neo4j.UpdateNeo4j(query, params)
```

2. **Delete Data:**
```go
// Neo4j - Hapus pasien
query := `
	MATCH (p:Pasien {email: $email})
	DETACH DELETE p
`
params := map[string]interface{}{"email": "john.doe@example.com"}
neo4j.DeleteNeo4j(query, params)
```

3. **Complex Relationship Query:**
```go
// Cari semua tenaga medis di rumah sakit tertentu
query := `
	MATCH (tm:TenagaMedis)-[:bekerja_di]->(d:Departemen)<-[:memiliki_departemen]-(rs:RumahSakit {id_rs: $id_rs})
	RETURN tm.nama_lengkap AS nama, tm.profesi AS profesi, d.nama_departemen AS departemen
`
params := map[string]interface{}{"id_rs": "RS001"}
results, _ := neo4j.ReadNeo4j(query, params)
```

## Troubleshooting

### Error: "Connection refused" saat `docker exec -it cassandra cqlsh`

**Penyebab:** Cassandra masih initializing (butuh 2-3 menit pertama kali dijalankan)

**Solusi:**
```powershell
# 1. Cek apakah Cassandra sudah running
docker ps | findstr cassandra

# 2. Cek logs untuk melihat progress
docker logs cassandra --tail 50

# 3. Tunggu hingga muncul "Startup complete" di logs
# Atau cek status node:
docker exec -it cassandra nodetool status

# 4. Jika sudah muncul "UN" (Up Normal), coba lagi:
docker exec -it cassandra cqlsh

# 5. Jika masih error, restart container:
docker-compose restart cassandra
# Tunggu 2-3 menit, lalu coba lagi
```

### Error: "Cassandra connection failed" saat run Go script
```powershell
# Cek apakah port 9042 accessible
docker exec -it cassandra nodetool status

# Jika output kosong atau error, restart:
docker-compose restart cassandra

# Tunggu 2-3 menit, lalu coba lagi
```

### Error: "Could not reach Neo4j" di Browser

**Penyebab:** URL connection salah atau Neo4j belum ready

**Solusi:**
```powershell
# 1. Pastikan Neo4j running
docker ps | findstr neo4j

# 2. Cek logs Neo4j
docker logs neo4j --tail 20

# 3. Di browser Neo4j, gunakan URL yang BENAR:
# 		BENAR: bolt://localhost:7687
# 		SALAH: http://localhost:7687
# 		SALAH: localhost:7474

# 4. Credentials:
# Username: neo4j
# Password: password123
```

**Jika masih gagal:**
```powershell
# Restart Neo4j container
docker-compose restart neo4j

# Tunggu 30 detik, refresh browser, connect lagi
```

### Error: "Neo4j connection failed" dari Go code
```powershell
# Cek logs Neo4j
docker logs neo4j --tail 30

# Pastikan port 7687 tidak digunakan aplikasi lain
netstat -ano | findstr 7687

# Test connection dengan Go
go run queries/getPasien.go
```

### Error: "module not found"
```powershell
# Pastikan kamu di root directory project
cd c:\Users\ACER\Documents\GitHub\Pemdal-Kel5-Project2

# Re-download dependencies
go mod tidy
go mod download
```

### Cassandra health check failed
```powershell
# Tunggu hingga Cassandra fully initialized (bisa 1-2 menit)
docker exec -it cassandra nodetool status
```

---

## 📚 Referensi

- [Cassandra CQL Documentation](https://cassandra.apache.org/doc/latest/cql/)
- [Neo4j Cypher Manual](https://neo4j.com/docs/cypher-manual/current/)
- [gocql Driver](https://github.com/gocql/gocql)
- [Neo4j Go Driver](https://github.com/neo4j/neo4j-go-driver)

---

## Contributors

**Kelompok 5 - Pemdal Project 2**

---