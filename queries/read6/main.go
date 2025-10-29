package main

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"src/cassandra"
)

type PatientOrderCost struct {
	Email      string
	TotalBiaya float64
}

func getPatientOrderCosts() ([]PatientOrderCost, error) {
	// Step 1: Get all orders
	query := `SELECT id_pesanan, email_pemesan FROM rumahsakit.pemesanan_obat`
	iter, err := cassandra.SelectCassandra(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query pemesanan_obat: %v", err)
	}
	defer iter.Close()

	// Collect all order IDs first
	type OrderInfo struct {
		ID    string
		Email string
	}
	orders := make([]OrderInfo, 0)
	var idPesanan, emailPemesan string
	for iter.Scan(&idPesanan, &emailPemesan) {
		orders = append(orders, OrderInfo{ID: idPesanan, Email: emailPemesan})
	}

	// Step 2: Cache all medication prices once
	priceCache := make(map[string]float64)
	priceQuery := `SELECT id_obat, harga FROM rumahsakit.obat`
	priceIter, err := cassandra.SelectCassandra(priceQuery)
	if err != nil {
		return nil, err
	}
	defer priceIter.Close()

	var idObat string
	var harga float64
	for priceIter.Scan(&idObat, &harga) {
		priceCache[idObat] = harga
	}

	// Step 3: Process orders using cached prices
	patientMap := make(map[string]float64)
	for _, order := range orders {
		detailQuery := `SELECT daftar_obat FROM rumahsakit.detail_pesanan_obat WHERE id_pesanan = ?`
		detailIter, err := cassandra.SelectCassandra(detailQuery, order.ID)
		if err != nil {
			log.Printf("Error getting details for order %s: %v", order.ID, err)
			continue
		}

		var daftarObat map[string]int
		if detailIter.Scan(&daftarObat) {
			for obatID, jumlah := range daftarObat {
				// Use cached price instead of querying
				harga := priceCache[obatID]
				subtotal := harga * float64(jumlah)
				patientMap[order.Email] += subtotal
			}
		}
		detailIter.Close()
	}

	// Convert to slice and sort
	patients := make([]PatientOrderCost, 0, len(patientMap))
	for email, totalBiaya := range patientMap {
		patients = append(patients, PatientOrderCost{
			Email:      email,
			TotalBiaya: totalBiaya,
		})
	}

	sort.Slice(patients, func(i, j int) bool {
		return patients[i].TotalBiaya > patients[j].TotalBiaya
	})

	return patients, nil
}

func displayTopPatients(patients []PatientOrderCost, limit int) {
	fmt.Println("\n" + strings.Repeat("=", 70))
	fmt.Println("     PASIEN DENGAN BIAYA PEMESANAN OBAT TERBESAR")
	fmt.Println(strings.Repeat("=", 70))
	fmt.Printf("%-5s %-40s %s\n", "No", "Email Pemesan", "Total Biaya")
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
		fmt.Printf("%-5d %-40s Rp %s\n",
			i+1,
			p.Email,
			formatRupiah(p.TotalBiaya))
	}

	fmt.Println(strings.Repeat("=", 70) + "\n")
}

func formatRupiah(amount float64) string {
	str := fmt.Sprintf("%.2f", amount)

	var intPart, decPart string
	for i := len(str) - 1; i >= 0; i-- {
		if str[i] == '.' {
			intPart = str[:i]
			decPart = str[i:]
			break
		}
	}

	var result string
	for i := len(intPart) - 1; i >= 0; i-- {
		if len(result) > 0 && (len(intPart)-i-1)%3 == 0 {
			result = "." + result
		}
		result = string(intPart[i]) + result
	}

	return result + decPart
}

func main() {
	cassandra.InitCassandra()
	defer cassandra.Close()

	start := time.Now()
	patients, err := getPatientOrderCosts()
	elapsed := time.Since(start)
	if err != nil {
		log.Fatalf("Error getting patient order costs: %v", err)
	}

	displayTopPatients(patients, 5)
	fmt.Printf("\nTime: %.3f seconds (%d ms)\n", elapsed.Seconds(), elapsed.Milliseconds())
}
