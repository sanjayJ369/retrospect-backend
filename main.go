package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/sanjayj369/retrospect-backend/api"
	db "github.com/sanjayj369/retrospect-backend/db/sqlc"
)

const (
	driverName   = "postgres"
	driverSource = "postgresql://root:root@localhost:5432/retrospect?sslmode=disable"
	address      = "0.0.0.0:8080"
)

func main() {
	conn, err := pgx.Connect(context.Background(), driverSource)
	if err != nil {
		log.Fatal("unable to connect to db:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(address)
	if err != nil {
		log.Fatal("cannot start server, err:", err)
	}
}
