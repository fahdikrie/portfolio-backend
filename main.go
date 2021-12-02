package main

import (
	"os"
	"portfolio-backend/routes"

	"github.com/gin-gonic/gin"
)

func main() {
    port := os.Getenv("PORT")

    if port == "" {
        port = "8000"
    }

    router := gin.New()
    router.Use(gin.Logger())
    routes.UserRoutes(router)
    router.Run(":" + port)
}
