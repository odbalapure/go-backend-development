# Build stage
FROM golang:1.24-alpine3.22 as builder
WORKDIR /app
COPY . .
RUN go build -o main main.go
# Install curl to download the migrate binary
# Its not present in the base image
# RUN apk add curl
# RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.18.3/migrate.linux-amd64.tar.gz | tar xvz

# Run stage
FROM alpine:3.22
WORKDIR /app
COPY --from=builder /app/main .
# COPY --from=builder /app/migrate ./migrate
COPY db/migration ./db/migration
COPY wait-for.sh .
COPY start.sh .
RUN chmod +x wait-for.sh start.sh
COPY app.env .

# Expose the port
EXPOSE 8080

# Run the application
CMD ["/app/main"]
ENTRYPOINT ["/app/start.sh"]
