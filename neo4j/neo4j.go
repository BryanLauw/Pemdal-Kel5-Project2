package neo4j

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

var (
	driver neo4j.DriverWithContext
	ctx    = context.Background()
)

// ====================================
// Init Neo4j connection
// ====================================
func InitNeo4j() {
	uri := getEnv("NEO4J_URI", "bolt://127.0.0.1:7687")
	user := getEnv("NEO4J_USER", "neo4j")
	pass := getEnv("NEO4J_PASSWORD", "password123")

	var err error
	driver, err = neo4j.NewDriverWithContext(uri, neo4j.BasicAuth(user, pass, ""))
	if err != nil {
		log.Fatalf("Failed to connect to Neo4j: %v", err)
	}
	fmt.Println("Connected to Neo4j")
}

// Close driver
func CloseNeo4j() {
	if driver != nil {
		driver.Close(ctx)
	}
}

// ====================================
// CRUD Functions
// ====================================

// Create Node or Relationship
func CreateNeo4j(query string, params map[string]interface{}) error {
	return runWrite(query, params)
}

// Create and Return data (for INSERT with RETURN clause)
func CreateAndReturnNeo4j(query string, params map[string]interface{}) ([]map[string]interface{}, error) {
	session := driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	result, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		res, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}

		var records []map[string]interface{}
		for res.Next(ctx) {
			records = append(records, res.Record().AsMap())
		}
		return records, res.Err()
	})

	if err != nil {
		return nil, err
	}
	return result.([]map[string]interface{}), nil
}

// Read / Query data
func ReadNeo4j(query string, params map[string]interface{}) ([]map[string]interface{}, error) {
	session := driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		res, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}

		var records []map[string]interface{}
		for res.Next(ctx) {
			records = append(records, res.Record().AsMap())
		}
		return records, res.Err()
	})

	if err != nil {
		return nil, err
	}
	return result.([]map[string]interface{}), nil
}

// Update Node
func UpdateNeo4j(query string, params map[string]interface{}) error {
	return runWrite(query, params)
}

// Delete Node or Relationship
func DeleteNeo4j(query string, params map[string]interface{}) error {
	return runWrite(query, params)
}

// --- Internal Helper ---
func runWrite(query string, params map[string]interface{}) error {
	session := driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, query, params)
		return nil, err
	})
	return err
}

// --- Utility ---
func getEnv(key, def string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return def
}
