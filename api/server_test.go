package api

import (
	db "simplebank/db/sqlc"
	"simplebank/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, store db.Store) *Server {
	config := util.Config{
		SYMMETRIC_KEY:   util.RandomString(32),
		ACCESS_DURATION: time.Minute,
	}

	server, err := NewServer(config, store)
	require.NoError(t, err)

	return server
}
