package usecases

import (
	"context"
	"errors"
	_ "github.com/AlecSmith96/faceit-user-service/docs"
	"github.com/AlecSmith96/faceit-user-service/internal/entities"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"time"
)

//go:generate mockgen --build_flags=--mod=mod -destination=../../mocks/userCreator.go  . "UserCreator"
type UserCreator interface {
	CreateUser(ctx context.Context, firstName, lastName, nickname, password, email, country string) (*entities.User, error)
}

// CreateUserRequestBody represents the request body for creating a new user
// @Description Request body for creating a new user
type CreateUserRequestBody struct {
	// FirstName represents the user's first name
	FirstName string `json:"first_name" binding:"required"`
	// LastName represents the user's last name
	LastName string `json:"last_name" binding:"required"`
	// Nickname represents the user's nickname
	Nickname string `json:"nickname" binding:"required"`
	// Password represents the user's password
	Password string `json:"password" binding:"required"`
	// Email represents the user's email address
	Email string `json:"email" binding:"required"`
	// Country represents the user's country
	Country string `json:"country" binding:"required"`
}

// CreateUserResponseBody represents the response body for a created user
// @Description Response body for a created user
type CreateUserResponseBody struct {
	// ID represents the user's unique identifier
	ID string `json:"id"`
	// FirstName represents the user's first name
	FirstName string `json:"first_name"`
	// LastName represents the user's last name
	LastName string `json:"last_name"`
	// Nickname represents the user's nickname
	Nickname string `json:"nickname"`
	// Password represents the user's password
	Password string `json:"password"`
	// Email represents the user's email address
	Email string `json:"email"`
	// Country represents the user's country
	Country string `json:"country"`
	// CreatedAt represents the timestamp when the user was created
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt represents the timestamp when the user was last updated
	UpdatedAt time.Time `json:"updated_at"`
}

// NewCreateUser creates a new user
// @Summary Create a new user
// @Description Create a new user with the provided details
// @Tags users
// @Accept json
// @Produce json
// @Param user body CreateUserRequestBody true "Create User Request Body"
// @Success 200 {object} CreateUserResponseBody
// @Failure 400
// @Failure 500
// @Router /user [post]
func NewCreateUser(userCreator UserCreator, changelogWriter ChangelogWriter) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request CreateUserRequestBody
		err := c.ShouldBindJSON(&request)
		if err != nil {
			slog.Warn("unable to bind request", "err", err)
			c.Status(http.StatusBadRequest)
			return
		}

		user, err := userCreator.CreateUser(
			c.Request.Context(),
			request.FirstName,
			request.LastName,
			request.Nickname,
			request.Password,
			request.Email,
			request.Country,
		)
		if err != nil {
			if errors.Is(err, entities.ErrEmailAlreadyUsed) {
				slog.Warn("email already registered to a uer", "err", err)
				c.Status(http.StatusBadRequest)
				return
			}

			slog.Error("creating user", "err", err)
			c.Status(http.StatusInternalServerError)
			return
		}

		entry := entities.ChangelogEntry{
			UserID:     user.ID,
			CreatedAt:  user.CreatedAt,
			ChangeType: "POST",
		}
		err = changelogWriter.PublishChangelogEntry(entry)
		if err != nil {
			// deliberately not returning error here as request didn't fail
			slog.Error("publishing changelog event", "err", err, "changelogEntry", entry)
		}

		c.JSON(http.StatusOK, CreateUserResponseBody{
			ID:        user.ID.String(),
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Nickname:  user.Nickname,
			Password:  user.Password,
			Email:     user.Email,
			Country:   user.Country,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		})
	}
}
