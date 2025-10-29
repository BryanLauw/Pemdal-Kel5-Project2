package main

import (
	"fmt"
	"log"
	"time"

	"src/neo4j"
)

type Pasien struct {
	Email string
}

func insertPasienFromUserByName() (*Pasien, error) {
	// Step 1: Find user by nama_lengkap (equivalent to SELECT email FROM user WHERE nama_lengkap = 'Andi Setiawan' LIMIT 1)
	// Step 2: Create Pasien node with that email (equivalent to INSERT INTO pasien (email))
	// Note: Di Neo4j, kita buat label tambahan :Pasien untuk user yang sudah ada

	query := `
		MATCH (u:Pasien {nama_lengkap: $nama_lengkap})
		WITH u LIMIT 1
		SET u:PasienTerdaftar
		RETURN u.email AS email
	`

	params := map[string]interface{}{
		"nama_lengkap": "Andi Setiawan",
	}

	results, err := neo4j.CreateAndReturnNeo4j(query, params)
	if err != nil {
		return nil, fmt.Errorf("gagal menambahkan pasien: %v", err)
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("user dengan nama tersebut tidak ditemukan")
	}

	record := results[0]
	pasien := &Pasien{
		Email: record["email"].(string),
	}

	return pasien, nil
}

func displayResult(pasien *Pasien) {
	fmt.Println("\n=== INSERT 2: Tambah Pasien Berdasarkan Nama dari User ===")
	fmt.Printf("Email Pasien : %s\n", pasien.Email)
	fmt.Println("\n✓ Pasien berhasil ditambahkan berdasarkan user yang ada!")
	fmt.Println("   (Di Neo4j: menambahkan label :PasienTerdaftar pada node yang sudah ada)")
}

func main() {
	neo4j.InitNeo4j()
	defer neo4j.CloseNeo4j()

	start := time.Now()
	pasien, err := insertPasienFromUserByName()
	elapsed := time.Since(start)

	if err != nil {
		log.Fatalf("Error inserting patient from user: %v", err)
	}

	displayResult(pasien)
	fmt.Printf("\nQuery Execution Time: %.3f seconds (%d ms)\n", elapsed.Seconds(), elapsed.Milliseconds())
}
