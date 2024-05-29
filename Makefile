# command for running postgres by itself, useful when developing
start-postgres:
	-make stop-postgres
	podman run -d -p 5432:5432 --name postgres -e POSTGRES_PASSWORD="postgres" -e POSTGRES_DB="users" postgres:13.6-alpine

stop-postgres:
	bash -c "podman stop postgres || true"
	bash -c "podman rm postgres || true"

# generate swagger documentation
swagger:
	swag init --generalInfo ../../cmd/main.go --dir ./internal/usecases

run-service:
	-make stop-postgres
	-make stop-service
	podman-compose up --build -d

stop-service:
	podman-compose down
	bash -c "podman stop faceit-user-service_faceit-user-service_1 || true"
	bash -c "podman rm faceit-user-service_faceit-user-service_1 || true"
