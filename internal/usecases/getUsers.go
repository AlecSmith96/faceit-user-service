package usecases

import "github.com/gin-gonic/gin"

type UserGetter interface {
}

func NewGetUsers(userCreator UserCreator) gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
