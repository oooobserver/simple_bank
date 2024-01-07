package db

import (
	"context"
	"log"
	"os"
	"simplebank/util"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

var testQueries *Queries
var testStore Store

func TestMain(m *testing.M) {
	// Load the config
	config, err := util.LoadConfig("../../.")
	if err != nil {
		log.Fatal("failed to read the config: ", err)
	}

	connPool, err := pgxpool.New(context.Background(), config.DB_SOURCE)
	if err != nil {
		log.Fatal("failed to connect to the db: ", err)
	}

	testQueries = New(connPool)
	testStore = NewStore(connPool)
	os.Exit(m.Run())
}
