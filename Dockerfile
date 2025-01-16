FROM golang:1.21-alpine AS builder

RUN go version
ENV GOPATH=/

COPY ./ ./

# install psql
RUN apk update
RUN apk add postgresql-client

# make wait-for-postgres.sh executable
RUN chmod +x wait-for-postgres.sh

# build go app
# RUN go mod download
RUN go build -o jwt_auth cmd/auth/main.go

# #EXPOSE the port
EXPOSE 3000

# Run the executable
CMD ["./jwt_auth"]