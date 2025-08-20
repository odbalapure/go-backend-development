## gRPC

It is an RPC framework developed by Google. And is now part of the CNCF:
- The remote interaction code is handled by gRPC.
- Support multiple programming languages.

### Why gRPC
- High performance: HTTP/2 with binary framing, multiplexing, header compression, bidirectional communication.
- Strong API contract: Server and client share the same protobuf RPC definition with strongly typed data.
- Automatic code generation: Code that serialize/deserialze data, transfer data b/w client and server are auto generated.

### Types of gRPC
- Unary: Typical client/server model.
- Client streaming: send stream of requests and server responds with one response.
- Server streaming: opposite of client streaming.
- Bidirectional streaming: sending multiple requests b/w client and server in abritrary order.

### gRPC gateway

It serves both HTTP and gRPC requests at once.
- Plugin of protobf compiler
- Generate proxy code from protobuf
- Translate HTTP JSON calls to gRPC
    - In process translation: only for unary (HTTP to gRPC conversion happens within the same application process)
    - Separate proxy server: both unary and streaming (Separate proxy server means the HTTP to gRPC conversion happens in a completely different service)
- Write code once; server both HTTP and gRPC requests


## Generate code from protobuf

Install it using

```
brew install protobuf
```

Check if its installed via

```
protoc --version
```

We need two more plugins
- **protoc-gen-go**: Here are the Go structs for your data
- **protoc-gen-go-grpc**: Here are the Go interfaces for your gRPC service

