package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"src/cassandra"
)

type PatientOrderCount struct {
	Email        string
	TotalPesanan int
}

func getPatientOrderCountFromCassandra() ([]PatientOrderCount, error) {
	query := "SELECT email_pemesan FROM pemesanan_obat"

	iter, err := cassandra.SelectCassandra(query)
	if err != nil {
		return nil, fmt.Errorf("Terjadi kesalahan: %v", err)
	}

	var email string
	orderCountMap := make(map[string]int)

	for iter.Scan(&email) {
		orderCountMap[email]++
	}

	if err := iter.Close(); err != nil {
		return nil, fmt.Errorf("gagal membaca data: %v", err)
	}

	result := make([]PatientOrderCount, 0, len(orderCountMap))
	for email, count := range orderCountMap {
		result = append(result, PatientOrderCount{
			Email:        email,
			TotalPesanan: count,
		})
	}

	return result, nil
}

func displayResult(patients []PatientOrderCount, limit int) {
	fmt.Println("\n" + strings.Repeat("=", 70))
	fmt.Println("     Jumlah Pesanan Obat per Pasien")
	fmt.Println(strings.Repeat("=", 70))
	fmt.Printf("%-5s %-40s %s\n", "No", "Email Pemesan", "Total Pesanan")
	fmt.Println(strings.Repeat("-", 70))

	if len(patients) == 0 {
		fmt.Println("Tidak ada data pemesanan obat.")
		return
	}

	displayLimit := limit
	if displayLimit > len(patients) {
		displayLimit = len(patients)
	}

	for i := 0; i < displayLimit; i++ {
		p := patients[i]
		fmt.Printf("%-5d %-40s %d\n",
			i+1,
			p.Email,
			p.TotalPesanan)
	}

	fmt.Println(strings.Repeat("=", 70) + "\n")
}

func main() {
	cassandra.InitCassandra()
	defer cassandra.Close()

	start := time.Now()
	patients, err := getPatientOrderCountFromCassandra()
	elapsed := time.Since(start)
	if err != nil {
		log.Fatalf("Error getting patient order costs: %v", err)
	}

	displayResult(patients, 10)
	fmt.Printf("\nTime: %.3f seconds (%d ms)\n", elapsed.Seconds(), elapsed.Milliseconds())
}
