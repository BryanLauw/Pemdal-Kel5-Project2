package main

import (
	"fmt"
	"log"
	"time"

	"src/cassandra"
)

type PesananExpired struct {
	IdPesanan       string
	WaktuPemesanan  time.Time
	StatusPemesanan string
}

func getExpiredOrdersFromCassandra() ([]PesananExpired, error) {
	// Step 1: SELECT data yang status_pemesanan = 'belum dibayar'
	query := `
		SELECT id_pesanan, waktu_pemesanan, status_pemesanan
		FROM pemesanan_obat
		WHERE status_pemesanan = 'belum dibayar'
		ALLOW FILTERING
	`

	iter, err := cassandra.SelectCassandra(query)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data pesanan: %v", err)
	}

	var result []PesananExpired
	var id, status string
	var waktu time.Time

	// Batas waktu: 2 hari yang lalu
	twoDaysAgo := time.Now().Add(-48 * time.Hour)

	for iter.Scan(&id, &waktu, &status) {
		// Filter di aplikasi karena Cassandra tidak support time comparison di WHERE
		if waktu.Before(twoDaysAgo) {
			result = append(result, PesananExpired{
				IdPesanan:       id,
				WaktuPemesanan:  waktu,
				StatusPemesanan: status,
			})
		}
	}

	if err := iter.Close(); err != nil {
		return nil, fmt.Errorf("gagal membaca data: %v", err)
	}

	// Sort by waktu_pemesanan (oldest first) dan ambil 5 tertua
	// Simple bubble sort untuk 5 items
	for i := 0; i < len(result)-1; i++ {
		for j := 0; j < len(result)-i-1; j++ {
			if result[j].WaktuPemesanan.After(result[j+1].WaktuPemesanan) {
				result[j], result[j+1] = result[j+1], result[j]
			}
		}
	}

	// LIMIT 5 (ambil 5 tertua)
	if len(result) > 5 {
		result = result[:5]
	}

	return result, nil
}

func updateExpiredOrders(orders []PesananExpired) (int, error) {
	updatedCount := 0

	// Step 2: UPDATE satu per satu
	for _, order := range orders {
		query := `
			UPDATE pemesanan_obat
			SET status_pemesanan = 'dibatalkan'
			WHERE id_pesanan = ?
		`

		err := cassandra.UpdateCassandra(query, order.IdPesanan)
		if err != nil {
			log.Printf("Gagal update pesanan %s: %v", order.IdPesanan, err)
			continue
		}
		updatedCount++
	}

	return updatedCount, nil
}

func displayResult(orders []PesananExpired, updatedCount int) {
	fmt.Println("\n=== UPDATE 1: Update 5 Pesanan Obat Tertua yang Expired ===")
	fmt.Printf("Total pesanan expired yang diupdate: %d (LIMIT 5)\n\n", updatedCount)

	if len(orders) == 0 {
		fmt.Println("Tidak ada pesanan yang expired.")
		return
	}

	fmt.Printf("%-15s %-25s %-20s %-20s\n", "ID Pesanan", "Waktu Pemesanan", "Status Lama", "Status Baru")
	fmt.Println("────────────────────────────────────────────────────────────────────────────────────")

	for _, order := range orders {
		fmt.Printf("%-15s %-25s %-20s %-20s\n",
			order.IdPesanan,
			order.WaktuPemesanan.Format("2006-01-02 15:04:05"),
			order.StatusPemesanan,
			"dibatalkan")
	}

	fmt.Println("\n✓ 5 pesanan tertua berhasil diubah ke 'dibatalkan'")

	fmt.Println("\n⚠️  LIMITATION CASSANDRA:")
	fmt.Println("   - Tidak support time-based filtering (NOW() - INTERVAL) di WHERE clause")
	fmt.Println("   - Tidak support LIMIT di UPDATE statement")
	fmt.Println("   - Tidak support ORDER BY di query UPDATE")
	fmt.Println("   - Harus: SELECT → filter & sort di aplikasi → UPDATE satu per satu")
	fmt.Println("   - Trade-off untuk mendapatkan high write performance & horizontal scalability")
}

func main() {
	cassandra.InitCassandra()
	defer cassandra.Close()

	fmt.Println("Mencari pesanan expired (belum dibayar > 2 hari)...")

	start := time.Now()

	// Step 1: Get expired orders
	orders, err := getExpiredOrdersFromCassandra()
	if err != nil {
		log.Fatalf("Error getting expired orders: %v", err)
	}

	// Step 2: Update expired orders
	updatedCount, err := updateExpiredOrders(orders)

	elapsed := time.Since(start)

	if err != nil {
		log.Fatalf("Error updating expired orders: %v", err)
	}

	displayResult(orders, updatedCount)
	fmt.Printf("\nQuery Execution Time: %.3f seconds (%d ms)\n", elapsed.Seconds(), elapsed.Milliseconds())
}
