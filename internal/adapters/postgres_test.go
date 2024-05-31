package adapters_test

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/AlecSmith96/faceit-user-service/internal/adapters"
	"github.com/AlecSmith96/faceit-user-service/internal/entities"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	. "github.com/onsi/gomega"
	"testing"
	"time"
)

func TestNewPostgresAdapter(t *testing.T) {
	g := NewWithT(t)
	db, _, err := sqlmock.New()
	g.Expect(err).ToNot(HaveOccurred())

	adapter := adapters.NewPostgresAdapter(db)
	g.Expect(adapter).To(BeAssignableToTypeOf(&adapters.PostgresAdapter{}))

	defer db.Close()
}

func TestPostgresAdapter_CreateUser(t *testing.T) {
	g := NewWithT(t)
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	adapter := adapters.NewPostgresAdapter(db)

	userEntity := entities.User{
		ID:        uuid.New(),
		FirstName: "alec",
		LastName:  "smith",
		Nickname:  "alecsmith",
		Password:  "somepasword",
		Email:     "alec@email.com",
		Country:   "UK",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	mock.ExpectQuery(`INSERT INTO platform_user \(first_name, last_name, nickname, password, email, country\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6\) RETURNING *`).
		WithArgs("alec", "smith", "alecsmith", "somepassword", "alec@email.com", "UK").
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "first_name", "last_name", "nickname", "password", "email", "country", "created_at", "updated_at"}).
				AddRow(userEntity.ID, userEntity.FirstName, userEntity.LastName, userEntity.Nickname, userEntity.Password, userEntity.Email, userEntity.Country, userEntity.CreatedAt, userEntity.UpdatedAt))

	user, err := adapter.CreateUser(
		context.Background(),
		"alec",
		"smith",
		"alecsmith",
		"somepassword",
		"alec@email.com",
		"UK",
	)
	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(*user).To(Equal(userEntity))
}

func TestPostgresAdapter_CreateUser_EmailUniqueConstraint(t *testing.T) {
	g := NewWithT(t)
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	adapter := adapters.NewPostgresAdapter(db)

	mock.ExpectQuery(`INSERT INTO platform_user \(first_name, last_name, nickname, password, email, country\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6\) RETURNING *`).
		WithArgs("alec", "smith", "alecsmith", "somepassword", "alec@email.com", "UK").
		WillReturnError(errors.New("pq: duplicate key value violates unique constraint \"platform_user_email_key\""))

	user, err := adapter.CreateUser(
		context.Background(),
		"alec",
		"smith",
		"alecsmith",
		"somepassword",
		"alec@email.com",
		"UK",
	)
	g.Expect(err).To(MatchError(entities.ErrEmailAlreadyUsed))
	g.Expect(user).To(BeNil())
}

func TestPostgresAdapter_CreateUser_ExecError(t *testing.T) {
	g := NewWithT(t)
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	adapter := adapters.NewPostgresAdapter(db)

	mock.ExpectQuery(`INSERT INTO platform_user \(first_name, last_name, nickname, password, email, country\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6\) RETURNING *`).
		WithArgs("alec", "smith", "alecsmith", "somepassword", "alec@email.com", "UK").
		WillReturnError(errors.New("an error occurred"))

	user, err := adapter.CreateUser(
		context.Background(),
		"alec",
		"smith",
		"alecsmith",
		"somepassword",
		"alec@email.com",
		"UK",
	)
	g.Expect(err).To(MatchError("an error occurred"))
	g.Expect(user).To(BeNil())
}

