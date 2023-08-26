package main

import (
	routes "workspace/goproject/routers"

	"workspace/goproject/utils"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var AccessKeyID string
var SecretAccessKey string
var MyRegion string
var MyBucket string
var filepath string

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
