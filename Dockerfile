# Build stage
FROM golang:1.24-alpine3.22 as builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

# Run stage
FROM alpine:3.22
WORKDIR /app
COPY --from=builder /app/main .
COPY app.env .

# Expose the port
EXPOSE 8080

# Run the application
CMD ["/app/main"]
