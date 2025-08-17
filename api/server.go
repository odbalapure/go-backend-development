package api

import (
	"fmt"
	db "simple-bank/db/sqlc"
	"simple-bank/token"
	"simple-bank/util"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Server servers HTTP requests
type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	// Router will send each API request to correct handler
	router *gin.Engine
}

// Creates a new HTTP server and setup routing
func NewServer(config util.Config, store db.Store) (*Server, error) {
	// We can choose from JWT or Paseto that implement the same Maker interface
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	server.setupRouter()

	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()
	// Create user
	router.POST("/users", server.createUser)
	// Login
	router.POST("/users/login", server.loginUser)
	router.POST("/tokens/renew", server.renewAccessToken)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	// Create account
	authRoutes.POST("/accounts", server.createAccount)
	// Get an account
	authRoutes.GET("/accounts/:id", server.getAccount)
	// Get accounts
	authRoutes.GET("/accounts", server.getAccounts)
	// Transfer
	authRoutes.POST("/transfers", server.createTransfer)

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
