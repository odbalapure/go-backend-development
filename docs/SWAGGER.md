## Swagger Doc Generation

### Swagger generation

We are trying to generate swagger doc for our gRCP gateway with minimal dependencies.

Instead of installing the dependencies we can utilize the source from the `grpc-ecosystem` repositories.

```sh
https://github.com/grpc-ecosystem/grpc-gateway
````

Copy this snippet from the [a_bit_of_everything.proto](https://github.com/grpc-ecosystem/grpc-gateway/blob/main/examples/internal/proto/examplepb/a_bit_of_everything.proto
) file and paste it in the `service_simple_bank.proto`. Keep the same folder structure i.e. `proto/protoc-gen-openapiv2/options/*`.

```proto
import "protoc-gen-openapiv2/options/annotations.proto";

option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_swagger) = {
  info: {
    title: "Simple Bank API"
    version: "1.0"
    contact: {
      name: "Personal Project"
      url: "https://github.com/go-backend-development"
      email: "none@example.com"
    }
    license: {
      name: "BSD 3-Clause License"
      url: "https://github.com/grpc-ecosystem/grpc-gateway/blob/main/LICENSE"
    }
    extensions: {
      key: "x-something-something"
      value: {string_value: "yadda"}
    }
  }
};
```

Copy the **proto** files under [options](https://github.com/grpc-ecosystem/grpc-gateway/tree/main/protoc-gen-openapiv2/options) folder and paste them under our projects' `proto/protoc-gen-openapiv2/options`.

You don't need the same folder/path structure. This is just for clarity and keeping things similar to the original source code.

## Auto generate docs

Copy the contents of the `dist` folder from the [dist](https://github.com/swagger-api/swagger-ui/tree/master/dist) to our `doc/swagger` folder.

> We can install `swagger-ui-react` to generate the swagger documentation - but to keep things minimal we just copy the essentials file required to generate it.

```sh
cp -r ../swagger-ui/dist/* doc/swagger/
```

Now update the url of the destination file here [init.js](../doc/swagger/swagger-initializer.js) to `simple_bank.swagger.json`.

We need to serve the contents of the swagger folder. For that we need to create an endpoint that will server these static files. 

Update the main.go `runGatewayServer` function.

```go
fs := http.FileServer(http.Dir("./doc/swagger"))
mux.Handle("/swagger/", http.StripPrefix("/swagger/", fs))
```

### Embed front end code in the Go binary

If we copy all the FE files into our docker build, we will have lot of FE files along side our Go binary.

Its better we embeds a directory of static files into our _Go binary_ to be later served from an `http.FileSystem`.

This avoids updating our Dockerfile. And all these files are in the memory of the server not the disk; serving them becomes much faster.

> This apporach makes sense if we the static files are not too large in size and will be browsed by a few people.

We can install a tool called as [statik](https://github.com/rakyll/statik) to do so:

```sh
go get github.com/rakyll/statik@latest
```

Now update the make proto command; the command generates a `statik.go` file. This file contains binary data, but it's encoded in a specific way that makes it usable in Go:

```sh
statik -src=./doc/swagger -dest=./doc
```

### Serving the files from the Go binary

Update the main.go file to serve the files from `doc/statik/statik.go`.

```go
// This serves static files from the memory of the server not the disk.
// This is much faster than serving from the disk.
statikFS, err := fs.New()
if err != nil {
  log.Fatal("cannot create statik file system:", err)
}

swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statikFS))
mux.Handle("/swagger/", swaggerHandler)
```

### Testing the change

Now update the `service_simple_bank.proto` file, this change will add the metadata to the swagger documentation:

```proto
service SimpleBank {
    rpc CreateUser(CreateUserRequest) returns (CreateUserResponse) {
        option (google.api.http) = {
            post: "/v1/create_user"
            body: "*"
        };
        // Meta data for create user
        option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_operation) = {
          description: "This API creates a new user using gRPC"
          summary: "Create new User"
          tags: "create_user"
        };
    }
    rpc LoginUser(LoginUserRequest) returns (LoginUserResponse) {
        option (google.api.http) = {
            post: "/v1/login_user"
            body: "*"
        };
        // Meta data for login user
        option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_operation) = {
          description: "This API logs in a user using gRPC"
          summary: "Login User"
          tags: "login_user"
        };
    }
}
```

> We need to run `make proto` and `make server` everytime to see the changes in our swagger documentation.
