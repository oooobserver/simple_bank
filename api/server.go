package api

import (
	db "simplebank/db/sqlc"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator"
)

// This struct serves HTTP requests
type Server struct {
	// config     util.Config

	// Use this to perform db operation
	store db.Store
	// tokenMaker token.Maker
	router *gin.Engine
}

func NewServer(store db.Store) (*Server, error) {
	server := &Server{
		store: store,
	}
	router := gin.Default()

	// Add validator middlewares
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	create_route(router, server)

	server.router = router
	return server, nil
}

func create_route(router *gin.Engine, server *Server) {
	// Account
	router.POST("/account", server.createAccount)
	router.GET("/account/:id", server.getAccount)
	router.GET("/account", server.listAccounts)

	// Transfer
	router.POST("/transfer", server.createTransfer)
}

// Start the server
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
