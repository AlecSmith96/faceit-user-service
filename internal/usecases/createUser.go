package usecases

import (
	"context"
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

type CreateUserRequestBody struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Nickname  string `json:"nickname" binding:"required"`
	Password  string `json:"password" binding:"required"`
	Email     string `json:"email" binding:"required"`
	Country   string `json:"country" binding:"required"`
}

type CreateUserResponseBody struct {
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

func NewCreateUser(userCreator UserCreator) gin.HandlerFunc {
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
			slog.Error("creating user", "err", err)
			c.Status(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, user)
	}
}
