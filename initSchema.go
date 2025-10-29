package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gocql/gocql"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// ===============================================
// SCHEMA CASSANDRA
// ===============================================
func createCassandraSchema() {
	host := getEnv("CASSANDRA_HOST", "127.0.0.1")
	port := getEnvInt("CASSANDRA_PORT", 9042)

	fmt.Println("Creating Cassandra keyspace and tables ...")

	// Koneksi sementara tanpa keyspace
	tempCluster := gocql.NewCluster(host)
	tempCluster.Port = port
	tempCluster.Consistency = gocql.Quorum

	tempSession, err := tempCluster.CreateSession()
	if err != nil {
		log.Fatalf("Cassandra initial connection failed: %v", err)
	}
	defer tempSession.Close()

	// Buat keyspace
	err = tempSession.Query(`
		CREATE KEYSPACE IF NOT EXISTS rumahsakit
		WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1};
	`).Exec()
	if err != nil {
		log.Fatalf("Failed to create keyspace rumahsakit: %v", err)
	}
	fmt.Println("Keyspace 'rumahsakit' ready.")

	// Koneksi ulang ke keyspace rumahsakit
	cluster := gocql.NewCluster(host)
	cluster.Port = port
	cluster.Keyspace = "rumahsakit"
	cluster.Consistency = gocql.Quorum

	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatalf("Cassandra connection to keyspace failed: %v", err)
	}
	defer session.Close()

	// Buat tabel-tabel
	queries := []string{
		`CREATE TABLE IF NOT EXISTS log_aktivitas (
			id_perangkat TEXT,
			waktu_aktivitas TIMESTAMP,
			detail_aktivitas TEXT,
			PRIMARY KEY ((id_perangkat), waktu_aktivitas)
		) WITH CLUSTERING ORDER BY (waktu_aktivitas DESC);`,

		`CREATE TABLE IF NOT EXISTS pemesanan_obat (
			id_pesanan TEXT PRIMARY KEY,
			email_pemesan TEXT,
			waktu_pemesanan TIMESTAMP,
			status_pemesanan TEXT
		);`,

		`CREATE TABLE IF NOT EXISTS detail_pesanan_obat (
			id_pesanan TEXT PRIMARY KEY,
			daftar_obat MAP<TEXT, INT>
		);`,

		`CREATE TABLE IF NOT EXISTS obat (
			id_obat TEXT PRIMARY KEY,
			nama TEXT,
			label TEXT,
			harga DOUBLE,
			stok INT
		);`,

		`CREATE TABLE IF NOT EXISTS pemesanan_layanan (
			id_pesanan TEXT PRIMARY KEY,
			email_pemesan TEXT,
			waktu_pemesanan TIMESTAMP,
			jadwal_pelaksanaan TIMESTAMP,
			status_pemesanan TEXT
		);`,

		`CREATE TABLE IF NOT EXISTS lokasi_layanan (
			id_rs TEXT,
			id_layanan TEXT,
			nama_layanan TEXT,
			biaya_layanan DOUBLE,
			PRIMARY KEY (id_rs, id_layanan)
		);`,
	}

	for _, q := range queries {
		if err := session.Query(q).Exec(); err != nil {
			log.Println("Error executing query:", err)
		}
	}

	fmt.Println("Cassandra denormalized schema created successfully.")
}

// ===============================================
// BAGIAN 2 — SCHEMA NEO4J
// ===============================================
func createNeo4jSchema() {
	fmt.Println("Creating Neo4j constraints and relationships...")

	uri := getEnv("NEO4J_URI", "bolt://localhost:7687")
	user := getEnv("NEO4J_USER", "neo4j")
	pass := getEnv("NEO4J_PASSWORD", "password123")

	driver, err := neo4j.NewDriverWithContext(uri, neo4j.BasicAuth(user, pass, ""))
	if err != nil {
		log.Fatalf("Cannot connect to Neo4j: %v", err)
	}
	defer driver.Close(context.Background())

	session := driver.NewSession(context.Background(), neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(context.Background())

	queries := []string{
		"CREATE CONSTRAINT IF NOT EXISTS FOR (p:Pasien) REQUIRE p.email IS UNIQUE;",
		"CREATE CONSTRAINT IF NOT EXISTS FOR (t:TenagaMedis) REQUIRE t.email IS UNIQUE;",
		"CREATE CONSTRAINT IF NOT EXISTS FOR (r:RumahSakit) REQUIRE r.id_rs IS UNIQUE;",
		"CREATE CONSTRAINT IF NOT EXISTS FOR (d:Departemen) REQUIRE d.nama_departemen IS UNIQUE;",
		"CREATE CONSTRAINT IF NOT EXISTS FOR (l:LayananMedis) REQUIRE l.id_layanan IS UNIQUE;",
		"CREATE CONSTRAINT IF NOT EXISTS FOR (b:Baymin) REQUIRE b.id_perangkat IS UNIQUE;",
		"CREATE CONSTRAINT IF NOT EXISTS FOR (j:JanjiTemu) REQUIRE j.id_janji_temu IS UNIQUE;",
		"CREATE CONSTRAINT IF NOT EXISTS FOR (r:Resep) REQUIRE r.id_resep IS UNIQUE;",
		"CREATE CONSTRAINT IF NOT EXISTS FOR (dr:DetailResep) REQUIRE dr.id_obat IS UNIQUE;",
	}

	for _, q := range queries {
		_, err := session.ExecuteWrite(context.Background(),
			func(tx neo4j.ManagedTransaction) (any, error) {
				_, err := tx.Run(context.Background(), q, nil)
				return nil, err
			})
		if err != nil {
			log.Printf("Neo4j Query failed: %s\nError: %v\n", q, err)
		}
	}

	fmt.Println("Neo4j constraints created successfully.")
}

// ===============================================
// MAIN FUNCTION
// ===============================================
func main() {
	createCassandraSchema()
	createNeo4jSchema()
}

// ===============================================
// HELPER FUNCTIONS
// ===============================================
func getEnv(key, def string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return def
}

func getEnvInt(key string, def int) int {
	if val, ok := os.LookupEnv(key); ok {
		var i int
		fmt.Sscanf(val, "%d", &i)
		return i
	}
	return def
}
