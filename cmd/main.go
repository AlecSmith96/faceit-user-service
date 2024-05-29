package main

import (
	"database/sql"
	"github.com/AlecSmith96/faceit-user-service/internal/adapters"
	"github.com/AlecSmith96/faceit-user-service/internal/drivers"
	_ "github.com/lib/pq"
	"log/slog"
	"os"
)

const (
	gooseDir = "./db/goose"
)

func main() {
	conf, err := adapters.NewConfig()
	if err != nil {
		slog.Error("reading config", "err", err)
		os.Exit(1)
	}

	db, err := sql.Open("postgres", conf.PostgresConnectionURI)
	if err != nil {
		slog.Error("connecting to postgres", "err", err)
		os.Exit(1)
	}

	postgresAdapter := adapters.NewPostgresAdapter(db)

	err = postgresAdapter.PerformDataMigration(gooseDir)
	if err != nil {
		slog.Error("running migrations", "err", err)
		os.Exit(1)
	}

	router := drivers.NewRouter(postgresAdapter)

	err = router.Run()
	if err != nil {
		slog.Error("running gin router", "err", err)
	}

}
