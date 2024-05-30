package adapters_test

import (
	"context"
	"errors"
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
