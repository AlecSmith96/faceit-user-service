package usecases

import (
	"context"
	"errors"
	"github.com/AlecSmith96/faceit-user-service/internal/entities"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
	"time"
)

//go:generate mockgen --build_flags=--mod=mod -destination=../../mocks/userUpdater.go  . "UserUpdater"
type UserUpdater interface {
	UpdateUser(ctx context.Context, userID uuid.UUID, firstName, lastName, nickname, password, email, country string) (*entities.User, error)
}

// UpdateUserRequestBody represents the request body for updating a user
// @Description Request body for updating a user
type UpdateUserRequestBody struct {
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

// UpdateUserResponseBody represents the response body for an updated user
// @Description Response body for an updated user
type UpdateUserResponseBody struct {
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

// NewUpdateUser updates a users information
// @Summary Update User
// @Description Updates user information for the provided userId
// @Tags users
// @Accept json
// @Produce json
// @Param userId path string true "User ID"
// @Param user body UpdateUserRequestBody true "Create User Request Body"
// @Success 200 {object} UpdateUserResponseBody
// @Failure 400
// @Failure 500
// @Router /users/{userId} [put]
func NewUpdateUser(userUpdater UserUpdater) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("userId")

		userIDUUID, err := uuid.Parse(userID)
		if err != nil {
			slog.Error("invalid userID", "err", err)
			c.Status(http.StatusBadRequest)
			return
		}

		var request CreateUserRequestBody
		err = c.ShouldBindJSON(&request)
		if err != nil {
			slog.Warn("unable to bind request", "err", err)
			c.Status(http.StatusBadRequest)
			return
		}

		user, err := userUpdater.UpdateUser(
			c.Request.Context(),
			userIDUUID,
			request.FirstName,
			request.LastName,
			request.Nickname,
			request.Password,
			request.Email,
			request.Country,
		)
		if err != nil {
			if errors.Is(err, entities.ErrUserNotFound) {
				slog.Warn("user not found", "err", err)
				c.Status(http.StatusBadRequest)
				return
			}

			slog.Error("updating user", "err", err)
			c.Status(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, UpdateUserResponseBody{
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
