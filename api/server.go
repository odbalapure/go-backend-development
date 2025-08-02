package api

import (
	db "simple-bank/db/sqlc"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Server servers HTTP requests
type Server struct {
	store db.Store
	// Router will send each API request to correct handler
	router *gin.Engine
}

// Creates a new HTTP server and setup routing
func NewServer(store db.Store) *Server {
	// Create server
	server := &Server{store: store}
	// Create router
	router := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	// Create user
	router.POST("/users", server.createUser)
	// Create account
	router.POST("/accounts", server.createAccount)
	// Get an account
	router.GET("/accounts/:id", server.getAccount)
	// Get accounts
	router.GET("/accounts", server.getAccounts)
	// Transfer
	router.POST("/transfers", server.createTransfer)

	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
