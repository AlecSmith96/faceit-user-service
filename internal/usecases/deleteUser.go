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

//go:generate mockgen --build_flags=--mod=mod -destination=../../mocks/userDeleter.go  . "UserDeleter"
type UserDeleter interface {
	DeleteUser(ctx context.Context, userID uuid.UUID) error
}

// NewDeleteUser deletes a user
// @Summary Delete user
// @Description Deletes a user from the system
// @Tags users
// @Accept json
// @Produce json
// @Param userId path string true "User ID"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /users/{userId} [delete]
func NewDeleteUser(userDeleter UserDeleter, changelogWriter ChangelogWriter) gin.HandlerFunc {
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

		entry := entities.ChangelogEntry{
			UserID:     userIDUUID,
			CreatedAt:  time.Now(),
			ChangeType: "DELETE",
		}
		err = changelogWriter.PublishChangelogEntry(entry)
		if err != nil {
			// deliberately not returning error here as request didn't fail
			slog.Error("publishing changelog event", "err", err, "changelogEntry", entry)
		}

		c.Status(http.StatusOK)
	}
}
