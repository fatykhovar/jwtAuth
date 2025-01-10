# use official Golang image
FROM golang:1.21-alpine AS builder

# set working directory
WORKDIR /build

# # Download and install the dependencies
# RUN go get -d -v ./...

RUN apk --no-cache add make

# Copy the source code
COPY ["go.mod", "go.sum", "./"]
RUN go mod download

COPY . .

# # Build the Go app
RUN go build -o ./bin/jwt_auth cmd/auth/main.go

FROM alpine
# WORKDIR c:/src
COPY --from=builder /build/bin/jwt_auth /

# #EXPOSE the port
EXPOSE 3000

# # Run the executable
CMD ["/jwt_auth"]