func TestNewPostgresAdapter_GetPaginatedUsers_firstName(t *testing.T) {
	g := NewWithT(t)
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	adapter := adapters.NewPostgresAdapter(db)

	userEntities := []entities.User{
		{
			ID:        uuid.New(),
			FirstName: "alec",
			LastName:  "smith",
			Nickname:  "alecsmith",
			Password:  "somepasword",
			Email:     "alec@email.com",
			Country:   "UK",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        uuid.New(),
			FirstName: "alecsmith",
			LastName:  "smith",
			Nickname:  "alecsmith",
			Password:  "somepasword",
			Email:     "alec2@email.com",
			Country:   "UK",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	mock.ExpectQuery(`SELECT \* FROM platform_user WHERE 1=1 AND first_name ILIKE \$1 ORDER BY created_at, id LIMIT 2;`).
		WithArgs("%alec%").
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "first_name", "last_name", "nickname", "password", "email", "country", "created_at", "updated_at"}).
				AddRow(userEntities[0].ID, userEntities[0].FirstName, userEntities[0].LastName, userEntities[0].Nickname, userEntities[0].Password, userEntities[0].Email, userEntities[0].Country, userEntities[0].CreatedAt, userEntities[0].UpdatedAt).
				AddRow(userEntities[1].ID, userEntities[1].FirstName, userEntities[1].LastName, userEntities[1].Nickname, userEntities[1].Password, userEntities[1].Email, userEntities[1].Country, userEntities[1].CreatedAt, userEntities[1].UpdatedAt))

	users, nextPageToken, err := adapter.GetPaginatedUsers(context.Background(), "alec", "", "", "", "", entities.PageInfo{
		NextPageToken: "",
		PageSize:      2,
	})
	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(nextPageToken).ToNot(BeEmpty())
	g.Expect(users).To(HaveLen(2))
}

func TestNewPostgresAdapter_GetPaginatedUsers_firstNameWithLastName(t *testing.T) {
	g := NewWithT(t)
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	adapter := adapters.NewPostgresAdapter(db)

	userEntities := []entities.User{
		{
			ID:        uuid.New(),
			FirstName: "alec",
			LastName:  "smith",
			Nickname:  "alecsmith",
			Password:  "somepasword",
			Email:     "alec@email.com",
			Country:   "UK",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        uuid.New(),
			FirstName: "alecsmith",
			LastName:  "smith",
			Nickname:  "alecsmith",
			Password:  "somepasword",
			Email:     "alec2@email.com",
			Country:   "UK",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	mock.ExpectQuery(`SELECT \* FROM platform_user WHERE 1=1 AND first_name ILIKE \$1 AND last_name ILIKE \$2 ORDER BY created_at, id LIMIT 10;`).
		WithArgs("%alec%", "%smith%").
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "first_name", "last_name", "nickname", "password", "email", "country", "created_at", "updated_at"}).
				AddRow(userEntities[0].ID, userEntities[0].FirstName, userEntities[0].LastName, userEntities[0].Nickname, userEntities[0].Password, userEntities[0].Email, userEntities[0].Country, userEntities[0].CreatedAt, userEntities[0].UpdatedAt).
				AddRow(userEntities[1].ID, userEntities[1].FirstName, userEntities[1].LastName, userEntities[1].Nickname, userEntities[1].Password, userEntities[1].Email, userEntities[1].Country, userEntities[1].CreatedAt, userEntities[1].UpdatedAt))

	users, nextPageToken, err := adapter.GetPaginatedUsers(context.Background(), "alec", "smith", "", "", "", entities.PageInfo{
		NextPageToken: "",
		PageSize:      10,
	})
	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(nextPageToken).To(Equal(""))
	g.Expect(users).To(HaveLen(2))
}

func TestNewPostgresAdapter_GetPaginatedUsers_allParams(t *testing.T) {
	g := NewWithT(t)
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	adapter := adapters.NewPostgresAdapter(db)

	userEntities := []entities.User{
		{
			ID:        uuid.New(),
			FirstName: "alec",
			LastName:  "smith",
			Nickname:  "alecsmith",
			Password:  "somepasword",
			Email:     "alec@email.com",
			Country:   "UK",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	mock.ExpectQuery(`SELECT \* FROM platform_user WHERE 1=1 AND first_name ILIKE \$1 AND last_name ILIKE \$2 AND nickname ILIKE \$3 AND email ILIKE \$4 AND country ILIKE \$5 ORDER BY created_at, id LIMIT 10;`).
		WithArgs("%alec%", "%smith%", "%alecsmith%", "%alec@email.com%", "%UK%").
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "first_name", "last_name", "nickname", "password", "email", "country", "created_at", "updated_at"}).
				AddRow(userEntities[0].ID, userEntities[0].FirstName, userEntities[0].LastName, userEntities[0].Nickname, userEntities[0].Password, userEntities[0].Email, userEntities[0].Country, userEntities[0].CreatedAt, userEntities[0].UpdatedAt))

	users, nextPageToken, err := adapter.GetPaginatedUsers(context.Background(), "alec", "smith", "alecsmith", "alec@email.com", "UK", entities.PageInfo{
		NextPageToken: "",
		PageSize:      10,
	})
	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(nextPageToken).To(Equal(""))
	g.Expect(users).To(HaveLen(1))
}

