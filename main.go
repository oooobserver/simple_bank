package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"simplebank/api"
	db "simplebank/db/sqlc"
	"simplebank/gapi"
	"simplebank/pb"
	"simplebank/util"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
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

	// Start the GRPC server
	go runGRPCServer(config, store)

	// Start the Gatewway server
	runGatewayServer(config, store)
}

func runGRPCServer(con util.Config, store db.Store) {
	grpcServer := grpc.NewServer()
	server, err := gapi.NewServer(con, store)
	if err != nil {
		log.Fatal("failed to create the GRPC server: ", err)
	}
	pb.RegisterSimpleBankServer(grpcServer, server)
	// Registe the kind of the grpcServver
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", con.GRPC_ADDR)
	if err != nil {
		log.Fatal("failed to create the listener: ", err)
	}

	fmt.Println("start the gRPC server")

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("failed to listen at given port: ", err)
	}
}

func runGatewayServer(con util.Config, store db.Store) {
	server, err := gapi.NewServer(con, store)
	if err != nil {
		log.Fatal("failed to create the GRPC server: ", err)
	}

	grpcServeMux := runtime.NewServeMux()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = pb.RegisterSimpleBankHandlerServer(ctx, grpcServeMux, server)
	if err != nil {
		log.Fatal("can't register handle server: ", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcServeMux)

	listener, err := net.Listen("tcp", con.WEB_ADDR)
	if err != nil {
		log.Fatal("failed to create the listener: ", err)
	}

	log.Printf("start the gateway at %s", listener.Addr().String())

	err = http.Serve(listener, mux)
	if err != nil {
		log.Fatal("can't start gateway server: ", err)
	}
}

func runGinServer(con util.Config, store db.Store) {
	server, err := api.NewServer(con, store)
	if err != nil {
		log.Fatal("failed to create the Gin server: ", err)
	}

	// Start the server
	err = server.Start(con.WEB_ADDR)
	if err != nil {
		log.Fatal("failed to listen at given port: ", err)
	}
}
