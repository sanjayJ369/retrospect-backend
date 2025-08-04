package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/sanjayj369/retrospect-backend/api"
	db "github.com/sanjayj369/retrospect-backend/db/sqlc"
	mail "github.com/sanjayj369/retrospect-backend/mail"

	"github.com/sanjayj369/retrospect-backend/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config, err:", err)
	}
	conn, err := pgx.Connect(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("unable to connect to db:", err)
	}

	store := db.NewStore(conn)
	mailSender, err := mail.NewMailgunSender(config)
	if err != nil {
		log.Fatal("cannot create email sender, err:", err)
	}

	server, err := api.NewServer(config, store, mailSender)
	if err != nil {
		log.Fatal("cannot create server, err:", err)
	}
	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server, err:", err)
	}
}
