package main

import (
	"context"
	"log"
	"simplebank/api"
	db "simplebank/db/sqlc"
	"simplebank/util"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// try again
	// Read the config
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("failed to read the config: ", err)
	}

	// Create DB connection
	connPool, err := pgxpool.New(context.Background(), config.DB_SOURCE)
	if err != nil {
		log.Fatal("failed to connect to the db: ", err)
	}

	store := db.NewStore(connPool)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("failed to create the server: ", err)
	}

	// Start the server
	err = server.Start(config.WEB_ADDR)
	if err != nil {
		log.Fatal("failed to listen at given port: ", err)
	}

}
