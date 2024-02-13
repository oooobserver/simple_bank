package main

import (
	"context"
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"net"
	"net/http"
	db "simplebank/db/sqlc"
	"simplebank/gapi"
	"simplebank/pb"
	"simplebank/util"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Read the config
	config, err := util.LoadConfig(".")
	if err != nil {
		msg := fmt.Sprint("failed to read the config: ", err)
		log.Fatal().Msg(msg)
	}

	if config.ENVIROMENT == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	// Create DB connection
	connPool, err := pgxpool.New(context.Background(), config.DB_SOURCE)
	if err != nil {
		msg := fmt.Sprint("failed to connect to the db: ", err)
		log.Fatal().Msg(msg)
	}

	// Run db migration
	runDBMigration(config.MIGRATION_URL, config.DB_SOURCE)

	store := db.NewStore(connPool)

	// Start the GRPC server
	go runGRPCServer(config, store)

	// Start the Gatewway server
	runGatewayServer(config, store)
}

func runDBMigration(mURL string, DBsource string) {
	migration, err := migrate.New(mURL, DBsource)
	if err != nil {
		msg := fmt.Sprint("can't create migration instance: ", err)
		log.Fatal().Msg(msg)
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		msg := fmt.Sprint("can't migrate up: ", err)
		log.Fatal().Msg(msg)
	}

	log.Info().Msg("db migrate up success!")
}

func runGRPCServer(con util.Config, store db.Store) {
	server, err := gapi.NewServer(con, store)
	if err != nil {
		msg := fmt.Sprint("failed to create the GRPC server: ", err)
		log.Fatal().Msg(msg)
	}

	grpcLogger := grpc.UnaryInterceptor(gapi.GrpcLogger)
	grpcServer := grpc.NewServer(grpcLogger)

	pb.RegisterSimpleBankServer(grpcServer, server)
	// Registe the kind of the grpcServver
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", con.GRPC_ADDR)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create listener")
	}

	log.Info().Msgf("start gRPC server at %s", listener.Addr().String())

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to listen at given port")
	}
}

func runGatewayServer(con util.Config, store db.Store) {
	server, err := gapi.NewServer(con, store)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create the GRPC server")
	}

	grpcServeMux := runtime.NewServeMux()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = pb.RegisterSimpleBankHandlerServer(ctx, grpcServeMux, server)
	if err != nil {
		log.Fatal().Err(err).Msg("can't register handle server")
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcServeMux)

	listener, err := net.Listen("tcp", con.WEB_ADDR)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create the listener")
	}

	log.Info().Msgf("start the gateway at %s", listener.Addr().String())

	err = http.Serve(listener, mux)
	if err != nil {
		log.Fatal().Err(err).Msg("can't start gateway server")
	}
}

// func runGinServer(con util.Config, store db.Store) {
// 	server, err := api.NewServer(con, store)
// 	if err != nil {
// 		log.Fatal().Err(err).Msg("failed to create the Gin server")
// 	}

// 	// Start the server
// 	err = server.Start(con.WEB_ADDR)
// 	if err != nil {
// 		log.Fatal().Err(err).Msg("failed to listen at given port")
// 	}
// }
