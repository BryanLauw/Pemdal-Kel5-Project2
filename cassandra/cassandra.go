package cassandra

import (
	"fmt"
	"log"
	"syscall"
	"github.com/gocql/gocql"
)

var Session *gocql.Session

// ====================================
// Init Cassandra connection
// ====================================
func InitCassandra() {
	cluster := gocql.NewCluster(getEnv("CASSANDRA_HOST", "127.0.0.1"))
	cluster.Port = getEnvInt("CASSANDRA_PORT", 9042)
	cluster.Keyspace = "rumahsakit"
	cluster.Consistency = gocql.Quorum

	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatalf("Cassandra connection failed: %v", err)
	}
	Session = session
	fmt.Println("Connected to Cassandra")
}

// Close session
func Close() {
	if Session != nil {
		Session.Close()
	}
}


// ====================================
// CRUD Functions
// ====================================

// Generic Query Executor (CQL)
func ExecCassandra(query string, params ...interface{}) error {
	return Session.Query(query, params...).Exec()
}

// Create / Insert
func InsertCassandra(query string, params ...interface{}) error {
	return ExecCassandra(query, params...)
}

// Read (SELECT)
func SelectCassandra(query string, params ...interface{}) (*gocql.Iter, error) {
	iter := Session.Query(query, params...).Iter()
	return iter, nil
}

// Update
func UpdateCassandra(query string, params ...interface{}) error {
	return ExecCassandra(query, params...)
}

// Delete
func DeleteCassandra(query string, params ...interface{}) error {
	return ExecCassandra(query, params...)
}

// --- Helper ---
func getEnv(key string, def string) string {
	val, ok := lookupEnv(key)
	if !ok {
		return def
	}
	return val
}

func getEnvInt(key string, def int) int {
	val, ok := lookupEnv(key)
	if !ok {
		return def
	}
	var port int
	fmt.Sscanf(val, "%d", &port)
	return port
}

func lookupEnv(key string) (string, bool) {
	val := ""
	ok := false
	if v, exists := syscall.Getenv(key); exists {
		val, ok = v, true
	}
	return val, ok
}
