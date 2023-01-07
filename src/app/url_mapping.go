package app

import (
	controllers "github.com/Ayobami-00/booketo-mvc-go-postgres-gin/src/controllers/users"
	"github.com/gin-gonic/gin"
)

func mapUrls(router *gin.Engine) {

	router.POST("/users", controllers.CreateUser)
}
