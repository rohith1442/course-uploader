package routers

import (
	transcode "workspace/goproject/controllers"

	"github.com/gin-gonic/gin"
)

func TranscodeRoutes(incomingRoutes *gin.Engine) {

	v1 := incomingRoutes.Group("/v1")
	v1.POST("/fileupload", transcode.HandleTranscode())
	v1.GET("/baremetal/:video", transcode.HandleBaremetal())

}
