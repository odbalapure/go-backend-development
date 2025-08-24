## Swagger Doc Generation

### Swagger generation

We are trying to generate swagger doc for our gRCP gateway with minimal dependencies. Clone this repository:

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

**NOTE**: You don't the same folder/path structure. This is just for clarity and keeping things similar to the original source code.

## Auto generate docs

We have to run the `make proto` command every time to generate a new swagger doc. We can install a swagger react plugin to do the same.

But to keep things minimal just copy the contents of the `dist` folder from the [dist](https://github.com/swagger-api/swagger-ui/tree/master/dist) to our `doc/swagger` folder.
swagger-ui-react

```sh
cp -r ../swagger-ui/dist/* doc/swagger/
```

Now update the url of the destination file here [init.js](../doc/swagger/swagger-initializer.js) to `simple_bank.swagger.json`.

Now, we need to serve the contents of the swagger folder. For that we need to create an endpoint that will server these static files. 

Update the main.go `runGatewayServer` function.

```go
fs := http.FileServer(http.Dir("./doc/swagger"))
mux.Handle("/swagger/", http.StripPrefix("/swagger/", fs))
```
