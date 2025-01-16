FROM golang:1.21-alpine AS builder

WORKDIR /build

RUN apk --no-cache add make

# Copy the source code
COPY ["go.mod", "go.sum", "./"]
RUN go mod download

COPY . .

# Build the Go app
RUN go build -o ./bin/jwt_auth cmd/auth/main.go

FROM alpine

COPY --from=builder /build/bin/jwt_auth /

# #EXPOSE the port
EXPOSE 3000

# Run the executable
CMD ["/jwt_auth"]