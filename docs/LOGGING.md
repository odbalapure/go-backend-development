## Loggers

We need logging in our apps because we need to know:

- gRPC methods used
- Request duration
- Response status code

Also, the logs need to be strucutred so that can be parsed and indexed easily.

## Adding logs in gRPC

We have `grpc` interceptor that can log requests and responses. It can be called in the **main.go**.

```go
// gRPC logger
grpcLogger := grpc.UnaryInterceptor(gapi.GrpcLogger)

// gRPC server
grpcServer := grpc.NewServer(grpcLogger)
```

> Since our requests are unary we used the `grpc.UnaryInterceptor`.

## Using a 3rd party logger

Using the `Printf` statements is not scalable not can slow down the app.

We can use a tool like `zerolog`, which writes logs in a binary efficient JSON format which is cheaper than printf style formatting.

It introduces consistency, eg: dev A can log `userId=123` and dev B can log `uuid: 123`.

```sh
github.com/rs/zerolog
```

> Most JSON loggers uses Go's `encoding/json` which relies on reflection and its slow. `zerolog` uses preallocated buffers and type specific methods.
