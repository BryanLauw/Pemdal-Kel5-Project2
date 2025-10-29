package main

import (
	"fmt"
	"log"
	"time"

	"src/neo4j"
)

type DepartemenRS struct {
	NamaDepartemen string
	RumahSakit     string
	IdRS           string
}

func insertDepartemenToNeo4j() (*DepartemenRS, error) {
	query := `
		MATCH (rs:RumahSakit {nama_rumah_sakit: $nama_rumah_sakit})
		WITH rs LIMIT 1
		CREATE (d:Departemen {
			nama_departemen: $nama_departemen,
			gedung: $gedung
		})
		CREATE (rs)-[:memiliki_departemen]->(d)
		RETURN d.nama_departemen AS departemen, rs.nama_rumah_sakit AS rumah_sakit, rs.id_rs AS id_rs
	`

	params := map[string]interface{}{
		"nama_rumah_sakit": "RS Sehat Selalu",
		"nama_departemen":  "Kardiologi",
		"gedung":           "Gedung A",
	}

	results, err := neo4j.CreateAndReturnNeo4j(query, params)
	if err != nil {
		return nil, fmt.Errorf("gagal menambahkan departemen: %v", err)
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("tidak ada data yang dikembalikan atau rumah sakit tidak ditemukan")
	}

	record := results[0]
	dept := &DepartemenRS{
		NamaDepartemen: record["departemen"].(string),
		RumahSakit:     record["rumah_sakit"].(string),
		IdRS:           record["id_rs"].(string),
	}

	return dept, nil
}

func displayResult(dept *DepartemenRS) {
	fmt.Println("\n=== INSERT 4: Menambah Departemen Tertentu pada RS ===")
	fmt.Printf("Nama Departemen    : %s\n", dept.NamaDepartemen)
	fmt.Printf("Rumah Sakit        : %s\n", dept.RumahSakit)
	fmt.Printf("ID RS              : %s\n", dept.IdRS)
	fmt.Println("\n✓ Departemen berhasil ditambahkan!")
}

func main() {
	neo4j.InitNeo4j()
	defer neo4j.CloseNeo4j()

	start := time.Now()
	dept, err := insertDepartemenToNeo4j()
	elapsed := time.Since(start)

	if err != nil {
		log.Fatalf("Error inserting department: %v", err)
	}

	displayResult(dept)
	fmt.Printf("\nQuery Execution Time: %.3f seconds (%d ms)\n", elapsed.Seconds(), elapsed.Milliseconds())
}
