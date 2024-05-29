package usecases

import (
	"context"
	"errors"
	"github.com/AlecSmith96/faceit-user-service/internal/entities"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
)

type UserDeleter interface {
	DeleteUser(ctx context.Context, userID uuid.UUID) error
}

func NewDeleteUser(userDeleter UserDeleter) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("userId")

		userIDUUID, err := uuid.Parse(userID)
		if err != nil {
			slog.Error("invalid userID", "err", err)
			c.Status(http.StatusBadRequest)
			return
		}

		err = userDeleter.DeleteUser(c.Request.Context(), userIDUUID)
		if err != nil {
			if errors.Is(err, entities.ErrUserNotFound) {
				slog.Warn("user not found", "err", err)
				c.Status(http.StatusBadRequest)
				return
			}

			slog.Error("deleting user", "err", err)
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Status(http.StatusOK)
	}
}
