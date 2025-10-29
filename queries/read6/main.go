package main

import (
	"fmt"
	"log"
	"sort"
	"time"

	"src/cassandra"
)

type PatientOrderCost struct {
	Email      string
	TotalBiaya float64
}

func getPatientOrderCostsFromCassandra() ([]PatientOrderCost, error) {
	// Step 1: Get all medication orders
	query := `SELECT id_pesanan, email_pemesan FROM rumahsakit.pemesanan_obat`

	iter, err := cassandra.SelectCassandra(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query pemesanan_obat: %v", err)
	}
	defer iter.Close()

	patientMap := make(map[string]float64)

	var idPesanan, emailPemesan string

	for iter.Scan(&idPesanan, &emailPemesan) {
		// Step 2: Get order details (list of medications)
		detailQuery := `SELECT daftar_obat FROM rumahsakit.detail_pesanan_obat WHERE id_pesanan = ?`
		detailIter, err := cassandra.SelectCassandra(detailQuery, idPesanan)
		if err != nil {
			log.Printf("Error getting details for order %s: %v", idPesanan, err)
			continue
		}

		var daftarObat map[string]int
		if detailIter.Scan(&daftarObat) {
			// Step 3: Calculate costs
			for idObat, jumlah := range daftarObat {
				harga := getObatPrice(idObat)
				subtotal := harga * float64(jumlah)
				patientMap[emailPemesan] += subtotal
			}
		}
		detailIter.Close()
	}

	if err := iter.Close(); err != nil {
		return nil, err
	}

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

func getObatPrice(idObat string) float64 {
	query := `SELECT harga FROM rumahsakit.obat WHERE id_obat = ? LIMIT 1`
	iter, err := cassandra.SelectCassandra(query, idObat)
	if err != nil {
		log.Printf("Error querying obat %s: %v", idObat, err)
		return 0
	}
	defer iter.Close()

	var harga float64
	if iter.Scan(&harga) {
		return harga
	}
	return 0
}

func displayTopPatients(patients []PatientOrderCost, limit int) {
	fmt.Println("     PASIEN DENGAN BIAYA PEMESANAN OBAT TERBESAR")
	fmt.Printf("%-5s %-40s %s\n", "No", "Email Pemesan", "Total Biaya")

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
	fmt.Println("Initializing Cassandra connection...")
	cassandra.InitCassandra()
	defer cassandra.Close()

	fmt.Println("Fetching patient medication order costs...")

	start := time.Now()
	patients, err := getPatientOrderCostsFromCassandra()
	elapsed := time.Since(start)
	if err != nil {
		log.Fatalf("Error getting patient order costs: %v", err)
	}

	displayTopPatients(patients, 5)
	fmt.Printf("\nTime: %.3f seconds (%d ms)\n", elapsed.Seconds(), elapsed.Milliseconds())
}
