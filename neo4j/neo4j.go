package neo4j

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

var Driver neo4j.DriverWithContext

func InitNeo4j() {
	uri := os.Getenv("NEO4J_URI")
	user := os.Getenv("NEO4J_USER")
	pass := os.Getenv("NEO4J_PASS")

	var err error
	Driver, err = neo4j.NewDriverWithContext(uri, neo4j.BasicAuth(user, pass, ""))
	if err != nil {
		log.Fatal("Neo4j connection error:", err)
	}
	fmt.Println("Connected to Neo4j")
}

// ===========================================================
// GENERIC CRUD EXECUTOR
// ===========================================================

// ExecNeo4j - CREATE / UPDATE / DELETE
func ExecNeo4j(ctx context.Context, query string, params map[string]interface{}) error {
	session := Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	_, err := session.Run(ctx, query, params)
	return err
}

// QueryNeo4j - SELECT / MATCH RETURN
func QueryNeo4j(ctx context.Context, query string, params map[string]interface{}) ([]map[string]interface{}, error) {
	session := Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	result, err := session.Run(ctx, query, params)
	if err != nil {
		return nil, err
	}

	records := []map[string]interface{}{}
	for result.Next(ctx) {
		row := make(map[string]interface{})
		for _, key := range result.Record().Keys {
			row[key], _ = result.Record().Get(key)
		}
		records = append(records, row)
	}
	return records, result.Err()
}
