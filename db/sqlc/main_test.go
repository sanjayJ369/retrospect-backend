package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5"
	_ "github.com/lib/pq"
	"github.com/sanjayj369/retrospect-backend/util"
)

var testQueries *Queries

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config, err:", err)
	}
	conn, err := pgx.Connect(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("unable to connect to db:", err)
	}

	testQueries = New(conn)

	os.Exit(m.Run())
}
