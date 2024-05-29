package usecases

import "github.com/gin-gonic/gin"

type UserGetter interface {
}

func NewGetUsers(userCreator UserCreator, changelogWriter ChangelogWriter) gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
