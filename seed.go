package main

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	"src/cassandra"
	"src/neo4j"

	faker "github.com/go-faker/faker/v4"
)

// ===============================================
//   CONSTANTS (DARI PERMINTAAN USER)
// ===============================================

const (
	NumPasien      = 1000
	NumTenagaMedis = 500
	NumRumahSakit  = 100 // FIXED: Menjadi 100
	NumDepartemen  = 300 // FIXED: Menjadi 300
	NumLayanan     = 500
	NumObat        = 1000
)

func randomProvince() string {
	provinces := []string{
		"Jawa Barat", "Jawa Tengah", "Jawa Timur",
		"DKI Jakarta", "Banten", "Sumatera Utara",
		"Sumatera Barat", "Kalimantan Timur", "Sulawesi Selatan",
	}
	return provinces[rand.Intn(len(provinces))]
}

func randomCity() string {
	cities := []string{
		"Bandung", "Semarang", "Surabaya",
		"Jakarta", "Serang", "Medan",
		"Padang", "Balikpapan", "Makassar",
	}
	return cities[rand.Intn(len(cities))]
}


// ===============================================
//   📦 DATA GENERATION
// ===============================================

func generatePasienData() []map[string]interface{} {
	data := make([]map[string]interface{}, NumPasien)
	for i := 0; i < NumPasien; i++ {
		data[i] = map[string]interface{}{
			"email":         faker.Email(),
			"kata_sandi":    "pass123",
			"nama_lengkap":  faker.Name(),
			"tanggal_lahir": faker.Date(),
			"nomor_telepon": fmt.Sprintf("08%d", rand.Intn(900000000)+100000000),
			"provinsi":      randomProvince(),
			"kota":          randomCity(),
			"jalan":         fmt.Sprintf("Jl. %s No.%d", faker.Word(), rand.Intn(300)+1),
		}
	}
	return data
}

func generateTenagaMedisData() []map[string]interface{} {
	data := make([]map[string]interface{}, NumTenagaMedis)
	professions := []string{"Dokter Umum", "Dokter Spesialis Anak", "Perawat", "Bidan", "Ahli Gizi", "Dokter Gigi"}
	for i := 0; i < NumTenagaMedis; i++ {
		data[i] = map[string]interface{}{
			"email":         fmt.Sprintf("tm%d@rs.com", i+1),
			"NIKes":         fmt.Sprintf("%08d", rand.Intn(99999999)),
			"profesi":       professions[i%len(professions)],
			"kata_sandi":    "docpass",
			"nama_lengkap":  faker.Name(),
			"tanggal_lahir": faker.Date(),
			"nomor_telepon": fmt.Sprintf("08%d", rand.Intn(900000000)+100000000),
			"provinsi":      randomProvince(),
			"kota":          randomCity(),
			"jalan":         fmt.Sprintf("Jl. %s No.%d", faker.Word(), rand.Intn(300)+1),
		}
	}
	return data
}

func generateRumahSakitData() []map[string]interface{} {
	data := make([]map[string]interface{}, NumRumahSakit)
	for i := 0; i < NumRumahSakit; i++ {
		id_rs := fmt.Sprintf("RS%03d", i+1)
		data[i] = map[string]interface{}{
			"id_rs":            id_rs,
			"email":            "info@" + id_rs + ".com",
			"nama_rumah_sakit": "RSUD Sejahtera " + faker.LastName(),
			"no_telepon":       fmt.Sprintf("021-%06d", rand.Intn(999999)),
			"provinsi":         randomProvince(),
			"kota":             randomCity(),
			"jalan":            fmt.Sprintf("Jl. %s No.%d", faker.Word(), rand.Intn(300)+1),
		}
	}
	return data
}

func generateDepartemenData() []map[string]interface{} {
	names := []string{"Poli Umum", "Poli Anak", "Gawat Darurat", "Poli Gigi", "Poli Jantung", "Farmasi"}
	data := make([]map[string]interface{}, NumDepartemen)
	for i := 0; i < NumDepartemen; i++ {
		data[i] = map[string]interface{}{
			"nama_departemen": names[i%len(names)] + " " + strconv.Itoa(i+1),
			"gedung":          "Gedung " + string('A'+i%5),
		}
	}
	return data
}

