package main

import (
	routes "workspace/course-uploader/routers"

	"workspace/course-uploader/utils"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	port := utils.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	router := gin.New()
	router.Use(gin.Logger())
	routes.TranscodeRoutes(router)

	godotenv.Load(".env")
	router.Run(":" + port)
}
