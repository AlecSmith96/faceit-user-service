FROM golang:1.21-alpine as build

RUN apk add --no-progress --no-cache gcc musl-dev
WORKDIR /app

COPY . .

COPY go.mod go.sum ./
RUN go mod download

RUN go mod vendor

# Build binary
RUN  GO111MODULE="on" GOOS=linux GOARCH=amd64 go build -mod=vendor -tags musl -ldflags '-extldflags "-static"' -o main cmd/main.go

FROM alpine:latest

WORKDIR /app

# Install bash and curl
RUN apk add --no-cache bash curl

COPY --from=build /app/main .

# migration files are required to connect to db
COPY --from=build /app/db/ ./db/

# Expose port 8080 and run binary
EXPOSE 8080
CMD ["./main"]