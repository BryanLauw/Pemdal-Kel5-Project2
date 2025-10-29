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

type LayananStats struct {
	NamaLayanan   string
	JumlahPesanan int
}

// ===============================================
// QUERY NEO4J
// ===============================================

func getMostOrderedServices() ([]LayananStats, error) {
	// Count appointments at hospitals that offer each service
	query := `
		MATCH (l:LayananMedis)<-[:menawarkan_layanan]-(rs:RumahSakit)<-[:di_rs]-(j:JanjiTemu)
		WITH l, COUNT(j) as jumlah_pesanan
		RETURN l.nama_layanan AS nama_layanan,
		       jumlah_pesanan
		ORDER BY jumlah_pesanan DESC
	`

	results, err := neo4j.ReadNeo4j(query, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to query Neo4j: %v", err)
	}

	// Process results
	services := make([]LayananStats, 0)
	for _, record := range results {
		services = append(services, LayananStats{
			NamaLayanan:   record["nama_layanan"].(string),
			JumlahPesanan: int(record["jumlah_pesanan"].(int64)),
		})
	}

	return services, nil
}

// ===============================================
// DISPLAY FUNCTION
// ===============================================

func displayServices(services []LayananStats, limit int) {
	fmt.Println("     LAYANAN MEDIS YANG PALING SERING DIPESAN")
	fmt.Printf("%-5s %-45s %s\n", "No", "Nama Layanan", "Jumlah Pesanan")

	if len(services) == 0 {
		fmt.Println("Tidak ada data layanan medis.")
		return
	}

	displayLimit := limit
	if displayLimit > len(services) {
		displayLimit = len(services)
	}

	for i := 0; i < displayLimit; i++ {
		s := services[i]
		fmt.Printf("%-5d %-45s %d\n", i+1, truncateString(s.NamaLayanan, 45), s.JumlahPesanan)
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
	fmt.Println("\nFetching service statistics...")
	start := time.Now()
	services, err := getMostOrderedServices()
	elapsed := time.Since(start)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// Display results
	displayServices(services, 5)

	fmt.Printf("\nTime: %.3f seconds (%d ms)\n", elapsed.Seconds(), elapsed.Milliseconds())
}
