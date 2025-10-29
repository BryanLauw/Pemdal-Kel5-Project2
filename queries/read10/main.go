package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"src/neo4j"
)

// ===============================================
// DATA STRUCTURE
// ===============================================

type PasienNoResep struct {
	Email           string
	NamaLengkap     string
	JumlahJanjiTemu int
}

// ===============================================
// QUERY NEO4J
// ===============================================

func getPatientsWithoutPrescriptions() ([]PasienNoResep, error) {
	// Find patients who have appointments but those appointments didn't produce prescriptions
	query := `
		MATCH (p:Pasien)<-[:memiliki_janji]-(j:JanjiTemu)
		WHERE NOT (j)-[:menghasilkan_resep]->(:Resep)
		WITH p, COUNT(DISTINCT j) as jumlah_janji_temu
		RETURN p.email AS email,
		       p.nama_lengkap AS nama_lengkap,
		       jumlah_janji_temu
		ORDER BY jumlah_janji_temu DESC, p.nama_lengkap ASC
		LIMIT 5
	`

	results, err := neo4j.ReadNeo4j(query, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to query Neo4j: %v", err)
	}

	// Process results
	patients := make([]PasienNoResep, 0)
	for _, record := range results {
		patients = append(patients, PasienNoResep{
			Email:           record["email"].(string),
			NamaLengkap:     record["nama_lengkap"].(string),
			JumlahJanjiTemu: int(record["jumlah_janji_temu"].(int64)),
		})
	}

	return patients, nil
}

// ===============================================
// DISPLAY FUNCTION
// ===============================================

func displayPatients(patients []PasienNoResep) {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("     PASIEN DENGAN JANJI TEMU TANPA RESEP")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Printf("%-5s %-35s %-30s %s\n", "No", "Nama Lengkap", "Email", "Jml Janji Temu")
	fmt.Println(strings.Repeat("-", 80))

	if len(patients) == 0 {
		fmt.Println("Tidak ada data pasien.")
		return
	}

	for i, p := range patients {
		fmt.Printf("%-5d %-35s %-30s %d\n",
			i+1,
			truncateString(p.NamaLengkap, 35),
			truncateString(p.Email, 30),
			p.JumlahJanjiTemu)
	}

	fmt.Println(strings.Repeat("=", 80) + "\n")
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
	fmt.Println("\nFetching patients without prescriptions...")

	start := time.Now()
	patients, err := getPatientsWithoutPrescriptions()
	elapsed := time.Since(start)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// Display results
	displayPatients(patients)

	// Summary
	fmt.Printf("Total patients with appointments but no prescriptions: %d\n", len(patients))

	// Timing summary
	fmt.Printf("\nTime: %.3f seconds (%d ms)\n", elapsed.Seconds(), elapsed.Milliseconds())
}
