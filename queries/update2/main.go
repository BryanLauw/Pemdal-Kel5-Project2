// =============================================================
// Query: Pindahtugaskan tenaga medis ke departemen lain.
// =============================================================

package main

import (
	"fmt"
	"log"
	"time"

	"src/neo4j"
)

func main() {
	neo4j.InitNeo4j()
	defer neo4j.CloseNeo4j()

	// 1. Ambil satu tenaga medis
	query := `
		MATCH (t:TenagaMedis)
		OPTIONAL MATCH (t)-[:BEKERJA_DI]->(d:Departemen)
		RETURN t.email AS email, d.nama_departemen AS departemen
		LIMIT 1
	`
	records, err := neo4j.ReadNeo4j(query, nil)
	if err != nil {
		log.Fatalf("Gagal membaca data: %v", err)
	}
	if len(records) == 0 {
		log.Fatalf("Tidak ditemukan node TenagaMedis di database.")
	}

	var email string
	var departemenBaru string
	if recEmail, ok := records[0]["email"].(string); ok {
		email = recEmail
	} else {
		log.Fatalf("Record tidak memiliki email yang valid.")
	}
	departemenBaru = "Departemen-Baru"

	// Print record sebelum pindah (lihat lokasi departemen lama)
	fmt.Println("=== Sebelum Pindah ===")
	fmt.Printf("Record sebelum pindah: %v\n", records[0])
	fmt.Printf("Email: %s ke departemen: %s\n", email, departemenBaru)

	// 2. Lakukan pemindahan
	start := time.Now()
	if err = PindahTenagaMedis(email, departemenBaru); err != nil {
		log.Fatalf("Gagal memindahkan tenaga medis: %v", err)
	}
	duration := time.Since(start)
	fmt.Printf("\nPindah berhasil (%.2f ms)\n\n", float64(duration.Milliseconds()))

	// 3. Ambil kembali dan print untuk melihat departemen baru
	fmt.Println("=== Setelah Pindah ===")
	queryAfter := `
		MATCH (t:TenagaMedis {email:$email})
		OPTIONAL MATCH (t)-[:BEKERJA_DI]->(d:Departemen)
		RETURN t.email AS email, d.nama_departemen AS departemen
	`
	recordsAfter, err := neo4j.ReadNeo4j(queryAfter, map[string]interface{}{"email": email})
	if err != nil {
		log.Fatalf("Gagal membaca data setelah pindah: %v", err)
	}
	if len(recordsAfter) == 0 {
		fmt.Println("Tidak ditemukan record setelah pindah.")
	} else {
		fmt.Printf("Record setelah pindah: %v\n", recordsAfter[0])
	}
}

func PindahTenagaMedis(email string, departemenBaru string) error {
	query := `
		MATCH (t:TenagaMedis {email:$email})
		OPTIONAL MATCH (t)-[r:BEKERJA_DI]->(d:Departemen)
		FOREACH (_ IN CASE WHEN r IS NULL THEN [] ELSE [1] END | DELETE r)
		WITH t
		MERGE (d2:Departemen {nama_departemen:$dept})
		MERGE (t)-[:BEKERJA_DI]->(d2)
	`
	params := map[string]interface{}{"email": email, "dept": departemenBaru}
	return neo4j.UpdateNeo4j(query, params)
}
