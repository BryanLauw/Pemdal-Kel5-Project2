package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"src/neo4j"
)

type RumahSakitStats struct {
	NamaRumahSakit    string
	JumlahTenagaMedis int
}

func getHospitalsByMedicalStaff() ([]RumahSakitStats, error) {
	// Count medical staff per hospital through departments
	query := `
		MATCH (rs:RumahSakit)-[:memiliki_departemen]->(d:Departemen)<-[:bekerja_di]-(t:TenagaMedis)
		WITH rs, COUNT(DISTINCT t) as jumlah_tenaga_medis
		RETURN rs.nama_rumah_sakit AS nama_rumah_sakit,
		       jumlah_tenaga_medis
		ORDER BY jumlah_tenaga_medis DESC
	`

	results, err := neo4j.ReadNeo4j(query, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to query Neo4j: %v", err)
	}

	// Process results
	hospitals := make([]RumahSakitStats, 0)
	for _, record := range results {
		hospitals = append(hospitals, RumahSakitStats{
			NamaRumahSakit:    record["nama_rumah_sakit"].(string),
			JumlahTenagaMedis: int(record["jumlah_tenaga_medis"].(int64)),
		})
	}

	return hospitals, nil
}

func displayHospitals(hospitals []RumahSakitStats, limit int) {
	fmt.Println("\n" + strings.Repeat("=", 70))
	fmt.Println("     RUMAH SAKIT DENGAN JUMLAH TENAGA MEDIS TERBANYAK")
	fmt.Println(strings.Repeat("=", 70))
	fmt.Printf("%-5s %-45s %s\n", "No", "Nama Rumah Sakit", "Jumlah Tenaga Medis")
	fmt.Println(strings.Repeat("-", 70))

	if len(hospitals) == 0 {
		fmt.Println("Tidak ada data rumah sakit.")
		return
	}

	displayLimit := limit
	if displayLimit > len(hospitals) {
		displayLimit = len(hospitals)
	}

	for i := 0; i < displayLimit; i++ {
		h := hospitals[i]
		fmt.Printf("%-5d %-45s %d\n", i+1, truncateString(h.NamaRumahSakit, 45), h.JumlahTenagaMedis)
	}

	fmt.Println(strings.Repeat("=", 70) + "\n")
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func main() {
	// Initialize Neo4j connection
	neo4j.InitNeo4j()
	defer neo4j.CloseNeo4j()

	// Execute query
	fmt.Println("\nFetching hospital statistics...")

	start := time.Now()
	hospitals, err := getHospitalsByMedicalStaff()
	elapsed := time.Since(start)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// Display results
	displayHospitals(hospitals, 10)

	// Timing summary
	fmt.Printf("\nTime: %.3f seconds (%d ms)\n", elapsed.Seconds(), elapsed.Milliseconds())
}
