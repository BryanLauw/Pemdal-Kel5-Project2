package main

import (
	"fmt"
	"log"
	"time"

	"src/neo4j"
)

type DokterSpesialis struct {
	NamaDokter string
	Telepon    string
	Departemen string
	RumahSakit string
	Alamat     string
}

func getDokterSpesialisFromNeo4j() ([]DokterSpesialis, error) {
	query := `
		MATCH (tm:TenagaMedis {profesi: $profesi})-[:bekerja_di]->(d:Departemen)
			  <-[:memiliki_departemen]-(rs:RumahSakit {kota: $kota})
		RETURN 
			tm.nama_lengkap AS nama_dokter,
			tm.nomor_telepon AS telepon,
			d.nama_departemen AS departemen,
			rs.nama_rumah_sakit AS rumah_sakit,
			rs.jalan AS alamat
		ORDER BY rs.nama_rumah_sakit, tm.nama_lengkap
		LIMIT 50
	`

	params := map[string]interface{}{
		"profesi": "Dokter Spesialis Anak",
		"kota":    "Bandung",
	}

	results, err := neo4j.ReadNeo4j(query, params)
	if err != nil {
		return nil, fmt.Errorf("gagal mencari dokter spesialis: %v", err)
	}

	var dokters []DokterSpesialis
	for _, record := range results {
		dokter := DokterSpesialis{
			NamaDokter: getStringValue(record, "nama_dokter"),
			Telepon:    getStringValue(record, "telepon"),
			Departemen: getStringValue(record, "departemen"),
			RumahSakit: getStringValue(record, "rumah_sakit"),
			Alamat:     getStringValue(record, "alamat"),
		}
		dokters = append(dokters, dokter)
	}

	return dokters, nil
}

func getStringValue(record map[string]interface{}, key string) string {
	if val, ok := record[key]; ok && val != nil {
		return val.(string)
	}
	return ""
}

func displayResult(dokters []DokterSpesialis) {
	fmt.Println("\n=== SPECIAL GRAPH: Cari Dokter Spesialis Anak di Bandung ===")
	fmt.Println("Keunggulan Graph: Relationship Traversal yang Efisien\n")

	if len(dokters) == 0 {
		fmt.Println("Tidak ada dokter spesialis yang ditemukan.")
		return
	}

	fmt.Printf("Total dokter ditemukan: %d\n\n", len(dokters))
	fmt.Printf("%-5s %-30s %-15s %-20s %-30s\n", "No", "Nama Dokter", "Telepon", "Departemen", "Rumah Sakit")
	fmt.Println("─────────────────────────────────────────────────────────────────────────────────────────────────────────")

	for i, d := range dokters {
		fmt.Printf("%-5d %-30s %-15s %-20s %-30s\n",
			i+1,
			d.NamaDokter,
			d.Telepon,
			d.Departemen,
			d.RumahSakit)
	}

	fmt.Println("\n✓ Query berhasil dijalankan dengan graph traversal!")
}

func main() {
	neo4j.InitNeo4j()
	defer neo4j.CloseNeo4j()

	start := time.Now()
	dokters, err := getDokterSpesialisFromNeo4j()
	elapsed := time.Since(start)

	if err != nil {
		log.Fatalf("Error getting specialist doctors: %v", err)
	}

	displayResult(dokters)
	fmt.Printf("\nQuery Execution Time: %.3f seconds (%d ms)\n", elapsed.Seconds(), elapsed.Milliseconds())
}
