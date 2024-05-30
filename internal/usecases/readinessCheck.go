package usecases

import (
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

type ReadinessChecker interface {
	CheckConnection() error
}

func NewReadinessCheck(readinessChecker ReadinessChecker) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := readinessChecker.CheckConnection()
		if err != nil {
			slog.Error("unable to establish connection with repository", "err", err)
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Status(http.StatusOK)
	}
}