func generateLayananMedisData() []map[string]interface{} {
	names := []string{"Konsultasi Umum", "Pemeriksaan Anak", "Tes Lab Darah", "Fisioterapi", "Rawat Inap", "Minor Surgery"}
	data := make([]map[string]interface{}, NumLayanan)
	for i := 0; i < NumLayanan; i++ {
		data[i] = map[string]interface{}{
			"id_layanan":    fmt.Sprintf("L%03d", i+1),
			"nama_layanan":  names[i%len(names)] + " " + strconv.Itoa(i+1),
			"biaya_layanan": float64(rand.Intn(400)+100) * 1000.0, // 100k - 500k
		}
	}
	return data
}

func generateBayminData(pasienData []map[string]interface{}) []map[string]interface{} {
	data := make([]map[string]interface{}, len(pasienData))
	colors := []string{"Merah", "Biru", "Hijau", "Kuning", "Putih", "Hitam"}
	for i, pasien := range pasienData {
		data[i] = map[string]interface{}{
			"email_pasien": pasien["email"],
			"id_perangkat": fmt.Sprintf("BAYMIN-%04d", i+1),
			"warna":        colors[rand.Intn(len(colors))],
		}
	}
	return data
}

func generateObatData() []map[string]interface{} {
	data := make([]map[string]interface{}, NumObat)
	labels := []string{"Pereda Nyeri", "Antibiotik", "Vitamin", "Obat Batuk", "Antasida"}
	for i := 0; i < NumObat; i++ {
		data[i] = map[string]interface{}{
			"id_obat": fmt.Sprintf("O%04d", i+1),
			"nama":    faker.Word() + " " + faker.Word(),
			"label":   labels[i%len(labels)],
			"harga":   float64(rand.Intn(50)+5) * 1000.0,
			"stok":    rand.Intn(200) + 50,
		}
	}
	return data
}

// ===============================================
//   📦 SEEDER CASSANDRA
// ===============================================

