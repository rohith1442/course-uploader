package routers

import (
	"workspace/goproject/controllers"
	"workspace/goproject/service"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func TranscodeRoutes(incomingRoutes *gin.Engine) {

	// Init Handlers
	videoService := &service.VideoService{
		Logger: logrus.New(),
	}
	videoController := controllers.InitVideoHandler("vidhyatech-course", "s3://vidhyatech-course/", videoService)

	// Video Handler Routes
	v1 := incomingRoutes.Group("/v1")
	v1.POST("/file/upload", videoController.HandleTranscode())
	v1.GET("/baremetal/:video", videoController.HandleBaremetal())

}
