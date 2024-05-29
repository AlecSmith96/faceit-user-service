package drivers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Router struct {
}

func NewRouter() (*Router, error) {

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	return r, nil
}
