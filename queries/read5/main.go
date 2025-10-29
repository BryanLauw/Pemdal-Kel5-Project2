package main

import (
	"fmt"
	"log"
	"time"

	"src/cassandra"
	"src/neo4j"
)

type DetailResep struct {
	IDJanjiTemu string
	Penyakit    string
	NamaObat    string
	LabelObat   string
	Dosis       string
}

func getDetailResepFromNeo4j(idJanjiTemu string) ([]DetailResep, error) {
	query := `
		MATCH (jt:JanjiTemu {id_janji_temu: $id_janji_temu})
		      -[:menghasilkan_resep]->
		      (r:Resep)-[:memiliki_detail]->(dr:DetailResep)
		RETURN jt.id_janji_temu AS id_janji_temu,
		       r.penyakit AS penyakit,
		       dr.id_obat AS id_obat,
		       dr.dosis AS dosis
	`

	params := map[string]interface{}{"id_janji_temu": idJanjiTemu}
	records, err := neo4j.ReadNeo4j(query, params)
	if err != nil {
		return nil, fmt.Errorf("gagal membaca dari Neo4j: %v", err)
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("tidak ditemukan detail resep untuk janji temu %s", idJanjiTemu)
	}

	var results []DetailResep
	for _, record := range records {
		idObat := fmt.Sprintf("%v", record["id_obat"])
		var namaObat, labelObat string

		iter, err := cassandra.SelectCassandra(
			"SELECT nama, label FROM obat WHERE id_obat = ?", idObat,
		)
		if err == nil && iter.Scan(&namaObat, &labelObat) {
			iter.Close()
		} else {
			namaObat = "-"
			labelObat = "-"
		}

		results = append(results, DetailResep{
			IDJanjiTemu: fmt.Sprintf("%v", record["id_janji_temu"]),
			Penyakit:    fmt.Sprintf("%v", record["penyakit"]),
			NamaObat:    namaObat,
			LabelObat:   labelObat,
			Dosis:       fmt.Sprintf("%v", record["dosis"]),
		})
	}

	return results, nil
}

func displayResult(details []DetailResep) {
	if len(details) == 0 {
		fmt.Println("Tidak ada data detail resep.")
		return
	}

	fmt.Printf("\nDetail Resep untuk Janji Temu %s:\n", details[0].IDJanjiTemu)
	fmt.Printf("%-25s %-25s %-25s %s\n", "Penyakit", "Nama Obat", "Label Obat", "Dosis")
	fmt.Println("-------------------------------------------------------------------------------")

	for _, d := range details {
		fmt.Printf("%-25s %-25s %-25s %s\n",
			d.Penyakit, d.NamaObat, d.LabelObat, d.Dosis)
	}
}

func main() {
	idJanjiTemu := "JT00011"

	neo4j.InitNeo4j()
	defer neo4j.CloseNeo4j()

	cassandra.InitCassandra()
	defer cassandra.Close()

	start := time.Now()
	results, err := getDetailResepFromNeo4j(idJanjiTemu)
	elapsed := time.Since(start)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	displayResult(results)
	fmt.Printf("\nTime: %.3f seconds (%d ms)\n", elapsed.Seconds(), elapsed.Milliseconds())
}
