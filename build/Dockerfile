FROM golang:1.21-alpine as build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go mod vendor

# Build binary
RUN  GO111MODULE="on" CGO_ENABLED=0 GOOS=linux go build -mod=vendor -o main cmd/main.go

#FROM scratch

# TODO: debugging, revert before finishing
FROM alpine:latest

WORKDIR /app

COPY --from=build /app/main .

# migration files are required to connect to db
COPY --from=build /app/db/ ./db/

# Expose port 8080 and run binary
EXPOSE 8080
CMD ["./main"]