package drivers

import (
	"github.com/AlecSmith96/faceit-user-service/internal/usecases"
	"github.com/gin-gonic/gin"
)

func NewRouter(
	userCreator usecases.UserCreator,
	userDeleter usecases.UserDeleter,
	userUpdater usecases.UserUpdater,
) *gin.Engine {
	r := gin.Default()

	r.POST("/user", usecases.NewCreateUser(userCreator))
	r.DELETE("/user/:userId", usecases.NewDeleteUser(userDeleter))
	r.PUT("/user/:userId", usecases.NewUpdateUser(userUpdater))

	return r
}
