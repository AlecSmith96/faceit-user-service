package main

import (
	"database/sql"
	_ "github.com/AlecSmith96/faceit-user-service/docs"
	"github.com/AlecSmith96/faceit-user-service/internal/adapters"
	"github.com/AlecSmith96/faceit-user-service/internal/drivers"
	_ "github.com/lib/pq"
	"log/slog"
	"os"
)

const (
	gooseDir = "./db/goose"
)

// @title faceit-user-service
// @version 1.0
// @description This is a simple REST server providing CRUD operations on a User object

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

	kafkaAdapter, err := adapters.NewKafkaAdapter(conf.KafkaHost)
	if err != nil {
		slog.Error("creating kafka adapter", "err", err)
		os.Exit(1)
	}
	defer kafkaAdapter.CloseConn()

	router := drivers.NewRouter(kafkaAdapter, postgresAdapter, postgresAdapter, postgresAdapter, postgresAdapter, postgresAdapter)

	err = router.Run()
	if err != nil {
		slog.Error("running gin router", "err", err)
	}

}
