## Docker

Create a Dockerfile in the root directory.

```Dockerfile
FROM golang:1.24-alpine3.22
WORKDIR /app
COPY . .
RUN go build -o main main.go

EXPOSE 8080
CMD ["/app/main"]
```

> `main` is the name of the executable present under the `/app`.

### Creating a docker image

Create a docker image using `docker build -t simplebank:latest`.

Now list the images using `docker images`.

```sh
ombalapure@Oms-MacBook-Air simple-bank % docker images
REPOSITORY   TAG         IMAGE ID       CREATED         SIZE
simplebank   latest      565cff4999a8   2 minutes ago   913MB```

NOTE: You might notice that the image size is 913MB. Heavier than the size of the alpine image.
```

To reduce the size of the image, we can use a multi-stage build.

```Dockerfile
FROM golang:1.24-alpine3.22 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go
```

### Multi stage build

The size of the resulting image will be much smaller this time

```Dockerfile
# Build stage
FROM golang:1.24-alpine3.22 as builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

# Run stage
FROM alpine:3.22
WORKDIR /app
COPY --from=builder /app/main .

EXPOSE 8080
CMD ["/app/main"]
```

- Single stage ships the entire build toolchain (go compiler, caches, headers, package manager, source code) inside the image.
- In multistage we build the source code on the `golang` base image. And copy the build from the `golang` image to the `alpine` base image.

## Running the docker image

```sh
docker run --name simplebank -p 8080:8080 simplebank:latest
```

This will fail with error "cannot load config file"; we need to copy the `app.env` file as well as part of run stage.

```Dockerfile
# Run stage
FROM alpine:3.22
WORKDIR /app
COPY --from=builder /app/main .
COPY app.env .
```

The `Gin` server can be run in release mode using

```sh
docker run --name simplebank -e GIN_MODE=release -p 8080:8080 simplebank:latest
```

> Calls to the postgres DB will result return `connection refused` error.

## Connecting two docker containers

// postgres container
```json
"Networks": {
    "bridge": {
        "IPAMConfig": null,
        "Links": null,
        "Aliases": null,
        "MacAddress": "96:0a:15:b8:ed:a1",
        "DriverOpts": null,
        "GwPriority": 0,
        "NetworkID": "e53314eb569e9cf5d6e5e677b4ec5c9f0d083b4225fd3408030f8be47d47e9e0",
        "EndpointID": "533c32f999244352c57cfc64b641d0f466d846a9c61df6dfc944793e010755ec",
        "Gateway": "172.17.0.1",
        "IPAddress": "172.17.0.2",
        "IPPrefixLen": 16,
        "IPv6Gateway": "",
        "GlobalIPv6Address": "",
        "GlobalIPv6PrefixLen": 0,
        "DNSNames": null
    }
}
```

// simplebank container
```json
"Networks": {
    "bridge": {
        "IPAMConfig": null,
        "Links": null,
        "Aliases": null,
        "MacAddress": "ca:c8:95:d9:02:ba",
        "DriverOpts": null,
        "GwPriority": 0,
        "NetworkID": "e53314eb569e9cf5d6e5e677b4ec5c9f0d083b4225fd3408030f8be47d47e9e0",
        "EndpointID": "f4252fe08806cf3fd154bd162039777137ce919e0a333f01004721ebe387e083",
        "Gateway": "172.17.0.1",
        "IPAddress": "172.17.0.3",
        "IPPrefixLen": 16,
        "IPv6Gateway": "",
        "GlobalIPv6Address": "",
        "GlobalIPv6PrefixLen": 0,
        "DNSNames": null
    }
}
```

Note that the IP Addresses are assigned by the docker network. They are different for each container.

We cant set the IP address explcitly overriding Viper config but theres a better way
```sh
-DB_SOURCE=postgres://root:secret@172.17.0.2:5432/simple_bank?sslmode=disabled
```

### Create a new network

Docker has a `bridge` network, which is the default network for docker. It is a virtual network that is created by docker.

docker network inspect bridge
```json
"Containers": {
    "0e74ea3a98885bd1680fa49588861ca3d1a1d9073f86e5e33c26759ed66ffcbc": {
        "Name": "simplebank",
        "EndpointID": "f4252fe08806cf3fd154bd162039777137ce919e0a333f01004721ebe387e083",
        "MacAddress": "ca:c8:95:d9:02:ba",
        "IPv4Address": "172.17.0.3/16",
        "IPv6Address": ""
    },
    "bb8104e685db737383284462da977d55f0907fd57eb8771e7bbc2b34e4d23390": {
        "Name": "postgres12",
        "EndpointID": "533c32f999244352c57cfc64b641d0f466d846a9c61df6dfc944793e010755ec",
        "MacAddress": "96:0a:15:b8:ed:a1",
        "IPv4Address": "172.17.0.2/16",
        "IPv6Address": ""
    }
}
```

So we create our own virtual network
```sh
docker network create bank-network
```

Connect to the network
```sh
docker network connect bank-network postgres12
```

Now connect the postgres12 to this newly created network

```sh
ombalapure@Oms-MacBook-Air simple-bank % docker network connect bank-network postgres12
ombalapure@Oms-MacBook-Air simple-bank % docker network inspect bank-network
```

```json
[
    {
        "Name": "bank-network",
        "Id": "8f71afa808467ad86797bcbc1ea3d0369f806ce9066af9a90d5c493f15b63159",
        "Created": "2025-08-12T04:09:35.188983842Z",
        "Scope": "local",
        "Driver": "bridge",
        "EnableIPv4": true,
        "EnableIPv6": false,
        "IPAM": {
            "Driver": "default",
            "Options": {},
            "Config": [
                {
                    "Subnet": "172.18.0.0/16",
                    "Gateway": "172.18.0.1"
                }
            ]
        },
        "Internal": false,
        "Attachable": false,
        "Ingress": false,
        "ConfigFrom": {
            "Network": ""
        },
        "ConfigOnly": false,
        "Containers": {
            "bb8104e685db737383284462da977d55f0907fd57eb8771e7bbc2b34e4d23390": {
                "Name": "postgres12",
                "EndpointID": "c7f4e1e0584a897d3817aaadccc3d8943ccae123175a161d140eb1c9f228905e",
                "MacAddress": "0e:cd:17:4c:8f:93",
                "IPv4Address": "172.18.0.2/16",
                "IPv6Address": ""
            }
        },
        "Options": {
            "com.docker.network.enable_ipv4": "true",
            "com.docker.network.enable_ipv6": "false"
        },
        "Labels": {}
    }
]
```

We can see that the `postgres12` container is connected to the `bank-network` network.

Now after **inspect** we can see that `postgres12` is now connected to two different networks including our own; which is completely ok.

Now while running the docker container mention the name of the `--network bank-network`

```sh
ombalapure@Oms-MacBook-Air simple-bank % docker run --name simple-bank --network bank-network -p 8080:8080 -e GIN_MODE=release -e DB_SOURCE="postgresql://root:secret@postgres12:5432/simple_bank?sslmode=false" simplebank:latest
```

Now the API calls will work

```curl
curl --location 'http://localhost:8080/users/login' \
--header 'Content-Type: application/json' \
--data '{
    "username": "ombalapure",
    "password": "secret"
}'
```

## Docker Compose

Docker compose lets us launch all our services at once on a single network.

```Dockerfile
version: "3.9"

services:
  postgres:
    image: postgres:12-alpine
    environment:
      - POSTGRES_USER=root
      - POSTGRES_PASSWORD=root
      - POSTGRES_DB=simple_bank
    ports:
      - 5432:5432
  
  api:
    build:
      context: .
      dockerfile: Dockerfile
    ports:
      - 8080:8080
    environment:
      - DB_SOURCE=postgres://root:secret@postgres:5432/simple_bank?sslmode=disable
```

- `build` means that this is a custom image unlike postgres:12-alpine; and the `.` means current directory.
- All services in the Docker compose will run on the same network they can commnuicate with each other via name. Hence DB_SOURCE=postgres://root:secret@`postgres`:5432/simple_bank?sslmode=disable.

**NOTE**:
- The image name will `simple-bank-api`, the name of the folder followed by the service name i.e. `api`.
- The container names will be `postgres-1` and `api-1` under `simple-bank`.

> Looking at the logs carefully you can find a new network being created by the name `simple-bank_default`


### Testing the services

For the database to work, download the `golang-migrate`.

```
RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.18.3/migrate.linux-amd64.tar.gz | tar xvz

depends_on:
    - postgres
```

> `depends_on` does not guarantee the database is up, we can use the `wait-for.sh` script.
