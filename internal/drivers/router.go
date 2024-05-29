package drivers

import (
	_ "github.com/AlecSmith96/faceit-user-service/docs"
	"github.com/AlecSmith96/faceit-user-service/internal/usecases"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

func NewRouter(
	userCreator usecases.UserCreator,
	userDeleter usecases.UserDeleter,
	userUpdater usecases.UserUpdater,
) *gin.Engine {
	r := gin.Default()

	// Swagger endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.POST("/user", usecases.NewCreateUser(userCreator))
	r.DELETE("/user/:userId", usecases.NewDeleteUser(userDeleter))
	r.PUT("/user/:userId", usecases.NewUpdateUser(userUpdater))

	return r
}
