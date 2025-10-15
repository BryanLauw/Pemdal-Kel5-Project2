package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gocql/gocql"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func testCassandra() {
	fmt.Println("\n=== 🧩 Testing Cassandra CRUD ===")

	// Gunakan "cassandra" kalau pakai Docker Compose
	cluster := gocql.NewCluster("localhost")
	cluster.Port = 9042
	cluster.Consistency = gocql.Quorum
	cluster.Timeout = 5 * time.Second

	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatal("❌ Cassandra connect error:", err)
	}
	defer session.Close()
	fmt.Println("✅ Connected to Cassandra")

	// Buat keyspace
	err = session.Query(`
		CREATE KEYSPACE IF NOT EXISTS test_app
		WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1}
	`).Exec()
	if err != nil {
		log.Println("⚠️ Keyspace creation warning:", err)
	}

	// Ganti keyspace
	cluster.Keyspace = "test_app"
	session, _ = cluster.CreateSession()
	defer session.Close()

	// Buat tabel
	err = session.Query(`
		CREATE TABLE IF NOT EXISTS users (
			email text PRIMARY KEY,
			name text,
			city text
		)
	`).Exec()
	if err != nil {
		log.Fatal("❌ Create table error:", err)
	}

	// CREATE
	if err := session.Query(
		"INSERT INTO users (email, name, city) VALUES (?, ?, ?)",
		"budi@example.com", "Budi", "Bandung",
	).Exec(); err != nil {
		log.Fatal("❌ Insert error:", err)
	}
	fmt.Println("✅ Data inserted")

	// READ
	iter := session.Query("SELECT email, name, city FROM users").Iter()
	var email, name, city string
	for iter.Scan(&email, &name, &city) {
		fmt.Printf("📄 Read: %s | %s | %s\n", email, name, city)
	}
	iter.Close()

	// UPDATE
	if err := session.Query(
		"UPDATE users SET city=? WHERE email=?",
		"Jakarta", "budi@example.com",
	).Exec(); err != nil {
		log.Fatal("❌ Update error:", err)
	}
	fmt.Println("✏️ Data updated")

	// READ AGAIN
	iter = session.Query("SELECT email, name, city FROM users WHERE email=?", "budi@example.com").Iter()
	for iter.Scan(&email, &name, &city) {
		fmt.Printf("📄 After update: %s | %s | %s\n", email, name, city)
	}
	iter.Close()

	// DELETE
	if err := session.Query("DELETE FROM users WHERE email=?", "budi@example.com").Exec(); err != nil {
		log.Fatal("❌ Delete error:", err)
	}
	fmt.Println("🗑️ Data deleted")

	fmt.Println("✅ Cassandra test done.")
}

func testNeo4j() {
	fmt.Println("\n=== 🧩 Testing Neo4j CRUD ===")

	// Gunakan "neo4j" kalau pakai Docker Compose
	uri := "bolt://localhost:7687"
	auth := neo4j.BasicAuth("neo4j", "password123", "")

	driver, err := neo4j.NewDriverWithContext(uri, auth)
	if err != nil {
		log.Fatal("❌ Neo4j connect error:", err)
	}
	defer driver.Close(context.Background())

	ctx := context.Background()
	session := driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)

	// CREATE
	_, err = session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, `
			CREATE (u:User {email: $email, name: $name, city: $city})
		`, map[string]any{
			"email": "andi@example.com",
			"name":  "Andi",
			"city":  "Surabaya",
		})
		return nil, err
	})
	if err != nil {
		log.Fatal("❌ Neo4j insert error:", err)
	}
	fmt.Println("✅ Node created")

	// READ
	_, _ = session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, _ := tx.Run(ctx, "MATCH (u:User) RETURN u.email, u.name, u.city", nil)
		for result.Next(ctx) {
			record := result.Record()
			fmt.Printf("📄 Read: %s | %s | %s\n",
				record.Values[0], record.Values[1], record.Values[2])
		}
		return nil, nil
	})

	// UPDATE
	_, _ = session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, `
			MATCH (u:User {email: $email})
			SET u.city = $city
		`, map[string]any{"email": "andi@example.com", "city": "Yogyakarta"})
		return nil, err
	})
	fmt.Println("✏️ Node updated")

	// DELETE
	_, _ = session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, "MATCH (u:User {email: $email}) DELETE u", map[string]any{"email": "andi@example.com"})
		return nil, err
	})
	fmt.Println("🗑️ Node deleted")

	fmt.Println("✅ Neo4j test done.")
}

func main() {
	testCassandra()
	testNeo4j()
}
