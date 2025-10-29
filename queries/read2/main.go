package main

import (
	"fmt"
	"log"
	"time"

	"src/cassandra"
)

type MedicineStock struct {
	IDObat string
	Nama string
	Label string
	Stok int
}

func getMedicineStockFromCassandra() ([]MedicineStock, error) {
	query := "SELECT id_obat, nama, label, stok FROM obat WHERE stok < 55 ALLOW FILTERING"

	iter, err := cassandra.SelectCassandra(query)
	if err != nil {
		return nil, fmt.Errorf("Terjadi kesalahan: %v", err)
	}

	var medicines []MedicineStock

	var idObat, nama, label string
	var stok int

	for iter.Scan(&idObat, &nama, &label, &stok) {
		med := MedicineStock{
			IDObat: idObat,
			Nama:   nama,
			Label:  label,
			Stok:   stok,
		}
		medicines = append(medicines, med)
	}

	if err := iter.Close(); err != nil {
		return nil, fmt.Errorf("gagal membaca data: %v", err)
	}

	return medicines, nil
}

func displayResult(medicines []MedicineStock) {
	fmt.Println("     Daftar Obat dengan Stok Kurang dari 55")
	fmt.Printf("%-10s %-30s %-20s %s\n", "ID Obat", "Nama", "Label", "Stok")

	if len(medicines) == 0 {
		fmt.Println("Tidak ada data obat dengan stok kurang dari 55.")
		return
	}

	for _, m := range medicines {
		fmt.Printf("%-10s %-30s %-20s %d\n", // <- stok pakai %d (integer)
			m.IDObat,
			m.Nama,
			m.Label,
			m.Stok)
	}
}

func main() {
	cassandra.InitCassandra()
	defer cassandra.Close()

	start := time.Now()
	scanResult, err := getMedicineStockFromCassandra()
	elapsed := time.Since(start)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	displayResult(scanResult)
	fmt.Printf("\nTime: %.3f seconds (%d ms)\n", elapsed.Seconds(), elapsed.Milliseconds())
}