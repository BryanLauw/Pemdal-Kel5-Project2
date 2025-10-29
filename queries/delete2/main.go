// Query: Hapus semua log aktivitas Baymin yang lebih tua dari 6 bulan.

package main

import (
	"fmt"
	"log"
	"time"

	"src/cassandra"
)

type LogAktivitas struct {
	IDPerangkat     string
	WaktuAktivitas  time.Time
	DetailAktivitas string
}

func main() {
	cassandra.InitCassandra()
	defer cassandra.Close()

	fmt.Println("=== Sebelum Delete ===")
	before, _ := getLogAktivitasLama()
	fmt.Printf("Jumlah row sebelum dihapus: %d\n", len(before))

	start := time.Now()
	err := HapusLogAktivitasLama()
	duration := time.Since(start)
	if err != nil {
		log.Fatalf("Gagal hapus log lama: %v", err)
	}
	fmt.Printf("\nLog lama dihapus (%.2f ms)\n\n", float64(duration.Milliseconds()))

	fmt.Println("=== Setelah Delete ===")
	after, _ := getLogAktivitasLama()
	fmt.Printf("Jumlah row setelah dihapus: %d\n", len(after))
}

func getLogAktivitasLama() ([]LogAktivitas, error) {
	sixMonthsAgo := time.Now().AddDate(0, -6, 0)
	query := `SELECT id_perangkat, waktu_aktivitas, detail_aktivitas FROM log_aktivitas WHERE waktu_aktivitas < ? ALLOW FILTERING`

	iter, err := cassandra.SelectCassandra(query, sixMonthsAgo)
	if err != nil {
		return nil, err
	}
	defer iter.Close()

	var logs []LogAktivitas
	var log LogAktivitas
	for iter.Scan(&log.IDPerangkat, &log.WaktuAktivitas, &log.DetailAktivitas) {
		logs = append(logs, log)
	}

	if err := iter.Close(); err != nil {
		return nil, err
	}

	return logs, nil
}

func HapusLogAktivitasLama() error {
	sixMonthsAgo := time.Now().AddDate(0, -6, 0)

	query := `SELECT id_perangkat, waktu_aktivitas FROM log_aktivitas`
	iter, err := cassandra.SelectCassandra(query)
	if err != nil {
		return err
	}

	var idPerangkat string
	var waktu time.Time
	for iter.Scan(&idPerangkat, &waktu) {
		if waktu.Before(sixMonthsAgo) {
			delQuery := `
				DELETE FROM log_aktivitas 
				WHERE id_perangkat = ? AND waktu_aktivitas = ?
			`
			if err := cassandra.DeleteCassandra(delQuery, idPerangkat, waktu); err != nil {
				log.Printf("Gagal hapus log untuk %s: %v\n", idPerangkat, err)
			}
		}
	}
	return iter.Close()
}
