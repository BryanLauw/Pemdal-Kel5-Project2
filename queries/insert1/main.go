package main

import (
	"fmt"
	"log"
	"time"

	"src/neo4j"
)

type Pasien struct {
	Email       string
	NamaLengkap string
}

func insertPasienToNeo4j() (*Pasien, error) {
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

	params := map[string]interface{}{
		"email":         "andi@example.com",
		"kata_sandi":    "hashed_password",
		"nama_lengkap":  "Andi Setiawan",
		"tanggal_lahir": "1995-04-21",
		"nomor_telepon": "08123456789",
		"provinsi":      "Jawa Barat",
		"kota":          "Bandung",
		"jalan":         "Jl. Merdeka 123",
	}

	results, err := neo4j.CreateAndReturnNeo4j(query, params)
	if err != nil {
		return nil, fmt.Errorf("gagal menambahkan pasien: %v", err)
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("tidak ada data yang dikembalikan")
	}

	record := results[0]
	pasien := &Pasien{
		Email:       record["email"].(string),
		NamaLengkap: record["nama"].(string),
	}

	return pasien, nil
}

func displayResult(pasien *Pasien) {
	fmt.Println("\n=== INSERT 1: Menambahkan Pengguna Baru ===")
	fmt.Printf("Email        : %s\n", pasien.Email)
	fmt.Printf("Nama Lengkap : %s\n", pasien.NamaLengkap)
	fmt.Println("\n✓ Pasien berhasil ditambahkan!")
}

func main() {
	neo4j.InitNeo4j()
	defer neo4j.CloseNeo4j()

	start := time.Now()
	pasien, err := insertPasienToNeo4j()
	elapsed := time.Since(start)

	if err != nil {
		log.Fatalf("Error inserting patient: %v", err)
	}

	displayResult(pasien)
	fmt.Printf("\nQuery Execution Time: %.3f seconds (%d ms)\n", elapsed.Seconds(), elapsed.Milliseconds())
}
