package cassandra

import (
	"fmt"
	"log"

	"github.com/gocql/gocql"
)

var Session *gocql.Session

func InitCassandra() {
	cluster := gocql.NewCluster("cassandra")
	cluster.Keyspace = "rumahsakit"
	cluster.Consistency = gocql.Quorum
	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatal("Cassandra connection error:", err)
	}
	Session = session
	fmt.Println("Connected to Cassandra")
}

// ===========================================================
// GENERIC CRUD EXECUTOR
// ===========================================================

// ExecCassandra - digunakan untuk CREATE, UPDATE, DELETE
func ExecCassandra(query string, params ...interface{}) error {
	if Session == nil {
		return fmt.Errorf("Cassandra session not initialized")
	}
	return Session.Query(query, params...).Exec()
}

// QueryCassandra - digunakan untuk SELECT, hasil dalam slice of map
func QueryCassandra(query string, params ...interface{}) ([]map[string]interface{}, error) {
	if Session == nil {
		return nil, fmt.Errorf("Cassandra session not initialized")
	}

	iter := Session.Query(query, params...).Iter()
	columns := iter.Columns()
	results := []map[string]interface{}{}
	row := make([]interface{}, len(columns))
	pointers := make([]interface{}, len(columns))

	for i := range row {
		pointers[i] = &row[i]
	}

	for iter.Scan(pointers...) {
		rowMap := make(map[string]interface{})
		for i, col := range columns {
			rowMap[col.Name] = row[i]
		}
		results = append(results, rowMap)
		row = make([]interface{}, len(columns))
		for i := range row {
			pointers[i] = &row[i]
		}
	}

	if err := iter.Close(); err != nil {
		return nil, err
	}
	return results, nil
}
