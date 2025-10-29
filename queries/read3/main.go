package main

import (
	"fmt"
	"log"
	"time"

	"src/cassandra"
	"src/neo4j"
)

type BayminLogs struct {
	Nama            string
	WaktuAktivitas  time.Time
	DetailAktivitas string
}

func getDevice(email string) (string, string, error) {
	query := `
		MATCH (p:Pasien {email: $email})-[:memiliki_perangkat]->(b:Baymin)
		RETURN b.id_perangkat AS id_perangkat, p.nama_lengkap AS nama
	`

	params := map[string]interface{}{"email": email}

	records, err := neo4j.ReadNeo4j(query, params)
	if err != nil {
		return "", "", fmt.Errorf("gagal membaca dari Neo4j: %v", err)
	}

	if len(records) == 0 {
		return "", "", fmt.Errorf("tidak ditemukan Baymin untuk pasien dengan email %s", email)
	}

	idPerangkat := fmt.Sprintf("%v", records[0]["id_perangkat"])
	nama := fmt.Sprintf("%v", records[0]["nama"])
	return idPerangkat, nama, nil
}

func getLogs(idPerangkat string, namaPasien string) ([]BayminLogs, error) {
	query := "SELECT waktu_aktivitas, detail_aktivitas FROM log_aktivitas WHERE id_perangkat = ?"

	iter, err := cassandra.SelectCassandra(query, idPerangkat)
	if err != nil {
		return nil, fmt.Errorf("gagal membaca dari Cassandra: %v", err)
	}
	defer iter.Close()

	var logs []BayminLogs
	var waktu time.Time
	var detail string

	for iter.Scan(&waktu, &detail) {
		logs = append(logs, BayminLogs{
			Nama:            namaPasien,
			WaktuAktivitas:  waktu,
			DetailAktivitas: detail,
		})
	}

	if err := iter.Close(); err != nil {
		return nil, fmt.Errorf("gagal menutup iterasi Cassandra: %v", err)
	}

	return logs, nil
}

func displayLogs(logs []BayminLogs) {
	fmt.Println("     Log Aktivitas Baymin per Pasien")
	fmt.Printf("%-25s %-25s %s\n", "Waktu Aktivitas", "Nama Pasien", "Detail Aktivitas")
	fmt.Println("----------------------------------------------------------------------------")

	if len(logs) == 0 {
		fmt.Println("Tidak ada log aktivitas untuk pasien ini.")
		return
	}

	for _, logData := range logs {
		fmt.Printf("%-25s %-25s %s\n",
			logData.WaktuAktivitas.Format("2006-01-02 15:04:05"),
			logData.Nama,
			logData.DetailAktivitas)
	}
}

func main() {
	email := "qXbIbDK@bInvBVI.net"

	cassandra.InitCassandra()
	defer cassandra.Close()

	neo4j.InitNeo4j()
	defer neo4j.CloseNeo4j()

	start := time.Now()

	idPerangkat, namaPasien, err := getDevice(email)
	if err != nil {
		log.Fatalf("Error Neo4j: %v", err)
	}

	logs, err := getLogs(idPerangkat, namaPasien)
	if err != nil {
		log.Fatalf("Error Cassandra: %v", err)
	}
	elapsed := time.Since(start)
	displayLogs(logs)

	fmt.Printf("\nTime: %.3f seconds (%d ms)\n", elapsed.Seconds(), elapsed.Milliseconds())
}