func TestNewPostgresAdapter_GetPaginatedUsers_WithNextPageToken(t *testing.T) {
	g := NewWithT(t)
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	adapter := adapters.NewPostgresAdapter(db)

	userEntities := []entities.User{
		{
			ID:        uuid.New(),
			FirstName: "alec",
			LastName:  "smith",
			Nickname:  "alecsmith",
			Password:  "somepasword",
			Email:     "alec@email.com",
			Country:   "UK",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        uuid.New(),
			FirstName: "alecsmith",
			LastName:  "smith",
			Nickname:  "alecsmith",
			Password:  "somepasword",
			Email:     "alec2@email.com",
			Country:   "UK",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	createdAt := time.Now().UTC()
	userID := uuid.New()
	token := fmt.Sprintf("%s|%s", createdAt.Format(time.RFC3339Nano), userID.String())
	nextPageToken := base64.URLEncoding.EncodeToString([]byte(token))

	mock.ExpectQuery(`SELECT \* FROM platform_user WHERE 1=1 AND first_name ILIKE \$1 AND \(created_at, id\) > \(\$2, \$3\) ORDER BY created_at, id LIMIT 10;`).
		WithArgs("%alec%", createdAt.UTC(), userID).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "first_name", "last_name", "nickname", "password", "email", "country", "created_at", "updated_at"}).
				AddRow(userEntities[0].ID, userEntities[0].FirstName, userEntities[0].LastName, userEntities[0].Nickname, userEntities[0].Password, userEntities[0].Email, userEntities[0].Country, userEntities[0].CreatedAt, userEntities[0].UpdatedAt).
				AddRow(userEntities[1].ID, userEntities[1].FirstName, userEntities[1].LastName, userEntities[1].Nickname, userEntities[1].Password, userEntities[1].Email, userEntities[1].Country, userEntities[1].CreatedAt, userEntities[1].UpdatedAt))

	users, nextPageToken, err := adapter.GetPaginatedUsers(context.Background(), "alec", "", "", "", "", entities.PageInfo{
		NextPageToken: nextPageToken,
		PageSize:      10,
	})
	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(nextPageToken).To(Equal(""))
	g.Expect(users).To(HaveLen(2))
}

func TestNewPostgresAdapter_GetPaginatedUsers_InvalidNextPageToken(t *testing.T) {
	g := NewWithT(t)
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	adapter := adapters.NewPostgresAdapter(db)

	users, nextPageToken, err := adapter.GetPaginatedUsers(context.Background(), "alec", "", "", "", "", entities.PageInfo{
		NextPageToken: "invalid-page-token",
		PageSize:      10,
	})
	g.Expect(err).To(MatchError("illegal base64 data at input byte 16"))
	g.Expect(nextPageToken).To(Equal(""))
	g.Expect(users).To(BeEmpty())
}

func TestNewPostgresAdapter_GetPaginatedUsers_QueryReturnsErr(t *testing.T) {
	g := NewWithT(t)
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	adapter := adapters.NewPostgresAdapter(db)

	createdAt := time.Now().UTC()
	userID := uuid.New()
	token := fmt.Sprintf("%s|%s", createdAt.Format(time.RFC3339Nano), userID.String())
	nextPageToken := base64.URLEncoding.EncodeToString([]byte(token))

	mock.ExpectQuery(`SELECT \* FROM platform_user WHERE 1=1 AND first_name ILIKE \$1 ORDER BY created_at, id LIMIT 10;`).
		WithArgs("%alec%").
		WillReturnError(errors.New("an error occurred"))

	users, nextPageToken, err := adapter.GetPaginatedUsers(context.Background(), "alec", "", "", "", "", entities.PageInfo{
		NextPageToken: "",
		PageSize:      10,
	})
	g.Expect(err).To(MatchError("an error occurred"))
	g.Expect(nextPageToken).To(Equal(""))
	g.Expect(users).To(BeEmpty())
}
