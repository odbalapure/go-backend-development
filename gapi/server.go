package gapi

import (
	"fmt"
	db "simple-bank/db/sqlc"
	"simple-bank/pb"
	"simple-bank/token"
	"simple-bank/util"
)

type Server struct {
	// Default implementations for any methods you don't implement
	// The `mustEmbedUnimplementedSimpleBankServer()` method that satisfies the interface requirement
	// If someone adds a new RPC method to your `.proto` file,
	// Go code will still compile because `pb.UnimplementedSimpleBankServer` provides a default implementation for "new method".
	pb.UnimplementedSimpleBankServer
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
}

// Creates a new gRPC server and setup routing
func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	return server, nil
}
