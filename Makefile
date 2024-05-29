start_postgres:
	-make stop_postgres
	podman run -d -p 5432:5432 --name postgres -e POSTGRES_PASSWORD="postgres" -e POSTGRES_DB="users" postgres:13.6-alpine

stop_postgres:
	bash -c "podman stop postgres || true"
	bash -c "podman rm postgres || true"

# generate swagger documentation
swagger:
	swag init --generalInfo ../../cmd/main.go --dir ./internal/usecases

