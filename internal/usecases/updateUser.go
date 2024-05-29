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

type UserUpdater interface {
	UpdateUser(ctx context.Context, userID uuid.UUID, firstName, lastName, nickname, password, email, country string) (*entities.User, error)
}

type UpdateUserRequestBody struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Nickname  string `json:"nickname" binding:"required"`
	Password  string `json:"password" binding:"required"`
	Email     string `json:"email" binding:"required"`
	Country   string `json:"country" binding:"required"`
}

type UpdateUserResponseBody struct {
	ID        string    `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Nickname  string    `json:"nickname"`
	Password  string    `json:"password"`
	Email     string    `json:"email"`
	Country   string    `json:"country"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

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
