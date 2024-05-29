package adapters

import (
	"context"
	"database/sql"
	"github.com/AlecSmith96/faceit-user-service/internal/entities"
	"github.com/AlecSmith96/faceit-user-service/internal/usecases"
	"github.com/google/uuid"
	"github.com/pressly/goose"
	"log/slog"
	"time"
)

type PostgresAdapter struct {
	db *sql.DB
}

var _ usecases.UserCreator = &PostgresAdapter{}
var _ usecases.UserDeleter = &PostgresAdapter{}
var _ usecases.UserUpdater = &PostgresAdapter{}

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

func (p *PostgresAdapter) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	result, err := p.db.Exec("DELETE FROM platform_user WHERE id = $1;", userID)
	if err != nil {
		slog.Debug("unable to delete user", "err", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		slog.Debug("unable to get rows affected", "err", err)
		return err
	}

	if rowsAffected == 0 {
		slog.Debug("user not found", "userID", userID)
		return entities.ErrUserNotFound
	}

	return nil
}

func (p *PostgresAdapter) UpdateUser(ctx context.Context, userID uuid.UUID, firstName, lastName, nickname, password, email, country string) (*entities.User, error) {
	result, err := p.db.Query(
		"UPDATE platform_user  SET first_name = $2, last_name = $3, nickname = $4, password = $5, email = $6, country = $7, updated_at = $8 WHERE id = $1 RETURNING *",
		userID,
		firstName,
		lastName,
		nickname,
		password,
		email,
		country,
		time.Now(),
	)
	if err != nil {
		slog.Debug("error updating creator", "err", err)
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
	} else {
		slog.Debug("user not found", "userID", userID)
		return nil, entities.ErrUserNotFound
	}

	return &user, nil
}
