package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5"
	_ "github.com/lib/pq"
)

const (
	driverName   = "postgres"
	driverSource = "postgresql://root:root@localhost:5432/retrospect?sslmode=disable"
)

var testQueries *Queries

func TestMain(m *testing.M) {
	conn, err := pgx.Connect(context.Background(), driverSource)
	if err != nil {
		log.Fatal("unable to connect to db:", err)
	}

	testQueries = New(conn)

	os.Exit(m.Run())
}
