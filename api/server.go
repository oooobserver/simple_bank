package api

import (
	"fmt"
	db "simplebank/db/sqlc"
	"simplebank/token"
	"simplebank/util"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator"
)

// This struct serves HTTP requests
type Server struct {
	config util.Config

	// Use this to perform db operation
	store db.Store

	// Use this to create the token
	tokenMaker token.Maker
	router     *gin.Engine
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.SYMMETRIC_KEY)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	// Add validator middlewares
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	create_route(server)

	return server, nil
}

func create_route(server *Server) {
	router := gin.Default()

	// User
	router.POST("/user", server.createUser)
	// router.GET("/user/:name", server.getUser)

	// Login
	router.POST("/user/login", server.loginUser)

	// Create a group to add the auth middleware
	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	// Account
	authRoutes.POST("/account", server.createAccount)
	authRoutes.GET("/account/:id", server.getAccount)
	authRoutes.GET("/account", server.listAccounts)

	// Transfer
	authRoutes.POST("/transfer", server.createTransfer)

	server.router = router
}

// Start the server
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
