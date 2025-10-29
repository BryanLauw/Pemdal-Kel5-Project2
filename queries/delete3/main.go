// Query: Hapus semua janji temu yang sudah lewat lebih dari 30 hari dan masih belum selesai.

package main

import (
	"fmt"
	"log"
	"time"

	"src/neo4j"
)

func main() {
	neo4j.InitNeo4j()
	defer neo4j.CloseNeo4j()

	fmt.Println("=== Sebelum Delete ===")
	before, _ := neo4j.ReadNeo4j(`
		MATCH (j:JanjiTemu)
		WHERE datetime(replace(j.waktu_pelaksanaan, ' ', 'T')) < datetime() - duration('P30D')
			  AND NOT (j)-[:MENGHASILKAN_RESEP]->(:Resep)
		RETURN j.id_janji_temu AS id, j.waktu_pelaksanaan AS waktu
	`, nil)
	fmt.Printf("Jumlah sebelum: %d\n", len(before))

	start := time.Now()
	err := HapusJanjiTemuLama()
	duration := time.Since(start)
	if err != nil {
		log.Fatalf("Gagal hapus janji temu lama: %v", err)
	}
	fmt.Printf("Janji temu lama dihapus (%.2f ms)\n", float64(duration.Milliseconds()))

	fmt.Println("=== Setelah Delete ===")
	after, _ := neo4j.ReadNeo4j(`
		MATCH (j:JanjiTemu)
		WHERE datetime(replace(j.waktu_pelaksanaan, ' ', 'T')) < datetime() - duration('P30D')
			  AND NOT (j)-[:MENGHASILKAN_RESEP]->(:Resep)
		RETURN j.id_janji_temu AS id, j.waktu_pelaksanaan AS waktu
	`, nil)
	fmt.Printf("Jumlah setelah: %d\n", len(after))
}

func HapusJanjiTemuLama() error {
	query := `
		MATCH (j:JanjiTemu)
		WHERE datetime(replace(j.waktu_pelaksanaan, ' ', 'T')) < datetime() - duration('P30D')
			  AND NOT (j)-[:MENGHASILKAN_RESEP]->(:Resep)
		DETACH DELETE j
	`
	return neo4j.DeleteNeo4j(query, nil)
}
