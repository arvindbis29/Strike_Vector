package handler

import (
	insightsGenerateController "voice-hack-backend/modules/tripPlanner/controller/insightsGenerateController"

	"github.com/gin-gonic/gin"
)

func RouteRequests(ginServer *gin.Engine) {
	apiGroup := ginServer.Group("insights")
	apiGroup.POST("/generate", insightsGenerateController.FindDestination)
	apiGroup.POST("/generate/", insightsGenerateController.FindDestination)
	apiGroup.GET("/generate", insightsGenerateController.FindDestination)
	apiGroup.GET("/generate/", insightsGenerateController.FindDestination)
	apiGroup.POST("/final", insightsGenerateController.FinalSummaryGenerate)
	apiGroup.POST("/final/", insightsGenerateController.FinalSummaryGenerate)
	apiGroup.GET("/final", insightsGenerateController.FinalSummaryGenerate)
	apiGroup.GET("/final/", insightsGenerateController.FinalSummaryGenerate)

}
