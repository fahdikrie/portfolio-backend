package routes

import (
	"portfolio-backend/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/user/login", controllers.Login())
	incomingRoutes.POST("/user/__signup", controllers.SignUp())
}