The quick start page is [here](https://grpc.io/docs/languages/go/quickstart/)

## Define gRPC API

Create a `user.proto` and `rpc_create_user.proto` file and a service `service_simple_bank.proto` that will use these types.

```proto3
// user.proto
syntax = "proto3";

package pb;

import "google/protobuf/timestamp.proto";

option go_package = "simple-bank/pb";

message User {
    string username = 1;
    string full_name = 2;
    string email = 3;
    google.protobuf.Timestamp password_changed_at = 4;
    google.protobuf.Timestamp created_at = 5;
}

// rpc_create_user.proto
syntax = "proto3";

package pb;

import "user.proto";

option go_package = "simple-bank/pb";

message CreateUserRequest {
    string username = 1;
    string password = 2;
    string email = 3;
    string full_name = 4;
}

message CreateUserResponse {
    User user = 1;
}

// service_simple_bank.proto
syntax = "proto3";

package pb;

import "rpc_create_user.proto";
import "rpc_login_user.proto";

option go_package = "simple-bank/pb";

service SimpleBank {
    rpc CreateUser(CreateUserRequest) returns (CreateUserResponse) {}
    rpc LoginUser(LoginUserRequest) returns (LoginUserResponse) {}
}
```

For the imports to work we need to configure a path in the `settings.json` of the IDE.

```json
"protoc": {
    "options": [
        "--protoc_path=proto",
    ]
}
```

Now generate the code from files present under **proto** folder:

```bash
protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
    proto/*.proto
```

> Run `go mod tidy` to add any missing dependencies.

## Running gRPC server

The [NewServer](../gapi/server.go) method will be similar to [Gin's](../api/server.go).

The only difference is the struct `pb.UnimplementedSimpleBankServer`. This is a placeholder for unimplemented service RPCs so that the GRPC server does not `panic`.

```go
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
```

If you don't add it, the gRPC server won't start:
```sh
./main.go:41:42: cannot use server (variable of type *gapi.Server) as pb.SimpleBankServer value in argument to pb.RegisterSimpleBankServer: *gapi.Server does not implement pb.SimpleBankServer (missing method mustEmbedUnimplementedSimpleBankServer)
```

Create a function that will start the gRPC server

```go
func runGrpcServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterSimpleBankServer(grpcServer, server)
	// This allows a gRPC client to explore what RPC are available in the server
	// Sort of a self documentation for a server
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal("cannot create listener:", err)
	}

	log.Printf("gRPC server started at %s", listener.Addr().String())

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start gRPC server:", err)
	}
}
```

The RPCs can be listed using:

```sh
ombalapure@Oms-MacBook-Air simple-bank % grpcurl -plaintext localhost:9090 list pb.SimpleBank
pb.SimpleBank.CreateUser
pb.SimpleBank.LoginUser
```

## Create a gRPC endpoint

Refer an un-implemented method from [service_simple_bank_grpc.pb.go](../pb/service_simple_bank_grpc.pb.go) for eg: `CreateUser`.

```go
func (UnimplementedSimpleBankServer) CreateUser(context.Context, *CreateUserRequest) (*CreateUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateUser not implemented")
}
```

This will accept a `Server` struct in its receiver. The logic will same as that of HTTP login user. Only the response will change:

```go
func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	hashedPassword, err := util.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password: %v", err)
	}
	// ...
	rsp, nil
}
```

Now, test this gRPC endpoint:

```sh
ombalapure@Oms-MacBook-Air simple-bank % grpcurl -plaintext -d '{
  "username": "john_doe",
  "password": "secret123",
  "email": "john@example.com",
  "full_name": "John Doe"
}' localhost:9090 pb.SimpleBank/CreateUser
{
  "user": {
    "username": "john_doe",
    "fullName": "John Doe",
    "email": "john@example.com",
    "passwordChangedAt": "0001-01-01T00:00:00Z",
    "createdAt": "2025-08-19T00:12:56.917757Z"
}
```

## protobuf V protoc

**protobuf**
- This is the data format/serialization protocol itself
- It's Google's binary data format for serializing structured data
- Think of it as the "language" or "format" for data exchange

**protoc**
- This is the compiler tool that processes `.proto` files
- `protoc` = Protocol buffer Compiler
- It's the executable that reads your `.proto` files and generates code

Inshort:
- Protocol Buffers = The specification/format (like JSON or XML)
- protoc = The compiler (like gcc for C or javac for Java)
- protoc-gen-* = Plugins that tell the compiler what language to generate

> `protoc` is the tool, and Protocol Buffers is what it processes.

## gRPC Gateway

### Required plugins

To create a gateway server that translates your HTTP calls to gRPC calls. We need these:

- "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway"
- "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2"
- "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
- "google.golang.org/protobuf/cmd/protoc-gen-go"

Then clone the [googleapis](https://github.com/googleapis/googleapis) repo. Copy these four files from `googleapis/google/api/*`
- annotations.proto
- field_behavior.proto
- http.proto
- httpbody.proto

### Updating the service proto file

Now update the `proto/service_simple_bank.proto` file.

```proto
import "google/api/annotations.proto";

option go_package = "simple-bank/pb";

service SimpleBank {
    rpc CreateUser(CreateUserRequest) returns (CreateUserResponse) {
        option (google.api.http) = {
            post: "/v1/create_user"
            body: "*"
        };
    }
    rpc LoginUser(LoginUserRequest) returns (LoginUserResponse) {
        option (google.api.http) = {
            post: "/v1/login_user"
            body: "*"
        };
    }
}
```

This import will let use define options for our gateway server.

### Config

Now update the proto command

```sh
proto:
	rm -f pb/*.go
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	# To generate `service_simple_bank.pb.gw.go` gateway file
    --grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
    proto/*.proto
```

### Updating the main.go point 

The make command will generate a gateway file `service_simple_bank.pb.gw.go` under `pb` folder. We need `RegisterSimpleBankHandlerServer` function to setup the gateway server.

```go
import (
	"net/http"
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/protobuf/encoding/protojson"
)

func runGatewayServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	grpcMux := runtime.NewServeMux(
		// Convert camel case response to snake case
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames: true,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		}),
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = pb.RegisterSimpleBankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal("cannot register http -> grpc gateway handler:", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	listener, err := net.Listen("tcp", config.HTTPServerAddress)
	if err != nil {
		log.Fatal("cannot create http -> grpc gateway listener:", err)
	}

	err = http.Serve(listener, mux)
	if err != nil {
		log.Fatal("cannot start http -> grpc gateway server:", err)
	}
}
```

NOTE: We cannot run both the servers simultaneously so one has to run as a go-routine:
```go
go runGatewayServer(config, store)
runGrpcServer(config, store)
```

## Extracting meta data

The metadata can be extracted from the headers and in some cases the context object of the request.

```go
func (s *Server) extractMetadata(ctx context.Context) *Metadata {
	mtd := &Metadata{}

	// NOTE: `md` is a map[string][]string
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		// This check is for the gateway server
		if userAgents := md.Get(grpcGatewayUserAgentHeader); len(userAgents) > 0 {
			mtd.UserAgent = userAgents[0]
		}

		// This check is for the plain GRPC server
		if userAgents := md.Get(userAgentHeader); len(userAgents) > 0 {
			mtd.UserAgent = userAgents[0]
		}

		// This check is for the gateway server
		if clients := md.Get(xForwardedForHeader); len(clients) > 0 {
			mtd.ClientIp = clients[0]
		}
	}

	// This check is for the plain GRPC server
	// NOTE: The IPAddress is stored in the context object of the request  not in the metadata
	if p, ok := peer.FromContext(ctx); ok {
		mtd.ClientIp = p.Addr.String()
	}

	return mtd
}
```