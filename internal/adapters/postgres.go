package adapters

import (
	"context"
	"database/sql"
	"github.com/AlecSmith96/faceit-user-service/internal/entities"
	"github.com/AlecSmith96/faceit-user-service/internal/usecases"
	"github.com/pressly/goose"
	"log/slog"
)

type PostgresAdapter struct {
	db *sql.DB
}

var _ usecases.UserCreator = &PostgresAdapter{}

func NewPostgresAdapter(db *sql.DB) *PostgresAdapter {
	return &PostgresAdapter{db: db}
}

// PerformDataMigration is a function that ensure that the database has had all migration ran against it on startup
func (p *PostgresAdapter) PerformDataMigration(gooseDir string) error {
	return goose.Up(p.db, gooseDir)
}

func (p *PostgresAdapter) CreateUser(ctx context.Context, firstName, lastName, nickname, password, email, country string) (*entities.User, error) {
	result, err := p.db.Query(
		"INSERT INTO platform_user (first_name, last_name, nickname, password, email, country) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *",
		firstName,
		lastName,
		nickname,
		password,
		email,
		country,
	)
	if err != nil {
		slog.Debug("error inserting user record", "err", err)
		return nil, err
	}

	var user entities.User
	if result.Next() {
		err = result.Scan(
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.Nickname,
			&user.Password,
			&user.Email,
			&user.Country,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			slog.Debug("marshalling user to struct", "err", err)
			return nil, err
		}
	}

	return &user, nil
}
