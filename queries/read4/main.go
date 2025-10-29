package main

import (
	"fmt"
	"log"
	"time"

	"src/neo4j"
)

type MedikJanjiTemu struct {
	Email string
	Nama string
	Profesi string
	JumlahJanjiTemu int
}

func getMedikJanjiTemuFromNeo4j() ([]MedikJanjiTemu, error) {
	query := `
		MATCH (tm:TenagaMedis)<-[:dengan_dokter]-(jt:JanjiTemu)
		RETURN tm.email AS email, tm.nama_lengkap AS nama, tm.profesi AS profesi, COUNT(jt) AS jumlah_janji_temu
	`

	records, err := neo4j.ReadNeo4j(query, nil)
	if err != nil {
		return nil, fmt.Errorf("gagal membaca dari Neo4j: %v", err)
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("tidak ditemukan janji temu dengan tenaga medis")
	}

	var results []MedikJanjiTemu
	for _, record := range records {
		results = append(results, MedikJanjiTemu{
			Email:           fmt.Sprintf("%v", record["email"]),
			Nama:            fmt.Sprintf("%v", record["nama"]),
			Profesi:         fmt.Sprintf("%v", record["profesi"]),
			JumlahJanjiTemu: int(record["jumlah_janji_temu"].(int64)),
		})
	}

	return results, nil
}

func displayResult(mediks []MedikJanjiTemu) {
	fmt.Println("     Jumlah Janji Temu per Tenaga Medis")
	fmt.Printf("%-5s %-40s %-30s %-25s %s\n", "No", "Email", "Nama Lengkap", "Profesi", "Jumlah Janji Temu")

	if len(mediks) == 0 {
		fmt.Println("Tidak ada data janji temu dengan tenaga medis.")
		return
	}

	for i, m := range mediks {
		fmt.Printf("%-5d %-40s %-30s %-25s %d\n",
			i+1,
			m.Email,
			m.Nama,
			m.Profesi,
			m.JumlahJanjiTemu)
	}
}

func main() {
	neo4j.InitNeo4j()
	defer neo4j.CloseNeo4j()
	
	start := time.Now()

	result, err := getMedikJanjiTemuFromNeo4j()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	displayResult(result)

	elapsed := time.Since(start)
	fmt.Printf("\nTime: %.3f seconds (%d ms)\n", elapsed.Seconds(), elapsed.Milliseconds())
}