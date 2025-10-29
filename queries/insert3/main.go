package main

import (
	"fmt"
	"log"
	"time"

	"src/neo4j"
)

type RumahSakit struct {
	IdRS           string
	NamaRumahSakit string
}

func insertRumahSakitToNeo4j() (*RumahSakit, error) {
	query := `
		CREATE (r:RumahSakit {
			id_rs: $id_rs,
			email: $email,
			nama_rumah_sakit: $nama_rumah_sakit,
			no_telepon: $no_telepon,
			provinsi: $provinsi,
			kota: $kota,
			jalan: $jalan
		})
		RETURN r.id_rs AS id_rs, r.nama_rumah_sakit AS nama
	`

	params := map[string]interface{}{
		"id_rs":            "RS999",
		"email":            "rs@example.com",
		"nama_rumah_sakit": "RS Sehat Selalu",
		"no_telepon":       "0221234567",
		"provinsi":         "Jawa Barat",
		"kota":             "Bandung",
		"jalan":            "Jl. Kesehatan 10",
	}

	results, err := neo4j.CreateAndReturnNeo4j(query, params)
	if err != nil {
		return nil, fmt.Errorf("gagal menambahkan rumah sakit: %v", err)
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("tidak ada data yang dikembalikan")
	}

	record := results[0]
	rs := &RumahSakit{
		IdRS:           record["id_rs"].(string),
		NamaRumahSakit: record["nama"].(string),
	}

	return rs, nil
}

func displayResult(rs *RumahSakit) {
	fmt.Println("\n=== INSERT 3: Menambah Rumah Sakit Baru ===")
	fmt.Printf("ID RS              : %s\n", rs.IdRS)
	fmt.Printf("Nama Rumah Sakit   : %s\n", rs.NamaRumahSakit)
	fmt.Println("\n✓ Rumah Sakit berhasil ditambahkan!")
}

func main() {
	neo4j.InitNeo4j()
	defer neo4j.CloseNeo4j()

	start := time.Now()
	rs, err := insertRumahSakitToNeo4j()
	elapsed := time.Since(start)

	if err != nil {
		log.Fatalf("Error inserting hospital: %v", err)
	}

	displayResult(rs)
	fmt.Printf("\nQuery Execution Time: %.3f seconds (%d ms)\n", elapsed.Seconds(), elapsed.Milliseconds())
}
