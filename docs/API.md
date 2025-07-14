## API Development

Most popular web frameworks for Golang:
- Gin
- Beego
- Echo
- Revel
- Martini
- Fiber
- Buffalo

Features - routing, parameter binding, validation, middle ware, some have built in ROM

If we prefer a light weight framework with support for routing only:
- FastHttp
- Groilla Mux
- HttpRouter
- Chi

[Gin](https://github.com/gin-gonic/gin) is the most popular one with most stars.

Install the gin server
```sh
go get -u github.com/gin-gonic/gin
```

### Create server

The `server.go` file can contain methods for starting the server, returning error response, call methods on `router`.

```go
package api

import (
	db "simple-bank/db/sqlc"
	"github.com/gin-gonic/gin"
)

// Server servers HTTP requests
type Server struct {
	store *db.Store
	// Router will send each API request to correct handler
	router *gin.Engine
}

// Creates a new HTTP server and setup routing
func NewServer(store *db.Store) *Server {
	// Create server
	server := &Server{store: store}
	// Create router
	router := gin.Default()

	// Create account
	router.POST("/accounts", server.createAccount)
	router.GET("/health", server.healthCheck)

	server.router = router
	return server
}

// Public method to access the private `router`
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

// Send error response
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
```

### Create account handler

```go
package api

import (
	"net/http"
	db "simple-bank/db/sqlc"

	"github.com/gin-gonic/gin"
)

type createAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=USD EUR INR"`
}

func (server *Server) healthCheck(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "ok")
}

func (server *Server) createAccount(ctx *gin.Context) {
	// The function needs a Gin context, look for the POST definiton
	// The function expects a function with a Context - "func(*Context)"
	var req createAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateAccountParams{
		Owner:    req.Owner,
		Balance:  0,
		Currency: req.Currency,
	}

	account, err := server.store.CreateAccount(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	ctx.JSON(http.StatusOK, account)
}
```

### Start the server

The `main.go` will initialize the server and open a database connection.

```go
package main

import (
	"database/sql"
	"log"
	"simple-bank/api"
	db "simple-bank/db/sqlc"

	_ "github.com/lib/pq"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgres://root:secret@localhost:5432/simple_bank?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
```

NOTE: The get all sqlc function returns null if not records are found.

So we need to configure the sqlc.yml by adding the following config property

```yaml
emit_empty_slices: true
```