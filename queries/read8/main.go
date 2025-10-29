package main

import (
	"fmt"
	"log"
	"time"

	"src/neo4j"
)

// ===============================================
// DATA STRUCTURE
// ===============================================

type RumahSakitStats struct {
	NamaRumahSakit  string
	JumlahJanjiTemu int
}

// ===============================================
// QUERY NEO4J
// ===============================================

func getTopHospitalsByAppointments() ([]RumahSakitStats, error) {
	// Count appointments per hospital
	query := `
		MATCH (rs:RumahSakit)<-[:di_rs]-(j:JanjiTemu)
		WITH rs, COUNT(j) as jumlah_janji_temu
		RETURN rs.nama_rumah_sakit AS nama_rumah_sakit,
		       jumlah_janji_temu
		ORDER BY jumlah_janji_temu DESC
		LIMIT 10
	`

	results, err := neo4j.ReadNeo4j(query, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to query Neo4j: %v", err)
	}

	// Process results
	hospitals := make([]RumahSakitStats, 0)
	for _, record := range results {
		hospitals = append(hospitals, RumahSakitStats{
			NamaRumahSakit:  record["nama_rumah_sakit"].(string),
			JumlahJanjiTemu: int(record["jumlah_janji_temu"].(int64)),
		})
	}

	return hospitals, nil
}

// ===============================================
// DISPLAY FUNCTION
// ===============================================

func displayHospitals(hospitals []RumahSakitStats) {
	fmt.Println("     10 RUMAH SAKIT DENGAN JUMLAH JANJI TEMU TERBANYAK")
	fmt.Printf("%-5s %-45s %s\n", "No", "Nama Rumah Sakit", "Jumlah Janji Temu")

	if len(hospitals) == 0 {
		fmt.Println("Tidak ada data rumah sakit.")
		return
	}

	for i, h := range hospitals {
		fmt.Printf("%-5d %-45s %d\n", i+1, truncateString(h.NamaRumahSakit, 45), h.JumlahJanjiTemu)
	}

}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// ===============================================
// MAIN FUNCTION
// ===============================================

func main() {
	// Initialize Neo4j connection
	fmt.Println("Initializing Neo4j connection...")
	neo4j.InitNeo4j()
	defer neo4j.CloseNeo4j()

	// Execute query
	fmt.Println("\nFetching hospital statistics...")

	start := time.Now()
	hospitals, err := getTopHospitalsByAppointments()
	elapsed := time.Since(start)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// Display results
	displayHospitals(hospitals)

	// Timing summary
	fmt.Printf("\nTime: %.3f seconds (%d ms)\n", elapsed.Seconds(), elapsed.Milliseconds())
}
