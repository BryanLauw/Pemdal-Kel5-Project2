// =============================================================
// Query: Hapus semua pemesanan obat yang dibatalkan.
// =============================================================

package main

import (
	"fmt"
	"log"
	"time"

	"src/cassandra"
)

type PemesananObat struct {
	IDPesanan       string
	StatusPemesanan string
}

func main() {
	cassandra.InitCassandra()
	defer cassandra.Close()

	fmt.Println("=== Sebelum Delete ===")
	before, _ := getPemesananObatDibatalkan()
	fmt.Printf("Jumlah row sebelum dihapus: %d\n", len(before))

	start := time.Now()
	err := HapusPemesananObatDibatalkan()
	duration := time.Since(start)

	if err != nil {
		log.Fatalf("Gagal hapus pesanan dibatalkan: %v", err)
	}
	fmt.Printf("\nHapus selesai (%.2f ms)\n\n", float64(duration.Milliseconds()))

	fmt.Println("=== Setelah Delete ===")
	after, _ := getPemesananObatDibatalkan()
	fmt.Printf("Jumlah row setelah dihapus: %d\n", len(after))
}

func getPemesananObatDibatalkan() ([]PemesananObat, error) {
	query := `SELECT id_pesanan, status_pemesanan FROM pemesanan_obat WHERE status_pemesanan = 'dibatalkan' ALLOW FILTERING`

	iter, err := cassandra.SelectCassandra(query)
	if err != nil {
		return nil, err
	}

	var results []PemesananObat
	var pesanan PemesananObat
	for iter.Scan(&pesanan.IDPesanan, &pesanan.StatusPemesanan) {
		results = append(results, pesanan)
		pesanan = PemesananObat{}
	}

	if err := iter.Close(); err != nil {
		return nil, err
	}
	return results, nil
}

func HapusPemesananObatDibatalkan() error {
	// Ambil semua id_pesanan yang statusnya 'dibatalkan'
	selectQuery := `
		SELECT id_pesanan FROM pemesanan_obat WHERE status_pemesanan = 'dibatalkan' ALLOW FILTERING
	`
	iter, err := cassandra.SelectCassandra(selectQuery)
	if err != nil {
		return err
	}

	var ids []string
	var id string
	for iter.Scan(&id) {
		ids = append(ids, id)
		id = ""
	}

	// Periksa error iterator
	if err := iter.Close(); err != nil {
		return err
	}

	// Hapus satu per satu berdasarkan primary key id_pesanan
	for _, pid := range ids {
		delQuery := fmt.Sprintf("DELETE FROM pemesanan_obat WHERE id_pesanan = '%s'", pid)
		if err := cassandra.DeleteCassandra(delQuery); err != nil {
			return err
		}
	}

	return nil
}
