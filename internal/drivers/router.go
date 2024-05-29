package drivers

import (
	"github.com/AlecSmith96/faceit-user-service/internal/usecases"
	"github.com/gin-gonic/gin"
	"net/http"
)

func NewRouter(
	userCreator usecases.UserCreator,
) *gin.Engine {
	r := gin.Default()

	r.POST("/user", usecases.NewCreateUser(userCreator))
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	return r
}
