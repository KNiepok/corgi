# Builder container
FROM golang:1.14-alpine AS build_base

RUN apk add --no-cache git

WORKDIR /tmp/corgi

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN apk add --no-cache gcc musl-dev


RUN go build -o ./bin/corgi cmd/corgi/main.go

# Runner container
FROM alpine:3.9
RUN apk add ca-certificates

COPY --from=build_base /tmp/corgi/bin/corgi /app/corgi

# todo make that variable
EXPOSE 8080

# Run the binary
CMD ["/app/corgi"]