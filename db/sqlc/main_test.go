package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	dbSource = "postgresql://root:123@localhost:5432/root?sslmode=disable"
)

var testQueries *Queries
var testStore Store

func TestMain(m *testing.M) {
	connPool, err := pgxpool.New(context.Background(), dbSource)
	if err != nil {
		log.Fatal("failed to connect to the db: ", err)
	}

	testQueries = New(connPool)
	testStore = NewStore(connPool)
	os.Exit(m.Run())
}
