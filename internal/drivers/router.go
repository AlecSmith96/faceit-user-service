package drivers

import (
	_ "github.com/AlecSmith96/faceit-user-service/docs"
	"github.com/AlecSmith96/faceit-user-service/internal/usecases"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

func NewRouter(
	changelogWriter usecases.ChangelogWriter,
	userGetter usecases.UserGetter,
	userCreator usecases.UserCreator,
	userDeleter usecases.UserDeleter,
	userUpdater usecases.UserUpdater,
	readinessChecker usecases.ReadinessChecker,
) *gin.Engine {
	r := gin.Default()

	// docs endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/users", usecases.NewGetUsers(userGetter, changelogWriter))
	r.POST("/user", usecases.NewCreateUser(userCreator, changelogWriter))
	r.DELETE("/user/:userId", usecases.NewDeleteUser(userDeleter, changelogWriter))
	r.PUT("/user/:userId", usecases.NewUpdateUser(userUpdater, changelogWriter))

	// health check
	r.GET("/health/readiness", usecases.NewReadinessCheck(readinessChecker))

	return r
}