func seedCassandra(obatData, rsData, layananData []map[string]interface{}, bayminData []map[string]interface{}) {
	fmt.Println("\n📦 Seeding Cassandra tables...")

	// --- 4️⃣ MASTER OBAT ---
	for _, data := range obatData {
		query := `INSERT INTO rumahsakit.obat (id_obat, nama, label, harga, stok) VALUES (?, ?, ?, ?, ?)`
		if err := cassandra.InsertCassandra(query, data["id_obat"], data["nama"], data["label"], data["harga"], data["stok"]); err != nil {
			log.Printf("❌ Error inserting obat %s: %v", data["id_obat"], err)
		}
	}

	// --- 6️⃣ PELAKSANAAN LAYANAN MEDIS (lokasi_layanan) ---
	// FIXED: Menggunakan multiple RS. Kita asumsikan setiap RS menawarkan 5-10 layanan secara acak.
	numServicesToOffer := rand.Intn(6) + 5 // 5 to 10 services per RS
	
	// Gunakan map untuk memastikan tidak ada duplikasi RS-Layanan yang dimasukkan ke Cassandra
	seededLocations := make(map[string]bool)

	for _, rs := range rsData {
		idRs := rs["id_rs"].(string)
		
		// Acak Layanan untuk setiap RS
		rand.Shuffle(len(layananData), func(i, j int) { layananData[i], layananData[j] = layananData[j], layananData[i] })
		
		for i := 0; i < numServicesToOffer && i < len(layananData); i++ {
			layanan := layananData[i]
			idLayanan := layanan["id_layanan"].(string)
			
			key := idRs + "-" + idLayanan
			if _, exists := seededLocations[key]; exists {
				continue // Skip if already seeded (shouldn't happen with shuffle, but good practice)
			}
			seededLocations[key] = true

			query := `INSERT INTO rumahsakit.lokasi_layanan (id_rs, id_layanan, nama_layanan, biaya_layanan) VALUES (?, ?, ?, ?)`
			err := cassandra.InsertCassandra(query, idRs, idLayanan, layanan["nama_layanan"], layanan["biaya_layanan"])
			if err != nil {
				log.Printf("❌ Error inserting lokasi_layanan %s-%s: %v", idRs, idLayanan, err)
			}
		}
	}
	
	// --- Data Transaksional Dummy (Hanya 5 contoh) ---
	now := time.Now()
	for i := 1; i <= 5; i++ {
		// 1️⃣ LOG AKTIVITAS (dari Baymin ID random)
		if len(bayminData) > 0 {
			idPerangkat := bayminData[rand.Intn(len(bayminData))]["id_perangkat"]
			cassandra.InsertCassandra(`INSERT INTO rumahsakit.log_aktivitas (id_perangkat, waktu_aktivitas, detail_aktivitas) VALUES (?, ?, ?)`, 
				idPerangkat, now.Add(-time.Duration(i)*time.Hour), "Status perangkat: "+faker.Sentence())
		}
		
		// 2️⃣ & 3️⃣ PEMESANAN OBAT & DETAIL PEMESANAN OBAT
		if len(obatData) > 1 {
			poID := fmt.Sprintf("POB%05d", i)
			obatMap := map[string]int{
				obatData[rand.Intn(len(obatData))]["id_obat"].(string): rand.Intn(5) + 1,
				obatData[rand.Intn(len(obatData))]["id_obat"].(string): rand.Intn(5) + 1,
			}
			emailPemesan := faker.Email()
			
			cassandra.InsertCassandra(`INSERT INTO rumahsakit.pemesanan_obat (id_pesanan, email_pemesan, waktu_pemesanan, status_pemesanan) VALUES (?, ?, ?, ?)`, 
				poID, emailPemesan, now.Add(time.Duration(i)*time.Hour), []string{"TERKIRIM", "DIPROSES"}[rand.Intn(2)])
			
			cassandra.InsertCassandra(`INSERT INTO rumahsakit.detail_pesanan_obat (id_pesanan, daftar_obat) VALUES (?, ?)`, poID, obatMap)
		}
	}

	fmt.Println("✅ Cassandra tables seeded successfully.")
}

// ===============================================
//   🕸️  SEEDER NEO4J
// ===============================================

