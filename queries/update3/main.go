// =============================================================
// Query: Ubah status suatu pemesanan layanan menjadi 'dibatalkan'.
// =============================================================

package main

import (
	"fmt"
	"log"
	"time"

	"src/cassandra"
)

func main() {
	cassandra.InitCassandra()
	defer cassandra.Close()

	// 1. Ambil satu pemesanan layanan yang belum dibatalkan
	selectQuery := `SELECT id_pesanan, status_pemesanan FROM pemesanan_layanan ALLOW FILTERING`
	iter, err := cassandra.SelectCassandra(selectQuery)
	if err != nil {
		log.Fatalf("Gagal membaca data sebelum update: %v", err)
	}

	var idPesanan string
	var status string
	found := false
	for iter.Scan(&idPesanan, &status) {
		fmt.Printf("Ditemukan baris: id_pesanan=%s, status_pemesanan=%s\n", idPesanan, status)
		if status != "dibatalkan" {
			found = true
			break
		}
		idPesanan = ""
		status = ""
	}

	if cerr := iter.Close(); cerr != nil {
		log.Fatalf("Error saat menutup iterator: %v", cerr)
	}

	if !found {
		// Tidak ada record yang belum dibatalkan
		fmt.Println("Tidak ditemukan pemesanan layanan yang belum dibatalkan. Tidak ada yang diubah.")
		return
	}

	// Print record sebelum update
	fmt.Println("=== Sebelum Update ===")
	fmt.Printf("Record sebelum update: id_pesanan=%s, status_pemesanan=%s\n", idPesanan, status)

	// 2. Lakukan perubahan status
	start := time.Now()
	if err := BatalkanPemesananLayanan(idPesanan); err != nil {
		log.Fatalf("Gagal ubah status: %v", err)
	}
	duration := time.Since(start)
	fmt.Printf("\nPemesanan dibatalkan (%.2f ms)\n\n", float64(duration.Milliseconds()))

	// 3. Ambil kembali record itu dan print untuk melihat perubahan (pakai SelectCassandra kembali)
	fmt.Println("=== Setelah Update ===")
	iter2, err := cassandra.SelectCassandra("SELECT id_pesanan, status_pemesanan FROM pemesanan_layanan WHERE id_pesanan = ?", idPesanan)
	if err != nil {
		log.Fatalf("Gagal membaca data setelah update: %v", err)
	}
	var id2, status2 string
	if iter2.Scan(&id2, &status2) {
		fmt.Printf("Record setelah update: id_pesanan=%s, status_pemesanan=%s\n", id2, status2)
	} else {
		if cerr := iter2.Close(); cerr != nil {
			log.Fatalf("Error saat menutup iterator: %v", cerr)
		}
		fmt.Println("Record tidak ditemukan setelah update.")
	}
	if cerr := iter2.Close(); cerr != nil {
		log.Fatalf("Error saat menutup iterator: %v", cerr)
	}
}

func BatalkanPemesananLayanan(idPesanan string) error {
	query := `UPDATE pemesanan_layanan SET status_pemesanan = 'dibatalkan' WHERE id_pesanan = ?`
	return cassandra.UpdateCassandra(query, idPesanan)
}
