package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	driverName   = "postgres"
	driverSource = "postgresql://root:root@localhost:5432/retrospect?sslmode=disable"
)

var testQueries *Queries

func TestMain(m *testing.M) {
	conn, err := sql.Open(driverName, driverSource)
	if err != nil {
		log.Fatal("unable to connect to db:", err)
	}

	testQueries = New(conn)

	os.Exit(m.Run())
}