func seedNeo4j(pasienData, tenagaMedisData, rsData, departemenData, layananMedisData, bayminData, obatData []map[string]interface{}) {
	fmt.Println("\n🕸️ Seeding Neo4j nodes and relationships...")

	// --- Create Nodes ---

	fmt.Println("   -> Creating Nodes...")
	// Pasien (1000)
	for _, data := range pasienData {
		query := `CREATE (p:Pasien {email: $email, kata_sandi: $kata_sandi, nama_lengkap: $nama_lengkap, tanggal_lahir: $tanggal_lahir, nomor_telepon: $nomor_telepon, provinsi: $provinsi, kota: $kota, jalan: $jalan})`
		if err := neo4j.CreateNeo4j(query, data); err != nil {
			// Log as warning since bulk constraint failure might happen
			log.Printf("⚠️ Error creating Pasien %s: %v", data["email"], err)
		}
	}

	// TenagaMedis (500)
	for _, data := range tenagaMedisData {
		query := `CREATE (t:TenagaMedis {email: $email, NIKes: $NIKes, profesi: $profesi, kata_sandi: $kata_sandi, nama_lengkap: $nama_lengkap, tanggal_lahir: $tanggal_lahir, nomor_telepon: $nomor_telepon, provinsi: $provinsi, kota: $kota, jalan: $jalan})`
		neo4j.CreateNeo4j(query, data)
	}

	// RumahSakit (100) - FIXED
	for _, data := range rsData {
		query := `CREATE (r:RumahSakit {id_rs: $id_rs, email: $email, nama_rumah_sakit: $nama_rumah_sakit, no_telepon: $no_telepon, provinsi: $provinsi, kota: $kota, jalan: $jalan})`
		neo4j.CreateNeo4j(query, data)
	}

	// Departemen (300)
	for _, data := range departemenData {
		query := `CREATE (d:Departemen {nama_departemen: $nama_departemen, gedung: $gedung})`
		neo4j.CreateNeo4j(query, data)
	}

	// LayananMedis (500)
	for _, data := range layananMedisData {
		query := `CREATE (l:LayananMedis {id_layanan: $id_layanan, nama_layanan: $nama_layanan, biaya_layanan: $biaya_layanan})`
		neo4j.CreateNeo4j(query, data)
	}

	// Baymin (1000)
	for _, data := range bayminData {
		query := `CREATE (b:Baymin {id_perangkat: $id_perangkat, warna: $warna, email_pasien: $email_pasien})`
		neo4j.CreateNeo4j(query, data)
	}
	
	// --- Create Relationships ---
	
	fmt.Println("   -> Creating Relationships...")

	// 1. Pasien memiliki_perangkat Baymin (1000 relasi)
	for _, data := range bayminData {
		query := `MATCH (p:Pasien {email: $email_pasien}), (b:Baymin {id_perangkat: $id_perangkat}) MERGE (p)-[:memiliki_perangkat]->(b)`
		params := map[string]interface{}{"email_pasien": data["email_pasien"], "id_perangkat": data["id_perangkat"]}
		neo4j.UpdateNeo4j(query, params)
	}

	// 2. TenagaMedis bekerja_di Departemen (500 relasi)
	for i, tm := range tenagaMedisData {
		// Distribusi 500 TM ke 300 Departemen
		dept := departemenData[i%len(departemenData)]
		query := `MATCH (t:TenagaMedis {email: $email_tm}), (d:Departemen {nama_departemen: $nama_dept}) MERGE (t)-[:bekerja_di]->(d)`
		params := map[string]interface{}{"email_tm": tm["email"], "nama_dept": dept["nama_departemen"]}
		neo4j.UpdateNeo4j(query, params)
	}

	// 3. RumahSakit memiliki_departemen Departemen (300 relasi)
	for i, dept := range departemenData {
		// Distribusi 300 Departemen ke 100 RS
		rs := rsData[i%len(rsData)]
		query := `MATCH (r:RumahSakit {id_rs: $id_rs}), (d:Departemen {nama_departemen: $nama_dept}) MERGE (r)-[:memiliki_departemen]->(d)`
		params := map[string]interface{}{"id_rs": rs["id_rs"], "nama_dept": dept["nama_departemen"]}
		neo4j.UpdateNeo4j(query, params)
	}

	// 4. RumahSakit menawarkan_layanan LayananMedis (5-10 relasi per RS)
	// Logika harus sama dengan yang digunakan di seedCassandra untuk konsistensi data
	numServicesToOffer := rand.Intn(6) + 5 
	for _, rs := range rsData {
		idRs := rs["id_rs"].(string)
		rand.Shuffle(len(layananMedisData), func(i, j int) { layananMedisData[i], layananMedisData[j] = layananMedisData[j], layananMedisData[i] })
		
		for i := 0; i < numServicesToOffer && i < len(layananMedisData); i++ {
			layanan := layananMedisData[i]
			query := `MATCH (r:RumahSakit {id_rs: $id_rs}), (l:LayananMedis {id_layanan: $id_layanan}) MERGE (r)-[:menawarkan_layanan]->(l)`
			params := map[string]interface{}{"id_rs": idRs, "id_layanan": layanan["id_layanan"]}
			neo4j.UpdateNeo4j(query, params)
		}
	}

	// 5. JanjiTemu, Resep, DetailResep (100 sample transactions)
	fmt.Println("   -> Creating 100 Sample Transactions (JanjiTemu/Resep)...")
	for i := 1; i <= 100; i++ {
		// Randomly select entities
		pasien := pasienData[rand.Intn(len(pasienData))]
		dokter := tenagaMedisData[rand.Intn(len(tenagaMedisData))]
		rs := rsData[rand.Intn(len(rsData))]
		
		// Create JanjiTemu
		jtID := fmt.Sprintf("JT%05d", i)
		janjiTemuData := map[string]interface{}{
			"id_janji_temu": jtID,
			"waktu_pelaksanaan": time.Now().Add(time.Duration(rand.Intn(30)*24) * time.Hour).Format("2006-01-02 15:04:05"),
			"alasan":            faker.Sentence(),
			"status":            []string{"SELESAI", "TERJADWAL", "BATAL"}[rand.Intn(3)],
		}
		neo4j.CreateNeo4j(`CREATE (j:JanjiTemu {id_janji_temu: $id_janji_temu, waktu_pelaksanaan: $waktu_pelaksanaan, alasan: $alasan, status: $status})`, janjiTemuData)
		
		// Link JanjiTemu
		neo4j.UpdateNeo4j(`MATCH (p:Pasien {email: $p_email}), (j:JanjiTemu {id_janji_temu: $jt_id}) MERGE (p)<-[:memiliki_janji]-(j)`, map[string]interface{}{"p_email": pasien["email"], "jt_id": jtID})
		neo4j.UpdateNeo4j(`MATCH (j:JanjiTemu {id_janji_temu: $jt_id}), (t:TenagaMedis {email: $t_email}) MERGE (j)-[:dengan_dokter]->(t)`, map[string]interface{}{"t_email": dokter["email"], "jt_id": jtID})
		neo4j.UpdateNeo4j(`MATCH (j:JanjiTemu {id_janji_temu: $jt_id}), (r:RumahSakit {id_rs: $id_rs}) MERGE (j)-[:di_rs]->(r)`, map[string]interface{}{"id_rs": rs["id_rs"], "jt_id": jtID})
		
		// Create Resep/DetailResep jika status SELESAI
		if janjiTemuData["status"] == "SELESAI" {
			resepID := fmt.Sprintf("R%05d", i)
			resepData := map[string]interface{}{"id_resep": resepID, "penyakit": faker.Word() + " " + faker.Word()}
			neo4j.CreateNeo4j(`CREATE (r:Resep {id_resep: $id_resep, penyakit: $penyakit})`, resepData)
			neo4j.UpdateNeo4j(`MATCH (j:JanjiTemu {id_janji_temu: $jt_id}), (r:Resep {id_resep: $resep_id}) MERGE (j)-[:menghasilkan_resep]->(r)`, map[string]interface{}{"jt_id": jtID, "resep_id": resepID})

			// Add 2 random DetailResep (Obat)
			rand.Shuffle(len(obatData), func(i, j int) { obatData[i], obatData[j] = obatData[j], obatData[i] })
			for j := 0; j < 2; j++ {
				obat := obatData[j]
				drData := map[string]interface{}{"id_obat": obat["id_obat"], "dosis": []string{"1x Sehari", "2x Sehari", "3x Sehari"}[rand.Intn(3)]}
				neo4j.CreateNeo4j(`CREATE (dr:DetailResep {id_obat: $id_obat, dosis: $dosis})`, drData)
				neo4j.UpdateNeo4j(`MATCH (r:Resep {id_resep: $resep_id}), (dr:DetailResep {id_obat: $id_obat}) MERGE (r)-[:memiliki_detail]->(dr)`, map[string]interface{}{"resep_id": resepID, "id_obat": obat["id_obat"]})
			}
		}
	}


	fmt.Println("✅ Neo4j nodes and relationships seeded successfully.")
}

// ===============================================
//   ⚙️  MAIN FUNCTION (untuk memanggil seeder)
// ===============================================

func main() {
	// --- Generate Data ---
	pasienData := generatePasienData()
	tenagaMedisData := generateTenagaMedisData()
	rsData := generateRumahSakitData()
	departemenData := generateDepartemenData()
	layananMedisData := generateLayananMedisData()
	bayminData := generateBayminData(pasienData)
	obatData := generateObatData()

	// --- Cassandra ---
	cassandra.InitCassandra()
	defer cassandra.Session.Close()
	// Asumsi createCassandraSchema() dipanggil di sini
	seedCassandra(obatData, rsData, layananMedisData, bayminData)

	// --- Neo4j ---
	neo4j.InitNeo4j()
	defer neo4j.CloseNeo4j()
	// Asumsi createNeo4jSchema() dipanggil di sini
	seedNeo4j(pasienData, tenagaMedisData, rsData, departemenData, layananMedisData, bayminData, obatData)
}
