package adapters

import (
	"context"
	"database/sql"
	"encoding/base64"
	"fmt"
	"github.com/AlecSmith96/faceit-user-service/internal/entities"
	"github.com/AlecSmith96/faceit-user-service/internal/usecases"
	"github.com/google/uuid"
	"github.com/pressly/goose"
	"log/slog"
	"strings"
	"time"
)

type PostgresAdapter struct {
	db *sql.DB
}

var _ usecases.UserCreator = &PostgresAdapter{}
var _ usecases.UserDeleter = &PostgresAdapter{}
var _ usecases.UserUpdater = &PostgresAdapter{}
var _ usecases.UserGetter = &PostgresAdapter{}
var _ usecases.ReadinessChecker = &PostgresAdapter{}

func NewPostgresAdapter(db *sql.DB) *PostgresAdapter {
	return &PostgresAdapter{db: db}
}

// PerformDataMigration is a function that ensure that the database has had all migration ran against it on startup
func (p *PostgresAdapter) PerformDataMigration(gooseDir string) error {
	return goose.Up(p.db, gooseDir)
}

func encodePageToken(userID uuid.UUID, createdAt time.Time) string {
	token := fmt.Sprintf("%s|%s", createdAt.Format(time.RFC3339Nano), userID.String())
	return base64.URLEncoding.EncodeToString([]byte(token))
}

func decodePageToken(token string) (uuid.UUID, time.Time, error) {
	decoded, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		return uuid.Nil, time.Time{}, err
	}
	parts := strings.SplitN(string(decoded), "|", 2)
	if len(parts) != 2 {
		return uuid.Nil, time.Time{}, fmt.Errorf("invalid token format")
	}

	createdAt, err := time.Parse(time.RFC3339Nano, parts[0])
	if err != nil {
		return uuid.Nil, time.Time{}, err
	}

	userID, err := uuid.Parse(parts[1])
	if err != nil {
		return uuid.Nil, time.Time{}, err
	}

	return userID, createdAt, nil
}

func (p *PostgresAdapter) CreateUser(ctx context.Context, firstName, lastName, nickname, password, email, country string) (*entities.User, error) {
	var user entities.User
	err := p.db.QueryRowContext(
		ctx,
		"INSERT INTO platform_user (first_name, last_name, nickname, password, email, country) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *",
		firstName,
		lastName,
		nickname,
		password,
		email,
		country,
	).Scan(
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
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint \"platform_user_email_key\"") {
			slog.Debug("email already registered to a user", "err", err)
			return nil, entities.ErrEmailAlreadyUsed

		}
		slog.Debug("error inserting user record", "err", err)
		return nil, err
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
		"UPDATE platform_user SET first_name = $2, last_name = $3, nickname = $4, password = $5, email = $6, country = $7, updated_at = $8 WHERE id = $1 RETURNING *",
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

func (p *PostgresAdapter) GetPaginatedUsers(
	ctx context.Context,
	firstName,
	lastName,
	nickname,
	email,
	country string,
	pageInfo entities.PageInfo,
) ([]entities.User, string, error) {
	var userID uuid.UUID
	var createdAtPageToken time.Time
	var err error
	if pageInfo.NextPageToken != "" {
		userID, createdAtPageToken, err = decodePageToken(pageInfo.NextPageToken)
		if err != nil {
			slog.Debug("decoding page token", "err", err)
			return nil, "", err
		}
	}

	queryString := `SELECT * FROM platform_user WHERE 1=1 `
	queryParamIndex := 1
	queryParams := make([]any, 0)
	if firstName != "" {
		queryString += fmt.Sprintf("AND first_name ILIKE $%d ", queryParamIndex)
		queryParamIndex++
		queryParams = append(queryParams, "%"+firstName+"%")
	}

	if lastName != "" {
		queryString += fmt.Sprintf("AND last_name ILIKE $%d ", queryParamIndex)
		queryParamIndex++
		queryParams = append(queryParams, "%"+lastName+"%")
	}

	if nickname != "" {
		queryString += fmt.Sprintf("AND nickname ILIKE $%d ", queryParamIndex)
		queryParamIndex++
		queryParams = append(queryParams, "%"+nickname+"%")
	}

	if email != "" {
		queryString += fmt.Sprintf("AND email ILIKE $%d ", queryParamIndex)
		queryParamIndex++
		queryParams = append(queryParams, "%"+email+"%")
	}

	if country != "" {
		queryString += fmt.Sprintf("AND country ILIKE $%d ", queryParamIndex)
		queryParamIndex++
		queryParams = append(queryParams, "%"+country+"%")
	}

	if pageInfo.NextPageToken != "" {
		queryString += fmt.Sprintf(`AND (created_at, id) > ($%d, $%d)`, queryParamIndex, queryParamIndex+1)
		queryParams = append(queryParams, createdAtPageToken, userID)
	}

	queryString += fmt.Sprintf(` ORDER BY created_at, id LIMIT %d;`, pageInfo.PageSize)
	rows, err := p.db.Query(queryString, queryParams...)
	if err != nil {
		slog.Debug("error updating creator", "err", err)
		return nil, "", err
	}

	users := make([]entities.User, 0)
	for rows.Next() {
		var user entities.User
		err = rows.Scan(
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
			return nil, "", err
		}

		users = append(users, user)
	}

	if len(users) == 0 || len(users) < pageInfo.PageSize {
		return users, "", nil
	}

	lastUser := users[len(users)-1]
	nextPageToken := encodePageToken(lastUser.ID, lastUser.CreatedAt)

	return users, nextPageToken, nil
}

func (p *PostgresAdapter) CheckConnection() error {
	err := p.db.Ping()
	if err != nil {
		slog.Debug("database connection lost")
		return err
	}

	return nil
}
